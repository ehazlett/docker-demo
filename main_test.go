package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetIndex(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(index))
	defer ts.Close()

	res, err := http.Get(ts.URL)
	if err != nil {
		t.Fatal(err)

	}

	assert.Equal(t, 200, res.StatusCode, "expected response code 200")
}

func TestGetPing(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(ping))
	defer ts.Close()

	res, err := http.Get(ts.URL)
	if err != nil {
		t.Fatal(err)

	}

	assert.Equal(t, 200, res.StatusCode, "expected response code 200")
}
