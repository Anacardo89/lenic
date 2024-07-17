package db

const (
	SelectUserById = `
	SELECT user_name FROM users
		WHERE id=?
	;`

	SelectUserByName = `
	SELECT * FROM users
		WHERE user_name = ?
	;`

	SelectUserByEmail = `
	SELECT * FROM users
		WHERE user_email = ?
	;`

	InsertUser = `
	INSERT INTO users
		SET user_name=?,
			user_email=?,
			user_password=?,
			user_active=?
	;`

	UpdateUserActive = `
	UPDATE users
		SET user_active=1
		WHERE user_name=?
	;`

	SelectLogin = `
	SELECT id, user_name, user_password, user_active FROM users
		WHERE user_name=?
	;`

	InsertSession = `
	INSERT INTO sessions
		SET session_active=?,
			session_id=?,
			user_id=?,
			session_update=?
		ON DUPLICATE KEY UPDATE user_id=?, session_update=?
	;`

	SelectUserFromSessions = `
	SELECT user_id FROM sessions
		WHERE session_id=?
	;`

	InsertPost = `
	INSERT INTO posts
		SET post_guid=?,
			post_title=?,
			post_user=?,
			post_content=?,
			post_active=?
	;`

	SelectPosts = `
	SELECT post_guid, post_title, post_user, post_content, created_at
		FROM posts
		ORDER BY created_at DESC
	;`

	SelectPostByGUID = `
	SELECT post_title, post_user, post_content, created_at
		FROM posts
		WHERE post_guid=?
	;`

	SelectComments = `
	SELECT id, comment_user, comment_text, created_at
		FROM comments
		WHERE post_guid=?
	;`

	InsertComment = `
	INSERT INTO comments
		SET post_guid=?,
			comment_user=?,
			comment_text=?,
			comment_active=?
	;`

	UpdateComment = `
	UPDATE comments
		SET comment_text=?
		WHERE id=?
	;`
)
