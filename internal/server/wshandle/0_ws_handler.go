package wshandle

import (
	"context"

	"github.com/Anacardo89/lenic/internal/repo"
	"github.com/Anacardo89/lenic/internal/session"
	"github.com/Anacardo89/lenic/internal/wsconnman"
	"github.com/Anacardo89/lenic/pkg/logger"
)

type WSHandler struct {
	ctx        context.Context
	db         repo.DBRepository
	log        *logger.Logger
	sm         *session.SessionManager
	wsConnMann *wsconnman.WSConnMan
}

func NewHandler(
	db repo.DBRepository,
	l *logger.Logger,
	sm *session.SessionManager,
	wsConnMan *wsconnman.WSConnMan,
) *WSHandler {
	return &WSHandler{
		db:         db,
		log:        l,
		sm:         sm,
		wsConnMann: wsConnMan,
	}
}
