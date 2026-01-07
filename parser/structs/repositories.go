package structs

import (
	"encoding/json"
	"log/slog"
)

type Repository struct {
	Name string
	Link string
}

func ParseToRepos(jsonData []byte) []Repository {

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
