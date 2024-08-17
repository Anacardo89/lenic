package orm

import (
	"time"

	"github.com/Anacardo89/tpsi25_blog/internal/handlers/data/query"
	"github.com/Anacardo89/tpsi25_blog/internal/model/database"
	"github.com/Anacardo89/tpsi25_blog/pkg/logger"
)

func (da *DataAccess) CreateUser(u *database.User) error {
	_, err := da.Db.Exec(query.InsertUser,
		u.UserName,
		u.UserEmail,
		u.UserPass,
		u.Active)
	return err
}

func (da *DataAccess) GetUserByID(id int) (*database.User, error) {
	u := database.User{}
	row := da.Db.QueryRow(query.SelectUserById, id)
	err := row.Scan(
		&u.Id,
		&u.UserName,
		&u.UserEmail,
		&u.UserPass,
		&u.CreatedAt,
		&u.UpdatedAt,
		&u.Active)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (da *DataAccess) GetUserByName(name string) (*database.User, error) {
	var (
		createdAt []byte
		updatedAt []byte
	)

	u := database.User{}
	row := da.Db.QueryRow(query.SelectUserByName, name)
	err := row.Scan(
		&u.Id,
		&u.UserName,
		&u.UserEmail,
		&u.UserPass,
		&createdAt,
		&updatedAt,
		&u.Active)
	if err != nil {
		logger.Error.Println(err)
		return nil, err
	}
	u.CreatedAt, err = time.Parse("2006-01-02 15:04:05", string(createdAt))
	if err != nil {
		logger.Error.Println(err)
		return nil, err
	}
	u.UpdatedAt, err = time.Parse("2006-01-02 15:04:05", string(updatedAt))
	if err != nil {
		logger.Error.Println(err)
		return nil, err
	}
	return &u, nil
}

func (da *DataAccess) GetUserByEmail(email string) (*database.User, error) {
	u := database.User{}
	row := da.Db.QueryRow(query.SelectUserByEmail, email)
	err := row.Scan(
		&u.Id,
		&u.UserName,
		&u.UserEmail,
		&u.UserPass,
		&u.CreatedAt,
		&u.UpdatedAt,
		&u.Active)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (da *DataAccess) SetUserAsActive(name string) error {
	_, err := da.Db.Exec(query.UpdateUserActive, name)
	if err != nil {
		return err
	}
	return nil
}