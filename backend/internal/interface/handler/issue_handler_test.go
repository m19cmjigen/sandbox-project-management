package handler

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/m19cmjigen/sandbox-project-management/backend/internal/domain"
	"github.com/m19cmjigen/sandbox-project-management/backend/pkg/logger"
)

// MockIssueUsecase はIssueUsecaseのモック
type MockIssueUsecase struct {
	mock.Mock
}

func (m *MockIssueUsecase) GetAll(ctx context.Context) ([]domain.Issue, error) {
	args := m.Called(ctx)
	return args.Get(0).([]domain.Issue), args.Error(1)
}

func (m *MockIssueUsecase) GetByID(ctx context.Context, id int64) (*domain.Issue, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Issue), args.Error(1)
}

func (m *MockIssueUsecase) GetByProjectID(ctx context.Context, projectID int64) ([]domain.Issue, error) {
	args := m.Called(ctx, projectID)
	return args.Get(0).([]domain.Issue), args.Error(1)
}

func (m *MockIssueUsecase) GetByFilter(ctx context.Context, filter domain.IssueFilter) ([]domain.Issue, error) {
	args := m.Called(ctx, filter)
	return args.Get(0).([]domain.Issue), args.Error(1)
}

func (m *MockIssueUsecase) GetByDelayStatus(ctx context.Context, delayStatus domain.DelayStatus) ([]domain.Issue, error) {
	args := m.Called(ctx, delayStatus)
	return args.Get(0).([]domain.Issue), args.Error(1)
}

func setupIssueHandler() (*IssueHandler, *MockIssueUsecase) {
	gin.SetMode(gin.TestMode)
	mockUsecase := new(MockIssueUsecase)
	log, _ := logger.New("error", "console")
	handler := NewIssueHandler(mockUsecase, log)
	return handler, mockUsecase
}

func createTestIssueForHandler(id int64, projectID int64, summary string, status string) domain.Issue {
	now := time.Now()
	dueDate := now.AddDate(0, 0, 7)
	priority := "Medium"
	return domain.Issue{
		ID:                id,
		JiraIssueID:       "JIRA-" + summary,
		JiraIssueKey:      "TEST-" + summary,
		ProjectID:         projectID,
		Summary:           summary,
		Status:            status,
		StatusCategory:    domain.StatusCategoryToDo,
		DueDate:           sql.NullTime{Time: dueDate, Valid: true},
		AssigneeName:      nil,
		AssigneeAccountID: nil,
		DelayStatus:       domain.DelayStatusGreen,
		Priority:          &priority,
		IssueType:         nil,
		LastUpdatedAt:     now,
		CreatedAt:         now,
		UpdatedAt:         now,
	}
}

// ListIssues Tests

func TestIssueHandler_ListIssues_Success_NoFilters(t *testing.T) {
	handler, mockUsecase := setupIssueHandler()

	issues := []domain.Issue{
		createTestIssueForHandler(1, 1, "Issue 1", "TODO"),
		createTestIssueForHandler(2, 2, "Issue 2", "IN_PROGRESS"),
	}

	expectedFilter := domain.IssueFilter{
		Limit:  100,
		Offset: 0,
	}

	mockUsecase.On("GetByFilter", mock.Anything, expectedFilter).Return(issues, nil)

	w := httptest.NewRecorder()
	_, router := gin.CreateTestContext(w)
	router.GET("/issues", handler.ListIssues)

	req := httptest.NewRequest(http.MethodGet, "/issues", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, float64(2), response["count"])

	mockUsecase.AssertExpectations(t)
}

func TestIssueHandler_ListIssues_Success_WithProjectIDFilter(t *testing.T) {
	handler, mockUsecase := setupIssueHandler()

	projectID := int64(1)
	issues := []domain.Issue{
		createTestIssueForHandler(1, projectID, "Issue 1", "TODO"),
	}

	expectedFilter := domain.IssueFilter{
		ProjectID: &projectID,
		Limit:     100,
		Offset:    0,
	}

	mockUsecase.On("GetByFilter", mock.Anything, expectedFilter).Return(issues, nil)

	w := httptest.NewRecorder()
	_, router := gin.CreateTestContext(w)
	router.GET("/issues", handler.ListIssues)

	req := httptest.NewRequest(http.MethodGet, "/issues?project_id=1", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, float64(1), response["count"])

	mockUsecase.AssertExpectations(t)
}

