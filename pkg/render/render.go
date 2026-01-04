package render

import (
	"html/template"
	"log"
	"net/http"
)

// Template renders an HTML template from the web directory
func Template(w http.ResponseWriter, logger *log.Logger, tmpl string) {
	parsedTemplate, err := template.ParseFiles("./web/" + tmpl)
	if err != nil {
		logger.Println(err)
		return
	}
	err = parsedTemplate.Execute(w, nil)
	if err != nil {
		logger.Println(err)
	}
}
