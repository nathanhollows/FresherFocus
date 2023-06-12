package models

import (
	"log"
	"regexp"
	"strconv"
	"time"
)

type StudentCase struct {
	ID        int `db:"id"`
	StudentID int `db:"student_id"`
	CaseID    int `db:"case_id"`
}

func (r StudentCase) createTable() {
	_, err := db.Exec("CREATE TABLE IF NOT EXISTS student_cases (id INTEGER PRIMARY KEY AUTOINCREMENT, student_id INTEGER, case_id INTEGER, UNIQUE(student_id, case_id));")
	if err != nil {
		log.Fatal(err)
	}
}

func (r StudentCase) MatchCaseAndStudents(caseID int) error {
	type CaseContent struct {
		Content string `db:"content"`
	}
	notes := []CaseContent{}
	err := db.Select(&notes, "SELECT content FROM notes WHERE case_id=? AND deleted_at=?", caseID, time.Time{})
	if err != nil {
		return err
	}

	// Remove any existing student_cases for this case, but roll back if there's an error
	db.Exec("BEGIN")
	_, err = db.Exec("DELETE FROM student_cases WHERE case_id=?", caseID)
	if err != nil {
		return err
	}

	ids := []int{}
	// Find the student ID's from the notes in the format of `#1234567#`
	reg := regexp.MustCompile(`@(\d{7})`)
	for _, note := range notes {
		matches := reg.FindAllStringSubmatch(note.Content, -1)
		for _, match := range matches {
			studentID, err := strconv.Atoi(match[1])
			if err != nil {
				return err
			}
			ids = append(ids, studentID)
		}
	}

	// Remove any duplicates
	for i := 0; i < len(ids); i++ {
		for j := i + 1; j < len(ids); j++ {
			if ids[i] == ids[j] {
				ids = append(ids[:j], ids[j+1:]...)
				j--
			}
		}
	}

	for _, id := range ids {
		_, err = db.Exec("INSERT INTO student_cases (student_id, case_id) VALUES (?, ?)", id, caseID)
		if err != nil {
			return err
		}
	}

	db.Exec("COMMIT")

	return nil
}
