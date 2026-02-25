package handler

import (
	"bytes"
	"context"
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

// MockProjectUsecase はProjectUsecaseのモック
type MockProjectUsecase struct {
	mock.Mock
}

func (m *MockProjectUsecase) GetAll(ctx context.Context) ([]domain.Project, error) {
	args := m.Called(ctx)
	return args.Get(0).([]domain.Project), args.Error(1)
}

func (m *MockProjectUsecase) GetAllWithStats(ctx context.Context) ([]domain.ProjectWithStats, error) {
	args := m.Called(ctx)
	return args.Get(0).([]domain.ProjectWithStats), args.Error(1)
}

func (m *MockProjectUsecase) GetByID(ctx context.Context, id int64) (*domain.Project, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Project), args.Error(1)
}

func (m *MockProjectUsecase) GetWithStats(ctx context.Context, id int64) (*domain.ProjectWithStats, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.ProjectWithStats), args.Error(1)
}

func (m *MockProjectUsecase) GetByOrganization(ctx context.Context, organizationID int64) ([]domain.Project, error) {
	args := m.Called(ctx, organizationID)
	return args.Get(0).([]domain.Project), args.Error(1)
}

func (m *MockProjectUsecase) GetUnassigned(ctx context.Context) ([]domain.Project, error) {
	args := m.Called(ctx)
	return args.Get(0).([]domain.Project), args.Error(1)
}

func (m *MockProjectUsecase) AssignToOrganization(ctx context.Context, projectID int64, organizationID *int64) error {
	args := m.Called(ctx, projectID, organizationID)
	return args.Error(0)
}

func setupProjectHandler() (*ProjectHandler, *MockProjectUsecase) {
	gin.SetMode(gin.TestMode)
	mockUsecase := new(MockProjectUsecase)
	log, _ := logger.New("error", "console")
	handler := NewProjectHandler(mockUsecase, log)
	return handler, mockUsecase
}

func createTestProjectForHandler(id int64, jiraProjectID string, name string) domain.Project {
	now := time.Now()
	key := "TEST-" + jiraProjectID
	return domain.Project{
		ID:              id,
		JiraProjectID:   jiraProjectID,
		Key:             key,
		Name:            name,
		LeadAccountID:   nil,
		LeadEmail:       nil,
		OrganizationID:  nil,
		CreatedAt:       now,
		UpdatedAt:       now,
	}
}

func createTestProjectWithStatsForHandler(id int64, jiraProjectID string, name string) domain.ProjectWithStats {
	project := createTestProjectForHandler(id, jiraProjectID, name)
	return domain.ProjectWithStats{
		Project:      project,
		TotalIssues:  10,
		RedIssues:    2,
		YellowIssues: 3,
		GreenIssues:  5,
		OpenIssues:   8,
		DoneIssues:   2,
	}
}

// ListProjects Tests

