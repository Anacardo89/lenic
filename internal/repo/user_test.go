package repo

import (
	"context"
	"errors"
	"path/filepath"
	"testing"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stretchr/testify/require"
)

func TestCreateUser(t *testing.T) {
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
			id, err := repo.CreateUser(ctx, tt.input)

			if tt.wantErr == nil {
				require.NoError(t, err)
				require.NotEqual(t, uuid.Nil, id)
				// Verify user exists in DB
				user, err := repo.GetUserByID(ctx, id)
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

			user, err := repo.GetUserByID(ctx, id)
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
			user, err := repo.GetUserByUserName(ctx, tt.inputUserName)

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
			user, err := repo.GetUserByEmail(ctx, tt.inputEmail)

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
			users, err := repo.SearchUsersByUserName(ctx, tt.inputPattern)

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
			users, err := repo.SearchUsersByDisplayName(ctx, tt.inputPattern)

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
			err := repo.SetUserActive(ctx, tt.username)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)

			// Verify the is_active field in the database
			user, getErr := repo.GetUserByUserName(ctx, tt.username)
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

			err = repo.SetNewPassword(ctx, userUUID, tt.newPass)
			if tt.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)

			// Verify password hash in DB
			user, err := repo.GetUserByID(ctx, userUUID)
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
			err := repo.UpdateProfilePic(ctx, tt.username, tt.newProfile)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)

			// Verify the profile pic in DB for existing users
			user, err := repo.GetUserByUserName(ctx, tt.username)
			if err != nil {
				// If the user doesn’t exist, QueryRow will return an error
				require.Equal(t, tt.username == "ghostuser", true)
				return
			}
			require.Equal(t, tt.newProfile, user.ProfilePic)
		})
	}
}
