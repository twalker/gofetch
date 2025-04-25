package apiclient

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

// package adapted from: https://gemini.google.com/app/6750213f3e04a7cd
var httpClient = &http.Client{
	Timeout: 30 * time.Second, // Default timeout
}

type APIClient struct {
	baseURL    *url.URL
	httpClient *http.Client
}

// NewClient creates a new instance of the Open API client.
// It requires the base URL string (e.g., "https://api.example.com") for the target host.
func NewClient(baseUrl string, client *http.Client) (*APIClient, error) {
	baseURL, err := url.Parse(baseUrl)
	if err != nil {
		return nil, fmt.Errorf("failed to parse base URL %q: %w", baseUrl, err)
	}
	if baseURL.Scheme == "" || baseURL.Host == "" {
		return nil, fmt.Errorf("invalid base URL %q: must include scheme and host", baseURL)
	}
	// Use the provided client or create a default one.
	if client == nil {
		client = httpClient
	}
	return &APIClient{
		baseURL:    baseURL,
		httpClient: client,
	}, nil
}

// newRequest creates an API request. A relative URL path can be provided in
// path, in which case it is resolved relative to the baseURL of the Client.
// Relative paths should always be specified without a preceding slash.
// If specified, the value pointed to by body is JSON encoded and included
// as the request body.
func (c *APIClient) newRequest(ctx context.Context, method, path string, body any) (*http.Request, error) {
	// Resolve the relative path against the base URL.
	relURL, err := url.Parse(path)
	if err != nil {
		return nil, fmt.Errorf("failed to parse relative path %q: %w", path, err)
	}
	fullURL := c.baseURL.ResolveReference(relURL)

	var buf io.ReadWriter
	if body != nil {
		buf = &bytes.Buffer{}
		// Encode the body to JSON.
		enc := json.NewEncoder(buf)
		enc.SetEscapeHTML(false) // Prevent encoding of <, >, &
		err := enc.Encode(body)
		if err != nil {
			return nil, fmt.Errorf("failed to encode request body: %w", err)
		}
	}

	// Create the HTTP request with context.
	req, err := http.NewRequestWithContext(ctx, method, fullURL.String(), buf)
	if err != nil {
		return nil, fmt.Errorf("failed to create new request: %w", err)
	}

	// Set standard headers.
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	req.Header.Set("Accept", "application/json")
	// Add any other common headers needed for your API (e.g., User-Agent, Authorization)
	// req.Header.Set("User-Agent", "my-app/1.0")
	// req.Header.Set("Authorization", "Bearer YOUR_API_KEY")

	return req, nil
}

// do sends an API request and returns the API response. The API response is
// JSON decoded and stored in the value pointed to by v, or returned as an
// error if an API error has occurred. If v implements the io.Writer
// interface, the raw response body will be written to v, without attempting to
// first decode it.
func (c *APIClient) do(req *http.Request, v any) (*http.Response, error) {
	// Execute the request using the configured http client.
	resp, err := c.httpClient.Do(req)
	if err != nil {
		// If the context was canceled, return that error.
		select {
		case <-req.Context().Done():
			return nil, req.Context().Err()
		default:
		}
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close() // Ensure the response body is closed.

	// Check for non-success status codes.
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		// Attempt to read the error body for more context, but don't fail if reading fails.
		bodyBytes, _ := io.ReadAll(resp.Body)
		return resp, fmt.Errorf("request failed with status %d: %s", resp.StatusCode, string(bodyBytes))
	}

	// If v is provided and is an io.Writer, write the raw body to it.
	if w, ok := v.(io.Writer); ok {
		_, err = io.Copy(w, resp.Body)
		if err != nil {
			return resp, fmt.Errorf("failed to write response body to writer: %w", err)
		}
	} else if v != nil {
		// Otherwise, decode the JSON response body into v.
		err = json.NewDecoder(resp.Body).Decode(v)
		// Handle EOF error specifically for empty bodies or non-JSON responses.
		if err == io.EOF {
			// Ignore EOF errors if the response body is empty.
			err = nil
		}
		if err != nil {
			return resp, fmt.Errorf("failed to decode response body: %w", err)
		}
	}

	return resp, nil
}

// --- Example Request Methods ---

// Get performs a GET request to the specified path.
// The response body is decoded into the value pointed to by `responsePayload`.
func (c *APIClient) Get(ctx context.Context, path string, responsePayload any) (*http.Response, error) {
	req, err := c.newRequest(ctx, http.MethodGet, path, nil) // No body for GET
	if err != nil {
		return nil, err
	}
	return c.do(req, responsePayload)
}

// Post performs a POST request to the specified path with the given request body.
// The request body is JSON encoded.
// The response body is decoded into the value pointed to by `responsePayload`.
func (c *APIClient) Post(ctx context.Context, path string, requestBody, responsePayload any) (*http.Response, error) {
	req, err := c.newRequest(ctx, http.MethodPost, path, requestBody)
	if err != nil {
		return nil, err
	}
	return c.do(req, responsePayload)
}

// Put performs a PUT request similar to Post.
func (c *APIClient) Put(ctx context.Context, path string, requestBody, responsePayload any) (*http.Response, error) {
	req, err := c.newRequest(ctx, http.MethodPut, path, requestBody)
	if err != nil {
		return nil, err
	}
	return c.do(req, responsePayload)
}

// Delete performs a DELETE request to the specified path.
// Often, DELETE requests don't have a request body or expect a specific response payload.
func (c *APIClient) Delete(ctx context.Context, path string, responsePayload any) (*http.Response, error) {
	req, err := c.newRequest(ctx, http.MethodDelete, path, nil) // No body for DELETE usually
	if err != nil {
		return nil, err
	}
	return c.do(req, responsePayload)
}
