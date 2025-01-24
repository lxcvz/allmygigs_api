package repository

import (
	"allmygigs/internal/application/port"
	"allmygigs/internal/infrastructure/adapters"
	"fmt"
)

type SpotifyRepository struct {
	AccessToken string
	adapter     *adapters.SpotifyAdapter
}

func NewSpotifyRepository(adapter *adapters.SpotifyAdapter) port.SpotifyRepository {
	return &SpotifyRepository{
		adapter: adapter,
	}
}

func (r *SpotifyRepository) GetArtists(ids []string) ([]map[string]interface{}, error) {
	result, err := r.adapter.GetArtistsById(ids)
	if err != nil {
		return nil, fmt.Errorf("erro ao obter artistas do Spotify: %v", err)
	}

	return result, nil
}
