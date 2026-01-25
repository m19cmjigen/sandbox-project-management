package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/m19cmjigen/sandbox-project-management/backend/internal/interface/repository"
	"github.com/m19cmjigen/sandbox-project-management/backend/internal/usecase"
)

// SyncHandler は同期APIのハンドラ
type SyncHandler struct {
	syncUsecase usecase.SyncUsecase
	syncLogRepo repository.SyncLogRepository
}

// NewSyncHandler は新しいSyncHandlerを作成
func NewSyncHandler(syncUsecase usecase.SyncUsecase, syncLogRepo repository.SyncLogRepository) *SyncHandler {
	return &SyncHandler{
		syncUsecase: syncUsecase,
		syncLogRepo: syncLogRepo,
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

// GetSyncLogs は同期ログ一覧を取得
func (h *SyncHandler) GetSyncLogs(c *gin.Context) {
	// デフォルトで最新20件を取得
	limitStr := c.DefaultQuery("limit", "20")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 20
	}
	if limit > 100 {
		limit = 100 // 最大100件まで
	}

	logs, err := h.syncLogRepo.FindAll(c.Request.Context(), limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"logs": logs,
	})
}

// GetSyncLog は同期ログ詳細を取得
func (h *SyncHandler) GetSyncLog(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid log ID"})
		return
	}

	log, err := h.syncLogRepo.FindByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if log == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Sync log not found"})
		return
	}

	c.JSON(http.StatusOK, log)
}

// GetLatestSyncLog は最新の同期ログを取得
func (h *SyncHandler) GetLatestSyncLog(c *gin.Context) {
	log, err := h.syncLogRepo.FindLatest(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if log == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "No sync log found"})
		return
	}

	c.JSON(http.StatusOK, log)
}
