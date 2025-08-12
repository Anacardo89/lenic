package models

import (
	"encoding/base64"
	"fmt"

	"github.com/Anacardo89/lenic/internal/db"
)

func FromDBUser(u *db.User) *User {
	return &User{
		ID:           u.ID,
		UserName:     u.UserName,
		EncodedName:  base64.URLEncoding.EncodeToString([]byte(u.UserName)),
		DisplayName:  u.DisplayName,
		Email:        u.Email,
		ProfilePic:   u.ProfilePic,
		Bio:          u.Bio,
		Followers:    u.Followers,
		Following:    u.Following,
		PasswordHash: u.PasswordHash,
		IsActive:     u.IsActive,
	}
}

func FromDBUserNotif(u *db.User) *UserNotif {
	return &UserNotif{
		ID:          u.ID,
		UserName:    u.UserName,
		EncodedName: base64.URLEncoding.EncodeToString([]byte(u.UserName)),
		ProfilePic:  u.ProfilePic,
	}
}

func FromDBFollows(f *db.Follows) *Follows {
	return &Follows{
		FollowerID:   f.FollowerID,
		FollowedID:   f.FollowedID,
		FollowStatus: FollowStatus(f.FollowStatus),
	}
}

func ToDBUser(u *User) *db.User {
	return &db.User{
		UserName:     u.UserName,
		Email:        u.Email,
		PasswordHash: u.PasswordHash,
		ProfilePic:   u.ProfilePic,
		IsActive:     u.IsActive,
		IsVerified:   u.IsVerified,
	}
}

func FromDBPost(p *db.Post, u *User) *Post {
	return &Post{
		ID:         p.ID,
		Author:     u,
		Title:      p.Title,
		RawContent: p.Content,
		Image:      p.PostImage,
		Date:       fmt.Sprint(p.CreatedAt.Format(dateLayout)),
		IsPublic:   p.IsPublic,
		Rating:     p.Rating,
	}
}

func FromDBComment(c *db.Comment, u *User) *Comment {
	return &Comment{
		ID:      c.ID,
		Author:  *u,
		Content: c.Content,
		Date:    fmt.Sprint(c.CreatedAt.Format(dateLayout)),
		Rating:  c.Rating,
	}
}

func FromDBNotification(n *db.Notification, u, fromU UserNotif) *Notification {
	return &Notification{
		ID:         n.ID,
		User:       u,
		FromUser:   fromU,
		NotifType:  NotifType(n.NotifType),
		NotifText:  n.NotifText,
		ResourceID: n.ResourceID,
		ParentID:   n.ParentID,
		IsRead:     n.IsRead,
	}
}

func FromDBConversation(c *db.Conversation, u1, u2 UserNotif, isRead bool) *Conversation {
	return &Conversation{
		ID:        c.ID,
		User1:     u1,
		User2:     u2,
		IsRead:    isRead,
		CreatedAt: c.CreatedAt,
		UpdatedAt: c.UpdatedAt,
	}
}

func FromDBDMessage(m *db.DMessage, u UserNotif) *DMessage {
	return &DMessage{
		ID:             m.ID,
		ConversationID: m.ConversationID,
		Sender:         u,
		Content:        m.Content,
		IsRead:         m.IsRead,
		CreatedAt:      m.CreatedAt,
	}
}
