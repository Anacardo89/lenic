package middleware

import (
	"net/http"
	"time"

	"github.com/Anacardo89/lenic/internal/session"
	"github.com/Anacardo89/lenic/pkg/logger"
)

type MiddlewareHandler struct {
	log          *logger.Logger
	sm           *session.SessionManager
	writeTimeout time.Duration
}

func NewMiddlewareHandler(sm *session.SessionManager, l *logger.Logger, wto time.Duration) *MiddlewareHandler {
	return &MiddlewareHandler{
		log:          l,
		sm:           sm,
		writeTimeout: wto - time.Second,
	}
}

func (m *MiddlewareHandler) Wrap(next http.Handler) http.Handler {
	return m.Log(m.Timeout(next))
}
