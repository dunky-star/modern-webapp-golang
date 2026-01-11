package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	appdata "github.com/dunky-star/modern-webapp-golang/internal/data"
	"github.com/dunky-star/modern-webapp-golang/internal/render"
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
	// Disable cache in dev mode to see template changes immediately
	useCache := app.config.env != "dev"

	remoteIPAddr := r.RemoteAddr
	app.logger.Printf("Remote address: %s", remoteIPAddr)
	app.session.Put(r.Context(), "remote_addr", remoteIPAddr)

	tmplData := appdata.NewTemplateData()
	tmplData.Data["Title"] = "Home, welcome!"
	tmplData.StringMap = map[string]string{
		"remote_addr": remoteIPAddr,
	}

	// CSRF token automatically injected by render.TemplateCache from middleware context
	render.TemplateCache(w, r, app.logger, "home.page.tmpl", useCache, tmplData)
}

func (app *application) aboutUsHandler(w http.ResponseWriter, r *http.Request) {
	// Disable cache in dev mode to see template changes immediately
	useCache := app.config.env != "dev"

	tmplData := appdata.NewTemplateData()
	tmplData.Data["Title"] = "About Us"

	// Retrieve remote address from session
	if remoteAddr := app.session.GetString(r.Context(), "remote_addr"); remoteAddr != "" {
		tmplData.StringMap = map[string]string{
			"remote_addr": remoteAddr,
		}
	}

	// CSRF token automatically injected by render.TemplateCache from middleware context
	render.TemplateCache(w, r, app.logger, "about.page.tmpl", useCache, tmplData)
}

func (app *application) generalsQuartersHandler(w http.ResponseWriter, r *http.Request) {
	useCache := app.config.env != "dev"
	tmplData := appdata.NewTemplateData()
	tmplData.Data["Title"] = "Generals Quarters"
	render.TemplateCache(w, r, app.logger, "generals.page.tmpl", useCache, tmplData)
}

func (app *application) majorsSuiteHandler(w http.ResponseWriter, r *http.Request) {
	useCache := app.config.env != "dev"
	tmplData := appdata.NewTemplateData()
	tmplData.Data["Title"] = "Majors Suite"
	render.TemplateCache(w, r, app.logger, "majors.page.tmpl", useCache, tmplData)
}

func (app *application) makeReservationHandler(w http.ResponseWriter, r *http.Request) {
	useCache := app.config.env != "dev"
	tmplData := appdata.NewTemplateData()
	tmplData.Data["Title"] = "Make Reservation"
	render.TemplateCache(w, r, app.logger, "make-reservation.page.tmpl", useCache, tmplData)
}

func (app *application) searchAvailabilityHandler(w http.ResponseWriter, r *http.Request) {
	useCache := app.config.env != "dev"
	tmplData := appdata.NewTemplateData()
	tmplData.Data["Title"] = "Search Availability"
	render.TemplateCache(w, r, app.logger, "search-availability.page.tmpl", useCache, tmplData)
}

func (app *application) postAvailabilityHandler(w http.ResponseWriter, r *http.Request) {
	start := r.FormValue("start")
	end := r.FormValue("end")
	w.Write([]byte(fmt.Sprintf("start date is %s and end date is %s", start, end)))
}
func (app *application) contactHandler(w http.ResponseWriter, r *http.Request) {
	useCache := app.config.env != "dev"
	tmplData := appdata.NewTemplateData()
	tmplData.Data["Title"] = "Contact Us"
	render.TemplateCache(w, r, app.logger, "contact.page.tmpl", useCache, tmplData)
}

type jsonResponse struct {
	OK      bool   `json:"ok"`
	Message string `json:"message"`
}

func (app *application) avialabilityJSONHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	response := jsonResponse{
		OK:      true,
		Message: "Available!",
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		app.logger.Println(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}
