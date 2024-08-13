package query

const (
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
