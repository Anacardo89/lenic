package model

import "time"

type User struct {
	Id        int
	UserName  string
	UserEmail string
	UserPass  string
	FirstName string
	LastName  string
	CreatedAt time.Time
	UpdatedAt time.Time
	Active    int
}