func TestIssueHandler_ListIssues_Success_WithDelayStatusFilter(t *testing.T) {
	handler, mockUsecase := setupIssueHandler()

	delayStatus := domain.DelayStatusRed
	issues := []domain.Issue{
		createTestIssueForHandler(1, 1, "Delayed Issue", "IN_PROGRESS"),
	}

	expectedFilter := domain.IssueFilter{
		DelayStatus: &delayStatus,
		Limit:       100,
		Offset:      0,
	}

	mockUsecase.On("GetByFilter", mock.Anything, expectedFilter).Return(issues, nil)

	w := httptest.NewRecorder()
	_, router := gin.CreateTestContext(w)
	router.GET("/issues", handler.ListIssues)

	req := httptest.NewRequest(http.MethodGet, "/issues?delay_status=RED", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockUsecase.AssertExpectations(t)
}

func TestIssueHandler_ListIssues_Success_WithStatusFilter(t *testing.T) {
	handler, mockUsecase := setupIssueHandler()

	status := "TODO"
	issues := []domain.Issue{
		createTestIssueForHandler(1, 1, "Todo Issue", status),
	}

	expectedFilter := domain.IssueFilter{
		Status: &status,
		Limit:  100,
		Offset: 0,
	}

	mockUsecase.On("GetByFilter", mock.Anything, expectedFilter).Return(issues, nil)

	w := httptest.NewRecorder()
	_, router := gin.CreateTestContext(w)
	router.GET("/issues", handler.ListIssues)

	req := httptest.NewRequest(http.MethodGet, "/issues?status=TODO", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockUsecase.AssertExpectations(t)
}

func TestIssueHandler_ListIssues_Success_WithAssigneeFilter(t *testing.T) {
	handler, mockUsecase := setupIssueHandler()

	assigneeID := "user123"
	issues := []domain.Issue{
		createTestIssueForHandler(1, 1, "Assigned Issue", "IN_PROGRESS"),
	}

	expectedFilter := domain.IssueFilter{
		AssigneeID: &assigneeID,
		Limit:      100,
		Offset:     0,
	}

	mockUsecase.On("GetByFilter", mock.Anything, expectedFilter).Return(issues, nil)

	w := httptest.NewRecorder()
	_, router := gin.CreateTestContext(w)
	router.GET("/issues", handler.ListIssues)

	req := httptest.NewRequest(http.MethodGet, "/issues?assignee_id=user123", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockUsecase.AssertExpectations(t)
}

func TestIssueHandler_ListIssues_Success_WithDueDateFilters(t *testing.T) {
	handler, mockUsecase := setupIssueHandler()

	dueDateFrom, _ := time.Parse("2006-01-02", "2026-01-01")
	dueDateTo, _ := time.Parse("2006-01-02", "2026-12-31")
	issues := []domain.Issue{
		createTestIssueForHandler(1, 1, "Issue in range", "TODO"),
	}

	expectedFilter := domain.IssueFilter{
		DueDateFrom: &dueDateFrom,
		DueDateTo:   &dueDateTo,
		Limit:       100,
		Offset:      0,
	}

	mockUsecase.On("GetByFilter", mock.Anything, expectedFilter).Return(issues, nil)

	w := httptest.NewRecorder()
	_, router := gin.CreateTestContext(w)
	router.GET("/issues", handler.ListIssues)

	req := httptest.NewRequest(http.MethodGet, "/issues?due_date_from=2026-01-01&due_date_to=2026-12-31", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockUsecase.AssertExpectations(t)
}

func TestIssueHandler_ListIssues_Success_WithPagination(t *testing.T) {
	handler, mockUsecase := setupIssueHandler()

	issues := []domain.Issue{
		createTestIssueForHandler(11, 1, "Issue 11", "TODO"),
	}

	expectedFilter := domain.IssueFilter{
		Limit:  10,
		Offset: 10,
	}

	mockUsecase.On("GetByFilter", mock.Anything, expectedFilter).Return(issues, nil)

	w := httptest.NewRecorder()
	_, router := gin.CreateTestContext(w)
	router.GET("/issues", handler.ListIssues)

	req := httptest.NewRequest(http.MethodGet, "/issues?limit=10&offset=10", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockUsecase.AssertExpectations(t)
}

func TestIssueHandler_ListIssues_InvalidProjectID(t *testing.T) {
	handler, _ := setupIssueHandler()

	w := httptest.NewRecorder()
	_, router := gin.CreateTestContext(w)
	router.GET("/issues", handler.ListIssues)

	req := httptest.NewRequest(http.MethodGet, "/issues?project_id=invalid", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]string
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, "Invalid project_id", response["error"])
}

func TestIssueHandler_ListIssues_InvalidDueDateFromFormat(t *testing.T) {
	handler, _ := setupIssueHandler()

	w := httptest.NewRecorder()
	_, router := gin.CreateTestContext(w)
	router.GET("/issues", handler.ListIssues)

	req := httptest.NewRequest(http.MethodGet, "/issues?due_date_from=invalid-date", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]string
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Contains(t, response["error"], "Invalid due_date_from format")
}

func TestIssueHandler_ListIssues_InvalidDueDateToFormat(t *testing.T) {
	handler, _ := setupIssueHandler()

	w := httptest.NewRecorder()
	_, router := gin.CreateTestContext(w)
	router.GET("/issues", handler.ListIssues)

	req := httptest.NewRequest(http.MethodGet, "/issues?due_date_to=2026/01/01", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]string
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Contains(t, response["error"], "Invalid due_date_to format")
}

func TestIssueHandler_ListIssues_UsecaseError(t *testing.T) {
	handler, mockUsecase := setupIssueHandler()

	mockUsecase.On("GetByFilter", mock.Anything, mock.AnythingOfType("domain.IssueFilter")).
		Return([]domain.Issue{}, errors.New("database error"))

	w := httptest.NewRecorder()
	_, router := gin.CreateTestContext(w)
	router.GET("/issues", handler.ListIssues)

	req := httptest.NewRequest(http.MethodGet, "/issues", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)

	var response map[string]string
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, "Failed to get issues", response["error"])

	mockUsecase.AssertExpectations(t)
}

// GetIssue Tests

func TestIssueHandler_GetIssue_Success(t *testing.T) {
	handler, mockUsecase := setupIssueHandler()

	issue := createTestIssueForHandler(1, 1, "Test Issue", "TODO")
	mockUsecase.On("GetByID", mock.Anything, int64(1)).Return(&issue, nil)

	w := httptest.NewRecorder()
	_, router := gin.CreateTestContext(w)
	router.GET("/issues/:id", handler.GetIssue)

	req := httptest.NewRequest(http.MethodGet, "/issues/1", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response domain.Issue
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, int64(1), response.ID)
	assert.Equal(t, "Test Issue", response.Summary)

	mockUsecase.AssertExpectations(t)
}

func TestIssueHandler_GetIssue_NotFound(t *testing.T) {
	handler, mockUsecase := setupIssueHandler()

	mockUsecase.On("GetByID", mock.Anything, int64(999)).Return(nil, nil)

	w := httptest.NewRecorder()
	_, router := gin.CreateTestContext(w)
	router.GET("/issues/:id", handler.GetIssue)

	req := httptest.NewRequest(http.MethodGet, "/issues/999", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)

	var response map[string]string
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, "Issue not found", response["error"])

	mockUsecase.AssertExpectations(t)
}

func TestIssueHandler_GetIssue_InvalidID(t *testing.T) {
	handler, _ := setupIssueHandler()

	w := httptest.NewRecorder()
	_, router := gin.CreateTestContext(w)
	router.GET("/issues/:id", handler.GetIssue)

	req := httptest.NewRequest(http.MethodGet, "/issues/invalid", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]string
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, "Invalid issue ID", response["error"])
}

func TestIssueHandler_GetIssue_UsecaseError(t *testing.T) {
	handler, mockUsecase := setupIssueHandler()

	mockUsecase.On("GetByID", mock.Anything, int64(1)).Return(nil, errors.New("database error"))

	w := httptest.NewRecorder()
	_, router := gin.CreateTestContext(w)
	router.GET("/issues/:id", handler.GetIssue)

	req := httptest.NewRequest(http.MethodGet, "/issues/1", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)

	var response map[string]string
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, "Failed to get issue", response["error"])

	mockUsecase.AssertExpectations(t)
}

