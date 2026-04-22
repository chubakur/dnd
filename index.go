package main

import (
	"context"
	"fmt"
	"os"

	"github.com/chubakur/dnd/llmcore"
	"github.com/chubakur/dnd/mcp"
	"github.com/chubakur/dnd/transport"
	"github.com/chubakur/dnd/types"
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
		Body:       req.Message,
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

	// Инициализация LLM клиента
	tools := mcp.MCPGetTools()
	client := llmcore.NewDeepSeekClient(apiKey, tools)

	// Создаем цепочку сообщений
	mc := llmcore.NewMessageChain()
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
