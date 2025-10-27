package repo

import (
	"context"
	"database/sql"
	"encoding/json"

	"github.com/google/uuid"

	"github.com/Anacardo89/lenic/internal/helpers"
)

// Conversations
func (db *dbHandler) CreateConversation(ctx context.Context, conv *Conversation) (uuid.UUID, error) {

	query := `
	INSERT INTO conversations (
		id,
		user1_id,
		user2_id
	)
	VALUES ($1, $2)
	RETURNING id
	;`

	ID := uuid.New()
	err := db.pool.QueryRow(ctx, query,
		ID,
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
	err := db.pool.QueryRow(ctx, query, conversationID, username).Scan(
		&c.ID,
		&c.User1ID,
		&c.User2ID,
		&c.UpdatedAt,
		&u.ID,
		&u.Username,
		&u.ProfilePic,
	)
	return &c, &u, err
}

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
	err := db.pool.QueryRow(ctx, query1, cID, user1, user2).Scan(
		&cID,
	)
	if err != nil {
		return nil, nil, err
	}
	err = db.pool.QueryRow(ctx, query2, cID).Scan(
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
	)
	return &c, users, err
}

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
		u.id,
		u.username,
		u.profile_pic,
		COALESCE(
			(
				SELECT json_agg(
					json_build_object(
						'id', m.id,
						'conversation_id', m.conversation_id,
						'sender_id', m.sender_id,
						'content', m.content,
						'is_read', m.is_read,
						'created_at', m.created_at
					) ORDER BY m.created_at DESC
				)
				FROM dmessages m
				WHERE m.conversation_id = c.id
				LIMIT 1000
			), '[]'::json
		) AS messages
	FROM users u
	JOIN conversations c
		ON u.id = CASE
			WHEN c.user1_id = $1 THEN c.user2_id
			ELSE c.user1_id
		END
	WHERE c.id = $2
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
	rows.Close()

	var completeConvos []*ConversationsWithDMs
	for _, c := range convos {
		rows, err = db.pool.Query(ctx, query3, u.ID, c.ID)
		if err != nil {
			return nil, nil, err
		}
		for rows.Next() {
			var dmJSON []byte
			var convo ConversationsWithDMs
			messages := make([]*DMessage, 0)
			var otherUser User
			convo.Messages = messages
			if err := rows.Scan(
				&otherUser.ID,
				&otherUser.Username,
				&otherUser.ProfilePic,
				&dmJSON,
			); err != nil {
				return nil, nil, err
			} else {
				convo.ID = c.ID
				convo.OtherUser = &otherUser
				if err := json.Unmarshal(dmJSON, &convo.Messages); err != nil {
					return nil, nil, err
				}
			}
			completeConvos = append(completeConvos, &convo)
		}
	}

	return &u, completeConvos, nil
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
		id,
		conversation_id,
		sender_id,
		content
	)
	VALUES ($1, $2, $3, $4)
	RETURNING id
	;`

	ID := uuid.New()
	err := db.pool.QueryRow(ctx, query,
		ID,
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
	ORDER BY m.created_at DESC
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
	return dms, nil
}

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

	_, err := db.pool.Exec(ctx, query, conversationID, username)
	return err
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
