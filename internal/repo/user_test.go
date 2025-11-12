package repo

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stretchr/testify/require"

	"github.com/Anacardo89/lenic/pkg/testutils"
)

func TestCreateUser(t *testing.T) {
	// Seed DB
	ctx := context.Background()
	err := SeedDB(ctx, TestDSN, SeedPath)
	require.NoError(t, err)

	tests := []struct {
		name        string
		input       *User
		wantErr     error
		wantNilUUID bool
	}{
		{
			name: "success - valid user",
			input: &User{
				Username:     "alice",
				Email:        "alice@example.com",
				PasswordHash: "hashed_password_123",
			},
			wantErr:     nil,
			wantNilUUID: false,
		},
		{
			name: "fail - duplicate username",
			input: &User{
				Username:     "alice",
				Email:        "alice2@example.com",
				PasswordHash: "hashed_password_123",
			},
			wantErr:     ErrUserExists,
			wantNilUUID: true,
		},
		{
			name: "fail - duplicate email",
			input: &User{
				Username:     "bob",
				Email:        "alice@example.com",
				PasswordHash: "hashed_password_123",
			},
			wantErr:     ErrUserExists,
			wantNilUUID: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			id, err := Repo.CreateUser(ctx, tt.input)

			if tt.wantErr == nil {
				require.NoError(t, err)
				require.NotEqual(t, uuid.Nil, id)
				// Verify user exists in DB
				user, err := Repo.GetUserByID(ctx, id)
				require.NoError(t, err)
				require.Equal(t, tt.input.Username, user.Username)
				require.Equal(t, tt.input.Email, user.Email)
			} else {
				if errors.Is(err, tt.wantErr) {
					require.Equal(t, uuid.Nil, id)
				} else {
					// If expecting a pg error type (violated NOT NULL or similar)
					_, isPgErr := err.(*pgconn.PgError)
					_, isWantPgErr := tt.wantErr.(*pgconn.PgError)
					require.Equal(t, isWantPgErr, isPgErr, "expected pg error type mismatch")
				}
			}
		})
	}
}

func TestGetUserByID(t *testing.T) {
	// Seed DB
	ctx := context.Background()
	err := SeedDB(ctx, TestDSN, SeedPath)
	require.NoError(t, err)

	tests := []struct {
		name     string
		inputID  string
		wantUser *User
		wantErr  bool
	}{
		{
			name:    "success - existing admin user",
			inputID: "a1f92e18-1d8f-4f0f-9a4d-3b9e3b26b7b1",
			wantUser: &User{
				ID:         uuid.MustParse("a1f92e18-1d8f-4f0f-9a4d-3b9e3b26b7b1"),
				Username:   "anacardo",
				Email:      "anacardo@example.com",
				UserRole:   "admin",
				IsActive:   true,
				IsVerified: true,
			},
			wantErr: false,
		},
		{
			name:    "success - existing moderator user",
			inputID: "cfa53179-9085-4f33-86b3-5dc5f7a1465f",
			wantUser: &User{
				ID:         uuid.MustParse("cfa53179-9085-4f33-86b3-5dc5f7a1465f"),
				Username:   "moderata",
				Email:      "moderata@example.com",
				UserRole:   "moderator",
				IsActive:   true,
				IsVerified: true,
			},
			wantErr: false,
		},
		{
			name:     "fail - user not found",
			inputID:  "11111111-2222-3333-4444-555555555555",
			wantUser: nil,
			wantErr:  true,
		},
		{
			name:     "fail - invalid UUID format",
			inputID:  "not-a-valid-uuid",
			wantUser: nil,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			id, parseErr := uuid.Parse(tt.inputID)
			if parseErr != nil {
				// Expected parse errors (invalid UUID)
				if tt.wantErr {
					t.Logf("expected parse error: %v", parseErr)
					return
				}
				t.Fatalf("unexpected parse error: %v", parseErr)
			}

			user, err := Repo.GetUserByID(ctx, id)
			if tt.wantErr {
				require.Error(t, err)
				require.Nil(t, user)
				return
			}

			require.NoError(t, err)
			require.NotNil(t, user)
			require.Equal(t, tt.wantUser.Username, user.Username)
			require.Equal(t, tt.wantUser.Email, user.Email)
			require.Equal(t, tt.wantUser.UserRole, user.UserRole)
			require.Equal(t, tt.wantUser.IsActive, user.IsActive)
			require.Equal(t, tt.wantUser.IsVerified, user.IsVerified)
		})
	}
}

