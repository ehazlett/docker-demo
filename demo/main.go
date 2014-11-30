package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
)

var (
	mux = http.NewServeMux()
)

type (
	Content struct {
		Title    string
		Hostname string
	}
)

func index(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("index.html.tmpl")
	if err != nil {
		fmt.Println("error loading template: %s", err)
		return
	}

	title := os.Getenv("TITLE")
	hostname, err := os.Hostname()
	if err != nil {
		hostname = "unknown"
	}

	cnt := &Content{
		Title:    title,
		Hostname: hostname,
	}

	t.Execute(w, cnt)
}

func main() {
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))
	mux.HandleFunc("/", index)

	log.Println("listening on :8080")

	if err := http.ListenAndServe(":8080", mux); err != nil {
		fmt.Errorf("error serving: %s", err)
	}
}
