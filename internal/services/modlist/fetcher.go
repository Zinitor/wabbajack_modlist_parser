package modlist

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type HTTPError struct {
	Code int
}

func (e *HTTPError) Error() string {
	return fmt.Sprintf("HTTP request failed with status: %d", e.Code)
}

func fetchAndParse[T any](ctx context.Context, client *http.Client, uri string) (T, error) {
	var data T
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return data, err
	}

	response, err := client.Do(request)
	if err != nil {
		return data, err
	}

	if response.StatusCode != http.StatusOK {
		response.Body.Close()
		return data, &HTTPError{Code: response.StatusCode}
	}
	defer response.Body.Close()

	err = json.NewDecoder(response.Body).Decode(&data)
	return data, err
}
