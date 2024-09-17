package mapper

import (
	"encoding/base64"
	"fmt"

	"github.com/Anacardo89/tpsi25_blog/internal/model/database"
	"github.com/Anacardo89/tpsi25_blog/internal/model/presentation"
	"github.com/Anacardo89/tpsi25_blog/pkg/db"
)

func User(u *database.User) *presentation.User {
	return &presentation.User{
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

func UserNotif(u *database.User) *presentation.UserNotif {
	return &presentation.UserNotif{
		Id:          u.Id,
		UserName:    u.UserName,
		EncodedName: base64.URLEncoding.EncodeToString([]byte(u.UserName)),
		ProfilePic:  u.ProfilePic,
	}
}

func Follows(f *database.Follows) *presentation.Follows {
	return &presentation.Follows{
		FollowerId: f.FollowerId,
		FollowedId: f.FollowedId,
		Status:     f.Status,
	}
}

func UserToDB(u *presentation.User) *database.User {
	return &database.User{
		UserName:   u.UserName,
		Email:      u.Email,
		HashPass:   u.HashPass,
		ProfilePic: u.ProfilePic,
		Active:     u.Active,
	}
}

func Session(s *database.Session) *presentation.Session {
	return &presentation.Session{
		Id:            s.Id,
		SessionId:     s.SessionId,
		Authenticated: true,
	}
}

func Post(p *database.Post, u *presentation.User) *presentation.Post {

	return &presentation.Post{
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

func Comment(c *database.Comment, u *presentation.User) *presentation.Comment {
	return &presentation.Comment{
		Id:      c.Id,
		Author:  *u,
		Content: c.Content,
		Date:    fmt.Sprint(c.CreatedAt.Format(db.DateLayout)),
		Rating:  c.Rating,
	}
}

func Notification(n *database.Notification, u presentation.UserNotif, from_u presentation.UserNotif) *presentation.Notification {
	return &presentation.Notification{
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

func Convesation(c *database.Conversation, u1 presentation.User, u2 presentation.User, dms []presentation.DMessage) *presentation.Conversation {
	return &presentation.Conversation{
		Id:        c.Id,
		User1:     u1,
		User2:     u2,
		CreatedAt: c.CreatedAt,
		UpdatedAt: c.UpdatedAt,
		Messages:  dms,
	}
}

func DMessage(m *database.DMessage, u presentation.User) *presentation.DMessage {
	return &presentation.DMessage{
		Id:             m.Id,
		ConversationId: m.ConversationId,
		Sender:         u,
		Content:        m.Content,
		CreatedAt:      m.CreatedAt,
	}
}
