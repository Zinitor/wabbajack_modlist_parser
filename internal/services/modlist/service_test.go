package modlist

import (
	"context"
	"testing"
	"wabbajackModlistParser/pkg/logger"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

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
		assert.Equal(t, repos, expectedRepos)
	})
}
