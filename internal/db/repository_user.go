package db

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
)

// User
func (c *dbClient) CreateUser(ctx context.Context, u *User) (uuid.UUID, error) {
	query := `
	INSERT INTO users (
		username, 
		display_name, 
		email, 
		password_hash
	)
	VALUES ($1, $2, $3, $4)
	RETURNING id
	;`

	var ID uuid.UUID
	err := c.Pool().QueryRow(ctx, query,
		u.UserName,
		u.DisplayName,
		u.Email,
		u.PasswordHash,
	).Scan(&ID)
	return ID, err
}

func (c *dbClient) GetUserByID(ctx context.Context, ID uuid.UUID) (*User, error) {

	query := `
	SELECT *
	FROM users
	WHERE id = $1
	;`

	u := User{}
	row := c.Pool().QueryRow(ctx, query, ID)
	err := row.Scan(
		&u.ID,
		&u.UserName,
		&u.DisplayName,
		&u.Email,
		&u.PasswordHash,
		&u.ProfilePic,
		&u.Bio,
		&u.Followers,
		&u.Following,
		&u.IsActive,
		&u.IsVerified,
		&u.UserRole,
		&u.CreatedAt,
		&u.UpdatedAt,
		&u.DeletedAt,
	)
	return &u, err
}

func (c *dbClient) GetUserByUserName(ctx context.Context, userName string) (*User, error) {

	query := `
	SELECT * 
	FROM users
	WHERE username = $1
	;`

	u := User{}
	row := c.Pool().QueryRow(ctx, query, userName)
	err := row.Scan(
		&u.ID,
		&u.UserName,
		&u.DisplayName,
		&u.Email,
		&u.PasswordHash,
		&u.ProfilePic,
		&u.Bio,
		&u.Followers,
		&u.Following,
		&u.IsActive,
		&u.IsVerified,
		&u.UserRole,
		&u.CreatedAt,
		&u.UpdatedAt,
		&u.DeletedAt,
	)
	return &u, err
}

func (c *dbClient) GetUserByEmail(ctx context.Context, email string) (*User, error) {

	query := `
	SELECT * 
	FROM users
	WHERE email = $1
	;`

	u := User{}
	row := c.Pool().QueryRow(ctx, query, email)
	err := row.Scan(
		&u.ID,
		&u.UserName,
		&u.DisplayName,
		&u.Email,
		&u.PasswordHash,
		&u.ProfilePic,
		&u.Bio,
		&u.Followers,
		&u.Following,
		&u.IsActive,
		&u.IsVerified,
		&u.UserRole,
		&u.CreatedAt,
		&u.UpdatedAt,
		&u.DeletedAt,
	)
	return &u, err
}

func (c *dbClient) SearchUsersByDisplayName(ctx context.Context, displayName string) ([]*User, error) {

	query := `
	SELECT * 
	FROM users
	WHERE display_name LIKE $1
	;`

	likeUser := "%" + displayName + "%"
	users := []*User{}
	rows, err := c.Pool().Query(ctx, query, likeUser)
	if err != nil {
		if err == sql.ErrNoRows {
			return users, nil
		}
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		u := User{}
		err = rows.Scan(
			&u.ID,
			&u.UserName,
			&u.DisplayName,
			&u.Email,
			&u.PasswordHash,
			&u.ProfilePic,
			&u.Bio,
			&u.Followers,
			&u.Following,
			&u.IsActive,
			&u.IsVerified,
			&u.UserRole,
			&u.CreatedAt,
			&u.UpdatedAt,
			&u.DeletedAt,
		)
		if err != nil {
			return nil, err
		}
		users = append(users, &u)
	}
	return users, nil
}

func (c *dbClient) SetUserActive(ctx context.Context, userName string) error {

	query := `
	UPDATE users
	SET is_active = TRUE
	WHERE username = $1
	;`

	_, err := c.Pool().Exec(ctx, query, userName)
	return err
}

func (c *dbClient) SetNewPassword(ctx context.Context, userName string, pass string) error {

	query := `
	UPDATE users
	SET password_hash = $2
	WHERE username = $1
	;`

	_, err := c.Pool().Exec(ctx, query, userName, pass)
	return err
}

func (c *dbClient) UpdateProfilePic(ctx context.Context, userName string, profilePic string) error {

	query := `
	UPDATE users
	SET profile_pic = $2,
	WHERE username = $1
	;`

	_, err := c.Pool().Exec(ctx, query, userName, profilePic)
	return err
}

// Follow
func (c *dbClient) FollowUser(ctx context.Context, followerID int, followedID int) error {

	query := `
	INSERT INTO follows (
		follower_id,
		followed_id
	)
	VALUES ($1, $2)
	;`

	_, err := c.Pool().Exec(ctx, query, followerID, followedID)
	return err
}

func (c *dbClient) AcceptFollow(ctx context.Context, followerID int, followedID int) error {

	query := `
	UPDATE follows
	SET follow_status = 'accepted'
	WHERE follower_id = $1 AND followed_id = $2
	;`

	_, err := c.pool.Exec(ctx, query, followerID, followedID)
	return err

}

func (c *dbClient) UnfollowUser(ctx context.Context, followerID int, followedID int) error {

	query := `
	DELETE FROM follows
	WHERE follower_id = $1 AND followed_id = $2
	;`

	_, err := c.Pool().Exec(ctx, query, followerID, followedID)
	return err
}

func (c *dbClient) GetUserFollows(ctx context.Context, followerID int, followedID int) (*Follows, error) {

	query := `
	SELECT *
	FROM follows
	WHERE follower_id = $1 AND followed_id = $2
	;`

	f := Follows{}
	row := c.Pool().QueryRow(ctx, query, followerID, followedID)
	err := row.Scan(
		&f.FollowerID,
		&f.FollowedID,
		&f.FollowStatus,
		&f.CreatedAt,
		&f.UpdatedAt,
	)
	return &f, err
}

func (c *dbClient) GetFollowers(ctx context.Context, followedID int) ([]*Follows, error) {

	query := `
	SELECT *
	FROM follows
	WHERE followed_id = $1 AND follow_status = 'accepted'
	;`

	follows := []*Follows{}
	rows, err := c.Pool().Query(ctx, query, followedID)
	if err != nil {
		if err == sql.ErrNoRows {
			return follows, nil
		}
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		f := Follows{}
		err = rows.Scan(
			&f.FollowerID,
			&f.FollowedID,
			&f.FollowStatus,
			&f.CreatedAt,
			&f.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		follows = append(follows, &f)
	}
	return follows, nil
}

func (c *dbClient) GetFollowing(ctx context.Context, followerID int) ([]*Follows, error) {

	query := `
	SELECT *
	FROM follows
	WHERE follower_id = $1 AND follow_status = 'acceoted'
	;`

	follows := []*Follows{}
	rows, err := c.Pool().Query(ctx, query, followerID)
	if err != nil {
		if err == sql.ErrNoRows {
			return follows, nil
		}
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		f := Follows{}
		err = rows.Scan(
			&f.FollowerID,
			&f.FollowedID,
			&f.FollowStatus,
			&f.CreatedAt,
			&f.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		follows = append(follows, &f)
	}
	return follows, nil
}
