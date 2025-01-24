package service

import (
	"allmygigs/internal/infrastructure/repository"
)

// SpotifyService estrutura para gerenciar as lógicas de negócio do Spotify
type SpotifyService struct {
	SpotifyRepo *repository.SpotifyRepository
}

// NewSpotifyService cria uma nova instância do serviço de Spotify
func NewSpotifyService(spotifyRepo *repository.SpotifyRepository) *SpotifyService {
	return &SpotifyService{
		SpotifyRepo: spotifyRepo,
	}
}
