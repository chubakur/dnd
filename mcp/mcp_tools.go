package mcp

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/chubakur/dnd/transport"
	"github.com/chubakur/dnd/types"
	"github.com/google/uuid"
)

type mcpToolFunc func(*transport.Transport, *types.DeepseekResponseToolCall) (string, error)

type MCPTool struct {
	Name        string                  `json:"name"`
	Description string                  `json:"description"`
	Parameters  types.MCPToolParameters `json:"parameters"`
	F           mcpToolFunc
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

type WorldDescription struct {
	Id          int64  `sql:"id" json:"id"`
	Status      int8   `sql:"status" json:"status"`
	Name        string `sql:"name" json:"name"`
	Description string `sql:"description" json:"description"`
}

type PlayerSession struct {
	PlayerId  uuid.UUID `sql:"player_id" json:"player_id"`
	SessionId uuid.UUID `sql:"session_id" json:"session_id"`
	DmId      uuid.UUID `sql:"dm_id" json:"dm_id"`
}

func mcptool_getWorldDescriptions(t *transport.Transport, p *types.DeepseekResponseToolCall) (string, error) {
	// Преобразуем кистинному контексту
	ctx, ok := t.Ctx.(context.Context)
	if !ok {
		ctx = context.Background()
	}

	descs, err := getWorldDescriptions(ctx, t)
	if err != nil {
		return "", err
	}

	bytes, err := json.Marshal(descs)
	if err != nil {
		return "", err
	}

	return string(bytes), nil
}

func getWorldDescriptions(ctx context.Context, t *transport.Transport) ([]WorldDescription, error) {
	var worldDescriptions = make([]WorldDescription, 0)
	ydbClient := t.YdbClient

	rows, err := ydbClient.Query().QueryResultSet(ctx, "SELECT id, status, name, description FROM world_descriptions")
	if err != nil {
		return worldDescriptions, err
	}
	defer rows.Close(ctx)

	// Используем правильный паттерн для YDB с каналом
	rowChan := rows.Rows(ctx)
	for row := range rowChan {
		var wd WorldDescription
		err = row.ScanStruct(&wd)
		if err != nil {
			return worldDescriptions, err
		}
		worldDescriptions = append(worldDescriptions, wd)
	}

	return worldDescriptions, nil
}

func mcptool_getActiveSessions(t *transport.Transport, d *types.DeepseekResponseToolCall) (string, error) {
	type reqt struct {
		PlayerId string `json:"player_id"`
	}
	var req reqt
	decoder := json.NewDecoder(strings.NewReader(d.Function.Arguments))
	err := decoder.Decode(&req)
	if err != nil {
		return "", err
	}
	uuid, err := uuid.Parse(req.PlayerId)
	if err != nil {
		return "", err
	}
	activeSessions, err := getActivePlayerSessions(t, uuid)
	if err != nil {
		return "", err
	}
	bytes, err := json.Marshal(activeSessions)
	if err != nil {
		return "", err
	}

	return string(bytes), nil
}

func getActivePlayerSessions(t *transport.Transport, playerId uuid.UUID) ([]PlayerSession, error) {
	ctx, ok := t.Ctx.(context.Context)
	if !ok {
		ctx = context.Background()
	}

	var result []PlayerSession = make([]PlayerSession, 0)
	query := fmt.Sprintf("SELECT player_id, session_id, dm_id FROM sessions WHERE player_id = Uuid('%s')", playerId.String())
	fmt.Println(query)

	// Type assertion to get the YDB client
	ydbClient := t.YdbClient

	res, err := ydbClient.Query().QueryResultSet(ctx, query)
	if err != nil {
		return result, err
	}

	defer res.Close(ctx)

	// Используем правильный паттерн для YDB с каналом
	rowChan := res.Rows(ctx)
	for row := range rowChan {
		var o PlayerSession
		err = row.ScanStruct(&o)
		if err != nil {
			return result, err
		}
		result = append(result, o)
	}

	return result, nil
}

func MCPGetTools() []*MCPTool {
	var tools []*MCPTool = []*MCPTool{
		{
			Name:        "get_worlds",
			Description: "Get available game settings for DnD party",
			Parameters: types.MCPToolParameters{
				Type:       "object",
				Properties: make([]types.MCPToolProperty, 0),
			},
			F: mcptool_getWorldDescriptions,
		},
		{
			Name:        "get_user_sessions",
			Description: "Get active sessions for current player",
			Parameters: types.MCPToolParameters{
				Type: "object",
				Properties: []types.MCPToolProperty{
					{
						Name:        "player_id",
						Type:        "string",
						Description: "UUID of current player.",
						IsRequired:  true,
					},
				},
			},
			F: mcptool_getActiveSessions,
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

func MCPCall(t *transport.Transport, toolQuery types.DeepseekResponseToolCall) MCPResult {
	tools := MCPGetTools()
	for _, tool := range tools {
		if tool.Name == toolQuery.Function.Name {
			toolResult, err := tool.F(t, &toolQuery)
			if err != nil {
				return MCPResult{
					Function:   toolQuery.Function.Name,
					Error:      err,
					ToolCallId: toolQuery.Id,
				}
			}
			return MCPResult{
				Function:   toolQuery.Function.Name,
				Result:     toolResult,
				ToolCallId: toolQuery.Id,
			}
		}
	}

	return MCPResult{
		Function:   toolQuery.Function.Name,
		Error:      fmt.Errorf("Error: tool %s not found", toolQuery.Function.Name),
		ToolCallId: toolQuery.Id,
	}
}
