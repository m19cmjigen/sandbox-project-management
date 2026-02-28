package router

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"

	"github.com/m19cmjigen/sandbox-project-management/backend/pkg/auth"
)

type userListItem struct {
	ID       int64  `db:"id"        json:"id"`
	Email    string `db:"email"     json:"email"`
	Role     string `db:"role"      json:"role"`
	IsActive bool   `db:"is_active" json:"is_active"`
}

type createUserRequest struct {
	Email    string `json:"email"    binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
	Role     string `json:"role"     binding:"required,oneof=admin project_manager viewer"`
}

type updateUserRequest struct {
	Role     *string `json:"role"`
	IsActive *bool   `json:"is_active"`
}

type changePasswordRequest struct {
	NewPassword string `json:"new_password" binding:"required,min=8"`
}

// listUsersHandlerWithDB handles GET /api/v1/users.
// Returns a list of all users (admin only).
func listUsersHandlerWithDB(db *sqlx.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var users []userListItem
		err := db.Select(&users, `SELECT id, email, role, is_active FROM users ORDER BY id`)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch users"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"data": users})
	}
}

// createUserHandlerWithDB handles POST /api/v1/users.
// Creates a new user with the specified email, password and role (admin only).
func createUserHandlerWithDB(db *sqlx.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req createUserRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		hash, err := auth.HashPassword(req.Password)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to process password"})
			return
		}

		var created userListItem
		err = db.QueryRowx(
			`INSERT INTO users (email, password_hash, role)
			 VALUES ($1, $2, $3)
			 RETURNING id, email, role, is_active`,
			req.Email, hash, req.Role,
		).StructScan(&created)
		if err != nil {
			// メールアドレス重複エラーを区別してわかりやすいメッセージを返す
			c.JSON(http.StatusConflict, gin.H{"error": "email already exists"})
			return
		}

		c.JSON(http.StatusCreated, created)
	}
}

// updateUserHandlerWithDB handles PUT /api/v1/users/:id.
// Updates role and/or is_active for the specified user (admin only).
func updateUserHandlerWithDB(db *sqlx.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		targetID, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
			return
		}

		var req updateUserRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if req.Role != nil {
			allowed := map[string]struct{}{"admin": {}, "project_manager": {}, "viewer": {}}
			if _, ok := allowed[*req.Role]; !ok {
				c.JSON(http.StatusBadRequest, gin.H{"error": "invalid role"})
				return
			}
		}

		// ロールと is_active を両方まとめて更新する
		_, err = db.Exec(
			`UPDATE users
			 SET role      = COALESCE($1, role),
			     is_active = COALESCE($2, is_active),
			     updated_at = CURRENT_TIMESTAMP
			 WHERE id = $3`,
			req.Role, req.IsActive, targetID,
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update user"})
			return
		}

		var updated userListItem
		err = db.QueryRowx(
			`SELECT id, email, role, is_active FROM users WHERE id = $1`,
			targetID,
		).StructScan(&updated)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
			return
		}

		c.JSON(http.StatusOK, updated)
	}
}

// changePasswordHandlerWithDB handles PUT /api/v1/users/:id/password.
// Resets the password for the specified user (admin only). Current password is not required.
func changePasswordHandlerWithDB(db *sqlx.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.ParseInt(c.Param("id"), 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
			return
		}
		var req changePasswordRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		hash, err := auth.HashPassword(req.NewPassword)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to hash password"})
			return
		}
		result, err := db.Exec(
			`UPDATE users SET password_hash = $1, updated_at = CURRENT_TIMESTAMP WHERE id = $2`,
			hash, id,
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update password"})
			return
		}
		rows, _ := result.RowsAffected()
		if rows == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "password updated"})
	}
}

// deleteUserHandlerWithDB handles DELETE /api/v1/users/:id.
// Deletes the specified user (admin only). Self-deletion is prohibited.
func deleteUserHandlerWithDB(db *sqlx.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		targetID, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
			return
		}

		// 自分自身の削除を禁止する
		claims := auth.GetClaims(c)
		if claims != nil && claims.UserID == targetID {
			c.JSON(http.StatusForbidden, gin.H{"error": "cannot delete yourself"})
			return
		}

		result, err := db.Exec(`DELETE FROM users WHERE id = $1`, targetID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete user"})
			return
		}

		rowsAffected, _ := result.RowsAffected()
		if rowsAffected == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
			return
		}

		c.JSON(http.StatusNoContent, nil)
	}
}
