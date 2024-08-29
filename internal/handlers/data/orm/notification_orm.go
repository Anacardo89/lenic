package orm

import (
	"database/sql"
	"time"

	"github.com/Anacardo89/tpsi25_blog/internal/handlers/data/query"
	"github.com/Anacardo89/tpsi25_blog/internal/model/database"
	"github.com/Anacardo89/tpsi25_blog/pkg/db"
)

func (da *DataAccess) CreateNotification(n *database.Notification) (sql.Result, error) {
	result, err := da.Db.Exec(query.InsertNotification,
		n.UserID,
		n.FromUserId,
		n.NotifType,
		n.NotifMsg,
		n.ResourceId)
	return result, err
}

func (da *DataAccess) GetNotificationById(id int) (*database.Notification, error) {
	var (
		createdAt []byte
		updatedAt []byte
	)
	n := database.Notification{}
	row := da.Db.QueryRow(query.SelectNotificationById, id)
	err := row.Scan(
		&n.Id,
		&n.UserID,
		&n.FromUserId,
		&n.NotifType,
		&n.NotifMsg,
		&n.ResourceId,
		&n.IsRead,
		&createdAt,
		&updatedAt)
	if err != nil {
		return nil, err
	}
	n.CreatedAt, err = time.Parse(db.DateLayout, string(createdAt))
	if err != nil {
		return nil, err
	}
	n.UpdatedAt, err = time.Parse(db.DateLayout, string(updatedAt))
	if err != nil {
		return nil, err
	}

	return &n, nil
}

func (da *DataAccess) GetNotificationsByUser(user_id int, limit int, offset int) ([]*database.Notification, error) {
	notifs := []*database.Notification{}
	rows, err := da.Db.Query(query.SelectNotificationsByUser, user_id, limit, offset)
	if err != nil {
		if err == sql.ErrNoRows {
			return notifs, nil
		}
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var (
			createdAt []byte
			updatedAt []byte
		)
		n := database.Notification{}
		err = rows.Scan(
			&n.Id,
			&n.UserID,
			&n.FromUserId,
			&n.NotifType,
			&n.NotifMsg,
			&n.ResourceId,
			&n.IsRead,
			&createdAt,
			&updatedAt,
		)
		if err != nil {
			return nil, err
		}
		n.CreatedAt, err = time.Parse(db.DateLayout, string(createdAt))
		if err != nil {
			return nil, err
		}
		n.UpdatedAt, err = time.Parse(db.DateLayout, string(updatedAt))
		if err != nil {
			return nil, err
		}
		notifs = append(notifs, &n)
	}
	return notifs, nil
}

func (da *DataAccess) UpdateNotificationRead(id int) error {
	_, err := da.Db.Exec(query.UpdateNotificationRead, id)
	if err != nil {
		return err
	}
	return nil
}
