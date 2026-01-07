package client_http

import (
	"bytes"
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"wrench/app/startup/tls_load"

	"go.opentelemetry.io/otel/trace"
)

var httpClient *http.Client = new(http.Client)
var httpClientInsecure *http.Client

var clients map[string]*http.Client = map[string]*http.Client{}

func GetHttpClient(request *HttpClientRequestData) (*http.Client, error) {
	clientKey := fmt.Sprintf("%s-%s", strconv.FormatBool(request.Insecure), request.TlsId)

	client := clients[clientKey]
	if client != nil {
		return client, nil
	}

	var tlsConfig *tls.Config

	if len(request.TlsId) > 0 {
		tlsConfigById, tlsErr := tls_load.GetTlsConfigById(request.TlsId)

		if tlsErr != nil {
			return nil, tlsErr
		}

		tlsConfig = tlsConfigById
	}

	if tlsConfig == nil {
		tlsConfig = &tls.Config{}
	}

	tlsConfig.InsecureSkipVerify = request.Insecure

	client = &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: tlsConfig,
		},
	}

	clients[clientKey] = client
	return client, nil
}

type HttpClientRequestData struct {
	Url      string
	Method   string
	Body     []byte
	Headers  map[string]string
	Insecure bool
	TlsId    string
}

type HttpClientResponseData struct {
	Body               []byte
	Headers            map[string]string
	StatusCode         int
	HttpClientResponse *http.Response
}

func (httpClientRequestData *HttpClientRequestData) SetDefaultHeader(ctx context.Context) {
	spanContext := trace.SpanContextFromContext(ctx)
	traceId := spanContext.TraceID().String()
	traceparent := fmt.Sprintf("00-%s-%s-%s", traceId, spanContext.SpanID(), "01")
	httpClientRequestData.SetHeader("traceparent", traceparent)

	if _, err := httpClientRequestData.GetHeaderValue("Content-Type"); err != nil {
		httpClientRequestData.SetHeader("Content-Type", "application/json")
	}
}

func (httpClientRequestData *HttpClientRequestData) SetHeaders(headers map[string]interface{}) {
	if headers != nil {
		if httpClientRequestData.Headers == nil {
			httpClientRequestData.Headers = make(map[string]string)
		}

		for key, value := range headers {
			httpClientRequestData.Headers[key] = fmt.Sprint(value)
		}
	}
}

func (httpResponse *HttpClientResponseData) StatusCodeSuccess() bool {
	return httpResponse.StatusCode <= 399
}

func (httpClientRequestData *HttpClientRequestData) SetHeader(key string, value string) {
	if len(key) > 0 {
		if httpClientRequestData.Headers == nil {
			httpClientRequestData.Headers = make(map[string]string)
		}

		httpClientRequestData.Headers[key] = value
	}
}

func (httpClientRequestData *HttpClientRequestData) GetHeaderValue(key string) (string, error) {
	if len(key) > 0 && httpClientRequestData.Headers != nil {
		value := httpClientRequestData.Headers[key]
		if value != "" {
			return value, nil
		}
	}
	return "", fmt.Errorf("header key '%s' not found ", key)
}

func HttpClientDo(ctx context.Context, request *HttpClientRequestData) (*HttpClientResponseData, error) {
	client, err := GetHttpClient(request)

	if err != nil {
		fmt.Println("Error creating client:", err)
		return nil, err
	}

	method := strings.ToUpper(request.Method)
	var body io.Reader = nil

	if request.Body != nil {
		body = bytes.NewBuffer(request.Body)
	}

	req, err := http.NewRequestWithContext(ctx, method, request.Url, body)

	if err != nil {
		fmt.Println("Error creating request:", err)
		return nil, err
	}

	header := req.Header
	if request.Headers != nil {
		for key, value := range request.Headers {
			header.Set(key, value)
		}
	}

	resp, err := client.Do(req)

	if err != nil {
		fmt.Println("Error making request:", err)
		return nil, err
	}

	defer resp.Body.Close()
	respBody, err := io.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}

	response := new(HttpClientResponseData)
	response.Body = respBody
	response.StatusCode = resp.StatusCode
	response.HttpClientResponse = resp
	return response, nil
}
