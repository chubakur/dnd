package messages

import (
	"fmt"
	"time"

	"github.com/chubakur/dnd/llmcore"
	"github.com/chubakur/dnd/transport"
	"github.com/google/uuid"
)

type Message struct {
	PlayerId   uuid.UUID `sql:"player_id"`
	ChatId     uuid.UUID `sql:"chat_id"`
	MessageId  uuid.UUID `sql:"message_id"`
	Time       time.Time `sql:"time" ydb:"time"`
	Role       string    `sql:"role"`
	Content    string    `sql:"content"`
	IsArchived bool      `sql:"is_archived"`
}

type Chat struct {
	PlayerId uuid.UUID `sql:"player_id"`
	ChatId   uuid.UUID `sql:"chat_id"`
	Status   int8      `sql:"status"`
}

func GetMessagesByChatId(t *transport.Transport, playerId, chatId uuid.UUID) ([]Message, error) {
	messages := make([]Message, 0)
	sqlQuery := fmt.Sprintf("SELECT player_id, chat_id, time, message_id, role, content, is_archived FROM messages WHERE player_id = Uuid('%s') AND chat_id = Uuid('%s') AND is_archived = false ORDER BY message_id", playerId, chatId)
	ydbResult, err := t.YdbClient.Query().QueryResultSet(t.Ctx, sqlQuery)
	if err != nil {
		return messages, err
	}
	defer ydbResult.Close(t.Ctx)
	rowsIter := ydbResult.Rows(t.Ctx)
	for row := range rowsIter {
		var msg Message
		err = row.ScanStruct(&msg)
		if err != nil {
			return messages, err
		}
		messages = append(messages, msg)
	}

	return messages, nil
}

func Write(t *transport.Transport, playerId, chatId uuid.UUID, message llmcore.DeepSeekRoleContent) (bool, error) {
	newMsgId := uuid.New()
	sqlQuery := fmt.Sprintf("INSERT INTO messages (player_id, chat_id, message_id, time, role, content) VALUES (Uuid(\"%s\"), Uuid(\"%s\"), Uuid(\"%s\"), CurrentUtcDatetime(), \"%s\", \"%s\")",
		playerId,
		chatId,
		newMsgId,
		message.Role,
		message.Content,
	)
	err := t.YdbClient.Query().Exec(t.Ctx, sqlQuery)
	if err != nil {
		return false, err
	}
	return true, nil
}
