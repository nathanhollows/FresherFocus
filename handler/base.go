package handler

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
	"github.com/nathanhollows/FresherFocus/models"
)

// Error represents a handler error. It provides methods for a HTTP status
// code and embeds the built-in error interface.
type Error interface {
	error
	Status() int
}

// StatusError represents an error with an associated HTTP status code.
type StatusError struct {
	Code int
	Err  error
}

// Allows StatusError to satisfy the error interface.
func (se StatusError) Error() string {
	return se.Err.Error()
}

// Returns our HTTP status code.
func (se StatusError) Status() int {
	return se.Code
}

// The Handler struct that takes a configured Env and a function matching
// our useful signature.
type Handler struct {
	H func(w http.ResponseWriter, r *http.Request) error
}

// ServeHTTP allows our Handler type to satisfy http.Handler.
func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	err := h.H(w, r)
	if err != nil {
		switch e := err.(type) {
		case Error:
			// We can retrieve the status here and write out a specific
			// HTTP status code.
			log.Printf("HTTP %d - %s", e.Status(), e)
			http.Error(w, e.Error(), e.Status())
		default:
			// Any error types we don't specifically look out for default
			// to serving a HTTP 500
			http.Error(w, http.StatusText(http.StatusInternalServerError),
				http.StatusInternalServerError)
		}
	}
}

func data(r *http.Request) map[string]interface{} {
	data := make(map[string]interface{})
	data["title"] = os.Getenv("APP_NAME")
	data["layout"] = "htmx"
	data["hxrequest"] = r.Header.Get("HX-Request") == "true"
	return data
}

func render(w http.ResponseWriter, data map[string]interface{}, patterns ...string) error {
	w.Header().Set("Content-Type", "text/html")
	err := parse(data, patterns...).ExecuteTemplate(w, "base", data)
	if err != nil {
		http.Error(w, err.Error(), 0)
		log.Print("Template executing error: ", err)
	}
	return err
}

func parse(data map[string]interface{}, patterns ...string) *template.Template {
	if data["title"] != os.Getenv("APP_NAME") {
		data["title"] = fmt.Sprintf("%s | %s", data["title"], os.Getenv("APP_NAME"))
	}
	for i := 0; i < len(patterns); i++ {
		patterns[i] = "html/templates/" + patterns[i] + ".html"
	}
	patterns = append(patterns, "html/layouts/"+data["layout"].(string)+".html")
	patterns = append(patterns, "html/partials/flash.html")
	patterns = append(patterns, "html/partials/notes.html")
	return template.Must(template.New("base").Funcs(funcs).ParseFiles(patterns...))
}

var funcs = template.FuncMap{
	"upper": func(v string) string {
		return strings.ToUpper(v)
	},
	"lower": func(v string) string {
		return strings.ToLower(v)
	},
	"date": func(t time.Time) string {
		if t.Year() == time.Now().Year() {
			return t.Format("2 January")
		}
		return t.Format("2 January 2006")
	},
	"time": func(t time.Time) string {
		return t.Format("15:04")
	},
	"divide": func(a, b int) float32 {
		if a == 0 || b == 0 {
			return 0
		}
		return float32(a) / float32(b)
	},
	"progress": func(a, b int) float32 {
		if a == 0 || b == 0 {
			return 0
		}
		return float32(a) / float32(b) * 100
	},
	"add": func(a, b int) int {
		return a + b
	},
	"year": func() string {
		return time.Now().Format("2006")
	},
	"static": func(filename string) string {
		filename = strings.TrimPrefix(filename, "/")
		// get last modified time
		file, err := os.Stat("static/" + filename)

		if err != nil {
			return "/static/" + filename
		}

		modifiedtime := file.ModTime()
		return "/static/" + filename + "?v=" + modifiedtime.Format("20060102150405")
	},
	"md": func(v string) template.HTML {
		extensions := parser.CommonExtensions | parser.AutoHeadingIDs | parser.NoEmptyLineBeforeBlock | parser.HardLineBreak

		// Find and capture student ID in string in the form of `@1234567`
		// and replace with a link to the student's profile
		reg := regexp.MustCompile(`@(\d{7})`)
		v = reg.ReplaceAllStringFunc(v, func(s string) string {
			id := strings.Trim(s, "@")
			student, ok := models.Students{}.Get(id)
			if !ok {
				return fmt.Sprintf(`<a href="/student/%s" class="badge">%s</a>`, id, id)
			}
			student.Hydrate()
			return student.Name() + " " + string(student.Badge())
		})

		// Find and capture case ID in string in the form of `#1234567`
		// and replace with a link to the case
		reg = regexp.MustCompile(`#(\d+)#`)
		v = reg.ReplaceAllStringFunc(v, func(s string) string {
			id := strings.Trim(s, "#")
			if id == "" {
				return "#"
			}
			c, err := models.Case{}.Get(id)
			if err != nil {
				return fmt.Sprintf(`<a href="/case/%s" class="secondary badge">Case #%s</a>`, id, id)
			}
			return fmt.Sprintf(`<a href="/case/%s" class="secondary badge">Case #%s</a> <em>%s</em>`, id, id, c.Title)
		})

		// Find and iterate through all instances of '[ ]' and replace with '<input type="checkbox" disabled>' and '[x]' with '<input type="checkbox" disabled checked>' with a value equal to the number of times the string is found
		reg = regexp.MustCompile(`\[(x| )\]`)
		matches := reg.FindAllStringSubmatch(v, -1)
		for i, match := range matches {
			if match[1] == "x" {
				v = strings.Replace(v, match[0], fmt.Sprintf(`<input type="checkbox" checked hx-post="/note/checkbox" hx-vars="box:%d" hx-swap="none">`, i), 1)
			} else {
				v = strings.Replace(v, match[0], fmt.Sprintf(`<input type="checkbox" hx-post="/note/checkbox" hx-vars="box:%d" hx-swap="none">`, i), 1)
			}
		}

		parser := parser.NewWithExtensions(extensions)
		doc := parser.Parse([]byte(v))
		htmlFlags := html.CommonFlags | html.HrefTargetBlank
		opts := html.RendererOptions{Flags: htmlFlags}
		renderer := html.NewRenderer(opts)

		return template.HTML(markdown.Render(doc, renderer))

	},
}

func OnlyLocal(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		// Check if the request is local
		if !strings.HasPrefix(r.RemoteAddr, "[::1]") {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}
		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}
