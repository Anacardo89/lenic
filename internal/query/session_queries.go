package query

const (
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
)
