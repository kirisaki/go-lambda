package main

import (
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type Headers struct {}
type Response events.APIGatewayProxyResponse

func Handler() (Response, error) {
	return Response{
		IsBase64Encoded: false,
		StatusCode: 200,
		Headers: map[string]string{
			"Content-Type":           "application/json",
		},
		Body: "Go Serverless v1.0! Your function executed successfully!",
	}, nil
}

func main() {
	lambda.Start(Handler)
}
