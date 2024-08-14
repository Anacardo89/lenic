package query

const (
	InsertSession = `
	INSERT INTO sessions
		SET session_id=?,
			user_id=?,
			session_active=?
		ON DUPLICATE KEY UPDATE user_id=?, session_update=?
	;`

	SelectSessionById = `
	SELECT * FROM sessions
		WHERE session_id=?
	;`

	SelectSessionBySessionId = `
	SELECT * FROM sessions
		WHERE session_id=?
	;`
)
