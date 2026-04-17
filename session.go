package main

import (
	"fmt"

	"github.com/google/uuid"
)

type PlayerSession struct {
	PlayerId  uuid.UUID `sql:"player_id"`
	SessionId uuid.UUID `sql:"session_id"`
}

func GetActivePlayerSessions(t *transport, playerId uuid.UUID) ([]PlayerSession, error) {
	var result []PlayerSession = make([]PlayerSession, 0)
	query := fmt.Sprintf("SELECT player_id, session_id FROM sessions WHERE player_id = Uuid('%s')", playerId.String())
	fmt.Println(query)
	res, err := t.ydbClient.Query().QueryResultSet(t.ctx, query)
	if err != nil {
		return result, err
	}

	defer res.Close(t.ctx)
	for row, err := range res.Rows(t.ctx) {
		if err != nil {
			return result, err
		}
		var o PlayerSession
		err = row.ScanStruct(&o)
		if err != nil {
			return result, err
		}
		result = append(result, o)
	}

	return result, nil
}
