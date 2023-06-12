package models

import (
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type Note struct {
	ID        int       `db:"id"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
	DeletedAt time.Time `db:"deleted_at"`
	CaseID    int       `db:"case_id"`
	StudentID int       `db:"student_id"`
	Content   string    `db:"content"`
}

type Notes []Note

type TOCEntry struct {
	Title  string
	Anchor string
	Break  bool
}

type TOC []TOCEntry

func (r Note) createTable() {
	_, err := db.Exec("CREATE TABLE IF NOT EXISTS notes (id INTEGER PRIMARY KEY AUTOINCREMENT, created_at DATETIME, updated_at DATETIME, deleted_at DATETIME, content TEXT, case_id INTEGER, student_id INTEGER);")
	if err != nil {
		log.Fatal(err)
	}
}

func (n *Note) Delete() error {
	n.DeletedAt = time.Now()
	_, err := db.NamedExec("UPDATE notes SET deleted_at=:deleted_at WHERE id=:id", n)
	StudentCase{}.MatchCaseAndStudents(n.CaseID)
	return err
}

func (n *Note) Save() error {
	var err error
	if n.CreatedAt.IsZero() {
		n.CreatedAt = time.Now()
		n.UpdatedAt = n.CreatedAt
	} else {
		n.UpdatedAt = time.Now()
	}

	query := "UPDATE notes SET created_at=:created_at, updated_at=:updated_at, deleted_at=:deleted_at, content=:content, case_id=:case_id, student_id=:student_id WHERE id=:id"
	if n.ID == 0 {
		query = "INSERT INTO notes (created_at, updated_at, deleted_at, content, case_id, student_id) VALUES (:created_at, :updated_at, :deleted_at, :content, :case_id, :student_id)"
	}

	res, err := db.NamedExec(query, n)
	if err != nil {
		return err
	}

	var id int64
	id, err = res.LastInsertId()
	if err != nil {
		return err
	}
	if n.ID == 0 {
		n.ID = int(id)
	}

	n.AfterSave()
	return err
}

// Get gets a note from the database
func (n Note) Get(id string) (Note, error) {
	err := db.Get(&n, "SELECT * FROM notes WHERE id=?", id)
	return n, err
}

func (n *Note) AfterSave() {
	// Update the case's updated_at
	c, _ := Case{}.Get(fmt.Sprintf("%d", n.CaseID))
	c.Update()
}

func (n Note) ToggleCheckbox(noteID string, boxID string) {
	n, _ = n.Get(noteID)
	id, _ := strconv.Atoi(boxID)

	v := n.Content

	// Iterate through v character by character and find the nth checkbox
	found := -1
	for i, c := range v {
		if c == '[' && v[i+2] == ']' {
			found++
			if found == id {
				if v[i+1] == 'x' {
					v = v[:i+1] + " " + v[i+2:]
				} else {
					v = v[:i+1] + "x" + v[i+2:]
				}
				break
			}
		}
	}

	n.Content = v
	n.Save()
}

func (n Notes) TOC() (toc TOC) {
	if len(n) == 0 {
		return toc
	}

	// For each note
	for i, note := range n {
		if i > 0 {
			if len(toc) > 0 && !toc[len(toc)-1].Break {
				toc = append(toc, TOCEntry{Break: true})
			}
		}
		// Extract the titles
		reg := regexp.MustCompile(`(?m)^#{1,6} (.*)$`)
		matches := reg.FindAllStringSubmatch(note.Content, -1)
		// For each title...
		if len(matches) > 0 {
			for _, match := range matches {
				reg := regexp.MustCompile(`[^a-zA-Z0-9 ]+`)
				match[0] = reg.ReplaceAllString(match[0], "")
				match[0] = strings.TrimSpace(match[0])
				anchor := strings.ToLower(match[0])
				anchor = strings.ReplaceAll(anchor, " ", "-")
				tocEntry := TOCEntry{
					Anchor: anchor,
					Title:  match[0],
					Break:  false,
				}
				toc = append(toc, tocEntry)
			}
		}
	}

	return toc
}
