package query

const (
	InsertUserTag = `
	INSERT INTO user_tags (
		user_id,
		tagged_resource_id,
		resource_type
	)
	VALUES ($1, $2, $3)
	;`

	SelectUserTagsByPostId = `
	SELECT *
	FROM user_tags
	WHERE tagged_resource_id = $1 AND resource_type = 'post'
	;`

	SelectUserTagsByCommentId = `
	SELECT * 
	FROM user_tags
	WHERE tagged_resource_id = $1 AND resource_type = 'comment'
	;`

	DeleteUserTagById = `
	DELETE FROM user_tags
	WHERE id = $1
	;`

	InsertHashtag = `
	INSERT INTO hashtags (tag_name)
	VALUES ($1)
	ON CONFLICT (tag_name) DO NOTHING
	;`

	SelectHashtagByName = `
	SELECT *
	FROM hashtags
	WHERE tag_name = $1
	;`

	InsertHashtagResource = `
	INSERT INTO hashtag_resources (
		tag_id,
		tagged_resource_id,
		resource_type
	)
	VALUES ($1, $2, $3)
	;`

	SelectReferenceTagById = `
	SELECT *
	FROM hashtag_resources
	WHERE tag_id = $1
	;`
)
