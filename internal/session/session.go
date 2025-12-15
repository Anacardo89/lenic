package session

import (
	"context"
	"net/http"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/sessions"

	"github.com/Anacardo89/lenic/config"
	"github.com/Anacardo89/lenic/internal/models"
	"github.com/Anacardo89/lenic/internal/repo"
)

// TODO - Get this to Redis

type SessionManager struct {
	ctx      context.Context
	cfg      *config.Session
	mu       sync.Mutex
	db       repo.DBRepo
	store    *sessions.CookieStore
	sessions map[uuid.UUID]*Session
}

func NewSessionManager(ctx context.Context, cfg *config.Session, db repo.DBRepo) *SessionManager {
	store := sessions.NewCookieStore([]byte(cfg.Secret))
	store.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   int(cfg.Duration.Seconds()),
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
	}
	return &SessionManager{
		ctx:      ctx,
		cfg:      cfg,
		db:       db,
		store:    store,
		sessions: make(map[uuid.UUID]*Session),
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

func (s *SessionManager) Store() *sessions.CookieStore {
	return s.store
}

func (s *SessionManager) CreateSession(w http.ResponseWriter, r *http.Request, userID uuid.UUID) (*Session, error) {
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
	s.sessions[sessionID] = session
	s.mu.Unlock()
	return session, nil
}

func (s *SessionManager) ValidateSession(w http.ResponseWriter, r *http.Request) *Session {
	// Error handling
	deleteSession := func(sessionID uuid.UUID) {
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
	sessionIDVal, ok := lenicSession.Values["session_id"]
	if !ok {
		return session
	}
	sessionID, ok := sessionIDVal.(uuid.UUID)
	if !ok {
		return session
	}
	s.mu.Lock()
	session, ok = s.sessions[sessionID]
	s.mu.Unlock()
	if !ok {
		return session
	}
	if time.Now().After(session.UpdatedAt.Add(s.cfg.Duration)) {
		deleteSession(sessionID)
		lenicSession.Options.MaxAge = -1
		lenicSession.Save(r, w)
		return session
	}
	dbUser, err := s.db.GetUserByID(r.Context(), session.User.ID)
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
	s.mu.Unlock()

	return session
}

func (s *SessionManager) DeleteSession(w http.ResponseWriter, r *http.Request) error {
	lenicSession, err := s.store.Get(r, "lenic_session")
	if err != nil {
		return err
	}
	sessionIDVal := lenicSession.Values["session_id"]
	sessionID := sessionIDVal.(uuid.UUID)
	if err != nil {
		return err
	}
	s.mu.Lock()
	delete(s.sessions, sessionID)
	s.mu.Unlock()

	lenicSession.Options.MaxAge = -1
	err = lenicSession.Save(r, w)
	return err
}
