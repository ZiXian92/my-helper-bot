package main

import (
	"context"
	"io"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
	"time"

	"github.com/ZiXian92/my-helper-bot/plugins"
	"github.com/ZiXian92/my-helper-bot/plugins/interfaces"
	"github.com/gorilla/mux"
	"github.com/hashicorp/go-plugin"
	"github.com/rs/zerolog/log"
)

func main() {
	// Initialize the logger
	log.Info().Msg("Starting my-helper-bot")

	// Create a context for graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Load plugins
	// TODO: Implement a loop to load all enabled plugins found in the application's
	// installed plugins folder.
	healthCheckPluginClient := plugin.NewClient(&plugin.ClientConfig{
		HandshakeConfig:  plugins.HandshakeConfig,
		Plugins:          plugins.PluginMap,
		Cmd:              exec.Command("./healthcheck"),
		AllowedProtocols: []plugin.Protocol{plugin.ProtocolGRPC},
	})
	defer healthCheckPluginClient.Kill()

	// Create a router using gorilla/mux
	router := mux.NewRouter()

	// Start an HTTP server
	httpServer := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	// Process plugins
	healthCheckDispenser, err := healthCheckPluginClient.Client()
	if err != nil {
		log.Warn().Err(err).Msg("Failed to get healthcheck plugin dispenser.")
	} else {
		webHandlerPlugin, err := healthCheckDispenser.Dispense(plugins.PluginMapKeyWebHandler)
		if err != nil {
			log.Warn().Err(err).Msg("Failed to get dispense web handler plugin.")
		} else if webHandler, ok := webHandlerPlugin.(interfaces.WebHandler); !ok {
			log.Warn().Msg("Web handler plugin does not implement web handler interface.")
		} else {
			eps := webHandler.GetEndpoints()
			log.Info().Any("web endpoints", eps).Msg("Found web endpoints.")
			for _, ep := range eps {
				router.HandleFunc(ep.Path, func(w http.ResponseWriter, r *http.Request) {
					rBodyBytes, err := io.ReadAll(r.Body)
					if err != nil {
						log.Err(err).Msg("Error reading HTTP request body.")
						http.Error(w, "Error processing request.", http.StatusInternalServerError)
						return
					}

					res := webHandler.HandleRequest(interfaces.WebRequest{
						EndPointName: ep.Name,
						Headers:      r.Header,
						URIParams:    mux.Vars(r),
						QueryParams:  r.URL.Query(),
						Body:         rBodyBytes,
					})
					for k, vList := range res.Headers {
						for _, v := range vList {
							w.Header().Add(k, v)
						}
					}
					w.Write(res.Body)
					w.WriteHeader(res.Code)
				}).Methods(ep.Methods...)
			}
		}
	}
	defer healthCheckDispenser.Close()

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
