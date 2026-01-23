package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Logger はアプリケーションロガー
type Logger struct {
	*zap.Logger
}

// New は新しいLoggerを作成する
func New(level, format string) (*Logger, error) {
	var config zap.Config

	// ログフォーマットの設定
	if format == "json" {
		config = zap.NewProductionConfig()
	} else {
		config = zap.NewDevelopmentConfig()
	}

	// ログレベルの設定
	var zapLevel zapcore.Level
	if err := zapLevel.UnmarshalText([]byte(level)); err != nil {
		zapLevel = zapcore.InfoLevel
	}
	config.Level = zap.NewAtomicLevelAt(zapLevel)

	// タイムスタンプのフォーマット
	config.EncoderConfig.TimeKey = "timestamp"
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	logger, err := config.Build()
	if err != nil {
		return nil, err
	}

	return &Logger{logger}, nil
}

// WithFields は追加のフィールドを持つロガーを返す
func (l *Logger) WithFields(fields ...zap.Field) *Logger {
	return &Logger{l.With(fields...)}
}

// Info はInfoレベルのログを出力
func (l *Logger) Info(msg string, fields ...zap.Field) {
	l.Logger.Info(msg, fields...)
}

// Debug はDebugレベルのログを出力
func (l *Logger) Debug(msg string, fields ...zap.Field) {
	l.Logger.Debug(msg, fields...)
}

// Warn はWarnレベルのログを出力
func (l *Logger) Warn(msg string, fields ...zap.Field) {
	l.Logger.Warn(msg, fields...)
}

// Error はErrorレベルのログを出力
func (l *Logger) Error(msg string, fields ...zap.Field) {
	l.Logger.Error(msg, fields...)
}

// Fatal はFatalレベルのログを出力し、アプリケーションを終了
func (l *Logger) Fatal(msg string, fields ...zap.Field) {
	l.Logger.Fatal(msg, fields...)
}
