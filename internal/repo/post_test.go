package repo

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/Anacardo89/lenic/pkg/testutils"
)

func TestCreatePost(t *testing.T) {
	// Seed DB
	ctx := context.Background()
	err := SeedDB(ctx, TestDSN, SeedPath)
	require.NoError(t, err)

	// UUIDs from seed
	anacardoID := uuid.MustParse("a1f92e18-1d8f-4f0f-9a4d-3b9e3b26b7b1")
	moderataID := uuid.MustParse("cfa53179-9085-4f33-86b3-5dc5f7a1465f")

	tests := []struct {
		name    string
		input   *Post
		wantErr bool
	}{
		{
			name: "success - create post for anacardo",
			input: &Post{
				ID:        uuid.New(),
				AuthorID:  anacardoID,
				Title:     "New Thoughts",
				Content:   "Sharing new ideas about politics.",
				PostImage: "",
				IsPublic:  true,
			},
			wantErr: false,
		},
		{
			name: "success - create post for moderata",
			input: &Post{
				ID:        uuid.New(),
				AuthorID:  moderataID,
				Title:     "Moderation Insights",
				Content:   "Tips for community moderation.",
				PostImage: "",
				IsPublic:  false,
			},
			wantErr: false,
		},
		{
			name: "fail - invalid author",
			input: &Post{
				ID:       uuid.New(),
				AuthorID: uuid.New(), // non-existent user
				Title:    "Ghost Author",
				Content:  "This author doesn't exist",
				IsPublic: true,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			id, err := Repo.CreatePost(ctx, tt.input)
			if tt.wantErr {
				require.Error(t, err)
				require.Equal(t, uuid.Nil, id)
				return
			}

			require.NoError(t, err)
			require.NotEqual(t, uuid.Nil, id)

			// Verify the post exists in DB
			post, err := Repo.GetPost(ctx, id)
			require.NoError(t, err)
			require.Equal(t, tt.input.Title, post.Title)
			require.Equal(t, tt.input.Content, post.Content)
			require.Equal(t, tt.input.AuthorID, post.AuthorID)
			require.Equal(t, tt.input.PostImage, post.PostImage)
			require.Equal(t, tt.input.IsPublic, post.IsPublic)
		})
	}
}

func TestGetFeed(t *testing.T) {
	// Seed DB
	ctx := context.Background()
	err := SeedDB(ctx, TestDSN, SeedPath)
	require.NoError(t, err)

	tests := []struct {
		name          string
		username      string
		expectedPosts []string // expected post titles in feed
	}{
		{
			name:     "feed for admin user anacardo",
			username: "anacardo",
			expectedPosts: []string{
				"Thoughts on Local Politics", // own post
				"Moderation Tips",            // followed posts? (anacardo follows moderata)
				// "Weekend Football Recap" (soccerpunk post not is public)
			},
		},
		{
			name:     "feed for moderator user moderata",
			username: "moderata",
			expectedPosts: []string{
				"Moderation Tips",            // own post
				"Thoughts on Local Politics", // public post from anacardo
				// "Weekend Football Recap" is not followed and public? Check seed
			},
		},
		{
			name:     "feed for regular user soccerpunk",
			username: "soccerpunk",
			expectedPosts: []string{
				"Weekend Football Recap",     // own post
				"Thoughts on Local Politics", // public post
				"Moderation Tips",            // followed by soccerpunk? No, only follows anacardo accepted
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			posts, err := Repo.GetFeed(ctx, tt.username)
			if tt.username == "ghostuser" {
				require.NoError(t, err)
				require.Len(t, posts, 0)
				return
			}

			require.NoError(t, err)
			require.NotNil(t, posts)

			titles := make([]string, len(posts))
			for i, p := range posts {
				titles[i] = p.Title
			}

			for _, expected := range tt.expectedPosts {
				require.Contains(t, titles, expected, "expected post title %q not found", expected)
			}
		})
	}
}

