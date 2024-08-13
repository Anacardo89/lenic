package orm

import (
	"database/sql"

	"github.com/Anacardo89/tpsi25_blog/internal/data/query"
	"github.com/Anacardo89/tpsi25_blog/internal/model"
)

var (
	Uda UserDataAccess
)

type UserDataAccess struct {
	db *sql.DB
}

func (uda *UserDataAccess) CreateUser(u *model.User) error {
	_, err := uda.db.Exec(
		query.InsertUser,
		u.UserName,
		u.UserEmail,
		u.UserPass,
		0)
	return err
}

func (uda *UserDataAccess) GetUserByID(id int) (*model.User, error) {
	u := model.User{}
	row := uda.db.QueryRow(query.SelectUserById, id)
	err := row.Scan(
		u.Id,
		u.UserName,
		u.UserEmail,
		u.CreatedAt,
		u.UpdatedAt,
		u.Active,
	)
	if err != nil {
		return nil, err
	}
	return &u, nil
}
