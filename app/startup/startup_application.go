package startup

import (
	handlers "wrench/app/handlers"
	settings "wrench/app/manifest/application_settings"

	"github.com/gorilla/mux"
)

func LoadApplicationSettings(settings *settings.ApplicationSettings) *mux.Router {
	var chain = handlers.ChainStatic.GetStatic()
	chain.BuildChain(settings)

	var endpoints = settings.Api.Endpoints
	if len(endpoints) > 0 {
		return LoadApiEndpoint(endpoints)
	} else {
		return nil
	}
}
