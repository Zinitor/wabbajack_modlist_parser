// Package v1 provides version 1 endpoints
package v1

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
	"wabbajackModlistParser/internal/services/modlist"
	"wabbajackModlistParser/pkg/logger"

	"github.com/go-chi/chi/v5"
)

// V1 -.
type V1 struct {
	l logger.Interface
}

// RegisterRoutes -.
func RegisterRoutes(apiV1Group chi.Router, l logger.Interface) {
	h := &V1{l: l}

	apiV1Group.Get("/status", h.apiStatus)
	apiV1Group.Get("/health", h.healthCheck)
	apiV1Group.Get("/modlists", h.getModlists)
}

// Health check endpoint
//
//	@Summary		Health check
//	@Description	Check if the API is running
//	@Tags			health
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	map[string]string
//	@Router			/api/v1/health [get]
func (h *V1) healthCheck(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status": "ok", "timestamp": "` + time.Now().Format(time.RFC3339) + `"}`))
}

// API status endpoint
//
//	@Summary		API status
//	@Description	Get API version and status
//	@Tags			api
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	map[string]string
//	@Router			/api/v1/status [get]
func (h *V1) apiStatus(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"api": "v1", "status": "running", "version": "1.0.0"}`))
}

// Get all modlists
//
//	@Summary		Get all modlists
//	@Description	Retrieve a list of all available modlists with their details
//	@Tags			modlists
//	@Accept			json
//	@Produce		json
//	@Success		200	{array}		ModlistSummaryResponse
//	@Failure		500	{object}	map[string]string	"Internal server error"
//	@Router			/api/v1/modlists [get]
func (h *V1) getModlists(w http.ResponseWriter, _ *http.Request) {
	service := modlist.NewModlistService(h.l)
	modlists, err := service.GetModlists(context.TODO())

	if err != nil {
		http.Error(w, fmt.Sprintf("fetcher failed: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if err = json.NewEncoder(w).Encode(modlists); err != nil {
		http.Error(w, fmt.Sprintf("JSON encode failed: %v", err), http.StatusInternalServerError)
		return
	}
}
