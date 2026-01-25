package usecase

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/m19cmjigen/sandbox-project-management/backend/internal/domain"
	"github.com/m19cmjigen/sandbox-project-management/backend/internal/infrastructure/auth"
)

var (
	ErrInvalidCredentials = errors.New("invalid username or password")
	ErrUserNotFound       = errors.New("user not found")
	ErrUserNotActive      = errors.New("user account is not active")
	ErrUsernameExists     = errors.New("username already exists")
	ErrEmailExists        = errors.New("email already exists")
	ErrInvalidRole        = errors.New("invalid user role")
)

// AuthUsecase defines authentication use cases
type AuthUsecase interface {
	// Login authenticates a user and returns a token
	Login(ctx context.Context, req domain.LoginRequest) (*domain.LoginResponse, error)

	// ValidateToken validates a JWT token and returns user information
	ValidateToken(ctx context.Context, tokenString string) (*domain.UserInfo, error)

	// RefreshToken generates a new token for an existing valid token
	RefreshToken(ctx context.Context, tokenString string) (*domain.LoginResponse, error)

	// CreateUser creates a new user account
	CreateUser(ctx context.Context, req domain.CreateUserRequest) (*domain.User, error)

	// GetUser retrieves a user by ID
	GetUser(ctx context.Context, userID int64) (*domain.User, error)

	// UpdateUser updates an existing user
	UpdateUser(ctx context.Context, userID int64, req domain.UpdateUserRequest) (*domain.User, error)

	// ChangePassword changes a user's password
	ChangePassword(ctx context.Context, userID int64, req domain.ChangePasswordRequest) error

	// ListUsers retrieves all users with optional filtering
	ListUsers(ctx context.Context, filter *domain.UserFilter) ([]*domain.User, error)

	// DeleteUser deletes (deactivates) a user
	DeleteUser(ctx context.Context, userID int64) error
}

type authUsecase struct {
	userRepo        domain.UserRepository
	jwtService      *auth.JWTService
	passwordService *auth.PasswordService
}

// NewAuthUsecase creates a new authentication usecase
func NewAuthUsecase(
	userRepo domain.UserRepository,
	jwtService *auth.JWTService,
	passwordService *auth.PasswordService,
) AuthUsecase {
	return &authUsecase{
		userRepo:        userRepo,
		jwtService:      jwtService,
		passwordService: passwordService,
	}
}

func (u *authUsecase) Login(ctx context.Context, req domain.LoginRequest) (*domain.LoginResponse, error) {
	// Get user by username
	user, err := u.userRepo.GetByUsername(ctx, req.Username)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	if user == nil {
		return nil, ErrInvalidCredentials
	}

	// Check if user is active
	if !user.IsActive {
		return nil, ErrUserNotActive
	}

	// Verify password
	if err := u.passwordService.VerifyPassword(req.Password, user.PasswordHash); err != nil {
		return nil, ErrInvalidCredentials
	}

	// Generate JWT token
	token, expiresAt, err := u.jwtService.GenerateToken(user)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	// Update last login time
	if err := u.userRepo.UpdateLastLogin(ctx, user.ID); err != nil {
		// Log error but don't fail the login
		fmt.Printf("Warning: failed to update last login time: %v\n", err)
	}

	return &domain.LoginResponse{
		Token:     token,
		ExpiresAt: expiresAt,
		User:      user.ToUserInfo(),
	}, nil
}

