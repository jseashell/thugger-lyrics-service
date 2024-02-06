package main

import (
	"bytes"
	"context"
	"encoding/json"
	"log/slog"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/google/uuid"
)

// Response is of type APIGatewayProxyResponse since we're leveraging the
// AWS Lambda Proxy Request functionality (default behavior)
//
// https://serverless.com/framework/docs/providers/aws/events/apigateway/#lambda-proxy-integration
type Response events.APIGatewayProxyResponse

type Song struct {
	ArtistNames              string   `json:"artist_names"`
	FullTitle                string   `json:"full_title"`
	HeaderImageThumbnailURL  string   `json:"header_image_thumbnail_url"`
	HeaderImageURL           string   `json:"header_image_url"`
	SongID                   int      `json:"id"`
	ID                       string   `json:"uuid"`
	Path                     string   `json:"path"`
	ReleaseDateForDisplay    string   `json:"release_date_for_display"`
	SongArtImageThumbnailURL string   `json:"song_art_image_thumbnail_url"`
	SongArtImageURL          string   `json:"song_art_image_url"`
	Title                    string   `json:"title"`
	URL                      string   `json:"url"`
	Lyrics                   []string `json:"lyrics"`
}

type RandomSong struct {
	ArtistNames              string   `dynamodbav:"artist_names"`
	FullTitle                string   `dynamodbav:"full_title"`
	HeaderImageThumbnailURL  string   `dynamodbav:"header_image_thumbnail_url"`
	HeaderImageURL           string   `dynamodbav:"header_image_url"`
	SongID                   int      `dynamodbav:"id"`
	ID                       string   `dynamodbav:"uuid"`
	Path                     string   `dynamodbav:"path"`
	ReleaseDateForDisplay    string   `dynamodbav:"release_date_for_display"`
	SongArtImageThumbnailURL string   `dynamodbav:"song_art_image_thumbnail_url"`
	SongArtImageURL          string   `dynamodbav:"song_art_image_url"`
	Title                    string   `dynamodbav:"title"`
	URL                      string   `dynamodbav:"url"`
	Lyrics                   []string `dynamodbav:"lyrics"`
}

// Handler is our lambda handler invoked by the `lambda.Start` function call
func Handler(ctx context.Context) (Response, error) {
	var buf bytes.Buffer

	lyric := randomSong()
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

func randomSong() Song {
	dbClient := newDbClient()

	av, _ := attributevalue.MarshalMap(map[string]interface{}{
		"ID": uuid.NewString(),
	})

	limit := int32(1)
	limitPtr := &limit

	res, err := dbClient.Scan(context.TODO(), &dynamodb.ScanInput{
		ExclusiveStartKey: av,
		TableName:         aws.String("thugger-songs"),
		Limit:             limitPtr,
	})

	if err != nil {
		panic(err)
	}

	randomSong := &RandomSong{}
	attributevalue.UnmarshalMap(res.Items[0], randomSong)

	song := Song{
		ArtistNames:              randomSong.ArtistNames,
		FullTitle:                randomSong.FullTitle,
		HeaderImageThumbnailURL:  randomSong.HeaderImageThumbnailURL,
		HeaderImageURL:           randomSong.HeaderImageURL,
		SongID:                   randomSong.SongID,
		ID:                       randomSong.ID,
		Path:                     randomSong.Path,
		ReleaseDateForDisplay:    randomSong.ReleaseDateForDisplay,
		SongArtImageThumbnailURL: randomSong.SongArtImageThumbnailURL,
		SongArtImageURL:          randomSong.SongArtImageURL,
		Title:                    randomSong.Title,
		URL:                      randomSong.URL,
		Lyrics:                   randomSong.Lyrics,
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
