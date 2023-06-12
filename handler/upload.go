package handler

import (
	"io"
	"net/http"
	"os"
	"time"

	"github.com/nathanhollows/FresherFocus/models"
)

func Upload(w http.ResponseWriter, r *http.Request) {
	data := data(r)
	data["layout"] = "main"

	r.ParseMultipartForm(10 << 20)

	path := "student"
	if r.Form.Has("caseID") {
		path = "case/" + r.Form.Get("caseID")
	} else if r.Form.Has("studentID") {
		path = "student/" + r.Form.Get("studentID")
	}

	// Save each file to static/files/
	for _, v := range r.MultipartForm.File {
		for _, h := range v {
			file, err := h.Open()
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			defer file.Close()

			// Create directory if it doesn't exist
			if _, err := os.Stat("static/files/" + path); os.IsNotExist(err) {
				os.MkdirAll("static/files/"+path, 0755)
			}

			// Create file
			dst, err := os.Create("static/files/" + path + "/" + h.Filename)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			defer dst.Close()

			// Copy file
			if _, err := io.Copy(dst, file); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}

			fileRecord := models.File{
				Name:      h.Filename,
				Path:      path + "/" + h.Filename,
				Type:      h.Header.Get("Content-Type"),
				CaseID:    r.Form.Get("caseID"),
				StudentID: r.Form.Get("studentID"),
			}
			fileRecord.Create()
		}
	}
	w.Header().Add("HX-Trigger", "upload")
}

func UploadCSV(w http.ResponseWriter, r *http.Request) {
	data := data(r)
	data["layout"] = "htmx"

	r.ParseMultipartForm(10 << 20)

	file, handler, err := r.FormFile("file")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	defer file.Close()

	// Save the file to the server
	f, err := os.OpenFile("./static/databases/"+handler.Filename, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	defer f.Close()
	io.Copy(f, file)

	record := models.CSV{
		Name:     handler.Filename,
		Path:     "./static/databases/" + handler.Filename,
		Imported: time.Time{},
	}
	err = record.Save()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	data["csvs"], err = models.CSV{}.List()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	render(w, data, "databases/list_csvs")
}
