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
