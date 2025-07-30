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

type HashTag struct {
	ID        uuid.UUID
	TadName   string
	CreatedAt time.Time
}