func TestGetUserByUserName(t *testing.T) {
	// Seed DB
	ctx := context.Background()
	err := SeedDB(ctx, TestDSN, SeedPath)
	require.NoError(t, err)

	tests := []struct {
		name          string
		inputUserName string
		wantUser      *User
		wantErr       bool
	}{
		{
			name:          "success - existing admin user",
			inputUserName: "anacardo",
			wantUser: &User{
				ID:         uuid.MustParse("a1f92e18-1d8f-4f0f-9a4d-3b9e3b26b7b1"),
				Username:   "anacardo",
				Email:      "anacardo@example.com",
				UserRole:   "admin",
				IsActive:   true,
				IsVerified: true,
			},
			wantErr: false,
		},
		{
			name:          "success - existing moderator user",
			inputUserName: "moderata",
			wantUser: &User{
				ID:         uuid.MustParse("cfa53179-9085-4f33-86b3-5dc5f7a1465f"),
				Username:   "moderata",
				Email:      "moderata@example.com",
				UserRole:   "moderator",
				IsActive:   true,
				IsVerified: true,
			},
			wantErr: false,
		},
		{
			name:          "fail - user not found",
			inputUserName: "nonexistentuser",
			wantUser:      nil,
			wantErr:       true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user, err := Repo.GetUserByUserName(ctx, tt.inputUserName)

			if tt.wantErr {
				require.Error(t, err)
				require.Nil(t, user)
				return
			}

			require.NoError(t, err)
			require.NotNil(t, user)
			require.Equal(t, tt.wantUser.ID, user.ID)
			require.Equal(t, tt.wantUser.Username, user.Username)
			require.Equal(t, tt.wantUser.Email, user.Email)
			require.Equal(t, tt.wantUser.UserRole, user.UserRole)
			require.Equal(t, tt.wantUser.IsActive, user.IsActive)
			require.Equal(t, tt.wantUser.IsVerified, user.IsVerified)
		})
	}
}

func TestGetUserByEmail(t *testing.T) {
	// Seed DB
	ctx := context.Background()
	err := SeedDB(ctx, TestDSN, SeedPath)
	require.NoError(t, err)

	tests := []struct {
		name       string
		inputEmail string
		wantUser   *User
		wantErr    bool
	}{
		{
			name:       "success - existing admin user",
			inputEmail: "anacardo@example.com",
			wantUser: &User{
				ID:         uuid.MustParse("a1f92e18-1d8f-4f0f-9a4d-3b9e3b26b7b1"),
				Username:   "anacardo",
				Email:      "anacardo@example.com",
				UserRole:   "admin",
				IsActive:   true,
				IsVerified: true,
			},
			wantErr: false,
		},
		{
			name:       "success - existing moderator user",
			inputEmail: "moderata@example.com",
			wantUser: &User{
				ID:         uuid.MustParse("cfa53179-9085-4f33-86b3-5dc5f7a1465f"),
				Username:   "moderata",
				Email:      "moderata@example.com",
				UserRole:   "moderator",
				IsActive:   true,
				IsVerified: true,
			},
			wantErr: false,
		},
		{
			name:       "fail - user not found",
			inputEmail: "nonexistent@example.com",
			wantUser:   nil,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user, err := Repo.GetUserByEmail(ctx, tt.inputEmail)

			if tt.wantErr {
				require.Error(t, err)
				require.Nil(t, user)
				return
			}

			require.NoError(t, err)
			require.NotNil(t, user)
			require.Equal(t, tt.wantUser.ID, user.ID)
			require.Equal(t, tt.wantUser.Username, user.Username)
			require.Equal(t, tt.wantUser.Email, user.Email)
			require.Equal(t, tt.wantUser.UserRole, user.UserRole)
			require.Equal(t, tt.wantUser.IsActive, user.IsActive)
			require.Equal(t, tt.wantUser.IsVerified, user.IsVerified)
		})
	}
}

