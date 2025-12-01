package repo

import (
	"context"

	"github.com/google/uuid"

	"github.com/Anacardo89/lenic/internal/helpers"
)

// Conversations

// Endpoints:
//
// POST /action/user/{user_encoded}/conversations/{conversation_id}/dms
func (db *dbHandler) GetConversationAndSender(ctx context.Context, conversationID uuid.UUID, username string) (*Conversation, *User, error) {
	query := `
		SELECT 
		c.id AS conversation_id,
		c.user1_id,
		c.user2_id,
		c.updated_at,
		u.id AS user_id,
		u.username,
		u.profile_pic
	FROM conversations c
	JOIN users u
		ON u.username = $2
	WHERE c.id = $1
		AND (u.id = c.user1_id OR u.id = c.user2_id);
	;`
	var c Conversation
	var u User
	if err := db.pool.QueryRow(ctx, query, conversationID, username).Scan(
		&c.ID,
		&c.User1ID,
		&c.User2ID,
		&c.UpdatedAt,
		&u.ID,
		&u.Username,
		&u.ProfilePic,
	); err != nil {
		return nil, nil, err
	}
	return &c, &u, nil
}

// Endpoints:
//
// POST /action/user/{user_encoded}/conversations
func (db *dbHandler) GetConversationAndUsers(ctx context.Context, user1, user2 string) (*Conversation, []*User, error) {
	query1 := `
	WITH ids AS (
		SELECT 
			u1.id AS u1_id,
			u2.id AS u2_id
		FROM users u1, users u2
		WHERE u1.username = $2 AND u2.username = $3
	)
	INSERT INTO conversations(
		id,
		user1_id,
		user2_id
	)
	SELECT
		$1,
		u1_id,
		u2_id
	FROM ids
	ON CONFLICT (user1_id, user2_id)
		DO UPDATE SET user1_id = conversations.user1_id
	RETURNING id
	;`
	query2 := `
	SELECT c.id, c.user1_id, c.user2_id, c.updated_at,
		u1.id, u1.username, u1.profile_pic,
		u2.id, u2.username, u2.profile_pic
	FROM conversations c
	JOIN users u1 ON u1.id = c.user1_id
	JOIN users u2 ON u2.id = c.user2_id
	WHERE c.id = $1
	;`
	cID := uuid.New()
	var c Conversation
	users := []*User{
		new(User),
		new(User),
	}
	if err := db.pool.QueryRow(ctx, query1, cID, user1, user2).Scan(
		&cID,
	); err != nil {
		return nil, nil, err
	}
	if err := db.pool.QueryRow(ctx, query2, cID).Scan(
		&c.ID,
		&c.User1ID,
		&c.User2ID,
		&c.CreatedAt,
		&users[0].ID,
		&users[0].Username,
		&users[0].ProfilePic,
		&users[1].ID,
		&users[1].Username,
		&users[1].ProfilePic,
	); err != nil {
		return nil, nil, err
	}
	return &c, users, nil
}

// Endpoints:
//
// GET /action/user/{user_encoded}/conversations
func (db *dbHandler) GetConversationsAndOwner(ctx context.Context, user string, limit, offset int) (*User, []*ConversationsWithDMs, error) {
	query1 := `
	SELECT
		id,
		username,
		profile_pic
	FROM users
	WHERE username = $1
	;`
	query2 := `
	SELECT
		id,
		user1_id,
		user2_id,
		updated_at
	FROM conversations
	WHERE user1_id = $1 OR user2_id = $1
	ORDER BY updated_at DESC
	LIMIT $2
	OFFSET $3
	;`
	query3 := `
	SELECT 
		id, 
		username, 
		profile_pic 
	FROM users
	WHERE id = $1
	;`
	query4 := `
	SELECT
		id,
		conversation_id,
		sender_id,
		content,
		is_read,
		created_at
	FROM dmessages 
	WHERE conversation_id = $1 
	ORDER BY created_at DESC
	LIMIT 1000
	;`
	var u User
	if err := db.pool.QueryRow(ctx, query1, user).Scan(
		&u.ID,
		&u.Username,
		&u.ProfilePic,
	); err != nil {
		return nil, nil, err
	}
	rows, err := db.pool.Query(ctx, query2, u.ID, limit, offset)
	if err != nil {
		return nil, nil, err
	}
	var convos []*Conversation
	for rows.Next() {
		var c Conversation
		if err := rows.Scan(
			&c.ID,
			&c.User1ID,
			&c.User2ID,
			&c.UpdatedAt,
		); err != nil {
			return nil, nil, err
		}
		convos = append(convos, &c)
	}
	if err := rows.Err(); err != nil {
		return nil, nil, err
	}
	rows.Close()
	var completeConvos = make([]*ConversationsWithDMs, len(convos))
	for i, c := range convos {
		var otherUserID uuid.UUID
		if c.User1ID == u.ID {
			otherUserID = c.User2ID
		} else {
			otherUserID = c.User1ID
		}
		completeConvos[i] = &ConversationsWithDMs{}
		completeConvos[i].ID = c.ID
		completeConvos[i].UpdatedAt = c.UpdatedAt
		completeConvos[i].OtherUser = new(User)
		if err := db.pool.QueryRow(ctx, query3, otherUserID).Scan(
			&completeConvos[i].OtherUser.ID,
			&completeConvos[i].OtherUser.Username,
			&completeConvos[i].OtherUser.ProfilePic,
		); err != nil {
			return nil, nil, err
		}
		rows, err := db.pool.Query(ctx, query4, c.ID)
		if err != nil {
			return nil, nil, err
		}
		defer rows.Close()
		messages := []*DMessage{}
		for rows.Next() {
			message := DMessage{}
			rows.Scan(
				&message.ID,
				&message.ConversationID,
				&message.SenderID,
				&message.Content,
				&message.IsRead,
				&message.CreatedAt,
			)
			messages = append(messages, &message)

		}
		if err := rows.Err(); err != nil {
			return nil, nil, err
		}
		completeConvos[i].Messages = messages
	}
	return &u, completeConvos, nil
}

