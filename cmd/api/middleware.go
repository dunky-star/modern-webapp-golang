package main

import (
	"context"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/dunky-star/modern-webapp-golang/internal/render"
	"github.com/dunky-star/modern-webapp-golang/pkg/csrf"
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
func logRequest(next http.Handler) http.Handler {
	// Initialize request logger on first use (also write to console in dev mode)
	alsoWriteToConsole := app.Env == "dev"
	if err := initRequestLogger(alsoWriteToConsole); err != nil {
		// Fallback to app logger if rotation init fails
		app.WarningLog.Printf("Failed to initialize request logger: %v", err)
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Create a response writer wrapper to capture status code
		rw := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}

		// Call the next handler
		next.ServeHTTP(rw, r)

		// Log the request details to rotating file (and console in dev)
		// Use string builder for better performance with multiple string operations
		duration := time.Since(start)
		if requestLogger != nil {
			requestLogger.Printf("%s %s %s %s %d %v",
				r.RemoteAddr,
				r.Proto,
				r.Method,
				r.URL.RequestURI(),
				rw.statusCode,
				duration,
			)
		} else {
			// Fallback to app logger if request logger not initialized
			app.InfoLog.Printf("%s %s %s %s %d %v",
				r.RemoteAddr,
				r.Proto,
				r.Method,
				r.URL.RequestURI(),
				rw.statusCode,
				duration,
			)
		}
	})
}

// responseWriter wraps http.ResponseWriter to capture status code
type responseWriter struct {
	http.ResponseWriter
	statusCode    int
	headerWritten bool
}

func (rw *responseWriter) WriteHeader(code int) {
	if rw.headerWritten {
		return // Prevent multiple WriteHeader calls
	}
	rw.statusCode = code
	rw.headerWritten = true
	rw.ResponseWriter.WriteHeader(code)
}

func (rw *responseWriter) Write(b []byte) (int, error) {
	if !rw.headerWritten {
		rw.WriteHeader(http.StatusOK)
	}
	return rw.ResponseWriter.Write(b)
}

// csrfTokenGenerator is middleware that generates and sets CSRF tokens for GET requests
// and stores the token in request context for use in templates
func csrfTokenGenerator(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Only generate tokens for GET requests (that might render templates)
		if r.Method == http.MethodGet {
			token, err := csrf.GenerateAndSetToken(w, r, app.Env)
			if err != nil {
				app.ErrorLog.Printf("Error generating CSRF token: %v", err)
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
// Parses form data for POST/PUT/PATCH requests before validation
func csrfProtect(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Skip CSRF validation for safe methods
		if csrf.IsSafeMethod(r.Method) {
			next.ServeHTTP(w, r)
			return
		}

		// Parse form for non-safe methods to extract CSRF token
		// This is safe because:
		// 1. Go's http.Request.ParseForm() has built-in size limits (10MB default)
		// 2. It's idempotent - safe to call multiple times
		// 3. We only parse for methods that need CSRF validation
		if err := r.ParseForm(); err != nil {
			app.ErrorLog.Printf("Error parsing form for CSRF validation: %v - %s %s from %s", err, r.Method, r.URL.Path, r.RemoteAddr)
			http.Error(w, "Bad Request: Invalid form data", http.StatusBadRequest)
			return
		}

		// Validate CSRF token for non-safe methods
		if err := csrf.ValidateToken(r); err != nil {
			app.ErrorLog.Printf("CSRF validation failed: %v - %s %s from %s", err, r.Method, r.URL.Path, r.RemoteAddr)
			// For form submissions, redirect back to the same path (as GET) with error message
			// This provides better UX than showing a 403 error page
			app.Session.Put(r.Context(), "error", "Your session has expired. Please fill out the form again.")
			http.Redirect(w, r, r.URL.Path, http.StatusSeeOther)
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

// htmlCacheControl adds appropriate cache headers for dynamic HTML pages
// Uses "no-cache" to allow conditional requests (ETags) while preventing stale content
// The CSRF token system already reuses valid tokens from cookies, so we don't need "no-store"
func htmlCacheControl(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Only apply to GET requests (HTML pages)
		if r.Method == http.MethodGet {
			// "no-cache" allows browsers to cache but requires revalidation
			// This enables ETag/If-None-Match for conditional requests
			// Vary on Cookie since CSRF tokens and sessions are cookie-based
			w.Header().Set("Cache-Control", "no-cache, must-revalidate")
			w.Header().Set("Vary", "Cookie")
		}
		next.ServeHTTP(w, r)
	})
}

// sessionMiddleware wraps the session manager's LoadAndSave middleware
// and injects the session manager into request context for automatic access
func sessionMiddleware(next http.Handler) http.Handler {
	return app.Session.LoadAndSave(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Inject session manager into context for automatic access by render package
		ctx := context.WithValue(r.Context(), render.SessionManagerKey{}, app.Session)
		next.ServeHTTP(w, r.WithContext(ctx))
	}))
}
