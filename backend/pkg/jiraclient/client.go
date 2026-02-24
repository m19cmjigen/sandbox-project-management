// Package jiraclient provides a Jira Cloud REST API v3 client with
// rate-limit handling and exponential backoff retry logic.
package jiraclient

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// Config holds the configuration for the Jira API client.
type Config struct {
	// BaseURL is the Jira Cloud instance URL, e.g. https://yourcompany.atlassian.net
	BaseURL string
	// Email is the Jira account email used for Basic authentication.
	Email string
	// APIToken is the Jira API token used for Basic authentication.
	APIToken string
}

// Client is a Jira Cloud REST API v3 client.
type Client struct {
	cfg        Config
	httpClient *http.Client
	authHeader string
	// sleepFn はリトライ前の待機に使用する。テストで差し替え可能。
	sleepFn func(time.Duration)
}

// New creates a new Jira API client with the given configuration.
func New(cfg Config) *Client {
	token := base64.StdEncoding.EncodeToString([]byte(cfg.Email + ":" + cfg.APIToken))
	return &Client{
		cfg:        cfg,
		httpClient: &http.Client{Timeout: 30 * time.Second},
		authHeader: "Basic " + token,
		sleepFn:    time.Sleep,
	}
}

// Ping calls GET /rest/api/3/myself and returns nil when the credentials are valid.
func (c *Client) Ping() error {
	var dest map[string]interface{}
	return c.get("/rest/api/3/myself", &dest)
}

// get performs an authenticated GET request to the given path and decodes the
// JSON response body into dest. Retries are applied on transient errors.
func (c *Client) get(path string, dest interface{}) error {
	url := c.cfg.BaseURL + path
	return c.doWithRetry(http.MethodGet, url, nil, dest)
}

// post performs an authenticated POST request with a JSON body and decodes the
// JSON response body into dest. Retries are applied on transient errors.
func (c *Client) post(path string, body interface{}, dest interface{}) error {
	url := c.cfg.BaseURL + path
	return c.doWithRetry(http.MethodPost, url, body, dest)
}

// doWithRetry executes an HTTP request, retrying on transient failures using
// exponential backoff and honouring Retry-After headers on 429 responses.
func (c *Client) doWithRetry(method, url string, body interface{}, dest interface{}) error {
	var lastErr error
	for attempt := 0; attempt <= maxRetries; attempt++ {
		if attempt > 0 {
			// リトライ前に待機する（前回のレスポンスがない場合は純粋なバックオフ）
			c.sleepFn(retryAfterDelay(nil, attempt-1))
		}

		resp, err := c.doOnce(method, url, body)
		if err != nil {
			lastErr = err
			continue
		}

		if retryableStatus(resp.StatusCode) {
			// レート制限・一時的サーバーエラーは待機してリトライ
			delay := retryAfterDelay(resp, attempt)
			resp.Body.Close()
			lastErr = fmt.Errorf("HTTP %d from %s", resp.StatusCode, url)
			if attempt < maxRetries {
				c.sleepFn(delay)
			}
			continue
		}

		if resp.StatusCode < 200 || resp.StatusCode >= 300 {
			body, _ := io.ReadAll(resp.Body)
			resp.Body.Close()
			return fmt.Errorf("jira API error: HTTP %d: %s", resp.StatusCode, string(body))
		}

		defer resp.Body.Close()
		if err := json.NewDecoder(resp.Body).Decode(dest); err != nil {
			return fmt.Errorf("decode response: %w", err)
		}
		return nil
	}
	return fmt.Errorf("jira API request failed after %d attempts: %w", maxRetries+1, lastErr)
}

// doOnce executes a single HTTP request without retry logic.
func (c *Client) doOnce(method, url string, body interface{}) (*http.Response, error) {
	var bodyReader io.Reader
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("marshal request body: %w", err)
		}
		bodyReader = bytes.NewReader(b)
	}

	req, err := http.NewRequest(method, url, bodyReader)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("Authorization", c.authHeader)
	req.Header.Set("Accept", "application/json")
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	return c.httpClient.Do(req)
}
