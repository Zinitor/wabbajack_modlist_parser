// File: internal/controller/restapi/v1/v1.go
package v1

import (
	"net/http"
	"time"
	"wabbajackModlistParser/pkg/logger"

	"github.com/go-chi/chi/v5"
)

// V1 -.
type V1 struct {
	l logger.Interface
}

// NewTestRoutes -.
func NewTestRoutes(apiV1Group chi.Router, l logger.Interface) {
	h := &V1{l: l}

	apiV1Group.Get("/status", h.apiStatus)
	apiV1Group.Get("/health", h.healthCheck) // Changed from Post to Get
}

// Health check endpoint
// @Summary Health check
// @Description Check if the API is running
// @Tags health
// @Accept json
// @Produce json
// @Success 200 {object} map[string]string
// @Router /api/v1/health [get]
func (h *V1) healthCheck(w http.ResponseWriter, r *http.Request) {
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
func (h *V1) apiStatus(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"api": "v1", "status": "running", "version": "1.0.0"}`))
}
