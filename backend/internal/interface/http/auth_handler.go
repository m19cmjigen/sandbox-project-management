package http

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/m19cmjigen/sandbox-project-management/backend/internal/domain"
	"github.com/m19cmjigen/sandbox-project-management/backend/internal/usecase"
)

// AuthHandler handles authentication HTTP requests
type AuthHandler struct {
	authUsecase usecase.AuthUsecase
}

// NewAuthHandler creates a new authentication handler
func NewAuthHandler(authUsecase usecase.AuthUsecase) *AuthHandler {
	return &AuthHandler{
		authUsecase: authUsecase,
	}
}

// Login handles user login requests
// POST /api/auth/login
func (h *AuthHandler) Login(c *gin.Context) {
	var req domain.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	response, err := h.authUsecase.Login(c.Request.Context(), req)
	if err != nil {
		if errors.Is(err, usecase.ErrInvalidCredentials) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
			return
		}
		if errors.Is(err, usecase.ErrUserNotActive) {
			c.JSON(http.StatusForbidden, gin.H{"error": "User account is not active"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to login"})
		return
	}

	c.JSON(http.StatusOK, response)
}

// RefreshToken handles token refresh requests
// POST /api/auth/refresh
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var req struct {
		Token string `json:"token" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	response, err := h.authUsecase.RefreshToken(c.Request.Context(), req.Token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Failed to refresh token"})
		return
	}

	c.JSON(http.StatusOK, response)
}

// GetCurrentUser returns the currently authenticated user
// GET /api/auth/me
func (h *AuthHandler) GetCurrentUser(c *gin.Context) {
	userInfo, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	c.JSON(http.StatusOK, userInfo)
}

// CreateUser handles user creation requests (admin only)
// POST /api/users
func (h *AuthHandler) CreateUser(c *gin.Context) {
	var req domain.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	user, err := h.authUsecase.CreateUser(c.Request.Context(), req)
	if err != nil {
		if errors.Is(err, usecase.ErrUsernameExists) {
			c.JSON(http.StatusConflict, gin.H{"error": "Username already exists"})
			return
		}
		if errors.Is(err, usecase.ErrEmailExists) {
			c.JSON(http.StatusConflict, gin.H{"error": "Email already exists"})
			return
		}
		if errors.Is(err, usecase.ErrInvalidRole) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user role"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	c.JSON(http.StatusCreated, user.ToUserInfo())
}

// GetUser retrieves a user by ID
// GET /api/users/:id
func (h *AuthHandler) GetUser(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	user, err := h.authUsecase.GetUser(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, usecase.ErrUserNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user"})
		return
	}

	c.JSON(http.StatusOK, user.ToUserInfo())
}

// UpdateUser updates an existing user
// PUT /api/users/:id
func (h *AuthHandler) UpdateUser(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	var req domain.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	user, err := h.authUsecase.UpdateUser(c.Request.Context(), id, req)
	if err != nil {
		if errors.Is(err, usecase.ErrUserNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}
		if errors.Is(err, usecase.ErrEmailExists) {
			c.JSON(http.StatusConflict, gin.H{"error": "Email already exists"})
			return
		}
		if errors.Is(err, usecase.ErrInvalidRole) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user role"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
		return
	}

	c.JSON(http.StatusOK, user.ToUserInfo())
}

// ChangePassword handles password change requests
// POST /api/users/:id/password
func (h *AuthHandler) ChangePassword(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// Users can only change their own password unless they're admin
	userInfo, _ := c.Get("user")
	currentUser := userInfo.(domain.UserInfo)
	if currentUser.ID != id && currentUser.Role != domain.RoleAdmin {
		c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden"})
		return
	}

	var req domain.ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	if err := h.authUsecase.ChangePassword(c.Request.Context(), id, req); err != nil {
		if errors.Is(err, usecase.ErrUserNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}
		if errors.Is(err, usecase.ErrInvalidCredentials) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid old password"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to change password"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Password changed successfully"})
}

// ListUsers retrieves all users with optional filtering
// GET /api/users
func (h *AuthHandler) ListUsers(c *gin.Context) {
	filter := &domain.UserFilter{}

	// Parse query parameters
	if roleStr := c.Query("role"); roleStr != "" {
		role := domain.UserRole(roleStr)
		filter.Role = &role
	}

	if activeStr := c.Query("is_active"); activeStr != "" {
		isActive := activeStr == "true"
		filter.IsActive = &isActive
	}

	if search := c.Query("search"); search != "" {
		filter.Search = search
	}

	users, err := h.authUsecase.ListUsers(c.Request.Context(), filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list users"})
		return
	}

	// Convert to UserInfo to avoid exposing password hashes
	userInfos := make([]domain.UserInfo, len(users))
	for i, user := range users {
		userInfos[i] = user.ToUserInfo()
	}

	c.JSON(http.StatusOK, userInfos)
}

// DeleteUser deletes (deactivates) a user
// DELETE /api/users/:id
func (h *AuthHandler) DeleteUser(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// Prevent users from deleting themselves
	userInfo, _ := c.Get("user")
	currentUser := userInfo.(domain.UserInfo)
	if currentUser.ID == id {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot delete your own account"})
		return
	}

	if err := h.authUsecase.DeleteUser(c.Request.Context(), id); err != nil {
		if errors.Is(err, usecase.ErrUserNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}
