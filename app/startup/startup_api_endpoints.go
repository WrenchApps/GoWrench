package startup

import (
	"context"
	"net/http"
	"os"
	"strings"

	"wrench/app/auth"
	handler "wrench/app/handlers"
	"wrench/app/manifest/api_settings"
	"wrench/app/manifest/application_settings"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func LoadApiEndpoint(ctx context.Context) http.Handler {
	app := application_settings.ApplicationSettingsStatic

	if app.Api == nil || app.Api.Endpoints == nil {
		return nil
	}

	hasAuthorization := app.Api.HasAuthorization()
	endpoints := app.Api.Endpoints
	muxRoute := mux.NewRouter()
	initialPage := new(InitialPage)
	initialPage.Append("<h2>Service: " + app.Service.Name + " version: " + app.Service.Version + "</h2>")
	initialPage.Append("<h2>Endpoints</h2>")

	for _, endpoint := range endpoints {
		shouldConfigureAuthorization := endpoint.ShouldConfigureAuthorization(hasAuthorization)
		var delegate = new(handler.RequestDelegate)
		delegate.SetEndpoint(&endpoint)

		if !endpoint.IsProxy {
			method := strings.ToUpper(string(endpoint.Method))
			route := endpoint.Route

			if shouldConfigureAuthorization {
				muxRoute.Handle(route, authMiddleware(ctx, app.Api.Authorization, endpoint, http.HandlerFunc(delegate.HttpHandler))).Methods(method)
			} else {
				muxRoute.HandleFunc(route, delegate.HttpHandler).Methods(method)
			}
			initialPage.Append("Route: <i>" + route + "</i> Method: <i>" + method + "</i> <b>Not is proxy</b>")
		} else {
			initialPage.Append("Route: <i>" + endpoint.Route + "</i> <b> IS PROXY</b>")
			if endpoint.Route == "/" {
				endpoint.Route = ""
			}

			if shouldConfigureAuthorization {
				muxRoute.Handle(endpoint.Route+"/{path:.*}", authMiddleware(ctx, app.Api.Authorization, endpoint, http.HandlerFunc(delegate.HttpHandler)))
			} else {
				muxRoute.HandleFunc(endpoint.Route+"/{path:.*}", delegate.HttpHandler)
			}
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

func authMiddleware(ctx context.Context, authorizationSettings *api_settings.AuthorizationSettings, endpoint api_settings.EndpointSettings, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if authorizationSettings.Type == api_settings.JWKSAuthorizationType {
			tokenString := r.Header.Get("Authorization")
			if len(tokenString) == 0 {
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte("Unauthorized"))
				return
			}
			tokenString = strings.Replace(tokenString, "Bearer ", "", 1)

			tokenIsValid := auth.JwksValidationAuthentication(ctx, tokenString, authorizationSettings)
			if tokenIsValid {
				tokenIsAuthorized := auth.JwksValidationAuthorization(tokenString, endpoint.Roles, endpoint.Scopes, endpoint.Claims)
				if !tokenIsAuthorized {
					w.WriteHeader(http.StatusForbidden)
					w.Write([]byte("Forbidden"))
					return
				}
			} else {
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte("Unauthorized"))
				return
			}

		}

		if authorizationSettings.Type == api_settings.HMACAuthorizationType {
			isHMACValid := auth.HMACValidate(ctx, authorizationSettings)
			if !isHMACValid {
				w.WriteHeader(http.StatusForbidden)
				w.Write([]byte("Forbidden"))
				return
			}
		}

		next.ServeHTTP(w, r)
	})
}
