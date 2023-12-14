package models

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type admin struct {
	id          int32
	name        string
	email       string
	au_password int32
	created_at  time.Time
}

// written by: Hipolito, Michael, Jahmur, Dennis, Rene
// tested by: Hipolito, Michael, Jahmur, Dennis, Rene
// debugged by: Hipolito, Michael, Jahmur, Dennis, Rene

// setup dependency injection
type AdminModel struct {
	DB *sql.DB //connection pool
}

// Inserts admins into table
func (m *AdminModel) InsertAdmin(name, email, password, auth string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)

	if err != nil {
		return err
	}

	if auth != "1234" {
		return ErrInvalidAuth
	}

	//Query to add admins to table
	query := ` 
					INSERT INTO admin_users(users_name, email, au_password_hash)
					VALUES($1, $2, $3)
		`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	_, err = m.DB.ExecContext(ctx, query, name, email, hashedPassword)
	//Making sure the query ran correct(Error handling )
	fmt.Println(err)
	if err != nil {
		switch {
		case err.Error() == `ERROR: duplicate key value violates unique constraint "admin_users_email_key" (SQLSTATE 23505)`:
			return ErrDuplicateEmail
		default:
			return ErrCouldNotMakeAdmin
		}
	}
	return nil
}

// Authenticates admin with data provided
func (m *AdminModel) AuthenticateAdmin(email, password string) (int, error) {
	var id int
	var hashedPassword []byte
	//Query to get admin ID and password hash from admins table
	query := `
		SELECT id, au_password_hash
		FROM admin_users
		WHERE email = $1
	`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	//Runs query and gets the admin ID and password hash
	err := m.DB.QueryRowContext(ctx, query, email).Scan(&id, &hashedPassword)

	//Error handling
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, ErrInvalidCredentials
		} else {
			return 0, err
		}
	}

	//compares password entered (hashes) and the stored password hash
	err = bcrypt.CompareHashAndPassword(hashedPassword, []byte(password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return 0, ErrInvalidCredentials
		} else {
			return 0, err
		}
	}
	//password is correct
	return id, nil

}

// Getting the details associated with an admin ID (email, username)
func (m *AdminModel) AdminData(ID int) (string, string, error) {
	var userName string
	var email string
	query := `
		SELECT users_name, email
		FROM admin_users
		WHERE ID = $1
	`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	//Executes query
	err := m.DB.QueryRowContext(ctx, query, ID).Scan(&userName, &email)
	//Making sure query ran properly
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", "", ErrAdminDetailsNotFound
		}
	}
	//we got the admin details to display
	return userName, email, nil

}
