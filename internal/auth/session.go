package auth

import (
	"crypto/rand"
	"encoding/base64"
	"io"
	"net/http"
	"time"

	"github.com/Anacardo89/tpsi25_blog/internal/handlers/data/orm"
	"github.com/Anacardo89/tpsi25_blog/internal/model/database"
	"github.com/Anacardo89/tpsi25_blog/internal/model/mapper"
	"github.com/Anacardo89/tpsi25_blog/internal/model/presentation"
	"github.com/Anacardo89/tpsi25_blog/pkg/logger"
	"github.com/gorilla/sessions"
)

type SessionStore struct {
	store    *sessions.CookieStore
	sessions map[string]presentation.Session
}

func NewSessionStore(store *sessions.CookieStore) *SessionStore {
	return &SessionStore{
		store:    store,
		sessions: make(map[string]presentation.Session),
	}
}

func (s *SessionStore) CreateSession(w http.ResponseWriter, r *http.Request) *presentation.Session {
	usrSession := presentation.Session{}
	session, err := SessionStore.Get(r, "lenic")
	if err != nil {
		logger.Error.Println(err)
	}
	newSID := generateSessionId()
	session.Values["sid"] = newSID
	err = session.Save(r, w)
	if err != nil {
		logger.Error.Println(err)
	}
	usrSession.SessionId = newSID
	usrSession.Authenticated = true
	return usrSession
}

func ValidateSession(w http.ResponseWriter, r *http.Request) presentation.Session {
	usrSession := presentation.Session{
		Authenticated: false,
	}
	session, err := SessionStore.Get(r, "lenic")
	if err != nil {
		logger.Error.Println(err)
	}
	if sid, valid := session.Values["sid"]; valid {
		dbsession, err := orm.Da.GetSessionBySessionID(sid.(string))
		if err != nil {
			logger.Error.Println(err)
			return usrSession
		}
		if time.Now().After(dbsession.UpdatedAt.Add(time.Duration(24) * time.Hour)) {
			session.Options.MaxAge = -1
			session.Save(r, w)
			return usrSession
		}
		dbuser, err := orm.GetUserBySessionID(sid.(string))
		if err != nil {
			session.Options.MaxAge = -1
			session.Save(r, w)
			return usrSession
		}
		u := mapper.User(dbuser)
		usrSession.User = *u
		UpdateSession(sid.(string), usrSession.User.Id)
		usrSession.SessionId = sid.(string)
		usrSession.Authenticated = true
	}
	return usrSession
}

func UpdateSession(sid string, uid int) {
	s := &database.Session{
		SessionId: sid,
		UserId:    uid,
		Active:    1,
	}
	if err := orm.Da.CreateSession(s); err != nil {
		logger.Error.Println(err)
	}
}

func generateSessionId() string {
	sid := make([]byte, 24)
	_, err := io.ReadFull(rand.Reader, sid)
	if err != nil {
		logger.Error.Println(err)
	}
	return base64.URLEncoding.EncodeToString(sid)
}