func TestProjectHandler_ListProjects_Success_Default(t *testing.T) {
	handler, mockUsecase := setupProjectHandler()

	projects := []domain.Project{
		createTestProjectForHandler(1, "PROJ1", "Project 1"),
		createTestProjectForHandler(2, "PROJ2", "Project 2"),
	}

	mockUsecase.On("GetAll", mock.Anything).Return(projects, nil)

	w := httptest.NewRecorder()
	_, router := gin.CreateTestContext(w)
	router.GET("/projects", handler.ListProjects)

	req := httptest.NewRequest(http.MethodGet, "/projects", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	projects_data := response["projects"].([]interface{})
	assert.Equal(t, 2, len(projects_data))

	mockUsecase.AssertExpectations(t)
}

func TestProjectHandler_ListProjects_Success_WithStats(t *testing.T) {
	handler, mockUsecase := setupProjectHandler()

	projects := []domain.ProjectWithStats{
		createTestProjectWithStatsForHandler(1, "PROJ1", "Project 1"),
		createTestProjectWithStatsForHandler(2, "PROJ2", "Project 2"),
	}

	mockUsecase.On("GetAllWithStats", mock.Anything).Return(projects, nil)

	w := httptest.NewRecorder()
	_, router := gin.CreateTestContext(w)
	router.GET("/projects", handler.ListProjects)

	req := httptest.NewRequest(http.MethodGet, "/projects?with_stats=true", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	projects_data := response["projects"].([]interface{})
	assert.Equal(t, 2, len(projects_data))

	mockUsecase.AssertExpectations(t)
}

func TestProjectHandler_ListProjects_Success_Unassigned(t *testing.T) {
	handler, mockUsecase := setupProjectHandler()

	projects := []domain.Project{
		createTestProjectForHandler(1, "PROJ1", "Unassigned Project"),
	}

	mockUsecase.On("GetUnassigned", mock.Anything).Return(projects, nil)

	w := httptest.NewRecorder()
	_, router := gin.CreateTestContext(w)
	router.GET("/projects", handler.ListProjects)

	req := httptest.NewRequest(http.MethodGet, "/projects?unassigned=true", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	projects_data := response["projects"].([]interface{})
	assert.Equal(t, 1, len(projects_data))

	mockUsecase.AssertExpectations(t)
}

func TestProjectHandler_ListProjects_Success_ByOrganization(t *testing.T) {
	handler, mockUsecase := setupProjectHandler()

	projects := []domain.Project{
		createTestProjectForHandler(1, "PROJ1", "Org Project 1"),
		createTestProjectForHandler(2, "PROJ2", "Org Project 2"),
	}

	mockUsecase.On("GetByOrganization", mock.Anything, int64(1)).Return(projects, nil)

	w := httptest.NewRecorder()
	_, router := gin.CreateTestContext(w)
	router.GET("/projects", handler.ListProjects)

	req := httptest.NewRequest(http.MethodGet, "/projects?organization_id=1", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	projects_data := response["projects"].([]interface{})
	assert.Equal(t, 2, len(projects_data))

	mockUsecase.AssertExpectations(t)
}

func TestProjectHandler_ListProjects_InvalidOrganizationID(t *testing.T) {
	handler, _ := setupProjectHandler()

	w := httptest.NewRecorder()
	_, router := gin.CreateTestContext(w)
	router.GET("/projects", handler.ListProjects)

	req := httptest.NewRequest(http.MethodGet, "/projects?organization_id=invalid", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]string
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, "Invalid organization ID", response["error"])
}

func TestProjectHandler_ListProjects_UsecaseError_Default(t *testing.T) {
	handler, mockUsecase := setupProjectHandler()

	mockUsecase.On("GetAll", mock.Anything).Return([]domain.Project{}, errors.New("database error"))

	w := httptest.NewRecorder()
	_, router := gin.CreateTestContext(w)
	router.GET("/projects", handler.ListProjects)

	req := httptest.NewRequest(http.MethodGet, "/projects", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)

	var response map[string]string
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, "Failed to get projects", response["error"])

	mockUsecase.AssertExpectations(t)
}

func TestProjectHandler_ListProjects_UsecaseError_WithStats(t *testing.T) {
	handler, mockUsecase := setupProjectHandler()

	mockUsecase.On("GetAllWithStats", mock.Anything).Return([]domain.ProjectWithStats{}, errors.New("database error"))

	w := httptest.NewRecorder()
	_, router := gin.CreateTestContext(w)
	router.GET("/projects", handler.ListProjects)

	req := httptest.NewRequest(http.MethodGet, "/projects?with_stats=true", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)

	var response map[string]string
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, "Failed to get projects", response["error"])

	mockUsecase.AssertExpectations(t)
}

func TestProjectHandler_ListProjects_UsecaseError_Unassigned(t *testing.T) {
	handler, mockUsecase := setupProjectHandler()

	mockUsecase.On("GetUnassigned", mock.Anything).Return([]domain.Project{}, errors.New("database error"))

	w := httptest.NewRecorder()
	_, router := gin.CreateTestContext(w)
	router.GET("/projects", handler.ListProjects)

	req := httptest.NewRequest(http.MethodGet, "/projects?unassigned=true", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)

	var response map[string]string
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, "Failed to get projects", response["error"])

	mockUsecase.AssertExpectations(t)
}

