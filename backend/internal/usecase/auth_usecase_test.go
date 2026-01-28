package usecase

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"

	"github.com/m19cmjigen/sandbox-project-management/backend/internal/domain"
	"github.com/m19cmjigen/sandbox-project-management/backend/internal/infrastructure/auth"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

// Mock implementations
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) Create(ctx context.Context, user *domain.User) error {
	args := m.Called(ctx, user)
	if args.Get(0) == nil {
		// Simulate ID assignment on create
		user.ID = 123
		user.CreatedAt = time.Now()
		user.UpdatedAt = time.Now()
	}
	return args.Error(0)
}

func (m *MockUserRepository) GetByID(ctx context.Context, id int64) (*domain.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockUserRepository) GetByUsername(ctx context.Context, username string) (*domain.User, error) {
	args := m.Called(ctx, username)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockUserRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockUserRepository) Update(ctx context.Context, user *domain.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepository) Delete(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockUserRepository) List(ctx context.Context, filter *domain.UserFilter) ([]*domain.User, error) {
	args := m.Called(ctx, filter)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.User), args.Error(1)
}

func (m *MockUserRepository) UpdateLastLogin(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockUserRepository) ChangePassword(ctx context.Context, id int64, passwordHash string) error {
	args := m.Called(ctx, id, passwordHash)
	return args.Error(0)
}

