package structs

import (
	"encoding/json"
	"log/slog"
)

type ModlistInfo struct {
	Title string `json:"title"`
	Game  string `json:"game"`
	//потом можно сделать статистику какие моды чаще всего используются в определенных типах модпаков
	Tags   []string `json:"tags"`
	IsNSFW bool     `json:"nsfw"`
}

func ParseToModlistInfo(jsonData []byte) []ModlistInfo {
	var parsedData []ModlistInfo
	err := json.Unmarshal(jsonData, &parsedData)
	if err != nil {
		slog.Error("unmarshal err", slog.Any("err", err))
	}

	return parsedData

}
