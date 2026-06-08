package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/chubakur/dnd/dndcore"
	"github.com/chubakur/dnd/llmcore"
	"github.com/chubakur/dnd/mcp"
	"github.com/chubakur/dnd/transport"
	"github.com/google/uuid"
)

type MasterJob interface {
	GetName() string
	Execute(t *transport.Transport) error
}

// GenerateWorldJob asks the LLM to generate a new world and saves it via MCP tools.
type GenerateWorldJob struct {
	PlayerId uuid.UUID `json:"player_id"`
	Prompt   string    `json:"prompt"`
}

func (j *GenerateWorldJob) GetName() string { return "generate_world" }

func (j *GenerateWorldJob) Execute(t *transport.Transport) error {
	apiKey := os.Getenv("DEEPSEEK_API_KEY")
	if apiKey == "" {
		return fmt.Errorf("DEEPSEEK_API_KEY not set")
	}
	tools := mcp.MCPGetTools()
	client := llmcore.NewDeepSeekClient(apiKey, tools)
	chain := llmcore.NewMessageChain()
	chain.AddSystemMessage(dndcore.DMSystemPrompt("", ""))
	chain.AddUserMessage(fmt.Sprintf(
		"Generate a new D&D world for player %s. Prompt: %s. Use the available tools to save the session.",
		j.PlayerId.String(), j.Prompt,
	))
	_, _, err := client.AgentExecutor(t, nil, chain, 5)
	return err
}

func ParseJob(data []byte) (MasterJob, error) {
	var raw struct {
		Name string `json:"name"`
	}
	if err := json.Unmarshal(data, &raw); err != nil {
		return nil, err
	}
	switch raw.Name {
	case "generate_world":
		var job GenerateWorldJob
		if err := json.Unmarshal(data, &job); err != nil {
			return nil, err
		}
		return &job, nil
	}
	return nil, fmt.Errorf("unknown job: %s", raw.Name)
}
