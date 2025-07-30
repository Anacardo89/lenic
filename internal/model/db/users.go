package db

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID          uuid.UUID
	UserName    string
	DisplayName string
	Email       string
	HashPass    string
	ProfilePic  string
	Bio         string
	Followers   int
	Following   int
	IsActive    bool
	IsVerified  bool
	LastLogin   time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   time.Time
}

type FollowStatus int

const (
	StatusPending FollowStatus = iota
	StatusAccepted
	StatusBlocked
)

var (
	followStatusList = []string{
		"pending",
		"accepted",
		"blocked",
	}
)

func (f FollowStatus) String() string {
	return followStatusList[f]
}

type Follows struct {
	FollowerID   uuid.UUID
	FollowedID   uuid.UUID
	FollowStatus string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
