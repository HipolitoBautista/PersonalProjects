package models

import (
	"database/sql"
	"testing"

	_ "github.com/jackc/pgx/v5/stdlib"
)

// This is the DSN to run tests on the server
// var DsnStorage = flag.String("dsn", "host=localhost port=5432 user=hipolito password=mypassword dbname=OSIPPDB_DB_DSN sslmode=disable", "PostgreSQL DSN (Data Source Name)")

func TestInsertCommentsPass(t *testing.T) {
	// Test case 1: Expected to pass
	dsn := *DsnStorage
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		t.Fatalf("Failed to open mock DB connection: %v", err)
	}
	defer db.Close()

	commentsModel := &CommentsModel{DB: db}

	CommentInstance := &Comments{}
	result := commentsModel.InsertComment(CommentInstance)

	if result != nil {
		t.Errorf("TestInsertCommentsPass: Expected Nil, got: %v", result)
	}
}

func TestUpdateCommentsPass(t *testing.T) {
	// Test case 1: Expected to pass
	dsn := *DsnStorage
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		t.Fatalf("Failed to open mock DB connection: %v", err)
	}
	defer db.Close()

	commentsModel := &CommentsModel{DB: db}

	CommentInstance := &Comments{}
	_, result := commentsModel.UpdateComments(CommentInstance)

	if result != nil {
		t.Errorf("TestUpdateCommentsPass: Expected Nil, got: %v", result)
	}
}

func TestExistsFalse(t *testing.T) {
	// Test case 1: Expected to pass
	dsn := *DsnStorage
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		t.Fatalf("Failed to open mock DB connection: %v", err)
	}
	defer db.Close()

	commentsModel := &CommentsModel{DB: db}

	result := commentsModel.Exists(1000)

	if result != false {
		t.Errorf("TestUpdateCommentsPass: Expected Nil, got: %v", result)
	}
}

func TestGetCommentsPass(t *testing.T) {
	// Test case 1: Expected to pass
	dsn := *DsnStorage
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		t.Fatalf("Failed to open mock DB connection: %v", err)
	}
	defer db.Close()

	commentsModel := &CommentsModel{DB: db}
	_, result := commentsModel.GetComments(1000)

	if result != nil {
		t.Errorf("TestUpdateCommentsPass: Expected Nil, got: %v", result)
	}
}
