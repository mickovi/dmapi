package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

const version = "1.0.0"

// config holds all configuration settings about the application.
// We will read it from command-line flags when the application starts.
type config struct {
	port int
	env  string
}

// application holds the dependecies for our HTTP handlers, helpers,
// and middleware.
type application struct {
	config config
	logger *log.Logger
}

func main() {
	var cfg config

	// Read the value of the port and env command-line flags into the config struct.
	// port 4000 and development environment are the default values.
	flag.IntVar(&cfg.port, "port", 4000, "API server port")
	flag.StringVar(&cfg.env, "env", "development", "Environment (development|staging|production)")
	flag.Parse()

	// logger writes messages to the standard out stream, prefixed
	// with the current date and time.
	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	// Instance of application.
	app := &application{
		config: cfg,
		logger: logger,
	}

	// HTTP server configuration.
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.port),
		Handler:      app.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	// Start the HTTP server.
	logger.Printf("starting %s server on %s", cfg.env, srv.Addr)
	err := srv.ListenAndServe()
	logger.Fatal(err)
}
