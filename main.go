package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/pocockn/spotify-api/spotify"
)

func main() {
	lambda.Start(spotify.Handle())
}
