package db

const (
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
)
