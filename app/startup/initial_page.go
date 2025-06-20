package startup

import (
	"encoding/json"
	"net/http"
	"wrench/app/cross_validation"
	"wrench/app/manifest/application_settings"
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
	app := application_settings.ApplicationSettingsStatic
	result := app.Valid()
	result.Append(cross_validation.Valid())

	w.Header().Set("Content-Type", "application/json")

	if result.IsSuccess() {
		body := make(map[string]interface{})
		body["status"] = "healthly"
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(body)
	} else {
		body := make(map[string]interface{})
		body["status"] = "unhealthly"
		body["erros"] = result.GetErrors()
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(body)
	}
}
