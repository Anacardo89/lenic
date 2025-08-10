package db

import (
	"time"

	"github.com/google/uuid"
)

// Users
type UserRole int

const (
	RoleUser UserRole = iota
	RoleMod
	RoleAdmin
)

var (
	userRoleList = []string{
		"user",
		"moderator",
		"admin",
	}
)

func (r UserRole) String() string {
	return userRoleList[r]
}

type User struct {
	ID           uuid.UUID
	UserName     string
	DisplayName  string
	Email        string
	PasswordHash string
	ProfilePic   string
	Bio          string
	Followers    int
	Following    int
	IsActive     bool
	IsVerified   bool
	UserRole     string
	LastLogin    time.Time
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    *time.Time
}

// Follows
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
