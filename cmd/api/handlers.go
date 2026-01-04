package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/dunky-star/modern-webapp-golang/pkg/render"
)

func (app *application) healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	// Calculate uptime dynamically
	uptime := time.Since(startTime).Truncate(time.Second)
	status := map[string]interface{}{
		"status":    "available",
		"uptime":    uptime.String(),
		"timestamp": time.Now().Format(time.RFC3339),
	}
	fmt.Fprintf(w, "Version: %s\n", version)
	if err := json.NewEncoder(w).Encode(status); err != nil {
		http.Error(w, fmt.Sprintf(`{"error":"%s"}`, err.Error()), http.StatusInternalServerError)
	}
}

func (app *application) homeHandler(w http.ResponseWriter, r *http.Request) {
	render.Template(w, app.logger, "home.page.tmpl")
}

func (app *application) aboutUsHandler(w http.ResponseWriter, r *http.Request) {
	render.Template(w, app.logger, "about.page.tmpl")
}
