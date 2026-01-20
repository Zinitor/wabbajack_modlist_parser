package v1

type ModlistSummaryResponse struct {
	ModlistName  string `json:"Name"`
	ArchivesLink string `json:"link"`
}

type RepositoryResponse struct {
	Name string
	Link string
}

type GamePopularity struct {
	Name       string `json:"name"`
	Popularity int    `json:"popularity"`
}

type ErrorResponse struct {
	Error   string `json:"error"`
	Code    int    `json:"code"`
	Message string `json:"message,omitempty"`
}
