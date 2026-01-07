package utils

import (
	"io"
	"log"
	"net/http"
	"sync"
)

func Fetch(baseUrl string) []byte {
	response, err := http.Get(baseUrl)
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

	return body
}

func FetchAndParse[T any](
	baseUrl string,
	parseFunc func([]byte) []T,
) []T {
	data := Fetch(baseUrl)
	parsed := parseFunc(data)
	return parsed
}

func ConcurrentFetchAndParse[T any](
	urls []string,
	parseFunc func([]byte) []T,
) <-chan []T {
	var wg sync.WaitGroup
	parsedChan := make(chan []T, len(urls))

	for _, url := range urls {
		u := url
		wg.Go(func() {
			parsedChan <- FetchAndParse(u, parseFunc)
		})
	}
	go func() {
		wg.Wait()
		close(parsedChan)
	}()

	return parsedChan
}
