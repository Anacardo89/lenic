package page

import (
	"context"

	"github.com/Anacardo89/lenic/internal/repo"
	"github.com/Anacardo89/lenic/internal/session"
	"github.com/Anacardo89/lenic/pkg/logger"
)

type PageHandler struct {
	ctx     context.Context
	homeDir string
	log     *logger.Logger
	db      repo.DBRepo
	sm      *session.SessionManager
}

func NewHandler(
	ctx context.Context,
	homeDir string,
	l *logger.Logger,
	db repo.DBRepo,
	sm *session.SessionManager,
) *PageHandler {
	return &PageHandler{
		ctx:     ctx,
		homeDir: homeDir,
		log:     l,
		db:      db,
		sm:      sm,
	}
}
