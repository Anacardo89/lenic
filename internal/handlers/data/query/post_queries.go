package query

const (
	InsertPost = `
	INSERT INTO posts
		SET post_guid=?,
			author_id=?,
			title=?,
			content=?,
			post_image=?,
			image_ext=?,
			is_public=?,
			vote_count=?,
			active=?
	;`

	SelectActivePosts = `
	SELECT * FROM posts
		WHERE active=1
		ORDER BY created_at DESC
	;`

	SelectPostByGUID = `
	SELECT * FROM posts
		WHERE post_guid=?
	;`

	UpdatePostText = `
	UPDATE posts
		SET content=?
		WHERE post_guid=?
	;`

	SetPostAsInactive = `
	UPDATE posts
		SET active=0
		WHERE post_guid=?
	;`
)
