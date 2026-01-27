package postgres

import (
	"context"
	"testing"

	"github.com/m19cmjigen/sandbox-project-management/backend/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestProjectRepository_Create(t *testing.T) {
	db := GetTestDB(t)
	if db == nil {
		return
	}
	defer db.Close()
	defer CleanupTestDB(t, db, "projects", "organizations")

	repo := NewProjectRepository(db)
	orgRepo := NewOrganizationRepository(db)
	ctx := context.Background()

	// Create test organization
	org := &domain.Organization{Name: "Test Org", ParentID: nil}
	require.NoError(t, orgRepo.Create(ctx, org))

	t.Run("Create project", func(t *testing.T) {
		orgID := org.ID
		leadID := "user123"
		project := &domain.Project{
			JiraProjectID:  "10001",
			Key:            "TEST",
			Name:           "Test Project",
			OrganizationID: &orgID,
			LeadAccountID:  &leadID,
		}

		err := repo.Create(ctx, project)
		require.NoError(t, err)
		assert.NotZero(t, project.ID)
		assert.NotZero(t, project.CreatedAt)
	})

	t.Run("Create project without organization", func(t *testing.T) {
		project := &domain.Project{
			JiraProjectID: "10002",
			Key:           "TEST2",
			Name:          "Test Project 2",
		}

		err := repo.Create(ctx, project)
		require.NoError(t, err)
		assert.NotZero(t, project.ID)
	})
}

func TestProjectRepository_FindAll(t *testing.T) {
	db := GetTestDB(t)
	if db == nil {
		return
	}
	defer db.Close()
	defer CleanupTestDB(t, db, "projects")

	repo := NewProjectRepository(db)
	ctx := context.Background()

	// Create test projects
	p1 := &domain.Project{JiraProjectID: "10001", Key: "P1", Name: "Project 1"}
	p2 := &domain.Project{JiraProjectID: "10002", Key: "P2", Name: "Project 2"}
	require.NoError(t, repo.Create(ctx, p1))
	require.NoError(t, repo.Create(ctx, p2))

	projects, err := repo.FindAll(ctx)
	require.NoError(t, err)
	assert.GreaterOrEqual(t, len(projects), 2)
}

func TestProjectRepository_FindByID(t *testing.T) {
	db := GetTestDB(t)
	if db == nil {
		return
	}
	defer db.Close()
	defer CleanupTestDB(t, db, "projects")

	repo := NewProjectRepository(db)
	ctx := context.Background()

	t.Run("Find existing project", func(t *testing.T) {
		project := &domain.Project{
			JiraProjectID: "10001",
			Key:           "TEST",
			Name:          "Test Project",
		}
		require.NoError(t, repo.Create(ctx, project))

		found, err := repo.FindByID(ctx, project.ID)
		require.NoError(t, err)
		assert.Equal(t, project.Key, found.Key)
		assert.Equal(t, project.Name, found.Name)
	})

	t.Run("Find non-existent project", func(t *testing.T) {
		_, err := repo.FindByID(ctx, 99999)
		assert.Error(t, err)
	})
}

func TestProjectRepository_Update(t *testing.T) {
	db := GetTestDB(t)
	if db == nil {
		return
	}
	defer db.Close()
	defer CleanupTestDB(t, db, "projects")

	repo := NewProjectRepository(db)
	ctx := context.Background()

	t.Run("Update project", func(t *testing.T) {
		project := &domain.Project{
			JiraProjectID: "10001",
			Key:           "TEST",
			Name:          "Original Name",
		}
		require.NoError(t, repo.Create(ctx, project))

		project.Name = "Updated Name"
		err := repo.Update(ctx, project)
		require.NoError(t, err)

		found, err := repo.FindByID(ctx, project.ID)
		require.NoError(t, err)
		assert.Equal(t, "Updated Name", found.Name)
	})
}

func TestProjectRepository_FindByOrganizationID(t *testing.T) {
	db := GetTestDB(t)
	if db == nil {
		return
	}
	defer db.Close()
	defer CleanupTestDB(t, db, "projects", "organizations")

	repo := NewProjectRepository(db)
	orgRepo := NewOrganizationRepository(db)
	ctx := context.Background()

	// Create organizations
	org1 := &domain.Organization{Name: "Org 1", ParentID: nil}
	org2 := &domain.Organization{Name: "Org 2", ParentID: nil}
	require.NoError(t, orgRepo.Create(ctx, org1))
	require.NoError(t, orgRepo.Create(ctx, org2))

	// Create projects
	orgID1 := org1.ID
	orgID2 := org2.ID
	p1 := &domain.Project{
		JiraProjectID:  "10001",
		Key:            "P1",
		Name:           "Project 1",
		OrganizationID: &orgID1,
	}
	p2 := &domain.Project{
		JiraProjectID:  "10002",
		Key:            "P2",
		Name:           "Project 2",
		OrganizationID: &orgID1,
	}
	p3 := &domain.Project{
		JiraProjectID:  "10003",
		Key:            "P3",
		Name:           "Project 3",
		OrganizationID: &orgID2,
	}
	require.NoError(t, repo.Create(ctx, p1))
	require.NoError(t, repo.Create(ctx, p2))
	require.NoError(t, repo.Create(ctx, p3))

	projects, err := repo.FindByOrganizationID(ctx, org1.ID)
	require.NoError(t, err)
	assert.Len(t, projects, 2)
}