func TestProjectHandler_ListProjects_UsecaseError_ByOrganization(t *testing.T) {
	handler, mockUsecase := setupProjectHandler()

	mockUsecase.On("GetByOrganization", mock.Anything, int64(1)).Return([]domain.Project{}, errors.New("database error"))

	w := httptest.NewRecorder()
	_, router := gin.CreateTestContext(w)
	router.GET("/projects", handler.ListProjects)

	req := httptest.NewRequest(http.MethodGet, "/projects?organization_id=1", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)

	var response map[string]string
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, "Failed to get projects", response["error"])

	mockUsecase.AssertExpectations(t)
}

// GetProject Tests

func TestProjectHandler_GetProject_Success(t *testing.T) {
	handler, mockUsecase := setupProjectHandler()

	project := createTestProjectForHandler(1, "PROJ1", "Test Project")
	mockUsecase.On("GetByID", mock.Anything, int64(1)).Return(&project, nil)

	w := httptest.NewRecorder()
	_, router := gin.CreateTestContext(w)
	router.GET("/projects/:id", handler.GetProject)

	req := httptest.NewRequest(http.MethodGet, "/projects/1", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response domain.Project
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, int64(1), response.ID)
	assert.Equal(t, "Test Project", response.Name)

	mockUsecase.AssertExpectations(t)
}

func TestProjectHandler_GetProject_Success_WithStats(t *testing.T) {
	handler, mockUsecase := setupProjectHandler()

	project := createTestProjectWithStatsForHandler(1, "PROJ1", "Test Project")
	mockUsecase.On("GetWithStats", mock.Anything, int64(1)).Return(&project, nil)

	w := httptest.NewRecorder()
	_, router := gin.CreateTestContext(w)
	router.GET("/projects/:id", handler.GetProject)

	req := httptest.NewRequest(http.MethodGet, "/projects/1?with_stats=true", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response domain.ProjectWithStats
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, int64(1), response.Project.ID)
	assert.Equal(t, "Test Project", response.Project.Name)
	assert.Equal(t, 10, response.TotalIssues)

	mockUsecase.AssertExpectations(t)
}

func TestProjectHandler_GetProject_NotFound(t *testing.T) {
	handler, mockUsecase := setupProjectHandler()

	mockUsecase.On("GetByID", mock.Anything, int64(999)).Return(nil, nil)

	w := httptest.NewRecorder()
	_, router := gin.CreateTestContext(w)
	router.GET("/projects/:id", handler.GetProject)

	req := httptest.NewRequest(http.MethodGet, "/projects/999", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)

	var response map[string]string
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, "Project not found", response["error"])

	mockUsecase.AssertExpectations(t)
}

func TestProjectHandler_GetProject_NotFound_WithStats(t *testing.T) {
	handler, mockUsecase := setupProjectHandler()

	mockUsecase.On("GetWithStats", mock.Anything, int64(999)).Return(nil, nil)

	w := httptest.NewRecorder()
	_, router := gin.CreateTestContext(w)
	router.GET("/projects/:id", handler.GetProject)

	req := httptest.NewRequest(http.MethodGet, "/projects/999?with_stats=true", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)

	var response map[string]string
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, "Project not found", response["error"])

	mockUsecase.AssertExpectations(t)
}

func TestProjectHandler_GetProject_InvalidID(t *testing.T) {
	handler, _ := setupProjectHandler()

	w := httptest.NewRecorder()
	_, router := gin.CreateTestContext(w)
	router.GET("/projects/:id", handler.GetProject)

	req := httptest.NewRequest(http.MethodGet, "/projects/invalid", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]string
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, "Invalid project ID", response["error"])
}

func TestProjectHandler_GetProject_UsecaseError(t *testing.T) {
	handler, mockUsecase := setupProjectHandler()

	mockUsecase.On("GetByID", mock.Anything, int64(1)).Return(nil, errors.New("database error"))

	w := httptest.NewRecorder()
	_, router := gin.CreateTestContext(w)
	router.GET("/projects/:id", handler.GetProject)

	req := httptest.NewRequest(http.MethodGet, "/projects/1", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)

	var response map[string]string
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, "Failed to get project", response["error"])

	mockUsecase.AssertExpectations(t)
}

