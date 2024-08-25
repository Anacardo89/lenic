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
			rating=?,
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

	RatePostUp = `
	INSERT INTO post_ratings
		SET post_id=?,
		user_id=?,
		rating_value=1
		ON DUPLICATE KEY UPDATE rating_value = CASE
        	WHEN rating_value = 1
				THEN 0
        	ELSE 1
    	END
	;`

	RatePostDown = `
	INSERT INTO post_ratings
		SET post_id=?,
		user_id=?,
		rating_value=-1
		ON DUPLICATE KEY UPDATE rating_value = CASE
        	WHEN rating_value = -1
				THEN 0
        	ELSE -1
    	END
	;`

	SelectPostUserRating = `
	SELECT * FROM post_ratings
		WHERE post_id=? AND user_id=?
	;`
)
