package query

const (
	InsertUser = `
	INSERT INTO users
		SET username=?,
			email=?,
			hashpass=?,
			active=?
	;`

	SelectUserById = `
	SELECT * FROM users
		WHERE id=?
	;`

	SelectUserByName = `
	SELECT * FROM users
		WHERE username = ?
	;`

	SelectUserByEmail = `
	SELECT * FROM users
		WHERE email = ?
	;`

	UpdateUserActive = `
	UPDATE users
		SET active=1
		WHERE username=?
	;`

	UpdatePassword = `
	UPDATE users
		SET hashpass = ?
		WHERE username = ?
	;`

	SelectUserFollows = `
	SELECT * FROM follows
		WHERE follower_id=? AND followed_id=?
	;`

	FollowUser = `
	INSERT INTO follows
		SET follower_id=?,
			followed_id=?
	;`

	UnfollowUser = `
	DELETE FROM follows
		WHERE follower_id=? AND followed_id=?
	;`
)
