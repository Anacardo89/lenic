package auth

import (
	"crypto/rand"
	"encoding/base64"
	"io"
	"net/http"
	"time"

	"github.com/Anacardo89/tpsi25_blog/internal/model"
	"github.com/Anacardo89/tpsi25_blog/internal/query"
	"github.com/Anacardo89/tpsi25_blog/pkg/db"
	"github.com/Anacardo89/tpsi25_blog/pkg/logger"
	"github.com/gorilla/sessions"
)

type Config struct {
	Pass string `yaml:"session_pass"`
}

type User struct {
	Id         int
	UserName   string
	UserEmail  string
	UserPass   string
	HashedPass string
	Active     int
}

type Session struct {
	Id            string
	Authenticated bool
	User          User
}

var (
	SessionStore *sessions.CookieStore
)

func CreateSession(w http.ResponseWriter, r *http.Request) Session {
	usrSession := Session{}
	session, err := SessionStore.Get(r, "tpsi25blog")
	if err != nil {
		logger.Error.Println(err)
	}
	newSID := generateSessionId()
	session.Values["sid"] = newSID
	err = session.Save(r, w)
	if err != nil {
		logger.Error.Println(err)
	}
	usrSession.Id = newSID
	usrSession.Authenticated = true
	return usrSession
}

func ValidateSession(r *http.Request) Session {
	usrSession := Session{}
	session, err := SessionStore.Get(r, "tpsi25blog")
	if err != nil {
		logger.Error.Println(err)
	}
	if sid, valid := session.Values["sid"]; valid {
		user := GetSessionUID(sid.(string))
		usrSession.User = User{
			Id:        user.Id,
			UserName:  user.UserName,
			UserEmail: user.UserEmail,
		}
		UpdateSession(sid.(string), user.Id)
		usrSession.Id = sid.(string)
		usrSession.Authenticated = true
	} else {
		usrSession.Authenticated = false
	}
	return usrSession
}

func generateSessionId() string {
	sid := make([]byte, 24)
	_, err := io.ReadFull(rand.Reader, sid)
	if err != nil {
		logger.Error.Println(err)
	}
	return base64.URLEncoding.EncodeToString(sid)
}

func UpdateSession(sid string, uid int) {
	const timeFmt = "2006-01-02T15:04:05.999999999"
	tstamp := time.Now().Format(timeFmt)
	_, err := db.Dbase.Exec(query.InsertSession, 1, sid, uid, tstamp, uid, tstamp)
	if err != nil {
		logger.Error.Println(err)
	}
}

func GetSessionUID(sid string) model.User {
	user := model.User{}
	err := db.Dbase.QueryRow(query.SelectUserFromSessions, sid).Scan(&user.Id)
	if err != nil {
		logger.Error.Println(err)
		return model.User{}
	}
	err = db.Dbase.QueryRow(query.SelectUserById, user.Id).Scan(&user.UserName)
	if err != nil {
		logger.Error.Println(err)
		return model.User{}
	}
	return user
}
