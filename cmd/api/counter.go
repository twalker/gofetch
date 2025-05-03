package main

import (
	"encoding/json"
	"net/http"
)

func (app *application) viewCounterHandler(w http.ResponseWriter, r *http.Request) {
	count := app.db.IncrementCounter()

	jsonResp, err := json.Marshal(map[string]int{"count": count})
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}

	err = writeJSON(w, http.StatusOK, jsonResp)
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, err.Error())
	}
}
