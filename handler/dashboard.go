package handler

import (
	"net/http"

	"github.com/nathanhollows/FresherFocus/models"
)

func Dashboard(w http.ResponseWriter, r *http.Request) {
	data := data(r)
	data["layout"] = "main"
	data["cases"], _ = models.Case{}.Search()
	render(w, data, "dashboard/index")
}
