package main

import "net/http"

func (app *application) routes() http.Handler {
	mux := http.NewServeMux()

	// Register routes
	mux.HandleFunc("GET /", app.homeHandler)
	mux.HandleFunc("GET /favicon.ico", app.faviconHandler)
	mux.HandleFunc("GET /v1/health", app.healthCheckHandler)
	mux.HandleFunc("GET /v1/about", app.aboutUsHandler)

	// Apply middleware chain (order matters: last middleware wraps first)
	// Health check typically doesn't need request logging in production
	return app.logRequest(mux)
}
