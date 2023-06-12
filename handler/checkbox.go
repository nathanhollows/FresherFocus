package handler

import (
	"net/http"

	"github.com/nathanhollows/FresherFocus/models"
)

func Checkbox(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	noteID := r.FormValue("note")
	boxID := r.FormValue("box")

	models.Note{}.ToggleCheckbox(noteID, boxID)

}
