package main

import (
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
)

var (
	mux        = http.NewServeMux()
	listenAddr string
)

type (
	Content struct {
		Title    string
		Hostname string
	}
)

func init() {
	flag.StringVar(&listenAddr, "listen", ":8080", "listen address")
}

func loadTemplate(filename string) (*template.Template, error) {
	return template.ParseFiles(filename)
}

func index(w http.ResponseWriter, r *http.Request) {
	t, err := loadTemplate("templates/index.html.tmpl")
	if err != nil {
		fmt.Printf("error loading template: %s\n", err)
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
	flag.Parse()

	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))
	mux.HandleFunc("/", index)

	log.Printf("listening on %s\n", listenAddr)

	if err := http.ListenAndServe(listenAddr, mux); err != nil {
		fmt.Errorf("error serving: %s", err)
	}
}
