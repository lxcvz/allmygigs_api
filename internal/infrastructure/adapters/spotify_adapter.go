package adapters

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type SpotifyAdapter struct {
	apiUrl         string
	apiAccessToken string
}

func NewSpotifyAdapter(apiUrl string, apiAccessToken string) *SpotifyAdapter {
	return &SpotifyAdapter{
		apiUrl:         apiUrl,
		apiAccessToken: apiAccessToken,
	}
}

func (a *SpotifyAdapter) GetArtistsById(ids []string) ([]map[string]interface{}, error) {

	// https://api.spotify.com/v1

	url := fmt.Sprintf("%s/artists?ids=%s", a.apiUrl, strings.Join(ids, ","))

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request for Spotify.: %v", err)
	}

	req.Header.Add("Authorization", "Bearer "+a.apiAccessToken)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error creating request for Spotify.: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error in the Spotify response: status: %d - %s", resp.StatusCode, http.StatusText(resp.StatusCode))
	}

	var result map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return nil, fmt.Errorf("error decoding Spotify JSON response: %v", err)
	}

	// Extrai as informações relevantes de cada artista
	var artistsInfo []map[string]interface{}
	artists := result["artists"].([]interface{})

	for _, artist := range artists {
		artistMap := artist.(map[string]interface{})
		artistData := map[string]interface{}{
			"artist_id":   artistMap["id"],
			"artist_name": artistMap["name"],
		}

		images := artistMap["images"].([]interface{})
		if len(images) > 0 {
			artistData["artist_image"] = images[0].(map[string]interface{})["url"]
		} else {
			artistData["artist_image"] = nil
		}

		artistsInfo = append(artistsInfo, artistData)
	}

	return artistsInfo, nil
}
