package handler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/m19cmjigen/sandbox-project-management/backend/internal/domain"
	"github.com/m19cmjigen/sandbox-project-management/backend/internal/usecase"
	"github.com/m19cmjigen/sandbox-project-management/backend/pkg/logger"
)

// IssueHandler はチケット管理のハンドラー
type IssueHandler struct {
	usecase usecase.IssueUsecase
	logger  *logger.Logger
}

// NewIssueHandler は新しいIssueHandlerを作成
func NewIssueHandler(uc usecase.IssueUsecase, log *logger.Logger) *IssueHandler {
	return &IssueHandler{
		usecase: uc,
		logger:  log,
	}
}

// ListIssues はチケット一覧を取得
func (h *IssueHandler) ListIssues(c *gin.Context) {
	// フィルタパラメータの構築
	filter := domain.IssueFilter{
		Limit:  100, // デフォルト100件
		Offset: 0,
	}

	// プロジェクトIDフィルタ
	if projectIDStr := c.Query("project_id"); projectIDStr != "" {
		projectID, err := strconv.ParseInt(projectIDStr, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project_id"})
			return
		}
		filter.ProjectID = &projectID
	}

	// 遅延ステータスフィルタ
	if delayStatusStr := c.Query("delay_status"); delayStatusStr != "" {
		delayStatus := domain.DelayStatus(delayStatusStr)
		filter.DelayStatus = &delayStatus
	}

	// ステータスフィルタ
	if status := c.Query("status"); status != "" {
		filter.Status = &status
	}

	// 担当者フィルタ
	if assigneeID := c.Query("assignee_id"); assigneeID != "" {
		filter.AssigneeID = &assigneeID
	}

	// 期限フィルタ
	if dueDateFromStr := c.Query("due_date_from"); dueDateFromStr != "" {
		dueDateFrom, err := time.Parse("2006-01-02", dueDateFromStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid due_date_from format. Use YYYY-MM-DD"})
			return
		}
		filter.DueDateFrom = &dueDateFrom
	}

	if dueDateToStr := c.Query("due_date_to"); dueDateToStr != "" {
		dueDateTo, err := time.Parse("2006-01-02", dueDateToStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid due_date_to format. Use YYYY-MM-DD"})
			return
		}
		filter.DueDateTo = &dueDateTo
	}

	// ページネーション
	if limitStr := c.Query("limit"); limitStr != "" {
		limit, err := strconv.Atoi(limitStr)
		if err == nil && limit > 0 {
			filter.Limit = limit
		}
	}

	if offsetStr := c.Query("offset"); offsetStr != "" {
		offset, err := strconv.Atoi(offsetStr)
		if err == nil && offset >= 0 {
			filter.Offset = offset
		}
	}

	// フィルタ条件でチケットを取得
	issues, err := h.usecase.GetByFilter(c.Request.Context(), filter)
	if err != nil {
		h.logger.Error("Failed to get issues", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get issues"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"issues": issues,
		"count":  len(issues),
	})
}

// GetIssue はIDでチケットを取得
func (h *IssueHandler) GetIssue(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid issue ID"})
		return
	}

	issue, err := h.usecase.GetByID(c.Request.Context(), id)
	if err != nil {
		h.logger.Error("Failed to get issue", zap.Int64("id", id), zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get issue"})
		return
	}

	if issue == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Issue not found"})
		return
	}

	c.JSON(http.StatusOK, issue)
}

// ListProjectIssues はプロジェクトのチケット一覧を取得
func (h *IssueHandler) ListProjectIssues(c *gin.Context) {
	projectID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID"})
		return
	}

	issues, err := h.usecase.GetByProjectID(c.Request.Context(), projectID)
	if err != nil {
		h.logger.Error("Failed to get project issues", zap.Int64("project_id", projectID), zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get issues"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"issues": issues,
		"count":  len(issues),
	})
}
