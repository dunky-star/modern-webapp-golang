package render

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
)

var templateCache = make(map[string]*template.Template)

func TemplateCache(w http.ResponseWriter, logger *log.Logger, t string) {
	var tmpl *template.Template
	var err error

	// Check to see if we already have the template in our cache
	tmpl, inCache := templateCache[t]
	if !inCache {
		// Need to check the template and add it to the cache
		logger.Println("Creating template and adding to cache")
		err = createTemplateCache(t)
		if err != nil {
			logger.Println(err)
			return
		}
		// Get the template from cache after creating it
		tmpl = templateCache[t]
	} else {
		// We have the template in cache, so we can use it
		logger.Println("Using template from cache")
	}

	err = tmpl.Execute(w, nil)
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

	// Add the template to the cache
	templateCache[t] = tmpl
	return nil
}
