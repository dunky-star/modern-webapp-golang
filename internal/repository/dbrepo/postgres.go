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
func (d *DBConnection) InsertReservation(res data.Reservation) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	defer cancel()

	var newId int

	stmt := `INSERT INTO reservations (first_name, last_name, email, phone, start_date,
	         end_date, room_id, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) returning id`

	err := d.DB.QueryRow(ctx, stmt,
		res.FirstName,
		res.LastName,
		res.Email,
		res.Phone,
		res.StartDate,
		res.EndDate,
		res.RoomId,
		time.Now(),
		time.Now(),
	).Scan(&newId)

	if err != nil {
		d.App.ErrorLog.Println(err)
		d.App.ErrorLog.Println("Error inserting reservation into database")
		return 0, err
	}

	return newId, nil
}

// InsertRoomRestriction inserts a room restriction into the database
func (d *DBConnection) InsertRoomRestriction(r data.RoomRestriction) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := `INSERT INTO room_restrictions (start_date, end_date, room_id, reservation_id,
		     created_at, updated_at, restriction_id) VALUES ($1, $2, $3, $4, $5, $6, $7)`
	_, err := d.DB.Exec(ctx, stmt,
		r.StartDate,
		r.EndDate,
		r.RoomId,
		r.ReservationId,
		time.Now(),
		time.Now(),
		r.RestrictionId,
	)

	if err != nil {
		d.App.ErrorLog.Println(err)
		d.App.ErrorLog.Println("Error inserting room restriction into database")
		return err
	}

	return nil
}

// SearchAvailabilityByDates searches for availability by dates and room id and returns true if available
func (d *DBConnection) SearchAvailabilityByDates(start, end time.Time, roomId int) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `SELECT COUNT(id) FROM room_restrictions WHERE room_id = $1 AND $2 < end_date AND $3 > start_date`

	var numRows int

	row := d.DB.QueryRow(ctx, query, roomId, start, end)

	err := row.Scan(&numRows)
	if err != nil {
		return false, err
	}
	if numRows == 0 {
		return true, nil
	}

	return false, nil
}
