// Package parser this is the parsing package
package parser

import (
	"io"
	"log"
	"net/http"
	"slices"
	"sort"
	"sync"

	"wabbajackModlistParser/parser/structs"
)

// Refactoring plan
// 1. Move parse, fetch and analyse into it's own packages
// 2. Improve concurrency for fetching archives
// 3. Use interfaces where possible
// 4. Store urls into the struct?

type ModPopularity struct {
	Name  string
	Count int
}

func CreateURLLinkForAPICall(archiveListPostfix string) string {
	urlPrefix := "https://raw.githubusercontent.com/wabbajack-tools/mod-lists/master/"

	return urlPrefix + archiveListPostfix
}

func GetTopPopularMods(modlists map[string]int, n int) []ModPopularity {
	popularity := make([]ModPopularity, 0, len(modlists))
	for modName, count := range modlists {
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

func ParseJSONFromAPIUrl[T any](apiURL string, parseTo func(r io.Reader) T) T {
	response, err := http.Get(apiURL)
	if err != nil {
		defer response.Body.Close()
		log.Fatal(err)
	}

	if response.StatusCode != http.StatusOK {
		defer response.Body.Close()
		log.Fatalf("API request failed with status: %d", response.StatusCode)
	}

	return parseTo(response.Body)
}

func CreateGameModlistTitleMap(apiUrls []structs.Repository, includedGameKeyNames []string) map[string][]string {
	var wg sync.WaitGroup
	modlistsChan := make(chan []structs.ModlistInfo, len(apiUrls))

	for _, obj := range apiUrls {
		u := obj.Link
		wg.Go(func() {
			modlistsChan <- ParseJSONFromAPIUrl(u, structs.ParseToModlistInfo)
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

func GetModpackArchives(modlistSummary []structs.ModlistSummary, modpackTitle string) map[string]int {
	var urlLink string
	for _, objModlist := range modlistSummary {
		if objModlist.ModlistName != modpackTitle {
			continue
		}
		urlLink = CreateURLLinkForAPICall(objModlist.ArchivesLink)
	}

	archiveList := ParseJSONFromAPIUrl(urlLink, structs.ParseToModlistArchiveMap)

	return archiveList
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

func MainParse(gameNames []string) {
	reposParser := structs.NewReposParser()
	reposCh := GetBase(reposParser)
	repositories := <-reposCh

	modlistSummaryParser := structs.NewModlistSummaryParser()
	summaryCh := GetBase(modlistSummaryParser)
	modlistSummary := <-summaryCh

	gameModlistTitlesMap := CreateGameModlistTitleMap(repositories, gameNames)
	var allModlistsLen int
	for key := range gameModlistTitlesMap {
		allModlistsLen = len(gameModlistTitlesMap[key])
	}

	allModlistsMap := make(map[string]int, allModlistsLen)

	for _, modpackTitles := range gameModlistTitlesMap {
		// wg внутри цикла чтобы мы собирали модпаки по каждой игре,
		var wg sync.WaitGroup
		archivesChan := make(chan map[string]int, len(modpackTitles))

		for _, title := range modpackTitles {
			// поскольку я ожидаю пока все горутины закончат работу и спарсят свои данные то здесь и происходит горлышко по памяти
			// ведь я жду пока я получу все архивы, а мне вовсе необязательно это делать
			t := title
			wg.Go(func() {
				archivesChan <- GetModpackArchives(modlistSummary, t)
			})
		}
		go func() {
			wg.Wait()
			close(archivesChan)
		}()

		for mObj := range archivesChan {
			for modName, quantity := range mObj {
				allModlistsMap[modName] += quantity
			}
		}

	}
	_ = GetTopPopularMods(allModlistsMap, 100)
	// fmt.Printf("modsCount: %v\n", modsCount)
}
