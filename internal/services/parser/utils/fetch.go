package utils

import (
	"io"
	"log"
	"net/http"
	"sync"
)

func Fetch(baseURL string) io.ReadCloser {
	response, err := http.Get(baseURL)
	if err != nil {
		log.Fatal(err)
	}

	if response.StatusCode != http.StatusOK {
		response.Body.Close()
		log.Fatalf("API request failed with status: %d", response.StatusCode)
	}

	return response.Body
}

func FetchAndParse[T any](
	baseURL string,
	parseFunc func(io.Reader) []T,
) []T {
	body := Fetch(baseURL)
	defer body.Close()
	parsed := parseFunc(body)
	return parsed
}

func ConcurrentFetchAndParse[T any](
	urls []string,
	parseFunc func(io.Reader) []T,
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
