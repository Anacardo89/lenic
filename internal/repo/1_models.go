package repo

import (
	"time"

	"github.com/google/uuid"
)

// Users

type User struct {
	ID           uuid.UUID  `json:"id"`
	Username     string     `json:"username"`
	DisplayName  string     `json:"display_name"`
	Email        string     `json:"email"`
	PasswordHash string     `json:"password_hash"`
	ProfilePic   string     `json:"profile_pic"`
	Bio          string     `json:"bio"`
	Followers    int        `json:"user_followers"`
	Following    int        `json:"user_following"`
	IsActive     bool       `json:"is_active"`
	IsVerified   bool       `json:"is_verified"`
	UserRole     string     `json:"user_role"`
	LastLogin    time.Time  `json:"last_login"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
	DeletedAt    *time.Time `json:"deleted_at"`
}

// Follows
type Follows struct {
	FollowerID   uuid.UUID
	FollowedID   uuid.UUID
	FollowStatus string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// Posts
type Post struct {
	ID        uuid.UUID  `json:"id"`
	AuthorID  uuid.UUID  `json:"author_id"`
	Title     string     `json:"title"`
	Content   string     `json:"content"`
	PostImage string     `json:"post_image"`
	Rating    int        `json:"rating"`
	IsPublic  bool       `json:"is_public"`
	IsActive  bool       `json:"is_active"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
}

// Post Ratings
type PostRatings struct {
	TargetID    uuid.UUID
	UserID      uuid.UUID
	RatingValue int
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// Comments
type Comment struct {
	ID        uuid.UUID  `json:"id"`
	PostID    uuid.UUID  `json:"post_id"`
	AuthorID  uuid.UUID  `json:"author_id"`
	Content   string     `json:"content"`
	Rating    int        `json:"rating"`
	IsActive  bool       `json:"is_active"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
}

// Comment Ratings
type CommentRatings struct {
	TargetID    uuid.UUID
	UserID      uuid.UUID
	RatingValue int
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type PostWithComments struct {
	Post
	UserRating int
	Author     User                 `json:"author"`
	Comments   []CommentsWithAuthor `json:"comments"`
}

type CommentsWithAuthor struct {
	Comment
	UserRating int  `json:"user_rating"`
	Author     User `json:"author"`
}

// Notifications

type Notification struct {
	ID         uuid.UUID
	UserID     uuid.UUID
	FromUserID uuid.UUID
	NotifType  string
	NotifText  string
	ResourceID uuid.UUID
	ParentID   uuid.UUID
	IsRead     bool
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

type NotificationWithUsers struct {
	Notification Notification
	User         User
	FromUser     User
}

// Tags
type ResourceType int

const (
	ResourcePost ResourceType = iota
	ResourceComment
)

var (
	resourceTypeList = []string{
		"post",
		"comment",
	}
)

func (r ResourceType) String() string {
	return resourceTypeList[r]
}

type UserTag struct {
	UserID      uuid.UUID
	TargetID    uuid.UUID
	ResourceTpe string
}

type HashTag struct {
	ID        uuid.UUID
	TagName   string
	CreatedAt time.Time
}

type HashTagResource struct {
	TagID       uuid.UUID
	TargetID    uuid.UUID
	ResourceTpe string
}

// Conversations
type Conversation struct {
	ID        uuid.UUID
	User1ID   uuid.UUID
	User2ID   uuid.UUID
	CreatedAt time.Time
	UpdatedAt time.Time
}

// DMs
type DMessage struct {
	ID             uuid.UUID
	ConversationID uuid.UUID
	SenderID       uuid.UUID
	Content        string
	IsRead         bool
	CreatedAt      time.Time
}

type ConversationsWithDMs struct {
	ID        uuid.UUID
	CreatedAt time.Time
	OtherUser *User
	Messages  []*DMessage
}

type DMessageWithUser struct {
	ID             uuid.UUID
	ConversationID uuid.UUID
	Sender         *User
	Content        string
	IsRead         bool
	CreatedAt      time.Time
}
