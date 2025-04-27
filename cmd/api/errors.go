package main

import (
	"net/http"
	"runtime/debug"
)

// internalServerError returns a 500 error response and logs the provided error.
func (app *application) internalServerError(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Error("internal error", "method", r.Method, "path", r.URL.Path, "error", err.Error(), string(correlationIDContextKey), app.getCorrelationID(r.Context()), "stack", string(debug.Stack()))

	writeJSONError(w, http.StatusInternalServerError, "the server encountered a problem")
}

// badRequestResponse returns a 400 error response and logs the provided error.
func (app *application) badRequestResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Warn("bad request", "method", r.Method, "path", r.URL.Path, "error", err.Error(), string(correlationIDContextKey), app.getCorrelationID(r.Context()))

	writeJSONError(w, http.StatusBadRequest, err.Error())
}

// forbiddenResponse returns a 403 response and logs the provided error.
func (app *application) forbiddenResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Warn("forbidden", "method", r.Method, "path", r.URL.Path, "error", err.Error(), string(correlationIDContextKey), app.getCorrelationID(r.Context()))

	writeJSONError(w, http.StatusForbidden, http.StatusText(http.StatusForbidden))
}

// unauthorizedResponse returns a 401 error response and logs the provided error.
func (app *application) unauthorizedResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Warn("unauthorized error", "method", r.Method, "path", r.URL.Path, "error", err.Error(), string(correlationIDContextKey), app.getCorrelationID(r.Context()))

	writeJSONError(w, http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized))
}

// notFoundResponse returns a 404 error respons and logs the provided error.
func (app *application) notFoundResponse(w http.ResponseWriter, r *http.Request) {
	app.logger.Warn("not found error", "method", r.Method, "path", r.URL.Path, string(correlationIDContextKey), app.getCorrelationID(r.Context()))

	writeJSONError(w, http.StatusNotFound, http.StatusText(http.StatusNotFound))
}

// apiClientErrorResponse returns the status code from the api client logs the error.
func (app *application) apiClientErrorResponse(w http.ResponseWriter, r *http.Request, status int, err error) {
	app.logger.Error("apiclient error", "method", r.Method, "path", r.URL.Path, "error", err.Error(), string(correlationIDContextKey), app.getCorrelationID(r.Context()), "stack", string(debug.Stack()))

	writeJSONError(w, status, "the server encountered a problem")
}
