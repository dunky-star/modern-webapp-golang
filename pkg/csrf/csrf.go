package csrf

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"fmt"
	"net/http"
	"time"

	"github.com/dunky-star/modern-webapp-golang/internal/config"
)

// TokenKey is the context key for CSRF tokens (string key for cross-package compatibility)
const TokenKey = "csrf_token"

const (
	tokenLength = 32
	CookieName  = "__csrf_token"
	HeaderName  = "X-CSRF-Token"
	FormField   = "csrf_token"
	MaxAge      = 12 * time.Hour // Token valid for 12 hours
)

// GenerateToken generates a cryptographically secure random token
func GenerateToken() (string, error) {
	b := make([]byte, tokenLength)
	if _, err := rand.Read(b); err != nil {
		return "", fmt.Errorf("failed to generate CSRF token: %w", err)
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

// SetCookie sets the CSRF token as an HTTP-only cookie
func SetCookie(w http.ResponseWriter, token string, env string) {
	cookie := &http.Cookie{
		Name:     CookieName,
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		Secure:   config.IsSecureCookie(env),
		SameSite: http.SameSiteStrictMode,
		MaxAge:   int(MaxAge.Seconds()),
	}
	http.SetCookie(w, cookie)
}

// GetCookie retrieves the CSRF token from the cookie
func GetCookie(r *http.Request) (string, error) {
	cookie, err := r.Cookie(CookieName)
	if err != nil {
		return "", err
	}
	return cookie.Value, nil
}

// GetTokenFromRequest retrieves CSRF token from header or form
// Assumes form is already parsed for POST requests (should be done in middleware)
func GetTokenFromRequest(r *http.Request) string {
	// Check header first (for AJAX requests)
	if token := r.Header.Get(HeaderName); token != "" {
		return token
	}

	// For POST/PUT/PATCH requests, check PostForm directly (form must be parsed by middleware)
	if r.Method == http.MethodPost || r.Method == http.MethodPut || r.Method == http.MethodPatch {
		if r.PostForm != nil {
			if token := r.PostForm.Get(FormField); token != "" {
				return token
			}
		}
	}

	// Fallback to FormValue (works for GET query params and already-parsed forms)
	if token := r.FormValue(FormField); token != "" {
		return token
	}
	return ""
}

// ValidateToken compares the token from request with the cookie token
func ValidateToken(r *http.Request) error {
	cookieToken, err := GetCookie(r)
	if err != nil {
		return fmt.Errorf("CSRF token cookie not found")
	}

	requestToken := GetTokenFromRequest(r)
	if requestToken == "" {
		return fmt.Errorf("CSRF token not found in request")
	}

	// Use constant-time comparison to prevent timing attacks
	if subtle.ConstantTimeCompare([]byte(cookieToken), []byte(requestToken)) != 1 {
		return fmt.Errorf("CSRF token mismatch")
	}

	return nil
}

// IsSafeMethod checks if the HTTP method is safe (doesn't modify state)
func IsSafeMethod(method string) bool {
	safeMethods := []string{http.MethodGet, http.MethodHead, http.MethodOptions, http.MethodTrace}
	for _, m := range safeMethods {
		if method == m {
			return true
		}
	}
	return false
}

// GenerateAndSetToken generates a new CSRF token and sets it as a cookie
// Returns the token for use in templates
func GenerateAndSetToken(w http.ResponseWriter, r *http.Request, env string) (string, error) {
	// Check if token already exists and is valid
	if existingToken, err := GetCookie(r); err == nil && existingToken != "" {
		return existingToken, nil
	}

	// Generate new token
	token, err := GenerateToken()
	if err != nil {
		return "", err
	}

	// Set cookie
	SetCookie(w, token, env)

	return token, nil
}
