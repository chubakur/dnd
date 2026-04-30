package main

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/chubakur/dnd/messages"
	"github.com/chubakur/dnd/transport"
	"github.com/google/uuid"
)

type messagesRequest struct {
	Action string `json:"action"`
	Params string `json:"params"`
}

type createMessageRequestParams struct {
	PlayerId uuid.UUID `json:"player_id"`
	ChatId   uuid.UUID `json:"chat_id"`
	Content  string    `json:"content"`
	Role     string    `json:"role"`
}

func MessagesHandler(ctx context.Context, req *messagesRequest) (*Response, error) {
	t, d, e := transport.InitTransport(ctx)
	if e != nil {
		return errorMsg(e)
	}
	defer d()
	switch req.Action {
	case "add":
		decoder := json.NewDecoder(strings.NewReader(req.Params))
		var params createMessageRequestParams
		err := decoder.Decode(&params)
		if err != nil {
			return errorMsg(err)
		}
		newMsg := messages.Message{
			MessageId: uuid.New(),
			PlayerId:  params.PlayerId,
			ChatId:    params.ChatId,
			Role:      params.Role,
			Content:   params.Content,
		}
		err = messages.AddToChat(t, &newMsg)
		if err != nil {
			return errorMsg(err)
		}
		jmsg, err := json.Marshal(newMsg)
		if err != nil {
			return errorMsg(err)
		}
		return &Response{
			StatusCode: 200,
			Body:       string(jmsg),
		}, nil
	}
	return errorMsg(fmt.Errorf("Invalid action: %s", req.Action))
}
