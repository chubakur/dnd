package main

import (
	"context"
	"encoding/json"
	"fmt"
)

type mcpToolFunc func(context.Context, *transport) (string, error)

type MCPTool struct {
	Name        string            `json:"name"`
	Description string            `json:"description"`
	Parameters  mcpToolParameters `json:"parameters"`
	f           mcpToolFunc
}

type mcpToolParameters struct {
	Type       string
	Properties []mcpToolProperty
}

type mcpToolProperty struct {
	Name        string
	Type        string
	Description string
	IsRequired  bool
}

type wrappedMCPFunction struct {
	Type        string                       `json:"type"`
	Name        string                       `json:"name"`
	Description string                       `json:"description"`
	Parameters  wrappedMCPFunctionParameters `json:"parameters"`
}

type wrappedMCPFunctionParametersProperty struct {
	Type        string `json:"type"`
	Description string `json:"description"`
}

type wrappedMCPFunctionParameters struct {
	Type       string                                          `json:"type"`
	Required   []string                                        `json:"required"`
	Properties map[string]wrappedMCPFunctionParametersProperty `json:"properties"`
}

type wrappedMCPTool struct {
	Type     string             `json:"type"`
	Function wrappedMCPFunction `json:"function"`
}

func (mc *MCPTool) wrapJson() wrappedMCPTool {
	result := wrappedMCPTool{
		Type: "function",
		Function: wrappedMCPFunction{
			Type:        mc.Parameters.Type,
			Name:        mc.Name,
			Description: mc.Description,
			Parameters: wrappedMCPFunctionParameters{
				Type:       "object",
				Properties: make(map[string]wrappedMCPFunctionParametersProperty),
				Required:   make([]string, 0),
			},
		},
	}

	for _, param := range mc.Parameters.Properties {
		result.Function.Parameters.Properties[param.Name] = wrappedMCPFunctionParametersProperty{
			Type:        param.Type,
			Description: param.Description,
		}
		if param.IsRequired {
			result.Function.Parameters.Required = append(result.Function.Parameters.Required, param.Name)
		}
	}

	return result
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
			Parameters: mcpToolParameters{
				Type:       "object",
				Properties: make([]mcpToolProperty, 0),
			},
			f: mcptool_getWorldDescriptions,
		},
		{
			Name:        "get_user_sessions",
			Description: "Get active sessions for current player",
			Parameters: mcpToolParameters{
				Type: "object",
				Properties: []mcpToolProperty{
					{
						Name:        "user_id",
						Type:        "string",
						Description: "UUID of current player.",
						IsRequired:  true,
					},
				},
			},
			f: func(ctx context.Context, t *transport) (string, error) {
				panic("Not implemented")
			},
		},
	}

	return tools
}

type MCPResult struct {
	Function   string
	Result     string
	Error      error
	ToolCallId string
}

func MCPCall(ctx context.Context, t *transport, tool_query deepseekResponseToolCall) MCPResult {
	tools := MCPGetTools()
	for _, tool := range tools {
		if tool.Name == tool_query.Function.Name {
			tool_result, err := tool.f(ctx, t)
			if err != nil {
				return MCPResult{
					Function:   tool_query.Function.Name,
					Error:      err,
					ToolCallId: tool_query.Id,
				}
			}
			return MCPResult{
				Function:   tool_query.Function.Name,
				Result:     tool_result,
				ToolCallId: tool_query.Id,
			}
		}
	}

	return MCPResult{
		Function:   tool_query.Function.Name,
		Error:      fmt.Errorf("Error: tool %s not found", tool_query.Function.Name),
		ToolCallId: tool_query.Id,
	}
}
