package mapper

import (
	"fmt"

	"github.com/Anacardo89/tpsi25_blog/internal/model/database"
	"github.com/Anacardo89/tpsi25_blog/internal/model/presentation"
	"github.com/Anacardo89/tpsi25_blog/pkg/db"
)

func User(u *database.User) *presentation.User {
	return &presentation.User{
		Id:         u.Id,
		UserName:   u.UserName,
		UserEmail:  u.UserEmail,
		HashedPass: u.UserPass,
		Active:     u.Active,
	}
}

func UserToDB(u *presentation.User) *database.User {
	return &database.User{
		UserName:  u.UserName,
		UserEmail: u.UserEmail,
		UserPass:  u.HashedPass,
		Active:    u.Active,
	}
}

func Session(s *database.Session) *presentation.Session {
	return &presentation.Session{
		Id:            s.Id,
		SessionId:     s.SessionId,
		Authenticated: true,
	}
}

func Post(p *database.Post) *presentation.Post {
	return &presentation.Post{
		Id:         p.Id,
		GUID:       p.GUID,
		User:       p.User,
		Title:      p.Title,
		RawContent: p.Content,
		Image:      p.Image,
		Date:       fmt.Sprint(p.CreatedAt.Format(db.DateLayout)),
	}
}

func Comment(c *database.Comment) *presentation.Comment {
	return &presentation.Comment{
		Id:          c.Id,
		Author:      c.CommentAuthor,
		CommentText: c.CommentText,
		Date:        fmt.Sprint(c.CreatedAt.Format(db.DateLayout)),
	}
}
