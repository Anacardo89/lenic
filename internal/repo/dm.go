package repo

import (
	"context"
	"database/sql"
	"encoding/json"

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

func (db *dbHandler) GetConversationAndSender(ctx context.Context, conversationID uuid.UUID, username string) (*Conversation, *User, error) {
	query := `
		SELECT 
		c.id AS conversation_id,
		c.user1_id,
		c.user2_id,
		c.created_at,
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
		&c.CreatedAt,
		&u.ID,
		&u.UserName,
		&u.ProfilePic,
	)
	return &c, &u, err
}

func (db *dbHandler) GetConversationAndUsers(ctx context.Context, user1, user2 string) (*Conversation, []*User, error) {
	query := `
	WITH u AS (
		SELECT id
		FROM users
		WHERE username = ANY($1::text[])
	),
	sorted AS (
		SELECT
			u1.id AS user1_id,
			u2.id AS user2_id
		FROM u u1
		CROSS JOIN u u2
		WHERE u1.id < u2.id
	),
	ins AS (
		INSERT INTO conversations (
			user1_id,
			user2_id
		)
		SELECT 
			user1_id,
			user2_id
		FROM sorted
		ON CONFLICT (user1_id, user2_id)
		DO UPDATE SET user1_id = EXCLUDED.user1_id
		RETURNING id
	)
	SELECT
		c.id,
		c.user1_id,
		c.user2_id,
		c.created_at,
		u1.id,
		u1.username,
		u1.profile_pic,
		u2.id,
		u2.username,
		u2.profile_pic
	FROM sorted
	LEFT JOIN ins
		ON true
	LEFT JOIN conversations c 
		ON c.id = ins.id
	JOIN users u1
		ON u1.id = sorted.user1_id
	JOIN users u2
		ON u2.id = sorted.user2_id
	;`

	usersIn := []string{user1, user2}
	var c Conversation
	users := make([]*User, 2)
	err := db.pool.QueryRow(ctx, query, usersIn).Scan(
		&c.ID,
		&c.User1ID,
		&c.User2ID,
		&c.CreatedAt,
		&users[0].ID,
		&users[0].UserName,
		&users[0].ProfilePic,
		&users[1].ID,
		&users[1].UserName,
		&users[1].ProfilePic,
	)
	return &c, users, err
}

func (db *dbHandler) GetConversationsAndOwner(ctx context.Context, user string, limit, offset int) (*User, []*ConversationsWithDMs, error) {

	query := `
	WITH u AS (
		SELECT
			id,
			username,
			profile_pic
		FROM users
		WHERE username = $1
	),
	user_convos AS (
		SELECT *
		FROM conversations c
		WHERE c.user1_id = (SELECT id FROM u)
			OR c.user2_id = (SELECT id FROM u)
		ORDER BY c.created_at DESC
		LIMIT $2
		OFFSET $3
	),
	convos_with_other AS (
		SELECT
			c.id,
			CASE
				WHEN c.user1_id = u.id
				THEN c.user2_id
				ELSE c.user1_id
			END AS other_user_id,
			c.created_at
		FROM user_convos c
		CROSS JOIN u
	)
	SELECT
		u.id,
		u.username,
		u.profile_pic,
		COALESCE(
			json_agg(
				json_build_object(
					'id', c.convo_id,
					'created_at', c.created_at,
					'other_user', json_build_object(
						'id', ou.id,
						'username', ou.username,
						'profile_pic', ou.profile_pic
					),
					'messages', COALESCE(
						(
							SELECT json_agg(json_build_object(
								'id', m.id,
								'conversation_id', m.conversation_id,
								'sender_id', m.sender_id,
								'content', m.content,
								'is_read', m.is_read,
								'created_at', m.created_at
							) ORDER BY m.created_at DESC)
							FROM (
								SELECT *
								FROM dmessages
								WHERE conversation_id = c.convo_id
								ORDER BY created_at DESC
								LIMIT 1000
							) m
						),
						'[]'
					)
				)
			) FILTER (WHERE c.convo_id IS NOT NULL),
			'[]'
		) AS conversations
	FROM u
	LEFT JOIN convos_with_other c
		ON true
	LEFT JOIN users ou
		ON ou.id = c.other_user_id
	GROUP BY u.id;
	;`

	rows, err := db.pool.Query(ctx, query, user, limit, offset)
	if err != nil {
		return nil, nil, err
	}
	defer rows.Close()
	var u User
	var convos []*ConversationsWithDMs
	var convosJSON []byte
	if rows.Next() {
		if err := rows.Scan(
			&u.ID,
			&u.UserName,
			&u.ProfilePic,
			&convosJSON,
		); err != nil {
			return nil, nil, err
		}
		if err := json.Unmarshal(convosJSON, &convos); err != nil {
			return &u, nil, err
		}
	} else {
		return nil, nil, sql.ErrNoRows
	}
	return &u, convos, nil
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

func (db *dbHandler) GetDMsByConversation(ctx context.Context, conersationID uuid.UUID, limit, offset int) ([]*DMessageWithUser, error) {

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
	rows, err := db.pool.Query(ctx, query, conersationID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		dm := DMessageWithUser{}
		err = rows.Scan(
			&dm.ID,
			&dm.ConversationID,
			&dm.Content,
			&dm.IsRead,
			&dm.CreatedAt,
			&dm.Sender.ID,
			&dm.Sender.UserName,
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
