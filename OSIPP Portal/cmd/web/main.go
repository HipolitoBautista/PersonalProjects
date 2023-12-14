package main

import (
	"context"
	"crypto/tls"
	"database/sql"
	"flag"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/HipolitoBautista/internal/models"
	"github.com/alexedwards/scs/v2"
	_ "github.com/jackc/pgx/v5/stdlib"
	// written by: Hipolito, Michael, Jahmur, Dennis, Rene
	// tested by: Hipolito, Michael, Jahmur, Dennis, Rene
	// debugged by: Hipolito, Michael, Jahmur, Dennis, Rene
)

type application struct {
	errorLog *log.Logger
	infoLog  *log.Logger

	sessionsManager *scs.SessionManager
	form            models.FormModel
	admin           models.AdminModel
	publicuser      models.PublicUserModel
	archive         models.ArchiveModel
	comments        models.CommentsModel
}

func main() {
	// configure our server
	addr := flag.String("addr", ":4000"+os.Getenv("PORT"), "HTTP network address")
	//This is the dsn to run the code on the server 
	//dsn := flag.String("dsn", "host=localhost port=5432 user=hipolito password=mypassword dbname=OSIPPDB_DB_DSN sslmode=disable", "PostgreSQL DSN (Data Source Name)")
	// This is the dsn to run the code locally
	dsn := flag.String("dsn", os.Getenv("OSIPPDB_DB_DSN"), "PostgreSQL DSN (Data Source Name)")
	flag.Parse()

	// get a database connection pool
	db, err := openDB(*dsn)
	if err != nil {
		log.Print(err)
		return
	}
	//create instances of errorLog and infoLog
	infoLog := log.New(os.Stdout, "INFO/t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR/t", log.Ldate|log.Ltime|log.Lshortfile)
	//setup a new session
	sessionManager := scs.New()
	sessionManager.Lifetime = 1 * time.Hour
	sessionManager.Cookie.Persist = true
	sessionManager.Cookie.Secure = true
	sessionManager.Cookie.SameSite = http.SameSiteLaxMode
	// share data across our handlers
	app := &application{
		errorLog:        errorLog,
		infoLog:         infoLog,
		sessionsManager: sessionManager,

		form:       models.FormModel{DB: db},
		admin:      models.AdminModel{DB: db},
		publicuser: models.PublicUserModel{DB: db},
		archive:    models.ArchiveModel{DB: db},
		comments:   models.CommentsModel{DB: db},
	}
	// cleanup the connection pool
	defer db.Close()
	// acquired a database connection pool
	log.Println("database connection pool established")
	// create and start a custom web server
	log.Printf("starting server on %s", *addr)
	//Tls config
	tlsConfig := &tls.Config{
		CurvePreferences: []tls.CurveID{tls.X25519, tls.CurveP256},
	}
	srv := &http.Server{
		Addr:         *addr,
		Handler:      app.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		TLSConfig:    tlsConfig,
	}
	err = srv.ListenAndServeTLS("./tls/cert.pem", "./tls/key.pem")
	log.Fatal(err)
}

// The openDB() function returns a database connection pool or error
func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}
	// create a context
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	// test the DB connection
	err = db.PingContext(ctx)
	if err != nil {
		return nil, err
	}
	return db, nil
}
