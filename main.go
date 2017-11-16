package main

import (
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/urfave/cli"
)

var (
	mux           = http.NewServeMux()
	sessionCookie = "session"
	waitGroup     = sync.WaitGroup{}
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
	waitGroup.Add(1)
	defer waitGroup.Done()
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

	current, _ := r.Cookie(sessionCookie)
	if current == nil {
		current = &http.Cookie{
			Name:    sessionCookie,
			Value:   fmt.Sprintf("%d", time.Now().UnixNano()),
			Path:    "/",
			Expires: time.Now().AddDate(0, 0, 1),
			MaxAge:  86400,
		}
	}
	fmt.Printf("cookie: %s\n", current.Value)

	http.SetCookie(w, current)

	if err := json.NewEncoder(w).Encode(p); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func waitForDone(ctx context.Context) {
	waitGroup.Wait()
	ctx.Done()
}

func main() {
	app := cli.NewApp()
	app.Name = "docker-demo"
	app.Usage = "docker demo application"
	app.Version = "1.0.0"
	app.Author = "@ehazlett"
	app.Email = ""
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "listen-addr, l",
			Usage: "listen address",
			Value: ":8080",
		},
		cli.StringFlag{
			Name:  "tls-cert, c",
			Usage: "tls certificate",
			Value: "",
		},
		cli.StringFlag{
			Name:  "tls-key, k",
			Usage: "tls certificate key",
			Value: "",
		},
	}
	app.Action = func(c *cli.Context) error {
		mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))
		mux.HandleFunc("/ping", ping)
		mux.HandleFunc("/", index)

		hostname := getHostname()
		listenAddr := c.String("listen-addr")
		tlsCert := c.String("tls-cert")
		tlsKey := c.String("tls-key")

		srv := &http.Server{
			Handler:      mux,
			Addr:         listenAddr,
			WriteTimeout: time.Second * 10,
			ReadTimeout:  time.Second * 10,
		}

		log.Printf("instance: %s\n", hostname)
		log.Printf("listening on %s\n", listenAddr)

		ch := make(chan os.Signal, 1)
		signal.Notify(ch, os.Interrupt)

		go func() {
			select {
			case <-ch:
				log.Println("stopping")
				ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
				defer cancel()

				waitForDone(ctx)

				if err := srv.Shutdown(ctx); err != nil {
					log.Fatal(err)
				}
			}
		}()

		var err error
		if tlsCert != "" && tlsKey != "" {
			err = srv.ListenAndServeTLS(tlsCert, tlsKey)
		} else {
			err = srv.ListenAndServe()
		}

		return err
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
