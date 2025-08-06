package query

const (
	InsertUser = `
	INSERT INTO users (
		id, 
		username, 
		display_name, 
		email, 
		password_hash
	)
	VALUES ($1, $2, $3, $4, $5)
	;`

	SelectUserById = `
	SELECT *
	FROM users
	WHERE id = $1
	;`

	SelectUserByName = `
	SELECT * 
	FROM users
	WHERE username = $1
	;`

	SelectSearchUsers = `
	SELECT * 
	FROM users
	WHERE username LIKE $1
	;`

	SelectUserByEmail = `
	SELECT * 
	FROM users
	WHERE email = $1
	;`

	UpdateUserActive = `
	UPDATE users
	SET is_active = 1
	WHERE username = $1
	;`

	UpdatePassword = `
	UPDATE users
	SET password_hash = $1
	WHERE username = $2
	;`

	UpdateProfilePic = `
	UPDATE users
	SET profile_pic = $1,
	WHERE username = $2
	;`

	SelectUserFollows = `
	SELECT *
	FROM follows
	WHERE follower_id = $1 AND followed_id = $2
	;`

	SelectUserFollowers = `
	SELECT *
	FROM follows
	WHERE followed_id = $1 AND follow_status = 1
	;`

	SelectUserFollowing = `
	SELECT *
	FROM follows
	WHERE follower_id = $1 AND follow_status = 1
	;`

	FollowUser = `
	INSERT INTO follows (
		follower_id,
		followed_id
	)
	VALUES ($1, $2)
	;`

	AcceptFollow = `
	UPDATE follows
	SET follow_status = 'accepted'
	WHERE follower_id = $1 AND followed_id = $2
	;`

	UnfollowUser = `
	DELETE FROM follows
	WHERE follower_id = $1 AND followed_id = $2
	;`
)
