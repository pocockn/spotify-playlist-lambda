package spotify

import (
	"encoding/json"
	"github.com/aws/aws-sdk-go/aws"
	dynamoDBLib "github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	dynamoDBWrapper "github.com/pocockn/awswrappers/dynamodb"
	spotifyClient "github.com/pocockn/spotify-api/client"
	"github.com/pocockn/spotify-api/models"
	"github.com/zmb3/spotify"
	"io"
)

type (
	// Client holds the dynamo client for interacting with APIs.
	Client struct {
		DynamoDB *dynamoDBWrapper.Client
	}
)

// Handle is the function run within Lambda. It fetches songs from the playlist and adds them into a Dynamo table,
// a song will only be added once, we can use the Spotify track ID too not add duplicates.
func (c Client) Handle() error {
	client := spotifyClient.NewClient()
	err := client.Auth()
	if err != nil {
		return errors.Wrap(err, "creating spotify client")
	}

	request, err := client.NewRequest("GET", "playlists/3qSco6IsRGt20VPM9XTGeh", nil)
	if err != nil {
		return errors.Wrap(err, "problem performing creating API request")
	}

	response, err := client.Do(request, nil)
	if err != nil {
		return errors.Wrap(err, "problem performing request against the Spotify API")
	}

	playlist, err := c.newPlaylistFromReader(response.Body)

	songs, err := c.fetchSongs()
	if err != nil {
		return err
	}

	err = c.addSongsToDynamo(playlist, songs)
	if err != nil {
		return err
	}

	return nil
}

func (c Client) addSongsToDynamo(playlist spotify.FullPlaylist, songs []models.Song) error {
	if len(songs) == len(playlist.Tracks.Tracks) {
		return nil
	}

	for _, track := range playlist.Tracks.Tracks {
		for _, song := range songs {
			if track.Track.ID == song.SpotifyID {
				continue
			}

			_, err := c.DynamoDB.PutItem(
				models.Song{
					SpotifyID: track.Track.ID,
					Name:      track.Track.Name,
					ID:        uuid.New().String(),
				},
			)
			if err != nil {
				return errors.Wrap(err, "problem putting item into dynamo")
			}
		}
	}

	return nil
}

func (c Client) fetchSongs() ([]models.Song, error) {
	scanParams := dynamoDBLib.ScanInput{
		TableName:     aws.String("brief_selected_watcher"),
		TotalSegments: aws.Int64(0),
	}

	var songs []models.Song
	err := c.DynamoDB.Scan(
		scanParams,
		&songs,
	)
	if err != nil {
		return nil, errors.Wrap(err, "problem fetching items from dyanmo")
	}

	return songs, nil
}

func (c Client) newPlaylistFromReader(response io.Reader) (spotify.FullPlaylist, error) {
	playlist := spotify.FullPlaylist{}

	err := json.NewDecoder(response).Decode(&playlist)
	if err != nil {
		return spotify.FullPlaylist{}, errors.Wrap(err, "problem decoding response into playlist")
	}

	return playlist, nil
}
