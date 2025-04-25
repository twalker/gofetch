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

func (app *application) getAlbumsFromApiClientHandler(w http.ResponseWriter, r *http.Request) {
	var albums []Album // Define a struct User corresponding to the API response
	resp, err := app.apiClient.Get(r.Context(), "/static/albums.json", &albums)
	app.logger.Info("Received response",
		"status", resp.StatusCode,
	)
	if err != nil {
		// TODO: return status code from api response and log error.
		// status code mapping?
		// app.internalServerError(w, r, err)
		// resp.Write(w)
		writeJSONError(w, resp.StatusCode, err.Error())
		app.logger.Error("apiclient returned an error:", "fetchUrl", resp.Request.URL.Path, "fetchStatus", resp.StatusCode)
		return
	}
	writeJSONData(w, resp.StatusCode, albums)
}

func (app *application) createAlbumHandler(w http.ResponseWriter, r *http.Request) {
	var newAlbum Album
	err := readJSON(w, r, &newAlbum)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	err = writeJSONData(w, http.StatusCreated, newAlbum)
	if err != nil {
		app.internalServerError(w, r, err)
	}
}
