// Package parser this is the parsing package
package parser

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"slices"
	"sort"
	"sync"
	"wabbajackModlistParser/parser/structs"
)

//Refactoring plan
// 1. Move parse, fetch and analyse into it's own packages
// 2. Improve concurrency for fetching archives
// 3. Use interfaces where possible
// 4. Store urls into the struct?

type ModPopularity struct {
	Name  string
	Count int
}

func CreateUrlLinkForApiCall(archiveListPostfix string) string {
	urlPrefix := "https://raw.githubusercontent.com/wabbajack-tools/mod-lists/master/"

	return urlPrefix + archiveListPostfix

}

func GetTopPopularMods(modlists []structs.BaseModlist, n int) []ModPopularity {
	counts := GetModsCountAcrossModpacks(modlists)

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

func GetModsCountAcrossModpacks(modlists []structs.BaseModlist) map[string]int {
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

func ParseMultipleApi(apiUrls []string) []structs.BaseModlist {
	modlists := make([]structs.BaseModlist, 0, len(apiUrls))
	for _, url := range apiUrls {
		modlists = append(modlists, ParseJsonFromApiURL(url, structs.ParseToBaseModlist))
	}
	return modlists
}

func ParseMultipleApiConcurrent(apiUrls []string) []structs.BaseModlist {
	var wg sync.WaitGroup
	modlistsChan := make(chan structs.BaseModlist, len(apiUrls))

	for _, url := range apiUrls {
		u := url // ← затеняем значение переменной цикла чтобы все вызовы получили нужное значение,
		// иначе они потенциально все могут получить последнее значение в range
		wg.Go(func() {
			modlistsChan <- ParseJsonFromApiURL(u, structs.ParseToBaseModlist)
		})
	}
	go func() {
		wg.Wait()
		close(modlistsChan)
	}()

	modlists := make([]structs.BaseModlist, 0, len(apiUrls))
	for i := range modlistsChan {
		modlists = append(modlists, i)
	}

	return modlists
}

func ParseJsonFromApiURL[T any](apiUrl string, parseTo func(jsonData []byte) T) T {
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

	return parseTo(body)
}

func CreateGameModlistTitleMap(apiUrls []structs.Repository, includedGameKeyNames []string) map[string][]string {
	var wg sync.WaitGroup
	modlistsChan := make(chan []structs.ModlistInfo, len(apiUrls))

	for _, obj := range apiUrls {
		u := obj.Link
		wg.Go(func() {
			modlistsChan <- ParseJsonFromApiURL(u, structs.ParseToModlistInfo)
		})
	}
	go func() {
		wg.Wait()
		close(modlistsChan)
	}()

	gameModlistTitleMap := make(map[string][]string, len(apiUrls))
	for range apiUrls {
		// we're making a pass through apiUrls to infer on the game the modpack belongs to
		for _, linkInfo := range <-modlistsChan {
			if slices.Contains(includedGameKeyNames, linkInfo.Game) {
				if slices.Contains(gameModlistTitleMap[linkInfo.Game], linkInfo.Title) {
					continue
				}
				gameModlistTitleMap[linkInfo.Game] = append(gameModlistTitleMap[linkInfo.Game], linkInfo.Title)
			}
		}

	}
	return gameModlistTitleMap
}

func GetModpackArchives(modlistSummary []structs.ModlistSummary, modpackTitle string) structs.BaseModlist {
	var urlLink string
	for _, objModlist := range modlistSummary {
		if objModlist.ModlistName != modpackTitle {
			continue
		}
		urlLink = CreateUrlLinkForApiCall(objModlist.ArchivesLink)
	}

	archiveList := ParseJsonFromApiURL(urlLink, structs.ParseToBaseModlist)

	return archiveList
}

func MainParse(gameNames []string) {
	reposParser := structs.NewReposParser()
	repositories := reposParser.Parse()

	modlistSummaryParser := structs.NewModlistSummaryParser()
	modlistSummary := modlistSummaryParser.Parse()

	includeGames := []string{"skyrimspecialedition", "fallout4"}
	gameModlistTitleMap := CreateGameModlistTitleMap(repositories, includeGames)

	allModlists := make([]structs.BaseModlist, 0, len(gameModlistTitleMap["skyrimspecialedition"]))

	for gameName, modpackTitles := range gameModlistTitleMap {
		if gameName != "skyrimspecialedition" { //temp
			continue
		}
		var wg sync.WaitGroup
		archivesChan := make(chan structs.BaseModlist, len(modpackTitles))

		for _, title := range modpackTitles {
			t := title
			wg.Go(func() {
				archivesChan <- GetModpackArchives(modlistSummary, t)
			})
		}
		go func() {
			wg.Wait()
			close(archivesChan)
		}()

		for i := range archivesChan {
			allModlists = append(allModlists, i)
		}

	}
	modsCount := GetTopPopularMods(allModlists, 100)
	fmt.Printf("modsCount: %v\n", modsCount)
}
