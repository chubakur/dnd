package chats

import (
	"errors"
	"fmt"
	"io"
	"time"

	"github.com/chubakur/dnd/transport"
	"github.com/google/uuid"
	"github.com/ydb-platform/ydb-go-sdk/v3/query"
)

type Chat struct {
	PlayerId   uuid.UUID `sql:"player_id" json:"player_id"`
	ChatId     uuid.UUID `sql:"chat_id" json:"chat_id"`
	UpdateTime time.Time `sql:"update_time" json:"update_time"`
}

func GetActive(t *transport.Transport, playerId uuid.UUID) (*Chat, error) {
	querySql := fmt.Sprintf("SELECT player_id, chat_id, update_time FROM chats WHERE player_id = Uuid('%s') ORDER BY update_time DESC LIMIT 1", playerId.String())
	res, err := t.YdbClient.Query().QueryResultSet(t.Ctx, querySql, query.WithIdempotent())
	if err != nil {
		return nil, err
	}
	defer res.Close(t.Ctx)
	row, err := res.NextRow(t.Ctx)
	if err != nil {
		if errors.Is(err, io.EOF) {
			return nil, nil
		}
		return nil, err
	}
	var activeChat Chat
	err = row.ScanStruct(&activeChat)
	if err != nil {
		return nil, err
	}

	return &activeChat, nil
}

func CreateNew(t *transport.Transport, playerId uuid.UUID) (*Chat, error) {
	newChatId := uuid.New()
	sqlQuery := fmt.Sprintf("INSERT INTO chats (player_id, chat_id, update_time) VALUES (Uuid('%s'), Uuid('%s'), CurrentUtcDatetime())",
		playerId.String(), newChatId.String())
	err := t.YdbClient.Query().Exec(t.Ctx, sqlQuery)
	if err != nil {
		return nil, err
	}

	// Return the created chat
	return &Chat{
		PlayerId:   playerId,
		ChatId:     newChatId,
		UpdateTime: time.Now().UTC(),
	}, nil
}

func MakeActive(playerId, chatId uuid.UUID) error {
	return nil
}
