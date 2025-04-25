package apiclient

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

/* test adapted from: https://gemini.google.com/app/6750213f3e04a7cd */

// setupTestServer creates a mock HTTP server for testing.
// It takes a handler function that defines the server's behavior.
func setupTestServer(handler http.HandlerFunc) (*httptest.Server, *APIClient) {
	server := httptest.NewServer(handler)
	client, err := NewClient(server.URL, nil) // Use the test server's URL
	if err != nil {
		server.Close() // Clean up server if client creation fails
		panic(fmt.Sprintf("Failed to create client for test server: %v", err))
	}
	return server, client
}

// TestNewClient tests the NewClient constructor function.
func TestNewClient(t *testing.T) {
	t.Run("ValidBaseURL", func(t *testing.T) {
		client, err := NewClient("https://example.com/api", nil)
		if err != nil {
			t.Fatalf("NewClient failed with valid URL: %v", err)
		}
		if client == nil {
			t.Fatal("NewClient returned nil client for valid URL")
		}
		if client.baseURL.String() != "https://example.com/api" {
			t.Errorf("Expected base URL %q, got %q", "https://example.com/api", client.baseURL.String())
		}
		if client.httpClient.Timeout != 30*time.Second {
			t.Errorf("Expected default timeout 30s, got %v", client.httpClient.Timeout)
		}
	})

	t.Run("ValidBaseURLWithCustomClient", func(t *testing.T) {
		customHTTPClient := &http.Client{Timeout: 60 * time.Second}
		client, err := NewClient("http://localhost:8080", customHTTPClient)
		if err != nil {
			t.Fatalf("NewClient failed with valid URL and custom client: %v", err)
		}
		if client.httpClient != customHTTPClient {
			t.Error("NewClient did not use the provided custom http.Client")
		}
		if client.httpClient.Timeout != 60*time.Second {
			t.Errorf("Expected custom timeout 60s, got %v", client.httpClient.Timeout)
		}
	})

	t.Run("InvalidURL_MissingScheme", func(t *testing.T) {
		_, err := NewClient("example.com/api", nil)
		if err == nil {
			t.Error("NewClient succeeded with missing scheme, expected error")
		} else if !strings.Contains(err.Error(), "must include scheme and host") {
			t.Errorf("Expected scheme/host error, got: %v", err)
		}
	})

	t.Run("InvalidURL_Unparseable", func(t *testing.T) {
		_, err := NewClient("://invalid-url", nil)
		if err == nil {
			t.Error("NewClient succeeded with unparseable URL, expected error")
		} else if !strings.Contains(err.Error(), "failed to parse base URL") {
			t.Errorf("Expected parse error, got: %v", err)
		}
	})
}

// --- Test Request Methods ---

// MockPayload is a simple struct for testing request/response bodies.
type MockPayload struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func TestClient_Get(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("Expected method GET, got %s", r.Method)
		}
		if r.URL.Path != "/test" {
			t.Errorf("Expected path /test, got %s", r.URL.Path)
		}
		if r.Header.Get("Accept") != "application/json" {
			t.Errorf("Expected Accept: application/json header, got %q", r.Header.Get("Accept"))
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(MockPayload{ID: 1, Name: "Test Get"})
	}

	server, client := setupTestServer(handler)
	defer server.Close()

	var responsePayload MockPayload
	resp, err := client.Get(context.Background(), "/test", &responsePayload)
	if err != nil {
		t.Fatalf("client.Get failed: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, resp.StatusCode)
	}
	if responsePayload.ID != 1 || responsePayload.Name != "Test Get" {
		t.Errorf("Unexpected response payload: %+v", responsePayload)
	}
}

