package models

import (
	"database/sql"
	"testing"

	_ "github.com/jackc/pgx/v5/stdlib"
)

// This is the DSN to run tests on the server
// var DsnStorage = flag.String("dsn", "host=localhost port=5432 user=hipolito password=mypassword dbname=OSIPPDB_DB_DSN sslmode=disable", "PostgreSQL DSN (Data Source Name)")

func TestArchiveInsertPass(t *testing.T) {
	// Test case 1: Expected to pass
	dsn := *DsnStorage
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		t.Fatalf("Failed to open mock DB connection: %v", err)
	}
	defer db.Close()

	archiveModel := &ArchiveModel{DB: db}

	result := archiveModel.ArchiveForm(1)

	if result != nil {
		t.Errorf("TestArchiveInsertError: Expected Nil, got: %v", result)
	}
}

func TestUnArchiveError(t *testing.T) {
	// Test case 1: Expected to pass
	dsn := *DsnStorage
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		t.Fatalf("Failed to open mock DB connection: %v", err)
	}
	defer db.Close()

	archiveModel := &ArchiveModel{DB: db}

	result := archiveModel.UnArchiveForm(0)

	if result != ErrFormIDDoesNotExist {
		t.Errorf("TestUnarchiveError: Expected ErrFormIDDoesNotExist, got: %v", result)
	}
}

