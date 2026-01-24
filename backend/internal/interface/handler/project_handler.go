package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/m19cmjigen/sandbox-project-management/backend/internal/usecase"
	"github.com/m19cmjigen/sandbox-project-management/backend/pkg/logger"
)

// ProjectHandler はプロジェクト管理のハンドラー
type ProjectHandler struct {
	usecase usecase.ProjectUsecase
	logger  *logger.Logger
}

// NewProjectHandler は新しいProjectHandlerを作成
func NewProjectHandler(uc usecase.ProjectUsecase, log *logger.Logger) *ProjectHandler {
	return &ProjectHandler{
		usecase: uc,
		logger:  log,
	}
}

// AssignProjectRequest はプロジェクト紐付けリクエスト
type AssignProjectRequest struct {
	OrganizationID *int64 `json:"organization_id"`
}

// ListProjects はプロジェクト一覧を取得
func (h *ProjectHandler) ListProjects(c *gin.Context) {
	// クエリパラメータで統計情報付きか判定
	withStats := c.Query("with_stats") == "true"
	orgIDStr := c.Query("organization_id")
	unassigned := c.Query("unassigned") == "true"

	if unassigned {
		// 未分類プロジェクトを取得
		projects, err := h.usecase.GetUnassigned(c.Request.Context())
		if err != nil {
			h.logger.Error("Failed to get unassigned projects", zap.Error(err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get projects"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"projects": projects})
		return
	}

	if orgIDStr != "" {
		// 組織IDでフィルタ
		orgID, err := strconv.ParseInt(orgIDStr, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid organization ID"})
			return
		}

		projects, err := h.usecase.GetByOrganization(c.Request.Context(), orgID)
		if err != nil {
			h.logger.Error("Failed to get projects by organization", zap.Int64("org_id", orgID), zap.Error(err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get projects"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"projects": projects})
		return
	}

	if withStats {
		// 統計情報付きで取得
		projects, err := h.usecase.GetAllWithStats(c.Request.Context())
		if err != nil {
			h.logger.Error("Failed to get projects with stats", zap.Error(err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get projects"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"projects": projects})
		return
	}

	// 通常の一覧取得
	projects, err := h.usecase.GetAll(c.Request.Context())
	if err != nil {
		h.logger.Error("Failed to get projects", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get projects"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"projects": projects})
}

// GetProject はIDでプロジェクトを取得
func (h *ProjectHandler) GetProject(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID"})
		return
	}

	withStats := c.Query("with_stats") == "true"

	if withStats {
		project, err := h.usecase.GetWithStats(c.Request.Context(), id)
		if err != nil {
			h.logger.Error("Failed to get project with stats", zap.Int64("id", id), zap.Error(err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get project"})
			return
		}
		if project == nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Project not found"})
			return
		}
		c.JSON(http.StatusOK, project)
		return
	}

	project, err := h.usecase.GetByID(c.Request.Context(), id)
	if err != nil {
		h.logger.Error("Failed to get project", zap.Int64("id", id), zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get project"})
		return
	}

	if project == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Project not found"})
		return
	}

	c.JSON(http.StatusOK, project)
}

// AssignProjectToOrganization はプロジェクトを組織に紐付け
func (h *ProjectHandler) AssignProjectToOrganization(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID"})
		return
	}

	var req AssignProjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = h.usecase.AssignToOrganization(c.Request.Context(), id, req.OrganizationID)
	if err != nil {
		h.logger.Error("Failed to assign project to organization",
			zap.Int64("project_id", id),
			zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Project assigned successfully"})
}
