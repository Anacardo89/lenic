package repo

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/Anacardo89/lenic/pkg/testutils"
)

func TestCreateComment(t *testing.T) {
	ctx := context.Background()

	// Re-seed the database before running
	err := SeedDB(ctx, TestDSN, SeedPath)
	require.NoError(t, err)

	// Connect test DB
	db, err := testutils.ConnectDB(ctx, TestDSN)
	require.NoError(t, err)

	validPostID := uuid.MustParse("b1d3c0f7-5a1a-4f9b-9b2a-2a8e4f8b9f01")   // Post 1
	validAuthorID := uuid.MustParse("f7a92b5b-4b7e-4787-9c0b-2b0b6cb86b4e") // "soccerpunk"

	tests := []struct {
		name      string
		comment   Comment
		expectErr bool
		verify    func(t *testing.T, c *Comment)
	}{
		{
			name: "Valid comment on existing post",
			comment: Comment{
				PostID:   validPostID,
				AuthorID: validAuthorID,
				Content:  "Looking forward to your next post!",
			},
			expectErr: false,
			verify: func(t *testing.T, c *Comment) {
				require.NotEqual(t, uuid.Nil, c.ID)
				require.Equal(t, validPostID, c.PostID)
				require.Equal(t, validAuthorID, c.AuthorID)
				require.Equal(t, "Looking forward to your next post!", c.Content)
				require.Equal(t, 0, c.Rating)          // default rating should be 0
				require.True(t, c.IsActive)            // should default to TRUE
				require.False(t, c.CreatedAt.IsZero()) // should have timestamp
			},
		},
		{
			name: "Invalid post_id should fail",
			comment: Comment{
				PostID:   uuid.MustParse("00000000-0000-0000-0000-000000000000"),
				AuthorID: validAuthorID,
				Content:  "This should not insert.",
			},
			expectErr: true,
		},
		{
			name: "Invalid author_id should fail",
			comment: Comment{
				PostID:   validPostID,
				AuthorID: uuid.MustParse("ffffffff-ffff-ffff-ffff-ffffffffffff"),
				Content:  "Invalid user posting comment.",
			},
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := tt.comment

			err := Repo.CreateComment(ctx, &c)

			if tt.expectErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			require.NotEqual(t, uuid.Nil, c.ID)

			// Optional: double-check comment actually exists in DB
			var exists bool
			err = db.QueryRow(ctx,
				`SELECT EXISTS(SELECT 1 FROM comments WHERE id = $1)`, c.ID,
			).Scan(&exists)
			require.NoError(t, err)
			require.True(t, exists, "comment should exist in database")

			if tt.verify != nil {
				tt.verify(t, &c)
			}
		})
	}
}

func TestGetComment(t *testing.T) {
	ctx := context.Background()

	// Seed DB before test
	err := SeedDB(ctx, TestDSN, SeedPath)
	require.NoError(t, err)

	// Known comment IDs from seed
	existingCommentID := uuid.MustParse("e1a1b2c3-d4f5-4a6b-9c7d-1e2f3a4b5c6d")
	nonExistentCommentID := uuid.MustParse("00000000-0000-0000-0000-000000000000")

	tests := []struct {
		name      string
		commentID uuid.UUID
		expectErr bool
		verify    func(t *testing.T, c *Comment)
	}{
		{
			name:      "Get existing comment",
			commentID: existingCommentID,
			expectErr: false,
			verify: func(t *testing.T, c *Comment) {
				require.Equal(t, existingCommentID, c.ID)
				require.Equal(t, uuid.MustParse("b1d3c0f7-5a1a-4f9b-9b2a-2a8e4f8b9f01"), c.PostID)
				require.Equal(t, uuid.MustParse("cfa53179-9085-4f33-86b3-5dc5f7a1465f"), c.AuthorID)
				require.Equal(t, "Great insights, really makes me think!", c.Content)
				require.Equal(t, 2, c.Rating)
				require.True(t, c.IsActive)
				require.False(t, c.CreatedAt.IsZero())
			},
		},
		{
			name:      "Get non-existent comment returns error",
			commentID: nonExistentCommentID,
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			comment, err := Repo.GetComment(ctx, tt.commentID)

			if tt.expectErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			require.NotNil(t, comment)
			if tt.verify != nil {
				tt.verify(t, comment)
			}
		})
	}
}

func TestUpdateComment(t *testing.T) {
	ctx := context.Background()

	// Seed DB
	err := SeedDB(ctx, TestDSN, SeedPath)
	require.NoError(t, err)

	// Connect test DB
	db, err := testutils.ConnectDB(ctx, TestDSN)
	require.NoError(t, err)

	// Known comment IDs from seed
	existingCommentID := uuid.MustParse("e1a1b2c3-d4f5-4a6b-9c7d-1e2f3a4b5c6d")
	nonExistentCommentID := uuid.MustParse("00000000-0000-0000-0000-000000000000")

	tests := []struct {
		name      string
		comment   *Comment
		expectErr bool
		verify    func(t *testing.T, c *Comment)
	}{
		{
			name: "Update existing comment content",
			comment: &Comment{
				ID:      existingCommentID,
				Content: "Updated comment content for testing",
			},
			expectErr: false,
			verify: func(t *testing.T, c *Comment) {
				require.Equal(t, existingCommentID, c.ID)
				require.Equal(t, "Updated comment content for testing", c.Content)
				require.True(t, c.IsActive)
			},
		},
		{
			name: "Update non-existent comment returns error",
			comment: &Comment{
				ID:      nonExistentCommentID,
				Content: "This should fail",
			},
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Repo.UpdateComment(ctx, tt.comment)

			if tt.expectErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			require.NotNil(t, tt.comment)

			if tt.verify != nil {
				tt.verify(t, tt.comment)
			}

			// Extra check: read from DB to ensure content is actually updated
			var content string
			err = db.QueryRow(ctx, `SELECT content FROM comments WHERE id = $1;`, tt.comment.ID).Scan(&content)
			require.NoError(t, err)
			require.Equal(t, tt.comment.Content, content)
		})
	}
}

