package structs

import (
	"io"
	"log"
	"log/slog"
	"os"

	"github.com/goccy/go-json"
)

type State struct {
	Name string `json:"Name"`
}

type Original struct {
	State State `json:"State"`
}

type Archive struct {
	Original Original `json:"Original"`
}

type BaseModlist struct {
	Archives []Archive `json:"Archives"`
	Name     string    `json:"Name"`
}

func ParseToBaseModlist(jsonData []byte) BaseModlist {
	var parsedData BaseModlist
	err := json.Unmarshal(jsonData, &parsedData)
	if err != nil {
		slog.Error("unmarshal err", slog.Any("err", err))
	}

	return parsedData
}

func ParseToModlistArchiveMap(r io.Reader) map[string]int {
	myMap := make(map[string]int, 0)
	decoder := json.NewDecoder(r)

	var fullModlist BaseModlist
	err := decoder.Decode(&fullModlist)
	if err != nil {
		log.Fatal(err)
	}

	for _, archive := range fullModlist.Archives {
		if archive.Original.State.Name != "" {
			myMap[archive.Original.State.Name] += 1
		}
	}

	return myMap
}

func ParseFromFile(filename string) BaseModlist {
	rawData, err := os.ReadFile(filename)
	if err != nil {
		slog.Error("read file err", slog.Any("err", err))
	}
	return ParseToBaseModlist(rawData)
}
