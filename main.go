package main

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"
)

type PageData struct {
	Content template.HTML
	Title   string
	Year    string
}

func main() {
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	http.HandleFunc("/homes", func(w http.ResponseWriter, r *http.Request) {
		renderPage(w, "views/homes.html", PageData{Title: "Haier the Realtor"})
	})

	http.HandleFunc("/media", func(w http.ResponseWriter, r *http.Request) {
		renderPage(w, "views/media.html", PageData{Title: "Haier the Creator"})
	})

	log.Fatal(http.ListenAndServe(":6969", nil))
}

func renderPage(w http.ResponseWriter, tmpl string, data PageData) {
	renderTemplate(w, "views/layout.html", PageData{
		Title:   data.Title,
		Content: renderTemplateToString(tmpl, data),
		Year:    fmt.Sprint(time.Now().Year()),
	})
}

func renderTemplate(w http.ResponseWriter, tmpl string, data PageData) {
	t, err := template.ParseFiles(tmpl)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = t.Execute(w, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func renderTemplateToString(tmpl string, data PageData) template.HTML {
	t, err := template.ParseFiles(tmpl)
	if err != nil {
		return ""
	}

	var buffer bytes.Buffer
	err = t.Execute(&buffer, data)
	if err != nil {
		return ""
	}

	return template.HTML(buffer.String())
}
