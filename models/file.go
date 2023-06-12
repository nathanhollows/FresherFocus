package models

import (
	"log"
	"time"
)

type File struct {
	ID        int
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
	DeletedAt time.Time `db:"deleted_at"`
	Name      string    `db:"name"`
	Path      string    `db:"path"`
	Type      string    `db:"type"`
	Shortcode string    `db:"shortcode"`
	CaseID    string    `db:"case_id"`
	StudentID string    `db:"student_id"`
}

type Files []File

// CreateTable creates the file table using the struct
func (f File) createTable() {
	_, err := db.Exec("CREATE TABLE IF NOT EXISTS files (id INTEGER PRIMARY KEY AUTOINCREMENT, created_at DATETIME, updated_at DATETIME, deleted_at DATETIME, name TEXT, path TEXT, type TEXT, shortcode TEXT, case_id INTEGER, student_id INTEGER)")
	if err != nil {
		log.Fatal(err)
	}
}

// Create inserts a new file into the database
func (f File) Create() {
	_, err := db.Exec("INSERT INTO files (created_at, updated_at, deleted_at, name, path, type, shortcode, case_id, student_id) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)", time.Now(), time.Now(), time.Time{}, f.Name, f.Path, f.Type, f.Shortcode, f.CaseID, f.StudentID)
	if err != nil {
		log.Fatal(err)
	}
}