func TestDisableComment(t *testing.T) {
	ctx := context.Background()

	// Seed DB
	err := SeedDB(ctx, TestDSN, SeedPath)
	require.NoError(t, err)

	// Connect test DB
	db, err := testutils.ConnectDB(ctx, TestDSN)
	require.NoError(t, err)

	// Known comment IDs from seed
	existingCommentID := uuid.MustParse("f2b2c3d4-e5f6-4b7c-8d9e-2f3a4b5c6d7e")
	nonExistentCommentID := uuid.MustParse("00000000-0000-0000-0000-000000000000")

	tests := []struct {
		name      string
		commentID uuid.UUID
		expectErr bool
		verify    func(t *testing.T, c *Comment)
	}{
		{
			name:      "Disable existing comment",
			commentID: existingCommentID,
			expectErr: false,
			verify: func(t *testing.T, c *Comment) {
				require.Equal(t, existingCommentID, c.ID)
				require.NotEmpty(t, c.Content)
			},
		},
		{
			name:      "Disable non-existent comment returns error",
			commentID: nonExistentCommentID,
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			comment, err := Repo.DisableComment(ctx, tt.commentID)

			if tt.expectErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			require.NotNil(t, comment)

			if tt.verify != nil {
				tt.verify(t, comment)
			}

			// Extra check: confirm comment is inactive in DB
			var isActive bool
			err = db.QueryRow(ctx, `SELECT is_active FROM comments WHERE id = $1;`, tt.commentID).Scan(&isActive)
			require.NoError(t, err)
			require.False(t, isActive)
		})
	}
}

func TestRateCommentUp(t *testing.T) {
	ctx := context.Background()

	// Seed DB
	err := SeedDB(ctx, TestDSN, SeedPath)
	require.NoError(t, err)

	// Connect test DB
	db, err := testutils.ConnectDB(ctx, TestDSN)
	require.NoError(t, err)

	// Known comment IDs from seed
	commentID := uuid.MustParse("f2b2c3d4-e5f6-4b7c-8d9e-2f3a4b5c6d7e") // comment by soccerpunk
	userID := uuid.MustParse("a1f92e18-1d8f-4f0f-9a4d-3b9e3b26b7b1")    // anacardo

	tests := []struct {
		name          string
		targetID      uuid.UUID
		userID        uuid.UUID
		expectedValue int
		expectErr     bool
	}{
		{
			name:          "First upvote creates rating = 1",
			targetID:      commentID,
			userID:        userID,
			expectedValue: 1,
		},
		{
			name:          "Second upvote toggles rating to 0",
			targetID:      commentID,
			userID:        userID,
			expectedValue: 0,
		},
		{
			name:          "Third upvote toggles rating back to 1",
			targetID:      commentID,
			userID:        userID,
			expectedValue: 1,
		},
		{
			name:      "Invalid comment ID returns error",
			targetID:  uuid.MustParse("00000000-0000-0000-0000-000000000000"),
			userID:    userID,
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Repo.RateCommentUp(ctx, tt.targetID, tt.userID)

			if tt.expectErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)

			var ratingValue int
			err = db.QueryRow(ctx, `SELECT rating_value FROM comment_ratings WHERE target_id = $1 AND user_id = $2;`,
				tt.targetID, tt.userID).Scan(&ratingValue)
			require.NoError(t, err)
			require.Equal(t, tt.expectedValue, ratingValue)
		})
	}
}

func TestRateCommentDown(t *testing.T) {
	ctx := context.Background()

	// Seed DB
	err := SeedDB(ctx, TestDSN, SeedPath)
	require.NoError(t, err)

	// Connect test DB
	db, err := testutils.ConnectDB(ctx, TestDSN)
	require.NoError(t, err)

	// Known comment IDs from seed
	commentID := uuid.MustParse("f2b2c3d4-e5f6-4b7c-8d9e-2f3a4b5c6d7e") // comment by soccerpunk
	userID := uuid.MustParse("a1f92e18-1d8f-4f0f-9a4d-3b9e3b26b7b1")    // anacardo

	tests := []struct {
		name          string
		targetID      uuid.UUID
		userID        uuid.UUID
		expectedValue int
		expectErr     bool
	}{
		{
			name:          "First downvote creates rating = -1",
			targetID:      commentID,
			userID:        userID,
			expectedValue: -1,
		},
		{
			name:          "Second downvote toggles rating to 0",
			targetID:      commentID,
			userID:        userID,
			expectedValue: 0,
		},
		{
			name:          "Third downvote toggles rating back to -1",
			targetID:      commentID,
			userID:        userID,
			expectedValue: -1,
		},
		{
			name:      "Invalid comment ID returns error",
			targetID:  uuid.MustParse("00000000-0000-0000-0000-000000000000"),
			userID:    userID,
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Repo.RateCommentDown(ctx, tt.targetID, tt.userID)

			if tt.expectErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)

			var ratingValue int
			err = db.QueryRow(ctx, `SELECT rating_value FROM comment_ratings WHERE target_id = $1 AND user_id = $2;`,
				tt.targetID, tt.userID).Scan(&ratingValue)
			require.NoError(t, err)
			require.Equal(t, tt.expectedValue, ratingValue)
		})
	}
}