func TestClient_Post(t *testing.T) {
	expectedRequestBody := MockPayload{ID: 10, Name: "Test Post Request"}
	expectedResponseBody := MockPayload{ID: 11, Name: "Test Post Response"}

	handler := func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("Expected method POST, got %s", r.Method)
		}
		if r.URL.Path != "/items" {
			t.Errorf("Expected path /items, got %s", r.URL.Path)
		}
		if r.Header.Get("Content-Type") != "application/json" {
			t.Errorf("Expected Content-Type: application/json header, got %q", r.Header.Get("Content-Type"))
		}

		// Decode request body to verify
		var receivedBody MockPayload
		if err := json.NewDecoder(r.Body).Decode(&receivedBody); err != nil {
			http.Error(w, "Failed to decode request body", http.StatusBadRequest)
			return
		}
		if receivedBody != expectedRequestBody {
			t.Errorf("Received unexpected request body. Got %+v, want %+v", receivedBody, expectedRequestBody)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated) // Typically 201 for POST
		json.NewEncoder(w).Encode(expectedResponseBody)
	}

	server, client := setupTestServer(handler)
	defer server.Close()

	var responsePayload MockPayload
	resp, err := client.Post(context.Background(), "/items", expectedRequestBody, &responsePayload)
	if err != nil {
		t.Fatalf("client.Post failed: %v", err)
	}
	if resp.StatusCode != http.StatusCreated {
		t.Errorf("Expected status %d, got %d", http.StatusCreated, resp.StatusCode)
	}
	if responsePayload != expectedResponseBody {
		t.Errorf("Unexpected response payload. Got %+v, want %+v", responsePayload, expectedResponseBody)
	}
}

func TestClient_Put(t *testing.T) {
	// Similar structure to TestClient_Post, adjust method and expected status code
	expectedRequestBody := MockPayload{ID: 20, Name: "Test Put Request"}
	expectedResponseBody := MockPayload{ID: 20, Name: "Test Put Response Updated"}

	handler := func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Errorf("Expected method PUT, got %s", r.Method)
		}
		if r.URL.Path != "/items/20" {
			t.Errorf("Expected path /items/20, got %s", r.URL.Path)
		}
		// ... (verify headers and request body as in POST) ...
		var receivedBody MockPayload
		if err := json.NewDecoder(r.Body).Decode(&receivedBody); err != nil {
			http.Error(w, "Failed to decode request body", http.StatusBadRequest)
			return
		}
		if receivedBody != expectedRequestBody {
			t.Errorf("Received unexpected request body. Got %+v, want %+v", receivedBody, expectedRequestBody)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK) // Typically 200 for successful PUT
		json.NewEncoder(w).Encode(expectedResponseBody)
	}

	server, client := setupTestServer(handler)
	defer server.Close()

	var responsePayload MockPayload
	resp, err := client.Put(context.Background(), "/items/20", expectedRequestBody, &responsePayload)
	if err != nil {
		t.Fatalf("client.Put failed: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, resp.StatusCode)
	}
	if responsePayload != expectedResponseBody {
		t.Errorf("Unexpected response payload. Got %+v, want %+v", responsePayload, expectedResponseBody)
	}
}

func TestClient_Delete(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("Expected method DELETE, got %s", r.Method)
		}
		if r.URL.Path != "/items/30" {
			t.Errorf("Expected path /items/30, got %s", r.URL.Path)
		}
		// DELETE often doesn't have a request body to check

		w.WriteHeader(http.StatusNoContent) // Typically 204 for successful DELETE
	}

	server, client := setupTestServer(handler)
	defer server.Close()

	// Often DELETE doesn't return a body, so responsePayload might be nil or an empty struct
	var responsePayload interface{} // Or specific empty struct if needed
	resp, err := client.Delete(context.Background(), "/items/30", &responsePayload)
	if err != nil {
		t.Fatalf("client.Delete failed: %v", err)
	}
	if resp.StatusCode != http.StatusNoContent {
		t.Errorf("Expected status %d, got %d", http.StatusNoContent, resp.StatusCode)
	}
	// Check if responsePayload remained nil or its expected empty state if applicable
}

