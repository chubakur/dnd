package dndcore

import "github.com/google/uuid"

type LlmAgent interface {
	// SendMessage adds a user message and runs the agent executor loop.
	// Returns the final assistant text response.
	SendMessage(playerId, chatId uuid.UUID, message string) (string, error)
	// GetHistory returns stored messages for a chat as a flat string slice (role: content).
	GetHistory(playerId, chatId uuid.UUID) ([]string, error)
	// ResetContext clears the in-memory message chain (does not delete DB records).
	ResetContext()
}
