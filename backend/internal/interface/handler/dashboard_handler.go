package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/m19cmjigen/sandbox-project-management/backend/internal/usecase"
	"github.com/m19cmjigen/sandbox-project-management/backend/pkg/logger"
)

// DashboardHandler はダッシュボードのハンドラー
type DashboardHandler struct {
	usecase usecase.DashboardUsecase
	logger  *logger.Logger
}

// NewDashboardHandler は新しいDashboardHandlerを作成
func NewDashboardHandler(uc usecase.DashboardUsecase, log *logger.Logger) *DashboardHandler {
	return &DashboardHandler{
		usecase: uc,
		logger:  log,
	}
}

// GetDashboardSummary はダッシュボード全体のサマリを取得
func (h *DashboardHandler) GetDashboardSummary(c *gin.Context) {
	summary, err := h.usecase.GetSummary(c.Request.Context())
	if err != nil {
		h.logger.Error("Failed to get dashboard summary", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get dashboard summary"})
		return
	}

	c.JSON(http.StatusOK, summary)
}

// GetOrganizationSummary は組織別のサマリを取得
func (h *DashboardHandler) GetOrganizationSummary(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid organization ID"})
		return
	}

	summary, err := h.usecase.GetOrganizationSummary(c.Request.Context(), id)
	if err != nil {
		h.logger.Error("Failed to get organization summary", zap.Int64("id", id), zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get organization summary"})
		return
	}

	if summary == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Organization not found"})
		return
	}

	c.JSON(http.StatusOK, summary)
}

// GetProjectSummary はプロジェクト別のサマリを取得
func (h *DashboardHandler) GetProjectSummary(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID"})
		return
	}

	summary, err := h.usecase.GetProjectSummary(c.Request.Context(), id)
	if err != nil {
		h.logger.Error("Failed to get project summary", zap.Int64("id", id), zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get project summary"})
		return
	}

	if summary == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Project not found"})
		return
	}

	c.JSON(http.StatusOK, summary)
}
