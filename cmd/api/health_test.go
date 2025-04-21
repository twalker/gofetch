package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHealthCheck(t *testing.T) {
	app := &application{}
	rr := httptest.NewRecorder()

	r, err := http.NewRequest(http.MethodGet, "/health", nil)
	if err != nil {
		t.Fatal(err)
	}

	app.healthCheckHandler(rr, r)

	rs := rr.Result()
	contentType := rs.Header.Get("Content-Type")

	if rs.StatusCode != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v", rs.StatusCode, http.StatusOK)
	}

	if expectedContentType := "application/json"; contentType != expectedContentType {
		t.Errorf("Handler return wrong content type: got %v want %v", contentType, expectedContentType)
	}
}
