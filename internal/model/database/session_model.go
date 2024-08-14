package database

import "time"

type Session struct {
	Id            int
	SessionId     string
	UserId        int
	SessionStart  time.Time
	SessionUpdate time.Time
	Active        int
}
