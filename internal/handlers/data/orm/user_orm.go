package orm

import (
	"database/sql"
	"time"

	"github.com/Anacardo89/tpsi25_blog/internal/handlers/data/query"
	"github.com/Anacardo89/tpsi25_blog/internal/model/database"
	"github.com/Anacardo89/tpsi25_blog/pkg/db"
)

func (da *DataAccess) CreateUser(u *database.User) error {
	_, err := da.Db.Exec(query.InsertUser,
		u.UserName,
		u.Email,
		u.HashPass,
		u.Active)
	return err
}

func (da *DataAccess) GetUserByID(id int) (*database.User, error) {
	var (
		createdAt []byte
		updatedAt []byte
	)
	u := database.User{}
	row := da.Db.QueryRow(query.SelectUserById, id)
	err := row.Scan(
		&u.Id,
		&u.UserName,
		&u.Email,
		&u.HashPass,
		&u.ProfilePic,
		&u.ProfilePicExt,
		&u.Followers,
		&u.Following,
		&createdAt,
		&updatedAt,
		&u.Active)
	if err != nil {
		return nil, err
	}
	u.CreatedAt, err = time.Parse(db.DateLayout, string(createdAt))
	if err != nil {
		return nil, err
	}
	u.UpdatedAt, err = time.Parse(db.DateLayout, string(updatedAt))
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
		&u.Email,
		&u.HashPass,
		&u.ProfilePic,
		&u.ProfilePicExt,
		&u.Followers,
		&u.Following,
		&createdAt,
		&updatedAt,
		&u.Active)
	if err != nil {
		return nil, err
	}
	u.CreatedAt, err = time.Parse(db.DateLayout, string(createdAt))
	if err != nil {
		return nil, err
	}
	u.UpdatedAt, err = time.Parse(db.DateLayout, string(updatedAt))
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (da *DataAccess) GetUserByEmail(email string) (*database.User, error) {
	var (
		createdAt []byte
		updatedAt []byte
	)
	u := database.User{}
	row := da.Db.QueryRow(query.SelectUserByEmail, email)
	err := row.Scan(
		&u.Id,
		&u.UserName,
		&u.Email,
		&u.HashPass,
		&u.ProfilePic,
		&u.ProfilePicExt,
		&u.Followers,
		&u.Following,
		&createdAt,
		&updatedAt,
		&u.Active)
	if err != nil {
		return nil, err
	}
	u.CreatedAt, err = time.Parse(db.DateLayout, string(createdAt))
	if err != nil {
		return nil, err
	}
	u.UpdatedAt, err = time.Parse(db.DateLayout, string(updatedAt))
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

func (da *DataAccess) SetNewPassword(user string, pass string) error {
	_, err := da.Db.Exec(query.UpdatePassword, pass, user)
	if err != nil {
		return err
	}
	return nil
}

func (da *DataAccess) GetUserFollows(follower_id int, followed_id int) (*database.Follows, error) {
	f := database.Follows{}
	row := da.Db.QueryRow(query.SelectUserFollows, follower_id, followed_id)
	err := row.Scan(
		&f.FollowerId,
		&f.FollowedId)
	if err != nil {
		return nil, err
	}
	return &f, nil
}

func (da *DataAccess) GetFollowers(followed_id int) (*[]database.Follows, error) {
	follows := []database.Follows{}
	rows, err := da.Db.Query(query.SelectUserFollowers, followed_id)
	if err != nil {
		if err == sql.ErrNoRows {
			return &follows, nil
		}
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		f := database.Follows{}
		err = rows.Scan(
			&f.FollowerId,
			&f.FollowedId,
		)
		if err != nil {
			return nil, err
		}
		follows = append(follows, f)
	}
	return &follows, nil
}

func (da *DataAccess) GetFollowing(follower_id int) (*[]database.Follows, error) {
	follows := []database.Follows{}
	rows, err := da.Db.Query(query.SelectUserFollowing, follower_id)
	if err != nil {
		if err == sql.ErrNoRows {
			return &follows, nil
		}
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		f := database.Follows{}
		err = rows.Scan(
			&f.FollowerId,
			&f.FollowedId,
		)
		if err != nil {
			return nil, err
		}
		follows = append(follows, f)
	}
	return &follows, nil
}

func (da *DataAccess) FollowUser(follower_id int, followed_id int) error {
	_, err := da.Db.Exec(query.FollowUser, follower_id, followed_id)
	if err != nil {
		return err
	}
	return nil
}

func (da *DataAccess) UnfollowUser(follower_id int, followed_id int) error {
	_, err := da.Db.Exec(query.UnfollowUser, follower_id, followed_id)
	if err != nil {
		return err
	}
	return nil
}

func (da *DataAccess) UpdateProfilePic(profile_pic string, profile_pic_ext string, username string) error {
	_, err := da.Db.Exec(query.UpdateProfilePic, profile_pic, profile_pic_ext, username)
	if err != nil {
		return err
	}
	return nil
}
