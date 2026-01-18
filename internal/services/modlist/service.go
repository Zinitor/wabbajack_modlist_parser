package modlist

import (
	"context"
	"net/http"
	"time"
	"wabbajackModlistParser/pkg/logger"
)

type Service struct {
	l logger.Interface
}

type Summary struct {
	ModlistName  string `json:"Name"`
	ArchivesLink string `json:"link"`
}

func NewModlistService(logger logger.Interface) Service {
	return Service{l: logger}
}

var DefaultTimeout time.Duration = 30 * time.Second

func (m *Service) GetModlists(ctx context.Context) ([]Summary, error) {
	modlists := make([]Summary, 0)
	restClient := http.Client{Timeout: DefaultTimeout}
	uri := "https://raw.githubusercontent.com/wabbajack-tools/mod-lists/master/reports/modListSummary.json"

	modlists, err := fetchAndParse[[]Summary](ctx, &restClient, uri)
	if err != nil {
		return modlists, err
	}

	return modlists, nil
}
