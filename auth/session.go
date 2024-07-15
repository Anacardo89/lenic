package auth

import (
	"crypto/rand"
	"encoding/base64"
	"io"
	"net/http"

	"github.com/Anacardo89/tpsi25_blog.git/logger"
	"github.com/gorilla/sessions"
)

type SessionConfig struct {
	Pass string `yaml:"pass"`
}

type User struct {
	Id         int
	UserName   string
	UserEmail  string
	UserPass   string
	HashedPass string
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
	session.Save(r, w)
	usrSession.Id = newSID
	usrSession.Authenticated = true
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
