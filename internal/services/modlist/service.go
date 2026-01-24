package modlist

import (
	"context"
	"fmt"
	"net/http"
	"sort"
	"strings"
	"sync"
	"wabbajackModlistParser/internal/adapters/rest"
	"wabbajackModlistParser/pkg/logger"

	"golang.org/x/sync/errgroup"
)

type Service struct {
	l      logger.Interface
	client *http.Client
}

func NewModlistService(logger logger.Interface, client *http.Client) Service {
	return Service{l: logger, client: client}
}

func (s *Service) GetModlistSummary(ctx context.Context) ([]Summary, error) {
	var modlists []Summary
	uri := "https://raw.githubusercontent.com/wabbajack-tools/mod-lists/master/reports/modListSummary.json"

	modlists, err := rest.Get[[]Summary](ctx, s.client, uri)
	if err != nil {
		return modlists, err
	}

	return modlists, nil
}

func (s *Service) GetUserRepos(ctx context.Context) ([]Repository, error) {
	var repositories []Repository

	uri := "https://raw.githubusercontent.com/wabbajack-tools/mod-lists/master/repositories.json"

	repoMaps, err := rest.Get[map[string]string](ctx, s.client, uri)
	if err != nil {
		return repositories, err
	}
	repositories = make([]Repository, 0, len(repoMaps))

	for name, link := range repoMaps {
		repositories = append(repositories,
			Repository{
				Name: name,
				Link: link,
			},
		)
	}

	return repositories, nil
}

func (s *Service) GetAllGamesFromModlists(ctx context.Context) ([]string, error) {
	gameSet := make(map[string]struct{})
	var mu sync.Mutex //for map

	repos, err := s.GetUserRepos(ctx)
	if err != nil {
		return nil, err
	}

	g, ctx := errgroup.WithContext(ctx)
	g.SetLimit(6)

	for _, repo := range repos {
		g.Go(func() error {
			modlistData, fErr := rest.Get[[]Data](ctx, s.client, repo.Link)
			if fErr != nil {
				return fmt.Errorf("failed to fetch modlists from %s:%w", repo.Link, fErr)
			}
			mu.Lock()
			for _, item := range modlistData {
				gameSet[item.Game] = struct{}{}
			}
			mu.Unlock()
			return nil
		})
	}

	if err = g.Wait(); err != nil {
		return nil, err
	}

	result := make([]string, 0, len(gameSet))
	for k := range gameSet {
		result = append(result, k)
	}
	// could probably optimize this bit
	sort.Strings(result)
	return result, nil
}

func (s *Service) GetTopPopularGames(ctx context.Context, gamesCount int, sortOrder string) ([]GamePopularity, error) {
	gameSet := make(map[string]int)
	var mu sync.Mutex //for map

	repos, err := s.GetUserRepos(ctx)
	if err != nil {
		return nil, err
	}
	g, ctx := errgroup.WithContext(ctx)
	g.SetLimit(6)

	for _, repo := range repos {
		g.Go(func() error {
			modlistData, err := rest.Get[[]Data](ctx, s.client, repo.Link)
			if err != nil {
				s.l.Error("failed to fetch modlists from %s:%w", repo.Link, err)
				return nil
			}
			mu.Lock()
			for _, item := range modlistData {
				gameSet[item.Game]++
			}
			mu.Unlock()
			return nil
		})
	}

	if err = g.Wait(); err != nil {
		return nil, err
	}

	resultSlice := make([]GamePopularity, 0, gamesCount)
	for key, value := range gameSet {
		item := GamePopularity{
			Name:       key,
			Popularity: value,
		}
		resultSlice = append(resultSlice, item)
	}

	switch strings.ToLower(sortOrder) {
	case "asc", "ascending":
		sort.Slice(resultSlice, func(i, j int) bool {
			if resultSlice[i].Popularity == resultSlice[j].Popularity {
				return resultSlice[i].Name < resultSlice[j].Name
			}
			return resultSlice[i].Popularity < resultSlice[j].Popularity
		})
	case "desc", "descending", "":
		fallthrough // Default to descending
	default:
		sort.Slice(resultSlice, func(i, j int) bool {
			if resultSlice[i].Popularity == resultSlice[j].Popularity {
				return resultSlice[i].Name < resultSlice[j].Name
			}
			return resultSlice[i].Popularity > resultSlice[j].Popularity
		})
	}

	if gamesCount < len(resultSlice) {
		resultSlice = resultSlice[:gamesCount]
	}

	return resultSlice, nil
}

func CreateURLLinkForAPICall(archiveListPostfix string) string {
	urlPrefix := "https://raw.githubusercontent.com/wabbajack-tools/mod-lists/master/"

	return urlPrefix + archiveListPostfix
}

// func GetModlistForGame(ctx context.Context, gameName string) (GameModlist, error) {

// }

// func GetModlistsForGames(ctx context.Context) ([]GameModlist, error) {

// }
