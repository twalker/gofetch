package main

import (
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
)

type MockDB struct{}

func (mdb *MockDB) IsHealthy() bool {
	return true
}

func (mdb *MockDB) IncrementCounter() int {
	return 42
}

// newTestApplication helper returns an instance of the
// application struct containing mocked dependencies.
func newTestApplication() *application {
	return &application{
		logger: slog.New(slog.DiscardHandler),
		db:     &MockDB{},
	}
}

func TestHealthCheck(t *testing.T) {
	app := newTestApplication()
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
