package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func (app *application) healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	cid := app.getCorrelationID(r.Context())
	resp := map[string]string{"health": "OK", "correlationID": cid}
	jsonResp, err := json.Marshal(resp)
	if err != nil {
		http.Error(w, "failed to marshal response", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if _, err := w.Write(jsonResp); err != nil {
		app.logger.Error(fmt.Sprintf("Failed to write response: %v", err))
	}
}
