package adapters

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type LastfmAdapter struct {
	apiKey string
	apiUrl string
}

func NewLastfmAdapter(apiUrl string, apiKey string) *LastfmAdapter {
	return &LastfmAdapter{
		apiKey: apiKey,
		apiUrl: apiUrl,
	}
}

func (a *LastfmAdapter) GetTopArtists(user string, period string, limit string) ([]string, error) {
	url := fmt.Sprintf(
		"%s?method=user.gettopartists&user=%s&api_key=%s&period=%s&limit=%s&format=json",
		a.apiUrl, user, a.apiKey, period, limit,
	)

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("error on lastfm request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error on lastfm response, status: %d - %s", resp.StatusCode, http.StatusText(resp.StatusCode))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error on reading lastfm response: %v", err)
	}

	var result map[string]interface{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, fmt.Errorf("error decoding JSON response: %v", err)
	}

	// Extraindo os artistas de uma maneira segura
	topArtistsData, ok := result["topartists"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("unexpected format in the LastFM response: topartists not found.")
	}

	artists, ok := topArtistsData["artist"].([]interface{})
	if !ok {
		return nil, fmt.Errorf("unexpected format in the LastFM response: artist is not a list.")
	}

	var artistIDs []string
	for _, artist := range artists {
		artistMap, ok := artist.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("unexpected format in the LastFM response: artist is not an object.")
		}
		if artistName, ok := artistMap["mbid"].(string); ok {
			artistIDs = append(artistIDs, artistName)
		}
	}

	return artistIDs, nil
}
