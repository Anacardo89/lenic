package models

import (
	"html/template"
	"time"
)

// User
type User struct {
	Id          int    `json:"id"`
	UserName    string `json:"username"`
	EncodedName string
	Email       string
	Pass        string
	ProfilePic  string
	Followers   int
	Following   int
	HashPass    string
	Active      int
}

type UserNotif struct {
	Id          int    `json:"id"`
	UserName    string `json:"username"`
	EncodedName string `json:"encoded"`
	ProfilePic  string `json:"profile_pic"`
}

type Follows struct {
	FollowerId int
	FollowedId int
	Status     int
}

// Post
type Post struct {
	Id         int
	GUID       string
	Author     User
	Title      string
	RawContent string
	Content    template.HTML
	Image      string
	Date       string
	IsPublic   bool
	Rating     int
	UserRating int
	Comments   []Comment
}

func (p Post) TruncatedText() string {
	chars := 0
	for i := range p.RawContent {
		chars++
		if chars > 150 {
			return p.RawContent[:i] + `...`
		}
	}
	return p.RawContent
}

// Comment
type Comment struct {
	Id         int
	Author     User
	Content    string
	Date       string
	Rating     int
	UserRating int
}

// Notification
type Notification struct {
	Id         int       `json:"id"`
	User       UserNotif `json:"user"`
	FromUser   UserNotif `json:"fromuser"`
	NotifType  string    `json:"type"`
	NotifMsg   string    `json:"msg"`
	ResourceId string    `json:"resource_id"`
	ParentId   string    `json:"parent_id"`
	IsRead     bool      `json:"is_read"`
}

// Conversation
type Conversation struct {
	Id        int       `json:"id"`
	User1     UserNotif `json:"user1"`
	User2     UserNotif `json:"user2"`
	IsRead    bool      `json:"is_read"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// DM
type DMessage struct {
	Id             int       `json:"id"`
	ConversationId int       `json:"conversation_id"`
	Sender         UserNotif `json:"sender"`
	Content        string    `json:"content"`
	IsRead         bool      `json:"is_read"`
	CreatedAt      time.Time `json:"created_at"`
}

// Session
type Session struct {
	Authenticated bool
	User          User
	Notifs        []Notification
	DMs           []Conversation
}
