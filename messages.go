package main

import (
	"github.com/google/uuid"
)

type Message struct {
	PlayerId  uuid.UUID `sql:"player_id"`
	ChatId    uuid.UUID `sql:"chat_id"`
	MessageId uuid.UUID `sql:"message_id"`
	// Time       datetime.DateTime `sql:"time"`
	Role       string `sql:"role"`
	Content    string `sql:"content"`
	IsArchived bool   `sql:"archived"`
}

type Chat struct {
	PlayerId uuid.UUID `sql:"player_id"`
	ChatId   uuid.UUID `sql:"chat_id"`
	Status   int8      `sql:"status"`
}
