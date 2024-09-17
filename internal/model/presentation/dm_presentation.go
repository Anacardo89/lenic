package presentation

import (
	"time"
)

type Conversation struct {
	Id        int        `json:"id"`
	User1     User       `json:"user1"`
	User2     User       `json:"user2"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	Messages  []DMessage `json:"messages"`
}

type DMessage struct {
	Id             int       `json:"id"`
	ConversationId int       `json:"conversation_id"`
	Sender         User      `json:"sender"`
	Content        string    `json:"content"`
	CreatedAt      time.Time `json:"created_at"`
}
