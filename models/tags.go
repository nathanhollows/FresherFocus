package models

import (
	"log"
	"strings"
	"time"
)

type Tag struct {
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
	DeletedAt time.Time `db:"deleted_at"`
	StudentID int       `db:"student_id"`
	Tag       string    `db:"tag"`
}

type Tags []Tag

// Create table
func (t Tag) createTable() {
	_, err := db.Exec("CREATE TABLE IF NOT EXISTS tags (created_at DATETIME, updated_at DATETIME, deleted_at DATETIME, student_id INTEGER, tag TEXT, PRIMARY KEY (student_id, tag));")
	if err != nil {
		log.Fatal(err)
	}
	t.Populate()
}

// Populate table using academic_year and student_id from students table and time.Now()
func (t Tag) Populate() {
	_, err := db.Exec(`INSERT OR IGNORE INTO tags (created_at, updated_at, deleted_at, student_id, tag) SELECT $1 as created_at, $1 as updated_at, $2 as deleted_at, student_id, academic_year FROM students;`, time.Now(), time.Time{}, time.Now())
	if err != nil {
		log.Fatal(err)
	}
}

// Create tag
func (t Tags) Create(studentID string, tag string) {
	_, err := db.Exec("INSERT INTO tags (created_at, updated_at, deleted_at, student_id, tag) VALUES ($1, $1, $2, $3, $4);", time.Now(), time.Time{}, studentID, tag)
	if err != nil {
		db.Exec("UPDATE tags SET deleted_at = $1 WHERE student_id = $2 AND tag = $3;", time.Time{}, studentID, tag)
	}
}

// Delete tag
func (t Tags) Delete(studentID string, tag string) {
	_, err := db.Exec("UPDATE tags SET deleted_at = $1 WHERE student_id = $2 AND tag = $3;", time.Now(), studentID, tag)
	if err != nil {
		log.Fatal(err)
	}
}

// Get Unique Tags
func (t Tags) List(exclude ...string) (tags []string, err error) {
	excludeTags := strings.Join(exclude, "','")
	excludeTags = "'" + excludeTags + "'"
	err = db.Select(&tags, "SELECT DISTINCT tag FROM tags WHERE tag NOT IN ('academic_year', "+excludeTags+") ORDER BY tag DESC;", excludeTags)
	if err != nil {
		return nil, err
	}

	return tags, nil
}

func (t Tags) TagNames() []string {
	var tags []string
	for _, tag := range t {
		tags = append(tags, tag.Tag)
	}
	return tags
}
