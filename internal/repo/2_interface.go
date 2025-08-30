package repo

import (
	"context"

	"github.com/google/uuid"
)

type DBRepository interface {
	// Convenience Methods
	Close()

	// User
	CreateUser(ctx context.Context, u *User) (uuid.UUID, error)
	GetUserByID(ctx context.Context, ID uuid.UUID) (*User, error)
	GetUserByUserName(ctx context.Context, userName string) (*User, error)
	GetUserByEmail(ctx context.Context, email string) (*User, error)
	SearchUsersByUserName(ctx context.Context, username string) ([]*User, error)
	SearchUsersByDisplayName(ctx context.Context, displayName string) ([]*User, error)
	SetUserActive(ctx context.Context, userName string) error
	SetNewPassword(ctx context.Context, userName, passHash string) error
	UpdateProfilePic(ctx context.Context, userName string, profilePic string) error

	// Follows
	FollowUser(ctx context.Context, followerID, followedID uuid.UUID) error
	AcceptFollow(ctx context.Context, followerID, followedID uuid.UUID) error
	UnfollowUser(ctx context.Context, followerID, followedID uuid.UUID) error
	GetUserFollows(ctx context.Context, followerID, followedID uuid.UUID) (*Follows, error)
	GetFollowers(ctx context.Context, followedID uuid.UUID) ([]*Follows, error)
	GetFollowing(ctx context.Context, followerID uuid.UUID) ([]*Follows, error)

	// Posts
	CreatePost(ctx context.Context, p *Post) (uuid.UUID, error)
	GetFeed(ctx context.Context, userID uuid.UUID) ([]*Post, error)
	GetUserPosts(ctx context.Context, userID uuid.UUID) ([]*Post, error)
	GetUserPublicPosts(ctx context.Context, userID uuid.UUID) ([]*Post, error)
	GetPost(ctx context.Context, ID uuid.UUID) (*Post, error)
	UpdatePost(ctx context.Context, post *Post) error
	DisablePost(ctx context.Context, ID uuid.UUID) error

	// Post Ratings
	RatePostUp(ctx context.Context, targetID, userID uuid.UUID) error
	RatePostDown(ctx context.Context, targetID, userID uuid.UUID) error
	GetPostUserRating(ctx context.Context, targetID, userID uuid.UUID) (*PostRatings, error)

	// Comments
	CreateComment(ctx context.Context, comment *Comment) (uuid.UUID, error)
	GetComment(ctx context.Context, ID uuid.UUID) (*Comment, error)
	GetCommentsByPost(ctx context.Context, postID uuid.UUID) ([]*Comment, error)
	UpdateCommentContent(ctx context.Context, ID uuid.UUID, content string) error
	DisableComment(ctx context.Context, ID uuid.UUID) error

	// CommentRatings
	RateCommentUp(ctx context.Context, targetID, userID uuid.UUID) error
	RateCommentDown(ctx context.Context, targetID, userID uuid.UUID) error
	GetCommentUserRating(ctx context.Context, targetID, userID uuid.UUID) (*CommentRatings, error)

	// Notifications
	CreateNotification(ctx context.Context, n *Notification) (uuid.UUID, error)
	GetFollowNotification(ctx context.Context, userID, fromUserID uuid.UUID) (*Notification, error)
	GetNotification(ctx context.Context, ID uuid.UUID) (*Notification, error)
	GetNotificationsByUser(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*Notification, error)
	UpdateNotificationRead(ctx context.Context, ID uuid.UUID) error
	DeleteNotification(ctx context.Context, ID uuid.UUID) error

	// UserTags
	CreateUserTag(ctx context.Context, t *UserTag) error
	GetUserTagByTarget(ctx context.Context, userID, targetID uuid.UUID) (*UserTag, error)
	DeleteUserTag(ctx context.Context, userID uuid.UUID, targetID uuid.UUID) error

	//HashTags
	CreateHashtag(ctx context.Context, t *HashTag) (uuid.UUID, error)
	GetHashTagByName(ctx context.Context, tagName string) (*HashTag, error)

	// HashTag Resources
	CreateHashTagResource(ctx context.Context, t *HashTagResource) error
	GetHashTagResourceByTarget(ctx context.Context, tagID, targetID uuid.UUID) (*HashTagResource, error)

	// Conversations
	CreateConversation(ctx context.Context, conv *Conversation) (uuid.UUID, error)
	GetConversation(ctx context.Context, ID uuid.UUID) (*Conversation, error)
	GetConversationByUsers(ctx context.Context, user1ID, user2ID uuid.UUID) (*Conversation, error)
	GetConversationsByUser(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*Conversation, error)
	UpdateConversation(ctx context.Context, ID uuid.UUID) error

	// DMs
	CreateDM(ctx context.Context, dm *DMessage) (uuid.UUID, error)
	GetDM(ctx context.Context, ID uuid.UUID) (*DMessage, error)
	GetConvoLastDMBySender(ctx context.Context, conversationID, senderID uuid.UUID) (*DMessage, error)
	GetDMsByConversation(ctx context.Context, conersationID uuid.UUID, limit, offset int) ([]*DMessage, error)
	UpdateDMRead(ctx context.Context, ID uuid.UUID) error
}
