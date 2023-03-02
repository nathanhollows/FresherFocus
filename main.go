// FresherFocus is a simple web application that records
// pastoral care interactions with students.

package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/joho/godotenv"
	"github.com/nathanhollows/FresherFocus/filesystem"
	"github.com/nathanhollows/FresherFocus/handler"
)

var router *chi.Mux

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

}

func main() {
	routes()
	fmt.Println(http.ListenAndServe(os.Getenv("PORT"), router))
}

func routes() {
	router = chi.NewRouter()
	router.Use(middleware.Recoverer)
	router.Use(middleware.StripSlashes)
	router.Use(middleware.Compress(5))
	router.Use(middleware.Logger)

	router.Get("/", handler.Index)

	router.NotFound(handler.NotFound)

	workDir, _ := os.Getwd()
	filesDir := filesystem.Myfs{Dir: http.Dir(filepath.Join(workDir, "assets"))}
	filesystem.FileServer(router, "/static", filesDir)
}
