package repo

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/Anacardo89/lenic/pkg/testutils"
)

func TestCreateUserTag(t *testing.T) {
	ctx := context.Background()

	// Connect to test DB
	db, err := testutils.ConnectDB(ctx, TestDSN)
	require.NoError(t, err)

	// Seed DB
	err = SeedDB(ctx, TestDSN, SeedPath)
	require.NoError(t, err)

	tests := []struct {
		name      string
		tag       *UserTag
		expectErr bool
	}{
		{
			name: "Tag user in a post",
			tag: &UserTag{
				UserID:       uuid.MustParse("a1f92e18-1d8f-4f0f-9a4d-3b9e3b26b7b1"), // anacardo
				TargetID:     uuid.MustParse("c2e4d1f8-6b2b-4a8c-8c3b-3b9f5a9c0d12"), // Moderation Tips
				ResourceType: "post",
			},
		},
		{
			name: "Tag user in a comment",
			tag: &UserTag{
				UserID:       uuid.MustParse("a1f92e18-1d8f-4f0f-9a4d-3b9e3b26b7b1"), // anacardo
				TargetID:     uuid.MustParse("e1a1b2c3-d4f5-4a6b-9c7d-1e2f3a4b5c6d"), // Comment on Post 1
				ResourceType: "comment",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Repo.CreateUserTag(ctx, tt.tag)
			if tt.expectErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)

			// Verify the tag exists in DB
			var count int
			err = db.QueryRow(ctx, `
				SELECT COUNT(*) 
				FROM user_tags 
				WHERE user_id = $1 
				  AND target_id = $2 
				  AND resource_type = $3
			`, tt.tag.UserID, tt.tag.TargetID, tt.tag.ResourceType).Scan(&count)
			require.NoError(t, err)
			require.Equal(t, 1, count)
		})
	}
}

func TestDeleteUserTag(t *testing.T) {
	ctx := context.Background()

	// Seed DB
	err := SeedDB(ctx, TestDSN, SeedPath)
	require.NoError(t, err)

	// Connect to test DB
	db, err := testutils.ConnectDB(ctx, TestDSN)
	require.NoError(t, err)

	// Insert a tag to be deleted
	tag := &UserTag{
		UserID:       uuid.MustParse("a1f92e18-1d8f-4f0f-9a4d-3b9e3b26b7b1"), // anacardo
		TargetID:     uuid.MustParse("c2e4d1f8-6b2b-4a8c-8c3b-3b9f5a9c0d12"), // Moderation Tips
		ResourceType: "post",
	}
	_, err = db.Exec(ctx, `
		INSERT INTO user_tags (user_id, target_id, resource_type)
		VALUES ($1, $2, $3)
	`, tag.UserID, tag.TargetID, tag.ResourceType)
	require.NoError(t, err)

	tests := []struct {
		name      string
		userID    uuid.UUID
		targetID  uuid.UUID
		expectErr bool
	}{
		{
			name:     "Delete existing tag",
			userID:   tag.UserID,
			targetID: tag.TargetID,
		},
		{
			name:      "Delete non-existent tag",
			userID:    uuid.MustParse("f7a92b5b-4b7e-4787-9c0b-2b0b6cb86b4e"),
			targetID:  uuid.MustParse("d3f5e2a9-7c3c-4b7d-9d4c-4c0a6b1d1e23"),
			expectErr: false, // should succeed but delete 0 rows
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Repo.DeleteUserTag(ctx, tt.userID, tt.targetID)
			if tt.expectErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)

			// Verify the tag no longer exists
			var count int
			err = db.QueryRow(ctx, `
				SELECT COUNT(*) 
				FROM user_tags 
				WHERE user_id = $1 AND target_id = $2
			`, tt.userID, tt.targetID).Scan(&count)
			require.NoError(t, err)
			require.Equal(t, 0, count)
		})
	}
}
