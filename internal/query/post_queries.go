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

	SelectPosts = `
	SELECT post_guid, post_title, post_user, post_content, created_at
		FROM posts
		ORDER BY created_at DESC
	;`

	SelectPostByGUID = `
	SELECT post_title, post_user, post_content, post_image, created_at
		FROM posts
		WHERE post_guid=?
	;`
)