// ListProjectIssues Tests

func TestIssueHandler_ListProjectIssues_Success(t *testing.T) {
	handler, mockUsecase := setupIssueHandler()

	issues := []domain.Issue{
		createTestIssueForHandler(1, 1, "Issue 1", "TODO"),
		createTestIssueForHandler(2, 1, "Issue 2", "IN_PROGRESS"),
	}

	mockUsecase.On("GetByProjectID", mock.Anything, int64(1)).Return(issues, nil)

	w := httptest.NewRecorder()
	_, router := gin.CreateTestContext(w)
	router.GET("/projects/:id/issues", handler.ListProjectIssues)

	req := httptest.NewRequest(http.MethodGet, "/projects/1/issues", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, float64(2), response["count"])

	mockUsecase.AssertExpectations(t)
}

func TestIssueHandler_ListProjectIssues_InvalidProjectID(t *testing.T) {
	handler, _ := setupIssueHandler()

	w := httptest.NewRecorder()
	_, router := gin.CreateTestContext(w)
	router.GET("/projects/:id/issues", handler.ListProjectIssues)

	req := httptest.NewRequest(http.MethodGet, "/projects/invalid/issues", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]string
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, "Invalid project ID", response["error"])
}

func TestIssueHandler_ListProjectIssues_UsecaseError(t *testing.T) {
	handler, mockUsecase := setupIssueHandler()

	mockUsecase.On("GetByProjectID", mock.Anything, int64(1)).Return([]domain.Issue{}, errors.New("database error"))

	w := httptest.NewRecorder()
	_, router := gin.CreateTestContext(w)
	router.GET("/projects/:id/issues", handler.ListProjectIssues)

	req := httptest.NewRequest(http.MethodGet, "/projects/1/issues", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)

	var response map[string]string
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, "Failed to get issues", response["error"])

	mockUsecase.AssertExpectations(t)
}
