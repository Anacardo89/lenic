package model

import "time"

type Comment struct {
	Id            int
	PostGUID      string
	CommentAuthor string
	CommentText   string
	CreatedAt     time.Time
	UpdatedAt     time.Time
}
