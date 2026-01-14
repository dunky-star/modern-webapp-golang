package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/dunky-star/modern-webapp-golang/internal/config"
	"github.com/dunky-star/modern-webapp-golang/internal/data"
	"github.com/dunky-star/modern-webapp-golang/internal/forms"
	"github.com/dunky-star/modern-webapp-golang/internal/helpers"
	"github.com/dunky-star/modern-webapp-golang/internal/render"
	"github.com/dunky-star/modern-webapp-golang/internal/repository"
	"github.com/dunky-star/modern-webapp-golang/internal/repository/dbrepo"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Repo the repository used by the handlers
var Repo *Repository

// Repository is the repository type
type Repository struct {
	app *config.AppConfig
	db  repository.DatabaseConn
}

// NewRepo creates a new repository
func NewRepo(a *config.AppConfig, dbconn *pgxpool.Pool) *Repository {
	return &Repository{
		app: a,
		db:  dbrepo.NewDBConnection(dbconn, a),
	}
}

// NewHandlers sets the repository for the handlers
func NewHandlers(r *Repository) {
	Repo = r
}

var appVersion string
var appStartTime time.Time

// SetAppVersion sets the application version for health check
func SetAppVersion(version string) {
	appVersion = version
}

// SetAppStartTime sets the application start time for health check
func SetAppStartTime(startTime time.Time) {
	appStartTime = startTime
}

func (m *Repository) HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	// Calculate uptime dynamically
	uptime := time.Since(appStartTime).Truncate(time.Second)
	status := map[string]interface{}{
		"version":   appVersion,
		"status":    "available",
		"uptime":    uptime.String(),
		"timestamp": time.Now().Format(time.RFC3339),
	}

	// Use Encoder with SetIndent for pretty-printed JSON that browsers will format nicely
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(status); err != nil {
		helpers.ServerError(w, err)
		return
	}
}

func (m *Repository) HomeHandler(w http.ResponseWriter, r *http.Request) {
	remoteIPAddr := r.RemoteAddr
	m.app.InfoLog.Printf("Remote address: %s", remoteIPAddr)
	m.app.Session.Put(r.Context(), "remote_addr", remoteIPAddr)

	render.TemplateCache(w, r, "home.page.tmpl", &data.TemplateData{
		Data: map[string]interface{}{
			"Title": "Home, welcome!",
		},
		StringMap: map[string]string{
			"remote_addr": remoteIPAddr,
		},
	})
}

func (m *Repository) AboutUsHandler(w http.ResponseWriter, r *http.Request) {
	dataMap := map[string]interface{}{
		"Title": "About Us",
	}
	stringMap := make(map[string]string)

	// Retrieve remote address from session
	if remoteAddr := m.app.Session.GetString(r.Context(), "remote_addr"); remoteAddr != "" {
		stringMap["remote_addr"] = remoteAddr
	}

	render.TemplateCache(w, r, "about.page.tmpl", &data.TemplateData{
		Data:      dataMap,
		StringMap: stringMap,
	})
}

func (m *Repository) GeneralsQuartersHandler(w http.ResponseWriter, r *http.Request) {
	render.TemplateCache(w, r, "generals.page.tmpl", &data.TemplateData{
		Data: map[string]interface{}{
			"Title": "Generals Quarters",
		},
	})
}

func (m *Repository) MajorsSuiteHandler(w http.ResponseWriter, r *http.Request) {
	render.TemplateCache(w, r, "majors.page.tmpl", &data.TemplateData{
		Data: map[string]interface{}{
			"Title": "Majors Suite",
		},
	})
}

func (m *Repository) SearchAvailabilityHandler(w http.ResponseWriter, r *http.Request) {
	render.TemplateCache(w, r, "search-availability.page.tmpl", &data.TemplateData{
		Data: map[string]interface{}{
			"Title": "Search Availability",
		},
	})
}

func (m *Repository) PostAvailabilityHandler(w http.ResponseWriter, r *http.Request) {
	start := r.FormValue("start")
	end := r.FormValue("end")
	w.Write([]byte(fmt.Sprintf("start date is %s and end date is %s", start, end)))
}

