package usecase

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"

	"github.com/m19cmjigen/sandbox-project-management/backend/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// MockIssueRepositoryForUsecase extends MockIssueRepository with FindByDelayStatus
type MockIssueRepositoryForUsecase struct {
	MockIssueRepository
}

func (m *MockIssueRepositoryForUsecase) FindByDelayStatus(ctx context.Context, status domain.DelayStatus) ([]domain.Issue, error) {
	args := m.Called(ctx, status)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.Issue), args.Error(1)
}

// Test helper functions for issues
func createTestIssue(id int64, projectID int64, delayStatus domain.DelayStatus) *domain.Issue {
	summary := "Test Issue"
	priority := "High"
	issueType := "Task"
	assignee := "test@example.com"

	return &domain.Issue{
		ID:                id,
		JiraIssueID:       "1000",
		JiraIssueKey:      "TEST-1",
		ProjectID:         projectID,
		Summary:           summary,
		Status:            "In Progress",
		StatusCategory:    domain.StatusCategoryInProgress,
		DueDate:           sql.NullTime{Time: time.Now().Add(7 * 24 * time.Hour), Valid: true},
		DelayStatus:       delayStatus,
		Priority:          &priority,
		IssueType:         &issueType,
		AssigneeName:      &assignee,
		AssigneeAccountID: &assignee,
		CreatedAt:         time.Now(),
		UpdatedAt:         time.Now(),
	}
}

func setupIssueUsecase(t *testing.T) (*issueUsecase, *MockIssueRepositoryForUsecase) {
	mockRepo := new(MockIssueRepositoryForUsecase)
	usecase := &issueUsecase{
		issueRepo: mockRepo,
	}
	return usecase, mockRepo
}

// GetAll tests
func TestIssueUsecase_GetAll_Success(t *testing.T) {
	uc, mockRepo := setupIssueUsecase(t)
	ctx := context.Background()

	expectedIssues := []domain.Issue{
		*createTestIssue(1, 1, domain.DelayStatusGreen),
		*createTestIssue(2, 1, domain.DelayStatusYellow),
		*createTestIssue(3, 2, domain.DelayStatusRed),
	}

	mockRepo.On("FindAll", ctx).Return(expectedIssues, nil)

	issues, err := uc.GetAll(ctx)

	require.NoError(t, err)
	assert.Len(t, issues, 3)
	mockRepo.AssertExpectations(t)
}

func TestIssueUsecase_GetAll_EmptyResult(t *testing.T) {
	uc, mockRepo := setupIssueUsecase(t)
	ctx := context.Background()

	mockRepo.On("FindAll", ctx).Return([]domain.Issue{}, nil)

	issues, err := uc.GetAll(ctx)

	require.NoError(t, err)
	assert.Len(t, issues, 0)
	mockRepo.AssertExpectations(t)
}

func TestIssueUsecase_GetAll_RepositoryError(t *testing.T) {
	uc, mockRepo := setupIssueUsecase(t)
	ctx := context.Background()

	mockRepo.On("FindAll", ctx).Return(nil, errors.New("database error"))

	issues, err := uc.GetAll(ctx)

	assert.Nil(t, issues)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "database error")
	mockRepo.AssertExpectations(t)
}

// GetByID tests
func TestIssueUsecase_GetByID_Success(t *testing.T) {
	uc, mockRepo := setupIssueUsecase(t)
	ctx := context.Background()

	expectedIssue := createTestIssue(1, 1, domain.DelayStatusGreen)
	mockRepo.On("FindByID", ctx, int64(1)).Return(expectedIssue, nil)

	issue, err := uc.GetByID(ctx, 1)

	require.NoError(t, err)
	assert.Equal(t, int64(1), issue.ID)
	assert.Equal(t, "Test Issue", issue.Summary)
	mockRepo.AssertExpectations(t)
}

