package client_http

import (
	"bytes"
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"strings"
)

var httpClient *http.Client = new(http.Client)
var httpClientInsecure *http.Client

func GetHttpClientStatic() *http.Client {
	return httpClient
}

func GetHttpClientInsecureStatic() *http.Client {

	if httpClientInsecure == nil {
		httpClientInsecure = &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: true,
				},
			},
		}
	}

	return httpClientInsecure
}

type HttpClientRequestData struct {
	Url      string
	Method   string
	Body     []byte
	Headers  map[string]string
	Insecure bool
}

type HttpClientResponseData struct {
	Body               []byte
	Headers            map[string]string
	StatusCode         int
	HttpClientResponse *http.Response
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
	if httpResponse.StatusCode <= 399 {
		return true
	}

	return false
}

func (httpClientRequestData *HttpClientRequestData) SetHeader(key string, value string) {
	if len(key) > 0 {
		if httpClientRequestData.Headers == nil {
			httpClientRequestData.Headers = make(map[string]string)
		}

		httpClientRequestData.Headers[key] = value
	}
}

func HttpClientDo(ctx context.Context, request *HttpClientRequestData) (*HttpClientResponseData, error) {
	var client *http.Client

	if request.Insecure == false {
		client = GetHttpClientStatic()
	} else {
		client = GetHttpClientInsecureStatic()
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
