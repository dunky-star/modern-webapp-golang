package data

import "encoding/gob"

// Reservation holds reservation data
type Reservation struct {
	FirstName string
	LastName  string
	Email     string
	Phone     string
}

// init registers custom types with gob for session serialization
func init() {
	// Register Reservation type for session storage (scs uses gob encoding)
	gob.Register(Reservation{})
}
