package query

const (
	InsertConversation = `
	INSERT INTO conversations
		SET user1_id=?,
			user2_id=?
	;`

	InsertDMessage = `
	INSERT INTO dmessages
		SET conversation_id=?,
			sender_id=?,
			content=?
	;`

	SelectConversationsByUserId = `
	SELECT * FROM conversations
		WHERE user1_id=? OR user2_id=?
			ORDER BY updated_at DESC
			LIMIT ? OFFSET ?
	;`

	SelectDMsByConversationId = `
	SELECT * FROM dmessages
		WHERE conversation_id=?
			ORDER BY created_at DESC
			LIMIT ? OFFSET ?
	;`

	UpdateConversationById = `
	UPDATE conversation
		SET updated_at=CURRENT_TIMESTAMP,
		WHERE id=?
	;`
)
