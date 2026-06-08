package main

import (
	"context"
	"fmt"
	"log/slog"
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
	slog.InfoContext(ctx, "webhook received", "tg_id", r.TgId)

	t, close, err := transport.InitTransport(ctx)
	if err != nil {
		slog.ErrorContext(ctx, "transport init failed", "err", err)
		return errorMsg(err)
	}
	defer close()

	apiKey := os.Getenv("DEEPSEEK_API_KEY")
	if apiKey == "" {
		return errorMsg(fmt.Errorf("DEEPSEEK_API_KEY not set"))
	}

	tools := mcp.MCPGetTools()
	client := llmcore.NewDeepSeekClient(apiKey, tools)
	mc := llmcore.NewMessageChain()
	mc.AddUserMessage(r.Message)

	res, err := client.Query(mc)
	if err != nil {
		slog.ErrorContext(ctx, "llm query failed", "err", err)
		return errorMsg(err)
	}

	if len(res.Choices) > 0 && len(res.Choices[0].Message.ToolCalls) > 0 {
		mc.AddMessage(res.Choices[0].Message)
		for _, toolCall := range res.Choices[0].Message.ToolCalls {
			slog.InfoContext(ctx, "mcp tool call", "tool", toolCall.Function.Name)
			mcpResult := mcp.MCPCall(t, toolCall)
			mc.AddToolMessage(mcpResult)
		}
		res, err = client.Query(mc)
		if err != nil {
			slog.ErrorContext(ctx, "llm query after tools failed", "err", err)
			return errorMsg(err)
		}
	}

	slog.InfoContext(ctx, "webhook response sent", "tg_id", r.TgId)
	return &Response{
		StatusCode: 200,
		Body:       res.Choices[0].Message.Content,
	}, nil
}
