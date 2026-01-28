package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/m19cmjigen/sandbox-project-management/backend/internal/domain"
	"github.com/stretchr/testify/assert"

	"github.com/stretchr/testify/require"
)

// MockProjectRepositoryForUsecase extends MockProjectRepository with additional methods
type MockProjectRepositoryForUsecase struct {
	MockProjectRepository
}

func (m *MockProjectRepositoryForUsecase) FindUnassigned(ctx context.Context) ([]domain.Project, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.Project), args.Error(1)
}

func (m *MockProjectRepositoryForUsecase) AssignToOrganization(ctx context.Context, projectID int64, organizationID *int64) error {
	args := m.Called(ctx, projectID, organizationID)
	return args.Error(0)
}

func setupProjectUsecase(t *testing.T) (*projectUsecase, *MockProjectRepositoryForUsecase) {
	mockRepo := new(MockProjectRepositoryForUsecase)
	usecase := &projectUsecase{
		projectRepo: mockRepo,
	}
	return usecase, mockRepo
}

// GetAll tests
func TestProjectUsecase_GetAll_Success(t *testing.T) {
	uc, mockRepo := setupProjectUsecase(t)
	ctx := context.Background()

	expectedProjects := []domain.Project{
		*createTestProject(1, "PROJ1", 1),
		*createTestProject(2, "PROJ2", 1),
		*createTestProject(3, "PROJ3", 2),
	}

	mockRepo.On("FindAll", ctx).Return(expectedProjects, nil)

	projects, err := uc.GetAll(ctx)

	require.NoError(t, err)
	assert.Len(t, projects, 3)
	mockRepo.AssertExpectations(t)
}

func TestProjectUsecase_GetAll_EmptyResult(t *testing.T) {
	uc, mockRepo := setupProjectUsecase(t)
	ctx := context.Background()

	mockRepo.On("FindAll", ctx).Return([]domain.Project{}, nil)

	projects, err := uc.GetAll(ctx)

	require.NoError(t, err)
	assert.Len(t, projects, 0)
	mockRepo.AssertExpectations(t)
}

func TestProjectUsecase_GetAll_RepositoryError(t *testing.T) {
	uc, mockRepo := setupProjectUsecase(t)
	ctx := context.Background()

	mockRepo.On("FindAll", ctx).Return(nil, errors.New("database error"))

	projects, err := uc.GetAll(ctx)

	assert.Nil(t, projects)
	assert.Error(t, err)
	mockRepo.AssertExpectations(t)
}

// GetAllWithStats tests
func TestProjectUsecase_GetAllWithStats_Success(t *testing.T) {
	uc, mockRepo := setupProjectUsecase(t)
	ctx := context.Background()

	expectedProjects := []domain.ProjectWithStats{
		createTestProjectWithStats(1, "PROJ1", 5, 3, 12),
		createTestProjectWithStats(2, "PROJ2", 0, 10, 5),
		createTestProjectWithStats(3, "PROJ3", 2, 0, 8),
	}

	mockRepo.On("FindAllWithStats", ctx).Return(expectedProjects, nil)

	projects, err := uc.GetAllWithStats(ctx)

	require.NoError(t, err)
	assert.Len(t, projects, 3)
	assert.Equal(t, 20, projects[0].TotalIssues)
	assert.Equal(t, 15, projects[1].TotalIssues)
	assert.Equal(t, 10, projects[2].TotalIssues)
	mockRepo.AssertExpectations(t)
}

func TestProjectUsecase_GetAllWithStats_EmptyResult(t *testing.T) {
	uc, mockRepo := setupProjectUsecase(t)
	ctx := context.Background()

	mockRepo.On("FindAllWithStats", ctx).Return([]domain.ProjectWithStats{}, nil)

	projects, err := uc.GetAllWithStats(ctx)

	require.NoError(t, err)
	assert.Len(t, projects, 0)
	mockRepo.AssertExpectations(t)
}

func TestProjectUsecase_GetAllWithStats_RepositoryError(t *testing.T) {
	uc, mockRepo := setupProjectUsecase(t)
	ctx := context.Background()

	mockRepo.On("FindAllWithStats", ctx).Return(nil, errors.New("database error"))

	projects, err := uc.GetAllWithStats(ctx)

	assert.Nil(t, projects)
	assert.Error(t, err)
	mockRepo.AssertExpectations(t)
}

// GetByID tests
func TestProjectUsecase_GetByID_Success(t *testing.T) {
	uc, mockRepo := setupProjectUsecase(t)
	ctx := context.Background()

	expectedProject := createTestProject(1, "PROJ1", 1)
	mockRepo.On("FindByID", ctx, int64(1)).Return(expectedProject, nil)

	project, err := uc.GetByID(ctx, 1)

	require.NoError(t, err)
	assert.Equal(t, int64(1), project.ID)
	assert.Equal(t, "PROJ1", project.Key)
	mockRepo.AssertExpectations(t)
}