func TestClient_ErrorHandling_Non2xxStatus(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound) // 404 Not Found
		fmt.Fprintln(w, `{"error": "Resource not found"}`)
	}

	server, client := setupTestServer(handler)
	defer server.Close()

	var responsePayload MockPayload
	resp, err := client.Get(context.Background(), "/notfound", &responsePayload)

	if err == nil {
		t.Fatalf("Expected an error for non-2xx status, but got nil")
	}
	if resp == nil {
		t.Fatalf("Expected non-nil response even on error, got nil")
	}
	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("Expected status %d in response, got %d", http.StatusNotFound, resp.StatusCode)
	}
	// Check if the error message contains the status code and potentially the body
	expectedErrorMsg := fmt.Sprintf("request failed with status %d", http.StatusNotFound)
	if !strings.Contains(err.Error(), expectedErrorMsg) {
		t.Errorf("Expected error message to contain %q, got %q", expectedErrorMsg, err.Error())
	}
	// Optional: check if error contains body snippet
	if !strings.Contains(err.Error(), "Resource not found") {
		t.Errorf("Expected error message to contain response body snippet, got %q", err.Error())
	}
}

func TestClient_ErrorHandling_DecodeError(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, `{"id": 1, "name": "Test Decode",}`) // Malformed JSON (trailing comma)
	}

	server, client := setupTestServer(handler)
	defer server.Close()

	var responsePayload MockPayload
	_, err := client.Get(context.Background(), "/malformed", &responsePayload)

	if err == nil {
		t.Fatalf("Expected a JSON decoding error, but got nil")
	}
	if !strings.Contains(err.Error(), "failed to decode response body") {
		t.Errorf("Expected error message to contain 'failed to decode response body', got %q", err.Error())
	}
}

func TestClient_ContextCancellation(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(100 * time.Millisecond) // Simulate a delay
		w.WriteHeader(http.StatusOK)
	}

	server, client := setupTestServer(handler)
	defer server.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond) // Very short timeout
	defer cancel()

	_, err := client.Get(ctx, "/slow", nil)

	if err == nil {
		t.Fatalf("Expected a context deadline exceeded error, but got nil")
	}
	if err != context.DeadlineExceeded {
		t.Errorf("Expected error context.DeadlineExceeded, got %v", err)
	}
}

func TestClient_Do_WriteToWriter(t *testing.T) {
	expectedBody := `{"raw": true}`
	handler := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, expectedBody)
	}

	server, client := setupTestServer(handler)
	defer server.Close()

	var buf strings.Builder // Use strings.Builder as an io.Writer
	req, _ := client.newRequest(context.Background(), http.MethodGet, "/raw", nil)
	resp, err := client.do(req, &buf)
	if err != nil {
		t.Fatalf("client.do failed: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, resp.StatusCode)
	}
	if buf.String() != expectedBody {
		t.Errorf("Expected body %q written to writer, got %q", expectedBody, buf.String())
	}
}

func TestClient_newRequest_PathResolution(t *testing.T) {
	client, err := NewClient("https://base.com/v1/", nil)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	tests := []struct {
		name        string
		path        string
		expectedURL string
		expectedErr bool
	}{
		{"SimplePath", "users", "https://base.com/v1/users", false},
		{"PathWithSlash", "/users", "https://base.com/users", false}, // ResolveReference treats leading slash as absolute path from host
		{"EmptyPath", "", "https://base.com/v1/", false},
		{"PathWithQuery", "items?id=123", "https://base.com/v1/items?id=123", false},
		{"AbsolutePath", "https://othersite.com/path", "https://othersite.com/path", false}, // Absolute URL overrides base
		{"InvalidRelativePath", ":invalid:", "", true},                                      // Causes url.Parse error
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := client.newRequest(context.Background(), http.MethodGet, tt.path, nil)

			if tt.expectedErr {
				if err == nil {
					t.Errorf("Expected an error for path %q, but got nil", tt.path)
				}
			} else {
				if err != nil {
					t.Fatalf("newRequest failed for path %q: %v", tt.path, err)
				}
				if req.URL.String() != tt.expectedURL {
					t.Errorf("Expected URL %q, got %q", tt.expectedURL, req.URL.String())
				}
			}
		})
	}
}
