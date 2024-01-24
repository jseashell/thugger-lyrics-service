package main

import (
	"bytes"
	"context"
	"encoding/json"
	"log/slog"
	"strconv"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go/aws"
)

// Response is of type APIGatewayProxyResponse since we're leveraging the
// AWS Lambda Proxy Request functionality (default behavior)
//
// https://serverless.com/framework/docs/providers/aws/events/apigateway/#lambda-proxy-integration
type Response events.APIGatewayProxyResponse

type DynamoSong struct {
	ArtistNames              string `dynamodbav:"artist_names"`
	FullTitle                string `dynamodbav:"full_title"`
	HeaderImageThumbnailURL  string `dynamodbav:"header_image_thumbnail_url"`
	HeaderImageURL           string `dynamodbav:"header_image_url"`
	ID                       int    `dynamodbav:"id"`
	Path                     string `dynamodbav:"path"`
	ReleaseDateForDisplay    string `dynamodbav:"release_date_for_display"`
	SongArtImageThumbnailURL string `dynamodbav:"song_art_image_thumbnail_url"`
	SongArtImageURL          string `dynamodbav:"song_art_image_url"`
	Title                    string `dynamodbav:"title"`
	URL                      string `dynamodbav:"url"`
}

type Song struct {
	ArtistNames              string `json:"artist_names"`
	FullTitle                string `json:"full_title"`
	HeaderImageThumbnailURL  string `json:"header_image_thumbnail_url"`
	HeaderImageURL           string `json:"header_image_url"`
	ID                       int    `json:"id"`
	Path                     string `json:"path"`
	ReleaseDateForDisplay    string `json:"release_date_for_display"`
	SongArtImageThumbnailURL string `json:"song_art_image_thumbnail_url"`
	SongArtImageURL          string `json:"song_art_image_url"`
	Title                    string `json:"title"`
	URL                      string `json:"url"`
}

// Handler is our lambda handler invoked by the `lambda.Start` function call
func Handler(ctx context.Context, event events.APIGatewayProxyRequest) (Response, error) {
	var buf bytes.Buffer

	songId, _ := strconv.Atoi(event.QueryStringParameters["song_id"])

	lyric := song(songId)
	body, err := json.Marshal(lyric)
	if err != nil {
		return Response{StatusCode: 404}, err
	}
	json.HTMLEscape(&buf, body)

	resp := Response{
		StatusCode:      200,
		IsBase64Encoded: false,
		Body:            buf.String(),
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	}

	return resp, nil
}

func main() {
	lambda.Start(Handler)
}

func song(songId int) Song {
	dbClient := newDbClient()

	av, _ := attributevalue.MarshalMap(map[string]interface{}{
		"ID": songId,
	})

	res, err := dbClient.GetItem(context.TODO(), &dynamodb.GetItemInput{
		Key:       av,
		TableName: aws.String("thugger-songs"),
	})

	if err != nil {
		panic(err)
	}

	dsong := &DynamoSong{}
	attributevalue.UnmarshalMap(res.Item, dsong)

	song := Song{
		ArtistNames:              dsong.ArtistNames,
		FullTitle:                dsong.FullTitle,
		HeaderImageThumbnailURL:  dsong.HeaderImageThumbnailURL,
		HeaderImageURL:           dsong.HeaderImageURL,
		ID:                       dsong.ID,
		Path:                     dsong.Path,
		ReleaseDateForDisplay:    dsong.ReleaseDateForDisplay,
		SongArtImageThumbnailURL: dsong.SongArtImageThumbnailURL,
		SongArtImageURL:          dsong.SongArtImageURL,
		Title:                    dsong.Title,
		URL:                      dsong.URL,
	}

	return song
}

func newDbClient() *dynamodb.Client {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		slog.Error("Unable to load AWS SDK config.")
		panic(err)
	}
	return dynamodb.NewFromConfig(cfg)
}
