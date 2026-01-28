package usecase

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/m19cmjigen/sandbox-project-management/backend/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// Mock implementations for dashboard tests
type MockOrganizationRepository struct {
	mock.Mock
}

func (m *MockOrganizationRepository) Create(ctx context.Context, org *domain.Organization) error {
	args := m.Called(ctx, org)
	return args.Error(0)
}

func (m *MockOrganizationRepository) FindByID(ctx context.Context, id int64) (*domain.Organization, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Organization), args.Error(1)
}

func (m *MockOrganizationRepository) FindAll(ctx context.Context) ([]domain.Organization, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.Organization), args.Error(1)
}

func (m *MockOrganizationRepository) Update(ctx context.Context, org *domain.Organization) error {
	args := m.Called(ctx, org)
	return args.Error(0)
}

func (m *MockOrganizationRepository) Delete(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockOrganizationRepository) FindByParentID(ctx context.Context, parentID int64) ([]domain.Organization, error) {
	args := m.Called(ctx, parentID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.Organization), args.Error(1)
}

type MockProjectRepository struct {
	mock.Mock
}

func (m *MockProjectRepository) Create(ctx context.Context, project *domain.Project) error {
	args := m.Called(ctx, project)
	return args.Error(0)
}

func (m *MockProjectRepository) FindByID(ctx context.Context, id int64) (*domain.Project, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Project), args.Error(1)
}

func (m *MockProjectRepository) FindAll(ctx context.Context) ([]domain.Project, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.Project), args.Error(1)
}

func (m *MockProjectRepository) FindAllWithStats(ctx context.Context) ([]domain.ProjectWithStats, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.ProjectWithStats), args.Error(1)
}

func (m *MockProjectRepository) FindWithStats(ctx context.Context, id int64) (*domain.ProjectWithStats, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.ProjectWithStats), args.Error(1)
}

func (m *MockProjectRepository) FindByOrganizationID(ctx context.Context, orgID int64) ([]domain.Project, error) {
	args := m.Called(ctx, orgID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.Project), args.Error(1)
}

func (m *MockProjectRepository) Update(ctx context.Context, project *domain.Project) error {
	args := m.Called(ctx, project)
	return args.Error(0)
}

func (m *MockProjectRepository) Delete(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

type MockIssueRepository struct {
	mock.Mock
}

func (m *MockIssueRepository) Create(ctx context.Context, issue *domain.Issue) error {
	args := m.Called(ctx, issue)
	return args.Error(0)
}

func (m *MockIssueRepository) FindByID(ctx context.Context, id int64) (*domain.Issue, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Issue), args.Error(1)
}

func (m *MockIssueRepository) FindAll(ctx context.Context) ([]domain.Issue, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.Issue), args.Error(1)
}

func (m *MockIssueRepository) FindByFilter(ctx context.Context, filter domain.IssueFilter) ([]domain.Issue, error) {
	args := m.Called(ctx, filter)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.Issue), args.Error(1)
}

func (m *MockIssueRepository) FindByProjectID(ctx context.Context, projectID int64) ([]domain.Issue, error) {
	args := m.Called(ctx, projectID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.Issue), args.Error(1)
}

func (m *MockIssueRepository) Update(ctx context.Context, issue *domain.Issue) error {
	args := m.Called(ctx, issue)
	return args.Error(0)
}

func (m *MockIssueRepository) Delete(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockIssueRepository) BulkUpdate(ctx context.Context, issues []domain.Issue) error {
	args := m.Called(ctx, issues)
	return args.Error(0)
}


