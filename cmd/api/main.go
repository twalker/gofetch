package main

import (
	"log"
	"log/slog"
	"os"
	"strconv"
	"sync"

	_ "github.com/joho/godotenv/autoload"
)

type config struct {
	port int
	env  string
}

type application struct {
	config config
	logger *slog.Logger
	wg     sync.WaitGroup
}

func main() {
	port, err := strconv.Atoi(os.Getenv("PORT"))
	if err != nil {
		log.Fatalf("Error reading PORT from environment variable: %v", err)
	}

	env := os.Getenv("ENV")
	if env == "" {
		env = "local"
	}
	cfg := config{
		port: port,
		env:  env,
	}

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
	})).With("pid", os.Getpid(), "name", "gofetch")

	app := &application{
		config: cfg,
		logger: logger,
	}
	mux := app.registerRoutes()
	err = app.serve(mux)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
}
