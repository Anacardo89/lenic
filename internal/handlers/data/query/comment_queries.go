package query

const (
	RateCommentUp = `
	INSERT INTO comment_ratings (
		target_id,
		user_id,
		rating_value = 1
	)
	VALUES ($1, $2)
	ON DUPLICATE KEY UPDATE rating_value = 
		CASE
			WHEN rating_value = 1
			THEN 0
			ELSE 1
		END
	;`

	RateCommentDown = `
	INSERT INTO comment_ratings (
		target_id,
		user_id,
		rating_value = -1
	)
	VALUES ($1, $2)
	ON DUPLICATE KEY UPDATE rating_value = 
		CASE
			WHEN rating_value = -1
			THEN 0
			ELSE 1
		END
	;`
)
