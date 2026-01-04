package main

import (
	"fmt"
	"html/template"
	"net/http"
)

func (app *application) healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	n, err := fmt.Fprintf(w, "Health check")
	if err != nil {
		app.logger.Println(err)
	}
	app.logger.Printf("Number of bytes written: %d\n", n)
}

func (app *application) homeHandler(w http.ResponseWriter, r *http.Request) {
	app.renderTemplate(w, "home.page.tmpl")
}

func (app *application) aboutUsHandler(w http.ResponseWriter, r *http.Request) {
	app.renderTemplate(w, "about.page.tmpl")
}

func (app *application) renderTemplate(w http.ResponseWriter, tmpl string) {
	parsedTemplate, err := template.ParseFiles("./web/" + tmpl)
	if err != nil {
		app.logger.Println(err)
		return
	}
	err = parsedTemplate.Execute(w, nil)
	if err != nil {
		app.logger.Println(err)
	}
}
