package repository

import (
	"allmygigs/internal/application/port"
	"allmygigs/internal/infrastructure/adapters"
)

// LastFMRepository estrutura responsável pela interação com a API do LastFM
type LastFMRepository struct {
	adapter *adapters.LastfmAdapter
}

// NewLastFMRepository cria uma nova instância do repositório
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
