package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/google/uuid"
)

// correlationIDHeaderKey is the custome http header name for correlation ids.
const correlationIDHeaderKey = "x-correlation-id"

// contextKey is a custom type used for context keys to avoid collisions.
type contextKey string

// correlationIDContextKey is the key used to store the correlation ID in the context.
const correlationIDContextKey contextKey = "correlationID"

// correlationIDMiddleware generates or retrieves a request ID (UUID v4)
// and adds it to the request header, response header, and request context.
func (app *application) correlationIDMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := r.Header.Get(correlationIDHeaderKey)

		if id == "" {
			newUUID, err := uuid.NewRandom()
			if err != nil {
				app.logger.Error(fmt.Sprintf("failed to generate UUID for correlationID: %v", err))
			}
			id = newUUID.String()
		}

		w.Header().Set(correlationIDHeaderKey, id)

		ctx := r.Context()
		ctx = context.WithValue(ctx, correlationIDContextKey, id)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}

// getCorrelationID retrieves the request ID from the context.
func (app *application) getCorrelationID(ctx context.Context) string {
	if id, ok := ctx.Value(correlationIDContextKey).(string); ok {
		return id
	}
	return ""
}

// recoverPanic middleware recovers from panics, logs the err, and prevents
// the server process from exiting.
func (app *application) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// The deferred function will always be run in the event of a panic
		// as Go unwinds the stack.
		defer func() {
			// Use the builtin recover function to check if there has been a panic or not.
			if err := recover(); err != nil {
				// If there was a panic, set a "Connection: close" header on the
				// response. This acts as a trigger to make Go's HTTP server
				// automatically close the current connection after a response has been
				// sent.
				w.Header().Set("Connection", "close")
				// The value returned by recover() has the type any, soo fmt.Errorf() is
				// used to normalize it into an error, return a 500 and, log the error.
				app.internalServerError(w, r, fmt.Errorf("%s", err))
			}
		}()

		next.ServeHTTP(w, r)
	})
}
