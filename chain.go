package main

type messageChain struct {
	chain []deepSeekRoleContent
}

func newMessageChain() *messageChain {
	return &messageChain{}
}

func (*messageChain) addSystemMessage(message string) {

}

func (mc *messageChain) addUserMessage(message string) {
	mc.chain = append(mc.chain, deepSeekRoleContent{Role: "user", Content: message})
}

func (*messageChain) addAssistantMessage(message string) {

}

func (mc *messageChain) addToolMessage(message, toolCallId string) {
	mc.chain = append(mc.chain, deepSeekRoleContent{Role: "tool", Content: message, ToolCallId: toolCallId})
}

func (mc *messageChain) addMessage(msg deepSeekRoleContent) {
	mc.chain = append(mc.chain, msg)
}