// Endpoints:
//
// ws - dm
func (db *dbHandler) GetConversationByUsers(ctx context.Context, user1ID, user2ID uuid.UUID) (*Conversation, error) {
	query := `
	SELECT
		id,
		user1_id,
		user2_id,
		created_at,
		updated_at
	FROM conversations
	WHERE user1_id = $1 AND user2_id = $2
		OR user1_id = $2 AND user2_id = $1
	;`
	min, max := helpers.OrderUUIDs(user1ID, user2ID)
	conv := Conversation{}
	if err := db.pool.QueryRow(ctx, query, min, max).
		Scan(
			&conv.ID,
			&conv.User1ID,
			&conv.User2ID,
			&conv.CreatedAt,
			&conv.UpdatedAt,
		); err != nil {
		return nil, err
	}
	return &conv, nil
}

// Endpoints:
//
// ws - dm
func (db *dbHandler) UpdateConversation(ctx context.Context, ID uuid.UUID) error {
	query := `
	UPDATE conversations
	SET updated_at = NOW()
	WHERE id = $1
	;`
	if _, err := db.pool.Exec(ctx, query, ID); err != nil {
		return err
	}
	return nil
}

// DMs

// Endpoints:
//
// POST /action/user/{user_encoded}/conversations/{conversation_id}/dms
func (db *dbHandler) CreateDM(ctx context.Context, dm *DMessage) (uuid.UUID, error) {
	query := `
	INSERT INTO dmessages (
		id,
		conversation_id,
		sender_id,
		content
	)
	VALUES ($1, $2, $3, $4)
	RETURNING id
	;`
	ID := uuid.New()
	if err := db.pool.QueryRow(ctx, query,
		ID,
		dm.ConversationID,
		dm.SenderID,
		dm.Content,
	).Scan(&ID); err != nil {
		return uuid.Nil, err
	}
	return ID, nil
}

// Endpoints:
//
// GET /action/user/{user_encoded}/conversations/{conversation_id}/dms
func (db *dbHandler) GetDMsByConversation(ctx context.Context, conversationID uuid.UUID, limit, offset int) ([]*DMessageWithUser, error) {
	query := `
	SELECT 
		m.id,
		m.conversation_id,
		m.content,
		m.is_read,
		m.created_at,
		u.id,
		u.username,
		u.profile_pic
	FROM dmessages m
	JOIN users u
		ON u.id = m.sender_id
	WHERE m.conversation_id = $1
	ORDER BY m.created_at ASC
	LIMIT $2
	OFFSET $3
	;`
	var dms []*DMessageWithUser
	rows, err := db.pool.Query(ctx, query, conversationID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		dm := DMessageWithUser{
			Sender: &User{},
		}
		err = rows.Scan(
			&dm.ID,
			&dm.ConversationID,
			&dm.Content,
			&dm.IsRead,
			&dm.CreatedAt,
			&dm.Sender.ID,
			&dm.Sender.Username,
			&dm.Sender.ProfilePic,
		)
		if err != nil {
			return nil, err
		}
		dms = append(dms, &dm)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return dms, nil
}

// Endpoints:
//
// PUT /action/user/{user_encoded}/conversations/{conversation_id}/read
func (db *dbHandler) ReadAllReceivedDMsInConvo(ctx context.Context, conversationID uuid.UUID, username string) error {
	query := `
		UPDATE dmessages m
		SET is_read = TRUE
		FROM users u
		WHERE u.username = $2
			AND m.conversation_id = $1
			AND m.sender_id != u.id
			AND m.is_read = FALSE;
	;`
	if _, err := db.pool.Exec(ctx, query, conversationID, username); err != nil {
		return err
	}
	return nil
}
