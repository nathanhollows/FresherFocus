package handler

import (
	"errors"
	"net/http"
	"strings"

	"github.com/go-chi/chi"
	"github.com/nathanhollows/FresherFocus/models"
)

func Cases(w http.ResponseWriter, r *http.Request) {
	data := data(r)
	data["layout"] = "main"
	data["title"] = "Cases"
	data["activeCount"] = len(models.Case{}.Active())
	data["cases"], _ = models.Case{}.Search()
	render(w, data, "cases/index")
}

// A list of all the case fragments that can be updated
var caseFragments = []string{"title", "summary"}

func SearchCases(w http.ResponseWriter, r *http.Request) {
	data := data(r)
	cases, _ := models.Case{}.Search()
	data["cases"] = cases
	render(w, data, "cases/fragments/results")
}

func NewCase(w http.ResponseWriter, r *http.Request) {
	data := data(r)

	c := models.Case{}
	c.Create()

	data["case"] = c
	data["title"] = "New Case"
	data["layout"] = "main"
	render(w, data, "cases/view", "cases/edit/title", "cases/edit/summary", "cases/fragments/sidebar")
}

func ViewCase(w http.ResponseWriter, r *http.Request) {
	data := data(r)
	c, err := models.Case{}.Get(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	if c.ID == 0 {
		http.Redirect(w, r, "/cases", http.StatusFound)
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	data["case"] = c
	data["title"] = c.Title
	data["layout"] = "main"
	render(w, data, "cases/view", "cases/edit/title", "cases/edit/summary", "cases/fragments/sidebar")
}

func EditCase(w http.ResponseWriter, r *http.Request) {
	data := data(r)
	id := chi.URLParam(r, "id")
	fragment := chi.URLParam(r, "fragment")

	// If the fragment is not allowed, return an error
	err := errors.New("invalid fragment")
	for _, f := range caseFragments {
		if f == fragment {
			err = error(nil)
			break
		}
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	c, err := models.Case{}.Get(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	data["case"] = c
	render(w, data, "cases/edit/"+fragment)
}

func ViewCaseFragment(w http.ResponseWriter, r *http.Request) {
	data := data(r)
	id := chi.URLParam(r, "id")
	fragment := chi.URLParam(r, "fragment")

	// If the fragment is not allowed, return an error
	err := errors.New("invalid fragment")
	for _, f := range caseFragments {
		if f == fragment {
			err = error(nil)
			break
		}
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	c, err := models.Case{}.Get(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	data["case"] = c
	render(w, data, "cases/fragments/"+fragment)
}

func UpdateCaseFragment(w http.ResponseWriter, r *http.Request) {
	data := data(r)
	id := chi.URLParam(r, "id")
	fragment := chi.URLParam(r, "fragment")

	// If the fragment is not allowed, return an error
	err := errors.New("invalid fragment")
	for _, f := range caseFragments {
		if f == fragment {
			err = error(nil)
			break
		}
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	c, err := models.Case{}.Get(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	c.UpdateCaseFragment(fragment, r.FormValue("value"))

	if strings.Contains(r.Header.Get("HX-Current-URL"), "new") {
		w.Header().Add("HX-Push", "/case/"+id)
	}
	w.Header().Add("HX-Trigger", "sidebar")
	data["case"] = c
	render(w, data, "cases/fragments/"+fragment)
}

func ToggleStatus(w http.ResponseWriter, r *http.Request) {
	data := data(r)
	c, err := models.Case{}.Get(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	c.ToggleStatus()
	data["case"] = c
	render(w, data, "cases/fragments/status")
}

func CaseSidebar(w http.ResponseWriter, r *http.Request) {
	data := data(r)
	c, err := models.Case{}.Get(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	data["case"] = c
	render(w, data, "cases/fragments/sidebar")
}
