package query

const (
	InsertPost = `
	INSERT INTO posts
	SET id = $1,
		author_id = $2,
		title = $3,
		content = $4,
		post_image = $5,
		is_public = $6,
	;`

	SelectFeed = `
	SELECT p.* 
	FROM posts p
	LEFT JOIN follows f 
		ON p.author_id = f.followed_id AND f.follower_id = $1
	WHERE 
		(p.is_public = TRUE AND p.is_active = 1) OR 
		(f.follower_id = $1 AND f.follow_status = 1 AND p.is_active = 1) OR 
		(p.author_id = $2 AND p.is_active = 1)
	ORDER BY 
		CASE 
			WHEN p.created_at >= NOW() - INTERVAL '24 hours' 
			THEN 1 
			ELSE 2 
		END
		ASC,
		p.rating DESC,
		p.created_at DESC
	;`

	SelectActivePosts = `
	SELECT * 
	FROM posts
	WHERE is_active = 1
	ORDER BY created_at DESC
	;`

	SelectUserActivePosts = `
	SELECT *
	FROM posts
	WHERE author_id = $1 AND is_active = 1
	ORDER BY created_at DESC
	;`

	SelectUserPublicPosts = `
	SELECT * 
	FROM posts
	WHERE author_id = $1 AND is_public = TRUE AND is_active = 1
	ORDER BY created_at DESC
	;`

	SelectPostByID = `
	SELECT * 
	FROM posts
	WHERE id = $1
	;`

	UpdatePost = `
	UPDATE posts
	SET title = $1,
		content = $2,
		is_public = $3
	WHERE id = $4
	;`

	SetPostAsInactive = `
	UPDATE posts
	SET is_active = 0
	WHERE id = $1
	;`

	RatePostUp = `
	INSERT INTO post_ratings
	SET target_id = $1,
		user_id = $2,
		rating_value = 1
	ON DUPLICATE KEY UPDATE rating_value = 
		CASE
			WHEN rating_value = 1
			THEN 0
			ELSE 1
		END
	;`

	RatePostDown = `
	INSERT INTO post_ratings
		SET target_id = $1,
		user_id = $2,
		rating_value = -1
	ON DUPLICATE KEY UPDATE rating_value = 
		CASE
			WHEN rating_value = -1
			THEN 0
			ELSE -1
		END
	;`

	SelectPostUserRating = `
	SELECT *
	FROM post_ratings
	WHERE post_id = $1 AND user_id = $2
	;`
)
