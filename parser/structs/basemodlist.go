package structs

import (
	"encoding/json"
	"log/slog"
	"os"
)

type State struct {
	Name    string `json:"Name"`
	Version string `json:"Version"`
}

type ArchiveData struct {
	Name  string `json:"Name"`
	Size  int    `json:"Size"`
	State State  `json:"State"`
}

type Archive struct {
	Status      string      `json:"Status"`
	ArchiveData ArchiveData `json:"Original"`
}

type BaseModlist struct {
	MachineURL     string    `json:"MachineURL"`
	ModlistName    string    `json:"Name"`
	ModListVersion string    `json:"Version"`
	Archives       []Archive `json:"Archives"`
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
