package page

import (
	"github.com/Anacardo89/lenic/internal/db"
	"github.com/Anacardo89/lenic/internal/session"
)

type PageHandler struct {
	db           db.DBRepository
	sessionStore *session.SessionStore
}

func NewHandler(db db.DBRepository, sessionStore *session.SessionStore) *PageHandler {
	return &PageHandler{
		db:           db,
		sessionStore: sessionStore,
	}
}
