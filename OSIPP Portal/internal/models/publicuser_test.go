package models

import (
	"database/sql"
	"testing"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func TestInsertUserError(t *testing.T) {
	// Test case 1: Expected to fail
	dsn := *DsnStorage
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		t.Fatalf("Failed to open mock DB connection: %v", err)
	}
	defer db.Close()

	publicuserModel := &PublicUserModel{DB: db}

	result := publicuserModel.Insert("", "", "")

	if result != ErrEmptyFields {
		t.Errorf("TestInsertUserError: Expected ErrEmptyFields, got: %v", result)
	}
}

func TestUserDetailsError(t *testing.T) {
	// Test case 3: Expected to pass
	dsn := *DsnStorage
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		t.Fatalf("Failed to open mock DB connection: %v", err)
	}
	defer db.Close()

	publicuserModel := &PublicUserModel{DB: db}

	_, _, err = publicuserModel.UserData(1000)

	if err != ErrAdminDoesNotExist {
		t.Errorf("TestUserDetailsError: Expected ErrAdminDoesNotExist, got: %v", err)
	}
}

func TestUserAuthenticationErr(t *testing.T) {
	// Test case 3: Expected to pass
	dsn := *DsnStorage
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		t.Fatalf("Failed to open mock DB connection: %v", err)
	}
	defer db.Close()

	publicuserModel := &PublicUserModel{DB: db}

	_, err = publicuserModel.Authenticate("dwadwadawd!!8@gmail.com", "pass")

	if err != ErrInvalidCredentials {
		t.Errorf("TestUserAuthenticationErr: Expected ErrInvalidCredentials, got: %v", err)
	}
}
