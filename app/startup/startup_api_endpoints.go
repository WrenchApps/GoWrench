package startup

import (
	"context"
	"net/http"
	"os"
	"strings"

	handler "wrench/app/handlers"
	"wrench/app/manifest/application_settings"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func LoadApiEndpoint(ctx context.Context) http.Handler {
	app := application_settings.ApplicationSettingsStatic

	if app.Api == nil || app.Api.Endpoints == nil {
		return nil
	}

	endpoints := app.Api.Endpoints
	muxRoute := mux.NewRouter()
	initialPage := new(InitialPage)
	initialPage.Append("<h2>Service: " + app.Service.Name + " version: " + app.Service.Version + "</h2>")
	initialPage.Append("<h2>Endpoints</h2>")

	for _, endpoint := range endpoints {
		var delegate = new(handler.RequestDelegate)
		delegate.SetEndpoint(&endpoint)
		delegate.Otel = app.Service.Otel

		if !endpoint.IsProxy {
			method := strings.ToUpper(string(endpoint.Method))
			route := endpoint.Route
			muxRoute.HandleFunc(route, delegate.HttpHandler).Methods(method)
			initialPage.Append("Route: <i>" + route + "</i> Method: <i>" + method + "</i> <b>Not is proxy</b>")
		} else {
			initialPage.Append("Route: <i>" + endpoint.Route + "</i> <b> IS PROXY</b>")
			if endpoint.Route == "/" {
				endpoint.Route = ""
			}
			muxRoute.HandleFunc(endpoint.Route+"/{path:.*}", delegate.HttpHandler)
		}
	}

	initialPage.Append("</br></br>")
	initialPage.Append("<h2>Envs</h2>")
	for _, env := range os.Environ() {
		envSplitted := strings.Split(env, "=")
		envName := envSplitted[0]
		initialPage.Append("Env: <i>" + envName + "</i>")
	}
	muxRoute.HandleFunc("/", initialPage.WriteInitialPage).Methods("GET")
	muxRoute.HandleFunc("/hc", initialPage.HealthCheckEndpoint).Methods("GET")

	if app.Api.Cors != nil {
		if len(app.Api.Cors.Origins) == 0 {
			app.Api.Cors.Origins = []string{"*"}
		}

		if len(app.Api.Cors.Methods) == 0 {
			app.Api.Cors.Methods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
		} else {
			for i, item := range app.Api.Cors.Methods {
				app.Api.Cors.Methods[i] = strings.ToUpper(item)
			}
		}

		if len(app.Api.Cors.Headers) == 0 {
			app.Api.Cors.Headers = []string{"Accept", "Accept-Language", "Content-Language", "Content-Type", "Authorization", "X-Requested-With", "X-Custom-Header"}
		}

		return handlers.CORS(
			handlers.AllowedOrigins(app.Api.Cors.Origins),
			handlers.AllowedMethods(app.Api.Cors.Methods),
			handlers.AllowedHeaders(app.Api.Cors.Headers),
		)(muxRoute)
	}

	return muxRoute
}
