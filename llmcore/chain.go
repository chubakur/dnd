package llmcore

import (
	"github.com/chubakur/dnd/types"
)

type MessageChain struct {
	chain []DeepSeekRoleContent
}

func NewMessageChain() *MessageChain {
	return &MessageChain{}
}

func (mc *MessageChain) AddUserMessage(message string) {
	mc.chain = append(mc.chain, DeepSeekRoleContent{Role: "user", Content: message})
}

func (mc *MessageChain) AddSystemMessage(message string) {
	mc.chain = append(mc.chain, DeepSeekRoleContent{Role: "system", Content: message})
}

func (mc *MessageChain) AddToolMessage(mcpRes types.MCPResult) {
	mc.chain = append(mc.chain, DeepSeekRoleContent{Role: "tool", Content: mcpRes.Result, ToolCallId: mcpRes.ToolCallId})
}

func (mc *MessageChain) AddMessage(msg DeepSeekRoleContent) {
	mc.chain = append(mc.chain, msg)
}
