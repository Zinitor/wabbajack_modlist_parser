package structs

import (
	"log/slog"
	"os"

	"github.com/goccy/go-json"
)

type State struct {
	Name string `json:"Name"`
	// Version string `json:"Version"`
}

type ArchiveData struct {
	State State `json:"State"`
}

type Archive struct {
	ArchiveData ArchiveData `json:"Original"`
}

type BaseModlist struct {
	Archives []Archive `json:"Archives"`
}

func ParseToBaseModlist(jsonData []byte) BaseModlist {
	var parsedData BaseModlist
	err := json.Unmarshal(jsonData, &parsedData)
	if err != nil {
		slog.Error("unmarshal err", slog.Any("err", err))
	}

	return parsedData
}

func ParseFromFile(filename string) BaseModlist {
	rawData, err := os.ReadFile(filename)
	if err != nil {
		slog.Error("read file err", slog.Any("err", err))
	}
	return ParseToBaseModlist(rawData)
}
