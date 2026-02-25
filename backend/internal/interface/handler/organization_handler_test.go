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

// MockOrganizationUsecase はOrganizationUsecaseのモック
type MockOrganizationUsecase struct {
	mock.Mock
}

func (m *MockOrganizationUsecase) GetAll(ctx context.Context) ([]domain.Organization, error) {
	args := m.Called(ctx)
	return args.Get(0).([]domain.Organization), args.Error(1)
}

func (m *MockOrganizationUsecase) GetByID(ctx context.Context, id int64) (*domain.Organization, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Organization), args.Error(1)
}

func (m *MockOrganizationUsecase) GetChildren(ctx context.Context, parentID int64) ([]domain.Organization, error) {
	args := m.Called(ctx, parentID)
	return args.Get(0).([]domain.Organization), args.Error(1)
}

func (m *MockOrganizationUsecase) GetRoots(ctx context.Context) ([]domain.Organization, error) {
	args := m.Called(ctx)
	return args.Get(0).([]domain.Organization), args.Error(1)
}

func (m *MockOrganizationUsecase) GetTree(ctx context.Context) ([]domain.OrganizationWithChildren, error) {
	args := m.Called(ctx)
	return args.Get(0).([]domain.OrganizationWithChildren), args.Error(1)
}

func (m *MockOrganizationUsecase) Create(ctx context.Context, name string, parentID *int64) (*domain.Organization, error) {
	args := m.Called(ctx, name, parentID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Organization), args.Error(1)
}

func (m *MockOrganizationUsecase) Update(ctx context.Context, id int64, name string, parentID *int64) (*domain.Organization, error) {
	args := m.Called(ctx, id, name, parentID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Organization), args.Error(1)
}

func (m *MockOrganizationUsecase) Delete(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func setupOrganizationHandler() (*OrganizationHandler, *MockOrganizationUsecase) {
	gin.SetMode(gin.TestMode)
	mockUsecase := new(MockOrganizationUsecase)
	log, _ := logger.New("error", "console")
	handler := NewOrganizationHandler(mockUsecase, log)
	return handler, mockUsecase
}

func createTestOrganizationForHandler(id int64, name string, parentID *int64) domain.Organization {
	now := time.Now()
	path := "/" + name
	level := 1
	if parentID != nil {
		level = 2
		path = "/parent/" + name
	}
	return domain.Organization{
		ID:        id,
		Name:      name,
		ParentID:  parentID,
		Path:      path,
		Level:     level,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

func createTestOrganizationWithChildren(id int64, name string, children []domain.Organization) domain.OrganizationWithChildren {
	org := createTestOrganizationForHandler(id, name, nil)
	return domain.OrganizationWithChildren{
		Organization: org,
		Children:     children,
	}
}

// ListOrganizations Tests

func TestOrganizationHandler_ListOrganizations_Success(t *testing.T) {
	handler, mockUsecase := setupOrganizationHandler()

	orgs := []domain.Organization{
		createTestOrganizationForHandler(1, "Org 1", nil),
		createTestOrganizationForHandler(2, "Org 2", nil),
	}

	mockUsecase.On("GetAll", mock.Anything).Return(orgs, nil)

	w := httptest.NewRecorder()
	_, router := gin.CreateTestContext(w)
	router.GET("/organizations", handler.ListOrganizations)

	req := httptest.NewRequest(http.MethodGet, "/organizations", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	orgs_data := response["organizations"].([]interface{})
	assert.Equal(t, 2, len(orgs_data))

	mockUsecase.AssertExpectations(t)
}

func TestOrganizationHandler_ListOrganizations_UsecaseError(t *testing.T) {
	handler, mockUsecase := setupOrganizationHandler()

	mockUsecase.On("GetAll", mock.Anything).Return([]domain.Organization{}, errors.New("database error"))

	w := httptest.NewRecorder()
	_, router := gin.CreateTestContext(w)
	router.GET("/organizations", handler.ListOrganizations)

	req := httptest.NewRequest(http.MethodGet, "/organizations", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)

	var response map[string]string
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, "Failed to get organizations", response["error"])

	mockUsecase.AssertExpectations(t)
}

// GetOrganization Tests

func TestOrganizationHandler_GetOrganization_Success(t *testing.T) {
	handler, mockUsecase := setupOrganizationHandler()

	org := createTestOrganizationForHandler(1, "Test Org", nil)
	mockUsecase.On("GetByID", mock.Anything, int64(1)).Return(&org, nil)

	w := httptest.NewRecorder()
	_, router := gin.CreateTestContext(w)
	router.GET("/organizations/:id", handler.GetOrganization)

	req := httptest.NewRequest(http.MethodGet, "/organizations/1", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response domain.Organization
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, int64(1), response.ID)
	assert.Equal(t, "Test Org", response.Name)

	mockUsecase.AssertExpectations(t)
}

func TestOrganizationHandler_GetOrganization_NotFound(t *testing.T) {
	handler, mockUsecase := setupOrganizationHandler()

	mockUsecase.On("GetByID", mock.Anything, int64(999)).Return(nil, nil)

	w := httptest.NewRecorder()
	_, router := gin.CreateTestContext(w)
	router.GET("/organizations/:id", handler.GetOrganization)

	req := httptest.NewRequest(http.MethodGet, "/organizations/999", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)

	var response map[string]string
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, "Organization not found", response["error"])

	mockUsecase.AssertExpectations(t)
}

func TestOrganizationHandler_GetOrganization_InvalidID(t *testing.T) {
	handler, _ := setupOrganizationHandler()

	w := httptest.NewRecorder()
	_, router := gin.CreateTestContext(w)
	router.GET("/organizations/:id", handler.GetOrganization)

	req := httptest.NewRequest(http.MethodGet, "/organizations/invalid", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]string
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, "Invalid organization ID", response["error"])
}

func TestOrganizationHandler_GetOrganization_UsecaseError(t *testing.T) {
	handler, mockUsecase := setupOrganizationHandler()

	mockUsecase.On("GetByID", mock.Anything, int64(1)).Return(nil, errors.New("database error"))

	w := httptest.NewRecorder()
	_, router := gin.CreateTestContext(w)
	router.GET("/organizations/:id", handler.GetOrganization)

	req := httptest.NewRequest(http.MethodGet, "/organizations/1", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)

	var response map[string]string
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, "Failed to get organization", response["error"])

	mockUsecase.AssertExpectations(t)
}

