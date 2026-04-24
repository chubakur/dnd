package main

import (
	"context"
	"fmt"
	"os"

	"github.com/chubakur/dnd/llmcore"
	"github.com/chubakur/dnd/mcp"
	"github.com/chubakur/dnd/transport"
)

type TgRequest struct {
	Message string `json:"message"`
	TgId    string `json:"tg_id"`
}

func WebhookHandler(ctx context.Context, r *TgRequest) (*Response, error) {
	t, close, err := transport.InitTransport(ctx)
	if err != nil {
		return errorMsg(err)
	}
	defer close()

	apiKey := os.Getenv("DEEPSEEK_API_KEY")
	if apiKey == "" {
		return errorMsg(fmt.Errorf("DEEPSEEK_API_KEY not set"))
	}

	// Получаем MCP tools
	tools := mcp.MCPGetTools()

	// Инициализируем LLM клиент
	client := llmcore.NewDeepSeekClient(apiKey, tools)

	// Создаем цепочку сообщений
	mc := llmcore.NewMessageChain()
	mc.AddUserMessage(r.Message)

	// Выполняем запрос
	res, err := client.Query(mc)
	if err != nil {
		return errorMsg(err)
	}

	// Проверяем, есть ли tool calls
	if len(res.Choices) > 0 && len(res.Choices[0].Message.ToolCalls) > 0 {
		mc.AddMessage(res.Choices[0].Message)

		// Обрабатываем tool calls
		for _, toolCall := range res.Choices[0].Message.ToolCalls {
			mcpResult := mcp.MCPCall(t, toolCall)
			mc.AddToolMessage(mcpResult)
		}

		// Второй запрос с результатами tools
		res, err = client.Query(mc)
		if err != nil {
			return errorMsg(err)
		}
	}

	return &Response{
		StatusCode: 200,
		Body:       res.Choices[0].Message.Content,
	}, nil
}
