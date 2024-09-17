package presentation

import (
	"time"
)

type Conversation struct {
	Id        int
	User1     User
	User2     User
	CreatedAt time.Time
	UpdatedAt time.Time
	Messages  []DMessage
}

type DMessage struct {
	Id             int
	ConversationId int
	Sender         User
	Content        string
	CreatedAt      time.Time
}
