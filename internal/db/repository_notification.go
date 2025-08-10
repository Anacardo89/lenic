package db

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
)

func (c *dbClient) CreateNotification(ctx context.Context, n *Notification) (uuid.UUID, error) {

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
	err := c.Pool().QueryRow(ctx, query,
		n.UserID,
		n.FromUserID,
		n.NotifType,
		n.NotifText,
		n.ResourceID,
		n.ParentID,
	).Scan(&ID)
	return ID, err
}

func (c *dbClient) GetFollowNotification(ctx context.Context, userID uuid.UUID, fromUserID uuid.UUID) (*Notification, error) {

	query := `
	SELECT *
	FROM notifications
	WHERE
		notif_type = 'follow_request' AND
		user_id = $1 AND from_user_id = $2
	;`

	n := Notification{}
	row := c.Pool().QueryRow(ctx, query, userID, fromUserID)
	err := row.Scan(
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

func (c *dbClient) GetNotificationById(ctx context.Context, ID uuid.UUID) (*Notification, error) {

	query := `
	SELECT *
	FROM notifications
	WHERE id = $1
	;`

	n := Notification{}
	row := c.Pool().QueryRow(ctx, query, ID)
	err := row.Scan(
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

func (c *dbClient) GetNotificationsByUser(ctx context.Context, userID uuid.UUID, limit int, offset int) ([]*Notification, error) {

	query := `
	SELECT *
	FROM notifications
	WHERE user_id = $1
	ORDER BY created_at DESC
	LIMIT $2
	OFFSET $3
	;`

	notifs := []*Notification{}
	rows, err := c.Pool().Query(ctx, query, userID, limit, offset)
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

func (c *dbClient) UpdateNotificationRead(ctx context.Context, ID uuid.UUID) error {

	query := `
	UPDATE notifications
	SET is_read = TRUE
	WHERE id = $1
	;`

	_, err := c.Pool().Exec(ctx, query, ID)
	return err
}

func (c *dbClient) DeleteNotificationByID(ctx context.Context, ID uuid.UUID) error {

	query := `
	DELETE FROM notifications
	WHERE id = $1
	;`

	_, err := c.Pool().Exec(ctx, query, ID)
	return err
}
