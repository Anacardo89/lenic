package query

const (
	InsertConversation = `
	INSERT INTO conversations (
		user1_id,
		user2_id
	)
	VALUES ($1, $2)
	;`

	SelectConversationById = `
	SELECT * 
	FROM conversations
	WHERE id = $1
	;`

	SelectConversationByUserIds = `
	SELECT *
	FROM conversations
	WHERE user1_id = $1 AND user2_id = $2
	;`

	SelectConversationsByUserId = `
	SELECT *
	FROM conversations
	WHERE user1_id = $1 OR user2_id = $2
	ORDER BY updated_at DESC
	LIMIT $3
	OFFSET $4
	;`

	InsertDMessage = `
	INSERT INTO dmessages (
		conversation_id,
		sender_id,
		content
	)
	VALUES ($1, $2, $3)
	;`

	SelectDMById = `
	SELECT *
	FROM dmessages
	WHERE id = $1
	;`

	SelectLastDMBySenderInConversation = `
	SELECT *
	FROM dmessages
	WHERE conversation_id = $1 AND sender_id = $2
	ORDER BY created_at DESC
	LIMIT 1
	;`

	SelectDMsByConversationId = `
	SELECT *
	FROM dmessages
	WHERE conversation_id = $1
	ORDER BY created_at
	LIMIT $2
	OFFSET $3
	;`

	UpdateConversationById = `
	UPDATE conversations
	SET updated_at = NOW()
	WHERE id = $1
	;`

	UpdateDMReadById = `
	UPDATE dmessages
	SET is_read = TRUE
	WHERE id = $1
	;`
)
