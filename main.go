package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
)

func main() {
	// Initialize the logger
	log.Info().Msg("Starting my-helper-bot")

	// Create a context for graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Create a router using gorilla/mux
	router := mux.NewRouter()
	router.HandleFunc("/health", healthHandler)

	// Start an HTTP server
	httpServer := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	go func() {
		log.Info().Msg("Starting HTTP server")
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal().Err(err).Msg("HTTP server stopped unexpectedly")
		}
	}()

	// Listen for OS signals to trigger graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	select {
	case sig := <-sigChan:
		log.Info().Msgf("Received signal: %s. Shutting down gracefully...", sig)
		cancel()
		// Close any resources related to your plugins here

		// Wait for a grace period before shutting down
		shutdownTimeout := 10 * time.Second
		ctx, cancel := context.WithTimeout(ctx, shutdownTimeout)
		defer cancel()

		if err := httpServer.Shutdown(ctx); err != nil {
			log.Fatal().Err(err).Msg("Error during HTTP server shutdown")
		}
	case <-ctx.Done():
		// Your plugin-specific shutdown logic can go here
	}

	log.Info().Msg("Your Go application has gracefully stopped")
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "OK")
}