// Test helper functions
func createTestOrganization(id int64, name string) *domain.Organization {
	return &domain.Organization{
		ID:        id,
		Name:      name,
		ParentID:  nil,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

func createTestProject(id int64, name string, orgID int64) *domain.Project {
	return &domain.Project{
		ID:             id,
		JiraProjectID:  "10000",
		Key:            name,
		Name:           name + " Project",
		OrganizationID: &orgID,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}
}

func createTestProjectWithStats(id int64, name string, red, yellow, green int) domain.ProjectWithStats {
	orgID := int64(1)
	return domain.ProjectWithStats{
		Project: domain.Project{
			ID:             id,
			JiraProjectID:  "10000",
			Key:            name,
			Name:           name + " Project",
			OrganizationID: &orgID,
			CreatedAt:      time.Now(),
			UpdatedAt:      time.Now(),
		},
		TotalIssues:  red + yellow + green,
		RedIssues:    red,
		YellowIssues: yellow,
		GreenIssues:  green,
	}
}

func setupDashboardUsecase(t *testing.T) (*dashboardUsecase, *MockOrganizationRepository, *MockProjectRepository, *MockIssueRepository) {
	mockOrgRepo := new(MockOrganizationRepository)
	mockProjectRepo := new(MockProjectRepository)
	mockIssueRepo := new(MockIssueRepository)

	usecase := &dashboardUsecase{
		orgRepo:     mockOrgRepo,
		projectRepo: mockProjectRepo,
		issueRepo:   mockIssueRepo,
	}

	return usecase, mockOrgRepo, mockProjectRepo, mockIssueRepo
}

// GetSummary tests
func TestDashboardUsecase_GetSummary_Success(t *testing.T) {
	uc, _, mockProjectRepo, _ := setupDashboardUsecase(t)
	ctx := context.Background()

	projects := []domain.ProjectWithStats{
		createTestProjectWithStats(1, "PROJ1", 5, 3, 2),  // Red project (has red issues)
		createTestProjectWithStats(2, "PROJ2", 0, 4, 6),  // Yellow project (no red, has yellow)
		createTestProjectWithStats(3, "PROJ3", 0, 0, 10), // Green project (no red/yellow)
	}

	mockProjectRepo.On("FindAllWithStats", ctx).Return(projects, nil)

	summary, err := uc.GetSummary(ctx)

	require.NoError(t, err)
	assert.Equal(t, 3, summary.TotalProjects)
	assert.Equal(t, 1, summary.DelayedProjects)  // 1 with red issues
	assert.Equal(t, 1, summary.WarningProjects)  // 1 with yellow issues (no red)
	assert.Equal(t, 1, summary.NormalProjects)   // 1 with only green
	assert.Equal(t, 30, summary.TotalIssues)     // 10 + 10 + 10
	assert.Equal(t, 5, summary.RedIssues)        // 5 + 0 + 0
	assert.Equal(t, 7, summary.YellowIssues)     // 3 + 4 + 0
	assert.Equal(t, 18, summary.GreenIssues)     // 2 + 6 + 10
	assert.Len(t, summary.ProjectsByStatus, 3)
	mockProjectRepo.AssertExpectations(t)
}

func TestDashboardUsecase_GetSummary_EmptyProjects(t *testing.T) {
	uc, _, mockProjectRepo, _ := setupDashboardUsecase(t)
	ctx := context.Background()

	mockProjectRepo.On("FindAllWithStats", ctx).Return([]domain.ProjectWithStats{}, nil)

	summary, err := uc.GetSummary(ctx)

	require.NoError(t, err)
	assert.Equal(t, 0, summary.TotalProjects)
	assert.Equal(t, 0, summary.DelayedProjects)
	assert.Equal(t, 0, summary.WarningProjects)
	assert.Equal(t, 0, summary.NormalProjects)
	assert.Equal(t, 0, summary.TotalIssues)
	mockProjectRepo.AssertExpectations(t)
}

func TestDashboardUsecase_GetSummary_RepositoryError(t *testing.T) {
	uc, _, mockProjectRepo, _ := setupDashboardUsecase(t)
	ctx := context.Background()

	mockProjectRepo.On("FindAllWithStats", ctx).Return(nil, errors.New("database error"))

	summary, err := uc.GetSummary(ctx)

	assert.Nil(t, summary)
	assert.Error(t, err)
	mockProjectRepo.AssertExpectations(t)
}

// GetOrganizationSummary tests
func TestDashboardUsecase_GetOrganizationSummary_Success(t *testing.T) {
	uc, mockOrgRepo, mockProjectRepo, _ := setupDashboardUsecase(t)
	ctx := context.Background()

	org := *createTestOrganization(1, "Test Organization")
	projects := []domain.Project{
		*createTestProject(1, "PROJ1", 1),
		*createTestProject(2, "PROJ2", 1),
	}

	projectStats1 := createTestProjectWithStats(1, "PROJ1", 3, 2, 5)
	projectStats2 := createTestProjectWithStats(2, "PROJ2", 0, 4, 6)

	mockOrgRepo.On("FindByID", ctx, int64(1)).Return(&org, nil)
	mockProjectRepo.On("FindByOrganizationID", ctx, int64(1)).Return(projects, nil)
	mockProjectRepo.On("FindWithStats", ctx, int64(1)).Return(&projectStats1, nil)
	mockProjectRepo.On("FindWithStats", ctx, int64(2)).Return(&projectStats2, nil)

	summary, err := uc.GetOrganizationSummary(ctx, 1)

	require.NoError(t, err)
	assert.Equal(t, org.ID, summary.Organization.ID)
	assert.Equal(t, 2, summary.TotalProjects)
	assert.Equal(t, 1, summary.DelayedProjects) // PROJ1 has red issues
	assert.Equal(t, 1, summary.WarningProjects) // PROJ2 has yellow issues
	assert.Len(t, summary.Projects, 2)
	mockOrgRepo.AssertExpectations(t)
	mockProjectRepo.AssertExpectations(t)
}

func TestDashboardUsecase_GetOrganizationSummary_OrganizationNotFound(t *testing.T) {
	uc, mockOrgRepo, _, _ := setupDashboardUsecase(t)
	ctx := context.Background()

	mockOrgRepo.On("FindByID", ctx, int64(999)).Return(nil, nil)

	summary, err := uc.GetOrganizationSummary(ctx, 999)

	require.NoError(t, err)
	assert.Nil(t, summary)
	mockOrgRepo.AssertExpectations(t)
}

func TestDashboardUsecase_GetOrganizationSummary_NoProjects(t *testing.T) {
	uc, mockOrgRepo, mockProjectRepo, _ := setupDashboardUsecase(t)
	ctx := context.Background()

	org := *createTestOrganization(1, "Test Organization")
	mockOrgRepo.On("FindByID", ctx, int64(1)).Return(&org, nil)
	mockProjectRepo.On("FindByOrganizationID", ctx, int64(1)).Return([]domain.Project{}, nil)

	summary, err := uc.GetOrganizationSummary(ctx, 1)

	require.NoError(t, err)
	assert.Equal(t, 0, summary.TotalProjects)
	assert.Equal(t, 0, summary.DelayedProjects)
	assert.Equal(t, 0, summary.WarningProjects)
	assert.Len(t, summary.Projects, 0)
	mockOrgRepo.AssertExpectations(t)
	mockProjectRepo.AssertExpectations(t)
}

func TestDashboardUsecase_GetOrganizationSummary_RepositoryError(t *testing.T) {
	uc, mockOrgRepo, _, _ := setupDashboardUsecase(t)
	ctx := context.Background()

	mockOrgRepo.On("FindByID", ctx, int64(1)).Return(nil, errors.New("database error"))

	summary, err := uc.GetOrganizationSummary(ctx, 1)

	assert.Nil(t, summary)
	assert.Error(t, err)
	mockOrgRepo.AssertExpectations(t)
}

// GetProjectSummary tests
func TestDashboardUsecase_GetProjectSummary_Success(t *testing.T) {
	uc, _, mockProjectRepo, _ := setupDashboardUsecase(t)
	ctx := context.Background()

	expectedStats := createTestProjectWithStats(1, "PROJ1", 5, 3, 12)
	mockProjectRepo.On("FindWithStats", ctx, int64(1)).Return(&expectedStats, nil)

	stats, err := uc.GetProjectSummary(ctx, 1)

	require.NoError(t, err)
	assert.Equal(t, int64(1), stats.ID)
	assert.Equal(t, 20, stats.TotalIssues)
	assert.Equal(t, 5, stats.RedIssues)
	assert.Equal(t, 3, stats.YellowIssues)
	assert.Equal(t, 12, stats.GreenIssues)
	mockProjectRepo.AssertExpectations(t)
}

func TestDashboardUsecase_GetProjectSummary_NotFound(t *testing.T) {
	uc, _, mockProjectRepo, _ := setupDashboardUsecase(t)
	ctx := context.Background()

	mockProjectRepo.On("FindWithStats", ctx, int64(999)).Return(nil, nil)

	stats, err := uc.GetProjectSummary(ctx, 999)

	require.NoError(t, err)
	assert.Nil(t, stats)
	mockProjectRepo.AssertExpectations(t)
}

func TestDashboardUsecase_GetProjectSummary_RepositoryError(t *testing.T) {
	uc, _, mockProjectRepo, _ := setupDashboardUsecase(t)
	ctx := context.Background()

	mockProjectRepo.On("FindWithStats", ctx, int64(1)).Return(nil, errors.New("database error"))

	stats, err := uc.GetProjectSummary(ctx, 1)

	assert.Nil(t, stats)
	assert.Error(t, err)
	mockProjectRepo.AssertExpectations(t)
}

// Edge case tests
func TestDashboardUsecase_GetSummary_AllRedProjects(t *testing.T) {
	uc, _, mockProjectRepo, _ := setupDashboardUsecase(t)
	ctx := context.Background()

	projects := []domain.ProjectWithStats{
		createTestProjectWithStats(1, "PROJ1", 5, 0, 0),
		createTestProjectWithStats(2, "PROJ2", 3, 2, 1),
		createTestProjectWithStats(3, "PROJ3", 10, 5, 0),
	}

	mockProjectRepo.On("FindAllWithStats", ctx).Return(projects, nil)

	summary, err := uc.GetSummary(ctx)

	require.NoError(t, err)
	assert.Equal(t, 3, summary.TotalProjects)
	assert.Equal(t, 3, summary.DelayedProjects)  // All have red issues
	assert.Equal(t, 0, summary.WarningProjects)
	assert.Equal(t, 0, summary.NormalProjects)
	mockProjectRepo.AssertExpectations(t)
}

func TestDashboardUsecase_GetSummary_MixedStatusProjects(t *testing.T) {
	uc, _, mockProjectRepo, _ := setupDashboardUsecase(t)
	ctx := context.Background()

	projects := []domain.ProjectWithStats{
		createTestProjectWithStats(1, "PROJ1", 0, 0, 100), // All green
		createTestProjectWithStats(2, "PROJ2", 0, 0, 0),   // No issues
		createTestProjectWithStats(3, "PROJ3", 1, 1, 1),   // Has red
	}

	mockProjectRepo.On("FindAllWithStats", ctx).Return(projects, nil)

	summary, err := uc.GetSummary(ctx)

	require.NoError(t, err)
	assert.Equal(t, 3, summary.TotalProjects)
	assert.Equal(t, 1, summary.DelayedProjects)  // PROJ3 with red
	assert.Equal(t, 0, summary.WarningProjects)  // PROJ1 and PROJ2 have no yellow (or are all green/none)
	assert.Equal(t, 2, summary.NormalProjects)   // PROJ1 and PROJ2
	mockProjectRepo.AssertExpectations(t)
}

// Additional methods for MockOrganizationRepository  
func (m *MockOrganizationRepository) ExistsByID(ctx context.Context, id int64) (bool, error) {
	args := m.Called(ctx, id)
	return args.Bool(0), args.Error(1)
}

func (m *MockOrganizationRepository) FindByPath(ctx context.Context, pathPrefix string) ([]domain.Organization, error) {
	args := m.Called(ctx, pathPrefix)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.Organization), args.Error(1)
}

func (m *MockOrganizationRepository) FindRoots(ctx context.Context) ([]domain.Organization, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.Organization), args.Error(1)
}