func TestIssueUsecase_GetByID_NotFound(t *testing.T) {
	uc, mockRepo := setupIssueUsecase(t)
	ctx := context.Background()

	mockRepo.On("FindByID", ctx, int64(999)).Return(nil, nil)

	issue, err := uc.GetByID(ctx, 999)

	require.NoError(t, err)
	assert.Nil(t, issue)
	mockRepo.AssertExpectations(t)
}

func TestIssueUsecase_GetByID_RepositoryError(t *testing.T) {
	uc, mockRepo := setupIssueUsecase(t)
	ctx := context.Background()

	mockRepo.On("FindByID", ctx, int64(1)).Return(nil, errors.New("database error"))

	issue, err := uc.GetByID(ctx, 1)

	assert.Nil(t, issue)
	assert.Error(t, err)
	mockRepo.AssertExpectations(t)
}

// GetByProjectID tests
func TestIssueUsecase_GetByProjectID_Success(t *testing.T) {
	uc, mockRepo := setupIssueUsecase(t)
	ctx := context.Background()

	expectedIssues := []domain.Issue{
		*createTestIssue(1, 1, domain.DelayStatusGreen),
		*createTestIssue(2, 1, domain.DelayStatusYellow),
		*createTestIssue(3, 1, domain.DelayStatusRed),
	}

	mockRepo.On("FindByProjectID", ctx, int64(1)).Return(expectedIssues, nil)

	issues, err := uc.GetByProjectID(ctx, 1)

	require.NoError(t, err)
	assert.Len(t, issues, 3)
	for _, issue := range issues {
		assert.Equal(t, int64(1), issue.ProjectID)
	}
	mockRepo.AssertExpectations(t)
}

func TestIssueUsecase_GetByProjectID_NoIssues(t *testing.T) {
	uc, mockRepo := setupIssueUsecase(t)
	ctx := context.Background()

	mockRepo.On("FindByProjectID", ctx, int64(999)).Return([]domain.Issue{}, nil)

	issues, err := uc.GetByProjectID(ctx, 999)

	require.NoError(t, err)
	assert.Len(t, issues, 0)
	mockRepo.AssertExpectations(t)
}

func TestIssueUsecase_GetByProjectID_RepositoryError(t *testing.T) {
	uc, mockRepo := setupIssueUsecase(t)
	ctx := context.Background()

	mockRepo.On("FindByProjectID", ctx, int64(1)).Return(nil, errors.New("database error"))

	issues, err := uc.GetByProjectID(ctx, 1)

	assert.Nil(t, issues)
	assert.Error(t, err)
	mockRepo.AssertExpectations(t)
}

// GetByFilter tests
func TestIssueUsecase_GetByFilter_Success(t *testing.T) {
	uc, mockRepo := setupIssueUsecase(t)
	ctx := context.Background()

	delayStatus := domain.DelayStatusRed
	filter := domain.IssueFilter{
		DelayStatus: &delayStatus,
	}

	expectedIssues := []domain.Issue{
		*createTestIssue(1, 1, domain.DelayStatusRed),
		*createTestIssue(2, 2, domain.DelayStatusRed),
	}

	mockRepo.On("FindByFilter", ctx, filter).Return(expectedIssues, nil)

	issues, err := uc.GetByFilter(ctx, filter)

	require.NoError(t, err)
	assert.Len(t, issues, 2)
	for _, issue := range issues {
		assert.Equal(t, domain.DelayStatusRed, issue.DelayStatus)
	}
	mockRepo.AssertExpectations(t)
}

func TestIssueUsecase_GetByFilter_MultipleFilters(t *testing.T) {
	uc, mockRepo := setupIssueUsecase(t)
	ctx := context.Background()

	projectID := int64(1)
	delayStatus := domain.DelayStatusYellow
	filter := domain.IssueFilter{
		ProjectID:   &projectID,
		DelayStatus: &delayStatus,
	}

	expectedIssues := []domain.Issue{
		*createTestIssue(1, 1, domain.DelayStatusYellow),
		*createTestIssue(2, 1, domain.DelayStatusYellow),
	}

	mockRepo.On("FindByFilter", ctx, filter).Return(expectedIssues, nil)

	issues, err := uc.GetByFilter(ctx, filter)

	require.NoError(t, err)
	assert.Len(t, issues, 2)
	mockRepo.AssertExpectations(t)
}

