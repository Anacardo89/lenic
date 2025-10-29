package repo

import (
	"context"
	"database/sql"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

var (
	ErrUserExists   = errors.New("user already exists")
	ErrUserNotFound = errors.New("user not found")
)

// User

// Endpoints:
//
// /action/register
func (db *dbHandler) CreateUser(ctx context.Context, u *User) (uuid.UUID, error) {
	query := `
	INSERT INTO users (
		id,
		username,
		email, 
		password_hash
	)
	VALUES ($1, $2, $3, $4)
	RETURNING id
	;`

	ID := uuid.New()
	if err := db.pool.QueryRow(ctx, query,
		ID,
		u.Username,
		u.Email,
		u.PasswordHash,
	).Scan(&ID); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" { // unique_violation
			return uuid.Nil, ErrUserExists
		}
		return uuid.Nil, err
	}
	return ID, nil
}

// Endpoints:
//
//	/user/{encoded_username}/feed
//	/user/{encoded_username}/followers
//	/user/{encoded_username}/following
//	/ws - rate_comment
//	/ws - rate_post
//	as well as CreateSession and ValidateSession
func (db *dbHandler) GetUserByID(ctx context.Context, ID uuid.UUID) (*User, error) {

	query := `
	SELECT
		id,
		username,
		email,
		password_hash,
		profile_pic,
		user_followers,
		user_following,
		is_active,
		is_verified,
		user_role,
		created_at,
		updated_at
	FROM users
	WHERE id = $1
	;`

	u := User{}
	err := db.pool.QueryRow(ctx, query, ID).
		Scan(
			&u.ID,
			&u.Username,
			&u.Email,
			&u.PasswordHash,
			&u.ProfilePic,
			&u.Followers,
			&u.Following,
			&u.IsActive,
			&u.IsVerified,
			&u.UserRole,
			&u.CreatedAt,
			&u.UpdatedAt,
		)
	return &u, err
}

// Endpoints:
//
// POST	/action/post/{post_id}/comment
// PUT /action/post/{post_id}/comment/{comment_id}
// DELETE /action/post/{post_id}/comment/{comment_id}
// /action/image
// /action/user/{user_encoded}/profile-pic
// /action/login
// /action/change-password
// POST /action/post
// PUT /action/post/{post_id}
// DELETE /action/post/{post_id}
// /user/{encoded_username}/followers
// /user/{encoded_username}/following
// /recover-password/{encoded_username}
// /user/{encoded_username}
// ws - comment_on_post
// ws - dm
// ws - follow_request
// ws - follow_accept
func (db *dbHandler) GetUserByUserName(ctx context.Context, userName string) (*User, error) {

	query := `
	SELECT
		id,
		username,
		email,
		password_hash,
		profile_pic,
		user_followers,
		user_following,
		is_active,
		is_verified,
		user_role,
		created_at,
		updated_at
	FROM users
	WHERE username = $1
	;`

	u := User{}
	err := db.pool.QueryRow(ctx, query, userName).Scan(
		&u.ID,
		&u.Username,
		&u.Email,
		&u.PasswordHash,
		&u.ProfilePic,
		&u.Followers,
		&u.Following,
		&u.IsActive,
		&u.IsVerified,
		&u.UserRole,
		&u.CreatedAt,
		&u.UpdatedAt,
	)
	return &u, err
}

// Endpoints:
//
// POST /action/forgot-password
func (db *dbHandler) GetUserByEmail(ctx context.Context, email string) (*User, error) {

	query := `
	SELECT
		id,
		username,
		email,
		password_hash,
		profile_pic,
		user_followers,
		user_following,
		is_active,
		is_verified,
		user_role,
		created_at,
		updated_at
	FROM users
	WHERE email = $1
	;`

	u := User{}
	err := db.pool.QueryRow(ctx, query, email).
		Scan(
			&u.ID,
			&u.Username,
			&u.Email,
			&u.PasswordHash,
			&u.ProfilePic,
			&u.Followers,
			&u.Following,
			&u.IsActive,
			&u.IsVerified,
			&u.UserRole,
			&u.CreatedAt,
			&u.UpdatedAt,
		)
	return &u, err
}

// Endpoints:
//
// GET /action/search/user
func (db *dbHandler) SearchUsersByUserName(ctx context.Context, username string) ([]*User, error) {

	query := `
	SELECT
		id,
		username,
		profile_pic
	FROM users
	WHERE username LIKE $1
	;`

	likeUser := "%" + username + "%"
	users := []*User{}
	rows, err := db.pool.Query(ctx, query, likeUser)
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
			&u.Username,
			&u.ProfilePic,
		)
		if err != nil {
			return nil, err
		}
		users = append(users, &u)
	}
	return users, nil
}

