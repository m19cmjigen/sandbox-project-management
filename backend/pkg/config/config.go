package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// Config はアプリケーション設定を保持する
type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Log      LogConfig
	JWT      JWTConfig
}

// ServerConfig はサーバー設定
type ServerConfig struct {
	Port    string
	GinMode string
}

// DatabaseConfig はデータベース設定
type DatabaseConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
	SSLMode  string
}

// LogConfig はログ設定
type LogConfig struct {
	Level  string
	Format string
}

// JWTConfig はJWT認証設定
type JWTConfig struct {
	SecretKey       string
	ExpirationHours int
}

// Load は環境変数から設定を読み込む
func Load() (*Config, error) {
	// .envファイルを読み込み（存在しない場合は無視）
	_ = godotenv.Load()

	dbPort, err := strconv.Atoi(getEnv("DB_PORT", "5432"))
	if err != nil {
		return nil, fmt.Errorf("invalid DB_PORT: %w", err)
	}

	jwtExpirationHours, err := strconv.Atoi(getEnv("JWT_EXPIRATION_HOURS", "24"))
	if err != nil {
		return nil, fmt.Errorf("invalid JWT_EXPIRATION_HOURS: %w", err)
	}

	config := &Config{
		Server: ServerConfig{
			Port:    getEnv("PORT", "8080"),
			GinMode: getEnv("GIN_MODE", "debug"),
		},
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     dbPort,
			User:     getEnv("DB_USER", "admin"),
			Password: getEnv("DB_PASSWORD", "admin123"),
			DBName:   getEnv("DB_NAME", "project_visualization"),
			SSLMode:  getEnv("DB_SSLMODE", "disable"),
		},
		Log: LogConfig{
			Level:  getEnv("LOG_LEVEL", "debug"),
			Format: getEnv("LOG_FORMAT", "json"),
		},
		JWT: JWTConfig{
			SecretKey:       getEnv("JWT_SECRET_KEY", "your-secret-key-change-in-production"),
			ExpirationHours: jwtExpirationHours,
		},
	}

	return config, nil
}

// GetDSN はデータベース接続文字列を返す
func (c *DatabaseConfig) GetDSN() string {
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		c.Host, c.Port, c.User, c.Password, c.DBName, c.SSLMode,
	)
}

// getEnv は環境変数を取得し、存在しない場合はデフォルト値を返す
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
