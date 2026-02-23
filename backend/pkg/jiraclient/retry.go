package jiraclient

import (
	"net/http"
	"strconv"
	"time"
)

const (
	maxRetries     = 3
	initialBackoff = 1 * time.Second
	maxBackoff     = 30 * time.Second
	backoffFactor  = 2
)

// retryableStatus returns true if the HTTP status code warrants a retry.
func retryableStatus(statusCode int) bool {
	switch statusCode {
	case http.StatusTooManyRequests,
		http.StatusInternalServerError,
		http.StatusBadGateway,
		http.StatusServiceUnavailable,
		http.StatusGatewayTimeout:
		return true
	}
	return false
}

// retryAfterDelay returns the duration to wait before retrying.
// If the response contains a Retry-After header, that value is used.
// Otherwise exponential backoff is applied based on the attempt number (0-indexed).
func retryAfterDelay(resp *http.Response, attempt int) time.Duration {
	// 429レスポンスの Retry-After ヘッダーを優先する
	if resp != nil && resp.StatusCode == http.StatusTooManyRequests {
		if v := resp.Header.Get("Retry-After"); v != "" {
			if secs, err := strconv.ParseFloat(v, 64); err == nil {
				d := time.Duration(secs * float64(time.Second))
				if d > maxBackoff {
					return maxBackoff
				}
				return d
			}
		}
	}

	// Exponential backoff: initialBackoff * backoffFactor^attempt
	delay := initialBackoff
	for i := 0; i < attempt; i++ {
		delay *= backoffFactor
		if delay > maxBackoff {
			return maxBackoff
		}
	}
	return delay
}
