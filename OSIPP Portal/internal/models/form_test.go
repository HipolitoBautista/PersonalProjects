package models

import (
	"database/sql"
	"testing"

	_ "github.com/jackc/pgx/v5/stdlib"
)

// This is the DSN to run tests on the server
// var DsnStorage = flag.String("dsn", "host=localhost port=5432 user=hipolito password=mypassword dbname=OSIPPDB_DB_DSN sslmode=disable", "PostgreSQL DSN (Data Source Name)")

func TestInsertFormError(t *testing.T) {
	// Test case 1: Expected to pass
	dsn := *DsnStorage
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		t.Fatalf("Failed to open mock DB connection: %v", err)
	}
	defer db.Close()

	formModel := &FormModel{DB: db}

	FormInstance := &Form{}
	_, result := formModel.Insert(FormInstance)

	if result != ErrCouldNotInsertForm {
		t.Errorf("TestInsertError: Expected ErrCouldNotInsertForm, got: %v", result)
	}
}

func TestDeleteFormError(t *testing.T) {
	// Test case 1: Expected to pass
	dsn := *DsnStorage
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		t.Fatalf("Failed to open mock DB connection: %v", err)
	}
	defer db.Close()

	formModel := &FormModel{DB: db}

	result := formModel.Delete(10000)

	if result != ErrCouldNotDeleteForm {
		t.Errorf("TestDeleteError: Expected ErrCouldNotDeleteForm , got: %v", result)
	}
}
