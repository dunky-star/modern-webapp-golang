package handlers

import (
	"encoding/json"
	"fmt"
	"html/template"
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
	err := r.ParseForm()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	// Create form with posted data
	form := forms.New(r.PostForm)

	// Validate required fields
	form.Required("start_date", "end_date")

	// Parse date strings to time.Time
	startDate, err := time.Parse("2006-01-02", r.Form.Get("start_date"))
	if err != nil {
		form.Errors.Add("start_date", "Invalid start date format")
	}

	endDate, err := time.Parse("2006-01-02", r.Form.Get("end_date"))
	if err != nil {
		form.Errors.Add("end_date", "Invalid end date format")
	}

	// If form is invalid, re-render the form with errors
	if !form.Valid() {
		render.TemplateCache(w, r, "search-availability.page.tmpl", &data.TemplateData{
			Form: form,
			Data: map[string]interface{}{
				"Title": "Search Availability",
			},
		})
		return
	}

	rooms, err := m.db.SearchAvailabilityForAllRooms(startDate, endDate)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	if len(rooms) == 0 {
		m.app.Session.Put(r.Context(), "error", "No rooms available for the selected dates")
		http.Redirect(w, r, "/search-availability", http.StatusSeeOther)
		return
	}

	// Store reservation with dates in session
	res := data.Reservation{
		StartDate: startDate,
		EndDate:   endDate,
	}
	m.app.Session.Put(r.Context(), "reservation", res)

	// Render choose-room template with available rooms
	dataMap := map[string]interface{}{
		"Title": "Choose Your Room",
		"rooms": rooms,
	}

	render.TemplateCache(w, r, "choose-room.page.tmpl", &data.TemplateData{
		Data: dataMap,
	})
}

// ChooseRoomHandler handles room selection from choose-room page
func (m *Repository) ChooseRoomHandler(w http.ResponseWriter, r *http.Request) {
	// Get reservation from session
	res, ok := m.app.Session.Get(r.Context(), "reservation").(data.Reservation)
	if !ok {
		m.app.Session.Put(r.Context(), "error", "Can't get reservation from session")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	// Get room_id from path parameter
	roomId, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		m.app.Session.Put(r.Context(), "error", "Invalid room selection")
		http.Redirect(w, r, "/search-availability", http.StatusSeeOther)
		return
	}

	// Set room_id in reservation
	res.RoomId = roomId

	// Store back in session
	m.app.Session.Put(r.Context(), "reservation", res)

	// Redirect to make-reservation
	http.Redirect(w, r, "/make-reservation", http.StatusSeeOther)
}

// BookRoom takes URL parameters, builds a sessional variable, and takes user to make res screen
func (m *Repository) BookRoomHandler(w http.ResponseWriter, r *http.Request) {
	roomID, _ := strconv.Atoi(r.URL.Query().Get("id"))
	sd := r.URL.Query().Get("s")
	ed := r.URL.Query().Get("e")

	layout := "2006-01-02"
	startDate, _ := time.Parse(layout, sd)
	endDate, _ := time.Parse(layout, ed)

	var res data.Reservation

	room, err := m.db.GetRoomByID(roomID)
	if err != nil {
		m.app.Session.Put(r.Context(), "error", "Can't get room from db!")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	res.Room.RoomName = room.RoomName
	res.RoomId = roomID
	res.StartDate = startDate
	res.EndDate = endDate

	m.app.Session.Put(r.Context(), "reservation", res)

	http.Redirect(w, r, "/make-reservation", http.StatusSeeOther)
}

func (m *Repository) ContactHandler(w http.ResponseWriter, r *http.Request) {
	render.TemplateCache(w, r, "contact.page.tmpl", &data.TemplateData{
		Data: map[string]interface{}{
			"Title": "Contact Us",
		},
	})
}

type jsonResponse struct {
	OK        bool   `json:"ok"`
	Message   string `json:"message"`
	RoomID    string `json:"room_id"`
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
}

func (m *Repository) AvialabilityJSONHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		resp := jsonResponse{
			OK:      false,
			Message: "Internal server error",
		}
		w.Header().Set("Content-Type", "application/json")
		encoder := json.NewEncoder(w)
		encoder.SetIndent("", "     ")
		encoder.Encode(resp)
		return
	}

	sd := r.Form.Get("start")
	ed := r.Form.Get("end")

	layout := "2006-01-02"
	startDate, _ := time.Parse(layout, sd)
	endDate, _ := time.Parse(layout, ed)

	roomID, _ := strconv.Atoi(r.Form.Get("room_id"))

	available, err := m.db.SearchAvailabilityByDatesByRoomId(startDate, endDate, roomID)
	if err != nil {
		resp := jsonResponse{
			OK:      false,
			Message: "Error querying database",
		}
		w.Header().Set("Content-Type", "application/json")
		encoder := json.NewEncoder(w)
		encoder.SetIndent("", "     ")
		encoder.Encode(resp)
		return
	}

	resp := jsonResponse{
		OK:        available,
		Message:   "",
		StartDate: sd,
		EndDate:   ed,
		RoomID:    strconv.Itoa(roomID),
	}

	w.Header().Set("Content-Type", "application/json")
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "     ")
	encoder.Encode(resp)
}

