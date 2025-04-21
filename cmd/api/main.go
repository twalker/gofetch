package main

import (
	"flag"
	"log/slog"
	"os"
	"sync"
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
	var cfg config

	flag.IntVar(&cfg.port, "port", 4000, "API server port")
	flag.StringVar(&cfg.env, "env", "development", "Environment (development|staging|production)")

	flag.Parse()

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
	})).With("pid", os.Getpid(), "name", "gofetch")

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
