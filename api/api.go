package api

import (
	"log/slog"
	"net/http"
	"wabbajackModlistParser/parser"

	"github.com/go-chi/chi/v5"
)

type Handler struct {
	logger *slog.Logger
}

func NewHandler(logger *slog.Logger) Handler {
	return Handler{
		logger: logger,
	}
}

func (h *Handler) RegisterRoutes(r chi.Router) {
	r.Get("/summary", h.GetModlistSummary)
}

func (h *Handler) handleError(w http.ResponseWriter, r *http.Request, err error) {
	var statusCode int
	var errorMsg string

	errorResponse := ErrorResponse{
		Error:   errorMsg,
		Code:    statusCode,
		Message: err.Error(),
	}
	ErrResponse(w, r, errorResponse)
}

func (h *Handler) GetModlistSummary(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	modlistSummary, err := parser.GetModlistSummary(ctx)
	if err != nil {
		h.handleError(w, r, err)
		return
	}

	JSONResponse(w, http.StatusOK, modlistSummary)
}