func TestGetPostAuthorFromComment(t *testing.T) {
	// Seed DB
	ctx := context.Background()
	err := SeedDB(ctx, TestDSN, SeedPath)
	require.NoError(t, err)

	tests := []struct {
		name          string
		commentID     uuid.UUID
		expectedID    uuid.UUID
		expectedName  string
		expectedError bool
	}{
		{
			name:         "Comment 1 author is anacardo",
			commentID:    uuid.MustParse("e1a1b2c3-d4f5-4a6b-9c7d-1e2f3a4b5c6d"),
			expectedID:   uuid.MustParse("a1f92e18-1d8f-4f0f-9a4d-3b9e3b26b7b1"),
			expectedName: "anacardo",
		},
		{
			name:         "Comment 2 author is anacardo",
			commentID:    uuid.MustParse("f2b2c3d4-e5f6-4b7c-8d9e-2f3a4b5c6d7e"),
			expectedID:   uuid.MustParse("a1f92e18-1d8f-4f0f-9a4d-3b9e3b26b7b1"),
			expectedName: "anacardo",
		},
		{
			name:         "Comment 3 author is moderata",
			commentID:    uuid.MustParse("a3c3d4e5-f6a7-4c8d-9e0f-3a4b5c6d7e8f"),
			expectedID:   uuid.MustParse("cfa53179-9085-4f33-86b3-5dc5f7a1465f"),
			expectedName: "moderata",
		},
		{
			name:         "Comment 4 author is soccerpunk",
			commentID:    uuid.MustParse("b4d4e5f6-a7b8-4d9e-0f1a-4b5c6d7e8f9a"),
			expectedID:   uuid.MustParse("f7a92b5b-4b7e-4787-9c0b-2b0b6cb86b4e"),
			expectedName: "soccerpunk",
		},
		{
			name:          "Non-existent comment returns error",
			commentID:     uuid.MustParse("00000000-0000-0000-0000-000000000000"),
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.NoError(t, err)

			user, err := Repo.GetPostAuthorFromComment(ctx, tt.commentID)
			if tt.expectedError {
				require.Error(t, err)
				require.Nil(t, user)
				return
			}

			require.NoError(t, err)
			require.NotNil(t, user)
			require.Equal(t, tt.expectedID, user.ID)
			require.Equal(t, tt.expectedName, user.Username)
		})
	}
}

func TestGetUserPosts(t *testing.T) {
	// Seed DB
	ctx := context.Background()
	err := SeedDB(ctx, TestDSN, SeedPath)
	require.NoError(t, err)

	tests := []struct {
		name          string
		userID        uuid.UUID
		expectedPosts []string
	}{
		{
			name:   "admin user anacardo posts",
			userID: uuid.MustParse("a1f92e18-1d8f-4f0f-9a4d-3b9e3b26b7b1"),
			expectedPosts: []string{
				"Thoughts on Local Politics",
			},
		},
		{
			name:   "moderator user moderata posts",
			userID: uuid.MustParse("cfa53179-9085-4f33-86b3-5dc5f7a1465f"),
			expectedPosts: []string{
				"Moderation Tips",
			},
		},
		{
			name:   "regular user soccerpunk posts",
			userID: uuid.MustParse("f7a92b5b-4b7e-4787-9c0b-2b0b6cb86b4e"),
			expectedPosts: []string{
				"Weekend Football Recap", // note: is_active TRUE, is_public FALSE
			},
		},
		{
			name:          "inactive user has no posts",
			userID:        uuid.MustParse("d8e3f4a5-b6c7-4d8e-9f0a-1b2c3d4e5f6d"),
			expectedPosts: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			posts, err := Repo.GetUserPosts(ctx, tt.userID)
			require.NoError(t, err)
			require.NotNil(t, posts)
			require.Len(t, posts, len(tt.expectedPosts))

			titles := make([]string, len(posts))
			for i, p := range posts {
				titles[i] = p.Title
			}

			for _, expected := range tt.expectedPosts {
				require.Contains(t, titles, expected)
			}
		})
	}
}

