package presentation

import (
	"time"
)

type Conversation struct {
	Id        int       `json:"id"`
	User1     UserNotif `json:"user1"`
	User2     UserNotif `json:"user2"`
	IsRead    bool      `json:"is_read"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type DMessage struct {
	Id             int       `json:"id"`
	ConversationId int       `json:"conversation_id"`
	Sender         UserNotif `json:"sender"`
	Content        string    `json:"content"`
	IsRead         bool      `json:"is_read"`
	CreatedAt      time.Time `json:"created_at"`
}
