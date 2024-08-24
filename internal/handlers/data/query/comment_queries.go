package query

const (
	InsertComment = `
	INSERT INTO comments
		SET post_guid=?,
			author_id=?,
			content=?,
			vote_count=?,
			active=?
	;`

	SelectActiveCommentsByPost = `
	SELECT * FROM comments
		WHERE post_guid=? AND active=1
	;`

	UpdateCommentText = `
	UPDATE comments
		SET content=?
		WHERE id=?
	;`

	SetCommentAsInactive = `
	UPDATE comments
		SET active=0
		WHERE id=?
	;`
)
