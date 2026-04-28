package async

import (
	"github.com/chubakur/dnd/llmcore"
	"github.com/chubakur/dnd/messages"
	"github.com/chubakur/dnd/transport"
	"github.com/chubakur/dnd/types"
	"github.com/google/uuid"
)

type AsyncTaskChatLlmStruct struct {
	AsyncTaskCommon
	PlayerId uuid.UUID `json:"player_id"`
	ChatId   uuid.UUID `json:"chat_id"`
}

func (a *AsyncTaskChatLlmStruct) Handle(t *transport.Transport) (*llmcore.MessageChain, error) {
	chain := llmcore.NewMessageChain()
	messages, err := messages.GetMessagesByChatId(t, a.PlayerId, a.ChatId)
	if err != nil {
		return chain, err
	}
	for _, message := range messages {
		contentWithRole := types.DeepSeekRoleContent{
			Role:    message.Role,
			Content: message.Content,
		}
		chain.AddMessage(contentWithRole)
	}

	return chain, err
}
