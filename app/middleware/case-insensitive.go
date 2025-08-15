package middleware

import (
	"net/http"
	"strings"
)

func CaseInsensitiveMux(mux http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.URL.Path = strings.ToLower(r.URL.Path)
		mux.ServeHTTP(w, r)
	})
}
