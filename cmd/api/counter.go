package main

import (
	"encoding/json"
	"net/http"
)

// viewCounterHandler shows and increments a counter stored in Redis.
func (app *application) viewCounterHandler(w http.ResponseWriter, r *http.Request) {
	count := app.db.IncrementCounter()

	jsonResp, err := json.Marshal(map[string]int{"count": count})
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.Write(jsonResp)
}
