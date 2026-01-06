package main

import (
	"log"
	"net/http"
	"sync"
	"time"
)

var (
	requestLogger     *log.Logger
	requestLogWriter  *rotatingLogWriter
	requestLoggerOnce sync.Once
)

// initRequestLogger initializes the rotating log writer for HTTP request logs
func initRequestLogger(alsoWriteToConsole bool) error {
	var initErr error
	requestLoggerOnce.Do(func() {
		rotatingWriter, err := newRotatingLogWriter(alsoWriteToConsole)
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
// Logs are written to a rotating file (log/access.log) that rotates at 5MB or 2 weeks
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
