package main

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/chubakur/dnd/async"
	"github.com/chubakur/dnd/llmcore"
	"github.com/chubakur/dnd/mcp"
	"github.com/chubakur/dnd/messages"
	"github.com/chubakur/dnd/transport"
	"github.com/chubakur/dnd/types"
	"github.com/google/uuid"
	"github.com/ydb-platform/ydb-go-sdk/v3/topic/topicwriter"
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
	response, err := client.Query(mc)
	if err != nil {
		panic(err)
	}
	fmt.Println(response)
	fmt.Println(response.Choices[0].Message.Content)
	os.Exit(0)

	pid := uuid.New()
	cid := uuid.New()
	message := llmcore.DeepSeekRoleContent{
		Role:    "user",
		Content: "Привет, подскажи кто ты есть и что умеешь.",
	}
	_, err = messages.Write(t, pid, cid, message)
	if err != nil {
		panic(err)
	}
	os.Exit(0)

	// err = AddConsumer(t)
	// if err != nil {
	// 	panic(err)
	// }
	// err = ConsumeMsg(t)
	// if err != nil {
	// 	panic(err)
	// }

	// os.Exit(0)

	writer, err := t.YdbClient.Topic().StartWriter("jobs")
	if err != nil {
		panic(err)
	}
	err = writer.Write(ctx, topicwriter.Message{Data: strings.NewReader("dsfsxcc123d")})
	if err != nil {
		panic(err)
	}
	err = writer.Flush(ctx)
	if err != nil {
		panic(err)
	}

	os.Exit(0)

	// Создаем цепочку сообщений
	mc.AddUserMessage("Привет, дружище, подскажи, какие сеттинги для игры ты знаешь?")
	res, err := client.Query(mc)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Response: %+v\n", res)

	// Обработка tool calls
	if len(res.Choices) > 0 && len(res.Choices[0].Message.ToolCalls) > 0 {
		mc.AddMessage(res.Choices[0].Message)
		for _, toolCall := range res.Choices[0].Message.ToolCalls {
			// Создаем transport adapter
			transport := &types.Transport{
				YdbClient: t.YdbClient,
				Ctx:       ctx,
			}
			mcpResult := mcp.MCPCall(transport, toolCall)
			fmt.Printf("MCP Result: %+v\n", mcpResult)
			mc.AddToolMessage(mcpResult)
		}
	}

	// Второй запрос с результатами tools
	res2, err := client.Query(mc)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Final response: %s\n", res2.Choices[0].Message.Content)
}
