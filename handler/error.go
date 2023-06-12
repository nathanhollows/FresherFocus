package handler

import (
	"net/http"
)

func NotFound(w http.ResponseWriter, r *http.Request) {
	data := data(r)
	data["layout"] = "main"
	render(w, data, "errors/404")
}
