package server

import (
	"net/http"

	"github.com/Anacardo89/tpsi25_blog/config"
	"github.com/Anacardo89/tpsi25_blog/internal/auth"
	"github.com/Anacardo89/tpsi25_blog/internal/db"
	"github.com/gorilla/mux"
)

type Server struct {
	cfg          *config.Config
	DB           db.DBRepository
	SessionStore *auth.SessionStore
	router       http.Handler
}

func NewServer(cfg *config.Config, db db.DBRepository, sessionStore *auth.SessionStore) *Server {
	s := &Server{
		cfg:          cfg,
		DB:           db,
		SessionStore: sessionStore,
	}

	r := mux.NewRouter()
	s.DeclareRoutes(r)
	s.router = r
	return s
}
