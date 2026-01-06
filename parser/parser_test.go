package parser_test

import (
	"testing"

	"wabbajackModlistParser/parser"

	"github.com/stretchr/testify/assert"
)

var ApiUrls []string = []string{
	"https://raw.githubusercontent.com/wabbajack-tools/mod-lists/master/reports/Geborgen/nordic-souls/status.json",
	"https://raw.githubusercontent.com/wabbajack-tools/mod-lists/master/reports/HoS/HoS/status.json",
	"https://raw.githubusercontent.com/wabbajack-tools/mod-lists/master/reports/wj-featured/living_skyrim/status.json",
}

func TestCompareParseFromApiAndFile(t *testing.T) {
	apiArchiveSumSize := parser.ParseJsonFromApiURL("https://raw.githubusercontent.com/wabbajack-tools/mod-lists/master/reports/Geborgen/nordic-souls/status.json")
	localArchiveSumSize := parser.ParseJsonFromFile("archiveData.json")
	assert.Equal(t, apiArchiveSumSize, localArchiveSumSize)
}

func TestParseMultipleApi(t *testing.T) {
	apiArchiveSumSize := parser.ParseMultipleApi(ApiUrls)
	assert.NotEmpty(t, apiArchiveSumSize)
}

func TestParseMultipleApiConurrent(t *testing.T) {
	apiArchiveSumSize := parser.ParseMultipleApiConcurrent(ApiUrls)
	assert.NotEmpty(t, apiArchiveSumSize)
}

func TestGetModsCountAcrossModpacks(t *testing.T) {
	// we cannot reasonably set a specific quantity a certain mod should appear cause  modlists change
	modsCountMap := parser.GetModsCountAcrossModpacks(ApiUrls)
	assert.Contains(t, modsCountMap, "Book Covers Skyrim")
}

func TestCreateUrlLinksForApiCall(t *testing.T) {
	// this creates url links but it doesn't give us the actual game for which modpack it is
	urlLinks := parser.CreateUrlLinksForApiCall()

	assert.NotNil(t, urlLinks)
}

func TestStoreModpacksBasedOnGame(t *testing.T) {
	urlLinks := parser.CreateUrlLinksForApiCall()

	gameModpackMap := make(map[string]string, len(urlLinks))
	gameModpackMap = parser.CreateModPackMap(urlLinks)

	assert.NotNil(t, urlLinks)

}

// func BenchmarkParse(b *testing.B) {
// 	for b.Loop() {
// 		_ = parser.ParseJSON()
// 	}
// }
