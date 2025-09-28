package page

import (
	"context"

	"github.com/Anacardo89/lenic/internal/repo"
	"github.com/Anacardo89/lenic/internal/session"
	"github.com/Anacardo89/lenic/pkg/logger"
)

type PageHandler struct {
	ctx context.Context
	log *logger.Logger
	db  repo.DBRepository
	sm  *session.SessionManager
}

func NewHandler(
	ctx context.Context,
	l *logger.Logger,
	db repo.DBRepository,
	sm *session.SessionManager,
) *PageHandler {
	return &PageHandler{
		ctx: ctx,
		log: l,
		db:  db,
		sm:  sm,
	}
}

func (h *PageHandler) decodeUser() (string, error) {
	return "", nil
}
