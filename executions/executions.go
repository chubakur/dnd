package executions

import (
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/chubakur/dnd/transport"
	"github.com/google/uuid"
	"github.com/ydb-platform/ydb-go-sdk/v3/query"
)

const (
	StatusPending = "pending"
	StatusDone    = "done"
	StatusError   = "error"
)

// Execution tracks an async LLM job so the workflow can poll for its result.
type Execution struct {
	ExecutionId uuid.UUID `sql:"execution_id" json:"execution_id"`
	PlayerId    uuid.UUID `sql:"player_id" json:"player_id"`
	ChatId      uuid.UUID `sql:"chat_id" json:"chat_id"`
	Status      string    `sql:"status" json:"status"`
	Response    string    `sql:"response" json:"response"`
	Error       string    `sql:"error" json:"error"`
}

// Create inserts a new pending execution and returns its id.
func Create(t *transport.Transport, playerId, chatId uuid.UUID) (uuid.UUID, error) {
	id := uuid.New()
	sql := fmt.Sprintf(
		"INSERT INTO executions (execution_id, player_id, chat_id, status, response, error, create_time, update_time) "+
			"VALUES (Uuid('%s'), Uuid('%s'), Uuid('%s'), '%s', '', '', CurrentUtcDatetime(), CurrentUtcDatetime())",
		id.String(), playerId.String(), chatId.String(), StatusPending,
	)
	if err := t.YdbClient.Query().Exec(t.Ctx, sql); err != nil {
		return uuid.Nil, err
	}
	return id, nil
}

// SetDone marks the execution finished with the LLM response text.
func SetDone(t *transport.Transport, id uuid.UUID, response string) error {
	esc := strings.ReplaceAll(response, "'", "''")
	sql := fmt.Sprintf(
		"UPDATE executions SET status = '%s', response = '%s', update_time = CurrentUtcDatetime() WHERE execution_id = Uuid('%s')",
		StatusDone, esc, id.String(),
	)
	return t.YdbClient.Query().Exec(t.Ctx, sql)
}

// SetError marks the execution failed with an error message.
func SetError(t *transport.Transport, id uuid.UUID, errMsg string) error {
	esc := strings.ReplaceAll(errMsg, "'", "''")
	sql := fmt.Sprintf(
		"UPDATE executions SET status = '%s', error = '%s', update_time = CurrentUtcDatetime() WHERE execution_id = Uuid('%s')",
		StatusError, esc, id.String(),
	)
	return t.YdbClient.Query().Exec(t.Ctx, sql)
}

// Get returns the execution by id, or nil if not found.
func Get(t *transport.Transport, id uuid.UUID) (*Execution, error) {
	sql := fmt.Sprintf(
		"SELECT execution_id, player_id, chat_id, status, response, error FROM executions WHERE execution_id = Uuid('%s') LIMIT 1",
		id.String(),
	)
	res, err := t.YdbClient.Query().QueryResultSet(t.Ctx, sql, query.WithIdempotent())
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
	var e Execution
	if err := row.ScanStruct(&e); err != nil {
		return nil, err
	}
	return &e, nil
}
