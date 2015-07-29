package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetIndex(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(index))
	defer ts.Close()

	res, err := http.Get(ts.URL)
	if err != nil {
		t.Fatal(err)

	}

	if res.StatusCode != 200 {
		t.Fatal("expected 200 status")
	}
}

func TestGetPing(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(ping))
	defer ts.Close()

	res, err := http.Get(ts.URL)
	if err != nil {
		t.Fatal(err)

	}

	if res.StatusCode != 200 {
		t.Fatal("expected 200 status")
	}
}
