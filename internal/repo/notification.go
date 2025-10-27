package repo

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
)

func (db *dbHandler) CreateNotification(ctx context.Context, n *Notification) error {

	query := `
	INSERT INTO notifications (
		id,
		user_id,
		from_user_id,
		notif_type,
		notif_text,
		resource_id,
		parent_id
	)
	VALUES ($1, $2, $3, $4, $5, $6, $7)
	RETURNING
		id,
		user_id,
		from_user_id,
		notif_type,
		notif_text,
		resource_id,
		parent_id,
		is_read,
		created_at,
		updated_at
	;`
	ID := uuid.New()
	err := db.pool.QueryRow(ctx, query,
		ID,
		n.UserID,
		n.FromUserID,
		n.NotifType,
		n.NotifText,
		n.ResourceID,
		n.ParentID,
	).Scan(
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
	return err
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

func (db *dbHandler) DeleteFollowNotification(ctx context.Context, username, fromUsername string) error {
	query := `
	DELETE FROM notifications
	WHERE notif_type = 'follow_request'
		AND user_id = 
			(
				SELECT id 
				FROM users 
				WHERE username = $1
			)
		AND from_user_id = 
			(
				SELECT id 
				FROM users 
				WHERE username = $2
			)
	;`
	_, err := db.pool.Exec(ctx, query, username, fromUsername)
	return err
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

func (db *dbHandler) GetUserNotifs(ctx context.Context, username string, limit, offset int) ([]*NotificationWithUsers, error) {

	query := `
	SELECT 
		n.id,
		n.notif_type,
		n.notif_text,
		n.resource_id,
		n.parent_id,
		n.is_read,
		n.created_at,
		n.updated_at,
		u.id          AS user_id,
		u.username    AS user_username,
		u.profile_pic AS user_profile_pic,
		fu.id          AS from_user_id,
		fu.username    AS from_user_username,
		fu.profile_pic AS from_user_profile_pic
	FROM notifications n
	JOIN users u 
		ON n.user_id = u.id
	JOIN users fu 
		ON n.from_user_id = fu.id
	WHERE u.username = $1
	ORDER BY n.created_at DESC
	LIMIT $2 OFFSET $3
	;`

	notifs := []*NotificationWithUsers{}
	rows, err := db.pool.Query(ctx, query, username, limit, offset)
	if err != nil {
		if err == sql.ErrNoRows {
			return notifs, nil
		}
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		n := NotificationWithUsers{}
		err = rows.Scan(
			&n.Notification.ID,
			&n.Notification.NotifType,
			&n.Notification.NotifText,
			&n.Notification.ResourceID,
			&n.Notification.ParentID,
			&n.Notification.IsRead,
			&n.Notification.CreatedAt,
			&n.Notification.UpdatedAt,
			&n.User.ID,
			&n.User.Username,
			&n.User.ProfilePic,
			&n.FromUser.ID,
			&n.FromUser.Username,
			&n.FromUser.ProfilePic,
		)
		if err != nil {
			return nil, err
		}
		notifs = append(notifs, &n)
	}
	return notifs, nil
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
