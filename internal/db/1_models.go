package db

import (
	"time"

	"github.com/google/uuid"
)

// Users
type UserRole int

const (
	RoleUser UserRole = iota
	RoleMod
	RoleAdmin
)

var (
	userRoleList = []string{
		"user",
		"moderator",
		"admin",
	}
)

func (r UserRole) String() string {
	return userRoleList[r]
}

type User struct {
	ID           uuid.UUID
	UserName     string
	DisplayName  string
	Email        string
	PasswordHash string
	ProfilePic   string
	Bio          string
	Followers    int
	Following    int
	IsActive     bool
	IsVerified   bool
	UserRole     string
	LastLogin    time.Time
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    *time.Time
}

// Follows
type FollowStatus int

const (
	StatusPending FollowStatus = iota
	StatusAccepted
	StatusBlocked
)

var (
	followStatusList = []string{
		"pending",
		"accepted",
		"blocked",
	}
)

func (f FollowStatus) String() string {
	return followStatusList[f]
}

type Follows struct {
	FollowerID   uuid.UUID
	FollowedID   uuid.UUID
	FollowStatus string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// Posts
type Post struct {
	ID        uuid.UUID
	AuthorID  uuid.UUID
	Title     string
	Content   string
	PostImage string
	Rating    int
	IsPublic  bool
	IsActive  bool
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
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
	ID        uuid.UUID
	PostID    uuid.UUID
	AuthorID  uuid.UUID
	Content   string
	Rating    int
	IsActive  bool
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}

// Comment Ratings
type CommentRatings struct {
	TargetID    uuid.UUID
	UserID      uuid.UUID
	RatingValue int
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// Notifications
type NotifType int

const (
	NotifFollowRequest NotifType = iota
	NotifFollowResponse
	NotifComment
	NotifPostMention
	NotifCommentMention
	NotifPostRating
	NotifCommentRating
)

var (
	notifTypeList = []string{
		"follow_request",
		"follow_response",
		"post_comment",
		"post_mention",
		"comment_mention",
		"post_rating",
		"comment_rating",
	}
)

func (n NotifType) String() string {
	return notifTypeList[n]
}

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
