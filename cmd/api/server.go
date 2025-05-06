package main

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"gofetch.timwalker.dev/internal/apiclient"
	"gofetch.timwalker.dev/internal/database"
)

type config struct {
	port int
	env  string
}

type application struct {
	config    config
	logger    *slog.Logger
	apiClient *apiclient.APIClient
	db        database.Service
}

func (app *application) serve(mux http.Handler) error {
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", app.config.port),
		Handler:      mux,
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		ErrorLog:     slog.NewLogLogger(app.logger.Handler(), slog.LevelError),
	}

	app.logger.Info("starting server", "addr", srv.Addr, "env", app.config.env)

	err := srv.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	app.logger.Info("stopped server", "addr", srv.Addr)

	return nil
}
