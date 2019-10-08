package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	dynamoDBWrapper "github.com/pocockn/awswrappers/dynamodb"
	"github.com/pocockn/spotify-api/spotify"
	"log"
)

func main() {
	dynamoDBClient, err := createDynamoDBClient()
	if err != nil {
		log.Fatal(err)
	}

	spotifyClient := spotify.Client{
		DynamoDB: dynamoDBClient,
	}

	lambda.Start(spotifyClient.Handle)
}

func createDynamoDBClient() (*dynamoDBWrapper.Client, error) {
	return dynamoDBWrapper.NewClient(
		&dynamoDBWrapper.ClientConfig{},
		false,
		nil,
		nil,
	)
}
