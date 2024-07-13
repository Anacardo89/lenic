package db

import "time"

type User struct {
	Id         int
	UserName   string
	UserEmail  string
	UserPasswd string
	UserSalt   string
	CreatedAt  time.Time
	UpdatedAt  time.Time
	IsActive   bool
}
