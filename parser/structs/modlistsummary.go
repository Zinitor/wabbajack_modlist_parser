package structs

import (
	"encoding/json"
	"log/slog"
)

type ModlistSummary struct {
	ModlistName  string `json:"Name"`
	MachineUrl   string `json:"MachineUrl"`
	ArchivesLink string `json:"link"`
}

func ParseToModlistSummary(jsonData []byte) []ModlistSummary {
	var parsedData []ModlistSummary
	err := json.Unmarshal(jsonData, &parsedData)
	if err != nil {
		slog.Error("unmarshal err", slog.Any("err", err))
	}

	return parsedData
}
