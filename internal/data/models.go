package data

import (
	"encoding/gob"
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

// init registers custom types with gob for session serialization
func init() {
	// Register Reservation type for session storage (scs uses gob encoding)
	gob.Register(Reservation{})
}
