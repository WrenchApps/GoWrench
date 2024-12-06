package contexts

import (
	"net/http"
)

type WrenchContext struct {
	ResponseWriter *http.ResponseWriter
	Request        *http.Request
	HasError       bool
}

func (wrenchContext *WrenchContext) SetHasError() {
	wrenchContext.HasError = true
}
