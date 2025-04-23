package main

import (
	"log/slog"
	"os"

	_ "github.com/joho/godotenv/autoload"
	"gofetch.timwalker.dev/internal/env"
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

	app := &application{
		config: cfg,
		logger: logger,
	}

	mux := app.registerRoutes()

	err := app.serve(mux)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
}
