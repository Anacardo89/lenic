package query

const (
	InsertNotification = `
	INSERT INTO notifications
		SET user_id=?,
			from_user_id=?,
			notif_type=?,
			notif_message=?,
			resource_id=?
	;`

	SelectNotificationById = `
	SELECT *
		FROM notifications
		WHERE id=?
	;`

	SelectNotificationsByUser = `
	SELECT *
		FROM notifications
		WHERE id=?
	;`

	UpdateNotificationRead = `
	UPDATE notifications
		SET is_read=TRUE
		WHERE id=?
	;`
)