// GetOrganizationChildren Tests

func TestOrganizationHandler_GetOrganizationChildren_Success(t *testing.T) {
	handler, mockUsecase := setupOrganizationHandler()

	parentID := int64(1)
	children := []domain.Organization{
		createTestOrganizationForHandler(2, "Child 1", &parentID),
		createTestOrganizationForHandler(3, "Child 2", &parentID),
	}

	mockUsecase.On("GetChildren", mock.Anything, int64(1)).Return(children, nil)

	w := httptest.NewRecorder()
	_, router := gin.CreateTestContext(w)
	router.GET("/organizations/:id/children", handler.GetOrganizationChildren)

	req := httptest.NewRequest(http.MethodGet, "/organizations/1/children", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	children_data := response["children"].([]interface{})
	assert.Equal(t, 2, len(children_data))

	mockUsecase.AssertExpectations(t)
}

func TestOrganizationHandler_GetOrganizationChildren_InvalidID(t *testing.T) {
	handler, _ := setupOrganizationHandler()

	w := httptest.NewRecorder()
	_, router := gin.CreateTestContext(w)
	router.GET("/organizations/:id/children", handler.GetOrganizationChildren)

	req := httptest.NewRequest(http.MethodGet, "/organizations/invalid/children", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]string
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, "Invalid organization ID", response["error"])
}

func TestOrganizationHandler_GetOrganizationChildren_UsecaseError(t *testing.T) {
	handler, mockUsecase := setupOrganizationHandler()

	mockUsecase.On("GetChildren", mock.Anything, int64(1)).Return([]domain.Organization{}, errors.New("database error"))

	w := httptest.NewRecorder()
	_, router := gin.CreateTestContext(w)
	router.GET("/organizations/:id/children", handler.GetOrganizationChildren)

	req := httptest.NewRequest(http.MethodGet, "/organizations/1/children", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)

	var response map[string]string
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, "Failed to get children", response["error"])

	mockUsecase.AssertExpectations(t)
}

