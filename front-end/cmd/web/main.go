package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		render(w, "index.page.html")
	})

	fmt.Println("Front end server starting on port 8000 ... ")

	err := http.ListenAndServe(":8000", nil)
	if err != nil {
		log.Panic(err)
	}
}

func render(w http.ResponseWriter, t string) {
	partials := []string{
		"./cmd/web/templates/layout/base.layout.html",
		"./cmd/web/templates/layout/header.partial.html",
		"./cmd/web/templates/layout/footer.partial.html",
	}

	var renderPage []string
	renderPage = append(renderPage, fmt.Sprintf("./cmd/web/templates/pages/%s", t))
	renderPage = append(renderPage, partials...)

	tpl, err := template.ParseFiles(renderPage...)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := tpl.Execute(w, nil); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
