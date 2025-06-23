package multi_resolution_transcoding

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type VideoMetadata struct {
	ID       int    `json:"id"`
	Title    string `json:"title"`
	Filename string `json:"filename"`
	Path     string `json:"path"`
}

func FetchVideosMetadata(jwtToken string) ([]VideoMetadata, error) {
	url := "http://localhost:4003/api/upload/videos"

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+jwtToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("unexpected response: %s", string(body))
	}

	var videos []VideoMetadata
	err = json.NewDecoder(resp.Body).Decode(&videos)
	if err != nil {
		return nil, err
	}

	return videos, nil
}
