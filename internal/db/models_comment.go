package db

import (
	"time"

	"github.com/google/uuid"
)

type Comment struct {
	ID        uuid.UUID
	PostID    uuid.UUID
	AuthorID  uuid.UUID
	Content   string
	Rating    int
	IsActive  bool
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}

type CommentRatings struct {
	TargetID    uuid.UUID
	UserID      uuid.UUID
	RatingValue int
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
