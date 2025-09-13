package session

import (
	"context"
	"net/http"
	"sync"
	"time"

	"github.com/Anacardo89/lenic/config"
	"github.com/Anacardo89/lenic/internal/repo"

	"github.com/Anacardo89/lenic/internal/models"
	"github.com/google/uuid"
	"github.com/gorilla/sessions"
)

// TODO - Get this to Redis

type SessionStore struct {
	ctx      context.Context
	cfg      config.Session
	mu       sync.Mutex
	db       repo.DBRepository
	store    *sessions.CookieStore
	sessions map[string]*Session
}

func NewSessionStore(ctx context.Context, cfg config.Session, db repo.DBRepository) *SessionStore {
	return &SessionStore{
		ctx:      ctx,
		cfg:      cfg,
		db:       db,
		store:    sessions.NewCookieStore([]byte(cfg.Secret)),
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

func (s *SessionStore) DeleteSession(w http.ResponseWriter, r *http.Request) error {
	lenicSession, err := s.store.Get(r, "lenic_session")
	if err != nil {
		return err
	}
	sessionID := lenicSession.Values["session_id"]
	s.mu.Lock()
	delete(s.sessions, sessionID.(uuid.UUID).String())
	s.mu.Unlock()

	lenicSession.Options.MaxAge = -1
	err = lenicSession.Save(r, w)
	return err
}

func (s *SessionStore) CreateSession(w http.ResponseWriter, r *http.Request, userID uuid.UUID) (*Session, error) {
	lenicSession, err := s.store.Get(r, "lenic_session")
	if err != nil {
		return nil, err
	}
	sessionID := uuid.New()
	lenicSession.Values["session_id"] = sessionID
	err = lenicSession.Save(r, w)
	if err != nil {
		return nil, err
	}
	dbUser, err := s.db.GetUserByID(s.ctx, userID)
	if err != nil {
		return nil, err
	}
	u := models.FromDBUser(dbUser)

	session := &Session{
		ID:              sessionID,
		IsAuthenticated: true,
		User:            u,
		UpdatedAt:       time.Now(),
	}
	s.mu.Lock()
	s.sessions[userID.String()] = session
	s.mu.Unlock()
	return session, nil
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
		return nil
	}
	sessionID, ok := lenicSession.Values["session_id"]
	if !ok {
		return session
	}
	s.mu.Lock()
	session, ok = s.sessions[sessionID.(uuid.UUID).String()]
	s.mu.Unlock()
	if !ok {
		return session
	}
	if time.Now().After(session.UpdatedAt.Add(time.Duration(24) * time.Hour)) {
		deleteSession(sessionID.(uuid.UUID).String())
		lenicSession.Options.MaxAge = -1
		lenicSession.Save(r, w)
		return session
	}
	dbUser, err := s.db.GetUserByID(r.Context(), session.User.ID)
	if err != nil {
		deleteSession(sessionID.(uuid.UUID).String())
		lenicSession.Options.MaxAge = -1
		lenicSession.Save(r, w)
		return session
	}
	u := models.FromDBUser(dbUser)
	session.User = u
	session.IsAuthenticated = true
	session.UpdatedAt = time.Now()

	s.mu.Lock()
	s.sessions[sessionID.(uuid.UUID).String()] = session

	return session
}
