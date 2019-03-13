package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestHomePage is test file for home.go
func TestHomePage(t *testing.T) {
	request, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}
	// record request
	record := httptest.NewRecorder()
	handler := http.HandlerFunc(Home)

	handler.ServeHTTP(record, request)

	if record.Code != http.StatusOK {
		t.Errorf("Get wrong response when try to connect to Home page, Status: %v", record.Code)
	}
}