func (u *authUsecase) ValidateToken(ctx context.Context, tokenString string) (*domain.UserInfo, error) {
	// Validate token
	claims, err := u.jwtService.ValidateToken(tokenString)
	if err != nil {
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	// Get user to ensure they still exist and are active
	user, err := u.userRepo.GetByID(ctx, claims.UserID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	if user == nil {
		return nil, ErrUserNotFound
	}
	if !user.IsActive {
		return nil, ErrUserNotActive
	}

	userInfo := user.ToUserInfo()
	return &userInfo, nil
}

func (u *authUsecase) RefreshToken(ctx context.Context, tokenString string) (*domain.LoginResponse, error) {
	// Validate current token
	claims, err := u.jwtService.ValidateToken(tokenString)
	if err != nil {
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	// Get user
	user, err := u.userRepo.GetByID(ctx, claims.UserID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	if user == nil {
		return nil, ErrUserNotFound
	}
	if !user.IsActive {
		return nil, ErrUserNotActive
	}

	// Generate new token
	newToken, expiresAt, err := u.jwtService.GenerateToken(user)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	return &domain.LoginResponse{
		Token:     newToken,
		ExpiresAt: expiresAt,
		User:      user.ToUserInfo(),
	}, nil
}

func (u *authUsecase) CreateUser(ctx context.Context, req domain.CreateUserRequest) (*domain.User, error) {
	// Validate role
	if !domain.IsValidRole(req.Role) {
		return nil, ErrInvalidRole
	}

	// Check if username already exists
	existingUser, err := u.userRepo.GetByUsername(ctx, req.Username)
	if err != nil {
		return nil, fmt.Errorf("failed to check username: %w", err)
	}
	if existingUser != nil {
		return nil, ErrUsernameExists
	}

	// Check if email already exists
	existingUser, err = u.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to check email: %w", err)
	}
	if existingUser != nil {
		return nil, ErrEmailExists
	}

	// Hash password
	passwordHash, err := u.passwordService.HashPassword(req.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// Create user
	user := &domain.User{
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: passwordHash,
		Role:         req.Role,
		IsActive:     true,
	}

	if req.FullName != "" {
		user.FullName = sql.NullString{String: req.FullName, Valid: true}
	}

	if err := u.userRepo.Create(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return user, nil
}

func (u *authUsecase) GetUser(ctx context.Context, userID int64) (*domain.User, error) {
	user, err := u.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	if user == nil {
		return nil, ErrUserNotFound
	}
	return user, nil
}

func (u *authUsecase) UpdateUser(ctx context.Context, userID int64, req domain.UpdateUserRequest) (*domain.User, error) {
	// Get existing user
	user, err := u.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	if user == nil {
		return nil, ErrUserNotFound
	}

	// Update fields
	if req.Email != nil {
		// Check if email is already in use by another user
		existingUser, err := u.userRepo.GetByEmail(ctx, *req.Email)
		if err != nil {
			return nil, fmt.Errorf("failed to check email: %w", err)
		}
		if existingUser != nil && existingUser.ID != userID {
			return nil, ErrEmailExists
		}
		user.Email = *req.Email
	}

	if req.FullName != nil {
		user.FullName = sql.NullString{String: *req.FullName, Valid: true}
	}

	if req.Role != nil {
		if !domain.IsValidRole(*req.Role) {
			return nil, ErrInvalidRole
		}
		user.Role = *req.Role
	}

	if req.IsActive != nil {
		user.IsActive = *req.IsActive
	}

	// Update user
	if err := u.userRepo.Update(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	return user, nil
}

func (u *authUsecase) ChangePassword(ctx context.Context, userID int64, req domain.ChangePasswordRequest) error {
	// Get user
	user, err := u.userRepo.GetByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}
	if user == nil {
		return ErrUserNotFound
	}

	// Verify old password
	if err := u.passwordService.VerifyPassword(req.OldPassword, user.PasswordHash); err != nil {
		return ErrInvalidCredentials
	}

	// Hash new password
	newPasswordHash, err := u.passwordService.HashPassword(req.NewPassword)
	if err != nil {
		return fmt.Errorf("failed to hash new password: %w", err)
	}

	// Update password
	if err := u.userRepo.ChangePassword(ctx, userID, newPasswordHash); err != nil {
		return fmt.Errorf("failed to change password: %w", err)
	}

	return nil
}

func (u *authUsecase) ListUsers(ctx context.Context, filter *domain.UserFilter) ([]*domain.User, error) {
	users, err := u.userRepo.List(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to list users: %w", err)
	}
	return users, nil
}

func (u *authUsecase) DeleteUser(ctx context.Context, userID int64) error {
	// Get user to verify existence
	user, err := u.userRepo.GetByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}
	if user == nil {
		return ErrUserNotFound
	}

	// Soft delete
	if err := u.userRepo.Delete(ctx, userID); err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	return nil
}
