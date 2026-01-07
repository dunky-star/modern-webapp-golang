package main

import (
	"net/http"
	"time"

	"github.com/alexedwards/scs/v2"
)

// newSessionManager creates and configures a new session manager
func newSessionManager(env string) *scs.SessionManager {
	sessionManager := scs.New()
	sessionManager.Lifetime = 24 * time.Hour
	sessionManager.Cookie.Secure = env == "prod"
	sessionManager.Cookie.HttpOnly = true
	sessionManager.Cookie.SameSite = http.SameSiteStrictMode
	return sessionManager
}
