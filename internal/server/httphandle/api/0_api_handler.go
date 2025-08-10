package api

import (
	"github.com/Anacardo89/lenic/internal/db"
	"github.com/Anacardo89/lenic/internal/session"
)

type APIHandler struct {
	db           db.DBRepository
	sessionStore *session.SessionStore
}

func NewHandler(db db.DBRepository, sessionStore *session.SessionStore) *APIHandler {
	return &APIHandler{
		db:           db,
		sessionStore: sessionStore,
	}
}
