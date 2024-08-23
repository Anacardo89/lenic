package query

const (
	InsertPost = `
	INSERT INTO posts
		SET post_guid=?,
			post_title=?,
			post_user=?,
			post_content=?,
			post_image=?,
			post_image_ext=?,
			post_active=?
	;`

	SelectActivePosts = `
	SELECT * FROM posts
		WHERE post_active=1
		ORDER BY created_at DESC
	;`

	SelectPostByGUID = `
	SELECT * FROM posts
		WHERE post_guid=?
	;`

	UpdatePostText = `
	UPDATE posts
		SET post_content=?
		WHERE post_guid=?
	;`

	SetPostAsInactive = `
	UPDATE posts
		SET post_active=0
		WHERE post_guid=?
	;`
)
