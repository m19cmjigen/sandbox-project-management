package domain

import (
	"database/sql"
	"time"
)

// UserRole represents the user's role in the system
type UserRole string

const (
	RoleAdmin   UserRole = "admin"
	RoleManager UserRole = "manager"
	RoleViewer  UserRole = "viewer"
)

// User represents a user account in the system
type User struct {
	ID           int64          `db:"id" json:"id"`
	Username     string         `db:"username" json:"username"`
	Email        string         `db:"email" json:"email"`
	PasswordHash string         `db:"password_hash" json:"-"` // Never expose password hash
	FullName     sql.NullString `db:"full_name" json:"full_name,omitempty"`
	Role         UserRole       `db:"role" json:"role"`
	IsActive     bool           `db:"is_active" json:"is_active"`
	LastLoginAt  sql.NullTime   `db:"last_login_at" json:"last_login_at,omitempty"`
	CreatedAt    time.Time      `db:"created_at" json:"created_at"`
	UpdatedAt    time.Time      `db:"updated_at" json:"updated_at"`
}

// LoginRequest represents a login request payload
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// LoginResponse represents a successful login response
type LoginResponse struct {
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expires_at"`
	User      UserInfo  `json:"user"`
}

// UserInfo represents public user information
type UserInfo struct {
	ID       int64    `json:"id"`
	Username string   `json:"username"`
	Email    string   `json:"email"`
	FullName string   `json:"full_name,omitempty"`
	Role     UserRole `json:"role"`
}

// CreateUserRequest represents a request to create a new user
type CreateUserRequest struct {
	Username string   `json:"username" binding:"required"`
	Email    string   `json:"email" binding:"required,email"`
	Password string   `json:"password" binding:"required,min=8"`
	FullName string   `json:"full_name"`
	Role     UserRole `json:"role" binding:"required"`
}

// UpdateUserRequest represents a request to update a user
type UpdateUserRequest struct {
	Email    *string   `json:"email,omitempty" binding:"omitempty,email"`
	FullName *string   `json:"full_name,omitempty"`
	Role     *UserRole `json:"role,omitempty"`
	IsActive *bool     `json:"is_active,omitempty"`
}

// ChangePasswordRequest represents a password change request
type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=8"`
}

// ToUserInfo converts a User to UserInfo (public information only)
func (u *User) ToUserInfo() UserInfo {
	fullName := ""
	if u.FullName.Valid {
		fullName = u.FullName.String
	}

	return UserInfo{
		ID:       u.ID,
		Username: u.Username,
		Email:    u.Email,
		FullName: fullName,
		Role:     u.Role,
	}
}

// IsAdmin checks if the user has admin role
func (u *User) IsAdmin() bool {
	return u.Role == RoleAdmin
}

// IsManager checks if the user has manager or admin role
func (u *User) IsManager() bool {
	return u.Role == RoleManager || u.Role == RoleAdmin
}

// CanManageUsers checks if the user can manage other users
func (u *User) CanManageUsers() bool {
	return u.IsAdmin()
}

// CanManageProjects checks if the user can manage projects
func (u *User) CanManageProjects() bool {
	return u.IsManager()
}

// IsValidRole checks if the given role is valid
func IsValidRole(role UserRole) bool {
	return role == RoleAdmin || role == RoleManager || role == RoleViewer
}
