package orm

import (
	"time"

	"github.com/Anacardo89/tpsi25_blog/internal/handlers/data/query"
	"github.com/Anacardo89/tpsi25_blog/internal/model/database"
	"github.com/Anacardo89/tpsi25_blog/pkg/db"
)

func (da *DataAccess) CreateSession(s *database.Session) error {
	_, err := da.Db.Exec(query.InsertSession,
		s.SessionId,
		s.UserId,
		s.Active,
		s.UserId)
	return err
}

func (da *DataAccess) GetSessionByID(id int) (*database.Session, error) {
	var (
		createdAt []byte
		updatedAt []byte
	)
	s := database.Session{}
	row := da.Db.QueryRow(query.SelectSessionById, id)
	err := row.Scan(
		&s.Id,
		&s.SessionId,
		&s.UserId,
		&createdAt,
		&updatedAt,
		&s.Active)
	if err != nil {
		return nil, err
	}
	s.CreatedAt, err = time.Parse(db.DateLayout, string(createdAt))
	if err != nil {
		return nil, err
	}
	s.UpdatedAt, err = time.Parse(db.DateLayout, string(updatedAt))
	if err != nil {
		return nil, err
	}
	return &s, nil
}

func (da *DataAccess) GetSessionBySessionID(sid string) (*database.Session, error) {
	var (
		createdAt []byte
		updatedAt []byte
	)
	s := database.Session{}
	row := da.Db.QueryRow(query.SelectSessionBySessionId, sid)
	err := row.Scan(
		&s.Id,
		&s.SessionId,
		&s.UserId,
		&createdAt,
		&updatedAt,
		&s.Active)
	if err != nil {
		return nil, err
	}
	s.CreatedAt, err = time.Parse(db.DateLayout, string(createdAt))
	if err != nil {
		return nil, err
	}
	s.UpdatedAt, err = time.Parse(db.DateLayout, string(updatedAt))
	if err != nil {
		return nil, err
	}
	return &s, nil
}
