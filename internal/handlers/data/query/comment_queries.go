package query

const (
	InsertComment = `
	INSERT INTO comments
		SET post_guid=?,
			comment_author=?,
			comment_text=?,
			comment_active=?
	;`

	SelectCommentsByPost = `
	SELECT * FROM comments
		WHERE post_guid=?
	;`

	UpdateCommentText = `
	UPDATE comments
		SET comment_text=?
		WHERE id=?
	;`
)
