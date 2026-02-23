package router

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

// --- Request body types ---

type createOrgRequest struct {
	Name     string `json:"name"      binding:"required"`
	ParentID *int64 `json:"parent_id"`
}

type updateOrgRequest struct {
	Name string `json:"name" binding:"required"`
}

type assignProjectOrgRequest struct {
	OrganizationID *int64 `json:"organization_id"`
}

// createOrganizationHandlerWithDB handles POST /organizations.
func createOrganizationHandlerWithDB(db *sqlx.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req createOrgRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "name is required"})
			return
		}
		req.Name = strings.TrimSpace(req.Name)
		if req.Name == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "name must not be empty"})
			return
		}

		var parentPath string
		var level int

		if req.ParentID != nil {
			// Fetch parent to derive path and level
			var parentOrg struct {
				Path  string `db:"path"`
				Level int    `db:"level"`
			}
			if err := db.QueryRowx(`SELECT path, level FROM organizations WHERE id = $1`, *req.ParentID).
				StructScan(&parentOrg); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "parent organization not found"})
				return
			}
			if parentOrg.Level >= 2 {
				c.JSON(http.StatusBadRequest, gin.H{"error": "maximum hierarchy depth (level 2) exceeded"})
				return
			}
			parentPath = parentOrg.Path
			level = parentOrg.Level + 1
		}

		// Insert the new organization (path computed after insert using the new ID)
		var newID int64
		err := db.QueryRowx(
			`INSERT INTO organizations (name, parent_id, path, level) VALUES ($1, $2, $3, $4) RETURNING id`,
			req.Name, req.ParentID, "/0/", level,
		).Scan(&newID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create organization"})
			return
		}

		// Update path using the real ID
		var newPath string
		if req.ParentID != nil {
			newPath = fmt.Sprintf("%s%d/", parentPath, newID)
		} else {
			newPath = fmt.Sprintf("/%d/", newID)
		}
		if _, err := db.Exec(`UPDATE organizations SET path = $1 WHERE id = $2`, newPath, newID); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update organization path"})
			return
		}

		// Return the created organization
		var org OrganizationRow
		if err := db.QueryRowx(orgQuery+` WHERE o.id = $1 GROUP BY o.id, o.name, o.parent_id, o.path, o.level, o.created_at, o.updated_at`, newID).
			StructScan(&org); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch created organization"})
			return
		}
		org.DelayStatus = orgDelayStatus(&org)
		c.JSON(http.StatusCreated, org)
	}
}

// updateOrganizationHandlerWithDB handles PUT /organizations/:id.
func updateOrganizationHandlerWithDB(db *sqlx.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid organization id"})
			return
		}

		var req updateOrgRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "name is required"})
			return
		}
		req.Name = strings.TrimSpace(req.Name)
		if req.Name == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "name must not be empty"})
			return
		}

		result, err := db.Exec(`UPDATE organizations SET name = $1 WHERE id = $2`, req.Name, id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update organization"})
			return
		}
		rows, _ := result.RowsAffected()
		if rows == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "organization not found"})
			return
		}

		var org OrganizationRow
		if err := db.QueryRowx(orgQuery+` WHERE o.id = $1 GROUP BY o.id, o.name, o.parent_id, o.path, o.level, o.created_at, o.updated_at`, id).
			StructScan(&org); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch updated organization"})
			return
		}
		org.DelayStatus = orgDelayStatus(&org)
		c.JSON(http.StatusOK, org)
	}
}

// deleteOrganizationHandlerWithDB handles DELETE /organizations/:id.
func deleteOrganizationHandlerWithDB(db *sqlx.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid organization id"})
			return
		}

		// Check for child organizations
		var childCount int
		if err := db.QueryRowx(`SELECT COUNT(*) FROM organizations WHERE parent_id = $1`, id).Scan(&childCount); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to check child organizations"})
			return
		}
		if childCount > 0 {
			c.JSON(http.StatusConflict, gin.H{"error": "cannot delete organization with child organizations"})
			return
		}

		// Check for assigned projects
		var projectCount int
		if err := db.QueryRowx(`SELECT COUNT(*) FROM projects WHERE organization_id = $1`, id).Scan(&projectCount); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to check assigned projects"})
			return
		}
		if projectCount > 0 {
			c.JSON(http.StatusConflict, gin.H{"error": fmt.Sprintf("cannot delete organization with %d assigned project(s)", projectCount)})
			return
		}

		result, err := db.Exec(`DELETE FROM organizations WHERE id = $1`, id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete organization"})
			return
		}
		rows, _ := result.RowsAffected()
		if rows == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "organization not found"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "organization deleted"})
	}
}

// assignProjectToOrganizationHandlerWithDB handles PUT /projects/:id/organization.
func assignProjectToOrganizationHandlerWithDB(db *sqlx.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		projectID, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid project id"})
			return
		}

		var req assignProjectOrgRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
			return
		}

		// Validate organization exists if provided
		if req.OrganizationID != nil {
			var exists bool
			if err := db.QueryRowx(`SELECT EXISTS(SELECT 1 FROM organizations WHERE id = $1)`, *req.OrganizationID).Scan(&exists); err != nil || !exists {
				c.JSON(http.StatusBadRequest, gin.H{"error": "organization not found"})
				return
			}
		}

		result, err := db.Exec(`UPDATE projects SET organization_id = $1 WHERE id = $2`, req.OrganizationID, projectID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to assign project"})
			return
		}
		rows, _ := result.RowsAffected()
		if rows == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "project not found"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "project assigned successfully"})
	}
}