func TestProjectUsecase_GetByID_NotFound(t *testing.T) {
	uc, mockRepo := setupProjectUsecase(t)
	ctx := context.Background()

	mockRepo.On("FindByID", ctx, int64(999)).Return(nil, nil)

	project, err := uc.GetByID(ctx, 999)

	require.NoError(t, err)
	assert.Nil(t, project)
	mockRepo.AssertExpectations(t)
}

func TestProjectUsecase_GetByID_RepositoryError(t *testing.T) {
	uc, mockRepo := setupProjectUsecase(t)
	ctx := context.Background()

	mockRepo.On("FindByID", ctx, int64(1)).Return(nil, errors.New("database error"))

	project, err := uc.GetByID(ctx, 1)

	assert.Nil(t, project)
	assert.Error(t, err)
	mockRepo.AssertExpectations(t)
}

// GetWithStats tests
func TestProjectUsecase_GetWithStats_Success(t *testing.T) {
	uc, mockRepo := setupProjectUsecase(t)
	ctx := context.Background()

	expectedStats := createTestProjectWithStats(1, "PROJ1", 5, 3, 12)
	mockRepo.On("FindWithStats", ctx, int64(1)).Return(&expectedStats, nil)

	stats, err := uc.GetWithStats(ctx, 1)

	require.NoError(t, err)
	assert.Equal(t, int64(1), stats.ID)
	assert.Equal(t, 20, stats.TotalIssues)
	assert.Equal(t, 5, stats.RedIssues)
	assert.Equal(t, 3, stats.YellowIssues)
	assert.Equal(t, 12, stats.GreenIssues)
	mockRepo.AssertExpectations(t)
}

func TestProjectUsecase_GetWithStats_NotFound(t *testing.T) {
	uc, mockRepo := setupProjectUsecase(t)
	ctx := context.Background()

	mockRepo.On("FindWithStats", ctx, int64(999)).Return(nil, nil)

	stats, err := uc.GetWithStats(ctx, 999)

	require.NoError(t, err)
	assert.Nil(t, stats)
	mockRepo.AssertExpectations(t)
}

func TestProjectUsecase_GetWithStats_RepositoryError(t *testing.T) {
	uc, mockRepo := setupProjectUsecase(t)
	ctx := context.Background()

	mockRepo.On("FindWithStats", ctx, int64(1)).Return(nil, errors.New("database error"))

	stats, err := uc.GetWithStats(ctx, 1)

	assert.Nil(t, stats)
	assert.Error(t, err)
	mockRepo.AssertExpectations(t)
}

// GetByOrganization tests
func TestProjectUsecase_GetByOrganization_Success(t *testing.T) {
	uc, mockRepo := setupProjectUsecase(t)
	ctx := context.Background()

	expectedProjects := []domain.Project{
		*createTestProject(1, "PROJ1", 1),
		*createTestProject(2, "PROJ2", 1),
	}

	mockRepo.On("FindByOrganizationID", ctx, int64(1)).Return(expectedProjects, nil)

	projects, err := uc.GetByOrganization(ctx, 1)

	require.NoError(t, err)
	assert.Len(t, projects, 2)
	for _, p := range projects {
		assert.Equal(t, int64(1), *p.OrganizationID)
	}
	mockRepo.AssertExpectations(t)
}

func TestProjectUsecase_GetByOrganization_NoProjects(t *testing.T) {
	uc, mockRepo := setupProjectUsecase(t)
	ctx := context.Background()

	mockRepo.On("FindByOrganizationID", ctx, int64(999)).Return([]domain.Project{}, nil)

	projects, err := uc.GetByOrganization(ctx, 999)

	require.NoError(t, err)
	assert.Len(t, projects, 0)
	mockRepo.AssertExpectations(t)
}

func TestProjectUsecase_GetByOrganization_RepositoryError(t *testing.T) {
	uc, mockRepo := setupProjectUsecase(t)
	ctx := context.Background()

	mockRepo.On("FindByOrganizationID", ctx, int64(1)).Return(nil, errors.New("database error"))

	projects, err := uc.GetByOrganization(ctx, 1)

	assert.Nil(t, projects)
	assert.Error(t, err)
	mockRepo.AssertExpectations(t)
}

// GetUnassigned tests
func TestProjectUsecase_GetUnassigned_Success(t *testing.T) {
	uc, mockRepo := setupProjectUsecase(t)
	ctx := context.Background()

	unassignedProject1 := *createTestProject(1, "PROJ1", 1)
	unassignedProject1.OrganizationID = nil
	unassignedProject2 := *createTestProject(2, "PROJ2", 1)
	unassignedProject2.OrganizationID = nil

	expectedProjects := []domain.Project{
		unassignedProject1,
		unassignedProject2,
	}

	mockRepo.On("FindUnassigned", ctx).Return(expectedProjects, nil)

	projects, err := uc.GetUnassigned(ctx)

	require.NoError(t, err)
	assert.Len(t, projects, 2)
	for _, p := range projects {
		assert.Nil(t, p.OrganizationID)
	}
	mockRepo.AssertExpectations(t)
}

func TestProjectUsecase_GetUnassigned_EmptyResult(t *testing.T) {
	uc, mockRepo := setupProjectUsecase(t)
	ctx := context.Background()

	mockRepo.On("FindUnassigned", ctx).Return([]domain.Project{}, nil)

	projects, err := uc.GetUnassigned(ctx)

	require.NoError(t, err)
	assert.Len(t, projects, 0)
	mockRepo.AssertExpectations(t)
}

