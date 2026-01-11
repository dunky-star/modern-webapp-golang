package render

import (
	"context"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"sync"

	"github.com/alexedwards/scs/v2"
	"github.com/dunky-star/modern-webapp-golang/internal/data"
	"github.com/dunky-star/modern-webapp-golang/pkg/csrf"
)

var (
	templateCache = make(map[string]*template.Template)
	mu            sync.RWMutex
)

// SessionManagerKey is the context key for the session manager (exported for middleware use)
type SessionManagerKey struct{}

// TemplateCache renders a template using cache. Set useCache to false to always reload templates (useful for development)
// If templateData is nil, no data will be passed to the template
// Automatically injects CSRF token and session data (flash, warning, error) from request context
func TemplateCache(w http.ResponseWriter, r *http.Request, logger *log.Logger, t string, useCache bool, templateData interface{}) {
	// Add default data (flash, warning, error, CSRF token) if templateData is TemplateData
	if r != nil && templateData != nil {
		if td, ok := templateData.(*data.TemplateData); ok {
			// Automatically extract session manager from context and add default data
			addDefaultData(td, r)
		} else {
			// For non-TemplateData types, try to inject CSRF token if possible
			if tmplData, ok := templateData.(interface {
				SetCSRFToken(string)
			}); ok {
				if token := getCSRFTokenFromContext(r.Context()); token != "" {
					tmplData.SetCSRFToken(token)
				}
			}
		}
	}
	var tmpl *template.Template
	var err error

	if useCache {
		// Check cache with read lock first
		mu.RLock()
		var inCache bool
		tmpl, inCache = templateCache[t]
		mu.RUnlock()

		if !inCache {
			// Acquire write lock to create template
			mu.Lock()
			// Double-check after acquiring lock (another goroutine might have created it)
			tmpl, inCache = templateCache[t]
			if !inCache {
				logger.Println("Creating template and adding to cache")
				err = createTemplateCache(t)
				if err != nil {
					mu.Unlock()
					logger.Println(err)
					return
				}
				tmpl = templateCache[t]
			}
			mu.Unlock()
		} else {
			logger.Println("Using template from cache")
		}
	} else {
		// Development mode: always reload template
		templates := []string{
			fmt.Sprintf("./web/%s", t),
			"./web/base.layout.tmpl",
		}
		var parseErr error
		tmpl, parseErr = template.ParseFiles(templates...)
		if parseErr != nil {
			logger.Println(parseErr)
			return
		}
	}

	err = tmpl.Execute(w, templateData)
	if err != nil {
		logger.Println(err)
		return
	}
}

func createTemplateCache(t string) error {
	templates := []string{
		fmt.Sprintf("./web/%s", t),
		"./web/base.layout.tmpl",
	}

	// Parse the templates
	tmpl, err := template.ParseFiles(templates...)
	if err != nil {
		return err
	}

	// Add the template to the cache (caller holds the lock)
	templateCache[t] = tmpl
	return nil
}

// addDefaultData adds default data for all templates (flash, warning, error messages, CSRF token)
// Automatically extracts session manager from request context
// Pops messages from session so they're only shown once
func addDefaultData(td *data.TemplateData, r *http.Request) {
	// Extract session manager from context (injected by middleware)
	if session := getSessionManagerFromContext(r.Context()); session != nil {
		// Pop flash, warning, and error messages from session (one-time display)
		td.Flash = session.PopString(r.Context(), "flash")
		td.Warning = session.PopString(r.Context(), "warning")
		td.Error = session.PopString(r.Context(), "error")
	}

	// Get CSRF token from context (injected by middleware)
	if token := getCSRFTokenFromContext(r.Context()); token != "" {
		td.SetCSRFToken(token)
	}
}

// getSessionManagerFromContext retrieves session manager from request context
func getSessionManagerFromContext(ctx context.Context) *scs.SessionManager {
	if session, ok := ctx.Value(SessionManagerKey{}).(*scs.SessionManager); ok {
		return session
	}
	return nil
}

// getCSRFTokenFromContext retrieves CSRF token from request context
// Uses the same string key as csrf middleware
func getCSRFTokenFromContext(ctx context.Context) string {
	if token, ok := ctx.Value(csrf.TokenKey).(string); ok {
		return token
	}
	return ""
}
