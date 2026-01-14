package repository

import (
	"time"

	"github.com/dunky-star/modern-webapp-golang/internal/data"
)

type DatabaseConn interface {
	AllUsers() bool
	InsertReservation(res data.Reservation) (int, error)
	InsertRoomRestriction(r data.RoomRestriction) error
	SearchAvailabilityByDates(start, end time.Time, roomId int) (bool, error)
}
