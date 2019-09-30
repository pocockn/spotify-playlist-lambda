package spotify

import (
	spotifyClient "github.com/pocockn/spotify-api/client"
)

func Handle() error {
	client := spotifyClient.NewClient()
	err := client.Auth()
	if err != nil {
		panic(err)
	}

	request, _ := client.NewRequest("GET", "playlists/3qSco6IsRGt20VPM9XTGeh", nil)
	_, err = client.Do(request, nil)
	if err != nil {
		panic(err)
	}

	return nil
}
