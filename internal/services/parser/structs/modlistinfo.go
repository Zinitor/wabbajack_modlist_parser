package structs

import (
	"io"
	"log/slog"

	"github.com/goccy/go-json"
)

type ModlistInfo struct {
	Title string `json:"title"`
	Game  string `json:"game"`
	// потом можно сделать статистику какие моды чаще всего используются в определенных типах модпаков
	// Tags   []string `json:"tags"`
	// IsNSFW bool `json:"nsfw"`
}

func ParseToModlistInfo(r io.Reader) []ModlistInfo {
	var data []ModlistInfo
	err := json.NewDecoder(r).Decode(&data)
	if err != nil {
		slog.Error("json decode failed", slog.Any("err", err))
	}

	return data
}
