package main

import (
	"encoding/json"
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

	Ping struct {
		Instance string `json:"instance"`
	}
)

func init() {
	flag.StringVar(&listenAddr, "listen", ":8080", "listen address")
}

func getHostname() string {
	hostname, err := os.Hostname()
	if err != nil {
		hostname = "unknown"
	}

	return hostname
}

func loadTemplate(filename string) (*template.Template, error) {
	return template.ParseFiles(filename)
}

func index(w http.ResponseWriter, r *http.Request) {
	remote := r.RemoteAddr

	forwarded := r.Header.Get("X-Forwarded-For")
	if forwarded != "" {
		remote = forwarded
	}

	log.Printf("request from %s\n", remote)

	t, err := loadTemplate("templates/index.html.tmpl")
	if err != nil {
		fmt.Printf("error loading template: %s\n", err)
		return
	}

	title := os.Getenv("TITLE")

	hostname := getHostname()

	cnt := &Content{
		Title:    title,
		Hostname: hostname,
	}

	t.Execute(w, cnt)
}

func ping(w http.ResponseWriter, r *http.Request) {
	hostname := getHostname()
	p := Ping{
		Instance: hostname,
	}

	if err := json.NewEncoder(w).Encode(p); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func main() {
	flag.Parse()

	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))
	mux.HandleFunc("/ping", ping)
	mux.HandleFunc("/", index)

	hostname := getHostname()

	log.Printf("instance: %s\n", hostname)
	log.Printf("listening on %s\n", listenAddr)

	if err := http.ListenAndServe(listenAddr, mux); err != nil {
		log.Fatalf("error serving: %s", err)
	}
}
