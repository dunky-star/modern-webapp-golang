package config

import (
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/dunky-star/modern-webapp-golang/pkg/helpers"
)

// AppConfig holds the application configuration
type AppConfig struct {
	Port          int
	Env           string
	InfoLog       *log.Logger
	Session       *scs.SessionManager
	UseCache      bool
	TemplateCache map[string]*template.Template
}

// New creates a new application configuration
func New(port int, env string) *AppConfig {
	infoLog := log.New(os.Stdout, "", log.Ldate|log.Ltime)
	session := newSessionManager(env)

	return &AppConfig{
		Port:    port,
		Env:     env,
		InfoLog: infoLog,
		Session: session,
	}
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
