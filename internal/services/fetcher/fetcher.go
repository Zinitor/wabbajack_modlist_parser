package fetcher

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"log/slog"
	"net/http"
)

type ModlistSummary struct {
	ModlistName  string `json:"Name"`
	ArchivesLink string `json:"link"`
}

type ModlistSummaryParser struct {
	baseUrl string
}

type Parser[T any] interface {
	Parse() []T
}

func GetBase[T any](p Parser[T]) <-chan []T {
	ch := make(chan []T, 1)
	go func() {
		defer close(ch)
		ch <- p.Parse()
	}()
	return ch
}

func NewModlistSummaryParser() *ModlistSummaryParser {
	return &ModlistSummaryParser{
		baseUrl: "https://raw.githubusercontent.com/wabbajack-tools/mod-lists/master/reports/modListSummary.json",
	}
}
func (m *ModlistSummaryParser) Parse() []ModlistSummary {
	responseBody := Fetch(m.baseUrl)
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

func Fetch(baseUrl string) io.ReadCloser {
	response, err := http.Get(baseUrl)
	if err != nil {
		log.Fatal(err)
	}

	if response.StatusCode != http.StatusOK {
		response.Body.Close()
		log.Fatalf("API request failed with status: %d", response.StatusCode)
	}

	return response.Body
}

func GetModlistSummary(ctx context.Context) ([]ModlistSummary, error) {
	modlistSummaryParser := NewModlistSummaryParser()
	summaryCh := GetBase(modlistSummaryParser)
	modlistSummary := <-summaryCh
	return modlistSummary, nil
}
