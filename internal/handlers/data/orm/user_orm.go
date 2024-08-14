package orm

import (
	"github.com/Anacardo89/tpsi25_blog/internal/handlers/data/query"
	"github.com/Anacardo89/tpsi25_blog/internal/model/database"
)

func (da *DataAccess) CreateUser(u *database.User) error {
	_, err := da.db.Exec(query.InsertUser,
		u.UserName,
		u.UserEmail,
		u.UserPass,
		u.Active)
	return err
}

func (da *DataAccess) GetUserByID(id int) (*database.User, error) {
	u := database.User{}
	row := da.db.QueryRow(query.SelectUserById, id)
	err := row.Scan(
		u.Id,
		u.UserName,
		u.UserEmail,
		u.CreatedAt,
		u.UpdatedAt,
		u.Active)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (da *DataAccess) GetUserByName(name string) (*database.User, error) {
	u := database.User{}
	row := da.db.QueryRow(query.SelectUserByName, name)
	err := row.Scan(
		u.Id,
		u.UserName,
		u.UserEmail,
		u.CreatedAt,
		u.UpdatedAt,
		u.Active)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (da *DataAccess) GetUserByEmail(email string) (*database.User, error) {
	u := database.User{}
	row := da.db.QueryRow(query.SelectUserByEmail, email)
	err := row.Scan(
		u.Id,
		u.UserName,
		u.UserEmail,
		u.CreatedAt,
		u.UpdatedAt,
		u.Active)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (da *DataAccess) SetUserAsActive(name string) error {
	_, err := da.db.Exec(query.UpdateUserActive, name)
	if err != nil {
		return err
	}
	return nil
}
