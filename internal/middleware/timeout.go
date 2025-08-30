package middleware

import (
	"context"
	"net/http"
)

func (m *MiddlewareHandler) Timeout(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), m.writeTimeout)
		defer cancel()
		done := make(chan struct{})
		go func() {
			next.ServeHTTP(w, r.WithContext(ctx))
			close(done)
		}()
		select {
		case <-ctx.Done():
			if ctx.Err() == context.DeadlineExceeded {
				w.WriteHeader(http.StatusGatewayTimeout)
				w.Write([]byte("request timed out"))
			}
		case <-done:
		}
	})
}