func (m *MockOrganizationRepository) HasChildren(ctx context.Context, id int64) (bool, error) {
	args := m.Called(ctx, id)
	return args.Bool(0), args.Error(1)
}

// Additional methods for MockProjectRepository
func (m *MockProjectRepository) AssignToOrganization(ctx context.Context, projectID int64, organizationID *int64) error {
	args := m.Called(ctx, projectID, organizationID)
	return args.Error(0)
}

func (m *MockProjectRepository) ExistsByJiraProjectID(ctx context.Context, jiraProjectID string) (bool, error) {
	args := m.Called(ctx, jiraProjectID)
	return args.Bool(0), args.Error(1)
}

func (m *MockProjectRepository) GetByKey(ctx context.Context, key string) (*domain.Project, error) {
	args := m.Called(ctx, key)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Project), args.Error(1)
}

func (m *MockProjectRepository) FindByJiraProjectID(ctx context.Context, jiraProjectID string) (*domain.Project, error) {
	args := m.Called(ctx, jiraProjectID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Project), args.Error(1)
}

func (m *MockProjectRepository) FindUnassigned(ctx context.Context) ([]domain.Project, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.Project), args.Error(1)
}

// Additional methods for MockIssueRepository
func (m *MockIssueRepository) BulkUpsert(ctx context.Context, issues []domain.Issue) error {
	args := m.Called(ctx, issues)
	return args.Error(0)
}

