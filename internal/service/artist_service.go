package service

import (
	"allmygigs/internal/application/port"
	"allmygigs/internal/domain/entity"
	"allmygigs/internal/infrastructure/repository"
	"fmt"
)

type TopArtistsService struct {
	lastfmRepo      port.LastfmRepository
	musicBrainzRepo *repository.MusicBrainzRepository
	spotifyRepo     port.SpotifyRepository
	artistRepo      *repository.ArtistRepository
}

func NewArtistService(
	lastfmRepo port.LastfmRepository,
	spotifyRepo port.SpotifyRepository,
	musicBrainzRepo *repository.MusicBrainzRepository,
	artistRepo *repository.ArtistRepository,
) *TopArtistsService {
	return &TopArtistsService{
		lastfmRepo:      lastfmRepo,
		spotifyRepo:     spotifyRepo,
		musicBrainzRepo: musicBrainzRepo,
		artistRepo:      artistRepo,
	}
}

func (s *TopArtistsService) GetTopArtists(user string, period string, limit string) ([]*entity.Artist, error) {
	mbids, err := s.lastfmRepo.GetTopArtists(user, period, limit)
	if err != nil {
		return nil, fmt.Errorf("error getting top artists: %v", err)
	}

	spotifyIDs, err := s.musicBrainzRepo.GetAllIds(mbids)
	if err != nil {
		return nil, fmt.Errorf("error getting Spotify IDs from MB: %v", err)
	}

	dbArtists, missingIds, err := s.artistRepo.GetArtistsBySpotifyIDsBatch(spotifyIDs)
	if err != nil {
		return nil, fmt.Errorf("error getting artists from database: %v", err)
	}

	if len(missingIds) == 0 {
		return dbArtists, nil
	}

	spotifyMissing, err := s.spotifyRepo.GetArtists(missingIds)
	if err != nil {
		return nil, fmt.Errorf("error getting artists from Spotify: %v", err)
	}

	err = s.artistRepo.SaveArtistsBatch(spotifyMissing)
	if err != nil {
		return nil, fmt.Errorf("error saving artists: %v", err)
	}

	artistsToSave := make([]*entity.Artist, len(spotifyMissing))
	for i, item := range spotifyMissing {
		artistsToSave[i] = &entity.Artist{
			ArtistID:    item["artist_id"].(string),
			ArtistName:  item["artist_name"].(string),
			ArtistImage: item["artist_image"].(string),
		}
	}

	artists := append(dbArtists, artistsToSave...)

	return artists, nil
}
