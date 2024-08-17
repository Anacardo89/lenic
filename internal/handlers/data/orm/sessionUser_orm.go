package orm

import (
	"github.com/Anacardo89/tpsi25_blog/internal/model/database"
	"github.com/Anacardo89/tpsi25_blog/pkg/logger"
)

func GetUserBySessionID(sid string) database.User {
	session, err := Da.GetSessionBySessionID(sid)
	if err != nil {
		logger.Error.Println(err)
	}
	user, err := Da.GetUserByID(session.UserId)
	if err != nil {
		logger.Error.Println(err)
	}
	return *user
}
