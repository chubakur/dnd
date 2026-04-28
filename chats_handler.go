package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/chubakur/dnd/chats"
	"github.com/chubakur/dnd/transport"
	"github.com/google/uuid"
)

type chatsRequest struct {
	PlayerId string `json:"player_id"`
}

func ChatsHandler(ctx context.Context, req *chatsRequest) (*Response, error) {
	if req.PlayerId == "" {
		return errorMsg(fmt.Errorf("Empty player_id"))
	}
	playerId, err := uuid.Parse(req.PlayerId)
	if err != nil {
		return errorMsg(err)
	}
	t, c, e := transport.InitTransport(ctx)
	if e != nil {
		return errorMsg(e)
	}
	defer c()
	chat, err := chats.GetActive(t, playerId)
	if err != nil {
		return errorMsg(err)
	}
	jbind, err := json.Marshal(chat)
	if err != nil {
		return errorMsg(err)
	}
	return &Response{
		StatusCode: 200,
		Body:       string(jbind),
	}, nil
}