func TestSearchUsersByUserName(t *testing.T) {
	// Seed DB
	ctx := context.Background()
	err := SeedDB(ctx, TestDSN, SeedPath)
	require.NoError(t, err)

	tests := []struct {
		name          string
		inputPattern  string
		wantUsernames []string
		wantErr       bool
	}{
		{
			name:          "match single user - exact",
			inputPattern:  "anacardo",
			wantUsernames: []string{"anacardo"},
			wantErr:       false,
		},
		{
			name:          "match multiple users - partial",
			inputPattern:  "a",
			wantUsernames: []string{"anacardo", "moderata", "inactiveuser"},
			wantErr:       false,
		},
		{
			name:          "match no users",
			inputPattern:  "nonexistent",
			wantUsernames: []string{},
			wantErr:       false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			users, err := Repo.SearchUsersByUserName(ctx, tt.inputPattern)

			if tt.wantErr {
				require.Error(t, err)
				require.Nil(t, users)
				return
			}

			require.NoError(t, err)
			require.Len(t, users, len(tt.wantUsernames))

			gotUsernames := make([]string, 0, len(users))
			for _, u := range users {
				gotUsernames = append(gotUsernames, u.Username)
				require.NotEqual(t, uuid.Nil, u.ID)
			}

			// Check that all expected usernames are present
			for _, expected := range tt.wantUsernames {
				require.Contains(t, gotUsernames, expected)
			}
		})
	}
}

func TestSearchUsersByDisplayName(t *testing.T) {
	// Seed DB
	ctx := context.Background()
	err := SeedDB(ctx, TestDSN, SeedPath)
	require.NoError(t, err)

	tests := []struct {
		name             string
		inputPattern     string
		wantDisplayNames []string
		wantErr          bool
	}{
		{
			name:             "match single display name - exact",
			inputPattern:     "Anacardo",
			wantDisplayNames: []string{"Anacardo"},
			wantErr:          false,
		},
		{
			name:             "match multiple display names - partial",
			inputPattern:     "a",
			wantDisplayNames: []string{"Anacardo", "Moderata Silva", "Inactive User"},
			wantErr:          false,
		},
		{
			name:             "match no display names",
			inputPattern:     "Nonexistent",
			wantDisplayNames: []string{},
			wantErr:          false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			users, err := Repo.SearchUsersByDisplayName(ctx, tt.inputPattern)

			if tt.wantErr {
				require.Error(t, err)
				require.Nil(t, users)
				return
			}

			require.NoError(t, err)
			require.Len(t, users, len(tt.wantDisplayNames))

			gotDisplayNames := make([]string, 0, len(users))
			for _, u := range users {
				gotDisplayNames = append(gotDisplayNames, u.DisplayName)
				require.NotEqual(t, uuid.Nil, u.ID)
			}

			// Check that all expected display names are present
			for _, expected := range tt.wantDisplayNames {
				found := false
				for _, got := range gotDisplayNames {
					if got == expected {
						found = true
						break
					}
				}
				require.True(t, found, "expected display name %s not found", expected)
			}
		})
	}
}

func TestSetUserActive(t *testing.T) {
	// Seed DB
	ctx := context.Background()
	err := SeedDB(ctx, TestDSN, SeedPath)
	require.NoError(t, err)

	tests := []struct {
		name       string
		username   string
		wantErr    bool
		expectTrue bool
	}{
		{
			name:       "activate existing inactive user",
			username:   "inactiveuser",
			wantErr:    false,
			expectTrue: true,
		},
		{
			name:       "activate already active user",
			username:   "anacardo",
			wantErr:    false,
			expectTrue: true,
		},
		{
			name:       "non-existent user",
			username:   "unknownuser",
			wantErr:    false, // Update won't fail, just affects 0 rows
			expectTrue: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Repo.SetUserActive(ctx, tt.username)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)

			// Verify the is_active field in the database
			user, getErr := Repo.GetUserByUserName(ctx, tt.username)
			if tt.expectTrue {
				require.NoError(t, getErr)
				require.True(t, user.IsActive)
			} else {
				// If user does not exist, GetUserByUserName returns error
				require.Error(t, getErr)
			}
		})
	}
}