func TestProjectHandler_GetProject_UsecaseError_WithStats(t *testing.T) {
	handler, mockUsecase := setupProjectHandler()

	mockUsecase.On("GetWithStats", mock.Anything, int64(1)).Return(nil, errors.New("database error"))

	w := httptest.NewRecorder()
	_, router := gin.CreateTestContext(w)
	router.GET("/projects/:id", handler.GetProject)

	req := httptest.NewRequest(http.MethodGet, "/projects/1?with_stats=true", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)

	var response map[string]string
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, "Failed to get project", response["error"])

	mockUsecase.AssertExpectations(t)
}

// AssignProjectToOrganization Tests

func TestProjectHandler_AssignProjectToOrganization_Success(t *testing.T) {
	handler, mockUsecase := setupProjectHandler()

	orgID := int64(1)
	assignReq := AssignProjectRequest{
		OrganizationID: &orgID,
	}

	mockUsecase.On("AssignToOrganization", mock.Anything, int64(1), &orgID).Return(nil)

	w := httptest.NewRecorder()
	_, router := gin.CreateTestContext(w)
	router.POST("/projects/:id/assign", handler.AssignProjectToOrganization)

	body, _ := json.Marshal(assignReq)
	req := httptest.NewRequest(http.MethodPost, "/projects/1/assign", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]string
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, "Project assigned successfully", response["message"])

	mockUsecase.AssertExpectations(t)
}

func TestProjectHandler_AssignProjectToOrganization_Success_Unassign(t *testing.T) {
	handler, mockUsecase := setupProjectHandler()

	assignReq := AssignProjectRequest{
		OrganizationID: nil,
	}

	mockUsecase.On("AssignToOrganization", mock.Anything, int64(1), (*int64)(nil)).Return(nil)

	w := httptest.NewRecorder()
	_, router := gin.CreateTestContext(w)
	router.POST("/projects/:id/assign", handler.AssignProjectToOrganization)

	body, _ := json.Marshal(assignReq)
	req := httptest.NewRequest(http.MethodPost, "/projects/1/assign", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	mockUsecase.AssertExpectations(t)
}

func TestProjectHandler_AssignProjectToOrganization_InvalidProjectID(t *testing.T) {
	handler, _ := setupProjectHandler()

	assignReq := AssignProjectRequest{
		OrganizationID: nil,
	}

	w := httptest.NewRecorder()
	_, router := gin.CreateTestContext(w)
	router.POST("/projects/:id/assign", handler.AssignProjectToOrganization)

	body, _ := json.Marshal(assignReq)
	req := httptest.NewRequest(http.MethodPost, "/projects/invalid/assign", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]string
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, "Invalid project ID", response["error"])
}

func TestProjectHandler_AssignProjectToOrganization_InvalidRequestBody(t *testing.T) {
	handler, _ := setupProjectHandler()

	w := httptest.NewRecorder()
	_, router := gin.CreateTestContext(w)
	router.POST("/projects/:id/assign", handler.AssignProjectToOrganization)

	req := httptest.NewRequest(http.MethodPost, "/projects/1/assign", bytes.NewBuffer([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestProjectHandler_AssignProjectToOrganization_UsecaseError(t *testing.T) {
	handler, mockUsecase := setupProjectHandler()

	orgID := int64(1)
	assignReq := AssignProjectRequest{
		OrganizationID: &orgID,
	}

	mockUsecase.On("AssignToOrganization", mock.Anything, int64(1), &orgID).Return(errors.New("project not found"))

	w := httptest.NewRecorder()
	_, router := gin.CreateTestContext(w)
	router.POST("/projects/:id/assign", handler.AssignProjectToOrganization)

	body, _ := json.Marshal(assignReq)
	req := httptest.NewRequest(http.MethodPost, "/projects/1/assign", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)

	var response map[string]string
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, "project not found", response["error"])

	mockUsecase.AssertExpectations(t)
}
