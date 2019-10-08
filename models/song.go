package models

import (
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/pkg/errors"
	"github.com/zmb3/spotify"
)

type (
	// Song is a song held in Dynamo, this corresponds to the Spotify ID and then some details we want
	// to pull through on the front end. If the song is already in Dynamo we don't bother updating.
	Song struct {
		ID        string     `json:"internal_id"`
		Name      string     `json:"name"`
		SpotifyID spotify.ID `json:"id"`
	}
)

// Marshal converts values stored in a Song, into a DynamoDB
// item to be created within the store.
func (s Song) Marshal() (*dynamodb.PutItemInput, error) {
	attributeValue, err := dynamodbattribute.MarshalMap(s)

	if err != nil {
		return nil, errors.Wrap(
			err,
			"Unable to convert dynamo.Song struct into DynamoDB attribute value",
		)
	}

	return &dynamodb.PutItemInput{
		Item:      attributeValue,
		TableName: spotifySongTableName(),
	}, nil
}

func spotifySongTableName() *string {
	t := "recs_for_the_decks"
	return &t
}
