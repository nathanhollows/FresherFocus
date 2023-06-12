package handler

import (
	"net/http"
)

func Markdown(w http.ResponseWriter, r *http.Request) {
	data := data(r)
	data["title"] = "Markdown"
	data["layout"] = "main"
	render(w, data, "help/markdown")
}
