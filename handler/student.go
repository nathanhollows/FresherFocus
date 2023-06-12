package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/go-chi/chi"
	"github.com/nathanhollows/FresherFocus/models"
)

func Students(w http.ResponseWriter, r *http.Request) {
	data := data(r)
	var students models.Students

	term := chi.URLParam(r, "term")
	data["term"] = term
	if term == "" {
		data["empty"] = "Begin typing to search"
	} else {
		params := make(map[string]string)
		params["search"] = term
		students = models.Students{}.Search(params)
	}
	switch len(term) {
	case 0:
		data["empty"] = "Begin typing to search"
	case 1, 2:
		data["empty"] = "Search must be at least 3 characters"
	default:
		data["empty"] = "No results found"
	}

	data["tags"], _ = models.Tags{}.List()
	data["title"] = "Students"
	data["layout"] = "main"
	data["students"] = students
	data["empty"] = "Begin typing to search"
	render(w, data, "student/index", "student/results")
}

func ViewStudent(w http.ResponseWriter, r *http.Request) {
	data := data(r)
	data["layout"] = "main"
	student, _ := models.Students{}.Get(chi.URLParam(r, "id"))
	student.Hydrate()
	data["nontags"] = student.GetNonTags()
	data["student"] = student
	data["title"] = student.Name()
	render(w, data, "student/view")
}

// AtStudents returns a list of students to be auto-completed
func AtStudents(w http.ResponseWriter, r *http.Request) {
	term := r.FormValue("q")
	if len(term) < 3 {
		return
	}
	students := models.Students{}.AtSearch(term)

	type at struct {
		Value string `json:"value"`
		Key   string `json:"key"`
	}

	var ats []at
	for _, s := range students {
		ats = append(ats, at{Value: s.ID, Key: s.Name()})
	}

	// write json
	w.Header().Set("Content-Type", "application/json")

	b, err := json.Marshal(ats)
	if err != nil {
		fmt.Println(err)
		return
	}

	w.Write(b)

}

func StudentSearch(w http.ResponseWriter, r *http.Request) {
	data := data(r)
	term := r.FormValue("search")
	data["search"] = term
	params := make(map[string]string)
	for k, v := range r.Form {
		if k == "tags" {
			// concat tags into a single string
			var tags string
			for _, t := range v {
				tags += t + ","
			}
			params[k] = strings.TrimSuffix(tags, ",")
			continue
		}
		params[k] = v[0]
	}
	students := models.Students{}.Search(params)
	data["students"] = students
	data["empty"] = "Begin typing to search"
	if len(r.FormValue("search")) > 2 {
		data["empty"] = "No results found"
	}
	if len(students) > 0 {
		w.Header().Add("HX-Push", "/students/search/"+term)
	}
	render(w, data, "student/results")
}

// Toggles a student's tag and returns the new list of tags
func StudentTag(w http.ResponseWriter, r *http.Request) {
	data := data(r)
	studentId := chi.URLParam(r, "id")
	tag := chi.URLParam(r, "tag")

	student, _ := models.Students{}.Get(studentId)
	student.Hydrate()
	student.ToggleTag(tag)

	data["student"] = student
	data["nontags"] = student.GetNonTags()
	render(w, data, "student/fragments/tags")
}

var studentFragments = []string{"name", "toc"}

func EditStudentFragment(w http.ResponseWriter, r *http.Request) {
	data := data(r)
	id := chi.URLParam(r, "id")
	fragment := chi.URLParam(r, "fragment")

	// If the fragment is not allowed, return an error
	err := errors.New("invalid fragment")
	for _, f := range studentFragments {
		if f == fragment {
			err = error(nil)
			break
		}
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	s, ok := models.Students{}.Get(id)
	if !ok {
		http.Error(w, "student not found", http.StatusNotFound)
		return
	}
	s.Hydrate()

	data["student"] = s
	render(w, data, "student/edit/"+fragment)
}

// Update a student's information
func UpdateStudentFragment(w http.ResponseWriter, r *http.Request) {

	data := data(r)
	id := chi.URLParam(r, "id")
	fragment := chi.URLParam(r, "fragment")

	// If the fragment is not allowed, return an error
	err := errors.New("invalid fragment")
	for _, f := range studentFragments {
		if f == fragment {
			err = error(nil)
			break
		}
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	s, ok := models.Students{}.Get(id)
	if !ok {
		http.Error(w, "student not found", http.StatusNotFound)
		return
	}
	s.Hydrate()

	err = s.UpdateFragment(fragment, r.FormValue("value"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	data["student"] = s
	render(w, data, "student/fragments/"+fragment)
}

func ViewStudentFragment(w http.ResponseWriter, r *http.Request) {
	data := data(r)
	id := chi.URLParam(r, "id")
	fragment := chi.URLParam(r, "fragment")

	// If the fragment is not allowed, return an error
	err := errors.New("invalid fragment")
	for _, f := range studentFragments {
		if f == fragment {
			err = error(nil)
			break
		}
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	s, ok := models.Students{}.Get(id)
	if !ok {
		http.Error(w, "student not found", http.StatusNotFound)
		return
	}
	s.Hydrate()

	data["student"] = s
	render(w, data, "student/fragments/"+fragment)
}
