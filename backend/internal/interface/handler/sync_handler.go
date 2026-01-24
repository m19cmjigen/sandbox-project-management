package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/m19cmjigen/sandbox-project-management/backend/internal/usecase"
)

// SyncHandler は同期APIのハンドラ
type SyncHandler struct {
	syncUsecase usecase.SyncUsecase
}

// NewSyncHandler は新しいSyncHandlerを作成
func NewSyncHandler(syncUsecase usecase.SyncUsecase) *SyncHandler {
	return &SyncHandler{
		syncUsecase: syncUsecase,
	}
}

// TriggerSyncRequest は同期トリガーリクエスト
type TriggerSyncRequest struct {
	OrganizationID int64 `json:"organization_id" binding:"required"`
}

// TriggerSync は手動で同期を実行
func (h *SyncHandler) TriggerSync(c *gin.Context) {
	var req TriggerSyncRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	syncLog, err := h.syncUsecase.SyncAllProjects(c.Request.Context(), req.OrganizationID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":  "Sync completed",
		"sync_log": syncLog,
	})
}

// SyncProjectRequest はプロジェクト同期リクエスト
type SyncProjectRequest struct {
	ProjectID int64 `json:"project_id" binding:"required"`
}

// SyncProject は特定プロジェクトの同期を実行
func (h *SyncHandler) SyncProject(c *gin.Context) {
	projectIDStr := c.Param("id")
	projectID, err := strconv.ParseInt(projectIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID"})
		return
	}

	if err := h.syncUsecase.SyncProjectIssues(c.Request.Context(), projectID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":    "Project sync completed",
		"project_id": projectID,
	})
}