// GetOrganizationTree Tests

func TestOrganizationHandler_GetOrganizationTree_Success(t *testing.T) {
	handler, mockUsecase := setupOrganizationHandler()

	tree := []domain.OrganizationWithChildren{
		createTestOrganizationWithChildren(1, "Root 1", []domain.Organization{}),
		createTestOrganizationWithChildren(2, "Root 2", []domain.Organization{}),
	}

	mockUsecase.On("GetTree", mock.Anything).Return(tree, nil)

	w := httptest.NewRecorder()
	_, router := gin.CreateTestContext(w)
	router.GET("/organizations/tree", handler.GetOrganizationTree)

	req := httptest.NewRequest(http.MethodGet, "/organizations/tree", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	tree_data := response["tree"].([]interface{})
	assert.Equal(t, 2, len(tree_data))

	mockUsecase.AssertExpectations(t)
}

func TestOrganizationHandler_GetOrganizationTree_UsecaseError(t *testing.T) {
	handler, mockUsecase := setupOrganizationHandler()

	mockUsecase.On("GetTree", mock.Anything).Return([]domain.OrganizationWithChildren{}, errors.New("database error"))

	w := httptest.NewRecorder()
	_, router := gin.CreateTestContext(w)
	router.GET("/organizations/tree", handler.GetOrganizationTree)

	req := httptest.NewRequest(http.MethodGet, "/organizations/tree", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)

	var response map[string]string
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, "Failed to get organization tree", response["error"])

	mockUsecase.AssertExpectations(t)
}

// CreateOrganization Tests

func TestOrganizationHandler_CreateOrganization_Success(t *testing.T) {
	handler, mockUsecase := setupOrganizationHandler()

	createReq := CreateOrganizationRequest{
		Name:     "New Org",
		ParentID: nil,
	}

	createdOrg := createTestOrganizationForHandler(1, "New Org", nil)
	mockUsecase.On("Create", mock.Anything, "New Org", (*int64)(nil)).Return(&createdOrg, nil)

	w := httptest.NewRecorder()
	_, router := gin.CreateTestContext(w)
	router.POST("/organizations", handler.CreateOrganization)

	body, _ := json.Marshal(createReq)
	req := httptest.NewRequest(http.MethodPost, "/organizations", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var response domain.Organization
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, int64(1), response.ID)
	assert.Equal(t, "New Org", response.Name)

	mockUsecase.AssertExpectations(t)
}

func TestOrganizationHandler_CreateOrganization_InvalidRequestBody(t *testing.T) {
	handler, _ := setupOrganizationHandler()

	w := httptest.NewRecorder()
	_, router := gin.CreateTestContext(w)
	router.POST("/organizations", handler.CreateOrganization)

	// Missing required "name" field
	body := []byte(`{"parent_id": 1}`)
	req := httptest.NewRequest(http.MethodPost, "/organizations", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestOrganizationHandler_CreateOrganization_UsecaseError(t *testing.T) {
	handler, mockUsecase := setupOrganizationHandler()

	createReq := CreateOrganizationRequest{
		Name:     "New Org",
		ParentID: nil,
	}

	mockUsecase.On("Create", mock.Anything, "New Org", (*int64)(nil)).Return(nil, errors.New("parent not found"))

	w := httptest.NewRecorder()
	_, router := gin.CreateTestContext(w)
	router.POST("/organizations", handler.CreateOrganization)

	body, _ := json.Marshal(createReq)
	req := httptest.NewRequest(http.MethodPost, "/organizations", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)

	var response map[string]string
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, "parent not found", response["error"])

	mockUsecase.AssertExpectations(t)
}

// UpdateOrganization Tests

func TestOrganizationHandler_UpdateOrganization_Success(t *testing.T) {
	handler, mockUsecase := setupOrganizationHandler()

	updateReq := UpdateOrganizationRequest{
		Name:     "Updated Org",
		ParentID: nil,
	}

	updatedOrg := createTestOrganizationForHandler(1, "Updated Org", nil)
	mockUsecase.On("Update", mock.Anything, int64(1), "Updated Org", (*int64)(nil)).Return(&updatedOrg, nil)

	w := httptest.NewRecorder()
	_, router := gin.CreateTestContext(w)
	router.PUT("/organizations/:id", handler.UpdateOrganization)

	body, _ := json.Marshal(updateReq)
	req := httptest.NewRequest(http.MethodPut, "/organizations/1", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response domain.Organization
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, int64(1), response.ID)
	assert.Equal(t, "Updated Org", response.Name)

	mockUsecase.AssertExpectations(t)
}

func TestOrganizationHandler_UpdateOrganization_InvalidID(t *testing.T) {
	handler, _ := setupOrganizationHandler()

	updateReq := UpdateOrganizationRequest{
		Name:     "Updated Org",
		ParentID: nil,
	}

	w := httptest.NewRecorder()
	_, router := gin.CreateTestContext(w)
	router.PUT("/organizations/:id", handler.UpdateOrganization)

	body, _ := json.Marshal(updateReq)
	req := httptest.NewRequest(http.MethodPut, "/organizations/invalid", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]string
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, "Invalid organization ID", response["error"])
}

