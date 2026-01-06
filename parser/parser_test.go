package parser_test

import (
	"fmt"
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

func TestCountEachModAppearanceInMultipleModpacks(t *testing.T) {
	// each of these modpacks has the Book Covers Skyrim
	// so if the total count for the mod Book Covers Skyrim is 3 then this this is working correctly
	modsCountMap := parser.GetModsCountAcrossModpacks(ApiUrls)
	if assert.Contains(t, modsCountMap, "Book Covers Skyrim") {
		assert.Equal(t, 3, modsCountMap["Book Covers Skyrim"], "Value for key 'Book Covers Skyrim' should be 3")
	}

	all := parser.GetTopPopularMods(ApiUrls, 10)
	fmt.Printf("all: %v\n", all)
}

// func BenchmarkParse(b *testing.B) {
// 	for b.Loop() {
// 		_ = parser.ParseJSON()
// 	}
// }
