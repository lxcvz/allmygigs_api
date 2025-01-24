package port

type ArtistsRepository interface {
	GetTopArtists(user string, period string, limit string) ([]string, error)
}
