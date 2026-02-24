package auth

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func TestGenerateAndValidate_Success(t *testing.T) {
	tm := NewTokenManager("test-secret")

	token, err := tm.GenerateAccessToken(1, "user@example.com", "viewer")
	if err != nil {
		t.Fatalf("expected no error generating token, got: %v", err)
	}
	if token == "" {
		t.Fatal("expected non-empty token")
	}

	claims, err := tm.ValidateAccessToken(token)
	if err != nil {
		t.Fatalf("expected no error validating token, got: %v", err)
	}
	if claims.UserID != 1 {
		t.Errorf("expected UserID=1, got %d", claims.UserID)
	}
	if claims.Email != "user@example.com" {
		t.Errorf("expected email, got %s", claims.Email)
	}
	if claims.Role != "viewer" {
		t.Errorf("expected role=viewer, got %s", claims.Role)
	}
}

func TestValidate_InvalidToken(t *testing.T) {
	tm := NewTokenManager("test-secret")

	_, err := tm.ValidateAccessToken("this.is.notvalid")
	if err != ErrTokenInvalid {
		t.Errorf("expected ErrTokenInvalid, got: %v", err)
	}
}

func TestValidate_WrongSecret(t *testing.T) {
	tm1 := NewTokenManager("secret-a")
	tm2 := NewTokenManager("secret-b")

	token, _ := tm1.GenerateAccessToken(1, "user@example.com", "admin")

	_, err := tm2.ValidateAccessToken(token)
	if err != ErrTokenInvalid {
		t.Errorf("expected ErrTokenInvalid for wrong secret, got: %v", err)
	}
}

func TestValidate_ExpiredToken(t *testing.T) {
	tm := NewTokenManager("test-secret")

	// 過去の有効期限を持つトークンを直接生成する
	past := time.Now().Add(-1 * time.Hour)
	claims := Claims{
		UserID: 1,
		Email:  "user@example.com",
		Role:   "viewer",
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(past.Add(-2 * time.Hour)),
			ExpiresAt: jwt.NewNumericDate(past),
		},
	}
	raw := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	expired, _ := raw.SignedString(tm.secret)

	_, err := tm.ValidateAccessToken(expired)
	if err != ErrTokenExpired {
		t.Errorf("expected ErrTokenExpired, got: %v", err)
	}
}

func TestPassword_HashAndCheck(t *testing.T) {
	hash, err := HashPassword("my-secure-password")
	if err != nil {
		t.Fatalf("expected no error hashing, got: %v", err)
	}
	if hash == "" {
		t.Fatal("expected non-empty hash")
	}

	if !CheckPassword(hash, "my-secure-password") {
		t.Error("expected CheckPassword to return true for correct password")
	}
	if CheckPassword(hash, "wrong-password") {
		t.Error("expected CheckPassword to return false for wrong password")
	}
}
