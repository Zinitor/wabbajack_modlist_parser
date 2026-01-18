package v1

type ModlistSummaryResponse struct {
	ModlistName  string `json:"Name"`
	ArchivesLink string `json:"link"`
}

type RepositoryResponse struct {
	Name string
	Link string
}
