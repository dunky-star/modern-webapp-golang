package main

import (
	"fmt"
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
	fmt.Fprintf(w, "This is the home page")
}

func (app *application) aboutUsHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "This is the about page")
}
