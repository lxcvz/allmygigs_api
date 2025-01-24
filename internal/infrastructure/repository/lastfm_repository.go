package repository

import (
	"allmygigs/internal/application/port"
	"allmygigs/internal/infrastructure/adapters"
)

type LastFMRepository struct {
	adapter *adapters.LastfmAdapter
}

func NewLastFMRepository(adapter *adapters.LastfmAdapter) port.LastfmRepository {
	return &LastFMRepository{
		adapter: adapter,
	}
}

func (r *LastFMRepository) GetTopArtists(user string, period string, limit string) ([]string, error) {
	artists, err := r.adapter.GetTopArtists(user, period, limit)
	if err != nil {
		return nil, err
	}

	return artists, nil
}
