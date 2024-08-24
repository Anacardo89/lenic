package database

import "time"

type Post struct {
	Id        int
	GUID      string
	AuthorId  int
	Title     string
	Content   string
	Image     string
	ImageExt  string
	CreatedAt time.Time
	UpdatedAt time.Time
	IsPublic  int
	VoteCount int
	Active    int
}

type PostVotes struct {
	PostId    int
	UserId    int
	VoteValue int
}
