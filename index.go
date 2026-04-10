package main

import (
	"context"
	"fmt"
	"log"
)

type Request struct {
	Action string `json:"action"`
}

type Response struct {
	StatusCode int         `json:"statusCode"`
	Body       interface{} `json:"body"`
}

func Handler(ctx context.Context, r *Request) (*Response, error) {
	if r.Action == "default" {
		return defaultHandler(ctx, r, nil)
	}
	return &Response{
		StatusCode: 200,
		Body:       "TestQ: invalid",
	}, nil
}

func defaultHandler(_ context.Context, r *Request, _ *transport) (*Response, error) {
	return &Response{
		StatusCode: 200,
		Body:       fmt.Sprintf("TestQ: [%s]", r.Action),
	}, nil
}

func main() {
	// apiKey := os.Getenv("DEEPSEEK_API_KEY")
	// client := NewDeepSeekClient(apiKey)
	// fmt.Println(client.Query("Привет, дружище, я тут пытаюсь с тобой коннектнуться."))

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	tp, closeFunc, err := InitTransport(ctx)
	if err != nil {
		log.Fatal(err)
	}

	defer closeFunc()
	wds, err := GetWorldDescriptions(ctx, tp)
	fmt.Println(wds)

}
