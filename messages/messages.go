package messages

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/chubakur/dnd/transport"
	"github.com/chubakur/dnd/types"
	"github.com/google/uuid"
)

type Message struct {
	PlayerId   uuid.UUID `sql:"player_id"`
	ChatId     uuid.UUID `sql:"chat_id"`
	MessageId  uuid.UUID `sql:"message_id"`
	Time       time.Time `sql:"time" ydb:"time"`
	Role       string    `sql:"role"`
	Content    string    `sql:"content"`
	ToolCalls  string    `sql:"tool_calls"`
	ToolCallId string    `sql:"tool_call_id"`
	IsArchived bool      `sql:"is_archived"`
}

type Chat struct {
	PlayerId uuid.UUID `sql:"player_id"`
	ChatId   uuid.UUID `sql:"chat_id"`
	Status   int8      `sql:"status"`
}

func GetMessagesByChatId(t *transport.Transport, playerId, chatId uuid.UUID) ([]Message, error) {
	messages := make([]Message, 0)
	sqlQuery := fmt.Sprintf("SELECT player_id, chat_id, time, message_id, role, content, is_archived, tool_calls, tool_call_id FROM messages WHERE player_id = Uuid('%s') AND chat_id = Uuid('%s') AND is_archived = false ORDER BY message_id", playerId, chatId)
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

func Write(t *transport.Transport, playerId, chatId uuid.UUID, message types.DeepSeekRoleContent) error {
	newMsgId := uuid.New()

	// Marshal tool_calls to JSON if present
	var toolCallsJSON string
	if len(message.ToolCalls) > 0 {
		toolCallsBytes, err := json.Marshal(message.ToolCalls)
		if err != nil {
			return fmt.Errorf("failed to marshal tool_calls: %w", err)
		}
		toolCallsJSON = string(toolCallsBytes)
		// Escape single quotes in JSON for YDB
		toolCallsJSON = strings.ReplaceAll(toolCallsJSON, "'", "''")
	}

	// Build INSERT query based on available fields
	if toolCallsJSON != "" && message.ToolCallId != "" {
		// Both tool_calls and tool_call_id are present
		// Escape single quotes in content and tool_call_id
		escapedContent := strings.ReplaceAll(message.Content, "'", "''")
		escapedToolCallId := strings.ReplaceAll(message.ToolCallId, "'", "''")
		sqlQuery := fmt.Sprintf("INSERT INTO messages (player_id, chat_id, message_id, time, role, content, tool_calls, tool_call_id) VALUES (Uuid('%s'), Uuid('%s'), Uuid('%s'), CurrentUtcDatetime(), '%s', '%s', '%s', '%s')",
			playerId.String(), chatId.String(), newMsgId.String(), message.Role, escapedContent, toolCallsJSON, escapedToolCallId)
		return t.YdbClient.Query().Exec(t.Ctx, sqlQuery)
	} else if toolCallsJSON != "" {
		// Only tool_calls is present
		escapedContent := strings.ReplaceAll(message.Content, "'", "''")
		sqlQuery := fmt.Sprintf("INSERT INTO messages (player_id, chat_id, message_id, time, role, content, tool_calls) VALUES (Uuid('%s'), Uuid('%s'), Uuid('%s'), CurrentUtcDatetime(), '%s', '%s', '%s')",
			playerId.String(), chatId.String(), newMsgId.String(), message.Role, escapedContent, toolCallsJSON)
		return t.YdbClient.Query().Exec(t.Ctx, sqlQuery)
	} else if message.ToolCallId != "" {
		// Only tool_call_id is present
		escapedContent := strings.ReplaceAll(message.Content, "'", "''")
		escapedToolCallId := strings.ReplaceAll(message.ToolCallId, "'", "''")
		sqlQuery := fmt.Sprintf("INSERT INTO messages (player_id, chat_id, message_id, time, role, content, tool_call_id) VALUES (Uuid('%s'), Uuid('%s'), Uuid('%s'), CurrentUtcDatetime(), '%s', '%s', '%s')",
			playerId.String(), chatId.String(), newMsgId.String(), message.Role, escapedContent, escapedToolCallId)
		return t.YdbClient.Query().Exec(t.Ctx, sqlQuery)
	} else {
		// Neither tool_calls nor tool_call_id is present
		escapedContent := strings.ReplaceAll(message.Content, "'", "''")
		sqlQuery := fmt.Sprintf("INSERT INTO messages (player_id, chat_id, message_id, time, role, content) VALUES (Uuid('%s'), Uuid('%s'), Uuid('%s'), CurrentUtcDatetime(), '%s', '%s')",
			playerId.String(), chatId.String(), newMsgId.String(), message.Role, escapedContent)
		return t.YdbClient.Query().Exec(t.Ctx, sqlQuery)
	}
}
