// Package v1 implements routing paths. Each services in own file.
package restapi

import (
	"net/http"
	"time"
	"wabbajackModlistParser/config"
	v1 "wabbajackModlistParser/internal/controller/restapi/v1"
	"wabbajackModlistParser/pkg/logger"

	_ "wabbajackModlistParser/docs" // This imports your generated docs

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	httpSwagger "github.com/swaggo/http-swagger"
)

// NewRouter -.
func NewRouter(router chi.Router, cfg *config.Config, l logger.Interface) {
	//TODO move elsewhere
	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.Compress(5))
	router.Use(middleware.Timeout(60 * time.Second))

	// Root endpoints (not under /api/v1)
	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Wabbajack Modlist Parser API"))
	})

	router.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status": "ok"}`))
	})

	// API v1 routes
	router.Route("/api/v1", func(r chi.Router) {
		v1.NewTestRoutes(r, l)
	})

	if cfg.Swagger.Enabled {
		router.Get("/swagger/*", httpSwagger.Handler(
			httpSwagger.URL("/swagger/doc.json"),
		))
		l.Info("Swagger UI enabled at /swagger/")
	}
}
