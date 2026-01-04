package render

import (
	"html/template"
	"log"
	"net/http"
)

// Template renders an HTML template from the web directory
func Template(w http.ResponseWriter, logger *log.Logger, tmpl string) {
	parsedTemplate, err := template.ParseFiles("./web/"+tmpl, "./web/base.layout.tmpl")
	if err != nil {
		logger.Println(err)
		return
	}
	err = parsedTemplate.Execute(w, nil)
	if err != nil {
		logger.Println(err)
	}
}

var templateCache = make(map[string]*template.Template)

func TemplateCache(w http.ResponseWriter, logger *log.Logger, t string) {
	var tmpl *template.Template
	var err error

	// Check to see if we already have the template in our cache
	tmpl, inCache := templateCache[t]
	if !inCache {
		// Need to check the template and add it to the cache
		tmpl, err = template.ParseFiles("./web/"+t, "./web/base.layout.tmpl")
		if err != nil {
			logger.Println(err)
			return
		}
		// Store in cache for future use
		templateCache[t] = tmpl
	}

	err = tmpl.Execute(w, nil)
	if err != nil {
		logger.Println(err)
	}
}
