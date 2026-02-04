package http

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
	"github.com/m19cmjigen/sandbox-project-management/backend/internal/domain"
	"github.com/m19cmjigen/sandbox-project-management/backend/internal/usecase"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// MockAuthUsecase is a mock implementation of AuthUsecase
type MockAuthUsecase struct {
	mock.Mock
}

func (m *MockAuthUsecase) Login(ctx context.Context, req domain.LoginRequest) (*domain.LoginResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.LoginResponse), args.Error(1)
}

func (m *MockAuthUsecase) ValidateToken(ctx context.Context, tokenString string) (*domain.UserInfo, error) {
	args := m.Called(ctx, tokenString)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.UserInfo), args.Error(1)
}

func (m *MockAuthUsecase) RefreshToken(ctx context.Context, tokenString string) (*domain.LoginResponse, error) {
	args := m.Called(ctx, tokenString)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.LoginResponse), args.Error(1)
}

func (m *MockAuthUsecase) CreateUser(ctx context.Context, req domain.CreateUserRequest) (*domain.User, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockAuthUsecase) GetUser(ctx context.Context, userID int64) (*domain.User, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockAuthUsecase) UpdateUser(ctx context.Context, userID int64, req domain.UpdateUserRequest) (*domain.User, error) {
	args := m.Called(ctx, userID, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockAuthUsecase) ChangePassword(ctx context.Context, userID int64, req domain.ChangePasswordRequest) error {
	args := m.Called(ctx, userID, req)
	return args.Error(0)
}

func (m *MockAuthUsecase) ListUsers(ctx context.Context, filter *domain.UserFilter) ([]*domain.User, error) {
	args := m.Called(ctx, filter)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.User), args.Error(1)
}

func (m *MockAuthUsecase) DeleteUser(ctx context.Context, userID int64) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}

// Helper functions
func setupAuthHandler() (*AuthHandler, *MockAuthUsecase) {
	gin.SetMode(gin.TestMode)
	mockUsecase := new(MockAuthUsecase)
	handler := NewAuthHandler(mockUsecase)
	return handler, mockUsecase
}

