package db

import (
	"time"

	"github.com/Anacardo89/tpsi25_blog.git/logger"
)

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

func UpdateSession(sid string, uid int) {
	const timeFmt = "2006-01-02T15:04:05.999999999"
	tstamp := time.Now().Format(timeFmt)
	_, err := Dbase.Exec(InsertSession, 1, sid, uid, tstamp, uid, tstamp)
	if err != nil {
		logger.Error.Println(err)
	}
}

func GetSessionUID(sid string) User {
	user := User{}
	err := Dbase.QueryRow(SelectUserFromSessions, sid).Scan(&user.Id)
	if err != nil {
		logger.Error.Println(err)
		return User{}
	}
	err = Dbase.QueryRow(SelectUserById, user.Id).Scan(&user.UserName)
	if err != nil {
		logger.Error.Println(err)
		return User{}
	}
	return user
}
