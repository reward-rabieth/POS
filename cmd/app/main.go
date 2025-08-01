package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"

	"syscall"

	_ "github.com/rustacean-dev/possystem/docs"
	"github.com/rustacean-dev/possystem/http"
	"golang.org/x/sync/errgroup"
	"maragu.dev/env"
)

// @title POS System API
// @version 1.0
// @description Modern POS built with PocketBase, Gomponents, HTMX
// @contact.name Rustacean
// @host localhost:8080
// @BasePath /

func main() {
	// Set up a logger that is used throughout the app
	log := slog.New(slog.NewTextHandler(os.Stderr, nil))

	// Start the app, exit with a non-zero exit code on errors
	if err := start(log); err != nil {
		log.Error("Error starting app", "error", err)
		os.Exit(1)
	}
}

func start(log *slog.Logger) error {
	log.Info("Starting app")

	// We load environment variables from .env if it exists
	_ = env.Load()

	// Catch signals to gracefully shut down the app
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer stop()

	// Set up the HTTP server, injecting the database and logger
	s := http.NewServer(http.NewServerOptions{
		Log: log,
	})

	// Use an errgroup to wait for separate goroutines which can error
	eg, ctx := errgroup.WithContext(ctx)

	// Start the server within the errgroup.
	// You can do this for other dependencies as well.
	eg.Go(func() error {
		return s.Start()
	})

	// Wait for the context to be done, which happens when a signal is caught
	<-ctx.Done()
	log.Info("Stopping app")

	// Stop the server gracefully
	eg.Go(func() error {
		return s.Stop()
	})

	// Wait for the server to stop
	if err := eg.Wait(); err != nil {
		return err
	}

	log.Info("Stopped app")

	return nil
}
