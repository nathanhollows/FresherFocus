package handler

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/nathanhollows/FresherFocus/models"
)

func Databases(w http.ResponseWriter, r *http.Request) {
	data := data(r)
	data["layout"] = "main"

	csvs, err := models.CSV{}.List()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	data["csvs"] = csvs

	render(w, data, "databases/index")
}

func ImportDatabase(w http.ResponseWriter, r *http.Request) {
	data := data(r)
	data["layout"] = "main"

	csvID := chi.URLParam(r, "id")
	csv, err := models.CSV{}.Get(csvID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	data["csv"] = csv

	render(w, data, "databases/import")
}
