package page

import (
	"context"

	"github.com/Anacardo89/lenic/internal/repo"
	"github.com/Anacardo89/lenic/internal/session"
	"github.com/Anacardo89/lenic/pkg/logger"
)

type PageHandler struct {
	ctx          context.Context
	db           repo.DBRepository
	sessionStore *session.SessionStore
	log          *logger.Logger
}

func NewHandler(ctx context.Context, l *logger.Logger, db repo.DBRepository, sessionStore *session.SessionStore) *PageHandler {
	return &PageHandler{
		ctx:          ctx,
		log:          l,
		db:           db,
		sessionStore: sessionStore,
	}
}

func (h *PageHandler) decodeUser() (string, error) {
	return "", nil
}
