package main

import (
	"net/http"
)

func (app *application) internalServerError(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Error("internal error", "method", r.Method, "path", r.URL.Path, "error", err.Error())

	writeJSONError(w, http.StatusInternalServerError, "the server encountered a problem")
}

func (app *application) badRequestResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Warn("bad request", "method", r.Method, "path", r.URL.Path, "error", err.Error())

	writeJSONError(w, http.StatusBadRequest, err.Error())
}

func (app *application) forbiddenResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Warn("forbidden", "method", r.Method, "path", r.URL.Path, "error", err.Error())

	writeJSONError(w, http.StatusForbidden, "forbidden")
}

func (app *application) unauthorizedErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Warn("unauthorized error", "method", r.Method, "path", r.URL.Path, "error", err.Error())

	writeJSONError(w, http.StatusUnauthorized, "unauthorized")
}

func (app *application) notFoundResponse(w http.ResponseWriter, r *http.Request) {
	app.logger.Warn("not found error", "method", r.Method, "path", r.URL.Path)

	writeJSONError(w, http.StatusNotFound, "not found")
}
