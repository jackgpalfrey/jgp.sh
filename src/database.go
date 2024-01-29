package main

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/lib/pq"
)

func connectToDatabase() *sql.DB {
	user := os.Getenv("POSTGRES_USER")
	pass := os.Getenv("POSTGRES_PASSWORD")
	host := "db"
	// dbname := "urls"
	query := "sslmode=disable"

	connStr := "postgres://" + user + ":" + pass + "@" + host + "?" + query // + "/" + dbname + "?" + query

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS urls (
			short text NOT NULL UNIQUE,
			long text NOT NULL,
			PRIMARY KEY (short)
		)
	`)
	if err != nil {
		log.Fatal(err)
	}

	return db
}
