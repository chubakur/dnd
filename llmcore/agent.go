package llmcore

import (
	"fmt"

	"github.com/chubakur/dnd/dndcore"
	"github.com/chubakur/dnd/messages"
	"github.com/chubakur/dnd/transport"
	"github.com/chubakur/dnd/types"
	"github.com/google/uuid"
)

const agentExecutorLimit = 10

// DeepSeekAgent implements dndcore.LlmAgent backed by the DeepSeek API.
type DeepSeekAgent struct {
	t      *transport.Transport
	client *deepSeekClient
	chain  *MessageChain
}

func NewDeepSeekAgent(t *transport.Transport, client *deepSeekClient, systemPrompt string) *DeepSeekAgent {
	chain := NewMessageChain()
	if systemPrompt != "" {
		chain.AddSystemMessage(systemPrompt)
	}
	return &DeepSeekAgent{t: t, client: client, chain: chain}
}

func (a *DeepSeekAgent) SendMessage(playerId, chatId uuid.UUID, message string) (string, error) {
	a.chain.AddUserMessage(message)
	err := messages.Write(a.t, playerId, chatId, types.DeepSeekRoleContent{Role: "user", Content: message})
	if err != nil {
		return "", fmt.Errorf("write user message: %w", err)
	}
	gc := &dndcore.GameContext{PlayerId: playerId, ChatId: chatId}
	_, resp, err := a.client.AgentExecutor(a.t, gc, a.chain, agentExecutorLimit)
	if err != nil {
		return "", err
	}
	choice := resp.GetFirstChoice()
	if choice == nil {
		return "", fmt.Errorf("empty response from LLM")
	}
	return choice.Message.Content, nil
}

func (a *DeepSeekAgent) GetHistory(playerId, chatId uuid.UUID) ([]string, error) {
	msgs, err := messages.GetMessagesByChatId(a.t, playerId, chatId)
	if err != nil {
		return nil, err
	}
	result := make([]string, 0, len(msgs))
	for _, m := range msgs {
		result = append(result, fmt.Sprintf("%s: %s", m.Role, m.Content))
	}
	return result, nil
}

func (a *DeepSeekAgent) ResetContext() {
	a.chain = NewMessageChain()
}
