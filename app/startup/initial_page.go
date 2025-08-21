package startup

import (
	"encoding/json"
	"net/http"
	"wrench/app"
	"wrench/app/cross_validation"
	"wrench/app/manifest/application_settings"
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

func (page *InitialPage) HealthCheckEndpoint(w http.ResponseWriter, r *http.Request) {
	application := application_settings.ApplicationSettingsStatic
	result := application.Valid()
	result.Append(cross_validation.Valid())

	w.Header().Set("Content-Type", "application/json")

	var errors []error

	if token_credentials.CredentialErrors != nil {
		errors = append(errors, token_credentials.CredentialErrors...)
	}

	if result.IsSuccess() && len(errors) == 0 {
		body := make(map[string]interface{})
		body["status"] = "healthly"
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(body)
	} else {

		result.AddErrors(errors)

		for _, err := range result.GetErrors() {
			app.LogError2(err, nil)
		}

		body := make(map[string]interface{})
		body["status"] = "unhealthly"
		body["erros"] = result.GetErrors()
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(body)
	}
}
