package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
)

func (app *application) checkErrorResponseHandler(w http.ResponseWriter, r *http.Request) {
	code := r.PathValue("code")

	statusCode, err := strconv.Atoi(code)
	if err != nil {
		app.internalServerError(w, r, fmt.Errorf("status code is required to be an integer: %w", err))
		return
	}

	switch statusCode {
	case 400:
		app.badRequestResponse(w, r, errors.New("logged 400"))
	case 401:
		app.unauthorizedResponse(w, r, errors.New("logged 401"))
	case 403:
		app.forbiddenResponse(w, r, errors.New("logged 403"))
	case 404:
		app.notFoundResponse(w, r)
	default:
		app.internalServerError(w, r, errors.New("unmapped error status code"))
	}
}
