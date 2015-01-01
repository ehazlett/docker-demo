package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/lib/pq"
)

type (
	Datastore struct {
		db *sql.DB
	}

	LogEntry struct {
		Date time.Time
		Addr string
		Path string
	}
)

func NewDatastore(host string, port int, username string, password string, dbName string, dbSSLMode string) (*Datastore, error) {
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s", username, password, host, port, dbName, dbSSLMode)
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}
	ds := &Datastore{
		db: db,
	}
	ds.init()
	return ds, nil
}

func (ds *Datastore) init() {
	if _, err := ds.db.Exec(`create table if not exists logs (date timestamp without time zone, addr varchar(64), path varchar(256));`); err != nil {
		log.Printf("unable to initdb: %s\n", err)
	}
}

func (ds *Datastore) Add(logEntry *LogEntry) error {
	ts := logEntry.Date.Format(time.RFC3339)
	if _, err := ds.db.Exec(`insert into logs (date, addr, path) values ($1, $2, $3)`, ts, logEntry.Addr, logEntry.Path); err != nil {
		return err
	}
	log.Printf("added entry %s %s\n", ts, logEntry.Path)
	return nil
}

func (ds *Datastore) GetAll() ([]*LogEntry, error) {
	logs := []*LogEntry{}

	rows, err := ds.db.Query(`select * from logs;`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var entry = LogEntry{}
		var date time.Time
		var addr, path string
		if err := rows.Scan(&date, &addr, &path); err != nil {
			return nil, err
		}

		entry.Date = date
		entry.Addr = addr
		entry.Path = path

		logs = append(logs, &entry)
	}
	log.Printf("retrieved %d entries\n", len(logs))
	return logs, nil
}

func (ds *Datastore) Close() {
	ds.db.Close()
}
