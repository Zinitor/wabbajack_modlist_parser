package structs

import (
	"io"
	"log/slog"
	"wabbajackModlistParser/internal/service/parser/utils"

	"github.com/goccy/go-json"
)

type ModlistSummary struct {
	ModlistName  string `json:"Name"`
	ArchivesLink string `json:"link"`
}

type ModlistSummaryParser struct {
	baseUrl string
}

func NewModlistSummaryParser() *ModlistSummaryParser {
	return &ModlistSummaryParser{
		baseUrl: "https://raw.githubusercontent.com/wabbajack-tools/mod-lists/master/reports/modListSummary.json",
	}
}

func (m *ModlistSummaryParser) Parse() []ModlistSummary {
	responseBody := utils.Fetch(m.baseUrl)
	defer responseBody.Close()
	return m.Transform(responseBody)
}

func (p *ModlistSummaryParser) Transform(r io.Reader) []ModlistSummary {
	var data []ModlistSummary
	err := json.NewDecoder(r).Decode(&data)
	if err != nil {
		slog.Error("json decode failed", slog.Any("err", err))
	}

	return data
}
