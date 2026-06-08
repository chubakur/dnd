package sessions

import (
	"errors"
	"fmt"
	"io"
	"time"

	"github.com/chubakur/dnd/transport"
	"github.com/google/uuid"
	"github.com/ydb-platform/ydb-go-sdk/v3/query"
)

type Session struct {
	PlayerId   uuid.UUID `sql:"player_id" json:"player_id"`
	SessionId  uuid.UUID `sql:"session_id" json:"session_id"`
	State      int       `sql:"state" json:"state"`
	Context    string    `sql:"context" json:"context"`
	CreateTime time.Time `sql:"create_time" json:"create_time"`
	UpdateTime time.Time `sql:"update_time" json:"update_time"`
}

func GetActive(t *transport.Transport, playerId uuid.UUID) (*Session, error) {
	querySql := fmt.Sprintf(
		"SELECT player_id, session_id, state, context, create_time, update_time FROM sessions WHERE player_id = Uuid('%s') AND state = 1 ORDER BY update_time DESC LIMIT 1",
		playerId.String(),
	)
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
	var s Session
	err = row.ScanStruct(&s)
	if err != nil {
		return nil, err
	}
	return &s, nil
}

func MakeActive(t *transport.Transport, playerId, sessionId uuid.UUID) error {
	// Deactivate all other sessions for this player first
	deactivateSql := fmt.Sprintf(
		"UPDATE sessions SET state = 0, update_time = CurrentUtcDatetime() WHERE player_id = Uuid('%s') AND state = 1",
		playerId.String(),
	)
	err := t.YdbClient.Query().Exec(t.Ctx, deactivateSql)
	if err != nil {
		return err
	}
	activateSql := fmt.Sprintf(
		"UPDATE sessions SET state = 1, update_time = CurrentUtcDatetime() WHERE player_id = Uuid('%s') AND session_id = Uuid('%s')",
		playerId.String(), sessionId.String(),
	)
	return t.YdbClient.Query().Exec(t.Ctx, activateSql)
}

func CreateSession(t *transport.Transport, playerId uuid.UUID, worldContext string) (*Session, error) {
	sessionId := uuid.New()
	insertSql := fmt.Sprintf(
		"INSERT INTO sessions (player_id, session_id, state, context, create_time, update_time) VALUES (Uuid('%s'), Uuid('%s'), 1, '%s', CurrentUtcDatetime(), CurrentUtcDatetime())",
		playerId.String(), sessionId.String(), worldContext,
	)
	err := t.YdbClient.Query().Exec(t.Ctx, insertSql)
	if err != nil {
		return nil, err
	}
	now := time.Now().UTC()
	return &Session{
		PlayerId:   playerId,
		SessionId:  sessionId,
		State:      1,
		Context:    worldContext,
		CreateTime: now,
		UpdateTime: now,
	}, nil
}
