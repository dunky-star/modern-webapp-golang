package main

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"fmt"
	"net/http"
	"time"
)

// CSRFTokenKey is the context key for CSRF tokens (string key for cross-package compatibility)
const CSRFTokenKey = "csrf_token"

const (
	csrfTokenLength = 32
	csrfCookieName  = "__csrf_token"
	csrfHeaderName  = "X-CSRF-Token"
	csrfFormField   = "csrf_token"
	csrfMaxAge      = 12 * time.Hour // Token valid for 12 hours
)

// generateCSRFToken generates a cryptographically secure random token
func generateCSRFToken() (string, error) {
	b := make([]byte, csrfTokenLength)
	if _, err := rand.Read(b); err != nil {
		return "", fmt.Errorf("failed to generate CSRF token: %w", err)
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

// setCSRFCookie sets the CSRF token as an HTTP-only cookie
func (app *application) setCSRFCookie(w http.ResponseWriter, token string) {
	secure := app.config.env == "prod" // Only use Secure flag in production (HTTPS)
	cookie := &http.Cookie{
		Name:     csrfCookieName,
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		Secure:   secure,
		SameSite: http.SameSiteStrictMode,
		MaxAge:   int(csrfMaxAge.Seconds()),
	}
	http.SetCookie(w, cookie)
}

// getCSRFCookie retrieves the CSRF token from the cookie
func getCSRFCookie(r *http.Request) (string, error) {
	cookie, err := r.Cookie(csrfCookieName)
	if err != nil {
		return "", err
	}
	return cookie.Value, nil
}

// getCSRFTokenFromRequest retrieves CSRF token from header or form
func getCSRFTokenFromRequest(r *http.Request) string {
	// Check header first (for AJAX requests)
	if token := r.Header.Get(csrfHeaderName); token != "" {
		return token
	}
	// Check form field (for traditional form submissions)
	if token := r.FormValue(csrfFormField); token != "" {
		return token
	}
	return ""
}

// validateCSRFToken compares the token from request with the cookie token
func validateCSRFToken(r *http.Request) error {
	cookieToken, err := getCSRFCookie(r)
	if err != nil {
		return fmt.Errorf("CSRF token cookie not found")
	}

	requestToken := getCSRFTokenFromRequest(r)
	if requestToken == "" {
		return fmt.Errorf("CSRF token not found in request")
	}

	// Use constant-time comparison to prevent timing attacks
	if subtle.ConstantTimeCompare([]byte(cookieToken), []byte(requestToken)) != 1 {
		return fmt.Errorf("CSRF token mismatch")
	}

	return nil
}

// isSafeMethod checks if the HTTP method is safe (doesn't modify state)
func isSafeMethod(method string) bool {
	safeMethods := []string{http.MethodGet, http.MethodHead, http.MethodOptions, http.MethodTrace}
	for _, m := range safeMethods {
		if method == m {
			return true
		}
	}
	return false
}

// generateAndSetCSRFToken generates a new CSRF token and sets it as a cookie
// Returns the token for use in templates
func (app *application) generateAndSetCSRFToken(w http.ResponseWriter, r *http.Request) (string, error) {
	// Check if token already exists and is valid
	if existingToken, err := getCSRFCookie(r); err == nil && existingToken != "" {
		return existingToken, nil
	}

	// Generate new token
	token, err := generateCSRFToken()
	if err != nil {
		return "", err
	}

	// Set cookie
	app.setCSRFCookie(w, token)

	return token, nil
}