func TestSetNewPassword(t *testing.T) {
	// Seed DB
	ctx := context.Background()
	err := SeedDB(ctx, TestDSN, SeedPath)
	require.NoError(t, err)

	tests := []struct {
		name       string
		userID     string
		newPass    string
		wantErr    bool
		verifyHash string
	}{
		{
			name:       "success - update password for admin",
			userID:     "a1f92e18-1d8f-4f0f-9a4d-3b9e3b26b7b1",
			newPass:    "new_hashed_password_123",
			wantErr:    false,
			verifyHash: "new_hashed_password_123",
		},
		{
			name:    "fail - non-existent user",
			userID:  "11111111-2222-3333-4444-555555555555",
			newPass: "whatever",
			wantErr: false,
		},
		{
			name:    "fail - invalid UUID",
			userID:  "not-a-uuid",
			newPass: "irrelevant",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userUUID, err := uuid.Parse(tt.userID)
			if err != nil {
				if tt.wantErr {
					t.Logf("expected parse error: %v", err)
					return
				}
				t.Fatalf("unexpected UUID parse error: %v", err)
			}

			err = Repo.SetNewPassword(ctx, userUUID, tt.newPass)
			if tt.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)

			// Verify password hash in DB
			user, err := Repo.GetUserByID(ctx, userUUID)
			if tt.name == "fail - non-existent user" {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.Equal(t, tt.verifyHash, user.PasswordHash)
		})
	}
}

func TestUpdateProfilePic(t *testing.T) {
	// Seed DB
	ctx := context.Background()
	err := SeedDB(ctx, TestDSN, SeedPath)
	require.NoError(t, err)

	tests := []struct {
		name       string
		username   string
		newProfile string
		wantErr    bool
	}{
		{
			name:       "success - update existing user",
			username:   "anacardo",
			newProfile: "https://example.com/newpic.jpg",
			wantErr:    false,
		},
		{
			name:       "success - update another user",
			username:   "soccerpunk",
			newProfile: "/avatars/soccerpunk.png",
			wantErr:    false,
		},
		{
			name:       "fail - non-existent user",
			username:   "ghostuser",
			newProfile: "/avatars/ghost.png",
			wantErr:    false, // still not an error in SQL, just zero rows affected
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Repo.UpdateProfilePic(ctx, tt.username, tt.newProfile)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)

			// Verify the profile pic in DB for existing users
			user, err := Repo.GetUserByUserName(ctx, tt.username)
			if err != nil {
				// If the user doesn’t exist, QueryRow will return an error
				require.Equal(t, tt.username == "ghostuser", true)
				return
			}
			require.Equal(t, tt.newProfile, user.ProfilePic)
		})
	}
}

func TestFollowUser(t *testing.T) {
	// Seed DB
	ctx := context.Background()
	err := SeedDB(ctx, TestDSN, SeedPath)
	require.NoError(t, err)

	// DB for test query
	db, err := testutils.ConnectDB(ctx, TestDSN)
	require.NoError(t, err)

	// UUIDs from seed
	anacardoID := uuid.MustParse("a1f92e18-1d8f-4f0f-9a4d-3b9e3b26b7b1")
	moderataID := uuid.MustParse("cfa53179-9085-4f33-86b3-5dc5f7a1465f")

	tests := []struct {
		name         string
		followerID   uuid.UUID
		followedUser string
		wantErr      bool
	}{
		{
			name:         "success - new follow request",
			followerID:   moderataID,
			followedUser: "anacardo",
			wantErr:      false,
		},
		{
			name:         "success - duplicate follow (DO NOTHING)",
			followerID:   moderataID,
			followedUser: "anacardo",
			wantErr:      false,
		},
		{
			name:         "fail - non-existent followed user",
			followerID:   anacardoID,
			followedUser: "ghostuser",
			wantErr:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Repo.FollowUser(ctx, tt.followerID, tt.followedUser)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)

			// Verify follow exists in DB if followed user exists
			if tt.followedUser != "ghostuser" {
				query := `
				SELECT follow_status
				FROM follows f
				JOIN users u ON u.id = f.followed_id
				WHERE f.follower_id = $1 AND u.username = $2
				;`
				var status string
				err := db.QueryRow(ctx, query, tt.followerID, tt.followedUser).Scan(&status)
				require.NoError(t, err)
				require.Equal(t, "pending", status)
			}
		})
	}
}

