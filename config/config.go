package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	AppEnv          string
	LastFmUrl       string
	LastFmApiKey    string
	SpotifyUrl      string
	SpotifyApiKey   string
	AwsRegion       string
	AwsDynamoId     string
	AwsDynamoSecret string
}

func LoadConfg() (*Config, error) {
	err := godotenv.Load()
	if err != nil {
		log.Printf("Warning: No .env file was found, using default environment variables.")
	}

	config := &Config{
		AppEnv:          os.Getenv("APP_ENV"),
		LastFmUrl:       os.Getenv("LASTFM_URL"),
		LastFmApiKey:    os.Getenv("LASTFM_API_KEY"),
		SpotifyUrl:      os.Getenv("SPOTIFY_URL"),
		SpotifyApiKey:   os.Getenv("SPOTIFY_API_KEY"),
		AwsRegion:       os.Getenv("AWS_REGION"),
		AwsDynamoId:     os.Getenv("AWS_DYNAMO_ID"),
		AwsDynamoSecret: os.Getenv("AWS_DYNAMO_SECRET"),
	}

	requiredEnvVars := []struct {
		name  string
		value string
	}{
		{"APP_ENV", config.AppEnv},
		{"LASTFM_URL", config.LastFmUrl},
		{"LASTFM_API_KEY", config.LastFmApiKey},
		{"SPOTIFY_URL", config.SpotifyUrl},
		{"SPOTIFY_API_KEY", config.SpotifyApiKey},
		{"AWS_REGION", config.AwsRegion},
		{"AWS_DYNAMO_ID", config.AwsDynamoId},
		{"AWS_DYNAMO_SECRET", config.AwsDynamoSecret},
	}

	// Verifica se alguma variável obrigatória está vazia
	for _, envVar := range requiredEnvVars {
		if envVar.value == "" {
			return nil, fmt.Errorf("Incomplete configuration variables: check your .env file, missing: %s", envVar.name)
		}
	}
	return config, nil
}
