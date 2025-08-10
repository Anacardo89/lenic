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
	UserID      uuid.UUID
	TargetID    uuid.UUID
	ResourceTpe string
}

type HashTag struct {
	ID        uuid.UUID
	TagName   string
	CreatedAt time.Time
}

type HashTagResource struct {
	TagID       uuid.UUID
	TargetID    uuid.UUID
	ResourceTpe string
}
