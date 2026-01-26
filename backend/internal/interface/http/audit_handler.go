package http

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/m19cmjigen/sandbox-project-management/backend/internal/domain"
	"github.com/m19cmjigen/sandbox-project-management/backend/internal/usecase"
)

// AuditHandler handles audit log HTTP requests
type AuditHandler struct {
	auditUsecase usecase.AuditUsecase
}

// NewAuditHandler creates a new audit handler
func NewAuditHandler(auditUsecase usecase.AuditUsecase) *AuditHandler {
	return &AuditHandler{
		auditUsecase: auditUsecase,
	}
}

// ListAuditLogs retrieves audit logs with filtering
// GET /api/v1/audit/logs
func (h *AuditHandler) ListAuditLogs(c *gin.Context) {
	filter := &domain.AuditLogFilter{
		Limit:  100,
		Offset: 0,
	}

	// Parse query parameters
	if limitStr := c.Query("limit"); limitStr != "" {
		if limit, err := strconv.Atoi(limitStr); err == nil && limit > 0 {
			filter.Limit = limit
		}
	}

	if offsetStr := c.Query("offset"); offsetStr != "" {
		if offset, err := strconv.Atoi(offsetStr); err == nil && offset >= 0 {
			filter.Offset = offset
		}
	}

	if userIDStr := c.Query("user_id"); userIDStr != "" {
		if userID, err := strconv.ParseInt(userIDStr, 10, 64); err == nil {
			filter.UserID = &userID
		}
	}

	if username := c.Query("username"); username != "" {
		filter.Username = &username
	}

	if actionStr := c.Query("action"); actionStr != "" {
		action := domain.AuditAction(actionStr)
		filter.Action = &action
	}

	if resourceTypeStr := c.Query("resource_type"); resourceTypeStr != "" {
		resourceType := domain.ResourceType(resourceTypeStr)
		filter.ResourceType = &resourceType
	}

	if resourceID := c.Query("resource_id"); resourceID != "" {
		filter.ResourceID = &resourceID
	}

	if method := c.Query("method"); method != "" {
		filter.Method = &method
	}

	if startDateStr := c.Query("start_date"); startDateStr != "" {
		if startDate, err := time.Parse(time.RFC3339, startDateStr); err == nil {
			filter.StartDate = &startDate
		}
	}

	if endDateStr := c.Query("end_date"); endDateStr != "" {
		if endDate, err := time.Parse(time.RFC3339, endDateStr); err == nil {
			filter.EndDate = &endDate
		}
	}

	logs, total, err := h.auditUsecase.ListAuditLogs(c.Request.Context(), filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve audit logs"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":   logs,
		"total":  total,
		"limit":  filter.Limit,
		"offset": filter.Offset,
	})
}

// GetAuditLog retrieves a specific audit log by ID
// GET /api/v1/audit/logs/:id
func (h *AuditHandler) GetAuditLog(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid audit log ID"})
		return
	}

	log, err := h.auditUsecase.GetAuditLog(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Audit log not found"})
		return
	}

	c.JSON(http.StatusOK, log)
}

// CleanupOldLogs deletes audit logs older than specified retention period
// DELETE /api/v1/audit/logs/cleanup
func (h *AuditHandler) CleanupOldLogs(c *gin.Context) {
	var req struct {
		RetentionDays int `json:"retention_days" binding:"required,min=1"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	deleted, err := h.auditUsecase.CleanupOldLogs(c.Request.Context(), req.RetentionDays)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to cleanup audit logs"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Audit logs cleaned up successfully",
		"deleted": deleted,
	})
}
