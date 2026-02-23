package router

import (
	"fmt"
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

// ProjectRow represents a project with aggregated issue counts.
type ProjectRow struct {
	ID            int64      `db:"id" json:"id"`
	JiraProjectID string     `db:"jira_project_id" json:"jira_project_id"`
	Key           string     `db:"key" json:"key"`
	Name          string     `db:"name" json:"name"`
	LeadAccountID *string    `db:"lead_account_id" json:"lead_account_id"`
	LeadEmail     *string    `db:"lead_email" json:"lead_email"`
	OrganizationID *int64    `db:"organization_id" json:"organization_id"`
	CreatedAt     time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt     time.Time  `db:"updated_at" json:"updated_at"`
	RedCount      int        `db:"red_count" json:"red_count"`
	YellowCount   int        `db:"yellow_count" json:"yellow_count"`
	GreenCount    int        `db:"green_count" json:"green_count"`
	OpenCount     int        `db:"open_count" json:"open_count"`
	TotalCount    int        `db:"total_count" json:"total_count"`
	DelayStatus   string     `json:"delay_status"`
}

// PaginationMeta holds pagination metadata for list responses.
type PaginationMeta struct {
	Page       int `json:"page"`
	PerPage    int `json:"per_page"`
	Total      int `json:"total"`
	TotalPages int `json:"total_pages"`
}

// ProjectListResponse is the response body for GET /projects.
type ProjectListResponse struct {
	Data       []ProjectRow   `json:"data"`
	Pagination PaginationMeta `json:"pagination"`
}

// listProjectsHandlerWithDB returns a Gin handler for listing projects with DB access.
func listProjectsHandlerWithDB(db *sqlx.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// --- Parse query params ---
		orgIDStr := c.Query("organization_id")
		unassigned := c.Query("unassigned") == "true"
		delayStatusFilter := c.Query("delay_status")
		sortParam := c.DefaultQuery("sort", "name")

		page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
		if err != nil || page < 1 {
			page = 1
		}
		perPage, err := strconv.Atoi(c.DefaultQuery("per_page", "20"))
		if err != nil || perPage < 1 || perPage > 100 {
			perPage = 20
		}

		// --- Build WHERE conditions ---
		var conditions []string
		var args []interface{}
		argIdx := 1

		if unassigned {
			conditions = append(conditions, "p.organization_id IS NULL")
		} else if orgIDStr != "" {
			orgID, err := strconv.ParseInt(orgIDStr, 10, 64)
			if err == nil {
				conditions = append(conditions, fmt.Sprintf("p.organization_id = $%d", argIdx))
				args = append(args, orgID)
				argIdx++
			}
		}

		whereClause := ""
		if len(conditions) > 0 {
			whereClause = "WHERE " + strings.Join(conditions, " AND ")
		}

		// --- ORDER BY clause (using subquery alias columns, no table prefix) ---
		var orderBy string
		switch sortParam {
		case "delay_count":
			orderBy = "red_count DESC, yellow_count DESC, name ASC"
		case "name_desc":
			orderBy = "name DESC"
		default:
			orderBy = "name ASC"
		}

		// --- Base query for aggregated project data ---
		// This subquery computes issue counts per project, then we filter by delay_status if requested.
		subquery := fmt.Sprintf(`
			SELECT
				p.id,
				p.jira_project_id,
				p.key,
				p.name,
				p.lead_account_id,
				p.lead_email,
				p.organization_id,
				p.created_at,
				p.updated_at,
				COALESCE(COUNT(CASE WHEN i.delay_status = 'RED' THEN 1 END), 0)    AS red_count,
				COALESCE(COUNT(CASE WHEN i.delay_status = 'YELLOW' THEN 1 END), 0) AS yellow_count,
				COALESCE(COUNT(CASE WHEN i.delay_status = 'GREEN' THEN 1 END), 0)  AS green_count,
				COALESCE(COUNT(CASE WHEN i.status_category != 'Done' THEN 1 END), 0) AS open_count,
				COALESCE(COUNT(i.id), 0) AS total_count
			FROM projects p
			LEFT JOIN issues i ON p.id = i.project_id
			%s
			GROUP BY p.id, p.jira_project_id, p.key, p.name, p.lead_account_id, p.lead_email,
			         p.organization_id, p.created_at, p.updated_at
		`, whereClause)

		// Wrap in outer query for delay_status filtering
		var outerConditions []string
		var outerArgs []interface{}
		outerArgIdx := argIdx

		switch delayStatusFilter {
		case "RED":
			outerConditions = append(outerConditions, fmt.Sprintf("red_count > 0"))
		case "YELLOW":
			outerConditions = append(outerConditions, fmt.Sprintf("yellow_count > 0 AND red_count = 0"))
		case "GREEN":
			outerConditions = append(outerConditions, fmt.Sprintf("red_count = 0 AND yellow_count = 0"))
		}

		outerWhere := ""
		if len(outerConditions) > 0 {
			outerWhere = "WHERE " + strings.Join(outerConditions, " AND ")
		}
		_ = outerArgs
		_ = outerArgIdx

		// --- COUNT query ---
		countQuery := fmt.Sprintf(`SELECT COUNT(*) FROM (%s) AS sub %s`, subquery, outerWhere)
		var total int
		if err := db.QueryRowx(countQuery, args...).Scan(&total); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to count projects"})
			return
		}

		// --- Main data query ---
		offset := (page - 1) * perPage
		limitIdx := argIdx
		offsetIdx := argIdx + 1

		mainQuery := fmt.Sprintf(`
			SELECT * FROM (%s) AS sub
			%s
			ORDER BY %s
			LIMIT $%d OFFSET $%d
		`, subquery, outerWhere, orderBy, limitIdx, offsetIdx)

		queryArgs := append(args, perPage, offset)

		projects := make([]ProjectRow, 0)
		if err := db.Select(&projects, mainQuery, queryArgs...); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch projects"})
			return
		}

		// Compute delay_status for each project
		for i := range projects {
			switch {
			case projects[i].RedCount > 0:
				projects[i].DelayStatus = "RED"
			case projects[i].YellowCount > 0:
				projects[i].DelayStatus = "YELLOW"
			default:
				projects[i].DelayStatus = "GREEN"
			}
		}

		totalPages := int(math.Ceil(float64(total) / float64(perPage)))
		if totalPages == 0 {
			totalPages = 1
		}

		c.JSON(http.StatusOK, ProjectListResponse{
			Data: projects,
			Pagination: PaginationMeta{
				Page:       page,
				PerPage:    perPage,
				Total:      total,
				TotalPages: totalPages,
			},
		})
	}
}

