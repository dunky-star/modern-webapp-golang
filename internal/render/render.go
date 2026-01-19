package render

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"

	"github.com/alexedwards/scs/v2"
	"github.com/dunky-star/modern-webapp-golang/internal/config"
	"github.com/dunky-star/modern-webapp-golang/internal/data"
	"github.com/dunky-star/modern-webapp-golang/internal/helpers"
	"github.com/dunky-star/modern-webapp-golang/pkg/csrf"
)

var app *config.AppConfig
var pathToTemplates = "./web"

// SessionManagerKey is the context key for the session manager (exported for middleware use)
type SessionManagerKey struct{}

// NewTemplates sets the config for the template package
func NewRender(a *config.AppConfig) {
	app = a
}

// AddDefaultData adds data for all templates
func AddDefaultData(td *data.TemplateData, r *http.Request) *data.TemplateData {
	// Extract session manager from context (injected by middleware)
	if session := getSessionManagerFromContext(r.Context()); session != nil {
		td.Flash = session.PopString(r.Context(), "flash")
		td.Warning = session.PopString(r.Context(), "warning")
		td.Error = session.PopString(r.Context(), "error")
	}

	// Get CSRF token from context (injected by middleware)
	if token := getCSRFTokenFromContext(r.Context()); token != "" {
		td.SetCSRFToken(token)
	}
	// Check if the user is authenticated
	td.IsAuthenticated = helpers.IsAuthenticated(r) || app.Session.Exists(r.Context(), "user_id")

	return td
}

// TemplateCache renders a template
func TemplateCache(w http.ResponseWriter, r *http.Request, tmpl string, td *data.TemplateData) error {
	var tc map[string]*template.Template

	if app != nil && app.TemplateCache != nil {
		// Always use the pre-built template cache from app config
		tc = app.TemplateCache
	} else {
		// Fallback: create cache on-the-fly if not initialized (shouldn't happen in normal operation)
		tc, _ = CreateTemplateCache()
	}

	t, ok := tc[tmpl]
	if !ok {
		return errors.New("could not get template from cache")
	}

	buf := new(bytes.Buffer)

	td = AddDefaultData(td, r)

	_ = t.Execute(buf, td)

	_, err := buf.WriteTo(w)
	if err != nil {
		fmt.Println("error writing template to browser", err)
		return err
	}

	return nil
}

// CreateTemplateCache creates a template cache as a map
func CreateTemplateCache() (map[string]*template.Template, error) {
	myCache := map[string]*template.Template{}

	pages, err := filepath.Glob(fmt.Sprintf("%s/*.page.tmpl", pathToTemplates))
	if err != nil {
		return myCache, err
	}

	for _, page := range pages {
		name := filepath.Base(page)
		ts, err := template.New(name).ParseFiles(page)
		if err != nil {
			return myCache, err
		}

		matches, err := filepath.Glob(fmt.Sprintf("%s/*.layout.tmpl", pathToTemplates))
		if err != nil {
			return myCache, err
		}

		if len(matches) > 0 {
			ts, err = ts.ParseGlob(fmt.Sprintf("%s/*.layout.tmpl", pathToTemplates))
			if err != nil {
				return myCache, err
			}
		}

		myCache[name] = ts
	}

	return myCache, nil
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
