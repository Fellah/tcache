package db

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/lib/pq"
)

const (
	DB_CONNECTION = "DB_CONNECTION"
)

var (
	db_connection string
	db            *sql.DB
)

func init() {
	var err error

	db_connection = os.Getenv(DB_CONNECTION)
	if db_connection == "" {
		log.Fatal()
	}

	log.Println(db_connection)

	db, err = sql.Open("postgres", db_connection)
	if err != nil {
		log.Fatal(err)
	}
}

func Close() {
	err := db.Close()
	if err != nil {
		log.Fatal(err)
	}
}
