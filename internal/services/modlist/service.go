package modlist

import (
	"context"
	"net/http"
	"sync"
	"wabbajackModlistParser/pkg/logger"
)

type Service struct {
	l      logger.Interface
	client *http.Client
}

func NewModlistService(logger logger.Interface, client *http.Client) Service {
	return Service{l: logger, client: client}
}

func (s *Service) GetModlistSummary(ctx context.Context) ([]Summary, error) {
	modlists := make([]Summary, 0)
	uri := "https://raw.githubusercontent.com/wabbajack-tools/mod-lists/master/reports/modListSummary.json"

	modlists, err := fetchAndParse[[]Summary](ctx, s.client, uri)
	if err != nil {
		return modlists, err
	}

	return modlists, nil
}

func (s *Service) GetUserRepos(ctx context.Context) ([]Repository, error) {
	repositories := make([]Repository, 0)

	uri := "https://raw.githubusercontent.com/wabbajack-tools/mod-lists/master/repositories.json"

	repoMaps, err := fetchAndParse[map[string]string](ctx, s.client, uri)
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

type ModlistData struct {
	Title string `json:"title"`
	Game  string `json:"game"`
}

type Semaphore chan struct{}

func (s Semaphore) Acquire() {
	s <- struct{}{}
}
func (s Semaphore) Release() {
	<-s
}

func (s *Service) GetAllGamesFromModlists(ctx context.Context) ([]string, error) {
	gameSet := make(map[string]struct{})
	var mu sync.Mutex //for map

	repos, err := s.GetUserRepos(ctx)
	if err != nil {
		return nil, err
	}

	//Should probably move elsewhere
	maxConcurrent := 6
	sem := make(Semaphore, maxConcurrent)
	var wg sync.WaitGroup
	errCh := make(chan error, len(repos))

	for _, repo := range repos {
		wg.Add(1)
		go func(repoLink string) {
			defer wg.Done()
			sem.Acquire()
			defer sem.Release()

			modlistData, err := fetchAndParse[[]ModlistData](ctx, s.client, repo.Link)
			if err != nil {
				errCh <- err
			}
			mu.Lock()
			for _, item := range modlistData {
				gameSet[item.Game] = struct{}{}
			}
			mu.Unlock()
		}(repo.Link)
	}

	wg.Wait()
	close(errCh)

	for err := range errCh {
		if err != nil {
			return nil, err
		}
	}

	result := make([]string, 0, len(gameSet))
	for k, _ := range gameSet {
		result = append(result, k)
	}
	return result, nil

}

func CreateURLLinkForAPICall(archiveListPostfix string) string {
	urlPrefix := "https://raw.githubusercontent.com/wabbajack-tools/mod-lists/master/"

	return urlPrefix + archiveListPostfix
}

// func GetModlistForGame(ctx context.Context, gameName string) (GameModlist, error) {

// }

// func GetModlistsForGames(ctx context.Context) ([]GameModlist, error) {

// }
