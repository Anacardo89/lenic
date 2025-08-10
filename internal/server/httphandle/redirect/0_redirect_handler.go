package redirect

import (
	"github.com/Anacardo89/lenic/internal/db"
	"github.com/Anacardo89/lenic/internal/session"
)

type RedirectHandler struct {
	db           db.DBRepository
	sessionStore *session.SessionStore
}

func NewHandler(db db.DBRepository, sessionStore *session.SessionStore) *RedirectHandler {
	return &RedirectHandler{
		db:           db,
		sessionStore: sessionStore,
	}
}
