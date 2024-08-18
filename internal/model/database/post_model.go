package database

import "time"

type Post struct {
	Id             int
	GUID           string
	Title          string
	User           string
	Content        string
	Image          string
	ImageExtention string
	CreatedAt      time.Time
	UpdatedAt      time.Time
	Active         int
}
