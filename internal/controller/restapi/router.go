// Package v1 implements routing paths. Each services in own file.
package restapi

import (
	"net/http"
	"time"
	"wabbajackModlistParser/config"
	"wabbajackModlistParser/pkg/logger"

	_ "wabbajackModlistParser/docs" // This imports your generated docs

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	httpSwagger "github.com/swaggo/http-swagger"
)

// NewRouter -.
// Swagger spec:
// @title       Go Clean Template API
// @description Using a translation service as an example
// @version     1.0
// @host        localhost:8080
// @BasePath    /v1
func NewRouter(router chi.Router, cfg *config.Config, l logger.Interface) {

	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.Compress(5))
	router.Use(middleware.Timeout(60 * time.Second))

	// Swagger
	// Swagger UI (only if enabled in config)
	if cfg.Swagger.Enabled {
		// Swagger UI will be available at /swagger/
		router.Get("/swagger/*", httpSwagger.Handler(
			httpSwagger.URL("/swagger/doc.json"), // The url pointing to API definition
		))
		l.Info("Swagger UI enabled at /swagger/")
	}

	// Routers
	router.Group(func(r chi.Router) {
		router.Get("/health", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("OK"))
		})
	})

	// API v1 routes
	router.Route("/api/v1", func(r chi.Router) {
		// API-specific middleware can go here
		r.Use(middleware.AllowContentType("application/json"))

		// Your API endpoints will go here
		r.Get("/status", apiStatus)
		// r.Get("/modlists", getModlists)
		// r.Post("/parse", parseModlist)
	})

}

// Health check endpoint
// @Summary Health check
// @Description Check if the API is running
// @Tags health
// @Accept json
// @Produce json
// @Success 200 {object} map[string]string
// @Router /health [get]
func healthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status": "ok", "timestamp": "` + time.Now().Format(time.RFC3339) + `"}`))
}

// API status endpoint
// @Summary API status
// @Description Get API version and status
// @Tags api
// @Accept json
// @Produce json
// @Success 200 {object} map[string]string
// @Router /api/v1/status [get]
func apiStatus(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"api": "v1", "status": "running", "version": "1.0.0"}`))
}
