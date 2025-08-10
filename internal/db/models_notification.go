package db

import (
	"time"

	"github.com/google/uuid"
)

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
