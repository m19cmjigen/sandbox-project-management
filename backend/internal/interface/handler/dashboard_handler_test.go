package handler

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/m19cmjigen/sandbox-project-management/backend/internal/domain"
	"github.com/m19cmjigen/sandbox-project-management/backend/internal/usecase"
	"github.com/m19cmjigen/sandbox-project-management/backend/pkg/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// MockDashboardUsecase is a mock implementation of DashboardUsecase
type MockDashboardUsecase struct {
	mock.Mock
}

func (m *MockDashboardUsecase) GetSummary(ctx context.Context) (*usecase.DashboardSummary, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*usecase.DashboardSummary), args.Error(1)
}

func (m *MockDashboardUsecase) GetOrganizationSummary(ctx context.Context, organizationID int64) (*usecase.OrganizationSummary, error) {
	args := m.Called(ctx, organizationID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*usecase.OrganizationSummary), args.Error(1)
}

func (m *MockDashboardUsecase) GetProjectSummary(ctx context.Context, projectID int64) (*domain.ProjectWithStats, error) {
	args := m.Called(ctx, projectID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.ProjectWithStats), args.Error(1)
}

// Helper functions
func setupDashboardHandler() (*DashboardHandler, *MockDashboardUsecase) {
	gin.SetMode(gin.TestMode)
	mockUsecase := new(MockDashboardUsecase)
	log, _ := logger.New("error", "console") // Use error level to suppress logs during tests
	handler := NewDashboardHandler(mockUsecase, log)
	return handler, mockUsecase
}