func (m *Repository) MakeReservationHandler(w http.ResponseWriter, r *http.Request) {
	// Get reservation from session
	res, ok := m.app.Session.Get(r.Context(), "reservation").(data.Reservation)
	if !ok {
		m.app.Session.Put(r.Context(), "error", "can't get reservation from session")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	// Get room from database
	room, err := m.db.GetRoomByID(res.RoomId)
	if err != nil {
		m.app.Session.Put(r.Context(), "error", "can't find room!")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	res.Room.RoomName = room.RoomName

	// Store reservation back in session
	m.app.Session.Put(r.Context(), "reservation", res)

	// Format dates for template
	sd := res.StartDate.Format("2006-01-02")
	ed := res.EndDate.Format("2006-01-02")

	stringMap := make(map[string]string)
	stringMap["start_date"] = sd
	stringMap["end_date"] = ed

	dataMap := make(map[string]interface{})
	dataMap["reservation"] = res

	render.TemplateCache(w, r, "make-reservation.page.tmpl", &data.TemplateData{
		Form:      forms.New(nil),
		Data:      dataMap,
		StringMap: stringMap,
	})
}

func (m *Repository) PostReservationHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		m.app.Session.Put(r.Context(), "error", "can't parse form!")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	sd := r.Form.Get("start_date")
	ed := r.Form.Get("end_date")

	layout := "2006-01-02"

	startDate, err := time.Parse(layout, sd)
	if err != nil {
		m.app.Session.Put(r.Context(), "error", "can't parse start date")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	endDate, err := time.Parse(layout, ed)
	if err != nil {
		m.app.Session.Put(r.Context(), "error", "can't parse end date")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	roomID, err := strconv.Atoi(r.Form.Get("room_id"))
	if err != nil {
		m.app.Session.Put(r.Context(), "error", "invalid data!")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	reservation := data.Reservation{
		FirstName: r.Form.Get("first_name"),
		LastName:  r.Form.Get("last_name"),
		Phone:     r.Form.Get("phone"),
		Email:     r.Form.Get("email"),
		StartDate: startDate,
		EndDate:   endDate,
		RoomId:    roomID,
	}

	form := forms.New(r.PostForm)

	form.Required("first_name", "last_name", "email")
	form.MinLength("first_name", 3, r)
	form.IsEmail("email")

	if !form.Valid() {
		dataMap := make(map[string]interface{})
		dataMap["reservation"] = reservation
		render.TemplateCache(w, r, "make-reservation.page.tmpl", &data.TemplateData{
			Form: form,
			Data: dataMap,
		})
		return
	}

	newReservationID, err := m.db.InsertReservation(reservation)
	if err != nil {
		m.app.Session.Put(r.Context(), "error", "can't insert reservation into database!")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	restriction := data.RoomRestriction{
		StartDate:     startDate,
		EndDate:       endDate,
		RoomId:        roomID,
		ReservationId: newReservationID,
		RestrictionId: 1,
	}

	err = m.db.InsertRoomRestriction(restriction)
	if err != nil {
		m.app.Session.Put(r.Context(), "error", "can't insert room restriction!")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	htmlMessage := fmt.Sprintf(`
	<strong>Reservation Confirmation</strong><br />
	Dear %s, <br /><br />
	Your reservation has been confirmed for %s to %s.<br /><br />
	`, reservation.FirstName, reservation.StartDate.Format("2006-01-02"), reservation.EndDate.Format("2006-01-02"))

	msg := data.MailData{
		To:       reservation.Email,
		From:     "me@here.com",
		Subject:  "Reservation Confirmation",
		Content:  template.HTML(htmlMessage),
		Template: "dunky.html",
	}
	m.app.MailChan <- msg

	m.app.Session.Put(r.Context(), "reservation", reservation)

	http.Redirect(w, r, "/reservation-summary", http.StatusSeeOther)
}

// ReservationSummary displays the reservation summary page
func (m *Repository) ReservationSummary(w http.ResponseWriter, r *http.Request) {
	reservation, ok := m.app.Session.Get(r.Context(), "reservation").(data.Reservation)
	if !ok {
		m.app.Session.Put(r.Context(), "error", "Can't get reservation from session")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	m.app.Session.Remove(r.Context(), "reservation")

	dataMap := make(map[string]interface{})
	dataMap["reservation"] = reservation

	sd := reservation.StartDate.Format("2006-01-02")
	ed := reservation.EndDate.Format("2006-01-02")
	stringMap := make(map[string]string)
	stringMap["start_date"] = sd
	stringMap["end_date"] = ed

	render.TemplateCache(w, r, "reservation-summary.page.tmpl", &data.TemplateData{
		Data:      dataMap,
		StringMap: stringMap,
	})
}

func (m *Repository) ShowLoginHandler(w http.ResponseWriter, r *http.Request) {
	render.TemplateCache(w, r, "login.page.tmpl", &data.TemplateData{
		Form: forms.New(nil),
		Data: map[string]interface{}{
			"Title": "Login",
		},
	})
}
