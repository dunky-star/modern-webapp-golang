package dbrepo

import (
	"context"
	"time"

	"github.com/dunky-star/modern-webapp-golang/internal/data"
)

func (d *DBConnection) AllUsers() bool {
	return true
}

// InsertReservation inserts a reservation into the database
func (d *DBConnection) InsertReservation(res data.Reservation) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	defer cancel()

	stmt := `INSERT INTO reservations (first_name, last_name, email, phone, start_date,
	 end_date, room_id, create_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`

	_, err := d.DB.Exec(ctx, stmt,
		res.FirstName,
		res.LastName,
		res.Email,
		res.Phone,
		res.StartDate,
		res.EndDate,
		res.RoomId,
		time.Now(),
		time.Now(),
	)
	if err != nil {
		d.App.ErrorLog.Println(err)
		d.App.ErrorLog.Println("Error inserting reservation into database")
		return err
	}

	return nil
}
