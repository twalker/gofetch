package main

import (
	"net/http"
)

// Example struct to hold API data
type Album struct {
	ID     int    `json:"id"`
	UserID int    `json:"userId"`
	Title  string `json:"title"`
}

func (app *application) getAlbumsWithClientHandler(w http.ResponseWriter, r *http.Request) {
	var albums []Album // Define a struct User corresponding to the API response
	resp, err := app.apiClient.Get(r.Context(), "/static/albums.json", &albums)
	app.logger.Info("Received response",
		"status", resp.StatusCode,
	)
	if err != nil {
		// TODO: return status code from api response. status code mapping?
		// app.internalServerError(w, r, err)
		// resp.Write(w)
		writeJSONError(w, resp.StatusCode, err.Error())
		return
	}
	writeJSONData(w, resp.StatusCode, albums)
}
