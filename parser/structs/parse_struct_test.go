package structs_test

import (
	"bytes"
	"embed"
	"io"
	"path/filepath"
	"testing"
	"wabbajackModlistParser/parser/structs"

	"github.com/stretchr/testify/assert"
)

//go:embed json_examples/*
var f embed.FS

func loadTestData(t testing.TB, filePath string) []byte {
	t.Helper()
	data, err := f.ReadFile(filePath)
	if err != nil {
		t.Fatal(err)
	}
	return data
}

func TestParseToModlistArchiveMap(t *testing.T) {
	//чисто технически мне не обязательно даже парсить в структуру мне достаточно чтобы у меня вернулась инфа о количестве модов
	// следовательно можно просто возвращать map
	data := loadTestData(t, "json_examples/mod_archive_example.json")
	reader := bytes.NewReader(data)
	modlistMap := structs.ParseToModlistArchiveMap(reader)
	assert.NotEmpty(t, modlistMap)
}

func BenchmarkParseToModlistArchiveMap(b *testing.B) {
	data := loadTestData(b, "json_examples/mod_archive_example.json")
	for b.Loop() {
		reader := bytes.NewReader(data)
		structs.ParseToModlistArchiveMap(reader)
	}
}

func benchmarkParseFile[T any](b *testing.B, filename string, foo func(io.Reader) T) {
	data, err := f.ReadFile(filepath.Join("json_examples", filename))
	if err != nil {
		b.Fatal(err)
	}

	for b.Loop() {
		reader := bytes.NewReader(data)
		foo(reader)
		// structs.ParseToModlistArchiveMap(reader)
	}
}

func BenchmarkParseModlist(b *testing.B) {
	testcases := []string{
		"mod_archive_example.json",
		"mod_archive_example_1.json",
	}

	for _, tc := range testcases {
		b.Run(tc, func(b *testing.B) {
			benchmarkParseFile(b, tc, structs.ParseToModlistArchiveMap)
		})
	}
}

func TestParseToModlistInfo(t *testing.T) {
	data := loadTestData(t, "json_examples/modlistinfo_example.json")
	reader := bytes.NewReader(data)
	_ = structs.ParseToModlistInfo(reader)
}

func BenchmarkParseToModlistInfo(b *testing.B) {
	testcases := []string{
		"modlistinfo_example.json",
	}

	for _, tc := range testcases {
		b.Run(tc, func(b *testing.B) {
			benchmarkParseFile(b, tc, structs.ParseToModlistInfo)
		})
	}
}
