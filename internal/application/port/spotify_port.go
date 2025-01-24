package port

type SpotifyRepository interface {
	GetArtists(ids []string) ([]map[string]interface{}, error)
}
