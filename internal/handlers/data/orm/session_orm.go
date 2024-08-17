package orm

import (
	"time"

	"github.com/Anacardo89/tpsi25_blog/internal/handlers/data/query"
	"github.com/Anacardo89/tpsi25_blog/internal/model/database"
	"github.com/Anacardo89/tpsi25_blog/pkg/db"
)

func (da *DataAccess) CreateSession(s *database.Session) error {
	_, err := da.Db.Exec(query.InsertUser,
		s.SessionId,
		s.UserId,
		s.Active)
	return err
}

func (da *DataAccess) GetSessionByID(id int) (*database.Session, error) {
	var (
		sessionStart  []byte
		sessionUpdate []byte
	)
	s := database.Session{}
	row := da.Db.QueryRow(query.SelectSessionById, id)
	err := row.Scan(
		&s.Id,
		&s.SessionId,
		&s.UserId,
		&sessionStart,
		&sessionUpdate,
		&s.Active)
	if err != nil {
		return nil, err
	}
	s.SessionStart, err = time.Parse(db.DateLayout, string(sessionStart))
	if err != nil {
		return nil, err
	}
	s.SessionUpdate, err = time.Parse(db.DateLayout, string(sessionUpdate))
	if err != nil {
		return nil, err
	}
	return &s, nil
}

func (da *DataAccess) GetSessionBySessionID(sid string) (*database.Session, error) {
	s := database.Session{}
	row := da.Db.QueryRow(query.SelectSessionBySessionId, sid)
	err := row.Scan(
		&s.Id,
		&s.SessionId,
		&s.UserId,
		&s.SessionStart,
		&s.SessionUpdate,
		&s.Active)
	if err != nil {
		return nil, err
	}
	return &s, nil
}
