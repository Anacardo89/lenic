package db

import (
	"context"
	"database/sql"

	"github.com/Anacardo89/tpsi25_blog/internal/helpers"
	"github.com/google/uuid"
)

// Conversations
func (c *dbClient) CreateConversation(ctx context.Context, conv *Conversation) (uuid.UUID, error) {

	query := `
	INSERT INTO conversations (
		user1_id,
		user2_id
	)
	VALUES ($1, $2)
	RETURNING id
	;`

	var ID uuid.UUID
	err := c.Pool().QueryRow(ctx, query,
		conv.User1ID,
		conv.User2ID,
	).Scan(&ID)
	return ID, err
}

func (c *dbClient) GetConversation(ctx context.Context, ID uuid.UUID) (*Conversation, error) {

	query := `
	SELECT * 
	FROM conversations
	WHERE id = $1
	;`

	conv := Conversation{}
	err := c.Pool().QueryRow(ctx, query, ID).
		Scan(
			&conv.ID,
			&conv.User1ID,
			&conv.User2ID,
			&conv.CreatedAt,
			&conv.UpdatedAt,
		)
	return &conv, err
}

func (c *dbClient) GetConversationByUsers(ctx context.Context, user1ID, user2ID uuid.UUID) (*Conversation, error) {

	query := `
	SELECT *
	FROM conversations
	WHERE user1_id = $1 AND user2_id = $2
	;`

	min, max := helpers.OrderUUIDs(user1ID, user2ID)
	conv := Conversation{}
	err := c.Pool().QueryRow(ctx, query, min, max).
		Scan(
			&conv.ID,
			&conv.User1ID,
			&conv.User2ID,
			&conv.CreatedAt,
			&conv.UpdatedAt,
		)
	return &conv, err
}

func (c *dbClient) GetConversationsByUser(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*Conversation, error) {

	query := `
	SELECT *
	FROM conversations
	WHERE user1_id = $1 OR user2_id = $1
	ORDER BY updated_at DESC
	LIMIT $3
	OFFSET $4
	;`

	convos := []*Conversation{}
	rows, err := c.Pool().Query(ctx, query, userID, limit, offset)
	if err != nil {
		if err == sql.ErrNoRows {
			return convos, nil
		}
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		convo := Conversation{}
		err = rows.Scan(
			&convo.ID,
			&convo.User1ID,
			&convo.User2ID,
			&convo.CreatedAt,
			&convo.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		convos = append(convos, &convo)
	}
	return convos, nil
}

func (c *dbClient) UpdateConversation(ctx context.Context, ID uuid.UUID) error {

	query := `
	UPDATE conversations
	SET updated_at = NOW()
	WHERE id = $1
	;`

	_, err := c.Pool().Exec(ctx, query, ID)
	return err
}

// DMs
func (c *dbClient) CreateDM(ctx context.Context, dm *DMessage) (uuid.UUID, error) {

	query := `
	INSERT INTO dmessages (
		conversation_id,
		sender_id,
		content
	)
	VALUES ($1, $2, $3)
	RETURNING id
	;`

	var ID uuid.UUID
	err := c.Pool().QueryRow(ctx, query,
		dm.ConversationID,
		dm.SenderID,
		dm.Content,
	).Scan(&ID)
	return ID, err
}

func (c *dbClient) GetDM(ctx context.Context, ID uuid.UUID) (*DMessage, error) {

	query := `
	SELECT *
	FROM dmessages
	WHERE id = $1
	;`

	dm := DMessage{}
	err := c.Pool().QueryRow(ctx, query, ID).
		Scan(
			&dm.ID,
			&dm.ConversationID,
			&dm.SenderID,
			&dm.Content,
			&dm.IsRead,
			&dm.CreatedAt,
		)
	return &dm, err
}

func (c *dbClient) GetConvoLastDMBySender(ctx context.Context, conversationID, senderID uuid.UUID) (*DMessage, error) {

	query := `
	SELECT *
	FROM dmessages
	WHERE conversation_id = $1 AND sender_id = $2
	ORDER BY created_at DESC
	LIMIT 1
	;`

	dm := DMessage{}
	err := c.Pool().QueryRow(ctx, query, conversationID, senderID).
		Scan(
			&dm.ID,
			&dm.ConversationID,
			&dm.SenderID,
			&dm.Content,
			&dm.IsRead,
			&dm.CreatedAt,
		)
	return &dm, err
}

func (c *dbClient) GetDMsByConversation(ctx context.Context, conersationID uuid.UUID, limit, offset int) ([]*DMessage, error) {

	query := `
	SELECT *
	FROM dmessages
	WHERE conversation_id = $1
	ORDER BY created_at
	LIMIT $2
	OFFSET $3
	;`

	dms := []*DMessage{}
	rows, err := c.Pool().Query(ctx, query, conersationID, limit, offset)
	if err != nil {
		if err == sql.ErrNoRows {
			return dms, nil
		}
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		dm := DMessage{}
		err = rows.Scan(
			&dm.ID,
			&dm.ConversationID,
			&dm.SenderID,
			&dm.Content,
			&dm.IsRead,
			&dm.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		dms = append(dms, &dm)
	}
	return dms, nil
}

func (c *dbClient) UpdateDMRead(ctx context.Context, ID uuid.UUID) error {

	query := `
	UPDATE dmessages
	SET is_read = TRUE
	WHERE id = $1
	;`

	_, err := c.Pool().Exec(ctx, query, ID)
	return err
}
