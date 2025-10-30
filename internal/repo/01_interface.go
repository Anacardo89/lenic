package repo

import (
	"context"

	"github.com/google/uuid"
)

type DBRepo interface {
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
	SetNewPassword(ctx context.Context, userID uuid.UUID, passHash string) error
	UpdateProfilePic(ctx context.Context, userName string, profilePic string) error

	// Follows
	FollowUser(ctx context.Context, followerID uuid.UUID, followedUsername string) error
	AcceptFollow(ctx context.Context, followerName, followedName string) error
	UnfollowUser(ctx context.Context, followerName, followedName string) error
	GetUserFollows(ctx context.Context, followerID, followedID uuid.UUID) (*Follows, error)
	GetFollowers(ctx context.Context, followedID uuid.UUID) ([]*Follows, error)
	GetFollowing(ctx context.Context, followerID uuid.UUID) ([]*Follows, error)

	// Posts
	CreatePost(ctx context.Context, p *Post) (uuid.UUID, error)
	GetFeed(ctx context.Context, username string) ([]*Post, error)
	GetUserPosts(ctx context.Context, userID uuid.UUID) ([]*Post, error)
	GetPostAuthorFromComment(ctx context.Context, commentID uuid.UUID) (*User, error)
	GetUserPublicPosts(ctx context.Context, userID uuid.UUID) ([]*Post, error)
	GetPost(ctx context.Context, ID uuid.UUID) (*Post, error)
	GetPostForPage(ctx context.Context, ID, userID uuid.UUID) (*PostWithComments, error)
	UpdatePost(ctx context.Context, post *Post) error
	DisablePost(ctx context.Context, ID uuid.UUID) (*Post, error)

	// Post Ratings
	RatePostUp(ctx context.Context, targetID, userID uuid.UUID) error
	RatePostDown(ctx context.Context, targetID, userID uuid.UUID) error

	// Comments
	CreateComment(ctx context.Context, comment *Comment) error
	GetComment(ctx context.Context, ID uuid.UUID) (*Comment, error)
	UpdateComment(ctx context.Context, comment *Comment) error
	DisableComment(ctx context.Context, ID uuid.UUID) (*Comment, error)

	// CommentRatings
	RateCommentUp(ctx context.Context, targetID, userID uuid.UUID) error
	RateCommentDown(ctx context.Context, targetID, userID uuid.UUID) error

	// Notifications
	CreateNotification(ctx context.Context, n *Notification) error
	DeleteFollowNotification(ctx context.Context, username, fromUsername string) error
	GetUserNotifs(ctx context.Context, username string, limit, offset int) ([]*NotificationWithUsers, error)
	UpdateNotificationRead(ctx context.Context, ID uuid.UUID) error

	// UserTags
	CreateUserTag(ctx context.Context, t *UserTag) error
	DeleteUserTag(ctx context.Context, userID uuid.UUID, targetID uuid.UUID) error

	//HashTags
	CreateHashtag(ctx context.Context, t *HashTag) (uuid.UUID, error)
	GetHashTagByName(ctx context.Context, tagName string) (*HashTag, error)

	// HashTag Resources
	CreateHashTagResource(ctx context.Context, t *HashTagResource) error
	GetHashTagResourceByTarget(ctx context.Context, tagID, targetID uuid.UUID) (*HashTagResource, error)

	// Conversations
	GetConversationAndSender(ctx context.Context, conversationID uuid.UUID, username string) (*Conversation, *User, error)
	GetConversationAndUsers(ctx context.Context, user1, user2 string) (*Conversation, []*User, error)
	GetConversationsAndOwner(ctx context.Context, user string, limit, offset int) (*User, []*ConversationsWithDMs, error)
	GetConversationByUsers(ctx context.Context, user1ID, user2ID uuid.UUID) (*Conversation, error)
	UpdateConversation(ctx context.Context, ID uuid.UUID) error

	// DMs
	CreateDM(ctx context.Context, dm *DMessage) (uuid.UUID, error)
	GetDMsByConversation(ctx context.Context, conersationID uuid.UUID, limit, offset int) ([]*DMessageWithUser, error)
	ReadAllReceivedDMsInConvo(ctx context.Context, conversationID uuid.UUID, username string) error
}
