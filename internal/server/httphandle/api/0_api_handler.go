package api

import (
	"context"

	"github.com/Anacardo89/lenic/internal/db"
	"github.com/Anacardo89/lenic/internal/server/wshandle"
	"github.com/Anacardo89/lenic/internal/session"
)

type APIHandler struct {
	ctx          context.Context
	db           db.DBRepository
	sessionStore *session.SessionStore
	wsHandler    *wshandle.WSHandler
}

func NewHandler(ctx context.Context, db db.DBRepository, sessionStore *session.SessionStore, wsHandler *wshandle.WSHandler) *APIHandler {
	return &APIHandler{
		ctx:          ctx,
		db:           db,
		sessionStore: sessionStore,
		wsHandler:    wsHandler,
	}
}
