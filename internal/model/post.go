package model

import "time"

type Post struct {
	Id             int
	GUID           string
	Title          string
	User           string
	Content        string
	Image          []byte
	ImageExtention string
	CreatedAt      time.Time
	UpdatedAt      time.Time
	Active         int
}
