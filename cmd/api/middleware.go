package main

import (
	"context"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/dunky-star/modern-webapp-golang/pkg/csrf"
	"github.com/dunky-star/modern-webapp-golang/pkg/helpers"
	"github.com/dunky-star/modern-webapp-golang/pkg/logging"
)

var (
	requestLogger     *log.Logger
	requestLogWriter  *logging.RotatingLogWriter
	requestLoggerOnce sync.Once
)

// initRequestLogger initializes the rotating log writer for HTTP request logs
func initRequestLogger(alsoWriteToConsole bool) error {
	var initErr error
	requestLoggerOnce.Do(func() {
		rotatingWriter, err := logging.NewRotatingLogWriter(alsoWriteToConsole)
		if err != nil {
			initErr = err
			return
		}
		requestLogWriter = rotatingWriter
		requestLogger = log.New(rotatingWriter, "", log.Ldate|log.Ltime)
	})
	return initErr
}

// closeRequestLogger closes the rotating log writer (call on application shutdown)
func closeRequestLogger() {
	if requestLogWriter != nil {
		requestLogWriter.Close()
	}
}

// logRequest logs HTTP request details (method, path, remote address, duration)
// Logs are written to a rotating file (output/logs/access.log) that rotates at 5MB or 2 weeks
func (app *application) logRequest(next http.Handler) http.Handler {
	// Initialize request logger on first use (also write to console in dev mode)
	alsoWriteToConsole := app.config.env == "dev"
	if err := initRequestLogger(alsoWriteToConsole); err != nil {
		// Fallback to app logger if rotation init fails
		app.logger.Printf("Warning: Failed to initialize request logger: %v", err)
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Create a response writer wrapper to capture status code
		rw := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}

		// Call the next handler
		next.ServeHTTP(rw, r)

		// Log the request details to rotating file (and console in dev)
		if requestLogger != nil {
			requestLogger.Printf("%s %s %s %s %d %v",
				r.RemoteAddr,
				r.Proto,
				r.Method,
				r.URL.RequestURI(),
				rw.statusCode,
				time.Since(start),
			)
		} else {
			// Fallback to app logger if request logger not initialized
			app.logger.Printf("%s %s %s %s %d %v",
				r.RemoteAddr,
				r.Proto,
				r.Method,
				r.URL.RequestURI(),
				rw.statusCode,
				time.Since(start),
			)
		}
	})
}

// responseWriter wraps http.ResponseWriter to capture status code
type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

// csrfTokenGenerator is middleware that generates and sets CSRF tokens for GET requests
// and stores the token in request context for use in templates
func (app *application) csrfTokenGenerator(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Only generate tokens for GET requests (that might render templates)
		if r.Method == http.MethodGet {
			token, err := csrf.GenerateAndSetToken(w, r, app.config.env)
			if err != nil {
				app.logger.Printf("Error generating CSRF token: %v", err)
				// Continue anyway - token generation failure shouldn't break the request
			} else {
				// Store token in context for handlers to access
				ctx := context.WithValue(r.Context(), csrf.TokenKey, token)
				r = r.WithContext(ctx)
			}
		}

		next.ServeHTTP(w, r)
	})
}

// csrfProtect is middleware that validates CSRF tokens for non-safe HTTP methods
func (app *application) csrfProtect(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Skip CSRF validation for safe methods
		if csrf.IsSafeMethod(r.Method) {
			next.ServeHTTP(w, r)
			return
		}

		// Validate CSRF token for non-safe methods
		if err := csrf.ValidateToken(r); err != nil {
			app.logger.Printf("CSRF validation failed: %v - %s %s from %s", err, r.Method, r.URL.Path, r.RemoteAddr)
			http.Error(w, "Forbidden: Invalid CSRF token", http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// secureHeaders adds security headers to all responses
// Should be applied globally to protect against common web vulnerabilities
func secureHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "deny")
		w.Header().Set("X-XSS-Protection", "1; mode=block")
		w.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")
		next.ServeHTTP(w, r)
	})
}

// cacheControl adds cache headers for static files
// Should only be applied to static file routes, not dynamic content
func cacheControl(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Cache static files for 1 hour
		w.Header().Set("Cache-Control", "public, max-age=3600")
		next.ServeHTTP(w, r)
	})
}

// newSessionManager creates and configures a new session manager
func newSessionManager(env string) *scs.SessionManager {
	sessionManager := scs.New()
	sessionManager.Lifetime = 24 * time.Hour
	sessionManager.Cookie.Secure = helpers.IsSecureCookie(env)
	sessionManager.Cookie.HttpOnly = true
	sessionManager.Cookie.SameSite = http.SameSiteStrictMode
	return sessionManager
}

// sessionMiddleware wraps the session manager's LoadAndSave middleware
func (app *application) sessionMiddleware(next http.Handler) http.Handler {
	return app.session.LoadAndSave(next)
}
