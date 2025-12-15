package server

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/Anacardo89/lenic/config"
	"github.com/Anacardo89/lenic/internal/middleware"
	"github.com/Anacardo89/lenic/internal/server/httphandle/api"
	"github.com/Anacardo89/lenic/internal/server/httphandle/page"
	"github.com/Anacardo89/lenic/internal/server/wshandle"
	"github.com/Anacardo89/lenic/pkg/logger"
)

type Server struct {
	httpSrv  *http.Server
	router   http.Handler
	addr     string
	log      *logger.Logger
	timeouts ServerTimeouts
}

type ServerTimeouts struct {
	ReadTimeout     time.Duration
	WriteTimeout    time.Duration
	ShutdownTimeout time.Duration
}

func NewServer(
	cfg *config.Server,
	l *logger.Logger,
	ah *api.APIHandler,
	ph *page.PageHandler,
	mw *middleware.MiddlewareHandler,
	wsh *wshandle.WSHandler,
) *Server {
	to := ServerTimeouts{
		ReadTimeout:     cfg.ReadTimeout,
		WriteTimeout:    cfg.WriteTimeout,
		ShutdownTimeout: cfg.ShutdownTimeout,
	}
	s := &Server{
		router:   NewRouter(ah, ph, wsh, mw),
		addr:     fmt.Sprintf(":%s", cfg.Port),
		log:      l,
		timeouts: to,
	}
	return s
}

func (s *Server) Start() error {
	s.httpSrv = &http.Server{
		Addr:         s.addr,
		Handler:      s.router,
		ReadTimeout:  s.timeouts.ReadTimeout,
		WriteTimeout: s.timeouts.WriteTimeout,
	}
	s.log.Info("Starting server on", "adress", s.addr)
	return s.httpSrv.ListenAndServe()
}

func (s *Server) Shutdown() error {
	ctx, cancel := context.WithTimeout(context.Background(), s.timeouts.ShutdownTimeout)
	defer cancel()
	if s.httpSrv != nil {
		return s.httpSrv.Shutdown(ctx)
	}
	return nil
}
