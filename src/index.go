package main

import (
	"context"
	"fmt"
)

type Response struct {
	StatusCode int         `json:"statusCode"`
	Body       interface{} `json:"body"`
}

func Handler(ctx context.Context) (*Response, error) {
	return &Response{
		StatusCode: 200,
		Body:       fmt.Sprintf("TestQ %d", GetRand()),
	}, nil
}

func main() {
	fmt.Println(Handler(context.TODO()))
}
