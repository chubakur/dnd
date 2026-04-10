package main

import (
	"context"
	"encoding/json"
	"fmt"
)

type Request struct {
	Action string `json:"action"`
}

type Response struct {
	StatusCode int         `json:"statusCode"`
	Body       interface{} `json:"body"`
}

func errorMsg(e error) (*Response, error) {
	return &Response{
		StatusCode: 500,
		Body:       e.Error(),
	}, nil
}

func QueueHandler(ctx context.Context, r queueRequest) (any, error) {

	fmt.Println(r)

	return r.Messages[0].Details.Message.Body, nil
}

func WebhookHandler(ctx context.Context, r *Request) (*Response, error) {
	if r.Action == "default" {
		return defaultHandler(ctx, r, nil)
	}
	connectors, close, err := InitTransport(ctx)
	if err != nil {
		return errorMsg(err)
	}
	defer close()
	if r.Action == "get_worlds" {
		return worldsHandler(ctx, r, connectors)
	}
	return &Response{
		StatusCode: 200,
		Body:       "TestQ: invalid",
	}, nil
}

func worldsHandler(ctx context.Context, r *Request, connections *transport) (*Response, error) {
	descs, err := GetWorldDescriptions(ctx, connections)
	if err != nil {
		return errorMsg(err)
	}
	b, err := json.Marshal(descs)
	if err != nil {
		return errorMsg(err)
	}
	return &Response{
		StatusCode: 200,
		Body:       string(b),
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
	req := &Request{
		Action: "get_worlds",
	}
	resp, err := WebhookHandler(ctx, req)
	if err != nil {
		panic(err)
	}
	fmt.Println(resp)
}
