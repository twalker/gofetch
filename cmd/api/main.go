package main

import (
	"flag"
	"fmt"
	"log/slog"

	"gofetch.timwalker.dev/internal/vcs"
)

type config struct {
	port int
	env  string
}

type application struct {
	config config
	logger *slog.Logger
}

var version = vcs.Version()

func main() {
	var cfg config

	// Read the value of the port and env command-line flags into the config struct. We
	// default to using the port number 4000 and the environment "development" if no
	// corresponding flags are provided.
	flag.IntVar(&cfg.port, "port", 4000, "API server port")
	flag.StringVar(&cfg.env, "env", "development", "Environment (development|staging|production)")

	fmt.Println("Hello world")
}
