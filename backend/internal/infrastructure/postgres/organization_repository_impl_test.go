package postgres

import (
	"context"
	"testing"

	"github.com/m19cmjigen/sandbox-project-management/backend/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOrganizationRepository_Create(t *testing.T) {
	db := GetTestDB(t)
	if db == nil {
		return
	}
	defer db.Close()
	defer CleanupTestDB(t, db, "organizations")

	repo := NewOrganizationRepository(db)
	ctx := context.Background()

	t.Run("Create root organization", func(t *testing.T) {
		org := &domain.Organization{
			Name:     "Test Company",
			ParentID: nil,
		}

		err := repo.Create(ctx, org)
		require.NoError(t, err)
		assert.NotZero(t, org.ID, "ID should be set after creation")
		assert.NotZero(t, org.CreatedAt, "CreatedAt should be set")
		assert.NotZero(t, org.UpdatedAt, "UpdatedAt should be set")
	})

	t.Run("Create child organization", func(t *testing.T) {
		// Create parent
		parent := &domain.Organization{
			Name:     "Parent Org",
			ParentID: nil,
		}
		err := repo.Create(ctx, parent)
		require.NoError(t, err)

		// Create child
		child := &domain.Organization{
			Name:     "Child Org",
			ParentID: &parent.ID,
		}
		err = repo.Create(ctx, child)
		require.NoError(t, err)
		assert.NotZero(t, child.ID)
		assert.Equal(t, parent.ID, *child.ParentID)
	})
}

func TestOrganizationRepository_FindAll(t *testing.T) {
	db := GetTestDB(t)
	if db == nil {
		return
	}
	defer db.Close()
	defer CleanupTestDB(t, db, "organizations")

	repo := NewOrganizationRepository(db)
	ctx := context.Background()

	// Create test data
	org1 := &domain.Organization{Name: "Org 1", ParentID: nil}
	org2 := &domain.Organization{Name: "Org 2", ParentID: nil}
	require.NoError(t, repo.Create(ctx, org1))
	require.NoError(t, repo.Create(ctx, org2))

	orgs, err := repo.FindAll(ctx)
	require.NoError(t, err)
	assert.GreaterOrEqual(t, len(orgs), 2, "Should have at least 2 organizations")
}

func TestOrganizationRepository_FindByID(t *testing.T) {
	db := GetTestDB(t)
	if db == nil {
		return
	}
	defer db.Close()
	defer CleanupTestDB(t, db, "organizations")

	repo := NewOrganizationRepository(db)
	ctx := context.Background()

	t.Run("Find existing organization", func(t *testing.T) {
		org := &domain.Organization{
			Name:     "Test Org",
			ParentID: nil,
		}
		err := repo.Create(ctx, org)
		require.NoError(t, err)

		found, err := repo.FindByID(ctx, org.ID)
		require.NoError(t, err)
		assert.Equal(t, org.ID, found.ID)
		assert.Equal(t, org.Name, found.Name)
	})

	t.Run("Find non-existent organization", func(t *testing.T) {
		_, err := repo.FindByID(ctx, 99999)
		assert.Error(t, err)
	})
}

func TestOrganizationRepository_Update(t *testing.T) {
	db := GetTestDB(t)
	if db == nil {
		return
	}
	defer db.Close()
	defer CleanupTestDB(t, db, "organizations")

	repo := NewOrganizationRepository(db)
	ctx := context.Background()

	t.Run("Update organization name", func(t *testing.T) {
		org := &domain.Organization{
			Name:     "Original Name",
			ParentID: nil,
		}
		err := repo.Create(ctx, org)
		require.NoError(t, err)

		org.Name = "Updated Name"
		err = repo.Update(ctx, org)
		require.NoError(t, err)

		found, err := repo.FindByID(ctx, org.ID)
		require.NoError(t, err)
		assert.Equal(t, "Updated Name", found.Name)
	})
}

func TestOrganizationRepository_Delete(t *testing.T) {
	db := GetTestDB(t)
	if db == nil {
		return
	}
	defer db.Close()
	defer CleanupTestDB(t, db, "organizations")

	repo := NewOrganizationRepository(db)
	ctx := context.Background()

	t.Run("Delete existing organization", func(t *testing.T) {
		org := &domain.Organization{
			Name:     "To Delete",
			ParentID: nil,
		}
		err := repo.Create(ctx, org)
		require.NoError(t, err)

		err = repo.Delete(ctx, org.ID)
		require.NoError(t, err)

		_, err = repo.FindByID(ctx, org.ID)
		assert.Error(t, err)
	})
}
