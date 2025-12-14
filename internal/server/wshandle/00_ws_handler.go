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
	cancel     context.CancelFunc
	db         repo.DBRepo
	log        *logger.Logger
	sm         *session.SessionManager
	wsConnMann *wsconnman.WSConnMan
}

func NewHandler(
	parentCtx context.Context,
	db repo.DBRepo,
	l *logger.Logger,
	sm *session.SessionManager,
	wsConnMan *wsconnman.WSConnMan,
) *WSHandler {
	ctx, cancel := context.WithCancel(parentCtx)
	return &WSHandler{
		ctx:        ctx,
		cancel:     cancel,
		db:         db,
		log:        l,
		sm:         sm,
		wsConnMann: wsConnMan,
	}
}
