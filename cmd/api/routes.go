package main

import "net/http"

func (app *application) routes() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /", app.homeHandler)
	mux.HandleFunc("GET /favicon.ico", app.faviconHandler)
	mux.HandleFunc("GET /v1/health", app.healthCheckHandler)
	mux.HandleFunc("GET /v1/about", app.aboutUsHandler)
	return mux
}
