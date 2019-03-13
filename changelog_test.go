package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestHomePage is test file for home.go
func TestChangeLogPage(t *testing.T) {
	response, err := http.NewRequest("GET", "/changelog", nil)
	if err != nil {
		t.Fatal(err)
	}
	// record response
	record := httptest.NewRecorder()
	handler := http.HandlerFunc(ChangeLog)

	handler.ServeHTTP(record, response)

	if record.Code != http.StatusOK {
		t.Errorf("Get wrong response when try to connect to ChangeLog page, Status: %v", record.Code)
	}
}
