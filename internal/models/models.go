package models

import (
	"html/template"
	"time"

	"github.com/google/uuid"
)

var (
	dateLayout = "2006-12-22 15:04:05"
)

// User
type UserRole string

const (
	RoleUser  UserRole = "user"
	RoleMod   UserRole = "moderator"
	RoleAdmin UserRole = "admin"
)

func (r UserRole) String() string {
	return string(r)
}

type User struct {
	ID           uuid.UUID `json:"id"`
	UserName     string    `json:"username"`
	EncodedName  string    `json:"encoded_name"`
	DisplayName  string    `json:"display_name"`
	Email        string    `json:"email"`
	Pass         string    `json:"pass"`
	ProfilePic   string    `json:"profile_pic"`
	Bio          string    `json:"bio"`
	Followers    int       `json:"followers"`
	Following    int       `json:"followiong"`
	PasswordHash string    `json:"password_hash"`
	IsActive     bool      `json:"is_active"`
	IsVerified   bool      `json:"is_verified"`
	UserRole     UserRole  `json:"user_role"`
}

type UserNotif struct {
	ID          uuid.UUID `json:"id"`
	UserName    string    `json:"username"`
	EncodedName string    `json:"encoded_name"`
	ProfilePic  string    `json:"profile_pic"`
}

// Follow
type FollowStatus string

const (
	StatusPending  FollowStatus = "pending"
	StatusAccepted FollowStatus = "accepted"
	StatusBlocked  FollowStatus = "blocked"
)

func (f FollowStatus) String() string {
	return string(f)
}

type Follows struct {
	FollowerID   uuid.UUID    `json:"follower_id"`
	FollowedID   uuid.UUID    `json:"followed_id"`
	FollowStatus FollowStatus `json:"follow_status"`
}

// Post
type Post struct {
	ID         uuid.UUID     `json:"id"`
	Author     User          `json:"author"`
	Title      string        `json:"title"`
	RawContent string        `json:"raw_content"`
	Content    template.HTML `json:"content"`
	Image      string        `json:"img"`
	Rating     int           `json:"rating"`
	UserRating int           `json:"user_ratinmg"`
	Date       string        `json:"date"`
	IsPublic   bool          `json:"is_public"`
	Comments   []Comment     `json:"comments"`
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
	ID         uuid.UUID `json:"id"`
	Author     User      `json:"author"`
	Content    string    `json:"content"`
	Date       string    `json:"date"`
	Rating     int       `json:"rating"`
	UserRating int       `json:"user_rating"`
}

// Notification
type NotifType string

const (
	NotifFollowRequest  NotifType = "follow_request"
	NotifFollowResponse NotifType = "follow_response"
	NotifComment        NotifType = "post_comment"
	NotifPostMention    NotifType = "post_mention"
	NotifCommentMention NotifType = "comment_mention"
	NotifPostRating     NotifType = "post_rating"
	NotifCommentRating  NotifType = "comment_rating"
)

func (n NotifType) String() string {
	return string(n)
}

type Notification struct {
	ID         uuid.UUID `json:"id"`
	User       UserNotif `json:"user"`
	FromUser   UserNotif `json:"from_user"`
	NotifType  NotifType `json:"notif_type"`
	NotifText  string    `json:"notif_text"`
	ResourceID string    `json:"resource_id"`
	ParentID   string    `json:"parent_id"`
	IsRead     bool      `json:"is_read"`
}

// Conversation
type Conversation struct {
	ID        uuid.UUID `json:"id"`
	User1     UserNotif `json:"user1"`
	User2     UserNotif `json:"user2"`
	IsRead    bool      `json:"is_read"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// DM
type DMessage struct {
	ID             uuid.UUID `json:"id"`
	ConversationID uuid.UUID `json:"conversation_id"`
	Sender         UserNotif `json:"sender"`
	Content        string    `json:"content"`
	IsRead         bool      `json:"is_read"`
	CreatedAt      time.Time `json:"created_at"`
}