func TestGetUserPublicPosts(t *testing.T) {
	// Seed DB
	ctx := context.Background()
	err := SeedDB(ctx, TestDSN, SeedPath)
	require.NoError(t, err)

	tests := []struct {
		name          string
		userID        uuid.UUID
		expectedPosts []string
	}{
		{
			name:   "admin user anacardo public posts",
			userID: uuid.MustParse("a1f92e18-1d8f-4f0f-9a4d-3b9e3b26b7b1"),
			expectedPosts: []string{
				"Thoughts on Local Politics",
			},
		},
		{
			name:   "moderator user moderata public posts",
			userID: uuid.MustParse("cfa53179-9085-4f33-86b3-5dc5f7a1465f"),
			expectedPosts: []string{
				"Moderation Tips",
			},
		},
		{
			name:          "regular user soccerpunk public posts (none)",
			userID:        uuid.MustParse("f7a92b5b-4b7e-4787-9c0b-2b0b6cb86b4e"),
			expectedPosts: []string{},
		},
		{
			name:          "inactive user has no public posts",
			userID:        uuid.MustParse("d8e3f4a5-b6c7-4d8e-9f0a-1b2c3d4e5f6d"),
			expectedPosts: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			posts, err := Repo.GetUserPublicPosts(ctx, tt.userID)
			require.NoError(t, err)
			require.NotNil(t, posts)
			require.Len(t, posts, len(tt.expectedPosts))

			titles := make([]string, len(posts))
			for i, p := range posts {
				titles[i] = p.Title
			}

			for _, expected := range tt.expectedPosts {
				require.Contains(t, titles, expected)
			}
		})
	}
}

func TestGetPost(t *testing.T) {
	// Seed DB
	ctx := context.Background()
	err := SeedDB(ctx, TestDSN, SeedPath)
	require.NoError(t, err)

	tests := []struct {
		name          string
		postID        uuid.UUID
		expectedTitle string
		expectedError bool
	}{
		{
			name:          "Get Post 1 - Thoughts on Local Politics",
			postID:        uuid.MustParse("b1d3c0f7-5a1a-4f9b-9b2a-2a8e4f8b9f01"),
			expectedTitle: "Thoughts on Local Politics",
		},
		{
			name:          "Get Post 2 - Moderation Tips",
			postID:        uuid.MustParse("c2e4d1f8-6b2b-4a8c-8c3b-3b9f5a9c0d12"),
			expectedTitle: "Moderation Tips",
		},
		{
			name:          "Get Post 3 - Weekend Football Recap",
			postID:        uuid.MustParse("d3f5e2a9-7c3c-4b7d-9d4c-4c0a6b1d1e23"),
			expectedTitle: "Weekend Football Recap",
		},
		{
			name:          "Non-existent post returns error",
			postID:        uuid.MustParse("00000000-0000-0000-0000-000000000000"),
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			post, err := Repo.GetPost(ctx, tt.postID)
			if tt.expectedError {
				require.Error(t, err)
				require.Nil(t, post)
				return
			}
			require.NoError(t, err)
			require.NotNil(t, post)
			require.Equal(t, tt.expectedTitle, post.Title)
			require.Equal(t, tt.postID, post.ID)
		})
	}
}

