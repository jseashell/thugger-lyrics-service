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

type Lyric struct {
	ID     string `json:"id"`
	SongID int    `json:"song_id"`
	Value  string `json:"value"`
}

type RandomLyric struct {
	ID     string `dynamodbav:"ID"`
	SongID int    `dynamodbav:"SongID"`
	Value  string `dynamodbav:"Value"`
}

// Handler is our lambda handler invoked by the `lambda.Start` function call
func Handler(ctx context.Context) (Response, error) {
	var buf bytes.Buffer

	lyric := randomLyric()
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

func randomLyric() Lyric {
	dbClient := newDbClient()

	av, _ := attributevalue.MarshalMap(map[string]interface{}{
		"ID": uuid.NewString(),
	})

	limit := int32(1)
	limitPtr := &limit

	res, err := dbClient.Scan(context.TODO(), &dynamodb.ScanInput{
		ExclusiveStartKey: av,
		TableName:         aws.String("thugger-lyrics"),
		Limit:             limitPtr,
	})

	if err != nil {
		panic(err)
	}

	randomLyric := &RandomLyric{}
	attributevalue.UnmarshalMap(res.Items[0], randomLyric)

	lyric := Lyric{
		ID:     randomLyric.ID,
		SongID: randomLyric.SongID,
		Value:  randomLyric.Value,
	}

	return lyric
}

func newDbClient() *dynamodb.Client {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		slog.Error("Unable to load AWS SDK config.")
		panic(err)
	}
	return dynamodb.NewFromConfig(cfg)
}
