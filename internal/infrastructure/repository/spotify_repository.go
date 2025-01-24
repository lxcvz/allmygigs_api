package repository

import (
	"allmygigs/internal/application/port"
	"allmygigs/internal/infrastructure/adapters"
	"fmt"
)

// SpotifyRepository estrutura responsável pela interação com a API do Spotify
type SpotifyRepository struct {
	AccessToken string
	adapter     *adapters.SpotifyAdapter
}

// NewSpotifyRepository cria uma nova instância do repositório do Spotify
func NewSpotifyRepository(adapter *adapters.SpotifyAdapter) port.SpotifyRepository {
	return &SpotifyRepository{
		adapter: adapter,
	}
}

// GetSeveralArtists faz uma requisição para obter vários artistas pela lista de IDs
func (r *SpotifyRepository) GetArtists(ids []string) ([]map[string]interface{}, error) {
	result, err := r.adapter.GetArtistsById(ids)
	if err != nil {
		return nil, fmt.Errorf("erro ao obter artistas do Spotify: %v", err)
	}

	return result, nil
}
