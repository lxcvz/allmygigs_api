package routes

import (
	"allmygigs/config"
	"allmygigs/internal/application/usecase"
	"allmygigs/internal/db"
	"allmygigs/internal/infrastructure/adapters"
	"allmygigs/internal/infrastructure/handler"
	"allmygigs/internal/infrastructure/repository"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine, cfg *config.Config, dynamoClient *db.DynamoDBClient) {

	lastfmAdapter := adapters.NewLastfmAdapter(cfg.LastFmUrl, cfg.LastFmApiKey)
	spotifyAdapter := adapters.NewSpotifyAdapter(cfg.SpotifyUrl, cfg.SpotifyApiKey)

	artistRepo := repository.NewArtistRepository(dynamoClient.Client)
	lastFMRepo := repository.NewLastFMRepository(lastfmAdapter)
	spotifyRepo := repository.NewSpotifyRepository(spotifyAdapter)
	musicBrainzRepo := repository.NewMusicBrainzRepository("https://musicbrainz.org/ws/2")

	artistService := usecase.NewArtistUsecase(lastFMRepo, spotifyRepo, musicBrainzRepo, artistRepo)
	healthUseCase := usecase.NewHealthCheckUsecase()

	artistHandler := handler.NewArtistHandler(artistService)
	healthHandler := handler.NewHealthHandler(healthUseCase)

	v1 := router.Group("/v1")
	{
		v1.GET("/top-artists", artistHandler.GetQuery)
		v1.GET("/health-check", healthHandler.CheckHealth)
	}

}
