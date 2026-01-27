package postgres

import (
	"context"
	"database/sql"
	"testing"

	"github.com/m19cmjigen/sandbox-project-management/backend/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUserRepository_Create(t *testing.T) {
	db := GetTestDB(t)
	if db == nil {
		return
	}
	defer db.Close()
	defer CleanupTestDB(t, db, "users")

	repo := NewUserRepository(db)
	ctx := context.Background()

	t.Run("Create user", func(t *testing.T) {
		user := &domain.User{
			Username:     "testuser",
			PasswordHash: "hashed_password",
			Email:        "test@example.com",
			FullName:     sql.NullString{String: "Test User", Valid: true},
			Role:         domain.RoleViewer,
			IsActive:     true,
		}

		err := repo.Create(ctx, user)
		require.NoError(t, err)
		assert.NotZero(t, user.ID)
		assert.NotZero(t, user.CreatedAt)
	})

	t.Run("Create user with duplicate username should fail", func(t *testing.T) {
		user1 := &domain.User{
			Username:     "duplicate",
			PasswordHash: "hash1",
			Email:        "user1@example.com",
			Role:         domain.RoleViewer,
			IsActive:     true,
		}
		require.NoError(t, repo.Create(ctx, user1))

		user2 := &domain.User{
			Username:     "duplicate",
			PasswordHash: "hash2",
			Email:        "user2@example.com",
			Role:         domain.RoleViewer,
			IsActive:     true,
		}
		err := repo.Create(ctx, user2)
		assert.Error(t, err)
	})
}

func TestUserRepository_GetByUsername(t *testing.T) {
	db := GetTestDB(t)
	if db == nil {
		return
	}
	defer db.Close()
	defer CleanupTestDB(t, db, "users")

	repo := NewUserRepository(db)
	ctx := context.Background()

	t.Run("Get existing user", func(t *testing.T) {
		user := &domain.User{
			Username:     "findme",
			PasswordHash: "hashed",
			Email:        "findme@example.com",
			Role:         domain.RoleViewer,
			IsActive:     true,
		}
		require.NoError(t, repo.Create(ctx, user))

		found, err := repo.GetByUsername(ctx, "findme")
		require.NoError(t, err)
		assert.Equal(t, user.Username, found.Username)
		assert.Equal(t, user.Email, found.Email)
	})

	t.Run("Get non-existent user", func(t *testing.T) {
		_, err := repo.GetByUsername(ctx, "nonexistent")
		assert.Error(t, err)
	})
}

func TestUserRepository_GetByID(t *testing.T) {
	db := GetTestDB(t)
	if db == nil {
		return
	}
	defer db.Close()
	defer CleanupTestDB(t, db, "users")

	repo := NewUserRepository(db)
	ctx := context.Background()

	t.Run("Get existing user", func(t *testing.T) {
		user := &domain.User{
			Username:     "testuser",
			PasswordHash: "hashed",
			Email:        "test@example.com",
			Role:         domain.RoleViewer,
			IsActive:     true,
		}
		require.NoError(t, repo.Create(ctx, user))

		found, err := repo.GetByID(ctx, user.ID)
		require.NoError(t, err)
		assert.Equal(t, user.Username, found.Username)
	})

	t.Run("Get non-existent user", func(t *testing.T) {
		_, err := repo.GetByID(ctx, 99999)
		assert.Error(t, err)
	})
}

func TestUserRepository_List(t *testing.T) {
	db := GetTestDB(t)
	if db == nil {
		return
	}
	defer db.Close()
	defer CleanupTestDB(t, db, "users")

	repo := NewUserRepository(db)
	ctx := context.Background()

	// Create test users
	users := []*domain.User{
		{Username: "user1", PasswordHash: "hash1", Email: "user1@example.com", Role: domain.RoleViewer, IsActive: true},
		{Username: "user2", PasswordHash: "hash2", Email: "user2@example.com", Role: domain.RoleManager, IsActive: true},
		{Username: "user3", PasswordHash: "hash3", Email: "user3@example.com", Role: domain.RoleAdmin, IsActive: false},
	}

	for _, user := range users {
		require.NoError(t, repo.Create(ctx, user))
	}

	allUsers, err := repo.List(ctx, nil)
	require.NoError(t, err)
	assert.GreaterOrEqual(t, len(allUsers), 3)
}

func TestUserRepository_Update(t *testing.T) {
	db := GetTestDB(t)
	if db == nil {
		return
	}
	defer db.Close()
	defer CleanupTestDB(t, db, "users")

	repo := NewUserRepository(db)
	ctx := context.Background()

	t.Run("Update user", func(t *testing.T) {
		user := &domain.User{
			Username:     "updateme",
			PasswordHash: "original_hash",
			Email:        "original@example.com",
			FullName:     sql.NullString{String: "Original Name", Valid: true},
			Role:         domain.RoleViewer,
			IsActive:     true,
		}
		require.NoError(t, repo.Create(ctx, user))

		user.Email = "updated@example.com"
		user.FullName = sql.NullString{String: "Updated Name", Valid: true}
		user.Role = domain.RoleManager

		err := repo.Update(ctx, user)
		require.NoError(t, err)

		found, err := repo.GetByID(ctx, user.ID)
		require.NoError(t, err)
		assert.Equal(t, "updated@example.com", found.Email)
		assert.Equal(t, "Updated Name", found.FullName.String)
		assert.Equal(t, domain.RoleManager, found.Role)
	})
}

func TestUserRepository_Delete(t *testing.T) {
	db := GetTestDB(t)
	if db == nil {
		return
	}
	defer db.Close()
	defer CleanupTestDB(t, db, "users")

	repo := NewUserRepository(db)
	ctx := context.Background()

	t.Run("Delete existing user", func(t *testing.T) {
		user := &domain.User{
			Username:     "deleteme",
			PasswordHash: "hash",
			Email:        "delete@example.com",
			Role:         domain.RoleViewer,
			IsActive:     true,
		}
		require.NoError(t, repo.Create(ctx, user))

		err := repo.Delete(ctx, user.ID)
		require.NoError(t, err)

		_, err = repo.GetByID(ctx, user.ID)
		assert.Error(t, err)
	})
}

