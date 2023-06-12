package models

import (
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"log"
	"regexp"
	"sort"
	"strings"
	"time"
)

type Student struct {
	AcademicYear         string `db:"academic_year"`
	ID                   string `db:"student_id"`
	Local                bool   `db:"local"`
	College              bool   `db:"college"`
	Admitted             bool   `db:"admitted"`
	SecondSemNotDeclared bool   `db:"second_semester_not_declared"`
	ProfessionalProgram  bool   `db:"professional_program"`
	Distance             bool   `db:"distance"`
	SecondPlusYear       bool   `db:"second_plus_year"`
	SummerSchool         bool   `db:"summer_school_only"`
	Nsn                  string `db:"nsn"`
	Surname              string `db:"surname"`
	First_name           string `db:"first_name"`
	Other_names          string `db:"other_names"`
	Known_name           string `db:"known_name"`
	Title                string `db:"title"`
	Gender               string `db:"gender"`
	Age                  string `db:"age"`
	Residency            string `db:"residency"`
	Citizenship          string `db:"citizenship"`
	Ethnic1              string `db:"ethnic1"`
	Ethnic2              string `db:"ethnic2"`
	Ethnic3              string `db:"ethnic3"`
	European             string `db:"european"`
	Maori                string `db:"maori"`
	Pacific              string `db:"pacific"`
	Asian                string `db:"asian"`
	Melaa                string `db:"melaa"`
	Dis_affect_study     string `db:"dis_affect_study"`
	Home_area_group      string `db:"home_area_group"`
	Home_area            string `db:"home_area"`
	School_id            string `db:"school_id"`
	School_name          string `db:"school_name"`
	Last_year_school     string `db:"last_year_school"`
	Highest_school_qual  string `db:"highest_school_qual"`
	Prior_activity       string `db:"prior_activity"`
	Admission            string `db:"admission"`
	Last_year_foun       string `db:"last_year_foun"`
	Accommodation        string `db:"accommodation"`
	Home_address1        string `db:"home.address1"`
	Home_address2        string `db:"home.address2"`
	Home_address3        string `db:"home.address3"`
	Home_address4        string `db:"home.address4"`
	Home_address5        string `db:"home.address5"`
	Home_postcode        string `db:"home.postcode"`
	Home_country_code    string `db:"home.country_code"`
	Home_country         string `db:"home.country"`
	Study_address1       string `db:"study.address1"`
	Study_address2       string `db:"study.address2"`
	Study_address3       string `db:"study.address3"`
	Study_address4       string `db:"study.address4"`
	Study_address5       string `db:"study.address5"`
	Study_postcode       string `db:"study.postcode"`
	Study_country_code   string `db:"study.country_code"`
	Study_country        string `db:"study_country"`
	Email                string `db:"student_email"`
	Mobile               string `db:"mobile"`
	Is_first_year_uni    string `db:"is_first_year_uni"`
	Prog1                string `db:"prog1"`
	Prog2                string `db:"prog2"`
	Prog3                string `db:"prog3"`
	Prog1_m1             string `db:"prog1_m1"`
	Prog1_m2             string `db:"prog1_m2"`
	Prog2_m1             string `db:"prog2_m1"`
	Prog2_m2             string `db:"prog2_m2"`
	Prog3_m1             string `db:"prog3_m1"`
	Prog3_m2             string `db:"prog3_m2"`
	Papers               string `db:"papers"`
	Declared             string `db:"declared"`
	Efts                 string `db:"efts"`
	S2_start             string `db:"s2_start"`
	CurrentPapers        string `db:"current_papers"`
	PaperCodes           Papers
	Cases                Cases
	Notes                Notes
	Tags                 Tags
}

type Students []Student

func (s Students) Len() int {
	return len(s)
}

func (s Students) Find(id string) *Student {
	for _, student := range s {
		if student.ID == id {
			return &student
		}
	}
	return nil
}

// Sort students by surname
func (s *Students) Sort() {
	sort.Slice(*s, func(i, j int) bool {
		return strings.Compare((*s)[i].Surname, (*s)[j].Surname) < 0
	})
}

// Search students by name or id using partial matching
func (s Students) Search(params map[string]string) Students {
	var results Students

	if len(params["search"]) < 3 {
		return nil
	}

	declared := "N"
	if params["declared"] == "on" {
		declared = "Y"
	}

	search := "%" + params["search"] + "%"

	tags := strings.Split(params["tags"], ",")
	// Remove any non-alphanumeric characters
	nonAlphanumericRegex := regexp.MustCompile(`[^a-zA-Z0-9 ]+`)
	for i, tag := range tags {
		tags[i] = nonAlphanumericRegex.ReplaceAllString(tag, "")
	}

	query := fmt.Sprintf(`
	SELECT DISTINCT students.*
	FROM students
	JOIN tags ON students.student_id = tags.student_id
	WHERE
		local = "true" AND
		(students.student_id LIKE $1
			OR first_name || ' ' || surname LIKE $1
			OR surname || ' ' || first_name LIKE $1
			OR known_name || ' ' || surname LIKE $1)
		AND declared = $2
		AND tags.tag IN ('%s')
		AND tags.deleted_at == $3
	ORDER BY 
	first_name ASC,
	surname ASC
		`, strings.Join(tags, "', '"))

	db.Select(&results, query, search, declared, time.Time{})

	if len(results) == len(s) {
		return nil
	}

	return results
}

