package dndcore

import "github.com/google/uuid"

type GameContext struct {
	PlayerId uuid.UUID
	ChatId   uuid.UUID
}
