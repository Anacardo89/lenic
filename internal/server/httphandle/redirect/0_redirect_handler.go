package redirect

import (
	"github.com/Anacardo89/lenic/internal/config"
	"github.com/Anacardo89/lenic/internal/db"
	"github.com/Anacardo89/lenic/internal/session"
)

type RedirectHandler struct {
	cfg          *config.ServerConfig
	db           db.DBRepository
	sessionStore *session.SessionStore
}

func NewHandler(cfg *config.ServerConfig, db db.DBRepository, sessionStore *session.SessionStore) *RedirectHandler {
	return &RedirectHandler{
		cfg:          cfg,
		db:           db,
		sessionStore: sessionStore,
	}
}
