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

//TODO
// 1. Разобраться с логгерами
// 2. Вынести трансформеры в отдельную подпапку вместе с структурой в которую они приводят
// 3. Разобраться с юрлками
// 4. Написать мейн

type ModPopularity struct {
	Name  string
	Count int
}

func CreateUrlLinksForApiCall() []string {
	modlists := ParseJsonFromApiURL("https://raw.githubusercontent.com/wabbajack-tools/mod-lists/master/reports/modListSummary.json", structs.ParseToModlistSummary)
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

func CreateModPackMap(apiUrls []structs.Repository, includedGameKeyNames []string) map[string][]string {
	gameModlistMap := make(map[string][]string, len(apiUrls))
	//тупая версия
	for idx, obj := range apiUrls {
		info := ParseJsonFromApiURL(obj.Link, structs.ParseToModlistInfo)
		for _, linkInfo := range info {
			if slices.Contains(includedGameKeyNames, linkInfo.Game) {
				gameModlistMap[linkInfo.Game] = append(gameModlistMap[linkInfo.Game], obj.Link)
			}
		}
		fmt.Printf("gameModpackMap: %v\n", idx)

	}

	//конкурентная версия
	//

	// for _, mInfo := range modlistsInfo {
	// if gameModlistMap[mInfo.Game]
	// }
	return gameModlistMap
}

// Забрать repositories.json,
// Пробежаться по каждой из представленных ссылок
// При парсинге оставлять только те где game == переменная
// Так собираем только модпаки для нужной игры а дальше уже разберемся как получать данные архивов
