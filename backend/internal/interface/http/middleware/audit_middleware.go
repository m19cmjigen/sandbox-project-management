package middleware

import (
	"bytes"
	"encoding/json"
	"io"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/m19cmjigen/sandbox-project-management/backend/internal/domain"
	"github.com/m19cmjigen/sandbox-project-management/backend/internal/usecase"
)

// responseWriter is a custom response writer that captures the response
type responseWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w *responseWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

// AuditMiddleware creates middleware that logs all API requests for audit purposes
func AuditMiddleware(auditUsecase usecase.AuditUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Skip health check and ready endpoints
		if c.Request.URL.Path == "/health" || c.Request.URL.Path == "/ready" {
			c.Next()
			return
		}

		startTime := time.Now()

		// Read and restore request body
		var requestBody string
		if c.Request.Body != nil && c.Request.Method != "GET" {
			bodyBytes, err := io.ReadAll(c.Request.Body)
			if err == nil {
				requestBody = string(bodyBytes)
				// Restore the body for further processing
				c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

				// Sanitize sensitive data from logs (passwords, tokens, etc.)
				if strings.Contains(c.Request.URL.Path, "/auth/login") ||
					strings.Contains(c.Request.URL.Path, "/password") {
					var bodyMap map[string]interface{}
					if err := json.Unmarshal(bodyBytes, &bodyMap); err == nil {
						if _, ok := bodyMap["password"]; ok {
							bodyMap["password"] = "***REDACTED***"
						}
						if _, ok := bodyMap["old_password"]; ok {
							bodyMap["old_password"] = "***REDACTED***"
						}
						if _, ok := bodyMap["new_password"]; ok {
							bodyMap["new_password"] = "***REDACTED***"
						}
						sanitized, _ := json.Marshal(bodyMap)
						requestBody = string(sanitized)
					}
				}
			}
		}

		// Capture response
		blw := &responseWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
		c.Writer = blw

		// Process request
		c.Next()

		// Calculate duration
		duration := time.Since(startTime)
		durationMs := int(duration.Milliseconds())

		// Get user information if available
		var userID *int64
		var username *string
		if userInterface, exists := c.Get("user"); exists {
			if user, ok := userInterface.(domain.UserInfo); ok {
				userID = &user.ID
				username = &user.Username
			}
		}

		// Determine action and resource type from path and method
		action, resourceType, resourceID := determineActionAndResource(c.Request.Method, c.Request.URL.Path, c)

		// Get IP address
		ipAddress := c.ClientIP()

		// Get user agent
		userAgent := c.Request.UserAgent()

		// Get response body (truncate if too long)
		responseBody := blw.body.String()
		if len(responseBody) > 10000 {
			responseBody = responseBody[:10000] + "...[truncated]"
		}

		// Truncate request body if too long
		if len(requestBody) > 10000 {
			requestBody = requestBody[:10000] + "...[truncated]"
		}

		// Get error message if any
		var errorMessage *string
		if len(c.Errors) > 0 {
			errMsg := c.Errors.String()
			errorMessage = &errMsg
		}

		// Get response status
		responseStatus := c.Writer.Status()

		// Create audit log
		auditReq := domain.AuditLogCreateRequest{
			UserID:         userID,
			Username:       username,
			Action:         action,
			ResourceType:   resourceType,
			ResourceID:     resourceID,
			Method:         c.Request.Method,
			Path:           c.Request.URL.Path,
			IPAddress:      &ipAddress,
			UserAgent:      &userAgent,
			ResponseStatus: &responseStatus,
			DurationMs:     &durationMs,
		}

		// Only log request/response body for non-GET requests
		if c.Request.Method != "GET" && requestBody != "" {
			auditReq.RequestBody = &requestBody
		}

		// Only log response body for errors or important operations
		if c.Writer.Status() >= 400 || c.Request.Method != "GET" {
			auditReq.ResponseBody = &responseBody
		}

		if errorMessage != nil {
			auditReq.ErrorMessage = errorMessage
		}

		// Log asynchronously to avoid blocking the response
		go func() {
			if err := auditUsecase.LogAction(c.Request.Context(), auditReq); err != nil {
				// Log error but don't fail the request
				println("Failed to create audit log:", err.Error())
			}
		}()
	}
}

// determineActionAndResource determines the action and resource type from the request
func determineActionAndResource(method, path string, c *gin.Context) (domain.AuditAction, domain.ResourceType, *string) {
	// Extract resource ID from path if present
	var resourceID *string
	if id := c.Param("id"); id != "" {
		resourceID = &id
	}

	// Determine action based on method
	var action domain.AuditAction
	switch method {
	case "POST":
		if strings.Contains(path, "/login") {
			action = domain.AuditActionLogin
		} else if strings.Contains(path, "/sync") {
			action = domain.AuditActionSync
		} else {
			action = domain.AuditActionCreate
		}
	case "PUT", "PATCH":
		action = domain.AuditActionUpdate
	case "DELETE":
		action = domain.AuditActionDelete
	case "GET":
		action = domain.AuditActionView
	default:
		action = domain.AuditActionView
	}

	// Determine resource type from path
	var resourceType domain.ResourceType
	switch {
	case strings.Contains(path, "/users"):
		resourceType = domain.ResourceTypeUser
	case strings.Contains(path, "/organizations"):
		resourceType = domain.ResourceTypeOrganization
	case strings.Contains(path, "/projects"):
		resourceType = domain.ResourceTypeProject
	case strings.Contains(path, "/issues"):
		resourceType = domain.ResourceTypeIssue
	case strings.Contains(path, "/sync"):
		resourceType = domain.ResourceTypeSyncLog
	case strings.Contains(path, "/auth"):
		resourceType = domain.ResourceTypeAuth
	case strings.Contains(path, "/dashboard"):
		resourceType = domain.ResourceTypeDashboard
	default:
		resourceType = domain.ResourceTypeDashboard // Default
	}

	return action, resourceType, resourceID
}