func createTestProjectWithStatsForHandler(id int64, name string, red, yellow, green int) *domain.ProjectWithStats {
	orgID := int64(1)
	return &domain.ProjectWithStats{
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

// GetDashboardSummary tests
func TestDashboardHandler_GetDashboardSummary_Success(t *testing.T) {
	handler, mockUsecase := setupDashboardHandler()

	summary := &usecase.DashboardSummary{
		TotalProjects:   10,
		DelayedProjects: 3,
		WarningProjects: 2,
		NormalProjects:  5,
		TotalIssues:     100,
		RedIssues:       15,
		YellowIssues:    25,
		GreenIssues:     60,
		ProjectsByStatus: []domain.ProjectWithStats{
			*createTestProjectWithStatsForHandler(1, "PROJ1", 5, 3, 12),
		},
	}

	mockUsecase.On("GetSummary", mock.Anything).Return(summary, nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/api/dashboard", nil)

	handler.GetDashboardSummary(c)

	assert.Equal(t, http.StatusOK, w.Code)

	var response usecase.DashboardSummary
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, 10, response.TotalProjects)
	assert.Equal(t, 3, response.DelayedProjects)
	assert.Equal(t, 100, response.TotalIssues)
	mockUsecase.AssertExpectations(t)
}

func TestDashboardHandler_GetDashboardSummary_Error(t *testing.T) {
	handler, mockUsecase := setupDashboardHandler()

	mockUsecase.On("GetSummary", mock.Anything).Return(nil, errors.New("database error"))

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/api/dashboard", nil)

	handler.GetDashboardSummary(c)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	mockUsecase.AssertExpectations(t)
}

// GetOrganizationSummary tests
func TestDashboardHandler_GetOrganizationSummary_Success(t *testing.T) {
	handler, mockUsecase := setupDashboardHandler()

	org := domain.Organization{
		ID:        1,
		Name:      "Test Organization",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	summary := &usecase.OrganizationSummary{
		Organization:    org,
		TotalProjects:   5,
		DelayedProjects: 2,
		WarningProjects: 1,
		Projects: []domain.ProjectWithStats{
			*createTestProjectWithStatsForHandler(1, "PROJ1", 5, 3, 12),
		},
	}

	mockUsecase.On("GetOrganizationSummary", mock.Anything, int64(1)).Return(summary, nil)

	w := httptest.NewRecorder()
	_, router := gin.CreateTestContext(w)

	router.GET("/organizations/:id/summary", handler.GetOrganizationSummary)

	req := httptest.NewRequest(http.MethodGet, "/organizations/1/summary", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response usecase.OrganizationSummary
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, int64(1), response.Organization.ID)
	assert.Equal(t, "Test Organization", response.Organization.Name)
	assert.Equal(t, 5, response.TotalProjects)
	mockUsecase.AssertExpectations(t)
}

func TestDashboardHandler_GetOrganizationSummary_NotFound(t *testing.T) {
	handler, mockUsecase := setupDashboardHandler()

	mockUsecase.On("GetOrganizationSummary", mock.Anything, int64(999)).Return(nil, nil)

	w := httptest.NewRecorder()
	_, router := gin.CreateTestContext(w)

	router.GET("/organizations/:id/summary", handler.GetOrganizationSummary)

	req := httptest.NewRequest(http.MethodGet, "/organizations/999/summary", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	mockUsecase.AssertExpectations(t)
}

func TestDashboardHandler_GetOrganizationSummary_InvalidID(t *testing.T) {
	handler, _ := setupDashboardHandler()

	w := httptest.NewRecorder()
	_, router := gin.CreateTestContext(w)

	router.GET("/organizations/:id/summary", handler.GetOrganizationSummary)

	req := httptest.NewRequest(http.MethodGet, "/organizations/invalid/summary", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestDashboardHandler_GetOrganizationSummary_Error(t *testing.T) {
	handler, mockUsecase := setupDashboardHandler()

	mockUsecase.On("GetOrganizationSummary", mock.Anything, int64(1)).Return(nil, errors.New("database error"))

	w := httptest.NewRecorder()
	_, router := gin.CreateTestContext(w)

	router.GET("/organizations/:id/summary", handler.GetOrganizationSummary)

	req := httptest.NewRequest(http.MethodGet, "/organizations/1/summary", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	mockUsecase.AssertExpectations(t)
}

// GetProjectSummary tests
func TestDashboardHandler_GetProjectSummary_Success(t *testing.T) {
	handler, mockUsecase := setupDashboardHandler()

	summary := createTestProjectWithStatsForHandler(1, "PROJ1", 5, 3, 12)

	mockUsecase.On("GetProjectSummary", mock.Anything, int64(1)).Return(summary, nil)

	w := httptest.NewRecorder()
	_, router := gin.CreateTestContext(w)

	router.GET("/projects/:id/summary", handler.GetProjectSummary)

	req := httptest.NewRequest(http.MethodGet, "/projects/1/summary", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response domain.ProjectWithStats
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, int64(1), response.ID)
	assert.Equal(t, "PROJ1", response.Key)
	assert.Equal(t, 20, response.TotalIssues)
	assert.Equal(t, 5, response.RedIssues)
	assert.Equal(t, 3, response.YellowIssues)
	assert.Equal(t, 12, response.GreenIssues)
	mockUsecase.AssertExpectations(t)
}

func TestDashboardHandler_GetProjectSummary_NotFound(t *testing.T) {
	handler, mockUsecase := setupDashboardHandler()

	mockUsecase.On("GetProjectSummary", mock.Anything, int64(999)).Return(nil, nil)

	w := httptest.NewRecorder()
	_, router := gin.CreateTestContext(w)

	router.GET("/projects/:id/summary", handler.GetProjectSummary)

	req := httptest.NewRequest(http.MethodGet, "/projects/999/summary", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	mockUsecase.AssertExpectations(t)
}

func TestDashboardHandler_GetProjectSummary_InvalidID(t *testing.T) {
	handler, _ := setupDashboardHandler()

	w := httptest.NewRecorder()
	_, router := gin.CreateTestContext(w)

	router.GET("/projects/:id/summary", handler.GetProjectSummary)

	req := httptest.NewRequest(http.MethodGet, "/projects/invalid/summary", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestDashboardHandler_GetProjectSummary_Error(t *testing.T) {
	handler, mockUsecase := setupDashboardHandler()

	mockUsecase.On("GetProjectSummary", mock.Anything, int64(1)).Return(nil, errors.New("database error"))

	w := httptest.NewRecorder()
	_, router := gin.CreateTestContext(w)

	router.GET("/projects/:id/summary", handler.GetProjectSummary)

	req := httptest.NewRequest(http.MethodGet, "/projects/1/summary", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	mockUsecase.AssertExpectations(t)
}
