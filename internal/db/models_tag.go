package db

import (
	"time"

	"github.com/google/uuid"
)

type ResourceType int

const (
	ResourcePost ResourceType = iota
	ResourceComment
)

var (
	resourceTypeList = []string{
		"post",
		"comment",
	}
)

func (r ResourceType) String() string {
	return resourceTypeList[r]
}

type UserTag struct {
	ID               uuid.UUID
	TaggedResourceID uuid.UUID
	ResourceTpe      string
}

type Hashtag struct {
	ID        uuid.UUID
	TagName   string
	CreatedAt time.Time
}

type HashtagResource struct {
	TagID            uuid.UUID
	TaggedResourceID uuid.UUID
	ResourceTpe      string
}
