package main

import (
	"errors"
	"html/template"
	"log"
	"net/http"
	"os"
	"regexp"
)

type Page struct {
	Title string
	Body  []byte
}

func (p *Page) save() error {
	filename := p.Title + ".txt"
	if _, err := os.Stat("pages"); err != nil {
		err := os.Mkdir("pages", 0700)
		if err != nil {
			return err
		}
	}
	return os.WriteFile("pages/"+filename, p.Body, 0600)
}

var validPath = regexp.MustCompile("^/(edit|save|view)/([a-zA-Z0-9]+)$")
var templates = template.Must(template.ParseFiles("edit.html", "view.html"))

func loadFile(title string) (*Page, error) {
	filename := "pages/" + title + ".txt"
	body, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return &Page{Title: title, Body: body}, nil
}

func renderTemplate(w http.ResponseWriter, tmpl string, p *Page) {
	err := templates.ExecuteTemplate(w, tmpl+".html", p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func getTitle(w http.ResponseWriter, r *http.Request) (string, error) {
	if m := validPath.FindStringSubmatch(r.URL.Path); m == nil {
		http.NotFound(w, r)
		return "", errors.New("invalid Page Title")
	}

	return r.PathValue("title"), nil
}

func makeHandler(fn func(http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		title, err := getTitle(w, r)
		if err != nil {
			return
		}
		fn(w, r, title)
	}
}

func viewHandler(w http.ResponseWriter, r *http.Request, title string) {
	p, err := loadFile(title)
	if err != nil {
		http.Redirect(w, r, "/edit/"+title, http.StatusFound)
		return
	}
	renderTemplate(w, "view", p)
}

func editHandler(w http.ResponseWriter, r *http.Request, title string) {
	p, err := loadFile(title)
	if err != nil {
		p = &Page{Title: title}
	}
	renderTemplate(w, "edit", p)
}

func saveHandler(w http.ResponseWriter, r *http.Request, title string) {
	body := r.FormValue("body")
	p := &Page{
		Title: title,
		Body:  []byte(body),
	}
	err := p.save()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/view/"+title, http.StatusFound)
}

func redirectFromIndex(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/view/home", http.StatusPermanentRedirect)

}

func main() {
	http.HandleFunc("/view", redirectFromIndex)
	http.HandleFunc("/view/{title}", makeHandler(viewHandler))
	http.HandleFunc("/edit/{title}", makeHandler(editHandler))
	http.HandleFunc("/save/{title}", makeHandler(saveHandler))
	log.Fatal(http.ListenAndServe(":8080", nil))
}
