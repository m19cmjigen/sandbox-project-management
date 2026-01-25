package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/m19cmjigen/sandbox-project-management/backend/internal/domain"
)

type userRepository struct {
	db *sqlx.DB
}

// NewUserRepository creates a new PostgreSQL user repository
func NewUserRepository(db *sqlx.DB) domain.UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(ctx context.Context, user *domain.User) error {
	query := `
		INSERT INTO users (username, email, password_hash, full_name, role, is_active)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, created_at, updated_at
	`

	return r.db.QueryRowContext(
		ctx,
		query,
		user.Username,
		user.Email,
		user.PasswordHash,
		user.FullName,
		user.Role,
		user.IsActive,
	).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)
}

func (r *userRepository) GetByID(ctx context.Context, id int64) (*domain.User, error) {
	var user domain.User
	query := `
		SELECT id, username, email, password_hash, full_name, role, is_active,
		       last_login_at, created_at, updated_at
		FROM users
		WHERE id = $1
	`

	err := r.db.GetContext(ctx, &user, query, id)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get user by ID: %w", err)
	}

	return &user, nil
}

func (r *userRepository) GetByUsername(ctx context.Context, username string) (*domain.User, error) {
	var user domain.User
	query := `
		SELECT id, username, email, password_hash, full_name, role, is_active,
		       last_login_at, created_at, updated_at
		FROM users
		WHERE username = $1
	`

	err := r.db.GetContext(ctx, &user, query, username)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get user by username: %w", err)
	}

	return &user, nil
}

func (r *userRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	var user domain.User
	query := `
		SELECT id, username, email, password_hash, full_name, role, is_active,
		       last_login_at, created_at, updated_at
		FROM users
		WHERE email = $1
	`

	err := r.db.GetContext(ctx, &user, query, email)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}

	return &user, nil
}

func (r *userRepository) Update(ctx context.Context, user *domain.User) error {
	query := `
		UPDATE users
		SET email = $1, full_name = $2, role = $3, is_active = $4, updated_at = CURRENT_TIMESTAMP
		WHERE id = $5
		RETURNING updated_at
	`

	return r.db.QueryRowContext(
		ctx,
		query,
		user.Email,
		user.FullName,
		user.Role,
		user.IsActive,
		user.ID,
	).Scan(&user.UpdatedAt)
}

func (r *userRepository) UpdateLastLogin(ctx context.Context, userID int64) error {
	query := `
		UPDATE users
		SET last_login_at = $1
		WHERE id = $2
	`

	_, err := r.db.ExecContext(ctx, query, time.Now(), userID)
	return err
}

func (r *userRepository) Delete(ctx context.Context, id int64) error {
	// Soft delete by setting is_active to false
	query := `
		UPDATE users
		SET is_active = FALSE, updated_at = CURRENT_TIMESTAMP
		WHERE id = $1
	`

	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

func (r *userRepository) List(ctx context.Context, filter *domain.UserFilter) ([]*domain.User, error) {
	query := `
		SELECT id, username, email, password_hash, full_name, role, is_active,
		       last_login_at, created_at, updated_at
		FROM users
		WHERE 1=1
	`

	args := []interface{}{}
	argCount := 1

	if filter != nil {
		if filter.Role != nil {
			query += fmt.Sprintf(" AND role = $%d", argCount)
			args = append(args, *filter.Role)
			argCount++
		}

		if filter.IsActive != nil {
			query += fmt.Sprintf(" AND is_active = $%d", argCount)
			args = append(args, *filter.IsActive)
			argCount++
		}

		if filter.Search != "" {
			searchPattern := "%" + strings.ToLower(filter.Search) + "%"
			query += fmt.Sprintf(" AND (LOWER(username) LIKE $%d OR LOWER(email) LIKE $%d OR LOWER(full_name) LIKE $%d)", argCount, argCount, argCount)
			args = append(args, searchPattern)
			argCount++
		}
	}

	query += " ORDER BY created_at DESC"

	var users []*domain.User
	err := r.db.SelectContext(ctx, &users, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to list users: %w", err)
	}

	return users, nil
}

func (r *userRepository) ChangePassword(ctx context.Context, userID int64, newPasswordHash string) error {
	query := `
		UPDATE users
		SET password_hash = $1, updated_at = CURRENT_TIMESTAMP
		WHERE id = $2
	`

	_, err := r.db.ExecContext(ctx, query, newPasswordHash, userID)
	return err
}
