package structs

import (
	"encoding/json"
	"log/slog"
	"wabbajackModlistParser/parser/utils"
)

type Repository struct {
	Name string
	Link string
}

type ReposParser struct {
	baseUrl string
}

func NewReposParser() *ReposParser {
	return &ReposParser{
		baseUrl: "https://raw.githubusercontent.com/wabbajack-tools/mod-lists/master/repositories.json",
	}
}

func (r *ReposParser) Parse() []Repository {
	responseBody := utils.Fetch(r.baseUrl)
	return r.Transform(responseBody)
}

func (r *ReposParser) Transform(jsonData []byte) []Repository {
	var result map[string]string
	err := json.Unmarshal(jsonData, &result)
	if err != nil {
		slog.Error("unmarshal err", slog.Any("err", err))
	}

	parsedData := make([]Repository, 0, len(result))

	for name, link := range result {
		parsedData = append(parsedData,
			Repository{
				Name: name,
				Link: link,
			},
		)

	}

	return parsedData
}