func (m *Repository) ContactHandler(w http.ResponseWriter, r *http.Request) {
	render.TemplateCache(w, r, "contact.page.tmpl", &data.TemplateData{
		Data: map[string]interface{}{
			"Title": "Contact Us",
		},
	})
}

type jsonResponse struct {
	OK      bool   `json:"ok"`
	Message string `json:"message"`
}

func (m *Repository) AvialabilityJSONHandler(w http.ResponseWriter, r *http.Request) {
	m.db.AllUsers()
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	response := jsonResponse{
		OK:      true,
		Message: "Available!",
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		helpers.ServerError(w, err)
		return
	}
}

func (m *Repository) MakeReservationHandler(w http.ResponseWriter, r *http.Request) {
	// Create empty reservation for initial form display
	var emptyReservation data.Reservation
	dataMap := make(map[string]interface{})
	dataMap["reservation"] = emptyReservation

	render.TemplateCache(w, r, "make-reservation.page.tmpl", &data.TemplateData{
		Form: forms.New(nil),
		Data: dataMap,
	})
}

func (m *Repository) PostReservationHandler(w http.ResponseWriter, r *http.Request) {
	// Parse form data
	err := r.ParseForm()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	// Create form with posted data
	form := forms.New(r.PostForm)

	// Validate required fields using forms package
	form.Required("first_name", "last_name", "email", "phone", "start_date", "end_date", "room_id")

	// Validate email format using forms package
	form.IsEmail("email")

	// Validate minimum length (e.g., phone should be at least 10 characters)
	form.MinLength("phone", 10, r)

	// Parse date strings to time.Time
	startDate, err := time.Parse("2006-01-02", r.Form.Get("start_date"))
	if err != nil {
		form.Errors.Add("start_date", "Invalid start date format")
	}

	endDate, err := time.Parse("2006-01-02", r.Form.Get("end_date"))
	if err != nil {
		form.Errors.Add("end_date", "Invalid end date format")
	}

	// Parse room_id string to int
	roomId, err := strconv.Atoi(r.Form.Get("room_id"))
	if err != nil {
		form.Errors.Add("room_id", "Invalid room ID")
	}

	// Create reservation from form data
	reservation := data.Reservation{
		FirstName: r.Form.Get("first_name"),
		LastName:  r.Form.Get("last_name"),
		Email:     r.Form.Get("email"),
		Phone:     r.Form.Get("phone"),
		StartDate: startDate,
		EndDate:   endDate,
		RoomId:    roomId,
	}

	dataMap := make(map[string]interface{})
	dataMap["reservation"] = reservation

	// If form is invalid, re-render the form with errors
	if !form.Valid() {
		render.TemplateCache(w, r, "make-reservation.page.tmpl", &data.TemplateData{
			Form: form,
			Data: dataMap,
		})
		return
	}

	// Form is valid - save reservation to database
	newReservationId, err := m.db.InsertReservation(reservation)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	restriction := data.RoomRestriction{
		StartDate:     startDate,
		EndDate:       endDate,
		RoomId:        roomId,
		ReservationId: newReservationId,
		RestrictionId: 1,
	}

	err = m.db.InsertRoomRestriction(restriction)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	m.app.InfoLog.Printf("Reservation created: %+v", reservation)

	// Store reservation in session for summary page
	m.app.Session.Put(r.Context(), "reservation", reservation)

	// Redirect to reservation summary page
	http.Redirect(w, r, "/reservation-summary", http.StatusSeeOther)
}

// ReservationSummary displays the res summary page
func (m *Repository) ReservationSummary(w http.ResponseWriter, r *http.Request) {
	// Get reservation from session
	reservation, ok := m.app.Session.Get(r.Context(), "reservation").(data.Reservation)
	if !ok {
		m.app.ErrorLog.Printf("ERROR\t can't get item from session")
		m.app.Session.Put(r.Context(), "error", "Can't get reservation from session")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	// Remove reservation from session after retrieving it
	m.app.Session.Remove(r.Context(), "reservation")

	dataMap := make(map[string]interface{})
	dataMap["reservation"] = reservation

	render.TemplateCache(w, r, "reservation-summary.page.tmpl", &data.TemplateData{
		Data: dataMap,
	})
}
