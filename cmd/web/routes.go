package main

import "net/http"

func (app *application) routes(static *string) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/snippet", app.showSnippet)
	mux.HandleFunc("/snippet/create", app.createSnippet)
	mux.HandleFunc("/", app.home)

	fileServer := http.FileServer(http.Dir(*static))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	return app.logRequest(secureHeaders(mux))
}
