package startup

import (
	"net/http"
	handlers "wrench/app/handlers"
	settings "wrench/app/manifest/application_settings"
)

func LoadApplicationSettings(settings *settings.ApplicationSettings) http.Handler {
	var chain = handlers.ChainStatic.GetStatic()
	chain.BuildChain(settings)
	return LoadApiEndpoint()
}
