package middleware

import (
	"context"
	"net/http"
)

type CtxKey string

const (
	CtxKeySession CtxKey = "session"
)

func (m *MiddlewareHandler) Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session := m.sm.ValidateSession(w, r)
		if session == nil {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		if !session.IsAuthenticated {
			m.log.Info("unauthorized action attempt",
				"method", r.Method,
				"path", r.URL.Path,
				"client_ip", r.RemoteAddr,
			)
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		ctx := context.WithValue(r.Context(), CtxKeySession, session)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
