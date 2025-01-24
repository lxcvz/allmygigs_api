package service

import (
	"allmygigs/internal/infrastructure/repository"
)

type SpotifyService struct {
	SpotifyRepo *repository.SpotifyRepository
}

func NewSpotifyService(spotifyRepo *repository.SpotifyRepository) *SpotifyService {
	return &SpotifyService{
		SpotifyRepo: spotifyRepo,
	}
}
