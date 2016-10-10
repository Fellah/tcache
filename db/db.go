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


func sendQueryParams(query string, params ...interface{}) error {
	CheckConnect()

	txn, err := startTransaction()

	stmt, err := txn.Prepare(query)
	if err != nil {
		return err
	}

	_, err = stmt.Exec(params...)
	if err != nil {
		return err
	}

	if err = stmt.Close(); err != nil {
		return err
	}

	commitTransaction(txn)

	return nil
}

func startTransaction() (*sql.Tx, error) {
	CheckConnect()

	return db.Begin()
}

func sendQueryParamsRaw(txn *sql.Tx, query string, params ...interface{}) error {
	stmt, err := txn.Prepare(query)
	if err != nil {
		return err
	}

	_, err = stmt.Exec(params...)
	if err != nil {
		return err
	}

	if err = stmt.Close(); err != nil {
		return err
	}

	return nil
}

func commitTransaction(txn *sql.Tx) error {
	if err := txn.Commit(); err != nil {
		return err
	}

	return nil
}

func CheckConnect() {
	err := db.Ping()

	if err != nil {
		db = Connect()
	}
}
