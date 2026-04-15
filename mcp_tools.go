package main

import "context"

type MCPTool struct {
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Parameters  map[string]any `json:"parameters"`
	// Func        func()
}

type WrappedMCPTool struct {
	Type     string   `json:"type"`
	Function *MCPTool `json:"function"`
}

func MCPGetTools(ctx context.Context, connections *transport) []*MCPTool {
	var tools []*MCPTool = []*MCPTool{
		{
			Name:        "get_worlds",
			Description: "Get available game settings for DnD party",
			Parameters: map[string]any{
				"type":       "object",
				"properties": map[string]string{},
				"required":   []string{},
			},
			// Func: func() {
			// 	GetWorldDescriptions(ctx, connections)
			// },
		},
	}

	return tools
}
