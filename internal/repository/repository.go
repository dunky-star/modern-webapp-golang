package repository

import (
	"time"

	"github.com/dunky-star/modern-webapp-golang/internal/data"
)

type DatabaseConn interface {
	AllUsers() bool
	InsertReservation(res data.Reservation) (int, error)
	InsertRoomRestriction(r data.RoomRestriction) error
	SearchAvailabilityByDatesByRoomId(start, end time.Time, roomId int) (bool, error)
	SearchAvailabilityForAllRooms(start, end time.Time) ([]data.Room, error)
	GetRoomByID(id int) (data.Room, error)
	GetUserByEmail(email string) (data.User, error)
	UpdateUser(u data.User) error
	Authenticate(email, testPassword string) (int, string, error)
}