func TestAcceptFollow(t *testing.T) {
	// Seed DB
	ctx := context.Background()
	err := SeedDB(ctx, TestDSN, SeedPath)
	require.NoError(t, err)

	// DB for test query
	db, err := testutils.ConnectDB(ctx, TestDSN)
	require.NoError(t, err)

	tests := []struct {
		name            string
		followerName    string
		followedName    string
		wantFinalStatus string
		wantErr         bool
	}{
		{
			name:            "success - accept pending follow",
			followerName:    "moderata",
			followedName:    "soccerpunk",
			wantFinalStatus: "accepted",
			wantErr:         false,
		},
		{
			name:            "success - accept already accepted follow (idempotent)",
			followerName:    "anacardo",
			followedName:    "moderata",
			wantFinalStatus: "accepted",
			wantErr:         false,
		},
		{
			name:            "fail - non-existent follow relation",
			followerName:    "soccerpunk",
			followedName:    "moderata",
			wantFinalStatus: "",
			wantErr:         false, // SQL executes but zero rows updated
		},
		{
			name:            "fail - non-existent user",
			followerName:    "ghostuser",
			followedName:    "anacardo",
			wantFinalStatus: "",
			wantErr:         false, // SQL executes but zero rows updated
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Repo.AcceptFollow(ctx, tt.followerName, tt.followedName)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)

			// Verify the follow status if both users exist
			query := `
			SELECT f.follow_status
			FROM follows f
			JOIN users uf ON uf.id = f.follower_id
			JOIN users ut ON ut.id = f.followed_id
			WHERE uf.username = $1 AND ut.username = $2
			;`
			var status string
			err = db.QueryRow(ctx, query, tt.followerName, tt.followedName).Scan(&status)
			if tt.wantFinalStatus == "" {
				require.Error(t, err) // no row found
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.wantFinalStatus, status)
			}
		})
	}
}

func TestUnfollowUser(t *testing.T) {
	// Seed DB
	ctx := context.Background()
	err := SeedDB(ctx, TestDSN, SeedPath)
	require.NoError(t, err)

	// DB for test query
	db, err := testutils.ConnectDB(ctx, TestDSN)
	require.NoError(t, err)

	tests := []struct {
		name         string
		followerName string
		followedName string
		wantErr      bool
	}{
		{
			name:         "success - delete existing follow",
			followerName: "anacardo",
			followedName: "moderata",
			wantErr:      false,
		},
		{
			name:         "success - delete pending follow",
			followerName: "moderata",
			followedName: "soccerpunk",
			wantErr:      false,
		},
		{
			name:         "success - delete non-existent follow",
			followerName: "soccerpunk",
			followedName: "moderata",
			wantErr:      false, // SQL executes but zero rows deleted
		},
		{
			name:         "success - non-existent user",
			followerName: "ghostuser",
			followedName: "anacardo",
			wantErr:      false, // SQL executes but zero rows deleted
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Repo.UnfollowUser(ctx, tt.followerName, tt.followedName)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)

			// Verify follow does not exist anymore
			query := `
			SELECT COUNT(*)
			FROM follows f
			JOIN users uf ON uf.id = f.follower_id
			JOIN users ut ON ut.id = f.followed_id
			WHERE uf.username = $1 AND ut.username = $2
			;`
			var count int
			err = db.QueryRow(ctx, query, tt.followerName, tt.followedName).Scan(&count)
			require.NoError(t, err)
			require.Equal(t, 0, count)
		})
	}
}

func TestGetUserFollows(t *testing.T) {
	// Seed DB
	ctx := context.Background()
	err := SeedDB(ctx, TestDSN, SeedPath)
	require.NoError(t, err)

	// UUIDs from seed
	anacardoID := uuid.MustParse("a1f92e18-1d8f-4f0f-9a4d-3b9e3b26b7b1")
	moderataID := uuid.MustParse("cfa53179-9085-4f33-86b3-5dc5f7a1465f")
	soccerpunkID := uuid.MustParse("f7a92b5b-4b7e-4787-9c0b-2b0b6cb86b4e")

	tests := []struct {
		name       string
		followerID uuid.UUID
		followedID uuid.UUID
		wantStatus string
		wantNil    bool
	}{
		{
			name:       "existing accepted follow",
			followerID: anacardoID,
			followedID: moderataID,
			wantStatus: "accepted",
			wantNil:    false,
		},
		{
			name:       "existing pending follow",
			followerID: moderataID,
			followedID: soccerpunkID,
			wantStatus: "pending",
			wantNil:    false,
		},
		{
			name:       "no follow exists",
			followerID: soccerpunkID,
			followedID: moderataID,
			wantNil:    true,
		},
		{
			name:       "non-existent users",
			followerID: uuid.New(),
			followedID: uuid.New(),
			wantNil:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			follow, err := Repo.GetUserFollows(ctx, tt.followerID, tt.followedID)
			require.NoError(t, err)
			if tt.wantNil {
				require.Nil(t, follow)
				return
			}
			require.NotNil(t, follow)
			require.Equal(t, tt.wantStatus, follow.FollowStatus)
			require.Equal(t, tt.followerID, follow.FollowerID)
			require.Equal(t, tt.followedID, follow.FollowedID)
		})
	}
}

