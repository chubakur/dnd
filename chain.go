package main

type messageChain struct {
	chain []deepSeekRoleContent
}

func newMessageChain() *messageChain {
	return &messageChain{}
}

func (mc *messageChain) addUserMessage(message string) {
	mc.chain = append(mc.chain, deepSeekRoleContent{Role: "user", Content: message})
}

func (mc *messageChain) addToolMessage(mcpRes MCPResult) {
	mc.chain = append(mc.chain, deepSeekRoleContent{Role: "tool", Content: mcpRes.Result, ToolCallId: mcpRes.ToolCallId})
}

func (mc *messageChain) addMessage(msg deepSeekRoleContent) {
	mc.chain = append(mc.chain, msg)
}
