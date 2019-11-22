package main

import (
	"fmt"
	//"context"
	//"encoding/json"

	//"github.com/aws/aws-lambda-go/events"
	//"github.com/aws/aws-lambda-go/lambda"
)

type Headers struct {}
//type Response events.APIGatewayProxyResponse
type Request struct {
	Src string `json:src`
}

type Name string

type Var struct {
	name string
}

type Lam struct {
	name string
	expr Expr
}

type App struct {
	f Expr
	arg Expr
}

type Expr interface {
	reduce() Expr
}

func (x Var) reduce() Expr {
	return x
}

func (x Lam) reduce() Expr {
	return x
}

func (x App) reduce() Expr {
	f := x.f.reduce()
	arg := x.arg.reduce()
	switch f.(type) {
	case Lam:
		return f
	default:
		return App{f, arg}
	}
}

/*
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

*/
func main(){
	x := App{Lam{"x", Var{"x"}}, Var{"y"}}
	fmt.Println(x.reduce())
}
