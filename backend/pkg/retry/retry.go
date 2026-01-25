package retry

import (
	"context"
	"fmt"
	"log"
	"net"
	"strings"
	"syscall"
	"time"
)

// RetryConfig holds retry configuration
type RetryConfig struct {
	MaxAttempts  int
	InitialDelay time.Duration
	MaxDelay     time.Duration
	Multiplier   float64
}

// DefaultRetryConfig returns the default retry configuration
func DefaultRetryConfig() *RetryConfig {
	return &RetryConfig{
		MaxAttempts:  3,
		InitialDelay: 1 * time.Second,
		MaxDelay:     10 * time.Second,
		Multiplier:   2.0,
	}
}

// RetryableFunc is a function that can be retried
type RetryableFunc func(ctx context.Context) error

// WithRetry executes a function with exponential backoff retry logic
func WithRetry(ctx context.Context, config *RetryConfig, operation string, fn RetryableFunc) error {
	var lastErr error
	delay := config.InitialDelay

	for attempt := 1; attempt <= config.MaxAttempts; attempt++ {
		// Execute the function
		err := fn(ctx)
		if err == nil {
			if attempt > 1 {
				log.Printf("%s succeeded on attempt %d", operation, attempt)
			}
			return nil
		}

		lastErr = err

		// Check if we should retry
		if !IsRetryableError(err) {
			log.Printf("%s failed with non-retryable error: %v", operation, err)
			return err
		}

		// If this was the last attempt, don't wait
		if attempt == config.MaxAttempts {
			break
		}

		// Log retry attempt
		log.Printf("%s failed (attempt %d/%d), retrying in %v: %v",
			operation, attempt, config.MaxAttempts, delay, err)

		// Wait before retry
		select {
		case <-time.After(delay):
			// Continue to next attempt
		case <-ctx.Done():
			return fmt.Errorf("%s cancelled: %w", operation, ctx.Err())
		}

		// Calculate next delay with exponential backoff
		delay = time.Duration(float64(delay) * config.Multiplier)
		if delay > config.MaxDelay {
			delay = config.MaxDelay
		}
	}

	return fmt.Errorf("%s failed after %d attempts: %w", operation, config.MaxAttempts, lastErr)
}

// IsRetryableError checks if an error is retryable
func IsRetryableError(err error) bool {
	if err == nil {
		return false
	}

	errStr := strings.ToLower(err.Error())

	// Network errors
	if strings.Contains(errStr, "timeout") ||
		strings.Contains(errStr, "connection refused") ||
		strings.Contains(errStr, "connection reset") ||
		strings.Contains(errStr, "temporary failure") ||
		strings.Contains(errStr, "no such host") ||
		strings.Contains(errStr, "network is unreachable") {
		return true
	}

	// HTTP 5xx errors
	if strings.Contains(errStr, "500") ||
		strings.Contains(errStr, "502") ||
		strings.Contains(errStr, "503") ||
		strings.Contains(errStr, "504") {
		return true
	}

	// Rate limiting
	if strings.Contains(errStr, "429") ||
		strings.Contains(errStr, "rate limit") {
		return true
	}

	// Check for specific error types
	var netErr net.Error
	if errors, ok := err.(interface{ As(interface{}) bool }); ok && errors.As(&netErr) {
		return netErr.Timeout() || netErr.Temporary()
	}

	// Check for syscall errors
	if errors, ok := err.(interface{ Unwrap() error }); ok {
		if unwrappedErr := errors.Unwrap(); unwrappedErr != nil {
			if unwrappedErr == syscall.ECONNREFUSED ||
				unwrappedErr == syscall.ECONNRESET ||
				unwrappedErr == syscall.ETIMEDOUT {
				return true
			}
		}
	}

	return false
}
