package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestFaviconHandler(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(faviconHandler))
	resp, err := http.Get(server.URL)
	if err != nil {
		t.Error(err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Ожидаемый результат OK, получено: %d", resp.StatusCode)
	}
}

func TestJsHandler(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(jsHandler))
	resp, err := http.Get(server.URL)
	if err != nil {
		t.Error(err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Ожидаемый результат OK, получено: %d", resp.StatusCode)
	}
}

func TestCssHandler(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(cssHandler))
	resp, err := http.Get(server.URL)
	if err != nil {
		t.Error(err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Ожидаемый результат OK, получено: %d", resp.StatusCode)
	}
}