func TestProjectUsecase_GetUnassigned_RepositoryError(t *testing.T) {
	uc, mockRepo := setupProjectUsecase(t)
	ctx := context.Background()

	mockRepo.On("FindUnassigned", ctx).Return(nil, errors.New("database error"))

	projects, err := uc.GetUnassigned(ctx)

	assert.Nil(t, projects)
	assert.Error(t, err)
	mockRepo.AssertExpectations(t)
}

// AssignToOrganization tests
func TestProjectUsecase_AssignToOrganization_Success(t *testing.T) {
	uc, mockRepo := setupProjectUsecase(t)
	ctx := context.Background()

	orgID := int64(1)
	mockRepo.On("AssignToOrganization", ctx, int64(1), &orgID).Return(nil)

	err := uc.AssignToOrganization(ctx, 1, &orgID)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestProjectUsecase_AssignToOrganization_Unassign(t *testing.T) {
	uc, mockRepo := setupProjectUsecase(t)
	ctx := context.Background()

	mockRepo.On("AssignToOrganization", ctx, int64(1), (*int64)(nil)).Return(nil)

	err := uc.AssignToOrganization(ctx, 1, nil)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestProjectUsecase_AssignToOrganization_RepositoryError(t *testing.T) {
	uc, mockRepo := setupProjectUsecase(t)
	ctx := context.Background()

	orgID := int64(1)
	mockRepo.On("AssignToOrganization", ctx, int64(1), &orgID).Return(errors.New("database error"))

	err := uc.AssignToOrganization(ctx, 1, &orgID)

	assert.Error(t, err)
	mockRepo.AssertExpectations(t)
}

func TestProjectUsecase_AssignToOrganization_ProjectNotFound(t *testing.T) {
	uc, mockRepo := setupProjectUsecase(t)
	ctx := context.Background()

	orgID := int64(1)
	mockRepo.On("AssignToOrganization", ctx, int64(999), &orgID).Return(errors.New("project not found"))

	err := uc.AssignToOrganization(ctx, 999, &orgID)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
	mockRepo.AssertExpectations(t)
}

// Edge case tests
func TestProjectUsecase_GetAllWithStats_MixedProjects(t *testing.T) {
	uc, mockRepo := setupProjectUsecase(t)
	ctx := context.Background()

	expectedProjects := []domain.ProjectWithStats{
		createTestProjectWithStats(1, "PROJ1", 0, 0, 0),    // No issues
		createTestProjectWithStats(2, "PROJ2", 100, 0, 0),  // All red
		createTestProjectWithStats(3, "PROJ3", 0, 100, 0),  // All yellow
		createTestProjectWithStats(4, "PROJ4", 0, 0, 100),  // All green
		createTestProjectWithStats(5, "PROJ5", 10, 20, 30), // Mixed
	}

	mockRepo.On("FindAllWithStats", ctx).Return(expectedProjects, nil)

	projects, err := uc.GetAllWithStats(ctx)

	require.NoError(t, err)
	assert.Len(t, projects, 5)
	assert.Equal(t, 0, projects[0].TotalIssues)
	assert.Equal(t, 100, projects[1].TotalIssues)
	assert.Equal(t, 100, projects[2].TotalIssues)
	assert.Equal(t, 100, projects[3].TotalIssues)
	assert.Equal(t, 60, projects[4].TotalIssues)
	mockRepo.AssertExpectations(t)
}

func TestProjectUsecase_AssignToOrganization_DifferentOrganizations(t *testing.T) {
	uc, mockRepo := setupProjectUsecase(t)
	ctx := context.Background()

	// Test reassigning from one org to another
	oldOrgID := int64(1)
	newOrgID := int64(2)

	mockRepo.On("AssignToOrganization", ctx, int64(1), &oldOrgID).Return(nil)
	mockRepo.On("AssignToOrganization", ctx, int64(1), &newOrgID).Return(nil)

	// First assignment
	err := uc.AssignToOrganization(ctx, 1, &oldOrgID)
	assert.NoError(t, err)

	// Reassignment
	err = uc.AssignToOrganization(ctx, 1, &newOrgID)
	assert.NoError(t, err)

	mockRepo.AssertExpectations(t)
}

// Additional methods for MockProjectRepositoryForUsecase
func (m *MockProjectRepositoryForUsecase) ExistsByJiraProjectID(ctx context.Context, jiraProjectID string) (bool, error) {
	args := m.Called(ctx, jiraProjectID)
	return args.Bool(0), args.Error(1)
}

func (m *MockProjectRepositoryForUsecase) GetByKey(ctx context.Context, key string) (*domain.Project, error) {
	args := m.Called(ctx, key)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Project), args.Error(1)
}

func (m *MockProjectRepositoryForUsecase) FindByJiraProjectID(ctx context.Context, jiraProjectID string) (*domain.Project, error) {
	args := m.Called(ctx, jiraProjectID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Project), args.Error(1)
}
