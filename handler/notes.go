package handler

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/nathanhollows/FresherFocus/models"
)

// Show the form to create a new note
func NewNoteForm(w http.ResponseWriter, r *http.Request) {
	data := data(r)
	data["source"] = chi.URLParam(r, "source")
	data["ID"] = chi.URLParam(r, "id")
	if r.Method == "GET" {
		render(w, data, "notes/new")
	}

}

func NewNoteButton(w http.ResponseWriter, r *http.Request) {
	data := data(r)
	data["source"] = chi.URLParam(r, "source")
	data["ID"] = chi.URLParam(r, "id")
	render(w, data, "notes/button")
}

func ViewNote(w http.ResponseWriter, r *http.Request) {
	data := data(r)
	id := chi.URLParam(r, "id")
	note, err := models.Note{}.Get(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	if note.CaseID != 0 {
		data["source"] = "case"
		data["ID"] = strconv.Itoa(note.CaseID)
	} else if note.StudentID != 0 {
		data["source"] = "student"
		data["ID"] = strconv.Itoa(note.StudentID)
	}

	data["note"] = note
	render(w, data, "notes/view")
}

func EditNote(w http.ResponseWriter, r *http.Request) {
	data := data(r)
	id := chi.URLParam(r, "id")
	note, err := models.Note{}.Get(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	if note.CaseID != 0 {
		data["source"] = "case"
		data["ID"] = strconv.Itoa(note.CaseID)
	} else if note.StudentID != 0 {
		data["source"] = "student"
		data["ID"] = strconv.Itoa(note.StudentID)
	}

	data["note"] = note
	w.Header().Add("HX-Trigger", "notes")
	render(w, data, "notes/edit")
}

func DeleteNote(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	note, err := models.Note{}.Get(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.Header().Add("HX-Trigger", "notes")
	note.Delete()

}

func SaveNote(w http.ResponseWriter, r *http.Request) {
	data := data(r)
	var err error
	note := models.Note{}

	r.ParseForm()

	// If we're editing an existing note, get it
	noteId := r.FormValue("id")
	if noteId != "" {
		note, err = models.Note{}.Get(noteId)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	} else {
		// Otherwise, get the source and id from the URL
		source := chi.URLParam(r, "source")
		id := chi.URLParam(r, "id")
		switch source {
		case "case":
			note.CaseID, err = strconv.Atoi(id)
		case "student":
			note.StudentID, err = strconv.Atoi(id)
		default:
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}

	// Update the note with the new content, or show the form again if there's no content
	content := r.FormValue("content")
	note.Content = content
	if content == "" {
		render(w, data, "notes/new")
	}

	err = note.Save()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Add("HX-Trigger", "notes")
	data["note"] = note
	render(w, data, "notes/view")
}
