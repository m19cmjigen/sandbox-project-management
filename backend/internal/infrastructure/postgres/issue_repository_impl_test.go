package postgres

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/m19cmjigen/sandbox-project-management/backend/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIssueRepository_Create(t *testing.T) {
	db := GetTestDB(t)
	if db == nil {
		return
	}
	defer db.Close()
	defer CleanupTestDB(t, db, "issues", "projects")

	issueRepo := NewIssueRepository(db)
	projectRepo := NewProjectRepository(db)
	ctx := context.Background()

	// Create test project
	project := &domain.Project{
		JiraProjectID: "10001",
		Key:           "TEST",
		Name:          "Test Project",
	}
	require.NoError(t, projectRepo.Create(ctx, project))

	t.Run("Create issue", func(t *testing.T) {
		priority := "High"
		assignee := "John Doe"
		issue := &domain.Issue{
			JiraIssueID:    "10001",
			JiraIssueKey:   "TEST-1",
			ProjectID:      project.ID,
			Summary:        "Test Issue",
			Status:         "To Do",
			StatusCategory: domain.StatusCategoryToDo,
			Priority:       &priority,
			DueDate:        sql.NullTime{Time: time.Now().Add(7 * 24 * time.Hour), Valid: true},
			DelayStatus:    domain.DelayStatusGreen,
			AssigneeName:   &assignee,
		}

		err := issueRepo.Create(ctx, issue)
		require.NoError(t, err)
		assert.NotZero(t, issue.ID)
		assert.NotZero(t, issue.CreatedAt)
	})
}

func TestIssueRepository_FindAll(t *testing.T) {
	db := GetTestDB(t)
	if db == nil {
		return
	}
	defer db.Close()
	defer CleanupTestDB(t, db, "issues", "projects")

	issueRepo := NewIssueRepository(db)
	projectRepo := NewProjectRepository(db)
	ctx := context.Background()

	// Create test project
	project := &domain.Project{
		JiraProjectID: "10001",
		Key:           "TEST",
		Name:          "Test Project",
	}
	require.NoError(t, projectRepo.Create(ctx, project))

	// Create test issues
	i1 := &domain.Issue{
		JiraIssueID:    "10001",
		JiraIssueKey:   "TEST-1",
		ProjectID:      project.ID,
		Summary:        "Issue 1",
		Status:         "To Do",
		StatusCategory: domain.StatusCategoryToDo,
		DelayStatus:    domain.DelayStatusGreen,
	}
	i2 := &domain.Issue{
		JiraIssueID:    "10002",
		JiraIssueKey:   "TEST-2",
		ProjectID:      project.ID,
		Summary:        "Issue 2",
		Status:         "In Progress",
		StatusCategory: domain.StatusCategoryInProgress,
		DelayStatus:    domain.DelayStatusYellow,
	}
	require.NoError(t, issueRepo.Create(ctx, i1))
	require.NoError(t, issueRepo.Create(ctx, i2))

	issues, err := issueRepo.FindAll(ctx)
	require.NoError(t, err)
	assert.GreaterOrEqual(t, len(issues), 2)
}

func TestIssueRepository_FindByID(t *testing.T) {
	db := GetTestDB(t)
	if db == nil {
		return
	}
	defer db.Close()
	defer CleanupTestDB(t, db, "issues", "projects")

	issueRepo := NewIssueRepository(db)
	projectRepo := NewProjectRepository(db)
	ctx := context.Background()

	// Create test project
	project := &domain.Project{
		JiraProjectID: "10001",
		Key:           "TEST",
		Name:          "Test Project",
	}
	require.NoError(t, projectRepo.Create(ctx, project))

	t.Run("Find existing issue", func(t *testing.T) {
		issue := &domain.Issue{
			JiraIssueID:    "10001",
			JiraIssueKey:   "TEST-1",
			ProjectID:      project.ID,
			Summary:        "Test Issue",
			Status:         "To Do",
			StatusCategory: domain.StatusCategoryToDo,
			DelayStatus:    domain.DelayStatusGreen,
		}
		require.NoError(t, issueRepo.Create(ctx, issue))

		found, err := issueRepo.FindByID(ctx, issue.ID)
		require.NoError(t, err)
		assert.Equal(t, issue.JiraIssueKey, found.JiraIssueKey)
		assert.Equal(t, issue.Summary, found.Summary)
	})

	t.Run("Find non-existent issue", func(t *testing.T) {
		_, err := issueRepo.FindByID(ctx, 99999)
		assert.Error(t, err)
	})
}

func TestIssueRepository_Update(t *testing.T) {
	db := GetTestDB(t)
	if db == nil {
		return
	}
	defer db.Close()
	defer CleanupTestDB(t, db, "issues", "projects")

	issueRepo := NewIssueRepository(db)
	projectRepo := NewProjectRepository(db)
	ctx := context.Background()

	// Create test project
	project := &domain.Project{
		JiraProjectID: "10001",
		Key:           "TEST",
		Name:          "Test Project",
	}
	require.NoError(t, projectRepo.Create(ctx, project))

	t.Run("Update issue", func(t *testing.T) {
		issue := &domain.Issue{
			JiraIssueID:    "10001",
			JiraIssueKey:   "TEST-1",
			ProjectID:      project.ID,
			Summary:        "Original Summary",
			Status:         "To Do",
			StatusCategory: domain.StatusCategoryToDo,
			DelayStatus:    domain.DelayStatusGreen,
		}
		require.NoError(t, issueRepo.Create(ctx, issue))

		issue.Summary = "Updated Summary"
		issue.Status = "Done"
		issue.StatusCategory = domain.StatusCategoryDone

		err := issueRepo.Update(ctx, issue)
		require.NoError(t, err)

		found, err := issueRepo.FindByID(ctx, issue.ID)
		require.NoError(t, err)
		assert.Equal(t, "Updated Summary", found.Summary)
		assert.Equal(t, "Done", found.Status)
	})
}

