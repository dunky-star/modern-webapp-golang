package data

import (
	"encoding/gob"
	"html/template"
	"time"
)

type User struct {
	Id          int       `json:"id"`
	FirstName   string    `json:"first_name"`
	LastName    string    `json:"last_name"`
	Email       string    `json:"email"`
	Password    string    `json:"password"`
	AccessLevel int       `json:"access_level"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type Room struct {
	Id        int       `json:"id"`
	RoomName  string    `json:"room_name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Restriction struct {
	Id              int       `json:"id"`
	RestrictionName string    `json:"restriction_name"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

type Reservation struct {
	Id        int       `json:"id"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Email     string    `json:"email"`
	Phone     string    `json:"phone"`
	StartDate time.Time `json:"start_date"`
	EndDate   time.Time `json:"end_date"`
	RoomId    int       `json:"room_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Room      Room      `json:"room"`
}

type RoomRestriction struct {
	Id            int         `json:"id"`
	StartDate     time.Time   `json:"start_date"`
	EndDate       time.Time   `json:"end_date"`
	RoomId        int         `json:"room_id"`
	ReservationId int         `json:"reservation_id"`
	RestrictionId int         `json:"restriction_id"`
	CreatedAt     time.Time   `json:"created_at"`
	UpdatedAt     time.Time   `json:"updated_at"`
	Room          Room        `json:"room"`
	Reservation   Reservation `json:"reservation"`
	Restriction   Restriction `json:"restriction"`
}

// MailData holds an email message
type MailData struct {
	To       string
	From     string
	Subject  string
	Content  template.HTML
	Template string
}

// init registers custom types with gob for session serialization
func init() {
	// Register types that are actually stored in sessions
	gob.Register(Reservation{})
	gob.Register(Room{})
	gob.Register([]Room{})
	gob.Register(time.Time{})
}
