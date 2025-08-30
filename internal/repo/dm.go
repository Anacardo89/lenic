package repo

import (
	"context"
	"database/sql"

	"github.com/Anacardo89/lenic/internal/helpers"
	"github.com/google/uuid"
)

// Conversations
func (db *dbHandler) CreateConversation(ctx context.Context, conv *Conversation) (uuid.UUID, error) {

	query := `
	INSERT INTO conversations (
		user1_id,
		user2_id
	)
	VALUES ($1, $2)
	RETURNING id
	;`

	var ID uuid.UUID
	err := db.pool.QueryRow(ctx, query,
		conv.User1ID,
		conv.User2ID,
	).Scan(&ID)
	return ID, err
}

func (db *dbHandler) GetConversation(ctx context.Context, ID uuid.UUID) (*Conversation, error) {

	query := `
	SELECT * 
	FROM conversations
	WHERE id = $1
	;`

	conv := Conversation{}
	err := db.pool.QueryRow(ctx, query, ID).
		Scan(
			&conv.ID,
			&conv.User1ID,
			&conv.User2ID,
			&conv.CreatedAt,
			&conv.UpdatedAt,
		)
	return &conv, err
}

func (db *dbHandler) GetConversationByUsers(ctx context.Context, user1ID, user2ID uuid.UUID) (*Conversation, error) {

	query := `
	SELECT *
	FROM conversations
	WHERE user1_id = $1 AND user2_id = $2
	;`

	min, max := helpers.OrderUUIDs(user1ID, user2ID)
	conv := Conversation{}
	err := db.pool.QueryRow(ctx, query, min, max).
		Scan(
			&conv.ID,
			&conv.User1ID,
			&conv.User2ID,
			&conv.CreatedAt,
			&conv.UpdatedAt,
		)
	return &conv, err
}

func (db *dbHandler) GetConversationsByUser(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*Conversation, error) {

	query := `
	SELECT *
	FROM conversations
	WHERE user1_id = $1 OR user2_id = $1
	ORDER BY updated_at DESC
	LIMIT $3
	OFFSET $4
	;`

	convos := []*Conversation{}
	rows, err := db.pool.Query(ctx, query, userID, limit, offset)
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

func (db *dbHandler) UpdateConversation(ctx context.Context, ID uuid.UUID) error {

	query := `
	UPDATE conversations
	SET updated_at = NOW()
	WHERE id = $1
	;`

	_, err := db.pool.Exec(ctx, query, ID)
	return err
}

// DMs
func (db *dbHandler) CreateDM(ctx context.Context, dm *DMessage) (uuid.UUID, error) {

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
	err := db.pool.QueryRow(ctx, query,
		dm.ConversationID,
		dm.SenderID,
		dm.Content,
	).Scan(&ID)
	return ID, err
}

func (db *dbHandler) GetDM(ctx context.Context, ID uuid.UUID) (*DMessage, error) {

	query := `
	SELECT *
	FROM dmessages
	WHERE id = $1
	;`

	dm := DMessage{}
	err := db.pool.QueryRow(ctx, query, ID).
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

func (db *dbHandler) GetConvoLastDMBySender(ctx context.Context, conversationID, senderID uuid.UUID) (*DMessage, error) {

	query := `
	SELECT *
	FROM dmessages
	WHERE conversation_id = $1 AND sender_id = $2
	ORDER BY created_at DESC
	LIMIT 1
	;`

	dm := DMessage{}
	err := db.pool.QueryRow(ctx, query, conversationID, senderID).
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

func (db *dbHandler) GetDMsByConversation(ctx context.Context, conersationID uuid.UUID, limit, offset int) ([]*DMessage, error) {

	query := `
	SELECT *
	FROM dmessages
	WHERE conversation_id = $1
	ORDER BY created_at
	LIMIT $2
	OFFSET $3
	;`

	dms := []*DMessage{}
	rows, err := db.pool.Query(ctx, query, conersationID, limit, offset)
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

func (db *dbHandler) UpdateDMRead(ctx context.Context, ID uuid.UUID) error {

	query := `
	UPDATE dmessages
	SET is_read = TRUE
	WHERE id = $1
	;`

	_, err := db.pool.Exec(ctx, query, ID)
	return err
}