func TestGetPostForPage(t *testing.T) {
	// Seed DB
	ctx := context.Background()
	err := SeedDB(ctx, TestDSN, SeedPath)
	require.NoError(t, err)

	tests := []struct {
		name               string
		postID             uuid.UUID
		userID             uuid.UUID
		expectedTitle      string
		expectedAuthor     string
		expectedComments   int
		expectedUserRating int
		expectedError      bool
	}{
		{
			name:               "Post 1 page for admin anacardo",
			postID:             uuid.MustParse("b1d3c0f7-5a1a-4f9b-9b2a-2a8e4f8b9f01"),
			userID:             uuid.MustParse("a1f92e18-1d8f-4f0f-9a4d-3b9e3b26b7b1"),
			expectedTitle:      "Thoughts on Local Politics",
			expectedAuthor:     "anacardo",
			expectedComments:   2,
			expectedUserRating: 0,
		},
		{
			name:               "Post 1 page for moderator moderata",
			postID:             uuid.MustParse("b1d3c0f7-5a1a-4f9b-9b2a-2a8e4f8b9f01"),
			userID:             uuid.MustParse("cfa53179-9085-4f33-86b3-5dc5f7a1465f"),
			expectedTitle:      "Thoughts on Local Politics",
			expectedAuthor:     "anacardo",
			expectedComments:   2,
			expectedUserRating: 1,
		},
		{
			name:          "Non-existent post returns error",
			postID:        uuid.MustParse("00000000-0000-0000-0000-000000000000"),
			userID:        uuid.MustParse("a1f92e18-1d8f-4f0f-9a4d-3b9e3b26b7b1"),
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			post, err := Repo.GetPostForPage(ctx, tt.postID, tt.userID)
			if tt.expectedError {
				require.Error(t, err)
				require.Nil(t, post)
				return
			}
			require.NoError(t, err)
			require.NotNil(t, post)
			require.Equal(t, tt.expectedTitle, post.Title)
			require.Equal(t, tt.expectedAuthor, post.Author.Username)
			require.Len(t, post.Comments, tt.expectedComments)
			require.Equal(t, tt.expectedUserRating, post.UserRating)
		})
	}
}

func TestUpdatePost(t *testing.T) {
	// Seed DB
	ctx := context.Background()
	err := SeedDB(ctx, TestDSN, SeedPath)
	require.NoError(t, err)

	tests := []struct {
		name       string
		postID     uuid.UUID
		newTitle   string
		newContent string
		newPublic  bool
		expectErr  bool
	}{
		{
			name:       "Update Post 1 title and content",
			postID:     uuid.MustParse("b1d3c0f7-5a1a-4f9b-9b2a-2a8e4f8b9f01"),
			newTitle:   "Updated Thoughts on Politics",
			newContent: "Updated content for the post.",
			newPublic:  true,
		},
		{
			name:       "Update Post 2 to private",
			postID:     uuid.MustParse("c2e4d1f8-6b2b-4a8c-8c3b-3b9f5a9c0d12"),
			newTitle:   "Moderation Tips Updated",
			newContent: "Tips updated for the forum.",
			newPublic:  false,
		},
		{
			name:      "Non-existent post returns error",
			postID:    uuid.MustParse("00000000-0000-0000-0000-000000000000"),
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			post := &Post{
				ID:       tt.postID,
				Title:    tt.newTitle,
				Content:  tt.newContent,
				IsPublic: tt.newPublic,
			}

			err := Repo.UpdatePost(ctx, post)
			if tt.expectErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)

			// Verify that the post was updated in DB
			updatedPost, err := Repo.GetPost(ctx, tt.postID)
			require.NoError(t, err)
			require.NotNil(t, updatedPost)
			require.Equal(t, tt.newTitle, updatedPost.Title)
			require.Equal(t, tt.newContent, updatedPost.Content)
			require.Equal(t, tt.newPublic, updatedPost.IsPublic)
		})
	}
}

func TestDisablePost(t *testing.T) {
	// Seed DB
	ctx := context.Background()
	err := SeedDB(ctx, TestDSN, SeedPath)
	require.NoError(t, err)

	// DB for test query
	db, err := testutils.ConnectDB(ctx, TestDSN)
	require.NoError(t, err)

	tests := []struct {
		name      string
		postID    uuid.UUID
		expectErr bool
	}{
		{
			name:   "Disable existing post (Post 1)",
			postID: uuid.MustParse("b1d3c0f7-5a1a-4f9b-9b2a-2a8e4f8b9f01"),
		},
		{
			name:      "Non-existent post returns error",
			postID:    uuid.MustParse("00000000-0000-0000-0000-000000000000"),
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			disabledPost, err := Repo.DisablePost(ctx, tt.postID)
			if tt.expectErr {
				require.Error(t, err)
				require.Nil(t, disabledPost)
				return
			}

			require.NoError(t, err)
			require.NotNil(t, disabledPost)
			require.Equal(t, tt.postID, disabledPost.ID)
			require.NotEmpty(t, disabledPost.Title)
			require.NotEmpty(t, disabledPost.Content)

			// Verify in DB that the post is inactive and has deleted_at set
			var (
				isActive  bool
				deletedAt *time.Time
			)
			err = db.QueryRow(ctx, `SELECT is_active, deleted_at FROM posts WHERE id = $1;`,
				tt.postID).Scan(&isActive, &deletedAt)
			require.NoError(t, err)
			require.False(t, isActive)
			require.NotNil(t, deletedAt)
		})
	}
}