func TestIssueUsecase_GetByFilter_EmptyResult(t *testing.T) {
	uc, mockRepo := setupIssueUsecase(t)
	ctx := context.Background()

	filter := domain.IssueFilter{}
	mockRepo.On("FindByFilter", ctx, filter).Return([]domain.Issue{}, nil)

	issues, err := uc.GetByFilter(ctx, filter)

	require.NoError(t, err)
	assert.Len(t, issues, 0)
	mockRepo.AssertExpectations(t)
}

func TestIssueUsecase_GetByFilter_RepositoryError(t *testing.T) {
	uc, mockRepo := setupIssueUsecase(t)
	ctx := context.Background()

	filter := domain.IssueFilter{}
	mockRepo.On("FindByFilter", ctx, filter).Return(nil, errors.New("database error"))

	issues, err := uc.GetByFilter(ctx, filter)

	assert.Nil(t, issues)
	assert.Error(t, err)
	mockRepo.AssertExpectations(t)
}

// GetByDelayStatus tests
func TestIssueUsecase_GetByDelayStatus_Red(t *testing.T) {
	uc, mockRepo := setupIssueUsecase(t)
	ctx := context.Background()

	expectedIssues := []domain.Issue{
		*createTestIssue(1, 1, domain.DelayStatusRed),
		*createTestIssue(2, 2, domain.DelayStatusRed),
		*createTestIssue(3, 3, domain.DelayStatusRed),
	}

	mockRepo.On("FindByDelayStatus", ctx, domain.DelayStatusRed).Return(expectedIssues, nil)

	issues, err := uc.GetByDelayStatus(ctx, domain.DelayStatusRed)

	require.NoError(t, err)
	assert.Len(t, issues, 3)
	for _, issue := range issues {
		assert.Equal(t, domain.DelayStatusRed, issue.DelayStatus)
	}
	mockRepo.AssertExpectations(t)
}

func TestIssueUsecase_GetByDelayStatus_Yellow(t *testing.T) {
	uc, mockRepo := setupIssueUsecase(t)
	ctx := context.Background()

	expectedIssues := []domain.Issue{
		*createTestIssue(1, 1, domain.DelayStatusYellow),
		*createTestIssue(2, 1, domain.DelayStatusYellow),
	}

	mockRepo.On("FindByDelayStatus", ctx, domain.DelayStatusYellow).Return(expectedIssues, nil)

	issues, err := uc.GetByDelayStatus(ctx, domain.DelayStatusYellow)

	require.NoError(t, err)
	assert.Len(t, issues, 2)
	for _, issue := range issues {
		assert.Equal(t, domain.DelayStatusYellow, issue.DelayStatus)
	}
	mockRepo.AssertExpectations(t)
}

func TestIssueUsecase_GetByDelayStatus_Green(t *testing.T) {
	uc, mockRepo := setupIssueUsecase(t)
	ctx := context.Background()

	expectedIssues := []domain.Issue{
		*createTestIssue(1, 1, domain.DelayStatusGreen),
	}

	mockRepo.On("FindByDelayStatus", ctx, domain.DelayStatusGreen).Return(expectedIssues, nil)

	issues, err := uc.GetByDelayStatus(ctx, domain.DelayStatusGreen)

	require.NoError(t, err)
	assert.Len(t, issues, 1)
	assert.Equal(t, domain.DelayStatusGreen, issues[0].DelayStatus)
	mockRepo.AssertExpectations(t)
}

func TestIssueUsecase_GetByDelayStatus_NoResults(t *testing.T) {
	uc, mockRepo := setupIssueUsecase(t)
	ctx := context.Background()

	mockRepo.On("FindByDelayStatus", ctx, domain.DelayStatusRed).Return([]domain.Issue{}, nil)

	issues, err := uc.GetByDelayStatus(ctx, domain.DelayStatusRed)

	require.NoError(t, err)
	assert.Len(t, issues, 0)
	mockRepo.AssertExpectations(t)
}