func TestIssueRepository_FindByFilter(t *testing.T) {
	db := GetTestDB(t)
	if db == nil {
		return
	}
	defer db.Close()
	defer CleanupTestDB(t, db, "issues", "projects", "organizations")

	issueRepo := NewIssueRepository(db)
	projectRepo := NewProjectRepository(db)
	orgRepo := NewOrganizationRepository(db)
	ctx := context.Background()

	// Create test organization
	org := &domain.Organization{Name: "Test Org", ParentID: nil}
	require.NoError(t, orgRepo.Create(ctx, org))

	// Create test projects
	orgID := org.ID
	p1 := &domain.Project{
		JiraProjectID:  "10001",
		Key:            "P1",
		Name:           "Project 1",
		OrganizationID: &orgID,
	}
	p2 := &domain.Project{
		JiraProjectID: "10002",
		Key:           "P2",
		Name:          "Project 2",
	}
	require.NoError(t, projectRepo.Create(ctx, p1))
	require.NoError(t, projectRepo.Create(ctx, p2))

	// Create issues with different statuses
	issues := []*domain.Issue{
		{
			JiraIssueID:    "10001",
			JiraIssueKey:   "P1-1",
			ProjectID:      p1.ID,
			Summary:        "Red Issue",
			Status:         "In Progress",
			StatusCategory: domain.StatusCategoryInProgress,
			DelayStatus:    domain.DelayStatusRed,
		},
		{
			JiraIssueID:    "10002",
			JiraIssueKey:   "P1-2",
			ProjectID:      p1.ID,
			Summary:        "Yellow Issue",
			Status:         "In Progress",
			StatusCategory: domain.StatusCategoryInProgress,
			DelayStatus:    domain.DelayStatusYellow,
		},
		{
			JiraIssueID:    "10003",
			JiraIssueKey:   "P1-3",
			ProjectID:      p1.ID,
			Summary:        "Green Issue",
			Status:         "In Progress",
			StatusCategory: domain.StatusCategoryInProgress,
			DelayStatus:    domain.DelayStatusGreen,
		},
		{
			JiraIssueID:    "10004",
			JiraIssueKey:   "P2-1",
			ProjectID:      p2.ID,
			Summary:        "Project 2 Issue",
			Status:         "In Progress",
			StatusCategory: domain.StatusCategoryInProgress,
			DelayStatus:    domain.DelayStatusRed,
		},
	}

	for _, issue := range issues {
		require.NoError(t, issueRepo.Create(ctx, issue))
	}

	t.Run("Filter by project ID", func(t *testing.T) {
		filter := domain.IssueFilter{
			ProjectID: &p1.ID,
		}
		result, err := issueRepo.FindByFilter(ctx, filter)
		require.NoError(t, err)
		assert.Len(t, result, 3)
	})

	t.Run("Filter by delay status", func(t *testing.T) {
		delayStatus := domain.DelayStatusRed
		filter := domain.IssueFilter{
			DelayStatus: &delayStatus,
		}
		result, err := issueRepo.FindByFilter(ctx, filter)
		require.NoError(t, err)
		assert.GreaterOrEqual(t, len(result), 2)
	})

	t.Run("Combine filters", func(t *testing.T) {
		delayStatus := domain.DelayStatusRed
		filter := domain.IssueFilter{
			ProjectID:   &p1.ID,
			DelayStatus: &delayStatus,
		}
		result, err := issueRepo.FindByFilter(ctx, filter)
		require.NoError(t, err)
		assert.Len(t, result, 1)
	})
}

func TestIssueRepository_CountByProjectID(t *testing.T) {
	db := GetTestDB(t)
	if db == nil {
		return
	}
	defer db.Close()
	defer CleanupTestDB(t, db, "issues", "projects")

	issueRepo := NewIssueRepository(db)
	projectRepo := NewProjectRepository(db)
	ctx := context.Background()

	// Create test project
	project := &domain.Project{
		JiraProjectID: "10001",
		Key:           "TEST",
		Name:          "Test Project",
	}
	require.NoError(t, projectRepo.Create(ctx, project))

	// Create 5 issues
	for i := 0; i < 5; i++ {
		issue := &domain.Issue{
			JiraIssueID:    string(rune('1' + i)),
			JiraIssueKey:   "TEST-" + string(rune('1'+i)),
			ProjectID:      project.ID,
			Summary:        "Issue " + string(rune('1'+i)),
			Status:         "To Do",
			StatusCategory: domain.StatusCategoryToDo,
			DelayStatus:    domain.DelayStatusGreen,
		}
		require.NoError(t, issueRepo.Create(ctx, issue))
	}

	count, err := issueRepo.CountByProjectID(ctx, project.ID)
	require.NoError(t, err)
	assert.Equal(t, 5, count)
}
