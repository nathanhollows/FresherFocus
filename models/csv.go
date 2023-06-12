package models

import (
	"log"
	"time"
)

// csv table
type CSV struct {
	ID         int
	UploadedAt time.Time
	Name       string
	Path       string
	Imported   time.Time
}

type CSVs []CSV

func (c CSV) createTable() {
	_, err := db.Exec("CREATE TABLE IF NOT EXISTS csv (id INTEGER PRIMARY KEY, uploaded_at DATETIME, name TEXT, path TEXT, imported DATETIME);")
	if err != nil {
		log.Fatal(err)
	}
}

func (c CSV) List() (CSVs, error) {
	rows, err := db.Query("SELECT id, uploaded_at, name, path, imported FROM csv ORDER BY uploaded_at DESC;")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var databases []CSV
	for rows.Next() {
		var d CSV
		err := rows.Scan(&d.ID, &d.UploadedAt, &d.Name, &d.Path, &d.Imported)
		if err != nil {
			return nil, err
		}
		databases = append(databases, d)
	}
	return databases, nil
}

func (c CSV) Get(id string) (CSV, error) {
	var d CSV
	err := db.QueryRow("SELECT id, uploaded_at, name, path, imported FROM csv WHERE id = ?;", id).Scan(&d.ID, &d.UploadedAt, &d.Name, &d.Path, &d.Imported)
	if err != nil {
		return d, err
	}
	return d, nil
}

func (c *CSV) Save() error {
	_, err := db.Exec("INSERT INTO csv (uploaded_at, name, path, imported) VALUES (?, ?, ?, ?)", time.Now(), c.Name, c.Path, c.Imported)
	return err
}

func (c CSV) IsImported() bool {
	return !c.Imported.IsZero()
}
