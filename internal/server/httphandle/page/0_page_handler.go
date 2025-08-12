package page

import (
	"context"

	"github.com/Anacardo89/lenic/internal/db"
	"github.com/Anacardo89/lenic/internal/session"
)

type PageHandler struct {
	ctx          context.Context
	db           db.DBRepository
	sessionStore *session.SessionStore
}

func NewHandler(ctx context.Context, db db.DBRepository, sessionStore *session.SessionStore) *PageHandler {
	return &PageHandler{
		ctx:          ctx,
		db:           db,
		sessionStore: sessionStore,
	}
}

func (h *PageHandler) decodeUser() (string, error) {
	return "", nil
}
