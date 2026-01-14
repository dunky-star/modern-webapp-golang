package repository

import "github.com/dunky-star/modern-webapp-golang/internal/data"

type DatabaseConn interface {
	AllUsers() bool
	InsertReservation(res data.Reservation) (int, error)
	InsertRoomRestriction(r data.RoomRestriction) error
}