func TestRatePostUp(t *testing.T) {
	// Seed DB
	ctx := context.Background()
	err := SeedDB(ctx, TestDSN, SeedPath)
	require.NoError(t, err)

	// DB for test query
	db, err := testutils.ConnectDB(ctx, TestDSN)
	require.NoError(t, err)

	postID := uuid.MustParse("b1d3c0f7-5a1a-4f9b-9b2a-2a8e4f8b9f01")
	userID := uuid.MustParse("a1f92e18-1d8f-4f0f-9a4d-3b9e3b26b7b1") // "anacardo"

	tests := []struct {
		name          string
		targetID      uuid.UUID
		userID        uuid.UUID
		expectedValue int
		expectErr     bool
		setup         func() // optional setup before test
	}{
		{
			name:          "First upvote creates rating = 1",
			targetID:      postID,
			userID:        userID,
			expectedValue: 1,
		},
		{
			name:          "Second upvote toggles rating to 0",
			targetID:      postID,
			userID:        userID,
			expectedValue: 0,
		},
		{
			name:          "Third upvote toggles rating back to 1",
			targetID:      postID,
			userID:        userID,
			expectedValue: 1,
		},
		{
			name:      "Invalid post ID returns error",
			targetID:  uuid.MustParse("00000000-0000-0000-0000-000000000000"),
			userID:    userID,
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setup != nil {
				tt.setup()
			}

			err := Repo.RatePostUp(ctx, tt.targetID, tt.userID)

			if tt.expectErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)

			var ratingValue int
			err = db.QueryRow(ctx, `SELECT rating_value FROM post_ratings WHERE target_id = $1 AND user_id = $2;`,
				tt.targetID, tt.userID).Scan(&ratingValue)
			require.NoError(t, err)
			require.Equal(t, tt.expectedValue, ratingValue)
		})
	}
}

func TestRatePostDown(t *testing.T) {
	// Seed DB
	ctx := context.Background()
	err := SeedDB(ctx, TestDSN, SeedPath)
	require.NoError(t, err)

	// Connect to DB for assertions
	db, err := testutils.ConnectDB(ctx, TestDSN)
	require.NoError(t, err)

	postID := uuid.MustParse("b1d3c0f7-5a1a-4f9b-9b2a-2a8e4f8b9f01")
	userID := uuid.MustParse("a1f92e18-1d8f-4f0f-9a4d-3b9e3b26b7b1") // "anacardo"

	tests := []struct {
		name          string
		targetID      uuid.UUID
		userID        uuid.UUID
		expectedValue int
		expectErr     bool
		setup         func()
	}{
		{
			name:          "First downvote creates rating = -1",
			targetID:      postID,
			userID:        userID,
			expectedValue: -1,
		},
		{
			name:          "Second downvote toggles rating to 0",
			targetID:      postID,
			userID:        userID,
			expectedValue: 0,
		},
		{
			name:          "Third downvote toggles rating back to -1",
			targetID:      postID,
			userID:        userID,
			expectedValue: -1,
		},
		{
			name:      "Invalid post ID returns error",
			targetID:  uuid.MustParse("00000000-0000-0000-0000-000000000000"),
			userID:    userID,
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setup != nil {
				tt.setup()
			}

			err := Repo.RatePostDown(ctx, tt.targetID, tt.userID)

			if tt.expectErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)

			var ratingValue int
			err = db.QueryRow(ctx, `SELECT rating_value FROM post_ratings WHERE target_id = $1 AND user_id = $2;`,
				tt.targetID, tt.userID).Scan(&ratingValue)
			require.NoError(t, err)
			require.Equal(t, tt.expectedValue, ratingValue)
		})
	}
}
