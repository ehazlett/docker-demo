package main

import (
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

var (
	mux       = http.NewServeMux()
	DS        *Datastore
	DbHost    string
	DbUser    string
	DbPass    string
	DbPort    int
	DbName    string
	DbSSLMode string
)

type (
	Content struct {
		Title    string
		Hostname string
		Entries  []*LogEntry
	}
)

func init() {
	flag.StringVar(&DbHost, "db-host", "", "database host")
	flag.IntVar(&DbPort, "db-port", 0, "database port")
	flag.StringVar(&DbUser, "db-user", "", "database username")
	flag.StringVar(&DbPass, "db-pass", "", "database password")
	flag.StringVar(&DbName, "db-name", "", "database name")
	flag.StringVar(&DbSSLMode, "db-ssl", "", "database ssl mode")
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

func dbIndex(w http.ResponseWriter, r *http.Request) {
	DS.Add(&LogEntry{
		Date: time.Now(),
		Addr: r.RemoteAddr,
		Path: r.URL.Path,
	})
	t, err := loadTemplate("templates/db.html.tmpl")
	if err != nil {
		fmt.Printf("error loading template: %s\n", err)
		return
	}

	title := os.Getenv("TITLE")

	hostname, err := os.Hostname()
	if err != nil {
		hostname = "unknown"
	}

	entries, err := DS.GetAll()
	if err != nil {
		log.Printf("error retrieving data: %s\n", err)
	}

	cnt := &Content{
		Title:    title,
		Hostname: hostname,
		Entries:  entries,
	}

	t.Execute(w, cnt)
}

func main() {
	flag.Parse()
	// load vars from env
	if DbHost == "" {
		DbHost = os.Getenv("DB_HOST")
	}

	if DbPort == 0 {
		port := os.Getenv("DB_PORT")
		if port != "" {
			p, err := strconv.Atoi(port)
			if err != nil {
				log.Fatalf("unable to parse database port: %s", err)
			}

			DbPort = p
		}
	}

	if DbUser == "" {
		DbUser = os.Getenv("DB_USER")
	}

	if DbPass == "" {
		DbPass = os.Getenv("DB_PASS")
	}

	if DbName == "" {
		DbName = os.Getenv("DB_NAME")
	}

	sslMode := os.Getenv("DB_SSL_MODE")
	if sslMode != "" && DbSSLMode == "" {
		DbSSLMode = sslMode
	}

	ds, err := NewDatastore(DbHost, DbPort, DbUser, DbPass, DbName, DbSSLMode)
	if err != nil {
		log.Fatalf("unable to connect to datastore: %s", err)
	}

	DS = ds
	defer DS.Close()

	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))
	mux.HandleFunc("/", index)
	mux.HandleFunc("/db", dbIndex)

	log.Println("listening on :8080")

	if err := http.ListenAndServe(":8080", mux); err != nil {
		fmt.Errorf("error serving: %s", err)
	}
}