func createTestUserEntity(id int64, role domain.UserRole) *domain.User {
	return &domain.User{
		ID:       id,
		Username: "testuser",
		Email:    "test@example.com",
		Role:     role,
		IsActive: true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

// Login tests
func TestAuthHandler_Login_Success(t *testing.T) {
	handler, mockUsecase := setupAuthHandler()

	loginReq := domain.LoginRequest{
		Username: "testuser",
		Password: "password123",
	}

	loginResp := &domain.LoginResponse{
		Token:     "test-token",
		ExpiresAt: time.Now().Add(24 * time.Hour),
		User: domain.UserInfo{
			ID:       1,
			Username: "testuser",
			Email:    "test@example.com",
			Role:     domain.RoleAdmin,
		},
	}

	mockUsecase.On("Login", mock.Anything, loginReq).Return(loginResp, nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	body, _ := json.Marshal(loginReq)
	c.Request = httptest.NewRequest(http.MethodPost, "/api/auth/login", bytes.NewBuffer(body))
	c.Request.Header.Set("Content-Type", "application/json")

	handler.Login(c)

	assert.Equal(t, http.StatusOK, w.Code)

	var response domain.LoginResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, "test-token", response.Token)
	assert.Equal(t, "testuser", response.User.Username)
	mockUsecase.AssertExpectations(t)
}

func TestAuthHandler_Login_InvalidCredentials(t *testing.T) {
	handler, mockUsecase := setupAuthHandler()

	loginReq := domain.LoginRequest{
		Username: "testuser",
		Password: "wrongpassword",
	}

	mockUsecase.On("Login", mock.Anything, loginReq).Return(nil, usecase.ErrInvalidCredentials)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	body, _ := json.Marshal(loginReq)
	c.Request = httptest.NewRequest(http.MethodPost, "/api/auth/login", bytes.NewBuffer(body))
	c.Request.Header.Set("Content-Type", "application/json")

	handler.Login(c)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	mockUsecase.AssertExpectations(t)
}

func TestAuthHandler_Login_UserNotActive(t *testing.T) {
	handler, mockUsecase := setupAuthHandler()

	loginReq := domain.LoginRequest{
		Username: "testuser",
		Password: "password123",
	}

	mockUsecase.On("Login", mock.Anything, loginReq).Return(nil, usecase.ErrUserNotActive)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	body, _ := json.Marshal(loginReq)
	c.Request = httptest.NewRequest(http.MethodPost, "/api/auth/login", bytes.NewBuffer(body))
	c.Request.Header.Set("Content-Type", "application/json")

	handler.Login(c)

	assert.Equal(t, http.StatusForbidden, w.Code)
	mockUsecase.AssertExpectations(t)
}

func TestAuthHandler_Login_InvalidRequestBody(t *testing.T) {
	handler, _ := setupAuthHandler()

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Request = httptest.NewRequest(http.MethodPost, "/api/auth/login", bytes.NewBufferString("invalid json"))
	c.Request.Header.Set("Content-Type", "application/json")

	handler.Login(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// RefreshToken tests
func TestAuthHandler_RefreshToken_Success(t *testing.T) {
	handler, mockUsecase := setupAuthHandler()

	refreshResp := &domain.LoginResponse{
		Token:     "new-test-token",
		ExpiresAt: time.Now().Add(24 * time.Hour),
		User: domain.UserInfo{
			ID:       1,
			Username: "testuser",
			Role:     domain.RoleAdmin,
		},
	}

	mockUsecase.On("RefreshToken", mock.Anything, "old-token").Return(refreshResp, nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	body, _ := json.Marshal(map[string]string{"token": "old-token"})
	c.Request = httptest.NewRequest(http.MethodPost, "/api/auth/refresh", bytes.NewBuffer(body))
	c.Request.Header.Set("Content-Type", "application/json")

	handler.RefreshToken(c)

	assert.Equal(t, http.StatusOK, w.Code)

	var response domain.LoginResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, "new-test-token", response.Token)
	mockUsecase.AssertExpectations(t)
}

func TestAuthHandler_RefreshToken_InvalidToken(t *testing.T) {
	handler, mockUsecase := setupAuthHandler()

	mockUsecase.On("RefreshToken", mock.Anything, "invalid-token").Return(nil, errors.New("invalid token"))

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	body, _ := json.Marshal(map[string]string{"token": "invalid-token"})
	c.Request = httptest.NewRequest(http.MethodPost, "/api/auth/refresh", bytes.NewBuffer(body))
	c.Request.Header.Set("Content-Type", "application/json")

	handler.RefreshToken(c)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	mockUsecase.AssertExpectations(t)
}

// GetCurrentUser tests
func TestAuthHandler_GetCurrentUser_Success(t *testing.T) {
	handler, _ := setupAuthHandler()

	userInfo := domain.UserInfo{
		ID:       1,
		Username: "testuser",
		Email:    "test@example.com",
		Role:     domain.RoleAdmin,
	}

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set("user", userInfo)
	c.Request = httptest.NewRequest(http.MethodGet, "/api/auth/me", nil)

	handler.GetCurrentUser(c)

	assert.Equal(t, http.StatusOK, w.Code)

	var response domain.UserInfo
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, "testuser", response.Username)
}

func TestAuthHandler_GetCurrentUser_Unauthorized(t *testing.T) {
	handler, _ := setupAuthHandler()

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/api/auth/me", nil)

	handler.GetCurrentUser(c)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

// CreateUser tests
func TestAuthHandler_CreateUser_Success(t *testing.T) {
	handler, mockUsecase := setupAuthHandler()

	createReq := domain.CreateUserRequest{
		Username: "newuser",
		Email:    "new@example.com",
		Password: "password123",
		Role:     domain.RoleViewer,
	}

	createdUser := createTestUserEntity(2, domain.RoleViewer)
	createdUser.Username = "newuser"
	createdUser.Email = "new@example.com"

	mockUsecase.On("CreateUser", mock.Anything, createReq).Return(createdUser, nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	body, _ := json.Marshal(createReq)
	c.Request = httptest.NewRequest(http.MethodPost, "/api/users", bytes.NewBuffer(body))
	c.Request.Header.Set("Content-Type", "application/json")

	handler.CreateUser(c)

	assert.Equal(t, http.StatusCreated, w.Code)

	var response domain.UserInfo
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, "newuser", response.Username)
	mockUsecase.AssertExpectations(t)
}

func TestAuthHandler_CreateUser_UsernameExists(t *testing.T) {
	handler, mockUsecase := setupAuthHandler()

	createReq := domain.CreateUserRequest{
		Username: "existinguser",
		Email:    "new@example.com",
		Password: "password123",
		Role:     domain.RoleViewer,
	}

	mockUsecase.On("CreateUser", mock.Anything, createReq).Return(nil, usecase.ErrUsernameExists)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	body, _ := json.Marshal(createReq)
	c.Request = httptest.NewRequest(http.MethodPost, "/api/users", bytes.NewBuffer(body))
	c.Request.Header.Set("Content-Type", "application/json")

	handler.CreateUser(c)

	assert.Equal(t, http.StatusConflict, w.Code)
	mockUsecase.AssertExpectations(t)
}

func TestAuthHandler_CreateUser_EmailExists(t *testing.T) {
	handler, mockUsecase := setupAuthHandler()

	createReq := domain.CreateUserRequest{
		Username: "newuser",
		Email:    "existing@example.com",
		Password: "password123",
		Role:     domain.RoleViewer,
	}

	mockUsecase.On("CreateUser", mock.Anything, createReq).Return(nil, usecase.ErrEmailExists)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	body, _ := json.Marshal(createReq)
	c.Request = httptest.NewRequest(http.MethodPost, "/api/users", bytes.NewBuffer(body))
	c.Request.Header.Set("Content-Type", "application/json")

	handler.CreateUser(c)

	assert.Equal(t, http.StatusConflict, w.Code)
	mockUsecase.AssertExpectations(t)
}

func TestAuthHandler_CreateUser_InvalidRole(t *testing.T) {
	handler, mockUsecase := setupAuthHandler()

	createReq := domain.CreateUserRequest{
		Username: "newuser",
		Email:    "new@example.com",
		Password: "password123",
		Role:     "invalid_role",
	}

	mockUsecase.On("CreateUser", mock.Anything, createReq).Return(nil, usecase.ErrInvalidRole)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	body, _ := json.Marshal(createReq)
	c.Request = httptest.NewRequest(http.MethodPost, "/api/users", bytes.NewBuffer(body))
	c.Request.Header.Set("Content-Type", "application/json")

	handler.CreateUser(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	mockUsecase.AssertExpectations(t)
}

// GetUser tests
func TestAuthHandler_GetUser_Success(t *testing.T) {
	handler, mockUsecase := setupAuthHandler()

	user := createTestUserEntity(1, domain.RoleAdmin)
	mockUsecase.On("GetUser", mock.Anything, int64(1)).Return(user, nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{{Key: "id", Value: "1"}}
	c.Request = httptest.NewRequest(http.MethodGet, "/api/users/1", nil)

	handler.GetUser(c)

	assert.Equal(t, http.StatusOK, w.Code)

	var response domain.UserInfo
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, "testuser", response.Username)
	mockUsecase.AssertExpectations(t)
}

func TestAuthHandler_GetUser_NotFound(t *testing.T) {
	handler, mockUsecase := setupAuthHandler()

	mockUsecase.On("GetUser", mock.Anything, int64(999)).Return(nil, usecase.ErrUserNotFound)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{{Key: "id", Value: "999"}}
	c.Request = httptest.NewRequest(http.MethodGet, "/api/users/999", nil)

	handler.GetUser(c)

	assert.Equal(t, http.StatusNotFound, w.Code)
	mockUsecase.AssertExpectations(t)
}

func TestAuthHandler_GetUser_InvalidID(t *testing.T) {
	handler, _ := setupAuthHandler()

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{{Key: "id", Value: "invalid"}}
	c.Request = httptest.NewRequest(http.MethodGet, "/api/users/invalid", nil)

	handler.GetUser(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// UpdateUser tests
func TestAuthHandler_UpdateUser_Success(t *testing.T) {
	handler, mockUsecase := setupAuthHandler()

	updateReq := domain.UpdateUserRequest{
		Email: func() *string { s := "updated@example.com"; return &s }(),
	}

	updatedUser := createTestUserEntity(1, domain.RoleAdmin)
	updatedUser.Email = "updated@example.com"

	mockUsecase.On("UpdateUser", mock.Anything, int64(1), updateReq).Return(updatedUser, nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{{Key: "id", Value: "1"}}

	body, _ := json.Marshal(updateReq)
	c.Request = httptest.NewRequest(http.MethodPut, "/api/users/1", bytes.NewBuffer(body))
	c.Request.Header.Set("Content-Type", "application/json")

	handler.UpdateUser(c)

	assert.Equal(t, http.StatusOK, w.Code)
	mockUsecase.AssertExpectations(t)
}

func TestAuthHandler_UpdateUser_NotFound(t *testing.T) {
	handler, mockUsecase := setupAuthHandler()

	updateReq := domain.UpdateUserRequest{}
	mockUsecase.On("UpdateUser", mock.Anything, int64(999), updateReq).Return(nil, usecase.ErrUserNotFound)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{{Key: "id", Value: "999"}}

	body, _ := json.Marshal(updateReq)
	c.Request = httptest.NewRequest(http.MethodPut, "/api/users/999", bytes.NewBuffer(body))
	c.Request.Header.Set("Content-Type", "application/json")

	handler.UpdateUser(c)

	assert.Equal(t, http.StatusNotFound, w.Code)
	mockUsecase.AssertExpectations(t)
}

// ChangePassword tests
func TestAuthHandler_ChangePassword_Success(t *testing.T) {
	handler, mockUsecase := setupAuthHandler()

	changeReq := domain.ChangePasswordRequest{
		OldPassword: "old123",
		NewPassword: "newpass123",
	}

	mockUsecase.On("ChangePassword", mock.Anything, int64(1), mock.AnythingOfType("domain.ChangePasswordRequest")).Return(nil)

	w := httptest.NewRecorder()
	_, router := gin.CreateTestContext(w)

	// Setup router with middleware to set user context
	router.Use(func(c *gin.Context) {
		c.Set("user", domain.UserInfo{ID: 1, Role: domain.RoleAdmin})
	})
	router.POST("/users/:id/password", handler.ChangePassword)

	body, _ := json.Marshal(changeReq)
	req := httptest.NewRequest(http.MethodPost, "/users/1/password", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockUsecase.AssertExpectations(t)
}

func TestAuthHandler_ChangePassword_Forbidden(t *testing.T) {
	handler, _ := setupAuthHandler()

	changeReq := domain.ChangePasswordRequest{
		OldPassword: "old123",
		NewPassword: "newpass123",
	}

	w := httptest.NewRecorder()
	_, router := gin.CreateTestContext(w)

	router.Use(func(c *gin.Context) {
		c.Set("user", domain.UserInfo{ID: 1, Role: domain.RoleViewer}) // Not admin, different user
	})
	router.POST("/users/:id/password", handler.ChangePassword)

	body, _ := json.Marshal(changeReq)
	req := httptest.NewRequest(http.MethodPost, "/users/2/password", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusForbidden, w.Code)
}

func TestAuthHandler_ChangePassword_InvalidOldPassword(t *testing.T) {
	handler, mockUsecase := setupAuthHandler()

	changeReq := domain.ChangePasswordRequest{
		OldPassword: "wrong",
		NewPassword: "newpass123",
	}

	mockUsecase.On("ChangePassword", mock.Anything, int64(1), mock.AnythingOfType("domain.ChangePasswordRequest")).Return(usecase.ErrInvalidCredentials)

	w := httptest.NewRecorder()
	_, router := gin.CreateTestContext(w)

	router.Use(func(c *gin.Context) {
		c.Set("user", domain.UserInfo{ID: 1, Role: domain.RoleAdmin})
	})
	router.POST("/users/:id/password", handler.ChangePassword)

	body, _ := json.Marshal(changeReq)
	req := httptest.NewRequest(http.MethodPost, "/users/1/password", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	mockUsecase.AssertExpectations(t)
}

// ListUsers tests
func TestAuthHandler_ListUsers_Success(t *testing.T) {
	handler, mockUsecase := setupAuthHandler()

	users := []*domain.User{
		createTestUserEntity(1, domain.RoleAdmin),
		createTestUserEntity(2, domain.RoleViewer),
	}

	mockUsecase.On("ListUsers", mock.Anything, mock.AnythingOfType("*domain.UserFilter")).Return(users, nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/api/users", nil)

	handler.ListUsers(c)

	assert.Equal(t, http.StatusOK, w.Code)

	var response []domain.UserInfo
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Len(t, response, 2)
	mockUsecase.AssertExpectations(t)
}

func TestAuthHandler_ListUsers_WithFilters(t *testing.T) {
	handler, mockUsecase := setupAuthHandler()

	users := []*domain.User{
		createTestUserEntity(1, domain.RoleAdmin),
	}

	mockUsecase.On("ListUsers", mock.Anything, mock.MatchedBy(func(filter *domain.UserFilter) bool {
		return filter.Role != nil && *filter.Role == domain.RoleAdmin &&
			filter.IsActive != nil && *filter.IsActive == true
	})).Return(users, nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/api/users?role=admin&is_active=true", nil)

	handler.ListUsers(c)

	assert.Equal(t, http.StatusOK, w.Code)
	mockUsecase.AssertExpectations(t)
}

// DeleteUser tests
func TestAuthHandler_DeleteUser_Success(t *testing.T) {
	handler, mockUsecase := setupAuthHandler()

	mockUsecase.On("DeleteUser", mock.Anything, int64(2)).Return(nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{{Key: "id", Value: "2"}}
	c.Set("user", domain.UserInfo{ID: 1, Role: domain.RoleAdmin})
	c.Request = httptest.NewRequest(http.MethodDelete, "/api/users/2", nil)

	handler.DeleteUser(c)

	assert.Equal(t, http.StatusOK, w.Code)
	mockUsecase.AssertExpectations(t)
}

func TestAuthHandler_DeleteUser_CannotDeleteSelf(t *testing.T) {
	handler, _ := setupAuthHandler()

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{{Key: "id", Value: "1"}}
	c.Set("user", domain.UserInfo{ID: 1, Role: domain.RoleAdmin})
	c.Request = httptest.NewRequest(http.MethodDelete, "/api/users/1", nil)

	handler.DeleteUser(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestAuthHandler_DeleteUser_NotFound(t *testing.T) {
	handler, mockUsecase := setupAuthHandler()

	mockUsecase.On("DeleteUser", mock.Anything, int64(999)).Return(usecase.ErrUserNotFound)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{{Key: "id", Value: "999"}}
	c.Set("user", domain.UserInfo{ID: 1, Role: domain.RoleAdmin})
	c.Request = httptest.NewRequest(http.MethodDelete, "/api/users/999", nil)

	handler.DeleteUser(c)

	assert.Equal(t, http.StatusNotFound, w.Code)
	mockUsecase.AssertExpectations(t)
}