func TestIssueUsecase_GetByDelayStatus_RepositoryError(t *testing.T) {
	uc, mockRepo := setupIssueUsecase(t)
	ctx := context.Background()

	mockRepo.On("FindByDelayStatus", ctx, domain.DelayStatusRed).Return(nil, errors.New("database error"))

	issues, err := uc.GetByDelayStatus(ctx, domain.DelayStatusRed)

	assert.Nil(t, issues)
	assert.Error(t, err)
	mockRepo.AssertExpectations(t)
}

// Edge case tests
func TestIssueUsecase_GetByFilter_WithStatusCategory(t *testing.T) {
	uc, mockRepo := setupIssueUsecase(t)
	ctx := context.Background()

	status := "In Progress"
	filter := domain.IssueFilter{
		Status: &status,
	}

	expectedIssues := []domain.Issue{
		*createTestIssue(1, 1, domain.DelayStatusGreen),
		*createTestIssue(2, 2, domain.DelayStatusYellow),
	}

	mockRepo.On("FindByFilter", ctx, filter).Return(expectedIssues, nil)

	issues, err := uc.GetByFilter(ctx, filter)

	require.NoError(t, err)
	assert.Len(t, issues, 2)
	for _, issue := range issues {
		assert.Equal(t, domain.StatusCategoryInProgress, issue.StatusCategory)
	}
	mockRepo.AssertExpectations(t)
}

func TestIssueUsecase_GetByFilter_ComplexFilter(t *testing.T) {
	uc, mockRepo := setupIssueUsecase(t)
	ctx := context.Background()

	projectID := int64(1)
	delayStatus := domain.DelayStatusRed
	status := "In Progress"
	assignee := "test@example.com"

	filter := domain.IssueFilter{
		ProjectID:      &projectID,
		DelayStatus:    &delayStatus,
		Status: &status,
		AssigneeID:  &assignee,
	}

	expectedIssues := []domain.Issue{
		*createTestIssue(1, 1, domain.DelayStatusRed),
	}

	mockRepo.On("FindByFilter", ctx, filter).Return(expectedIssues, nil)

	issues, err := uc.GetByFilter(ctx, filter)

	require.NoError(t, err)
	assert.Len(t, issues, 1)
	assert.Equal(t, int64(1), issues[0].ProjectID)
	assert.Equal(t, domain.DelayStatusRed, issues[0].DelayStatus)
	mockRepo.AssertExpectations(t)
}

// Additional methods for MockIssueRepositoryForUsecase
func (m *MockIssueRepositoryForUsecase) BulkUpsert(ctx context.Context, issues []domain.Issue) error {
	args := m.Called(ctx, issues)
	return args.Error(0)
}

func (m *MockIssueRepositoryForUsecase) ExistsByJiraIssueID(ctx context.Context, jiraIssueID string) (bool, error) {
	args := m.Called(ctx, jiraIssueID)
	return args.Bool(0), args.Error(1)
}

func (m *MockIssueRepositoryForUsecase) FindByJiraIssueID(ctx context.Context, jiraIssueID string) (*domain.Issue, error) {
	args := m.Called(ctx, jiraIssueID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Issue), args.Error(1)
}

func (m *MockIssueRepositoryForUsecase) GetByJiraKey(ctx context.Context, jiraKey string) (*domain.Issue, error) {
	args := m.Called(ctx, jiraKey)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Issue), args.Error(1)
}

func (m *MockIssueRepositoryForUsecase) CountByProjectID(ctx context.Context, projectID int64) (int, error) {
	args := m.Called(ctx, projectID)
	return args.Int(0), args.Error(1)
}

func (m *MockIssueRepositoryForUsecase) CountByDelayStatus(ctx context.Context, projectID int64, status domain.DelayStatus) (int, error) {
	args := m.Called(ctx, projectID, status)
	return args.Int(0), args.Error(1)
}
