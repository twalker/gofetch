package main

import (
	"net/http"
)

func (app *application) registerRoutes() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", app.healthCheckHandler)
	mux.HandleFunc("HEAD /health", app.healthCheckHandler)

	return app.correlationIDMiddleware(mux)
}
