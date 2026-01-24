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
	l       logger.Interface
	service *modlist.Service
}

func NewV1(l logger.Interface, service *modlist.Service) *V1 {
	return &V1{
		l:       l,
		service: service,
	}
}

// RegisterRoutes -.
func RegisterRoutes(r chi.Router, h *V1) {
	r.Get("/status", h.apiStatus)
	r.Get("/health", h.healthCheck)
	r.Get("/modlists", h.getModlists)
	r.Get("/repositories", h.getRepos)
	r.Get("/games", h.getAllGames)
	r.Get("/games/top-popular", h.getTopPopularByModlistsCount)
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
	modlists, err := h.service.GetModlistSummary(context.TODO())

	if err != nil {
		http.Error(w, fmt.Sprintf("service failed: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if err = json.NewEncoder(w).Encode(modlists); err != nil {
		http.Error(w, fmt.Sprintf("JSON encode failed: %v", err), http.StatusInternalServerError)
		return
	}
}

// Get all user repositories
//
//	@Summary		Get all repositories
//	@Description	Retrieve a list of all available repositories with their details
//	@Tags			repositories
//	@Accept			json
//	@Produce		json
//	@Success		200	{array}		RepositoryResponse
//	@Failure		500	{object}	map[string]string	"Internal server error"
//	@Router			/api/v1/repositories [get]
func (h *V1) getRepos(w http.ResponseWriter, _ *http.Request) {
	repos, err := h.service.GetUserRepos(context.TODO())

	if err != nil {
		http.Error(w, fmt.Sprintf("service failed: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err = json.NewEncoder(w).Encode(repos); err != nil {
		http.Error(w, fmt.Sprintf("JSON encode failed: %v", err), http.StatusInternalServerError)
		return
	}
}

// Get all games repositories
//
//	@Summary		Get all games
//	@Description	Retrieve a list of all available games
//	@Tags			games
//	@Accept			json
//	@Produce		json
//	@Success		200	{array}		string
//	@Failure		500	{object}	map[string]string	"Internal server error"
//	@Router			/api/v1/games [get]
func (h *V1) getAllGames(w http.ResponseWriter, _ *http.Request) {
	modlists, err := h.service.GetAllGamesFromModlists(context.TODO())

	if err != nil {
		http.Error(w, fmt.Sprintf("service failed: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err = json.NewEncoder(w).Encode(modlists); err != nil {
		http.Error(w, fmt.Sprintf("JSON encode failed: %v", err), http.StatusInternalServerError)
		return
	}
}

// Get top popular games by modpack count
//
//	@Summary		Get top popular games
//	@Description	Retrieves a list of games ranked by popularity based on the number of modpacks/modlists available for each game. The popularity score is determined by counting how many modpacks reference each game across all repositories.
//	@Tags			games
//	@Accept			json
//	@Produce		json
//	@Param			limit	query		int				false	"Number of top games to return (default: 10, max: 100)"
//	@Param			sort	query		string			false	"Sort order: 'asc' for ascending (least popular first), 'desc' for descending (most popular first) (default: 'desc')"
//	@Success		200		{array}		GamePopularity	"List of games with their popularity scores"
//	@Failure		400		{object}	ErrorResponse	"Invalid request parameters"
//	@Failure		500		{object}	ErrorResponse	"Internal server error"
//	@Failure		502		{object}	ErrorResponse	"Service temporarily unavailable"
//	@Router			/api/v1/games/top-popular [get]
//	@Example		curl -X GET "https://api.example.com/api/v1/games/top-popular?limit=5&sort=desc"
func (h *V1) getTopPopularByModlistsCount(w http.ResponseWriter, _ *http.Request) {
	modlists, err := h.service.GetTopPopularGames(context.TODO(), 10, "desc")

	if err != nil {
		http.Error(w, fmt.Sprintf("service failed: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err = json.NewEncoder(w).Encode(modlists); err != nil {
		http.Error(w, fmt.Sprintf("JSON encode failed: %v", err), http.StatusInternalServerError)
		return
	}
}
