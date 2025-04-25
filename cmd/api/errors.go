package main

import (
	"net/http"
)

// internalServerError returns a 500 error response and logs the provided error.
func (app *application) internalServerError(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Error("internal error", "method", r.Method, "path", r.URL.Path, "error", err.Error(), string(correlationIDContextKey), app.getCorrelationID(r.Context()))

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
