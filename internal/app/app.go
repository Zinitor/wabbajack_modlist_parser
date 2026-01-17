package app

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"wabbajackModlistParser/config"
	"wabbajackModlistParser/pkg/httpserver"
	"wabbajackModlistParser/pkg/logger"

	"github.com/go-chi/chi/v5/middleware"
)

func Run(cfg *config.Config) {

	l := logger.New(cfg.Log.Level)

	srv := httpserver.New(l,
		httpserver.Address(cfg.HTTP.Port),
		httpserver.ReadTimeout(10*time.Second),
		httpserver.WriteTimeout(10*time.Second),
	)

	router := srv.Router()
	//TODO move middlewares elsewhere
	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.Compress(5))
	router.Use(middleware.Timeout(60 * time.Second))

	//TODO move registering routes elsewhere
	router.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// Start the server
	srv.Start()

	// Create a channel to listen for OS signals
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	// Block until we receive a signal or server error
	select {
	case err := <-srv.Notify():
		l.Error("Server error: %v", err)
	case sig := <-quit:
		l.Info("Received signal: %s, shutting down...", sig)
	}

	// Shutdown the server gracefully
	_, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(); err != nil {
		l.Error("Error during server shutdown: %v", err)
	}

}
