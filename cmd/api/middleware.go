package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/google/uuid"
)

// contextKey is a custom type used for context keys to avoid collisions.
type contextKey string

const correlationIDHeader = "x-correlation-id"

// correlationIDKey is the key used to store the correlation ID in the context.
const correlationIDKey contextKey = "correlationID"

// correlationIDMiddleware generates or retrieves a request ID (UUID v4)
// and adds it to the request header, response header, and request context.
func (app *application) correlationIDMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := r.Header.Get(correlationIDHeader)

		if id == "" {
			newUUID, err := uuid.NewRandom()
			if err != nil {
				app.logger.Error(fmt.Sprintf("failed to generate UUID for correlationID: %v", err))
			} else {
				id = newUUID.String()
				r.Header.Set(correlationIDHeader, id)
			}
		}

		if id != "" {
			w.Header().Set(correlationIDHeader, id)
		}

		ctx := r.Context()
		if id != "" {
			ctx = context.WithValue(ctx, correlationIDKey, id)
		}
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}

// getCorrelationID retrieves the request ID from the context.
func (app *application) getCorrelationID(ctx context.Context) string {
	if id, ok := ctx.Value(correlationIDKey).(string); ok {
		return id
	}
	return ""
}
