package repo

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
)

func (db *dbHandler) CreateNotification(ctx context.Context, n *Notification) (uuid.UUID, error) {

	query := `
	INSERT INTO notifications (
		user_id,
		from_user_id,
		notif_type,
		notif_text,
		resource_id,
		parent_id
	)
	VALUES ($1, $2, $3, $4, $5, $6)
	RETURNING id
	;`

	var ID uuid.UUID
	err := db.pool.QueryRow(ctx, query,
		n.UserID,
		n.FromUserID,
		n.NotifType,
		n.NotifText,
		n.ResourceID,
		n.ParentID,
	).Scan(&ID)
	return ID, err
}

func (db *dbHandler) GetFollowNotification(ctx context.Context, userID, fromUserID uuid.UUID) (*Notification, error) {

	query := `
	SELECT *
	FROM notifications
	WHERE
		notif_type = 'follow_request' AND
		user_id = $1 AND from_user_id = $2
	;`

	n := Notification{}
	err := db.pool.QueryRow(ctx, query, userID, fromUserID).
		Scan(
			&n.ID,
			&n.UserID,
			&n.FromUserID,
			&n.NotifType,
			&n.NotifText,
			&n.ResourceID,
			&n.ParentID,
			&n.IsRead,
			&n.CreatedAt,
			&n.UpdatedAt,
		)
	return &n, err
}

func (db *dbHandler) GetNotification(ctx context.Context, ID uuid.UUID) (*Notification, error) {

	query := `
	SELECT *
	FROM notifications
	WHERE id = $1
	;`

	n := Notification{}
	err := db.pool.QueryRow(ctx, query, ID).
		Scan(
			&n.ID,
			&n.UserID,
			&n.FromUserID,
			&n.NotifType,
			&n.NotifText,
			&n.ResourceID,
			&n.ParentID,
			&n.IsRead,
			&n.CreatedAt,
			&n.UpdatedAt,
		)
	return &n, err
}

func (db *dbHandler) GetNotificationsByUser(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*Notification, error) {

	query := `
	SELECT *
	FROM notifications
	WHERE user_id = $1
	ORDER BY created_at DESC
	LIMIT $2
	OFFSET $3
	;`

	notifs := []*Notification{}
	rows, err := db.pool.Query(ctx, query, userID, limit, offset)
	if err != nil {
		if err == sql.ErrNoRows {
			return notifs, nil
		}
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		n := Notification{}
		err = rows.Scan(
			&n.ID,
			&n.UserID,
			&n.FromUserID,
			&n.NotifType,
			&n.NotifText,
			&n.ResourceID,
			&n.ParentID,
			&n.IsRead,
			&n.CreatedAt,
			&n.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		notifs = append(notifs, &n)
	}
	return notifs, nil
}

func (db *dbHandler) UpdateNotificationRead(ctx context.Context, ID uuid.UUID) error {

	query := `
	UPDATE notifications
	SET is_read = TRUE
	WHERE id = $1
	;`

	_, err := db.pool.Exec(ctx, query, ID)
	return err
}

func (db *dbHandler) DeleteNotification(ctx context.Context, ID uuid.UUID) error {

	query := `
	DELETE FROM notifications
	WHERE id = $1
	;`

	_, err := db.pool.Exec(ctx, query, ID)
	return err
}