func TestGetFollowers(t *testing.T) {
	// Seed DB
	ctx := context.Background()
	err := SeedDB(ctx, TestDSN, SeedPath)
	require.NoError(t, err)

	// UUIDs from seed
	anacardoID := uuid.MustParse("a1f92e18-1d8f-4f0f-9a4d-3b9e3b26b7b1")
	moderataID := uuid.MustParse("cfa53179-9085-4f33-86b3-5dc5f7a1465f")
	soccerpunkID := uuid.MustParse("f7a92b5b-4b7e-4787-9c0b-2b0b6cb86b4e")

	tests := []struct {
		name       string
		followedID uuid.UUID
		wantCount  int
		wantIDs    []uuid.UUID
	}{
		{
			name:       "followers of moderata",
			followedID: moderataID,
			wantCount:  1,
			wantIDs:    []uuid.UUID{anacardoID},
		},
		{
			name:       "followers of anacardo",
			followedID: anacardoID,
			wantCount:  1,
			wantIDs:    []uuid.UUID{soccerpunkID},
		},
		{
			name:       "followers of soccerpunk",
			followedID: soccerpunkID,
			wantCount:  0,
			wantIDs:    []uuid.UUID{},
		},
		{
			name:       "non-existent user",
			followedID: uuid.New(),
			wantCount:  0,
			wantIDs:    []uuid.UUID{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			follows, err := Repo.GetFollowers(ctx, tt.followedID)
			require.NoError(t, err)
			require.Len(t, follows, tt.wantCount)

			for i, f := range follows {
				require.Equal(t, tt.wantIDs[i], f.FollowerID)
				require.Equal(t, tt.followedID, f.FollowedID)
				require.Equal(t, "accepted", f.FollowStatus)
			}
		})
	}
}

func TestGetFollowing(t *testing.T) {
	// Seed DB
	ctx := context.Background()
	err := SeedDB(ctx, TestDSN, SeedPath)
	require.NoError(t, err)

	// UUIDs from seed
	anacardoID := uuid.MustParse("a1f92e18-1d8f-4f0f-9a4d-3b9e3b26b7b1")
	moderataID := uuid.MustParse("cfa53179-9085-4f33-86b3-5dc5f7a1465f")
	soccerpunkID := uuid.MustParse("f7a92b5b-4b7e-4787-9c0b-2b0b6cb86b4e")

	tests := []struct {
		name       string
		followerID uuid.UUID
		wantCount  int
		wantIDs    []uuid.UUID
	}{
		{
			name:       "anacardo following",
			followerID: anacardoID,
			wantCount:  1,
			wantIDs:    []uuid.UUID{moderataID}, // from seed, accepted
		},
		{
			name:       "moderata following",
			followerID: moderataID,
			wantCount:  0, // pending, not accepted
			wantIDs:    []uuid.UUID{},
		},
		{
			name:       "soccerpunk following",
			followerID: soccerpunkID,
			wantCount:  1,
			wantIDs:    []uuid.UUID{anacardoID}, // accepted
		},
		{
			name:       "non-existent follower",
			followerID: uuid.New(),
			wantCount:  0,
			wantIDs:    []uuid.UUID{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			follows, err := Repo.GetFollowing(ctx, tt.followerID)
			require.NoError(t, err)
			require.Len(t, follows, tt.wantCount)

			for i, f := range follows {
				require.Equal(t, tt.followerID, f.FollowerID)
				require.Equal(t, tt.wantIDs[i], f.FollowedID)
				require.Equal(t, "accepted", f.FollowStatus)
			}
		})
	}
}
