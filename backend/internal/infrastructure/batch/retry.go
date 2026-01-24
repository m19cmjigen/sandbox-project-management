package batch

import (
	"context"
	"fmt"
	"log"
	"time"
)

// RetryConfig holds retry configuration
type RetryConfig struct {
	MaxAttempts  int
	InitialDelay time.Duration
	MaxDelay     time.Duration
	Multiplier   float64
}

// DefaultRetryConfig returns default retry configuration
func DefaultRetryConfig() *RetryConfig {
	return &RetryConfig{
		MaxAttempts:  3,
		InitialDelay: 1 * time.Second,
		MaxDelay:     30 * time.Second,
		Multiplier:   2.0,
	}
}

// RetryableFunc is a function that can be retried
type RetryableFunc func(ctx context.Context) error

// WithRetry executes a function with exponential backoff retry
func WithRetry(ctx context.Context, config *RetryConfig, operation string, fn RetryableFunc) error {
	var lastErr error
	delay := config.InitialDelay

	for attempt := 1; attempt <= config.MaxAttempts; attempt++ {
		// Check if context is cancelled
		select {
		case <-ctx.Done():
			return fmt.Errorf("operation cancelled: %w", ctx.Err())
		default:
		}

		// Execute the function
		err := fn(ctx)
		if err == nil {
			if attempt > 1 {
				log.Printf("%s succeeded on attempt %d/%d", operation, attempt, config.MaxAttempts)
			}
			return nil
		}

		lastErr = err
		log.Printf("%s failed (attempt %d/%d): %v", operation, attempt, config.MaxAttempts, err)

		// Don't sleep after the last attempt
		if attempt == config.MaxAttempts {
			break
		}

		// Sleep before retry with exponential backoff
		log.Printf("Retrying %s in %v...", operation, delay)
		select {
		case <-time.After(delay):
			// Continue to next attempt
		case <-ctx.Done():
			return fmt.Errorf("operation cancelled during retry: %w", ctx.Err())
		}

		// Calculate next delay with exponential backoff
		delay = time.Duration(float64(delay) * config.Multiplier)
		if delay > config.MaxDelay {
			delay = config.MaxDelay
		}
	}

	return fmt.Errorf("%s failed after %d attempts: %w", operation, config.MaxAttempts, lastErr)
}

// IsRetryableError determines if an error should be retried
func IsRetryableError(err error) bool {
	if err == nil {
		return false
	}

	// Add logic to determine if error is retryable
	// For example, network errors, temporary failures, rate limits, etc.
	errMsg := err.Error()

	// Network-related errors
	if contains(errMsg, "connection refused") ||
		contains(errMsg, "timeout") ||
		contains(errMsg, "temporary failure") ||
		contains(errMsg, "network unreachable") {
		return true
	}

	// HTTP status codes that are retryable
	if contains(errMsg, "status 429") || // Too Many Requests
		contains(errMsg, "status 503") || // Service Unavailable
		contains(errMsg, "status 504") { // Gateway Timeout
		return true
	}

	return false
}

// contains checks if a string contains a substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) && containsAt(s, substr))
}

func containsAt(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
