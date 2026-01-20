package modlist

type Summary struct {
	ModlistName  string `json:"Name"`
	ArchivesLink string `json:"link"`
}

type Repository struct {
	Name string
	Link string
}

type GameModlist struct {
	GameName string
	Modlists []string
}

type ModlistData struct {
	Title string `json:"title"`
	Game  string `json:"game"`
}

type GamePopularity struct {
	Name       string
	Popularity int
}
