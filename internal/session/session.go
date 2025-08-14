package session

import (
	"context"
	"net/http"
	"sync"
	"time"

	"github.com/Anacardo89/lenic/internal/config"
	"github.com/Anacardo89/lenic/internal/db"

	"github.com/Anacardo89/lenic/internal/models"
	"github.com/Anacardo89/lenic/pkg/logger"
	"github.com/google/uuid"
	"github.com/gorilla/sessions"
)

type SessionStore struct {
	ctx      context.Context
	cfg      *config.SessionConfig
	mu       sync.Mutex
	db       db.DBRepository
	store    *sessions.CookieStore
	sessions map[string]*Session
}

func NewSessionStore(ctx context.Context, cfg *config.SessionConfig, store *sessions.CookieStore, db db.DBRepository) *SessionStore {
	return &SessionStore{
		ctx:      ctx,
		db:       db,
		store:    store,
		sessions: make(map[string]*Session),
	}
}

type Session struct {
	ID              uuid.UUID
	IsAuthenticated bool                   `json:"is_authenticated"`
	User            *models.User           `json:"user"`
	Notifs          []*models.Notification `json:"notifs"`
	DMs             []*models.Conversation `json:"dms"`
	UpdatedAt       time.Time              `json:"updated_at"`
}

func (s *SessionStore) Store() *sessions.CookieStore {
	return s.store
}

func (s *SessionStore) DeleteSession(r *http.Request) {
	lenicSession, err := s.store.Get(r, "lenic_session")
	if err != nil {
		logger.Error.Println(err)
	}
	sessionID := lenicSession.Values["session_id"]
	s.mu.Lock()
	delete(s.sessions, sessionID)
	s.mu.Unlock()

	lenicSession.Options.MaxAge = -1
	err = lenicSession.Save(r, w)
	if err != nil {
		logger.Error.Println(err)
	}

	return
}

func (s *SessionStore) CreateSession(w http.ResponseWriter, r *http.Request, userID uuid.UUID) *Session {
	lenicSession, err := s.store.Get(r, "lenic_session")
	if err != nil {
		logger.Error.Println(err)
	}
	sessionID := uuid.New()
	lenicSession.Values["session_id"] = sessionID
	err = lenicSession.Save(r, w)
	if err != nil {
		logger.Error.Println(err)
	}
	dbUser, err := s.db.GetUserByID(s.ctx, userID)
	if err != nil {
		logger.Error.Println(err)
	}
	u := models.FromDBUser(dbUser)

	session := &Session{
		ID:              sessionID,
		IsAuthenticated: true,
		User:            u,
		UpdatedAt:       time.Now(),
	}
	s.mu.Lock()
	s.sessions[userID] = session
	s.mu.Unlock()
	return session
}

func (s *SessionStore) ValidateSession(w http.ResponseWriter, r *http.Request) *Session {
	// Error handling
	deleteSession := func(sessionID string) {
		s.mu.Lock()
		delete(s.sessions, sessionID)
		s.mu.Unlock()
	}
	//

	// Execution
	session := &Session{IsAuthenticated: false}
	lenicSession, err := s.store.Get(r, "lenic_session")
	if err != nil {
		logger.Error.Println(err)
	}
	sessionID, ok := lenicSession.Values["session_id"]
	if !ok {
		return session
	}
	s.mu.Lock()
	session, ok = s.sessions[sessionID]
	s.mu.Unlock()
	if !ok {
		return session
	}
	if time.Now().After(session.UpdatedAt.Add(time.Duration(24) * time.Hour)) {
		deleteSession(sessionID)
		lenicSession.Options.MaxAge = -1
		lenicSession.Save(r, w)
		return session
	}
	dbUser, err := s.db.GetUserByID(s.ctx, &session.User.ID)
	if err != nil {
		deleteSession(sessionID)
		lenicSession.Options.MaxAge = -1
		lenicSession.Save(r, w)
		return session
	}
	u := models.FromDBUser(dbUser)
	session.User = u
	session.IsAuthenticated = true
	session.UpdatedAt = time.Now()

	s.mu.Lock()
	s.sessions[sessionID] = session

	return session
}
