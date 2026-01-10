package main

import "net/http"

func (app *application) routes() http.Handler {
	mux := http.NewServeMux()

	// Static file serving with caching headers for better performance
	mux.Handle("/static/", cacheControl(http.StripPrefix("/static/", http.FileServer(http.Dir("./static")))))

	// Register routes
	mux.HandleFunc("GET /", app.homeHandler)
	mux.HandleFunc("GET /favicon.ico", app.faviconHandler)
	mux.HandleFunc("GET /v1/health", app.healthCheckHandler)
	mux.HandleFunc("GET /v1/about", app.aboutUsHandler)

	// Apply middleware chain (order matters: last middleware wraps first)
	// Security headers (outermost - applies to all responses)
	// -> Request logging
	// -> Session management
	// -> CSRF protection
	// -> CSRF token generation
	// -> Routes
	return secureHeaders(
		app.logRequest(
			app.sessionMiddleware(
				app.csrfProtect(
					app.csrfTokenGenerator(mux),
				),
			),
		),
	)
}
