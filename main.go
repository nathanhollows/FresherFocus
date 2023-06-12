// FresherFocus is a simple web application that records
// pastoral care interactions with students.

package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/joho/godotenv"
	"github.com/nathanhollows/FresherFocus/filesystem"
	"github.com/nathanhollows/FresherFocus/handler"
)

var router *chi.Mux
var server *http.Server

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	routes()

	server := &http.Server{
		Addr:    os.Getenv("PORT"),
		Handler: router,
	}
	fmt.Println(server.ListenAndServe())
	// systray.Run(onReady, onExit)

	log.Printf("Finished")
}

// func onReady() {
// 	var iconBytes []byte
// 	icon := "static/favicon.ico"
// 	if f, err := os.Open(icon); err == nil {
// 		if fi, err := f.Stat(); err == nil {
// 			size := fi.Size()
// 			iconBytes = make([]byte, size)
// 			f.Read(iconBytes)
// 		}
// 	}

// 	systray.SetIcon(iconBytes)
// 	systray.SetTitle("FresherFocus")
// 	systray.SetTooltip("FresherFocus – Pastoral Care Management")
// 	open := systray.AddMenuItem("Open", "Open FresherFocus")
// 	open.ClickedCh = make(chan struct{})
// 	quit := systray.AddMenuItem("Quit", "Quit FresherFocus")
// 	quit.ClickedCh = make(chan struct{})
// 	http.ListenAndServe(os.Getenv("PORT"), router)
// 	go func() {
// 		<-quit.ClickedCh
// 		systray.Quit()
// 	}()
// 	go func() {
// 		<-open.ClickedCh
// 		openURL("http://localhost" + os.Getenv("PORT"))
// 	}()

// }

// func openURL(url string) {
// 	var err error
// 	switch runtime.GOOS {
// 	case "linux":
// 		err = exec.Command("xdg-open", url).Start()
// 	case "windows":
// 		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
// 	case "darwin":
// 		err = exec.Command("open", url).Start()
// 	default:
// 		err = fmt.Errorf("unsupported platform")
// 	}
// 	if err != nil {
// 		log.Panic(err)
// 	}
// }

// func onExit() {

// }

func routes() {
	router = chi.NewRouter()
	router.Use(middleware.Compress(5))
	router.Use(middleware.CleanPath)
	router.Use(middleware.StripSlashes)
	router.Use(middleware.RedirectSlashes)
	router.Use(middleware.Recoverer)

	router.Get("/", handler.Dashboard)
	router.Get("/dashboard", handler.Dashboard)
	router.Get("/markdown", handler.Markdown)

	router.Route("/cases", func(r chi.Router) {
		r.Get("/", handler.Cases)
		r.Get("/new", handler.NewCase)
		r.Get("/search", handler.SearchCases)
	})

	router.Route("/case", func(r chi.Router) {
		r.Get("/{id}", handler.ViewCase)
		r.Get("/{id}/edit/{fragment}", handler.EditCase)
		r.Post("/{id}/edit/{fragment}", handler.UpdateCaseFragment)
		r.Post("/{id}/togglestatus", handler.ToggleStatus)
		r.Get("/{id}/fragment/{fragment}", handler.ViewCaseFragment)
		r.Get("/{id}/sidebar", handler.CaseSidebar)
	})

	router.Route("/note", func(r chi.Router) {
		r.Get("/add/{source}/{id}", handler.NewNoteForm)
		r.Post("/add/{source}/{id}", handler.SaveNote)
		r.Get("/edit/{id}", handler.EditNote)
		r.Get("/view/{id}", handler.ViewNote)
		r.Delete("/{id}", handler.DeleteNote)
		r.Get("/button/{case}/{id}", handler.NewNoteButton)
		r.Post("/checkbox", handler.Checkbox)
	})

	router.Route("/students", func(r chi.Router) {
		r.Get("/", handler.Students)
		r.Post("/search", handler.StudentSearch)
		r.Get("/search/{term}", handler.Students)
		r.Get("/search/at", handler.AtStudents)
	})

	router.Route("/student", func(r chi.Router) {
		r.Get("/{id}", handler.ViewStudent)
		r.Post("/{id}/tag/{tag}", handler.StudentTag)
		r.Get("/{id}/edit/{fragment}", handler.EditStudentFragment)
		r.Post("/{id}/edit/{fragment}", handler.UpdateStudentFragment)
		r.Get("/{id}/fragment/{fragment}", handler.ViewStudentFragment)
	})

	router.Route("/databases", func(r chi.Router) {
		r.Get("/", handler.Databases)
	})

	router.Post("/upload", handler.Upload)
	router.Post("/upload/csv", handler.UploadCSV)

	router.HandleFunc("/shutdown", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
		go func() {
			if err := server.Shutdown(context.Background()); err != nil {
				log.Fatal(err)
			}
		}()
	})

	router.NotFound(handler.NotFound)

	workDir, _ := os.Getwd()
	filesDir := filesystem.Myfs{Dir: http.Dir(filepath.Join(workDir, "static"))}
	filesystem.FileServer(router, "/static", filesDir)
}
