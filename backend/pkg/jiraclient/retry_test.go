package jiraclient

import (
	"net/http"
	"testing"
	"time"
)

func TestRetryableStatus(t *testing.T) {
	retryable := []int{429, 500, 502, 503, 504}
	nonRetryable := []int{200, 201, 400, 401, 403, 404}

	for _, code := range retryable {
		if !retryableStatus(code) {
			t.Errorf("expected status %d to be retryable", code)
		}
	}
	for _, code := range nonRetryable {
		if retryableStatus(code) {
			t.Errorf("expected status %d to be non-retryable", code)
		}
	}
}

func TestRetryAfterDelay_Backoff(t *testing.T) {
	// attempt=0: initialBackoff (1s)
	// attempt=1: 2s
	// attempt=2: 4s
	cases := []struct {
		attempt  int
		expected time.Duration
	}{
		{0, 1 * time.Second},
		{1, 2 * time.Second},
		{2, 4 * time.Second},
	}

	for _, tc := range cases {
		got := retryAfterDelay(nil, tc.attempt)
		if got != tc.expected {
			t.Errorf("attempt=%d: expected %v, got %v", tc.attempt, tc.expected, got)
		}
	}
}

func TestRetryAfterDelay_MaxBackoff(t *testing.T) {
	// 多数のリトライ後は maxBackoff (30s) を超えない
	got := retryAfterDelay(nil, 10)
	if got > maxBackoff {
		t.Errorf("expected delay <= %v, got %v", maxBackoff, got)
	}
}

func TestRetryAfterDelay_RetryAfterHeader(t *testing.T) {
	resp := &http.Response{
		StatusCode: http.StatusTooManyRequests,
		Header:     http.Header{"Retry-After": []string{"5"}},
	}
	got := retryAfterDelay(resp, 0)
	if got != 5*time.Second {
		t.Errorf("expected 5s from Retry-After header, got %v", got)
	}
}

func TestRetryAfterDelay_RetryAfterHeaderCappedAtMax(t *testing.T) {
	// Retry-After が maxBackoff を超える場合は maxBackoff にキャップされる
	resp := &http.Response{
		StatusCode: http.StatusTooManyRequests,
		Header:     http.Header{"Retry-After": []string{"120"}},
	}
	got := retryAfterDelay(resp, 0)
	if got != maxBackoff {
		t.Errorf("expected %v (capped), got %v", maxBackoff, got)
	}
}
