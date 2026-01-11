package main

import "net/http"

func (app *application) routes() http.Handler {
	mux := http.NewServeMux()

	// Static file serving with caching headers (method-specific to avoid Go 1.22+ conflicts)
	staticHandler := cacheControl(http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))
	mux.HandleFunc("GET /static/", staticHandler.ServeHTTP)
	mux.HandleFunc("HEAD /static/", staticHandler.ServeHTTP)

	// Register application routes
	mux.HandleFunc("GET /", app.homeHandler)
	mux.HandleFunc("GET /health", app.healthCheckHandler)
	mux.HandleFunc("GET /about", app.aboutUsHandler)
	mux.HandleFunc("GET /contact", app.contactHandler)
	mux.HandleFunc("GET /search-availability", app.searchAvailabilityHandler)
	mux.HandleFunc("POST /search-availability", app.postAvailabilityHandler)
	mux.HandleFunc("GET /search-availability-json", app.avialabilityJSONHandler)
	mux.HandleFunc("GET /generals-quarters", app.generalsQuartersHandler)
	mux.HandleFunc("GET /majors-suite", app.majorsSuiteHandler)
	mux.HandleFunc("GET /make-reservation", app.makeReservationHandler)

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