func (m *MockIssueRepository) ExistsByJiraIssueID(ctx context.Context, jiraIssueID string) (bool, error) {
	args := m.Called(ctx, jiraIssueID)
	return args.Bool(0), args.Error(1)
}

func (m *MockIssueRepository) FindByJiraIssueID(ctx context.Context, jiraIssueID string) (*domain.Issue, error) {
	args := m.Called(ctx, jiraIssueID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Issue), args.Error(1)
}

func (m *MockIssueRepository) GetByJiraKey(ctx context.Context, jiraKey string) (*domain.Issue, error) {
	args := m.Called(ctx, jiraKey)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Issue), args.Error(1)
}

func (m *MockIssueRepository) CountByProjectID(ctx context.Context, projectID int64) (int, error) {
	args := m.Called(ctx, projectID)
	return args.Int(0), args.Error(1)
}

func (m *MockIssueRepository) CountByDelayStatus(ctx context.Context, projectID int64, status domain.DelayStatus) (int, error) {
	args := m.Called(ctx, projectID, status)
	return args.Int(0), args.Error(1)
}

func (m *MockIssueRepository) FindByDelayStatus(ctx context.Context, status domain.DelayStatus) ([]domain.Issue, error) {
	args := m.Called(ctx, status)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.Issue), args.Error(1)
}
