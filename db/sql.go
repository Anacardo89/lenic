package db

const (
	SelectUserById = `
	SELECT user_name FROM users
		WHERE id=?
	`

	SelectUserByName = `
	SELECT * FROM users
		WHERE user_name = ?;
	`

	SelectUserByEmail = `
	SELECT * FROM users
		WHERE user_email = ?;
	`

	InsertUser = `
	INSERT INTO users
		SET user_name=?,
			user_email=?,
			user_password=?,
			user_active=?
	`

	SelectLogin = `
	SELECT id, user_name, user_password FROM users
		WHERE user_name=?
	`

	InsertSession = `
	"INSERT INTO sessions
		SET session_active=?,
			session_id=?,
			user_id=?,
			session_update=?
		ON DUPLICATE KEY UPDATE user_id=?, session_update=?"
	`

	SelectUserFromSessions = `
	SELECT user_id FROM sessions
		WHERE session_id=?
	`
)
