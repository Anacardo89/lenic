package db

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type DBRepository interface {
	// Convenience Methods
	Pool() *pgxpool.Pool
	Close()

	// User
	CreateUser(ctx context.Context, u *User) (uuid.UUID, error)
	GetUserByID(ctx context.Context, ID uuid.UUID) (*User, error)
	GetUserByUserName(ctx context.Context, userName string) (*User, error)
	GetUserByEmail(ctx context.Context, email string) (*User, error)
	SearchUsersByDisplayName(ctx context.Context, displayName string) ([]*User, error)
	SetUserActive(ctx context.Context, userName string) error
	SetNewPassword(ctx context.Context, userName string, pass string) error
	UpdateProfilePic(ctx context.Context, userName string, profilePic string) error

	// Follows
	FollowUser(ctx context.Context, followerID int, followedID int) error
	AcceptFollow(ctx context.Context, followerID int, followedID int) error
	UnfollowUser(ctx context.Context, followerID int, followedID int) error
	GetUserFollows(ctx context.Context, followerID int, followedID int) (*Follows, error)
	GetFollowers(ctx context.Context, followedID int) ([]*Follows, error)
	GetFollowing(ctx context.Context, followerID int) ([]*Follows, error)

	// Posts
	CreatePost(ctx context.Context, p *Post) (uuid.UUID, error)
	GetFeed(ctx context.Context, userID uuid.UUID) ([]*Post, error)
	GetUserPosts(ctx context.Context, userID uuid.UUID) ([]*Post, error)
	GetUserPublicPosts(ctx context.Context, userID uuid.UUID) ([]*Post, error)
	GetPostByID(ctx context.Context, ID uuid.UUID) (*Post, error)
	UpdatePost(ctx context.Context, post *Post) error
	DisablePost(ctx context.Context, ID uuid.UUID) error

	// Post Ratings
	RatePostUp(ctx context.Context, targetID uuid.UUID, userID uuid.UUID) error
	RatePostDown(ctx context.Context, targetID uuid.UUID, userID uuid.UUID) error
	GetPostUserRating(ctx context.Context, targetID uuid.UUID, userID uuid.UUID) (*PostRatings, error)

	// Comments
	CreateComment(ctx context.Context, comment *Comment) (uuid.UUID, error)
	GetCommentById(ctx context.Context, ID uuid.UUID) (*Comment, error)
	GetCommentsByPost(ctx context.Context, postID uuid.UUID) ([]*Comment, error)
	UpdateCommentContent(ctx context.Context, ID uuid.UUID, content string) error
	DisableComment(ctx context.Context, ID uuid.UUID) error

	// CommentRatings
	RateCommentUp(ctx context.Context, targetID uuid.UUID, userID uuid.UUID) error
	RateCommentDown(ctx context.Context, targetID uuid.UUID, userID uuid.UUID) error
	GetCommentUserRating(ctx context.Context, targetID uuid.UUID, userID uuid.UUID) (*CommentRatings, error)

	// Notifications
	CreateNotification(ctx context.Context, n *Notification) (uuid.UUID, error)
	GetFollowNotification(ctx context.Context, userID uuid.UUID, fromUserID uuid.UUID) (*Notification, error)
	GetNotificationById(ctx context.Context, ID uuid.UUID) (*Notification, error)
	GetNotificationsByUser(ctx context.Context, userID uuid.UUID, limit int, offset int) ([]*Notification, error)
	UpdateNotificationRead(ctx context.Context, ID uuid.UUID) error
	DeleteNotificationByID(ctx context.Context, ID uuid.UUID) error
}

func (c *dbClient) Pool() *pgxpool.Pool {
	return c.pool
}

func (c *dbClient) Close() {
	c.pool.Close()
}
