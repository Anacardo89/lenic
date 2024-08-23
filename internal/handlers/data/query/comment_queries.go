package query

const (
	InsertComment = `
	INSERT INTO comments
		SET post_guid=?,
			comment_author=?,
			comment_text=?,
			comment_active=?
	;`

	SelectActiveCommentsByPost = `
	SELECT * FROM comments
		WHERE post_guid=? AND comment_active=1
	;`

	UpdateCommentText = `
	UPDATE comments
		SET comment_text=?
		WHERE id=?
	;`

	SetCommentAsInactive = `
	UPDATE comments
		SET comment_active=0
		WHERE id=?
	;`
)
