package presentation

import (
	"time"
)

type Conversation struct {
	Id        int       `json:"id"`
	User1     UserNotif `json:"user1"`
	User2     UserNotif `json:"user2"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type DMessage struct {
	Id             int       `json:"id"`
	ConversationId int       `json:"conversation_id"`
	Sender         UserNotif `json:"sender"`
	Content        string    `json:"content"`
	CreatedAt      time.Time `json:"created_at"`
}
