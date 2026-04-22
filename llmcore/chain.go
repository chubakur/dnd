package llmcore

import (
	"github.com/chubakur/dnd/types"
)

type messageChain struct {
	chain []deepSeekRoleContent
}

func NewMessageChain() *messageChain {
	return &messageChain{}
}

func (mc *messageChain) AddUserMessage(message string) {
	mc.chain = append(mc.chain, deepSeekRoleContent{Role: "user", Content: message})
}

func (mc *messageChain) AddSystemMessage(message string) {
	mc.chain = append(mc.chain, deepSeekRoleContent{Role: "system", Content: message})
}

func (mc *messageChain) AddToolMessage(mcpRes types.MCPResult) {
	mc.chain = append(mc.chain, deepSeekRoleContent{Role: "tool", Content: mcpRes.Result, ToolCallId: mcpRes.ToolCallId})
}

func (mc *messageChain) AddMessage(msg deepSeekRoleContent) {
	mc.chain = append(mc.chain, msg)
}
