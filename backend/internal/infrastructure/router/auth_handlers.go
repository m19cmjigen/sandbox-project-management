package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"

	"github.com/m19cmjigen/sandbox-project-management/backend/pkg/auth"
)

type loginRequest struct {
	Email    string `json:"email"    binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type loginResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
	User        userInfo `json:"user"`
}

type userInfo struct {
	ID    int64  `json:"id"`
	Email string `json:"email"`
	Role  string `json:"role"`
}

// loginHandler handles POST /api/v1/auth/login.
// Validates email/password and returns a JWT access token on success.
func loginHandler(db *sqlx.DB, tm *auth.TokenManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req loginRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "email and password are required"})
			return
		}

		// DB からユーザー取得
		var user struct {
			ID           int64  `db:"id"`
			Email        string `db:"email"`
			PasswordHash string `db:"password_hash"`
			Role         string `db:"role"`
			IsActive     bool   `db:"is_active"`
		}
		err := db.QueryRowx(
			`SELECT id, email, password_hash, role, is_active FROM users WHERE email = $1`,
			req.Email,
		).StructScan(&user)
		if err != nil {
			// タイミング攻撃対策: ユーザーが存在しない場合も同じエラーを返す
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid email or password"})
			return
		}

		if !user.IsActive {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "account is disabled"})
			return
		}

		if !auth.CheckPassword(user.PasswordHash, req.Password) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid email or password"})
			return
		}

		token, err := tm.GenerateAccessToken(user.ID, user.Email, user.Role)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate token"})
			return
		}

		c.JSON(http.StatusOK, loginResponse{
			AccessToken: token,
			TokenType:   "Bearer",
			ExpiresIn:   86400, // 24 hours in seconds
			User:        userInfo{ID: user.ID, Email: user.Email, Role: user.Role},
		})
	}
}

// meHandler handles GET /api/v1/auth/me.
// Returns the currently authenticated user's information.
func meHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		claims := auth.GetClaims(c)
		if claims == nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}
		c.JSON(http.StatusOK, userInfo{
			ID:    claims.UserID,
			Email: claims.Email,
			Role:  claims.Role,
		})
	}
}
