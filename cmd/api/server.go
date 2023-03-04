package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func (app *application) serve() error {
	// HTTP server configuration.
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", app.config.port),
		Handler:      app.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	// A shutdownError channel. We will use this to receive any errors returned
	// by the graceful Shutdown() function.
	shutdownError := make(chan error)

	go func() {
		// Intercept the signals.
		quit := make(chan os.Signal, 1)

		// Use signal.Notify() to listen for incoming SIGINT and SIGTERM signals
		// and relay them to the quit channel. Any ohter signals will not be caught
		// by signal.Notify() and will retain their default behavior.
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

		// Read the signal from the quit channer. This code will block a signal is
		// received.
		s := <-quit

		app.logger.PrintInfo("shutting down server", map[string]string{
			"signal": s.String(),
		})

		// A context with a 20-second timeout.
		ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
		defer cancel()

		// Shutdown() will return nil if the graceful shutdown was successful, or an
		// error. We relay this return value to the shutdownError channel.
		shutdownError <- srv.Shutdown(ctx)

		os.Exit(0)
	}()

	// Start the HTTP server.
	app.logger.PrintInfo("starting server", map[string]string{
		"addr": srv.Addr,
		"env":  app.config.env,
	})

	// Calling Shutdown() on our server will cause ListenAndServe() to immediately
	// return a http.ErrServerClosed error. So if we see this error, it is actually a
	// good thing and an indication that the graceful shutdown has started. So we check
	// specifically for this, only returning the error if it is NOT http.ErrServerClosed.
	err := srv.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	// Otherwise, we wait to receive the return value from Shutdown() on the
	// shutdownError channel. If return value is an error, we know that there was a
	// problem with the graceful shutdown and we return the error.
	err = <-shutdownError
	if err != nil {
		return err
	}

	// At this point we know that the graceful shutdown completed successfully and we
	// log a "stopped server" message.
	app.logger.PrintInfo("stopperd server", map[string]string{
		"addr": srv.Addr,
	})

	return nil
}
