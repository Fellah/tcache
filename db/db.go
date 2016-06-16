package db

import (
	"database/sql"
	"os"
	"time"

	_ "github.com/lib/pq"

	"github.com/fellah/tcache/log"
)

const DB_CONNECTION = "DB_CONNECTION"

var db *sql.DB

func init() {
	db = Connect()
}

func Connect() *sql.DB {
	dbConnection := os.Getenv(DB_CONNECTION)
	if dbConnection == "" {
		log.Error.Fatalln()
	}

	db, err := sql.Open("postgres", dbConnection)
	if err != nil {
		log.Error.Fatalln(err)
	}

	// Config connections.
	db.SetConnMaxLifetime(5 * time.Minute)
	//db.SetMaxIdleConns(20)
	db.SetMaxOpenConns(50)

	return db
}

func Close() {
	if db == nil {
		return
	}

	err := db.Close()
	if err != nil {
		log.Error.Fatal(err)
	}
}
