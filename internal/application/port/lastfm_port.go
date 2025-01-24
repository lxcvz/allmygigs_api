package port

type LastfmRepository interface {
	GetTopArtists(user string, period string, limit string) ([]string, error)
}
