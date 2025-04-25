package main

import (
	"net/http"
)

func (app *application) registerRoutes() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", app.healthCheckHandler)
	mux.HandleFunc("HEAD /health", app.healthCheckHandler)
	mux.HandleFunc("GET /errors/{code}", app.checkErrorResponseHandler)
	// Unmatched route patters receive a 404
	mux.HandleFunc("/", app.notFoundResponse)

	return app.correlationIDMiddleware(mux)
}
