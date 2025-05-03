package main

import (
	"log/slog"
	"os"

	"gofetch.timwalker.dev/internal/apiclient"
	"gofetch.timwalker.dev/internal/database"
	"gofetch.timwalker.dev/internal/env"

	_ "github.com/joho/godotenv/autoload"
)

func main() {
	cfg := config{
		port: env.GetInt("PORT", 4000),
		env:  env.GetString("ENV", "local"),
	}

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
	})).With("pid", os.Getpid(), "name", "gofetch")
	slog.SetDefault(logger)

	apiClient, err := apiclient.NewClient(env.GetString("API_BASE_URL", "http://localhost:4444"), nil)
	if err != nil {
		logger.Error(err.Error())
	}

	app := &application{
		config:    cfg,
		logger:    logger,
		apiClient: apiClient,
		db:        database.New(),
	}

	mux := app.registerRoutes()

	err = app.serve(mux)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
}
