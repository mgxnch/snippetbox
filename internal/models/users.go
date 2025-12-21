package models

import (
	"database/sql"
	"time"
)

// User hold the data from the users table.
type User struct {
	ID             int
	Name           string
	Email          string
	HashedPassword []byte
	Created        time.Time
}

// UserModel interacts with the database.
type UserModel struct {
	DB *sql.DB
}

// Insert inserts a user into the database.
func (m *UserModel) Insert(name, email, password string) error {
	return nil
}

// Authenticate verifies whether a user with the email and password exists.
// It returns the user's ID if they do. Otherwise, this returns 0 and an error.
func (m *UserModel) Authenticate(email, password string) (int, error) {
	return 0, nil
}

// Exists checks if a user with a specific ID exists.
func (m *UserModel) Exists(id int) (bool, error) {
	return false, nil
}
