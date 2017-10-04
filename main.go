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
		Title           string
		Version         string
		Hostname        string
		RefreshInterval string
		Metadata        string
		SkipErrors      bool
		ShowVersion     bool
	}

	Ping struct {
		Instance  string `json:"instance"`
		Version   string `json:"version"`
		Metadata  string `json:"metadata,omitempty"`
		RequestID string `json:"request_id,omitempty"`
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

func getVersion() string {
	ver := os.Getenv("VERSION")
	if ver == "" {
		ver = "0.1"
	}

	return ver
}

func loadTemplate(filename string) (*template.Template, error) {
	return template.ParseFiles(filename)
}

func getMetadata() string {
	return os.Getenv("METADATA")
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
	refreshInterval := os.Getenv("REFRESH_INTERVAL")
	if refreshInterval == "" {
		refreshInterval = "1000"
	}

	cnt := &Content{
		Title:           title,
		Version:         getVersion(),
		Hostname:        hostname,
		RefreshInterval: refreshInterval,
		Metadata:        getMetadata(),
		SkipErrors:      os.Getenv("SKIP_ERRORS") != "",
		ShowVersion:     os.Getenv("SHOW_VERSION") != "",
	}

	t.Execute(w, cnt)
}

func ping(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Connection", "close")

	hostname := getHostname()
	p := Ping{
		Instance: hostname,
		Version:  getVersion(),
		Metadata: getMetadata(),
	}

	requestID := r.Header.Get("X-Request-Id")
	if requestID != "" {
		p.RequestID = requestID
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
