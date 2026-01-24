package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/m19cmjigen/sandbox-project-management/backend/internal/usecase"
	"github.com/m19cmjigen/sandbox-project-management/backend/pkg/logger"
)

// OrganizationHandler は組織管理のハンドラー
type OrganizationHandler struct {
	usecase usecase.OrganizationUsecase
	logger  *logger.Logger
}

// NewOrganizationHandler は新しいOrganizationHandlerを作成
func NewOrganizationHandler(uc usecase.OrganizationUsecase, log *logger.Logger) *OrganizationHandler {
	return &OrganizationHandler{
		usecase: uc,
		logger:  log,
	}
}

// CreateOrganizationRequest は組織作成リクエスト
type CreateOrganizationRequest struct {
	Name     string `json:"name" binding:"required"`
	ParentID *int64 `json:"parent_id"`
}

// UpdateOrganizationRequest は組織更新リクエスト
type UpdateOrganizationRequest struct {
	Name     string `json:"name" binding:"required"`
	ParentID *int64 `json:"parent_id"`
}

// ListOrganizations は組織一覧を取得
func (h *OrganizationHandler) ListOrganizations(c *gin.Context) {
	orgs, err := h.usecase.GetAll(c.Request.Context())
	if err != nil {
		h.logger.Error("Failed to get organizations", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get organizations"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"organizations": orgs})
}

// GetOrganization はIDで組織を取得
func (h *OrganizationHandler) GetOrganization(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid organization ID"})
		return
	}

	org, err := h.usecase.GetByID(c.Request.Context(), id)
	if err != nil {
		h.logger.Error("Failed to get organization", zap.Int64("id", id), zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get organization"})
		return
	}

	if org == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Organization not found"})
		return
	}

	c.JSON(http.StatusOK, org)
}

// GetOrganizationChildren は子組織を取得
func (h *OrganizationHandler) GetOrganizationChildren(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid organization ID"})
		return
	}

	children, err := h.usecase.GetChildren(c.Request.Context(), id)
	if err != nil {
		h.logger.Error("Failed to get children", zap.Int64("id", id), zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get children"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"children": children})
}

// GetOrganizationTree は組織ツリーを取得
func (h *OrganizationHandler) GetOrganizationTree(c *gin.Context) {
	tree, err := h.usecase.GetTree(c.Request.Context())
	if err != nil {
		h.logger.Error("Failed to get organization tree", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get organization tree"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"tree": tree})
}

// CreateOrganization は新しい組織を作成
func (h *OrganizationHandler) CreateOrganization(c *gin.Context) {
	var req CreateOrganizationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	org, err := h.usecase.Create(c.Request.Context(), req.Name, req.ParentID)
	if err != nil {
		h.logger.Error("Failed to create organization", zap.String("name", req.Name), zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, org)
}

// UpdateOrganization は組織を更新
func (h *OrganizationHandler) UpdateOrganization(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid organization ID"})
		return
	}

	var req UpdateOrganizationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	org, err := h.usecase.Update(c.Request.Context(), id, req.Name, req.ParentID)
	if err != nil {
		h.logger.Error("Failed to update organization", zap.Int64("id", id), zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, org)
}

// DeleteOrganization は組織を削除
func (h *OrganizationHandler) DeleteOrganization(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid organization ID"})
		return
	}

	err = h.usecase.Delete(c.Request.Context(), id)
	if err != nil {
		h.logger.Error("Failed to delete organization", zap.Int64("id", id), zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}
