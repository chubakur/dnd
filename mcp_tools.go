package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
)

type mcpToolFunc func(context.Context, *transport) (string, error)

type MCPTool struct {
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Parameters  map[string]any `json:"parameters"`
	f           mcpToolFunc
}

type WrappedMCPTool struct {
	Type     string   `json:"type"`
	Function *MCPTool `json:"function"`
}

func mcptool_getWorldDescriptions(ctx context.Context, t *transport) (string, error) {
	descs, err := GetWorldDescriptions(ctx, t)
	if err != nil {
		return "", err
	}

	bytes, err := json.Marshal(descs)
	if err != nil {
		return "", err
	}

	return string(bytes), nil
}

func MCPGetTools() []*MCPTool {
	var tools []*MCPTool = []*MCPTool{
		{
			Name:        "get_worlds",
			Description: "Get available game settings for DnD party",
			Parameters: map[string]any{
				"type":       "object",
				"properties": map[string]string{},
				"required":   []string{},
			},
			f: mcptool_getWorldDescriptions,
		},
	}

	return tools
}

type MCPResult struct {
	Function string
	Result   string
	Error    error
}

func MCPCall(ctx context.Context, t *transport, tool_query deepseekResponseToolCall) MCPResult {
	tools := MCPGetTools()
	for _, tool := range tools {
		if tool.Name == tool_query.Function.Name {
			tool_result, err := tool.f(ctx, t)
			if err != nil {
				return MCPResult{
					Function: tool_query.Function.Name,
					Error:    err,
				}
			}
			return MCPResult{
				Function: tool_query.Function.Name,
				Result:   tool_result,
			}
		}
	}

	return MCPResult{
		Function: tool_query.Function.Name,
		Error:    errors.New(fmt.Sprintf("Error: tool %s not found", tool_query.Function.Name)),
	}
}
