package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestHomePage is test file for home.go
func TestImageRecognitionPage(t *testing.T) {
	response, err := http.NewRequest("GET", "/imageRecognition", nil)
	if err != nil {
		t.Fatal(err)
	}
	// record response
	record := httptest.NewRecorder()
	handler := http.HandlerFunc(UploadImage)

	handler.ServeHTTP(record, response)

	if record.Code != http.StatusOK {
		t.Errorf("Get wrong response when try to connect to ImageRecognition page, Status: %v", record.Code)
	}
}
