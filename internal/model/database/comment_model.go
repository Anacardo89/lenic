package database

import "time"

type Comment struct {
	Id        int
	PostGUID  string
	AuthorId  int
	Content   string
	CreatedAt time.Time
	UpdatedAt time.Time
	Rating    int
	Active    int
}

type CommentVotes struct {
	CommentId   int
	UserId      int
	RatingValue int
}