// Test helper functions
func createTestUser(role domain.UserRole) *domain.User {
	hash, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	return &domain.User{
		ID:           1,
		Username:     "testuser",
		Email:        "test@example.com",
		PasswordHash: string(hash),
		FullName:     sql.NullString{String: "Test User", Valid: true},
		Role:         role,
		IsActive:     true,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
}

func setupAuthUsecase(t *testing.T) (*authUsecase, *MockUserRepository, *auth.JWTService, *auth.PasswordService) {
	mockRepo := new(MockUserRepository)
	jwtService := auth.NewJWTService(auth.JWTConfig{
		SecretKey:       "test-secret-key-for-testing-only",
		ExpirationHours: 24,
	})
	passwordService := auth.NewPasswordService()

	usecase := &authUsecase{
		userRepo:        mockRepo,
		jwtService:      jwtService,
		passwordService: passwordService,
	}

	return usecase, mockRepo, jwtService, passwordService
}

// Login tests
func TestAuthUsecase_Login_Success(t *testing.T) {
	uc, mockRepo, _, _ := setupAuthUsecase(t)
	ctx := context.Background()

	user := createTestUser(domain.RoleAdmin)
	mockRepo.On("GetByUsername", ctx, "testuser").Return(user, nil)
	mockRepo.On("UpdateLastLogin", ctx, user.ID).Return(nil)

	req := domain.LoginRequest{
		Username: "testuser",
		Password: "password123",
	}

	resp, err := uc.Login(ctx, req)

	require.NoError(t, err)
	assert.NotEmpty(t, resp.Token)
	assert.NotZero(t, resp.ExpiresAt)
	assert.Equal(t, user.ID, resp.User.ID)
	assert.Equal(t, user.Username, resp.User.Username)
	assert.Equal(t, user.Role, resp.User.Role)
	mockRepo.AssertExpectations(t)
}

func TestAuthUsecase_Login_InvalidUsername(t *testing.T) {
	uc, mockRepo, _, _ := setupAuthUsecase(t)
	ctx := context.Background()

	mockRepo.On("GetByUsername", ctx, "nonexistent").Return(nil, nil)

	req := domain.LoginRequest{
		Username: "nonexistent",
		Password: "password123",
	}

	resp, err := uc.Login(ctx, req)

	assert.Nil(t, resp)
	assert.ErrorIs(t, err, ErrInvalidCredentials)
	mockRepo.AssertExpectations(t)
}

func TestAuthUsecase_Login_InvalidPassword(t *testing.T) {
	uc, mockRepo, _, _ := setupAuthUsecase(t)
	ctx := context.Background()

	user := createTestUser(domain.RoleAdmin)
	mockRepo.On("GetByUsername", ctx, "testuser").Return(user, nil)

	req := domain.LoginRequest{
		Username: "testuser",
		Password: "wrongpassword",
	}

	resp, err := uc.Login(ctx, req)

	assert.Nil(t, resp)
	assert.ErrorIs(t, err, ErrInvalidCredentials)
	mockRepo.AssertExpectations(t)
}

func TestAuthUsecase_Login_InactiveUser(t *testing.T) {
	uc, mockRepo, _, _ := setupAuthUsecase(t)
	ctx := context.Background()

	user := createTestUser(domain.RoleAdmin)
	user.IsActive = false
	mockRepo.On("GetByUsername", ctx, "testuser").Return(user, nil)

	req := domain.LoginRequest{
		Username: "testuser",
		Password: "password123",
	}

	resp, err := uc.Login(ctx, req)

	assert.Nil(t, resp)
	assert.ErrorIs(t, err, ErrUserNotActive)
	mockRepo.AssertExpectations(t)
}

// ValidateToken tests
func TestAuthUsecase_ValidateToken_Success(t *testing.T) {
	uc, mockRepo, jwtService, _ := setupAuthUsecase(t)
	ctx := context.Background()

	user := createTestUser(domain.RoleAdmin)
	token, _, err := jwtService.GenerateToken(user)
	require.NoError(t, err)

	mockRepo.On("GetByID", ctx, user.ID).Return(user, nil)

	userInfo, err := uc.ValidateToken(ctx, token)

	require.NoError(t, err)
	assert.Equal(t, user.ID, userInfo.ID)
	assert.Equal(t, user.Username, userInfo.Username)
	assert.Equal(t, user.Role, userInfo.Role)
	mockRepo.AssertExpectations(t)
}

func TestAuthUsecase_ValidateToken_InvalidToken(t *testing.T) {
	uc, mockRepo, _, _ := setupAuthUsecase(t)
	ctx := context.Background()

	userInfo, err := uc.ValidateToken(ctx, "invalid-token")

	assert.Nil(t, userInfo)
	assert.Error(t, err)
	mockRepo.AssertNotCalled(t, "GetByID")
}

func TestAuthUsecase_ValidateToken_UserNotFound(t *testing.T) {
	uc, mockRepo, jwtService, _ := setupAuthUsecase(t)
	ctx := context.Background()

	user := createTestUser(domain.RoleAdmin)
	token, _, err := jwtService.GenerateToken(user)
	require.NoError(t, err)

	mockRepo.On("GetByID", ctx, user.ID).Return(nil, nil)

	userInfo, err := uc.ValidateToken(ctx, token)

	assert.Nil(t, userInfo)
	assert.ErrorIs(t, err, ErrUserNotFound)
	mockRepo.AssertExpectations(t)
}

func TestAuthUsecase_ValidateToken_InactiveUser(t *testing.T) {
	uc, mockRepo, jwtService, _ := setupAuthUsecase(t)
	ctx := context.Background()

	user := createTestUser(domain.RoleAdmin)
	token, _, err := jwtService.GenerateToken(user)
	require.NoError(t, err)

	user.IsActive = false
	mockRepo.On("GetByID", ctx, user.ID).Return(user, nil)

	userInfo, err := uc.ValidateToken(ctx, token)

	assert.Nil(t, userInfo)
	assert.ErrorIs(t, err, ErrUserNotActive)
	mockRepo.AssertExpectations(t)
}

// RefreshToken tests
func TestAuthUsecase_RefreshToken_Success(t *testing.T) {
	uc, mockRepo, jwtService, _ := setupAuthUsecase(t)
	ctx := context.Background()

	user := createTestUser(domain.RoleAdmin)
	oldToken, _, err := jwtService.GenerateToken(user)
	require.NoError(t, err)

	mockRepo.On("GetByID", ctx, user.ID).Return(user, nil)

	resp, err := uc.RefreshToken(ctx, oldToken)

	require.NoError(t, err)
	assert.NotEmpty(t, resp.Token)
	assert.NotZero(t, resp.ExpiresAt)
	mockRepo.AssertExpectations(t)
}

func TestAuthUsecase_RefreshToken_InvalidToken(t *testing.T) {
	uc, mockRepo, _, _ := setupAuthUsecase(t)
	ctx := context.Background()

	resp, err := uc.RefreshToken(ctx, "invalid-token")

	assert.Nil(t, resp)
	assert.Error(t, err)
	mockRepo.AssertNotCalled(t, "GetByID")
}

// CreateUser tests
func TestAuthUsecase_CreateUser_Success(t *testing.T) {
	uc, mockRepo, _, _ := setupAuthUsecase(t)
	ctx := context.Background()

	mockRepo.On("GetByUsername", ctx, "newuser").Return(nil, nil)
	mockRepo.On("GetByEmail", ctx, "new@example.com").Return(nil, nil)
	mockRepo.On("Create", ctx, mock.AnythingOfType("*domain.User")).Return(nil)

	req := domain.CreateUserRequest{
		Username: "newuser",
		Email:    "new@example.com",
		Password: "password123",
		FullName: "New User",
		Role:     domain.RoleViewer,
	}

	user, err := uc.CreateUser(ctx, req)

	require.NoError(t, err)
	assert.Equal(t, req.Username, user.Username)
	assert.Equal(t, req.Email, user.Email)
	assert.Equal(t, req.Role, user.Role)
	assert.True(t, user.IsActive)
	assert.NotEmpty(t, user.PasswordHash)
	mockRepo.AssertExpectations(t)
}

func TestAuthUsecase_CreateUser_UsernameExists(t *testing.T) {
	uc, mockRepo, _, _ := setupAuthUsecase(t)
	ctx := context.Background()

	existingUser := createTestUser(domain.RoleViewer)
	mockRepo.On("GetByUsername", ctx, "testuser").Return(existingUser, nil)

	req := domain.CreateUserRequest{
		Username: "testuser",
		Email:    "new@example.com",
		Password: "password123",
		Role:     domain.RoleViewer,
	}

	user, err := uc.CreateUser(ctx, req)

	assert.Nil(t, user)
	assert.ErrorIs(t, err, ErrUsernameExists)
	mockRepo.AssertExpectations(t)
}

func TestAuthUsecase_CreateUser_EmailExists(t *testing.T) {
	uc, mockRepo, _, _ := setupAuthUsecase(t)
	ctx := context.Background()

	existingUser := createTestUser(domain.RoleViewer)
	mockRepo.On("GetByUsername", ctx, "newuser").Return(nil, nil)
	mockRepo.On("GetByEmail", ctx, "test@example.com").Return(existingUser, nil)

	req := domain.CreateUserRequest{
		Username: "newuser",
		Email:    "test@example.com",
		Password: "password123",
		Role:     domain.RoleViewer,
	}

	user, err := uc.CreateUser(ctx, req)

	assert.Nil(t, user)
	assert.ErrorIs(t, err, ErrEmailExists)
	mockRepo.AssertExpectations(t)
}

func TestAuthUsecase_CreateUser_InvalidRole(t *testing.T) {
	uc, mockRepo, _, _ := setupAuthUsecase(t)
	ctx := context.Background()

	req := domain.CreateUserRequest{
		Username: "newuser",
		Email:    "new@example.com",
		Password: "password123",
		Role:     "invalid_role",
	}

	user, err := uc.CreateUser(ctx, req)

	assert.Nil(t, user)
	assert.ErrorIs(t, err, ErrInvalidRole)
	mockRepo.AssertNotCalled(t, "GetByUsername")
	mockRepo.AssertNotCalled(t, "Create")
}

// GetUser tests
func TestAuthUsecase_GetUser_Success(t *testing.T) {
	uc, mockRepo, _, _ := setupAuthUsecase(t)
	ctx := context.Background()

	expectedUser := createTestUser(domain.RoleAdmin)
	mockRepo.On("GetByID", ctx, int64(1)).Return(expectedUser, nil)

	user, err := uc.GetUser(ctx, 1)

	require.NoError(t, err)
	assert.Equal(t, expectedUser.ID, user.ID)
	assert.Equal(t, expectedUser.Username, user.Username)
	mockRepo.AssertExpectations(t)
}

func TestAuthUsecase_GetUser_NotFound(t *testing.T) {
	uc, mockRepo, _, _ := setupAuthUsecase(t)
	ctx := context.Background()

	mockRepo.On("GetByID", ctx, int64(999)).Return(nil, nil)

	user, err := uc.GetUser(ctx, 999)

	assert.Nil(t, user)
	assert.ErrorIs(t, err, ErrUserNotFound)
	mockRepo.AssertExpectations(t)
}

// UpdateUser tests
func TestAuthUsecase_UpdateUser_Success(t *testing.T) {
	uc, mockRepo, _, _ := setupAuthUsecase(t)
	ctx := context.Background()

	existingUser := createTestUser(domain.RoleViewer)
	newEmail := "updated@example.com"
	newFullName := "Updated Name"
	newRole := domain.RoleManager

	mockRepo.On("GetByID", ctx, int64(1)).Return(existingUser, nil)
	mockRepo.On("GetByEmail", ctx, newEmail).Return(nil, nil)
	mockRepo.On("Update", ctx, mock.AnythingOfType("*domain.User")).Return(nil)

	req := domain.UpdateUserRequest{
		Email:    &newEmail,
		FullName: &newFullName,
		Role:     &newRole,
	}

	user, err := uc.UpdateUser(ctx, 1, req)

	require.NoError(t, err)
	assert.Equal(t, newEmail, user.Email)
	assert.Equal(t, newFullName, user.FullName.String)
	assert.Equal(t, newRole, user.Role)
	mockRepo.AssertExpectations(t)
}

func TestAuthUsecase_UpdateUser_EmailConflict(t *testing.T) {
	uc, mockRepo, _, _ := setupAuthUsecase(t)
	ctx := context.Background()

	existingUser := createTestUser(domain.RoleViewer)
	otherUser := createTestUser(domain.RoleViewer)
	otherUser.ID = 2
	otherUser.Email = "other@example.com"
	newEmail := "other@example.com"

	mockRepo.On("GetByID", ctx, int64(1)).Return(existingUser, nil)
	mockRepo.On("GetByEmail", ctx, newEmail).Return(otherUser, nil)

	req := domain.UpdateUserRequest{
		Email: &newEmail,
	}

	user, err := uc.UpdateUser(ctx, 1, req)

	assert.Nil(t, user)
	assert.ErrorIs(t, err, ErrEmailExists)
	mockRepo.AssertExpectations(t)
}

func TestAuthUsecase_UpdateUser_NotFound(t *testing.T) {
	uc, mockRepo, _, _ := setupAuthUsecase(t)
	ctx := context.Background()

	mockRepo.On("GetByID", ctx, int64(999)).Return(nil, nil)

	req := domain.UpdateUserRequest{}
	user, err := uc.UpdateUser(ctx, 999, req)

	assert.Nil(t, user)
	assert.ErrorIs(t, err, ErrUserNotFound)
	mockRepo.AssertExpectations(t)
}

// ChangePassword tests
func TestAuthUsecase_ChangePassword_Success(t *testing.T) {
	uc, mockRepo, _, _ := setupAuthUsecase(t)
	ctx := context.Background()

	user := createTestUser(domain.RoleAdmin)
	mockRepo.On("GetByID", ctx, int64(1)).Return(user, nil)
	mockRepo.On("ChangePassword", ctx, int64(1), mock.AnythingOfType("string")).Return(nil)

	req := domain.ChangePasswordRequest{
		OldPassword: "password123",
		NewPassword: "newpassword456",
	}

	err := uc.ChangePassword(ctx, 1, req)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestAuthUsecase_ChangePassword_WrongOldPassword(t *testing.T) {
	uc, mockRepo, _, _ := setupAuthUsecase(t)
	ctx := context.Background()

	user := createTestUser(domain.RoleAdmin)
	mockRepo.On("GetByID", ctx, int64(1)).Return(user, nil)

	req := domain.ChangePasswordRequest{
		OldPassword: "wrongpassword",
		NewPassword: "newpassword456",
	}

	err := uc.ChangePassword(ctx, 1, req)

	assert.ErrorIs(t, err, ErrInvalidCredentials)
	mockRepo.AssertExpectations(t)
	mockRepo.AssertNotCalled(t, "ChangePassword")
}

func TestAuthUsecase_ChangePassword_UserNotFound(t *testing.T) {
	uc, mockRepo, _, _ := setupAuthUsecase(t)
	ctx := context.Background()

	mockRepo.On("GetByID", ctx, int64(999)).Return(nil, nil)

	req := domain.ChangePasswordRequest{
		OldPassword: "password123",
		NewPassword: "newpassword456",
	}

	err := uc.ChangePassword(ctx, 999, req)

	assert.ErrorIs(t, err, ErrUserNotFound)
	mockRepo.AssertExpectations(t)
}

// ListUsers tests
func TestAuthUsecase_ListUsers_Success(t *testing.T) {
	uc, mockRepo, _, _ := setupAuthUsecase(t)
	ctx := context.Background()

	expectedUsers := []*domain.User{
		createTestUser(domain.RoleAdmin),
		createTestUser(domain.RoleViewer),
	}

	mockRepo.On("List", ctx, (*domain.UserFilter)(nil)).Return(expectedUsers, nil)

	users, err := uc.ListUsers(ctx, nil)

	require.NoError(t, err)
	assert.Len(t, users, 2)
	mockRepo.AssertExpectations(t)
}

func TestAuthUsecase_ListUsers_WithFilter(t *testing.T) {
	uc, mockRepo, _, _ := setupAuthUsecase(t)
	ctx := context.Background()

	filter := &domain.UserFilter{
		Role:     &[]domain.UserRole{domain.RoleAdmin}[0],
		IsActive: &[]bool{true}[0],
	}

	expectedUsers := []*domain.User{
		createTestUser(domain.RoleAdmin),
	}

	mockRepo.On("List", ctx, filter).Return(expectedUsers, nil)

	users, err := uc.ListUsers(ctx, filter)

	require.NoError(t, err)
	assert.Len(t, users, 1)
	assert.Equal(t, domain.RoleAdmin, users[0].Role)
	mockRepo.AssertExpectations(t)
}

func TestAuthUsecase_ListUsers_RepositoryError(t *testing.T) {
	uc, mockRepo, _, _ := setupAuthUsecase(t)
	ctx := context.Background()

	mockRepo.On("List", ctx, (*domain.UserFilter)(nil)).Return(nil, errors.New("database error"))

	users, err := uc.ListUsers(ctx, nil)

	assert.Nil(t, users)
	assert.Error(t, err)
	mockRepo.AssertExpectations(t)
}

// DeleteUser tests
func TestAuthUsecase_DeleteUser_Success(t *testing.T) {
	uc, mockRepo, _, _ := setupAuthUsecase(t)
	ctx := context.Background()

	user := createTestUser(domain.RoleViewer)
	mockRepo.On("GetByID", ctx, int64(1)).Return(user, nil)
	mockRepo.On("Delete", ctx, int64(1)).Return(nil)

	err := uc.DeleteUser(ctx, 1)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestAuthUsecase_DeleteUser_NotFound(t *testing.T) {
	uc, mockRepo, _, _ := setupAuthUsecase(t)
	ctx := context.Background()

	mockRepo.On("GetByID", ctx, int64(999)).Return(nil, nil)

	err := uc.DeleteUser(ctx, 999)

	assert.ErrorIs(t, err, ErrUserNotFound)
	mockRepo.AssertExpectations(t)
	mockRepo.AssertNotCalled(t, "Delete")
}

func TestAuthUsecase_DeleteUser_RepositoryError(t *testing.T) {
	uc, mockRepo, _, _ := setupAuthUsecase(t)
	ctx := context.Background()

	user := createTestUser(domain.RoleViewer)
	mockRepo.On("GetByID", ctx, int64(1)).Return(user, nil)
	mockRepo.On("Delete", ctx, int64(1)).Return(errors.New("database error"))

	err := uc.DeleteUser(ctx, 1)

	assert.Error(t, err)
	mockRepo.AssertExpectations(t)
}
