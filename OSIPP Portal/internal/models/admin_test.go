package models

import (
	"database/sql"
	"flag"
	"os"
	"testing"

	_ "github.com/jackc/pgx/v5/stdlib"
)

// This is the DSN to run tests on the server
// var DsnStorage = flag.String("dsn", "host=localhost port=5432 user=hipolito password=mypassword dbname=OSIPPDB_DB_DSN sslmode=disable", "PostgreSQL DSN (Data Source Name)")

// This is the DSN to run tests locally
var DsnStorage = flag.String("dsn", os.Getenv("OSIPPDB_DB_DSN"), "PostgreSQL DSN (Data Source Name)")

func TestInsertAdminError(t *testing.T) {
	// Test case 1: Expected to fail
	dsn := *DsnStorage
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		t.Fatalf("Failed to open mock DB connection: %v", err)
	}
	defer db.Close()

	adminModel := &AdminModel{DB: db}

	result := adminModel.InsertAdmin("Jamar", "Kukul", "password", "")

	if result != ErrInvalidAuth {
		t.Errorf("TestInsertAdminError: Expected ErrInvalidAuth, got: %v", result)
	}
}

func TestInsertAdminSemiAuthError(t *testing.T) {
	// Test case 2: Expected to pass
	dsn := *DsnStorage
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		t.Fatalf("Failed to open mock DB connection: %v", err)
	}
	defer db.Close()

	adminModel := &AdminModel{DB: db}

	result := adminModel.InsertAdmin("Jamar", "Kukul", "password", "12")

	if result != ErrInvalidAuth {
		t.Errorf("TestInsertAdminSemiAuthError: Expected ErrInvalidAuth, got: %v", result)
	}
}

func TestAdminDetailsError(t *testing.T) {
	// Test case 3: Expected to pass
	dsn := *DsnStorage
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		t.Fatalf("Failed to open mock DB connection: %v", err)
	}
	defer db.Close()

	adminModel := &AdminModel{DB: db}

	_, _, err = adminModel.AdminData(1000)

	if err != ErrAdminDetailsNotFound {
		t.Errorf("TestAuthnticateAdminError: Expected ErrAdminDetailsNotFound, got: %v", err)
	}
}
