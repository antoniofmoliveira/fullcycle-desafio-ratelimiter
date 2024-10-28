package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/antoniofmoliveira/fullcycle-desafio-ratelimiter/internal/config"
	"github.com/antoniofmoliveira/fullcycle-desafio-ratelimiter/internal/server"
)

// Main runs the server and waits for a signal to shutdown.
// It then calls the Shutdown method on the server with a 10 second timeout.
// If the shutdown fails, it logs an error and exits with a non-zero status code.
func main() {

	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, nil)))

	c := config.NewConfig()

	server := server.NewServer(*c)
	go server.Start()

	termChan := make(chan os.Signal, 1)
	signal.Notify(termChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)

	<-termChan
	slog.Info("server: shutting down")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := server.Server.Shutdown(ctx); err != nil {
		slog.Error("Could not shutdown the server", "error", err)
	}
	slog.Info("Server stopped")
	os.Exit(0)
}
