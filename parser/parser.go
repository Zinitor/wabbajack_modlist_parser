// Package parser this is the parsing package
package parser

import (
	"encoding/json"
	"io"
	"log"
	"log/slog"
	"net/http"
	"os"
	"sort"
	"sync"
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

type ModPopularity struct {
	Name  string
	Count int
}

type ModlistSummary struct {
	ModlistName string `json:"Name"`
	MachineUrl  string `json:"MachineUrl"`
}

func ParseJsonToModlistSummary(jsonData []byte) []ModlistSummary {
	var parsedData []ModlistSummary
	err := json.Unmarshal(jsonData, &parsedData)
	if err != nil {
		slog.Error("unmarshal err", slog.Any("err", err))
	}

	return parsedData

}

func GetModlistSummary() []ModlistSummary {
	modlistSummaryLink := "https://raw.githubusercontent.com/wabbajack-tools/mod-lists/master/reports/modListSummary.json"
	response, err := http.Get(modlistSummaryLink)
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		log.Fatalf("API request failed with status: %d", response.StatusCode)
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}

	return ParseJsonToModlistSummary(body)
}

func CreateUrlLinksForApiCall() []string {
	modlists := GetModlistSummary()
	urlPrefix := "https://raw.githubusercontent.com/wabbajack-tools/mod-lists/master/reports/"
	urlPostfix := "/status.json"

	urlLinks := make([]string, 0, len(modlists))

	for _, modpack := range modlists {
		archiveSearchString := urlPrefix + modpack.MachineUrl + urlPostfix
		urlLinks = append(urlLinks, archiveSearchString)

	}
	return urlLinks
}

func GetTopPopularMods(apiUrls []string, n int) []ModPopularity {
	counts := GetModsCountAcrossModpacks(apiUrls)

	popularity := make([]ModPopularity, 0, len(counts))
	for modName, count := range counts {
		popularity = append(popularity, ModPopularity{Name: modName, Count: count})
	}

	sort.Slice(popularity, func(i, j int) bool {
		return popularity[i].Count > popularity[j].Count
	})

	if n > len(popularity) {
		n = len(popularity)
	}
	return popularity[:n]
}

func GetModsCountAcrossModpacks(apiUrls []string) map[string]int {
	modlists := ParseMultipleApi(apiUrls)

	ModsCountMap := make(map[string]int)

	for _, mod := range modlists {
		for _, archive := range mod.Archives {
			modName := archive.ArchiveData.State.Name
			if modName == "" {
				continue
			}
			ModsCountMap[modName]++
		}
	}
	return ModsCountMap
}

func ParseMultipleApi(apiUrls []string) []BaseModlist {
	modlists := make([]BaseModlist, 0, len(apiUrls))
	for _, url := range apiUrls {
		modlists = append(modlists, ParseJsonFromApiURL(url))
	}
	return modlists
}

func ParseMultipleApiConcurrent(apiUrls []string) []BaseModlist {
	var wg sync.WaitGroup
	modlistsChan := make(chan BaseModlist, len(apiUrls))

	for _, url := range apiUrls {
		u := url // ← затеняем значение переменной цикла чтобы все вызовы получили нужное значение,
		// иначе они потенциально все могут получить последнее значение в range
		wg.Go(func() {
			modlistsChan <- ParseJsonFromApiURL(u)
		})
	}
	go func() {
		wg.Wait()
		close(modlistsChan)
	}()

	modlists := make([]BaseModlist, 0, len(apiUrls))
	for i := range modlistsChan {
		modlists = append(modlists, i)
	}

	return modlists
}

func ParseJsonFromApiURL(apiUrl string) BaseModlist {
	response, err := http.Get(apiUrl)
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		log.Fatalf("API request failed with status: %d", response.StatusCode)
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}

	return ParseJSONToBaseModlist(body)
}

func ParseJsonFromFile(filename string) BaseModlist {
	rawData, err := os.ReadFile(filename)
	if err != nil {
		slog.Error("read file err", slog.Any("err", err))
	}
	return ParseJSONToBaseModlist(rawData)
}

func ParseJSONToBaseModlist(jsonData []byte) BaseModlist {
	var parsedData BaseModlist
	err := json.Unmarshal(jsonData, &parsedData)
	if err != nil {
		slog.Error("unmarshal err", slog.Any("err", err))
	}

	var totalSize int
	for _, archive := range parsedData.Archives {
		totalSize += archive.ArchiveData.Size
	}
	return parsedData
}