// TODO:
// 1- change username to displayName as seen name in frontend
// 2- implement search based on display name
func (db *dbHandler) SearchUsersByDisplayName(ctx context.Context, displayName string) ([]*User, error) {

	query := `
	SELECT
		id,
		username,
		email,
		password_hash,
		profile_pic,
		user_followers,
		user_following,
		is_active,
		is_verified,
		user_role,
		created_at,
		updated_at 
	FROM users
	WHERE display_name LIKE $1
	;`

	likeUser := "%" + displayName + "%"
	users := []*User{}
	rows, err := db.pool.Query(ctx, query, likeUser)
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
			&u.Username,
			&u.Email,
			&u.PasswordHash,
			&u.ProfilePic,
			&u.Followers,
			&u.Following,
			&u.IsActive,
			&u.IsVerified,
			&u.UserRole,
			&u.CreatedAt,
			&u.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		users = append(users, &u)
	}
	return users, nil
}

// Endpoints:
//
// /action/activate
//
// TODO: set is_active to be equivalent to deleted and is_verified to be the current is_active
func (db *dbHandler) SetUserActive(ctx context.Context, userName string) error {

	query := `
	UPDATE users
	SET is_active = TRUE
	WHERE username = $1
	;`

	_, err := db.pool.Exec(ctx, query, userName)
	return err
}

// Endpoints:
//
// /action/recover-password
// /action/change-password
func (db *dbHandler) SetNewPassword(ctx context.Context, userID uuid.UUID, passHash string) error {
	query := `
	UPDATE users
	SET password_hash = $2
	WHERE id = $1
	;`
	_, err := db.pool.Exec(ctx, query, userID, passHash)
	return err
}

// Endpoints:
//
// /action/user/{user_encoded}/profile-pic
func (db *dbHandler) UpdateProfilePic(ctx context.Context, username string, profilePic string) error {

	query := `
	UPDATE users
	SET profile_pic = $2
	WHERE username = $1
	;`

	_, err := db.pool.Exec(ctx, query, username, profilePic)
	return err
}

// Follow

// Endpoints:
//
// POST /action/user/{user_encoded}/follow
func (db *dbHandler) FollowUser(ctx context.Context, followerID uuid.UUID, followedUsername string) error {
	query := `
	INSERT INTO follows (
		follower_id,
		followed_id
	)
	VALUES (
		$1,
		(
			SELECT id
			FROM users 
			WHERE username = $2
		)
	)
	ON CONFLICT (follower_id, followed_id) DO NOTHING;
	;`
	_, err := db.pool.Exec(ctx, query, followerID, followedUsername)
	return err
}

// Endpoints:
//
// PUT /action/user/{user_encoded}/accept
func (db *dbHandler) AcceptFollow(ctx context.Context, followerName, followedName string) error {
	query := `
	UPDATE follows
	SET follow_status = 'accepted'
	WHERE follower_id = 
		(
			SELECT id
			FROM users 
			WHERE username = $1
		)
		AND followed_id = 
			(
			SELECT id
			FROM users 
			WHERE username = $2
		)
	;`
	_, err := db.pool.Exec(ctx, query, followerName, followedName)
	return err
}

// Endpoints:
//
// DELETE /action/user/{user_encoded}/unfollow
func (db *dbHandler) UnfollowUser(ctx context.Context, followerName, followedName string) error {
	query := `
	DELETE FROM follows
	WHERE follower_id = 
		(
			SELECT id 
			FROM users 
			WHERE username = $1
		)
		AND followed_id = 
			(
				SELECT id 
				FROM users 
				WHERE username = $2
			)
	;`
	_, err := db.pool.Exec(ctx, query, followerName, followedName)
	return err
}

// Endpoints:
//
// /user/{encoded_username}
func (db *dbHandler) GetUserFollows(ctx context.Context, followerID, followedID uuid.UUID) (*Follows, error) {

	query := `
	SELECT
		follower_id,
		followed_id,
		follow_status,
		created_at,
		updated_at
	FROM follows
	WHERE follower_id = $1 AND followed_id = $2
	;`

	f := Follows{}
	err := db.pool.QueryRow(ctx, query, followerID, followedID).
		Scan(
			&f.FollowerID,
			&f.FollowedID,
			&f.FollowStatus,
			&f.CreatedAt,
			&f.UpdatedAt,
		)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &f, nil
}

// Endpoints:
//
// /user/{encoded_username}/followers
func (db *dbHandler) GetFollowers(ctx context.Context, followedID uuid.UUID) ([]*Follows, error) {

	query := `
	SELECT *
	FROM follows
	WHERE followed_id = $1 AND follow_status = 'accepted'
	;`

	follows := []*Follows{}
	rows, err := db.pool.Query(ctx, query, followedID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
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

// Endpoints:
//
// /user/{encoded_username}/following
func (db *dbHandler) GetFollowing(ctx context.Context, followerID uuid.UUID) ([]*Follows, error) {

	query := `
	SELECT *
	FROM follows
	WHERE follower_id = $1 AND follow_status = 'acceoted'
	;`

	follows := []*Follows{}
	rows, err := db.pool.Query(ctx, query, followerID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
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
