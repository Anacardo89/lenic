package database

import "time"

type Notification struct {
	Id         int
	UserID     int
	FromUserId int
	NotifType  string
	NotifMsg   string
	ResourceId int
	IsRead     bool
	CreatedAt  time.Time
	UpdatedAt  time.Time
}
