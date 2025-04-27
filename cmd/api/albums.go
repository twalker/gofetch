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
	var albums []Album
	resp, err := app.apiClient.Get(r.Context(), "/static/albums.json", &albums)
	if err != nil {
		if resp != nil {
			app.apiClientErrorResponse(w, r, resp.StatusCode, err)
			return
		}
		app.internalServerError(w, r, err)
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
