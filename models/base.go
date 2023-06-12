package models

import (
	"log"

	"github.com/jmoiron/sqlx"
)

var db *sqlx.DB

func init() {
	var err error
	db, err = sqlx.Open("sqlite3", "fresher.db")
	if err != nil {
		log.Fatal(err)
	}

	// Create the databases
	Case{}.createTable()
	Case{}.clearEmpty()
	File{}.createTable()
	Note{}.createTable()
	CSV{}.createTable()
	StudentCase{}.createTable()
	Tag{}.createTable()
}
