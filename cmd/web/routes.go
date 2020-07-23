package main

import (
	"net/http"

	"github.com/justinas/alice"
)

func (app *application) routes(static *string) http.Handler {
	standardMiddleware := alice.New(app.recoverPanic, app.logRequest, secureHeaders)
	mux := http.NewServeMux()
	mux.HandleFunc("/snippet", app.showSnippet)
	mux.HandleFunc("/snippet/create", app.createSnippet)
	mux.HandleFunc("/", app.home)

	fileServer := http.FileServer(http.Dir(*static))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	return standardMiddleware.Then(mux)
}
