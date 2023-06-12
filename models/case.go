package models

import (
	"errors"
	"log"
	"time"
)

// Case is a struct that represents notes and files for a student
type Case struct {
	ID        int       `db:"id"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
	DeletedAt time.Time `db:"deleted_at"`
	ClosedAt  time.Time `db:"closed_at"`
	Title     string    `db:"title"`
	Summary   string    `db:"summary"`
	Files     Files
	Notes     Notes
	Students  Students
}

type Cases []Case

func (c Case) clearEmpty() error {
	_, err := db.Exec("UPDATE cases SET deleted_at=? WHERE title=='' AND summary=='' AND deleted_at==?", time.Now(), time.Time{})
	return err
}

func (c Case) Status() string {
	if c.ClosedAt.IsZero() {
		return "Open"
	}
	return "Closed"
}

func (c *Case) ToggleStatus() {
	if c.ID == 0 {
		return
	}
	if c.ClosedAt.IsZero() {
		c.ClosedAt = time.Now()
	} else {
		c.ClosedAt = time.Time{}
	}
	c.Update()
}

// Save saves a case to the database
func (c *Case) Create() error {
	c.CreatedAt = time.Now()
	c.UpdatedAt = time.Now()
	res, err := db.NamedExec("INSERT INTO cases (created_at, updated_at, deleted_at, closed_at, title, summary) VALUES (:created_at, :updated_at, :deleted_at, :closed_at, :title, :summary)", c)
	if err != nil {
		return err
	}
	id, _ := res.LastInsertId()
	c.ID = int(id)

	// Delete the case if it's empty after 5 minutes
	go func() {
		for {
			time.Sleep(5 * time.Minute)
			Case{}.clearEmpty()
		}
	}()

	return err
}

// Update updates a case in the database
func (c *Case) Update() error {
	c.UpdatedAt = time.Now()
	c.DeletedAt = time.Time{}
	_, err := db.NamedExec("UPDATE cases SET updated_at=:updated_at, deleted_at=:deleted_at, closed_at=:closed_at, title=:title, summary=:summary WHERE id=:id", c)
	StudentCase{}.MatchCaseAndStudents(c.ID)
	return err
}

// Delete deletes a case from the database
func (c *Case) Delete() error {
	c.DeletedAt = time.Now()
	_, err := db.NamedExec("UPDATE cases SET deleted_at=:deleted_at WHERE id=:id", c)
	return err
}

// Search searches for cases in the database
func (c Case) Search() ([]Case, error) {
	var cases []Case
	err := db.Select(&cases, "SELECT * FROM cases WHERE deleted_at==? ORDER BY CASE WHEN closed_at==? then 0 else 1 end, created_at DESC", time.Time{}, time.Time{})
	for i := range cases {
		err = cases[i].GetNotes()
		if err != nil {
			continue
		}
		err = cases[i].GetStudents()
		if err != nil {
			continue
		}
	}
	return cases, err
}

// Search searches for cases in the database
func (c Case) Active() Cases {
	var cases Cases
	err := db.Select(&cases, "SELECT * FROM cases WHERE deleted_at==? AND closed_at==? ORDER BY created_at DESC", time.Time{}, time.Time{})
	for i := range cases {
		err = cases[i].GetNotes()
		if err != nil {
			continue
		}
		err = cases[i].GetStudents()
		if err != nil {
			continue
		}
	}
	return cases
}

// Get gets a case from the database
func (c Case) Get(id string) (Case, error) {
	err := db.Get(&c, "SELECT * FROM cases WHERE id=?", id)
	if err != nil {
		return Case{}, err
	}
	err = c.GetNotes()
	if err != nil {
		return Case{}, err
	}
	err = c.GetStudents()
	if err != nil {
		return Case{}, err
	}
	err = c.GetFiles()
	if err != nil {
		return Case{}, err
	}

	return c, err
}

func (c *Case) GetNotes() error {
	err := db.Select(&c.Notes, "SELECT * FROM notes WHERE case_id=? and deleted_at==?", c.ID, time.Time{})
	return err
}

func (c *Case) GetStudents() error {
	err := db.Select(&c.Students, "SELECT * FROM students WHERE student_id IN (SELECT student_id FROM student_cases WHERE case_id = ?)", c.ID)
	return err
}

func (c *Case) GetFiles() error {
	err := db.Select(&c.Files, "SELECT * FROM files WHERE case_id=? and deleted_at==?", c.ID, time.Time{})
	return err
}

// CreateCaseTable creates the case table using the struct
func (c Case) createTable() {
	// Check if table exists
	_, err := db.Exec("CREATE TABLE IF NOT EXISTS cases (id INTEGER PRIMARY KEY AUTOINCREMENT , created_at DATETIME, updated_at DATETIME, deleted_at DATETIME, closed_at DATETIME, title TEXT, summary TEXT);")
	if err != nil {
		log.Fatal(err)
	}
}

func (c *Case) UpdateCaseFragment(fragment string, value string) error {
	switch fragment {
	case "title":
		c.Title = value
	case "summary":
		c.Summary = value
	default:
		return errors.New("invalid fragment")
	}
	return c.Update()
}

func (c Case) TOC() TOC {
	return c.Notes.TOC()
}
