package db

import (
	"database/sql"
	"os"

	_ "github.com/lib/pq"

	"github.com/fellah/tcache/log"
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
		log.Error.Fatalln()
	}

	db, err = sql.Open("postgres", db_connection)
	if err != nil {
		log.Error.Fatalln(err)
	}
}

func Close() {
	err := db.Close()
	if err != nil {
		log.Error.Fatal(err)
	}
}
