package domain

import (
	"context"
)

// UserRepository defines the interface for user data access
type UserRepository interface {
	// Create creates a new user
	Create(ctx context.Context, user *User) error

	// GetByID retrieves a user by ID
	GetByID(ctx context.Context, id int64) (*User, error)

	// GetByUsername retrieves a user by username
	GetByUsername(ctx context.Context, username string) (*User, error)

	// GetByEmail retrieves a user by email
	GetByEmail(ctx context.Context, email string) (*User, error)

	// Update updates an existing user
	Update(ctx context.Context, user *User) error

	// UpdateLastLogin updates the last login timestamp
	UpdateLastLogin(ctx context.Context, userID int64) error

	// Delete deletes a user (soft delete by setting is_active to false)
	Delete(ctx context.Context, id int64) error

	// List retrieves all users with optional filtering
	List(ctx context.Context, filter *UserFilter) ([]*User, error)

	// ChangePassword updates a user's password
	ChangePassword(ctx context.Context, userID int64, newPasswordHash string) error
}

// UserFilter represents filtering options for listing users
type UserFilter struct {
	Role     *UserRole
	IsActive *bool
	Search   string // Search in username, email, or full name
}