func (s Students) AtSearch(term string) Students {
	var results Students

	if len(term) < 3 {
		return nil
	}

	term = term + "%"

	query := `
	SELECT DISTINCT students.*
FROM students
LEFT JOIN tags ON students.student_id = tags.student_id AND tags.tag = 'Leader'
WHERE local = 'true' AND 
(students.student_id LIKE $1
    OR first_name || ' ' || surname LIKE $1
    OR known_name || ' ' || surname LIKE $1)
ORDER BY (CASE WHEN tags.tag IS NOT NULL THEN 0 ELSE 1 END), first_name ASC;
		`

	err := db.Select(&results, query, term)
	if err != nil {
		log.Println(err)
	}

	return results
}

func (s *Student) ParsePaperCodes() {
	var arr []string
	_ = json.Unmarshal([]byte(s.CurrentPapers), &arr)
	// Split text in the format of "HIST108 (S1)" into "HIST108" and "S1"
	for _, paper := range arr {
		paper = strings.TrimSpace(paper)
		paper = strings.TrimSuffix(paper, ")")
		parts := strings.Split(paper, " (")
		s.PaperCodes = append(s.PaperCodes, Paper{PaperCode: parts[0], Semester: parts[1]})
	}
	// Sort papers by semester and then paper code
	sort.Slice(s.PaperCodes, func(i, j int) bool {
		if s.PaperCodes[i].Semester == s.PaperCodes[j].Semester {
			return s.PaperCodes[i].PaperCode < s.PaperCodes[j].PaperCode
		}
		return s.PaperCodes[i].Semester < s.PaperCodes[j].Semester
	})

}

func (s Students) Get(id string) (Student, bool) {
	var student Student

	err := db.Get(&student, "SELECT * FROM students WHERE student_id = $1", id)
	if err != nil {
		return Student{}, false
	}

	return student, true
}

func (s *Student) Hydrate() {
	s.GetTags()
	s.GetNotes()
	s.GetCases()
	s.ParsePaperCodes()
}

func (s *Student) GetTags() {
	var tags Tags
	_ = db.Select(&tags, "SELECT * FROM tags WHERE student_id = $1 AND deleted_at == $2", s.ID, time.Time{})
	s.Tags = tags
}

func (s *Student) GetNonTags() []string {
	s.Tags.TagNames()
	tags, _ := Tags{}.List(s.Tags.TagNames()...)
	return tags
}

func (s *Student) ToggleTag(tag string) {
	// Regex to check if tag is alphanumeric
	nonAlphanumericRegex := regexp.MustCompile(`[^a-zA-Z0-9]+`)
	if nonAlphanumericRegex.MatchString(tag) {
		return
	}
	for _, t := range s.Tags {
		if t.Tag == tag {
			Tags{}.Delete(s.ID, tag)
			s.GetTags()
			return
		}
	}

	Tags{}.Create(s.ID, tag)
	s.GetTags()
}

func (s *Student) GetNotes() {
	var notes Notes
	_ = db.Select(&notes, "SELECT * FROM notes WHERE student_id = $1 AND deleted_at == $2", s.ID, time.Time{})
	s.Notes = notes
}

func (s *Student) GetCases() {
	var cases Cases
	db.Select(&cases, "SELECT c.* FROM cases c JOIN student_cases s ON s.case_id == c.id WHERE s.student_id = $1 AND c.deleted_at == $2 ORDER BY CASE WHEN closed_at==? then 0 else 1 end, created_at DESC", s.ID, time.Time{}, time.Time{})
	s.Cases = cases
}

func (s Student) Name() string {
	if s.Known_name != "" && s.Known_name != s.First_name {
		return s.First_name + " (" + s.Known_name + ") " + s.Surname
	}
	return s.First_name + " " + s.Surname
}

func (s Student) KnownAs() string {
	if strings.EqualFold(s.Known_name, s.First_name) {
		return ""
	}
	return s.Known_name
}

func (s *Student) Badge() template.HTML {
	badge := `<a href="/student/%s" class="badge %s">%s</a>`
	leader := false

	if s.Tags == nil {
		s.GetTags()
	}

	for _, tag := range s.Tags {
		if tag.Tag == "Leader" {
			leader = true
			break
		}
	}

	if leader {
		badge = fmt.Sprintf(badge, s.ID, "leader", s.ID)
	} else {
		badge = fmt.Sprintf(badge, s.ID, "secondary", s.ID)
	}

	return template.HTML(badge)
}

// Get all students from the students table
func GetAllStudents() (Students, error) {
	var students Students

	err := db.Select(&students, "SELECT * FROM students WHERE local = 'true'")
	if err != nil {
		return nil, err
	}

	return students, nil
}

func (s *Student) UpdateFragment(fragment string, value string) error {
	switch fragment {
	case "name":
		s.Known_name = value
	default:
		return errors.New("invalid fragment")
	}
	_, err := db.NamedExec("UPDATE students SET known_name=:known_name WHERE student_id=:student_id", s)
	return err
}

func (s Student) TOC() TOC {
	return s.Notes.TOC()
}
