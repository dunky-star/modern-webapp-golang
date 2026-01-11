package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/dunky-star/modern-webapp-golang/internal/data"
	appdata "github.com/dunky-star/modern-webapp-golang/internal/data"
	"github.com/dunky-star/modern-webapp-golang/internal/forms"
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

func (app *application) makeReservationHandler(w http.ResponseWriter, r *http.Request) {
	useCache := app.config.env != "dev"

	// Create empty reservation for initial form display
	var emptyReservation data.Reservation

	// Create template data with Form and Data
	tmplData := appdata.NewTemplateData()
	tmplData.Data["Title"] = "Make Reservation"
	tmplData.Data["reservation"] = emptyReservation
	tmplData.Form = forms.New(nil) // Empty form for GET request (no errors initially)

	render.TemplateCache(w, r, app.logger, "make-reservation.page.tmpl", useCache, tmplData)
}

func (app *application) postReservationHandler(w http.ResponseWriter, r *http.Request) {
	useCache := app.config.env != "dev"

	// Parse form data
	err := r.ParseForm()
	if err != nil {
		app.logger.Printf("Error parsing form: %v", err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	// Create form with posted data
	form := forms.New(r.PostForm)

	// Validate required fields using forms package
	form.Required("first_name", "last_name", "email", "phone")

	// Validate email format using forms package
	form.IsEmail("email")

	// Validate minimum length (e.g., phone should be at least 10 characters)
	form.MinLength("phone", 10, r)

	// Create reservation from form data
	reservation := data.Reservation{
		FirstName: r.Form.Get("first_name"),
		LastName:  r.Form.Get("last_name"),
		Email:     r.Form.Get("email"),
		Phone:     r.Form.Get("phone"),
	}

	// Create template data
	tmplData := appdata.NewTemplateData()
	tmplData.Data["Title"] = "Make Reservation"
	tmplData.Data["reservation"] = reservation
	tmplData.Form = form // Pass form with validation errors (if any)

	// If form is invalid, re-render the form with errors
	if !form.Valid() {
		render.TemplateCache(w, r, app.logger, "make-reservation.page.tmpl", useCache, tmplData)
		return
	}

	// Form is valid - process the reservation
	// TODO: Save reservation to database, send confirmation email, etc.
	app.logger.Printf("Reservation created: %+v", reservation)

	// Store reservation in session for summary page
	app.session.Put(r.Context(), "reservation", reservation)

	// Redirect to reservation summary page
	http.Redirect(w, r, "/reservation-summary", http.StatusSeeOther)
}

// ReservationSummary displays the res summary page
func (app *application) reservationSummary(w http.ResponseWriter, r *http.Request) {
	useCache := app.config.env != "dev"

	// Get reservation from session
	reservation, ok := app.session.Get(r.Context(), "reservation").(data.Reservation)
	if !ok {
		app.logger.Println("can't get item from session")
		app.session.Put(r.Context(), "error", "Can't get reservation from session")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	// Remove reservation from session after retrieving it
	app.session.Remove(r.Context(), "reservation")

	// Create template data with reservation
	tmplData := appdata.NewTemplateData()
	tmplData.Data["Title"] = "Reservation Summary"
	tmplData.Data["reservation"] = reservation

	render.TemplateCache(w, r, app.logger, "reservation-summary.page.tmpl", useCache, tmplData)
}
