package startup

import (
	"strings"
	handler "wrench/app/handlers"
	settings "wrench/app/manifest/api_settings"

	"github.com/gorilla/mux"
)

func LoadApiEndpoint(endpoints []settings.EndpointSettings) *mux.Router {
	r := mux.NewRouter()
	initialPage := new(InitialPage)

	for _, endpoint := range endpoints {
		method := strings.ToUpper(string(endpoint.Method))
		route := endpoint.Route

		var delegate = new(handler.RequestDelegate)
		delegate.SetEndpoint(&endpoint)
		r.HandleFunc(route, delegate.FirstHttp).Methods(method)
		initialPage.Append("Route: <i>" + route + "</i> Method: <i>" + method + "</i>")
	}

	r.HandleFunc("/", initialPage.WriteInitialPage).Methods("GET")
	return r
}
