package render

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"sync"
)

// TemplateData holds data sent from handlers to templates
type TemplateData struct {
	StringMap map[string]string
	IntMap    map[string]int
	FloatMap  map[string]float32
	Data      map[string]interface{}
	CSRFToken string
	Flash     string
	Warning   string
	Error     string
}

// NewTemplateData creates a new TemplateData with initialized maps
func NewTemplateData() *TemplateData {
	return &TemplateData{
		StringMap: make(map[string]string),
		IntMap:    make(map[string]int),
		FloatMap:  make(map[string]float32),
		Data:      make(map[string]interface{}),
	}
}

var (
	templateCache = make(map[string]*template.Template)
	mu            sync.RWMutex
)

// TemplateCache renders a template using cache. Set useCache to false to always reload templates (useful for development)
// If data is nil, an empty TemplateData will be used
func TemplateCache(w http.ResponseWriter, logger *log.Logger, t string, useCache bool, data *TemplateData) {
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

	// Use provided data or empty TemplateData if nil
	if data == nil {
		data = &TemplateData{}
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
