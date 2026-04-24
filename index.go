package main

import (
	"context"
	"fmt"
	"os"

	"github.com/chubakur/dnd/async"
	"github.com/chubakur/dnd/dndcore"
	"github.com/chubakur/dnd/llmcore"
	"github.com/chubakur/dnd/mcp"
	"github.com/chubakur/dnd/transport"
	"github.com/google/uuid"
)

type Request struct {
	Message string `json:"message"`
}

type Response struct {
	StatusCode int `json:"statusCode"`
	Body       any `json:"body"`
}

func QueueHandler(ctx context.Context, req *Request) (*Response, error) {
	return &Response{
		StatusCode: 200,
		Body:       req.Message + " Ok",
	}, nil
}

func errorMsg(e error) (*Response, error) {
	return &Response{
		StatusCode: 500,
		Body:       e.Error(),
	}, nil
}

func main() {
	apiKey := os.Getenv("DEEPSEEK_API_KEY")
	if apiKey == "" {
		panic("Set DEEPSEEK_API_KEY")
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	t, close, err := transport.InitTransport(ctx)
	if err != nil {
		panic(err)
	}
	defer close()

	asyncTask := async.AsyncTaskChatLlmStruct{
		PlayerId: uuid.MustParse("4e54ac3e-9c91-4dc0-a582-9439f8756a3a"),
		ChatId:   uuid.MustParse("40fce1fe-5b56-424c-b1e2-1da6e5b4422d"),
	}
	mc, err := asyncTask.Handle(t)
	if err != nil {
		panic(err)
	}
	tools := mcp.MCPGetTools()
	client := llmcore.NewDeepSeekClient(apiKey, tools)
	pc := &dndcore.GameContext{
		PlayerId: uuid.MustParse("4e54ac3e-9c91-4dc0-a582-9439f8756a3a"),
		ChatId:   uuid.MustParse("40fce1fe-5b56-424c-b1e2-1da6e5b4422d"),
	}
	mc, response, err := client.AgentExecutor(t, pc, mc, 3)
	if err != nil {
		panic(err)
	}
	fmt.Println(response)
	fmt.Println(response.Choices[0].Message.Content)
	os.Exit(0)
}
