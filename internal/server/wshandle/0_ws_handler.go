package wshandle

import (
	"context"

	"github.com/Anacardo89/lenic/internal/db"
	"github.com/Anacardo89/lenic/internal/session"
	"github.com/Anacardo89/lenic/internal/wsconnman"
)

type WSHandler struct {
	ctx          context.Context
	db           db.DBRepository
	sessionStore *session.SessionStore
	wsConnMann   *wsconnman.WSConnMan
}

func NewHandler(db db.DBRepository, sessionStore *session.SessionStore, wsConnMan *wsconnman.WSConnMan) *WSHandler {
	return &WSHandler{
		db:           db,
		sessionStore: sessionStore,
		wsConnMann:   wsConnMan,
	}
}
