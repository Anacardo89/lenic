package database

import (
	"time"
)

type User struct {
	Id            int
	UserName      string
	Email         string
	HashPass      string
	ProfilePic    string
	ProfilePicExt string
	CreatedAt     time.Time
	UpdatedAt     time.Time
	Active        int
}

type Follows struct {
	FollowerId int
	FollowedId int
}
