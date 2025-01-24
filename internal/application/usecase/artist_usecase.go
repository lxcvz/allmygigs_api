package usecase

import (
	"allmygigs/internal/application/port"
	"allmygigs/internal/domain/entity"
	"allmygigs/internal/infrastructure/repository"
	"fmt"
)

type ArtistUsecase struct {
	lastfmRepo      port.LastfmRepository
	spotifyRepo     port.SpotifyRepository
	musicBrainzRepo *repository.MusicBrainzRepository
	artistRepo      *repository.ArtistRepository
}

func NewArtistUsecase(
	lastfmRepo port.LastfmRepository,
	spotifyRepo port.SpotifyRepository,
	musicBrainzRepo *repository.MusicBrainzRepository,
	artistRepo *repository.ArtistRepository,
) *ArtistUsecase {
	return &ArtistUsecase{
		lastfmRepo:      lastfmRepo,
		spotifyRepo:     spotifyRepo,
		musicBrainzRepo: musicBrainzRepo,
		artistRepo:      artistRepo,
	}
}

func (uc *ArtistUsecase) GetUserTopArtists(user string, period string, limit string) ([]*entity.Artist, error) {

	mbids, err := uc.lastfmRepo.GetTopArtists(user, period, limit)
	if err != nil {
		return nil, fmt.Errorf("error to get top artists: %v", err)
	}

	spotifyIDs, err := uc.musicBrainzRepo.GetAllIds(mbids)
	if err != nil {
		return nil, fmt.Errorf("error to get Spotify ids from MB: %v", err)
	}

	dbArtists, missingIds, err := uc.artistRepo.GetArtistsBySpotifyIDsBatch(spotifyIDs)
	if err != nil {
		return nil, fmt.Errorf("error on get artists from database: %v", err)
	}

	if len(missingIds) == 0 {
		return dbArtists, nil
	}

	spotifyMissing, err := uc.spotifyRepo.GetArtists(missingIds)
	if err != nil {
		return nil, fmt.Errorf("error to get artists from Spotify: %v", err)
	}

	err = uc.artistRepo.SaveArtistsBatch(spotifyMissing)
	if err != nil {
		return nil, fmt.Errorf("error to save artists: %v", err)
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
