package main

import (
	"context"
	"fmt"
	"os"
)

type Response struct {
	StatusCode int         `json:"statusCode"`
	Body       interface{} `json:"body"`
}

func Handler(ctx context.Context) (*Response, error) {
	return &Response{
		StatusCode: 200,
		Body:       fmt.Sprintf("TestQ %d", 123),
	}, nil
}

func main() {
	apiKey := os.Getenv("DEEPSEEK_API_KEY")
	client := NewDeepSeekClient(apiKey)
	fmt.Println(client.Query("Привет, дружище, я тут пытаюсь с тобой коннектнуться."))
}
