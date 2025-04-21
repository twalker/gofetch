package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
)

// Helper function to check if a string is a valid UUID
func isValidUUID(u string) bool {
	_, err := uuid.Parse(u)
	return err == nil
}

func TestCorrelationIDMiddleware(t *testing.T) {
	t.Run("generates_id_when_header_is_missing", func(t *testing.T) {
		rr := httptest.NewRecorder()
		req, err := http.NewRequest(http.MethodConnect, "/", nil)
		if err != nil {
			t.Fatal(err)
		}
		var handlerIDHeader string
		var handlerIDContext string
		var handlerCalled bool
		app := &application{}

		nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			handlerCalled = true
			handlerIDHeader = w.Header().Get(correlationIDHeader)
			handlerIDContext = app.getCorrelationID(r.Context())
			w.WriteHeader(http.StatusOK)
		})

		app.correlationIDMiddleware(nextHandler).ServeHTTP(rr, req)

		if !handlerCalled {
			t.Fatal("Next handler was not called by the middleware")
		}

		if handlerIDHeader == "" {
			t.Errorf("Middleware failed to set '%s' header in response for next handler", correlationIDHeader)
		}
		if handlerIDContext == "" {
			t.Errorf("Middleware failed to set '%s' in context for next handler", correlationIDKey)
		}
		if !isValidUUID(handlerIDContext) {
			t.Errorf("Expected context key '%s' ('%s') to be a valid UUID, but it was not", correlationIDKey, handlerIDContext)
		}
		if status := rr.Code; status != http.StatusOK {
			t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusAccepted)
		}
	})

	t.Run("reuses_id_when_header_is_present", func(t *testing.T) {
		rr := httptest.NewRecorder()
		req, err := http.NewRequest(http.MethodConnect, "/", nil)
		if err != nil {
			t.Fatal(err)
		}
		existingID := "my-predefined-request-id-123"
		req.Header.Set(correlationIDHeader, existingID)
		var handlerIDHeader string
		var handlerIDContext string
		var handlerCalled bool
		app := &application{}

		nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			handlerCalled = true
			handlerIDHeader = w.Header().Get(correlationIDHeader)
			handlerIDContext = app.getCorrelationID(r.Context())
			w.WriteHeader(http.StatusOK)
		})

		app.correlationIDMiddleware(nextHandler).ServeHTTP(rr, req)

		if !handlerCalled {
			t.Fatal("Next handler was not called by the middleware")
		}
		if handlerIDHeader != existingID {
			t.Errorf("Middleware did not put existing header ID in context: expected '%s', got '%s'",
				existingID, handlerIDContext)
		}
		if handlerIDContext != existingID {
			t.Errorf("Middleware did not put existing header ID in context: expected '%s', got '%s'", existingID, handlerIDContext)
		}
		if status := rr.Code; status != http.StatusOK {
			t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusAccepted)
		}
	})
}
