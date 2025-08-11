package api

import (
	"context"

	"github.com/Anacardo89/lenic/internal/db"
	"github.com/Anacardo89/lenic/internal/session"
)

type APIHandler struct {
	ctx          context.Context
	db           db.DBRepository
	sessionStore *session.SessionStore
}

func NewHandler(ctx context.Context, db db.DBRepository, sessionStore *session.SessionStore) *APIHandler {
	return &APIHandler{
		ctx:          ctx,
		db:           db,
		sessionStore: sessionStore,
	}
}
