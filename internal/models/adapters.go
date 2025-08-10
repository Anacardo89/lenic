package models

import (
	"encoding/base64"
	"fmt"

	"github.com/Anacardo89/lenic/internal/db"
)

func FromDBUser(u *db.User) *User {
	return &User{
		Id:          u.Id,
		UserName:    u.UserName,
		EncodedName: base64.URLEncoding.EncodeToString([]byte(u.UserName)),
		Email:       u.Email,
		HashPass:    u.HashPass,
		ProfilePic:  u.ProfilePic,
		Followers:   u.Followers,
		Following:   u.Following,
		Active:      u.Active,
	}
}

func FromDBUserNotif(u *db.User) *UserNotif {
	return &UserNotif{
		Id:          u.Id,
		UserName:    u.UserName,
		EncodedName: base64.URLEncoding.EncodeToString([]byte(u.UserName)),
		ProfilePic:  u.ProfilePic,
	}
}

func FromDBFollows(f *db.Follows) *Follows {
	return &Follows{
		FollowerId: f.FollowerId,
		FollowedId: f.FollowedId,
		Status:     f.Status,
	}
}

func ToDBUser(u *User) *db.User {
	return &db.User{
		UserName:   u.UserName,
		Email:      u.Email,
		HashPass:   u.HashPass,
		ProfilePic: u.ProfilePic,
		Active:     u.Active,
	}
}

func FromDBPost(p *db.Post, u *User) *Post {
	return &Post{
		Id:         p.Id,
		GUID:       p.GUID,
		Author:     *u,
		Title:      p.Title,
		RawContent: p.Content,
		Image:      p.Image,
		Date:       fmt.Sprint(p.CreatedAt.Format(db.DateLayout)),
		IsPublic:   p.IsPublic,
		Rating:     p.Rating,
	}
}

func FromDBComment(c *db.Comment, u *User) *Comment {
	return &Comment{
		Id:      c.Id,
		Author:  *u,
		Content: c.Content,
		Date:    fmt.Sprint(c.CreatedAt.Format(db.DateLayout)),
		Rating:  c.Rating,
	}
}

func FromDBNotification(n *db.Notification, u, from_u UserNotif) *Notification {
	return &Notification{
		Id:         n.Id,
		User:       u,
		FromUser:   from_u,
		NotifType:  n.NotifType,
		NotifMsg:   n.NotifMsg,
		ResourceId: n.ResourceId,
		ParentId:   n.ParentId,
		IsRead:     n.IsRead,
	}
}

func FromDBConversation(c *db.Conversation, u1, u2 UserNotif, is_read bool) *Conversation {
	return &Conversation{
		Id:        c.Id,
		User1:     u1,
		User2:     u2,
		IsRead:    is_read,
		CreatedAt: c.CreatedAt,
		UpdatedAt: c.UpdatedAt,
	}
}

func FromDBDMessage(m *db.DMessage, u UserNotif) *DMessage {
	return &DMessage{
		Id:             m.Id,
		ConversationId: m.ConversationId,
		Sender:         u,
		Content:        m.Content,
		IsRead:         m.IsRead,
		CreatedAt:      m.CreatedAt,
	}
}
