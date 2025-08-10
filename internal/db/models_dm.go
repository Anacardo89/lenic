package db

import (
	"time"

	"github.com/google/uuid"
)

type Conversation struct {
	ID        uuid.UUID
	User1ID   uuid.UUID
	User2ID   uuid.UUID
	CreatedAt time.Time
	UpdatedAt time.Time
}

type DMessage struct {
	ID             uuid.UUID
	ConversationID uuid.UUID
	SenderID       uuid.UUID
	Content        string
	IsRead         bool
	CreatedAt      time.Time
}
