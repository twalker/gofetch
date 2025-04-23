package main

import (
	"net/http"
)

func (app *application) healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	resp := map[string]string{
		"health":        "OK",
		"env":           app.config.env,
		"correlationID": app.getCorrelationID(r.Context()),
	}

	if err := writeJSONData(w, http.StatusOK, resp); err != nil {
		writeJSONError(w, http.StatusInternalServerError, err.Error())
	}
}
