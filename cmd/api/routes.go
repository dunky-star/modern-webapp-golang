package main

import (
	"net/http"

	"github.com/dunky-star/modern-webapp-golang/internal/handlers"
)

func routes() http.Handler {
	mux := http.NewServeMux()

	// Static file serving with caching headers (method-specific to avoid Go 1.22+ conflicts)
	staticHandler := cacheControl(http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))
	mux.HandleFunc("GET /static/", staticHandler.ServeHTTP)
	mux.HandleFunc("HEAD /static/", staticHandler.ServeHTTP)

	// Register application routes
	mux.HandleFunc("GET /", handlers.Repo.HomeHandler)
	mux.HandleFunc("GET /health", handlers.Repo.HealthCheckHandler)
	mux.HandleFunc("GET /about", handlers.Repo.AboutUsHandler)
	mux.HandleFunc("GET /contact", handlers.Repo.ContactHandler)
	mux.HandleFunc("GET /user/login", handlers.Repo.ShowLoginHandler)
	mux.HandleFunc("GET /search-availability", handlers.Repo.SearchAvailabilityHandler)
	mux.HandleFunc("POST /search-availability", handlers.Repo.PostAvailabilityHandler)
	mux.HandleFunc("POST /search-availability-json", handlers.Repo.AvialabilityJSONHandler)
	mux.HandleFunc("GET /choose-room/{id}", handlers.Repo.ChooseRoomHandler)
	mux.HandleFunc("GET /book-room", handlers.Repo.BookRoomHandler)
	mux.HandleFunc("GET /generals-quarters", handlers.Repo.GeneralsQuartersHandler)
	mux.HandleFunc("GET /majors-suite", handlers.Repo.MajorsSuiteHandler)
	mux.HandleFunc("GET /make-reservation", handlers.Repo.MakeReservationHandler)
	mux.HandleFunc("POST /make-reservation", handlers.Repo.PostReservationHandler)
	mux.HandleFunc("GET /reservation-summary", handlers.Repo.ReservationSummary)

	// Apply middleware chain (order matters: last middleware wraps first)
	// Security headers (outermost - applies to all responses)
	// -> Request logging
	// -> HTML cache control (for dynamic pages)
	// -> Session management
	// -> CSRF protection
	// -> CSRF token generation
	// -> Routes
	return secureHeaders(
		logRequest(
			htmlCacheControl(
				sessionMiddleware(
					csrfProtect(
						csrfTokenGenerator(mux),
					),
				),
			),
		),
	)
}
