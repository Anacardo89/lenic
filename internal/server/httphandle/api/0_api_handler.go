package api

import (
	"context"

	"github.com/Anacardo89/lenic/config"
	"github.com/Anacardo89/lenic/internal/auth"
	"github.com/Anacardo89/lenic/internal/repo"
	"github.com/Anacardo89/lenic/internal/server/wshandle"
	"github.com/Anacardo89/lenic/internal/session"
	"github.com/Anacardo89/lenic/pkg/logger"
	"github.com/Anacardo89/lenic/pkg/mail"
)

type APIHandler struct {
	ctx          context.Context
	cfg          *config.Server
	db           repo.DBRepository
	tokenManager *auth.TokenManager
	sm           *session.SessionManager
	wsHandler    *wshandle.WSHandler
	mail         *mail.Client
	log          *logger.Logger
}

func NewHandler(
	ctx context.Context,
	l *logger.Logger,
	cfg *config.Server,
	db repo.DBRepository,
	tm *auth.TokenManager,
	sm *session.SessionManager,
	wsHandler *wshandle.WSHandler,
	mail *mail.Client,
) *APIHandler {
	return &APIHandler{
		ctx:          ctx,
		cfg:          cfg,
		log:          l,
		db:           db,
		tokenManager: tm,
		sm:           sm,
		wsHandler:    wsHandler,
		mail:         mail,
	}
}
