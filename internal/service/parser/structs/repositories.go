package structs

import (
	"io"
	"log/slog"
	"wabbajackModlistParser/internal/service/parser/utils"

	"github.com/goccy/go-json"
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
	defer responseBody.Close()
	return r.Transform(responseBody)
}

func (r *ReposParser) Transform(reader io.Reader) []Repository {
	var result map[string]string
	err := json.NewDecoder(reader).Decode(&result)
	if err != nil {
		slog.Error("JSON decode failed", slog.Any("err", err))
		return nil
	}

	parsedData := make([]Repository, 0, len(result))
	for name, link := range result {
		parsedData = append(parsedData, Repository{Name: name, Link: link})
	}
	return parsedData
}
