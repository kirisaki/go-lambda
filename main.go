package main

import (
	"context"
	"encoding/json"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type Headers struct {}
type Response events.APIGatewayProxyResponse
type Request struct {
	Src string `json:src`
}

func Handler(ctx context.Context, gwreq events.APIGatewayProxyRequest) (Response, error) {
	body := ([]byte)(gwreq.Body)
	req := new(Request)
	if err := json.Unmarshal(body, req); err != nil {
		panic("")
	}
	return Response{
		IsBase64Encoded: false,
		StatusCode: 200,
		Headers: map[string]string{
			"Content-Type":           "plain/text",
		},
		Body: req.Src,
	}, nil
}

func main() {
	lambda.Start(Handler)
}
