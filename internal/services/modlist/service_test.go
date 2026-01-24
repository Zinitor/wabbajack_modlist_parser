package modlist

import (
	"context"
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"
	"wabbajackModlistParser/internal/adapters/rest"
	"wabbajackModlistParser/pkg/logger"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func NewMockClient(bodyData string) *http.Client {
	return &http.Client{
		Transport: &MockRoundTripper{
			StatusCode: http.StatusOK,
			Body:       bodyData,
		},
	}
}

// MockRoundTripper to simulate HTTP responses.
type MockRoundTripper struct {
	StatusCode int
	Body       string
	Err        error // To simulate network errors
}

func (m *MockRoundTripper) RoundTrip(_ *http.Request) (*http.Response, error) {
	if m.Err != nil {
		return nil, m.Err
	}
	return &http.Response{
		StatusCode: m.StatusCode,
		Body:       io.NopCloser(strings.NewReader(m.Body)),
		Header:     make(http.Header),
	}, nil
}

func TestGetUserRepos(t *testing.T) {
	testJSON := `
	{
  "wj-featured": "https://raw.githubusercontent.com/wabbajack-tools/mod-lists/master/modlists.json"
  }
	`
	expectedRepos := []Repository{
		{
			Name: "wj-featured",
			Link: "https://raw.githubusercontent.com/wabbajack-tools/mod-lists/master/modlists.json",
		},
	}

	logger := logger.New("info")

	t.Run("success", func(t *testing.T) {
		mockClient := NewMockClient(testJSON)
		s := NewModlistService(logger, mockClient)

		repos, err := s.GetUserRepos(context.TODO())
		require.NoError(t, err)
		assert.Equal(t, expectedRepos, repos)
	})
}

func TestGetModlists(t *testing.T) {
	validJSON := `{
  "Name": "Skyrim Enhanced",
  "link": "https://nexusmods.com/skyrim/mods/archive.zip"
  }
`
	invalidJSON := `{m}`

	testURL := "https://example.com"
	testCtx := context.TODO()

	want := Summary{
		ModlistName:  "Skyrim Enhanced",
		ArchivesLink: "https://nexusmods.com/skyrim/mods/archive.zip",
	}

	t.Run("success", func(t *testing.T) {
		item, err := rest.Get[Summary](testCtx, NewMockClient(validJSON), testURL)
		require.NoError(t, err)
		assert.Equal(t, want, item)
	})

	t.Run("invalid json", func(t *testing.T) {
		item, err := rest.Get[Summary](testCtx, NewMockClient(invalidJSON), testURL)
		require.Error(t, err)
		assert.Equal(t, Summary{}, item)
	})

	t.Run("404 not found", func(t *testing.T) {
		client := &http.Client{
			Transport: &MockRoundTripper{
				StatusCode: http.StatusNotFound,
				Body:       "Not Found",
			},
		}
		item, err := rest.Get[Summary](testCtx, client, testURL)

		var httpErr *rest.HTTPError
		if assert.Error(t, err) && errors.As(err, &httpErr) {
			assert.Equal(t, http.StatusNotFound, httpErr.Code)
		}
		assert.Equal(t, Summary{}, item)
	})

	t.Run("network failure", func(t *testing.T) {
		client := &http.Client{
			Transport: &MockRoundTripper{Err: io.EOF}, // Simulate broken connection
		}
		item, err := rest.Get[Summary](testCtx, client, testURL)
		require.Error(t, err)
		var httpErr *rest.HTTPError
		assert.NotErrorAs(t, err, &httpErr, "Expected network error, not HTTP error")
		assert.Equal(t, Summary{}, item)
	})

	t.Run("500 server error", func(t *testing.T) {
		client := &http.Client{
			Transport: &MockRoundTripper{
				StatusCode: http.StatusInternalServerError,
				Body:       "Internal Server Error",
			},
		}
		item, err := rest.Get[Summary](testCtx, client, testURL)

		var httpErr *rest.HTTPError
		if assert.Error(t, err) && errors.As(err, &httpErr) {
			assert.Equal(t, http.StatusInternalServerError, httpErr.Code)
		}
		assert.Equal(t, Summary{}, item)
	})
}
