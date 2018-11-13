package main

import (
	"html/template"
	"log"
	"net/http"
	"regexp"
)

// Page includes info to be inserted in page
type Page struct {
	Title     string
	VersionID int
	ClusterID string
	BlockID   string
}

func loadPage(title string) (*Page, error) {
	return &Page{Title: title, VersionID: getVersionNo(), ClusterID: getClusterID(), BlockID: getBlockID()}, nil
}

func viewHandler(w http.ResponseWriter, r *http.Request, title string) {
	p, err := loadPage(title)
	if err != nil {
		print("err in handler")
	}
	renderTemplate(w, "view", p)
}

var templates = template.Must(template.ParseFiles("view.html"))

func renderTemplate(w http.ResponseWriter, tmpl string, p *Page) {
	err := templates.ExecuteTemplate(w, tmpl+".html", p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

var validPath = regexp.MustCompile("^/(edit|save|view)/([a-zA-Z0-9]+)$")

func makeHandler(fn func(http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		m := validPath.FindStringSubmatch(r.URL.Path)
		if m == nil {
			http.NotFound(w, r)
			return
		}
		fn(w, r, m[2])
	}
}

func main() {
	http.HandleFunc("/view/", makeHandler(viewHandler))
	log.Fatal(http.ListenAndServe(":8080", nil))
}
