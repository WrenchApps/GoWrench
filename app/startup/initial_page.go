package startup

import (
	"encoding/json"
	"net/http"
	"wrench/app"
	"wrench/app/cross_validation"
	"wrench/app/manifest/application_settings"
	"wrench/app/startup/connections"
	keys_load "wrench/app/startup/keys"
	"wrench/app/startup/token_credentials"
)

type InitialPage struct {
	Html string
}

func (page *InitialPage) Append(text string) {
	html := page.Html + "<p>" + text + "</p>"
	page.Html = html
}

func (page *InitialPage) WriteInitialPage(w http.ResponseWriter, r *http.Request) {
	htmlFirst := "<!DOCTYPE html><html><head><title>Initial Page</title></head><body>" + page.Html + "</body></html>"
	w.Write([]byte(htmlFirst))
}

var bodyHcResult map[string]interface{}
var statusCode int

func (page *InitialPage) HealthCheckEndpoint(w http.ResponseWriter, r *http.Request) {

	var errors []error

	if token_credentials.CredentialErrors != nil {
		errors = append(errors, token_credentials.CredentialErrors...)
	}

	if connections.ErrorLoadConnections != nil {
		errors = append(errors, connections.ErrorLoadConnections...)
	}

	if keys_load.ErrorLoadKeys != nil {
		errors = append(errors, keys_load.ErrorLoadKeys...)
	}

	if bodyHcResult == nil || len(errors) > 0 {
		application := application_settings.ApplicationSettingsStatic
		result := application.Valid()
		result.Append(cross_validation.Valid())

		w.Header().Set("Content-Type", "application/json")

		if result.IsSuccess() && len(errors) == 0 {
			statusCode = http.StatusOK
			bodyHcResult = make(map[string]interface{})
			bodyHcResult["status"] = "healthy"

		} else {
			statusCode = http.StatusInternalServerError
			result.AddErrors(errors)

			for _, err := range result.GetErrors() {
				app.LogError2(err, nil)
			}

			bodyHcResult = make(map[string]interface{})
			bodyHcResult["status"] = "unhealthy"
			bodyHcResult["errors"] = result.GetErrors()
		}
	}

	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(bodyHcResult)
}
