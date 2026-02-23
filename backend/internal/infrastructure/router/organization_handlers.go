package router

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

// OrganizationRow represents an organization with aggregated project delay stats.
type OrganizationRow struct {
	ID             int64     `db:"id" json:"id"`
	Name           string    `db:"name" json:"name"`
	ParentID       *int64    `db:"parent_id" json:"parent_id"`
	Path           string    `db:"path" json:"path"`
	Level          int       `db:"level" json:"level"`
	CreatedAt      time.Time `db:"created_at" json:"created_at"`
	UpdatedAt      time.Time `db:"updated_at" json:"updated_at"`
	TotalProjects  int       `db:"total_projects" json:"total_projects"`
	RedProjects    int       `db:"red_projects" json:"red_projects"`
	YellowProjects int       `db:"yellow_projects" json:"yellow_projects"`
	GreenProjects  int       `db:"green_projects" json:"green_projects"`
	DelayStatus    string    `json:"delay_status"`
}

// orgDelayStatus computes the delay status for an organization based on project counts.
func orgDelayStatus(o *OrganizationRow) string {
	switch {
	case o.RedProjects > 0:
		return "RED"
	case o.YellowProjects > 0:
		return "YELLOW"
	default:
		return "GREEN"
	}
}

// orgQuery is the shared SQL for fetching organizations with project delay stats.
const orgQuery = `
	SELECT
		o.id,
		o.name,
		o.parent_id,
		o.path,
		o.level,
		o.created_at,
		o.updated_at,
		COALESCE(COUNT(DISTINCT p.id), 0) AS total_projects,
		COALESCE(COUNT(DISTINCT CASE WHEN COALESCE(pds.red_count, 0) > 0
			THEN p.id END), 0) AS red_projects,
		COALESCE(COUNT(DISTINCT CASE WHEN COALESCE(pds.red_count, 0) = 0
			AND COALESCE(pds.yellow_count, 0) > 0
			THEN p.id END), 0) AS yellow_projects,
		COALESCE(COUNT(DISTINCT CASE WHEN p.id IS NOT NULL
			AND COALESCE(pds.red_count, 0) = 0
			AND COALESCE(pds.yellow_count, 0) = 0
			THEN p.id END), 0) AS green_projects
	FROM organizations o
	LEFT JOIN projects p ON o.id = p.organization_id
	LEFT JOIN (
		SELECT
			project_id,
			COUNT(CASE WHEN delay_status = 'RED'    THEN 1 END) AS red_count,
			COUNT(CASE WHEN delay_status = 'YELLOW' THEN 1 END) AS yellow_count
		FROM issues
		GROUP BY project_id
	) pds ON p.id = pds.project_id
`

// listOrganizationsHandlerWithDB returns a Gin handler for listing all organizations.
func listOrganizationsHandlerWithDB(db *sqlx.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		query := orgQuery + `
			GROUP BY o.id, o.name, o.parent_id, o.path, o.level, o.created_at, o.updated_at
			ORDER BY o.path
		`

		var orgs []OrganizationRow
		if err := db.Select(&orgs, query); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch organizations"})
			return
		}

		for i := range orgs {
			orgs[i].DelayStatus = orgDelayStatus(&orgs[i])
		}

		c.JSON(http.StatusOK, orgs)
	}
}

// getOrganizationHandlerWithDB returns a Gin handler for fetching a single organization by ID.
func getOrganizationHandlerWithDB(db *sqlx.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid organization id"})
			return
		}

		query := orgQuery + `
			WHERE o.id = $1
			GROUP BY o.id, o.name, o.parent_id, o.path, o.level, o.created_at, o.updated_at
		`

		var org OrganizationRow
		if err := db.Get(&org, query, id); err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "organization not found"})
			return
		}

		org.DelayStatus = orgDelayStatus(&org)
		c.JSON(http.StatusOK, org)
	}
}

// getChildOrganizationsHandlerWithDB returns a Gin handler for fetching children of an organization.
func getChildOrganizationsHandlerWithDB(db *sqlx.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid organization id"})
			return
		}

		query := orgQuery + `
			WHERE o.parent_id = $1
			GROUP BY o.id, o.name, o.parent_id, o.path, o.level, o.created_at, o.updated_at
			ORDER BY o.path
		`

		var orgs []OrganizationRow
		if err := db.Select(&orgs, query, id); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch child organizations"})
			return
		}

		for i := range orgs {
			orgs[i].DelayStatus = orgDelayStatus(&orgs[i])
		}

		c.JSON(http.StatusOK, orgs)
	}
}
