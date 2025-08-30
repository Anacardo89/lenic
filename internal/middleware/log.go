package middleware

import (
	"net/http"
	"time"
)

func (m *MiddlewareHandler) Log(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		m.log.Info("request received",
			"method", r.Method,
			"path", r.URL.Path,
			"query", r.URL.RawQuery,
			"time_received", start,
			"client_ip", r.RemoteAddr,
		)

		rw := newLogRW(w)
		next.ServeHTTP(rw, r)

		duration := time.Since(start)
		m.log.Info("request completed",
			"method", r.Method,
			"path", r.URL.Path,
			"query", r.URL.RawQuery,
			"status", rw.Status(),
			"size", rw.Size(),
			"duration_ms", duration.Milliseconds(),
			"client_ip", r.RemoteAddr,
		)
	})
}
