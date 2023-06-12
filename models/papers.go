package models

type Paper struct {
	StudentID     int
	Academic_year int
	PaperCode     string
	Semester      string
}

type Papers []Paper
