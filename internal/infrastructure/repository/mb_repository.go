package repository

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"regexp"
	"strings"
	"sync"
	"time"
)

type MusicBrainzRepository struct {
	APIBaseURL string
}

func NewMusicBrainzRepository(apiBaseURL string) *MusicBrainzRepository {
	return &MusicBrainzRepository{
		APIBaseURL: apiBaseURL,
	}
}

func (r *MusicBrainzRepository) GetAllIds(mbids []string) ([]string, error) {
	var wg sync.WaitGroup
	spotifyIDs := make(chan string, len(mbids))

	client := &http.Client{
		Timeout: 13 * time.Second,
	}

	for _, mbid := range mbids {
		wg.Add(1)
		go func(mbid string) {
			defer wg.Done()
			spotifyID, err := r.GetOneId(client, mbid)
			if err != nil {
				log.Printf("Erro ao buscar Spotify ID para MBID %s: %v", mbid, err)
				return
			}
			spotifyIDs <- spotifyID
		}(mbid)
	}

	wg.Wait()
	close(spotifyIDs)

	var result []string
	for id := range spotifyIDs {
		result = append(result, id)
	}

	if len(result) == 0 {
		return nil, fmt.Errorf("no Spotify ID found")
	}

	return result, nil
}

func (r *MusicBrainzRepository) GetOneId(client *http.Client, mbid string) (string, error) {
	url := fmt.Sprintf("%s/artist/%s?inc=url-rels&fmt=json", r.APIBaseURL, mbid)

	resp, err := r.makeRequest(client, url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("error in the MusicBrainz API response: status: %d - %s", resp.StatusCode, http.StatusText(resp.StatusCode))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading the MusicBrainz response: %v", err)
	}

	return r.extractSpotifyID(body)
}

func (r *MusicBrainzRepository) makeRequest(client *http.Client, url string) (*http.Response, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating the request: %v", err)
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "Allmygigs.live/1.0.0 (lucas-mateus.dc@hotmail.com)")

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making the request to MusicBrainz: %v", err)
	}

	return resp, nil
}

func (r *MusicBrainzRepository) extractSpotifyID(body []byte) (string, error) {
	var result map[string]interface{}
	err := json.Unmarshal(body, &result)
	if err != nil {
		return "", fmt.Errorf("error decoding JSON response from MusicBrainz: %v", err)
	}

	relations, ok := result["relations"].([]interface{})
	if !ok {
		return "", fmt.Errorf("unexpected format in the MusicBrainz response: relations not found or not a list")
	}

	for _, relation := range relations {
		relationMap, ok := relation.(map[string]interface{})
		if !ok {
			continue
		}

		urlMap, ok := relationMap["url"].(map[string]interface{})
		if !ok {
			continue
		}

		resource, ok := urlMap["resource"].(string)
		if ok && strings.Contains(resource, "https://open.spotify.com/artist") {
			re := regexp.MustCompile(`/artist/([A-Za-z0-9]+)`)
			matches := re.FindStringSubmatch(resource)
			if len(matches) > 1 {
				return matches[1], nil
			}
		}
	}

	return "", fmt.Errorf("spotify ID not found - %s", err)
}
