package query

const (
	InsertUser = `
	INSERT INTO users
		SET user_name=?,
			user_email=?,
			user_password=?,
			user_active=?
	;`

	SelectUserById = `
	SELECT * FROM users
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

	UpdateUserActive = `
	UPDATE users
		SET user_active=1
		WHERE user_name=?
	;`
)
