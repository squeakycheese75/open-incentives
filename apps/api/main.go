package main

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	engine "github.com/squeakycheese75/open-incentives-engine"
	"github.com/squeakycheese75/open-incentives/internal/handlers"
)

func main() {
	engine := engine.New()
	h := handlers.NewHandlers(engine)

	mux := http.NewServeMux()
	mux.HandleFunc("POST /v1/evaluate", h.Evaluate)
	mux.HandleFunc("GET /health", h.Health)

	srv := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	go func() {
		log.Println("api listening on :8080")

		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 2 seconds.
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	slog.Warn("Received request to shutdown server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		slog.Error("server shutdown failed", "error", err)
	}

	<-ctx.Done()
	slog.Warn("Server exiting ...")
}
