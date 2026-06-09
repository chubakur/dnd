package async

import (
	"github.com/chubakur/dnd/dndcore"
	"github.com/chubakur/dnd/llmcore"
	"github.com/chubakur/dnd/messages"
	"github.com/chubakur/dnd/transport"
	"github.com/chubakur/dnd/types"
	"github.com/google/uuid"
)

type AsyncTaskChatLlmStruct struct {
	AsyncTaskCommon
	ExecutionId uuid.UUID `json:"execution_id"`
	PlayerId    uuid.UUID `json:"player_id"`
	ChatId      uuid.UUID `json:"chat_id"`
}

// Handle builds a MessageChain for the LLM call: a DM system prompt followed
// by the stored chat history (the user's latest message is already persisted
// before this job is enqueued).
func (a *AsyncTaskChatLlmStruct) Handle(t *transport.Transport) (*llmcore.MessageChain, error) {
	chain := llmcore.NewMessageChain()
	chain.AddSystemMessage(dndcore.DMSystemPrompt("", ""))
	history, err := messages.GetMessagesByChatId(t, a.PlayerId, a.ChatId)
	if err != nil {
		return chain, err
	}
	for _, message := range history {
		chain.AddMessage(types.DeepSeekRoleContent{
			Role:    message.Role,
			Content: message.Content,
		})
	}
	return chain, nil
}
