package main

import (
	//"bytes"
	"context"
	//"encoding/json"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type Headers struct {}
type Response events.APIGatewayProxyResponse

func Handler(ctx context.Context, req events.APIGatewayProxyRequest) (Response, error) {
	return Response{
		IsBase64Encoded: false,
		StatusCode: 200,
		Headers: map[string]string{
			"Content-Type":           "plain/text",
		},
		Body: req.Body,
	}, nil
}

func main() {
	lambda.Start(Handler)
}
