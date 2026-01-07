package structs

import (
	"encoding/json"
	"log/slog"
	"wabbajackModlistParser/parser/utils"
)

type ModlistSummary struct {
	ModlistName  string `json:"Name"`
	MachineUrl   string `json:"MachineUrl"`
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

	return m.Transform(responseBody)
}

func (m *ModlistSummaryParser) Transform(jsonData []byte) []ModlistSummary {
	var parsedData []ModlistSummary
	err := json.Unmarshal(jsonData, &parsedData)
	if err != nil {
		slog.Error("unmarshal err", slog.Any("err", err))
	}

	return parsedData
}
