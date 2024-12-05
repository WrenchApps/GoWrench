package handlers

import (
	contexts "wrench/app/contexts"
)

type HttpLastHandler struct {
	Next Handler
}

func (httpLast *HttpLastHandler) Do(wrenchContext *contexts.WrenchContext, bodyContext *contexts.BodyContext) {
	var w = *wrenchContext.ResponseWriter

	header := w.Header()
	header.Set("Content-Type", bodyContext.ContentType)

	if bodyContext.Headers != nil {
		for key, value := range bodyContext.Headers {
			header.Set(key, value)
		}
	}

	w.WriteHeader(bodyContext.HttpStatusCode)
	w.Write([]byte(bodyContext.Body))
	if httpLast.Next != nil {
		httpLast.Next.Do(wrenchContext, bodyContext)
	}
}

func (httpLast *HttpLastHandler) SetNext(next Handler) {
	httpLast.Next = next
}
