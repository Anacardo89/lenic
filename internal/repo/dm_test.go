package repo

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/require"

	"github.com/Anacardo89/lenic/pkg/testutils"
)

func TestGetConversationAndSender(t *testing.T) {
	ctx := context.Background()

	// Seed DB
	err := SeedDB(ctx, TestDSN, SeedPath)
	require.NoError(t, err)

	conversations := []struct {
		name          string
		convoID       uuid.UUID
		user1Username string
		user1ID       uuid.UUID
		user2Username string
		user2ID       uuid.UUID
		expectErr     bool
	}{
		{
			name:          "Retrieve conversation for user1",
			convoID:       uuid.MustParse("70f31f9b-632b-4c4a-bb3a-1f2c6013f001"),
			user1Username: "anacardo",
			user1ID:       uuid.MustParse("a1f92e18-1d8f-4f0f-9a4d-3b9e3b26b7b1"),
			user2Username: "moderata",
			user2ID:       uuid.MustParse("cfa53179-9085-4f33-86b3-5dc5f7a1465f"),
		},
		{
			name:          "Retrieve conversation for user2",
			convoID:       uuid.MustParse("70f31f9b-632b-4c4a-bb3a-1f2c6013f001"),
			user1Username: "moderata",
			user1ID:       uuid.MustParse("cfa53179-9085-4f33-86b3-5dc5f7a1465f"),
			user2Username: "anacardo",
			user2ID:       uuid.MustParse("a1f92e18-1d8f-4f0f-9a4d-3b9e3b26b7b1"),
		},
		{
			name:          "Non-participant username returns error",
			convoID:       uuid.MustParse("70f31f9b-632b-4c4a-bb3a-1f2c6013f001"),
			user1Username: "soccerpunk",
			user1ID:       uuid.MustParse("f7a92b5b-4b7e-4787-9c0b-2b0b6cb86b4e"),
			user2Username: "moderata",
			user2ID:       uuid.MustParse("cfa53179-9085-4f33-86b3-5dc5f7a1465f"),
			expectErr:     true,
		},
	}

	for _, tt := range conversations {
		t.Run(tt.name, func(t *testing.T) {

			// Call the method
			c, u, err := Repo.GetConversationAndSender(ctx, tt.convoID, tt.user1Username)
			if tt.expectErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)

			// Check conversation fields
			require.Equal(t, tt.convoID, c.ID)
			require.True(t, (c.User1ID == tt.user1ID && c.User2ID == tt.user2ID) || (c.User1ID == tt.user2ID && c.User2ID == tt.user1ID))

			// Check returned user
			require.Equal(t, tt.user1Username, u.Username)
		})
	}
}

func TestGetConversationAndUsers(t *testing.T) {
	ctx := context.Background()

	// Seed DB
	err := SeedDB(ctx, TestDSN, SeedPath)
	require.NoError(t, err)

	tests := []struct {
		name      string
		user1     string
		user2     string
		expectErr bool
	}{
		{
			name:  "Retrieve conversation between anacardo and moderata",
			user1: "anacardo",
			user2: "moderata",
		},
		{
			name:  "Retrieve conversation between moderata and anacardo",
			user1: "moderata",
			user2: "anacardo",
		},
		{
			name:      "Non-existent user returns error",
			user1:     "soccerpunk",
			user2:     "nonexistent",
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, users, err := Repo.GetConversationAndUsers(ctx, tt.user1, tt.user2)

			if tt.expectErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)

			// Check conversation IDs are set
			require.NotZero(t, c.ID)
			require.NotZero(t, c.User1ID)
			require.NotZero(t, c.User2ID)

			// Users slice must have exactly 2 users
			require.Len(t, users, 2)

			usernames := []string{users[0].Username, users[1].Username}
			require.Contains(t, usernames, tt.user1)
			require.Contains(t, usernames, tt.user2)

			// Conversation user IDs must match the user IDs
			userIDs := []uuid.UUID{users[0].ID, users[1].ID}
			require.Contains(t, userIDs, c.User1ID)
			require.Contains(t, userIDs, c.User2ID)
		})
	}
}

