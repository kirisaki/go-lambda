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
	Src string `json:"src"`
	Step int `json:"step"`
}



func Handler(ctx context.Context, gwreq events.APIGatewayProxyRequest) (Response, error) {
	body := ([]byte)(gwreq.Body)
	req := new(Request)
	if err := json.Unmarshal(body, req); err != nil {
		return Response{
			IsBase64Encoded: false,
			StatusCode: 400,
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
			Body: "{\"msg\":\"invalid json\"}",
		}, nil
	}
	if req.Step > 1000 {
		return Response{
			IsBase64Encoded: false,
			StatusCode: 400,
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
			Body: "{\"msg\":\"too much steps\"}",
		}, nil
	}
	ast, perr := parse(req.Src)
	if perr != nil {
		return Response{
			IsBase64Encoded: false,
			StatusCode: 400,
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
			Body: "{\"msg\":\"failed parsing: " + perr.Error() + "\"}",
		}, nil
	}
	for i := 0; i < req.Step; i++ {
		ast = ast.reduce()
	}
	return Response{
		IsBase64Encoded: false,
		StatusCode: 200,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		Body: "{\"result\":\"" + prettify(ast) + "\"}",
	}, nil
}

func main() {
	lambda.Start(Handler)
}

