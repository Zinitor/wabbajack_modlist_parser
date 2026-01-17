package parser_test

import (
	"fmt"
	"testing"

	"wabbajackModlistParser/parser"
	"wabbajackModlistParser/parser/structs"

	"github.com/stretchr/testify/assert"
)

var ApiUrls []string = []string{
	"https://raw.githubusercontent.com/wabbajack-tools/mod-lists/master/reports/Geborgen/nordic-souls/status.json",
	"https://raw.githubusercontent.com/wabbajack-tools/mod-lists/master/reports/HoS/HoS/status.json",
	"https://raw.githubusercontent.com/wabbajack-tools/mod-lists/master/reports/wj-featured/living_skyrim/status.json",
}

// func TestCompareParseFromApiAndFile(t *testing.T) {
// 	apiArchiveSumSize := parser.ParseJsonFromApiURL("https://raw.githubusercontent.com/wabbajack-tools/mod-lists/master/reports/Geborgen/nordic-souls/status.json", structs.ParseToBaseModlist)
// 	localArchiveSumSize := structs.ParseFromFile("archiveData.json")
// 	assert.Equal(t, apiArchiveSumSize, localArchiveSumSize)
// }

// func TestParseMultipleApi(t *testing.T) {
// 	apiArchiveSumSize := parser.ParseMultipleApi(ApiUrls)
// 	assert.NotEmpty(t, apiArchiveSumSize)
// }

// func TestParseMultipleApiConurrent(t *testing.T) {
// 	apiArchiveSumSize := parser.ParseMultipleApiConcurrent(ApiUrls)
// 	assert.NotEmpty(t, apiArchiveSumSize)
// }

func TestCreateUrlLinksForApiCall(t *testing.T) {
	// this creates url links but it doesn't give us the actual game for which modpack it is
	urlLink := parser.CreateURLLinkForAPICall("reports/Wildlander/wildlander/status.json")

	wantLink := "https://raw.githubusercontent.com/wabbajack-tools/mod-lists/master/reports/Wildlander/wildlander/status.json"
	assert.Equal(t, wantLink, urlLink)
}

// func TestGetModlistSummary(t *testing.T) {
// 	modlistSummaryUsingGeneric := parser.ParseJsonFromApiURL("https://raw.githubusercontent.com/wabbajack-tools/mod-lists/master/reports/modListSummary.json", structs.ParseToModlistSummary)
// 	assert.NotEmpty(t, modlistSummaryUsingGeneric)
// }

// func TestStoreModpacksBasedOnGame(t *testing.T) {
// 	urlLinks := parser.CreateUrlLinksForApiCall()

// 	gameModpackMap := parser.CreateModPackMap(urlLinks)

// 	fmt.Printf("gameModpackMap: %v\n", gameModpackMap)

// 	// assert.NotNil(t, urlLinks)

// }

func TestParseJsonToModlistInfo(t *testing.T) {
	modlistInfo := parser.ParseJSONFromAPIUrl("https://raw.githubusercontent.com/tpartridge89/ElderTeej/main/modlists.json", structs.ParseToModlistInfo)
	assert.NotEmpty(t, modlistInfo)
}

func TestParseRepositories(t *testing.T) {
	reposParser := structs.NewReposParser()
	repositories := reposParser.Parse()
	assert.NotEmpty(t, repositories)
}

func TestParseModlistsFromRepositoryLinks(t *testing.T) {
	reposParser := structs.NewReposParser()
	repositories := reposParser.Parse()

	includeGames := []string{"skyrimspecialedition"}
	gameModlistTitleMap := parser.CreateGameModlistTitleMap(repositories, includeGames)

	for _, game := range includeGames {
		assert.Contains(t, gameModlistTitleMap, game, "Expected gameModlistMap to contain modlist for %s", game)
	}

	fmt.Printf("gameModlistMap: %v\n", gameModlistTitleMap)
}

// func TestGetModpackArchives(t *testing.T) {
// 	modlistSummary := parser.ParseJsonFromApiURL("https://raw.githubusercontent.com/wabbajack-tools/mod-lists/master/reports/modListSummary.json", structs.ParseToModlistSummary)

// 	modpackTitle := "Skyrim Modding Essentials"

// 	archives := parser.GetModpackArchives(modlistSummary, modpackTitle)

// 	fmt.Printf("archives: %v\n", archives)

// 	assert.NotEmpty(t, archives)

// }

func TestGetAllGameModpackArchives(t *testing.T) {
	includeGames := []string{"skyrimspecialedition"}

	parser.MainParse(includeGames)
}

func BenchmarkMainParse(b *testing.B) {
	includeGames := []string{"skyrimspecialedition"}

	for b.Loop() {
		parser.MainParse(includeGames)
	}
}