func TestGetConversationsAndOwner(t *testing.T) {
	ctx := context.Background()

	// Seed DB
	err := SeedDB(ctx, TestDSN, SeedPath)
	require.NoError(t, err)

	tests := []struct {
		name      string
		username  string
		limit     int
		offset    int
		expectErr bool
	}{
		{
			name:     "Get conversations for anacardo",
			username: "anacardo",
			limit:    10,
			offset:   0,
		},
		{
			name:     "Get conversations for moderata",
			username: "moderata",
			limit:    10,
			offset:   0,
		},
		{
			name:      "Non-existent user returns error",
			username:  "ghost",
			limit:     10,
			offset:    0,
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u, convos, err := Repo.GetConversationsAndOwner(ctx, tt.username, tt.limit, tt.offset)

			if tt.expectErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)

			// Check returned user
			require.Equal(t, tt.username, u.Username)
			require.NotZero(t, u.ID)
			// Check conversations
			for _, c := range convos {
				require.NotZero(t, c.ID)
				require.NotNil(t, c.OtherUser)
				require.NotZero(t, c.OtherUser.ID)
				require.NotEmpty(t, c.OtherUser.Username)

				// Messages should be unmarshaled
				for _, m := range c.Messages {
					require.NotZero(t, m.ID)
					require.NotZero(t, m.ConversationID)
					require.NotZero(t, m.SenderID)
				}
			}
		})
	}
}

func TestGetConversationByUsers(t *testing.T) {
	ctx := context.Background()

	// Seed DB
	err := SeedDB(ctx, TestDSN, SeedPath)
	require.NoError(t, err)

	// Example users and conversation from SeedDB
	user1ID := uuid.MustParse("a1f92e18-1d8f-4f0f-9a4d-3b9e3b26b7b1") // anacardo
	user2ID := uuid.MustParse("cfa53179-9085-4f33-86b3-5dc5f7a1465f") // moderata
	nonExistentID := uuid.New()

	tests := []struct {
		name      string
		u1        uuid.UUID
		u2        uuid.UUID
		expectErr bool
	}{
		{
			name:      "Existing conversation normal order",
			u1:        user1ID,
			u2:        user2ID,
			expectErr: false,
		},
		{
			name:      "Existing conversation reversed order",
			u1:        user2ID,
			u2:        user1ID,
			expectErr: false,
		},
		{
			name:      "Non-existent conversation",
			u1:        user1ID,
			u2:        nonExistentID,
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			conv, err := Repo.GetConversationByUsers(ctx, tt.u1, tt.u2)
			if tt.expectErr {
				require.Error(t, err)
				require.Nil(t, conv)
				return
			}

			require.NoError(t, err)
			require.NotNil(t, conv)
			// Check that the conversation contains both users
			require.True(t, (conv.User1ID == tt.u1 && conv.User2ID == tt.u2) ||
				(conv.User1ID == tt.u2 && conv.User2ID == tt.u1))
			require.NotZero(t, conv.ID)
		})
	}
}

func TestUpdateConversation_Table(t *testing.T) {
	ctx := context.Background()

	// Seed DB
	err := SeedDB(ctx, TestDSN, SeedPath)
	require.NoError(t, err)

	// Connect test DB
	db, err := testutils.ConnectDB(ctx, TestDSN)
	require.NoError(t, err)

	tests := []struct {
		name      string
		convoID   uuid.UUID
		expectErr bool
	}{
		{
			name:    "Update existing conversation",
			convoID: uuid.MustParse("70f31f9b-632b-4c4a-bb3a-1f2c6013f001"),
		},
		{
			name:      "Update non-existent conversation",
			convoID:   uuid.New(), // random UUID that likely does not exist
			expectErr: false,      // Exec does not return error if row not found
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Get updated_at before
			var before time.Time
			_ = db.QueryRow(ctx, "SELECT updated_at FROM conversations WHERE id = $1", tt.convoID).Scan(&before)

			// Call UpdateConversation
			err := Repo.UpdateConversation(ctx, tt.convoID)
			if tt.expectErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)

			// Check updated_at only if conversation exists
			var after time.Time
			err = db.QueryRow(ctx, "SELECT updated_at FROM conversations WHERE id = $1", tt.convoID).Scan(&after)
			if err == pgx.ErrNoRows {
				// row doesn't exist, nothing to check
				return
			}
			require.NoError(t, err)
			require.True(t, after.After(before), "updated_at should be later than before update")
		})
	}
}

