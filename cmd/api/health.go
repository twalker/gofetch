package main

import (
	"net/http"
)

func (app *application) healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodHead {
		return
	}

	type health struct {
		CorrelationID        string `json:"correlationId"`
		Env                  string `json:"env"`
		ApplicationIsHealthy bool   `json:"appIsHealthy"`
		DatabaseIsHealthy    bool   `json:"dbIsHealthy"`
	}

	response := health{
		CorrelationID:        app.getCorrelationID(r.Context()),
		Env:                  app.config.env,
		ApplicationIsHealthy: true,
		DatabaseIsHealthy:    app.db.IsHealthy(),
	}

	if err := writeJSON(w, http.StatusOK, response); err != nil {
		writeJSONError(w, http.StatusInternalServerError, err.Error())
	}
}
