package sessions

import (
	"time"

	"github.com/google/uuid"
)

type Session struct {
	PlayerId   uuid.UUID `sql:"player_id" json:"player_id"`
	SessionId  uuid.UUID `sql:"session_id" json:"session_id"`
	State      int       `sql:"state" json:"state"`
	Context    string    `sql:"context" json:"context"`
	CreateTime time.Time `sql:"create_time" json:"create_time"`
	UpdateTime time.Time `sql:"update_time" json:"update_time"`
}

func getActive(playerId uuid.UUID) (*Session, error) {
	return nil, nil
}

func makeActive(playerId, sessionId uuid.UUID) error {
	return nil
}