func TestCreateDM(t *testing.T) {
	ctx := context.Background()

	// Seed DB
	err := SeedDB(ctx, TestDSN, SeedPath)
	require.NoError(t, err)

	// Connect test DB
	db, err := testutils.ConnectDB(ctx, TestDSN)
	require.NoError(t, err)

	// Assume these IDs exist in your seeded DB
	conversationID := uuid.MustParse("70f31f9b-632b-4c4a-bb3a-1f2c6013f001")
	senderID := uuid.MustParse("a1f92e18-1d8f-4f0f-9a4d-3b9e3b26b7b1")

	tests := []struct {
		name      string
		content   string
		expectErr bool
	}{
		{
			name:    "Create valid DM",
			content: "Hello there!",
		},
		{
			name:      "Create DM with empty content",
			content:   "",
			expectErr: false, // depends on your DB schema; assuming empty allowed
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dm := &DMessage{
				ConversationID: conversationID,
				SenderID:       senderID,
				Content:        tt.content,
			}

			id, err := Repo.CreateDM(ctx, dm)
			if tt.expectErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.NotEqual(t, uuid.Nil, id)

			// Verify row exists in DB
			var dbContent string
			err = db.QueryRow(ctx, "SELECT content FROM dmessages WHERE id = $1", id).Scan(&dbContent)
			require.NoError(t, err)
			require.Equal(t, tt.content, dbContent)
		})
	}
}

func TestGetDMsByConversation(t *testing.T) {
	ctx := context.Background()

	// Seed DB
	err := SeedDB(ctx, TestDSN, SeedPath)
	require.NoError(t, err)

	// Assume this conversation exists in your seeded DB
	conversationID := uuid.MustParse("70f31f9b-632b-4c4a-bb3a-1f2c6013f001")

	tests := []struct {
		name           string
		conversationID uuid.UUID
		limit, offset  int
		expectErr      bool
	}{
		{
			name:           "Retrieve messages for existing conversation",
			conversationID: conversationID,
			limit:          10,
			offset:         0,
		},
		{
			name:           "Non-existent conversation returns empty slice",
			conversationID: uuid.New(),
			limit:          10,
			offset:         0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dms, err := Repo.GetDMsByConversation(ctx, tt.conversationID, tt.limit, tt.offset)
			if tt.expectErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)

			// If conversation exists, check that messages have valid fields
			for _, dm := range dms {
				require.NotZero(t, dm.ID)
				require.Equal(t, tt.conversationID, dm.ConversationID)
				require.NotZero(t, dm.Sender.ID)
				require.NotEmpty(t, dm.Sender.Username)
			}
		})
	}
}

func TestReadAllReceivedDMsInConvo(t *testing.T) {
	ctx := context.Background()

	// Seed DB
	err := SeedDB(ctx, TestDSN, SeedPath)
	require.NoError(t, err)

	// Assume these exist in your seeded DB
	conversationID := uuid.MustParse("70f31f9b-632b-4c4a-bb3a-1f2c6013f001")
	recipientUsername := "moderata"

	tests := []struct {
		name           string
		conversationID uuid.UUID
		username       string
		expectErr      bool
	}{
		{
			name:           "Mark received DMs as read for existing user",
			conversationID: conversationID,
			username:       recipientUsername,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Repo.ReadAllReceivedDMsInConvo(ctx, tt.conversationID, tt.username)
			if tt.expectErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)

			// Verify that all messages sent to this user are now marked as read
			dms, err := Repo.GetDMsByConversation(ctx, tt.conversationID, 100, 0)
			require.NoError(t, err)

			for _, dm := range dms {
				if dm.Sender.Username != tt.username {
					require.True(t, dm.IsRead, "Message from %s should be marked as read", dm.Sender.Username)
				}
			}
		})
	}
}
