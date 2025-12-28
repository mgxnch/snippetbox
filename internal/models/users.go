package models

import (
	"database/sql"
	"errors"
	"strings"
	"time"

	"github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
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

// Insert inserts a user into the database. It checks that the email is unique and that
// the password can be converted into a valid bcrypt hash.
func (m *UserModel) Insert(name, email, password string) error {
	// Create bcrupt hash of plaintext password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12) // 2^12 = 4096 iterations
	if err != nil {
		return err
	}

	stmt := `INSERT INTO users (name, email, hashed_password, created)
	VALUES (?, ?, ?, UTC_TIMESTAMP())`

	// Use Exec() method to insert user into users table
	_, err = m.DB.Exec(stmt, name, email, string(hashedPassword))
	if err != nil {
		// We need to already know that MySQL returns error number 1062 (ER_DUP_ENTRY) when
		// the unique key constraint is violated
		var mySQLError *mysql.MySQLError
		if errors.As(err, &mySQLError) {
			// "users_uc_email" is our constraint name
			if mySQLError.Number == 1062 && strings.Contains(mySQLError.Message, "users_uc_email") {
				return ErrDuplicateEmail
			}
		}
		return err
	}
	return nil
}

// Authenticate verifies whether a user with the email and password exists.
// It returns the user's ID if they do. Otherwise, this returns 0 and an error.
func (m *UserModel) Authenticate(email, password string) (int, error) {
	var id int
	var hashedPassword []byte

	// Check if email exists
	stmt := "SELECT id, hashed_password FROM users WHERE email = ?"
	err := m.DB.QueryRow(stmt, email).Scan(&id, &hashedPassword)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, ErrInvalidCredentials
		}
		return 0, err
	}

	// Check if password is correct
	err = bcrypt.CompareHashAndPassword(hashedPassword, []byte(password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return 0, ErrInvalidCredentials
		}
		return 0, err
	}

	// Email and password are correct
	return id, nil
}

// Exists checks if a user with a specific ID exists.
func (m *UserModel) Exists(id int) (bool, error) {
	var exists bool

	stmt := "SELECT EXISTS(SELECT true FROM users WHERE id = ?)"
	err := m.DB.QueryRow(stmt, id).Scan(&exists)
	return exists, err
}
