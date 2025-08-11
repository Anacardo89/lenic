package server

import (
	"net/http"

	"github.com/Anacardo89/lenic/config"
	"github.com/Anacardo89/lenic/internal/auth"
	"github.com/Anacardo89/lenic/internal/db"
	"github.com/Anacardo89/lenic/internal/server/handlers/api"
	"github.com/Anacardo89/lenic/internal/server/handlers/page"
	"github.com/Anacardo89/lenic/internal/server/handlers/redirect"
	"github.com/Anacardo89/lenic/internal/server/http/redirect"
	"github.com/Anacardo89/lenic/internal/server/websocket"
	"github.com/Anacardo89/lenic/internal/wsconnman"
	"github.com/gorilla/mux"
)

type Server struct {
	cfg          *config.Config
	db           db.DBRepository
	sessionStore *auth.SessionStore
	wsConnMann   *wsconnman.WSConnMan
	router       http.Handler

	apiHandler      *api.Handler
	pageHandler     *page.Handler
	redirectHandler *redirect.Handler
	websocketHanler *websocket.Handler
}

func NewServer(cfg *config.Config, db db.DBRepository, sessionStore *auth.SessionStore) *Server {
	s := &Server{
		cfg:             cfg,
		db:              db,
		sessionStore:    sessionStore,
		apiHandler:      api.NewHandler(db, sessionStore),
		pageHandler:     page.NewHandler(db, sessionStore),
		redirectHandler: redirect.NewHandler(cfg, db, sessionStore),
	}

	r := mux.NewRouter()
	s.DeclareRoutes(r)
	s.router = r
	return s
}