func TestOrganizationHandler_UpdateOrganization_InvalidRequestBody(t *testing.T) {
	handler, _ := setupOrganizationHandler()

	w := httptest.NewRecorder()
	_, router := gin.CreateTestContext(w)
	router.PUT("/organizations/:id", handler.UpdateOrganization)

	// Missing required "name" field
	body := []byte(`{"parent_id": 1}`)
	req := httptest.NewRequest(http.MethodPut, "/organizations/1", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestOrganizationHandler_UpdateOrganization_UsecaseError(t *testing.T) {
	handler, mockUsecase := setupOrganizationHandler()

	updateReq := UpdateOrganizationRequest{
		Name:     "Updated Org",
		ParentID: nil,
	}

	mockUsecase.On("Update", mock.Anything, int64(1), "Updated Org", (*int64)(nil)).Return(nil, errors.New("organization not found"))

	w := httptest.NewRecorder()
	_, router := gin.CreateTestContext(w)
	router.PUT("/organizations/:id", handler.UpdateOrganization)

	body, _ := json.Marshal(updateReq)
	req := httptest.NewRequest(http.MethodPut, "/organizations/1", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)

	var response map[string]string
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, "organization not found", response["error"])

	mockUsecase.AssertExpectations(t)
}

// DeleteOrganization Tests

func TestOrganizationHandler_DeleteOrganization_Success(t *testing.T) {
	handler, mockUsecase := setupOrganizationHandler()

	mockUsecase.On("Delete", mock.Anything, int64(1)).Return(nil)

	w := httptest.NewRecorder()
	_, router := gin.CreateTestContext(w)
	router.DELETE("/organizations/:id", handler.DeleteOrganization)

	req := httptest.NewRequest(http.MethodDelete, "/organizations/1", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNoContent, w.Code)

	mockUsecase.AssertExpectations(t)
}

func TestOrganizationHandler_DeleteOrganization_InvalidID(t *testing.T) {
	handler, _ := setupOrganizationHandler()

	w := httptest.NewRecorder()
	_, router := gin.CreateTestContext(w)
	router.DELETE("/organizations/:id", handler.DeleteOrganization)

	req := httptest.NewRequest(http.MethodDelete, "/organizations/invalid", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]string
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, "Invalid organization ID", response["error"])
}

func TestOrganizationHandler_DeleteOrganization_UsecaseError(t *testing.T) {
	handler, mockUsecase := setupOrganizationHandler()

	mockUsecase.On("Delete", mock.Anything, int64(1)).Return(errors.New("organization has children"))

	w := httptest.NewRecorder()
	_, router := gin.CreateTestContext(w)
	router.DELETE("/organizations/:id", handler.DeleteOrganization)

	req := httptest.NewRequest(http.MethodDelete, "/organizations/1", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)

	var response map[string]string
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, "organization has children", response["error"])

	mockUsecase.AssertExpectations(t)
}
