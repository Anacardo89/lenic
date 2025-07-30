package db

import (
	"time"

	"github.com/google/uuid"
)

type Post struct {
	ID        uuid.UUID
	AuthorID  uuid.UUID
	Title     string
	Content   string
	PostImage string
	Rating    int
	IsPublic  bool
	IsActive  bool
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt time.Time
}

type PostRatings struct {
	TargetID    uuid.UUID
	UserID      uuid.UUID
	RatingValue int
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