// getProjectHandlerWithDB returns a Gin handler for fetching a single project by ID.
func getProjectHandlerWithDB(db *sqlx.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid project id"})
			return
		}

		query := `
			SELECT
				p.id,
				p.jira_project_id,
				p.key,
				p.name,
				p.lead_account_id,
				p.lead_email,
				p.organization_id,
				p.created_at,
				p.updated_at,
				COALESCE(COUNT(CASE WHEN i.delay_status = 'RED' THEN 1 END), 0)    AS red_count,
				COALESCE(COUNT(CASE WHEN i.delay_status = 'YELLOW' THEN 1 END), 0) AS yellow_count,
				COALESCE(COUNT(CASE WHEN i.delay_status = 'GREEN' THEN 1 END), 0)  AS green_count,
				COALESCE(COUNT(CASE WHEN i.status_category != 'Done' THEN 1 END), 0) AS open_count,
				COALESCE(COUNT(i.id), 0) AS total_count
			FROM projects p
			LEFT JOIN issues i ON p.id = i.project_id
			WHERE p.id = $1
			GROUP BY p.id, p.jira_project_id, p.key, p.name, p.lead_account_id, p.lead_email,
			         p.organization_id, p.created_at, p.updated_at
		`

		var project ProjectRow
		if err := db.Get(&project, query, id); err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "project not found"})
			return
		}

		switch {
		case project.RedCount > 0:
			project.DelayStatus = "RED"
		case project.YellowCount > 0:
			project.DelayStatus = "YELLOW"
		default:
			project.DelayStatus = "GREEN"
		}

		c.JSON(http.StatusOK, project)
	}
}
