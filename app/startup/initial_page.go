package startup

import "net/http"

type InitialPage struct {
	Html string
}

func (page *InitialPage) Append(text string) {
	page.Html += page.Html + "<p>" + text + "</p>"
}

func (page *InitialPage) WriteInitialPage(w http.ResponseWriter, r *http.Request) {
	htmlFirst := "<!DOCTYPE html><html><head><title>Initial Page</title></head><body>" + page.Html + "</body></html>"
	w.Write([]byte(htmlFirst))
}
