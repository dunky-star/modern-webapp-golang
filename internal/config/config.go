package config

import (
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/alexedwards/scs/v2"
)

// AppConfig holds the application configuration
type AppConfig struct {
	Port          int
	Env           string
	DSN           string
	InfoLog       *log.Logger
	ErrorLog      *log.Logger
	WarningLog    *log.Logger
	Session       *scs.SessionManager
	UseCache      bool
	TemplateCache map[string]*template.Template
}

// New creates a new application configuration
func New(port int, env string, dsn string) *AppConfig {
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	warningLog := log.New(os.Stdout, "WARNING\t", log.Ldate|log.Ltime|log.Lshortfile)
	session := newSessionManager(env)

	return &AppConfig{
		Port:       port,
		Env:        env,
		DSN:        dsn,
		InfoLog:    infoLog,
		ErrorLog:   errorLog,
		WarningLog: warningLog,
		Session:    session,
	}
}

// IsSecureCookie returns true if cookies should use the Secure flag (HTTPS only)
func IsSecureCookie(env string) bool {
	return env == "prod"
}

// newSessionManager creates and configures a new session manager
func newSessionManager(env string) *scs.SessionManager {
	sessionManager := scs.New()
	sessionManager.Lifetime = 24 * time.Hour
	sessionManager.Cookie.Secure = IsSecureCookie(env)
	sessionManager.Cookie.HttpOnly = true
	sessionManager.Cookie.SameSite = http.SameSiteStrictMode
	return sessionManager
}
