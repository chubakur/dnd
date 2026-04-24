package llmcore

import (
	"github.com/chubakur/dnd/mcp"
	"github.com/chubakur/dnd/types"
)

type MessageChain struct {
	chain []types.DeepSeekRoleContent
}

func NewMessageChain() *MessageChain {
	return &MessageChain{}
}

func (mc *MessageChain) AddUserMessage(message string) {
	mc.chain = append(mc.chain, types.DeepSeekRoleContent{Role: "user", Content: message})
}

func (mc *MessageChain) AddSystemMessage(message string) {
	mc.chain = append(mc.chain, types.DeepSeekRoleContent{Role: "system", Content: message})
}

func (mc *MessageChain) AddToolMessage(mcpRes mcp.MCPResult) {
	mc.chain = append(mc.chain, types.DeepSeekRoleContent{Role: "tool", Content: mcpRes.Result, ToolCallId: mcpRes.ToolCallId})
}

func (mc *MessageChain) AddMessage(msg types.DeepSeekRoleContent) {
	mc.chain = append(mc.chain, msg)
}
