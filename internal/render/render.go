package render

import (
	"context"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"sync"

	"github.com/dunky-star/modern-webapp-golang/pkg/csrf"
)

var (
	templateCache = make(map[string]*template.Template)
	mu            sync.RWMutex
)

// TemplateCache renders a template using cache. Set useCache to false to always reload templates (useful for development)
// If data is nil, no data will be passed to the template
// Automatically injects CSRF token from request context if data is TemplateData
func TemplateCache(w http.ResponseWriter, r *http.Request, logger *log.Logger, t string, useCache bool, data interface{}) {
	// Automatically inject CSRF token from context if data is TemplateData
	if r != nil && data != nil {
		if tmplData, ok := data.(interface {
			SetCSRFToken(string)
		}); ok {
			if token := getCSRFTokenFromContext(r.Context()); token != "" {
				tmplData.SetCSRFToken(token)
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

	err = tmpl.Execute(w, data)
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

// getCSRFTokenFromContext retrieves CSRF token from request context
// Uses the same string key as csrf middleware
func getCSRFTokenFromContext(ctx context.Context) string {
	if token, ok := ctx.Value(csrf.TokenKey).(string); ok {
		return token
	}
	return ""
}
