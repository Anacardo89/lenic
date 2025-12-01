package repo

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/require"

	"github.com/Anacardo89/lenic/pkg/testutils"
)

func TestCreateNotification(t *testing.T) {
	ctx := context.Background()

	// Seed DB
	err := SeedDB(ctx, TestDSN, SeedPath)
	require.NoError(t, err)

	// Connect test DB
	db, err := testutils.ConnectDB(ctx, TestDSN)
	require.NoError(t, err)

	userID := uuid.MustParse("f7a92b5b-4b7e-4787-9c0b-2b0b6cb86b4e")     // soccerpunk
	fromUserID := uuid.MustParse("a1f92e18-1d8f-4f0f-9a4d-3b9e3b26b7b1") // anacardo
	resourceID := "b1d3c0f7-5a1a-4f9b-9b2a-2a8e4f8b9f01"                 // post ID

	tests := []struct {
		name      string
		notif     *Notification
		expectErr bool
	}{
		{
			name: "Follow request notification",
			notif: &Notification{
				UserID:     userID,
				FromUserID: fromUserID,
				NotifType:  "follow_request",
				NotifText:  "anacardo wants to follow you.",
				ResourceID: resourceID,
				ParentID:   nil,
			},
		},
		{
			name: "Comment notification",
			notif: &Notification{
				UserID:     userID,
				FromUserID: fromUserID,
				NotifType:  "post_comment",
				NotifText:  "anacardo commented on your post.",
				ResourceID: resourceID,
				ParentID:   nil,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Repo.CreateNotification(ctx, tt.notif)

			if tt.expectErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			require.NotEqual(t, uuid.Nil, tt.notif.ID)
			require.Equal(t, userID, tt.notif.UserID)
			require.Equal(t, fromUserID, tt.notif.FromUserID)
			require.False(t, tt.notif.IsRead)

			// Verify in DB
			var notifFromDB Notification
			err = db.QueryRow(ctx, `
				SELECT
					id,
					user_id, 
					from_user_id, 
					notif_type, 
					notif_text, 
					resource_id, 
					parent_id, 
					is_read 
				FROM notifications 
				WHERE id = $1;
			`, tt.notif.ID).Scan(
				&notifFromDB.ID,
				&notifFromDB.UserID,
				&notifFromDB.FromUserID,
				&notifFromDB.NotifType,
				&notifFromDB.NotifText,
				&notifFromDB.ResourceID,
				&notifFromDB.ParentID,
				&notifFromDB.IsRead,
			)
			require.NoError(t, err)
			require.Equal(t, tt.notif.ID, notifFromDB.ID)
			require.Equal(t, tt.notif.NotifText, notifFromDB.NotifText)
		})
	}
}

func TestDeleteFollowNotification(t *testing.T) {
	ctx := context.Background()
	// Seed the DB
	err := SeedDB(ctx, TestDSN, SeedPath)
	require.NoError(t, err)

	// Connect test DB
	db, err := testutils.ConnectDB(ctx, TestDSN)
	require.NoError(t, err)

	tests := []struct {
		name          string
		username      string
		fromUsername  string
		expectDeleted bool
	}{
		{
			name:          "Delete existing follow request",
			username:      "soccerpunk", // recipient
			fromUsername:  "moderata",   // sender
			expectDeleted: true,
		},
		{
			name:          "No notification to delete",
			username:      "anacardo",
			fromUsername:  "soccerpunk",
			expectDeleted: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Repo.DeleteFollowNotification(ctx, tt.username, tt.fromUsername)
			require.NoError(t, err)

			var count int
			err = db.QueryRow(ctx, `
				SELECT COUNT(*) 
				FROM notifications
				WHERE notif_type = 'follow_request'
					AND user_id = (
						SELECT id 
						FROM users 
						WHERE username = $1
					)
					AND from_user_id = (
						SELECT id 
						FROM users 
						WHERE username = $2
					)
			;`, tt.username, tt.fromUsername).Scan(&count)
			require.NoError(t, err)

			if tt.expectDeleted {
				require.Equal(t, 0, count)
			} else {
				require.GreaterOrEqual(t, count, 0)
			}
		})
	}
}

func TestGetUserNotifs(t *testing.T) {
	ctx := context.Background()
	// Seed DB
	err := SeedDB(ctx, TestDSN, SeedPath)
	require.NoError(t, err)

	tests := []struct {
		name          string
		username      string
		limit         int
		offset        int
		expectedCount int
		expectErr     bool
	}{
		{
			name:          "Get first 2 notifications for soccerpunk",
			username:      "soccerpunk",
			limit:         2,
			offset:        0,
			expectedCount: 2,
		},
		{
			name:          "Get all notifications for anacardo",
			username:      "anacardo",
			limit:         10,
			offset:        0,
			expectedCount: 1,
		},
		{
			name:          "Offset returns empty slice",
			username:      "soccerpunk",
			limit:         10,
			offset:        10,
			expectedCount: 0,
		},
		{
			name:          "Non-existent username returns empty slice",
			username:      "nonexistent",
			limit:         5,
			offset:        0,
			expectedCount: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			notifs, err := Repo.GetUserNotifs(ctx, tt.username, tt.limit, tt.offset)

			if tt.expectErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.Len(t, notifs, tt.expectedCount)

			// Optional: verify first notification fields
			if len(notifs) > 0 {
				for _, n := range notifs {
					require.NotEmpty(t, n.Notification.ID)
					require.NotEmpty(t, n.Notification.NotifType)
					require.NotEmpty(t, n.User.Username)
					require.NotEmpty(t, n.FromUser.Username)
				}
			}
		})
	}
}

func TestUpdateNotificationRead(t *testing.T) {
	ctx := context.Background()

	// Seed DB
	err := SeedDB(ctx, TestDSN, SeedPath)
	require.NoError(t, err)

	db, err := testutils.ConnectDB(ctx, TestDSN)
	require.NoError(t, err)

	notifID := uuid.MustParse("c3c3d4e5-f6a7-4b8c-9d0e-3e4f5a6b7c8d") // Example existing notification

	tests := []struct {
		name      string
		notifID   uuid.UUID
		expectErr bool
	}{
		{
			name:    "Mark existing notification as read",
			notifID: notifID,
		},
		{
			name:      "Non-existing notification ID",
			notifID:   uuid.MustParse("00000000-0000-0000-0000-000000000000"),
			expectErr: false, // Update affects 0 rows but no error
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Repo.UpdateNotificationRead(ctx, tt.notifID)
			if tt.expectErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)

			// Verify in DB if notification exists
			var isRead bool
			err = db.QueryRow(ctx, `SELECT is_read FROM notifications WHERE id = $1`, tt.notifID).Scan(&isRead)
			if err != nil {
				if tt.notifID == uuid.MustParse("00000000-0000-0000-0000-000000000000") {
					require.Equal(t, pgx.ErrNoRows, err) // expected no row
					return
				}
				require.NoError(t, err)
			}
			require.True(t, isRead)
		})
	}
}
