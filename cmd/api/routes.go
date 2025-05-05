package main

import (
	"net/http"
)

func (app *application) registerRoutes() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", app.healthCheckHandler)
	mux.HandleFunc("HEAD /health", app.healthCheckHandler)
	mux.HandleFunc("GET /errors/{code}", app.checkErrorResponseHandler)
	mux.HandleFunc("GET /apiclient/albums", app.getAlbumsFromApiClientHandler)
	mux.HandleFunc("POST /albums", app.createAlbumHandler)
	mux.HandleFunc("GET /counter", app.viewCounterHandler)
	mux.HandleFunc("GET /fibonacci/{num}", app.fibonacciHandler)
	// Unmatched route patters receive a 404
	mux.HandleFunc("/", app.notFoundResponse)

	return app.recoverPanic(app.correlationIDMiddleware(mux))
}
