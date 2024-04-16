package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/ucho456job/lgtmeme/internal/setup"
)

func main() {
	e := setup.SetupServer()
	defer setup.CloseConnection()

	// Graceful shutdown
	// Wait for interrupt signal to gracefully shutdown the server with a timeout of 10 seconds.
	// Use a buffered channel to avoid missing signals as recommended for signal.Notify
	port := ":" + os.Getenv("PORT")
	go func() {
		if err := e.Start(port); err != nil && err != http.ErrServerClosed {
			e.Logger.Fatal("shutting down the server")
		}
	}()
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}
}
