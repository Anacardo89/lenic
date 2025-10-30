package repo

import (
	"context"
	"path/filepath"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestCreatePost(t *testing.T) {
	// Init
	ctx := context.Background()
	repo, dsn, closeDB, seedPath, err := BuildTestDBEnv(ctx)
	require.NoError(t, err)
	defer closeDB()
	// Seed DB
	seed := filepath.Join(seedPath, "repo_tests.sql")
	err = SeedDB(ctx, dsn, seed)
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
			id, err := repo.CreatePost(ctx, tt.input)
			if tt.wantErr {
				require.Error(t, err)
				require.Equal(t, uuid.Nil, id)
				return
			}

			require.NoError(t, err)
			require.NotEqual(t, uuid.Nil, id)

			// Verify the post exists in DB
			post, err := repo.GetPost(ctx, id)
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
	// Init
	ctx := context.Background()
	repo, dsn, closeDB, seedPath, err := BuildTestDBEnv(ctx)
	require.NoError(t, err)
	defer closeDB()
	// Seed DB
	seed := filepath.Join(seedPath, "repo_tests.sql")
	err = SeedDB(ctx, dsn, seed)
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
			posts, err := repo.GetFeed(ctx, tt.username)
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
