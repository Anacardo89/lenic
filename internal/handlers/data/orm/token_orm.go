package orm

import (
	"time"

	"github.com/Anacardo89/tpsi25_blog/internal/handlers/data/query"
	"github.com/Anacardo89/tpsi25_blog/internal/model/database"
	"github.com/Anacardo89/tpsi25_blog/pkg/db"
)

func (da *DataAccess) CreateToken(t *database.Token) error {
	_, err := da.Db.Exec(query.InsertToken,
		t.Token,
		t.UserId,
		t.Token,
	)
	return err
}

func (da *DataAccess) GetTokenByUserId(id int) (*database.Token, error) {
	var (
		createdAt []byte
		updatedAt []byte
	)
	t := database.Token{}
	row := da.Db.QueryRow(query.SelectTokenByUserId, id)
	err := row.Scan(
		&t.Id,
		&t.Token,
		&t.UserId,
		&createdAt,
		&updatedAt,
	)
	if err != nil {
		return nil, err
	}
	t.CreatedAt, err = time.Parse(db.DateLayout, string(createdAt))
	if err != nil {
		return nil, err
	}
	t.UpdatedAt, err = time.Parse(db.DateLayout, string(updatedAt))
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func (da *DataAccess) DeleteTokenByUserId(id int) error {
	_, err := da.Db.Exec(query.DeleteTokenByUserId, id)
	if err != nil {
		return err
	}
	return nil
}
