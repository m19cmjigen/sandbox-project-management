package router

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

// DashboardOrg holds per-organization stats for the dashboard summary.
type DashboardOrg struct {
	ID             int64   `db:"id"              json:"id"`
	Name           string  `db:"name"            json:"name"`
	ParentID       *int64  `db:"parent_id"       json:"parent_id"`
	Level          int     `db:"level"           json:"level"`
	TotalProjects  int     `db:"total_projects"  json:"total_projects"`
	RedProjects    int     `db:"red_projects"    json:"red_projects"`
	YellowProjects int     `db:"yellow_projects" json:"yellow_projects"`
	GreenProjects  int     `db:"green_projects"  json:"green_projects"`
	DelayStatus    string  `db:"delay_status"    json:"delay_status"`
	DelayRate      float64 `json:"delay_rate"`
}

// DashboardSummaryResponse is the response body for GET /dashboard/summary.
type DashboardSummaryResponse struct {
	TotalProjects  int            `json:"total_projects"`
	RedProjects    int            `json:"red_projects"`
	YellowProjects int            `json:"yellow_projects"`
	GreenProjects  int            `json:"green_projects"`
	TotalIssues    int            `json:"total_issues"`
	RedIssues      int            `json:"red_issues"`
	YellowIssues   int            `json:"yellow_issues"`
	GreenIssues    int            `json:"green_issues"`
	Organizations  []DashboardOrg `json:"organizations"`
}

// projectStatusCTE computes delay_status per project from issues.
const projectStatusCTE = `
	WITH project_stats AS (
		SELECT
			p.id AS project_id,
			p.organization_id,
			CASE
				WHEN COUNT(i.id) FILTER (WHERE i.delay_status = 'RED')    > 0 THEN 'RED'
				WHEN COUNT(i.id) FILTER (WHERE i.delay_status = 'YELLOW') > 0 THEN 'YELLOW'
				ELSE 'GREEN'
			END AS delay_status
		FROM projects p
		LEFT JOIN issues i ON i.project_id = p.id
		GROUP BY p.id, p.organization_id
	)
`

// getDashboardSummaryHandlerWithDB returns the global dashboard summary.
func getDashboardSummaryHandlerWithDB(db *sqlx.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// --- Project counts (computed from issues) ---
		type projectCounts struct {
			Total  int `db:"total"`
			Red    int `db:"red"`
			Yellow int `db:"yellow"`
			Green  int `db:"green"`
		}
		var pc projectCounts
		err := db.QueryRowx(projectStatusCTE + `
			SELECT
				COUNT(*)                                                AS total,
				COUNT(*) FILTER (WHERE delay_status = 'RED')           AS red,
				COUNT(*) FILTER (WHERE delay_status = 'YELLOW')        AS yellow,
				COUNT(*) FILTER (WHERE delay_status = 'GREEN')         AS green
			FROM project_stats
		`).StructScan(&pc)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch project counts"})
			return
		}

		// --- Issue counts ---
		type issueCounts struct {
			Total  int `db:"total"`
			Red    int `db:"red"`
			Yellow int `db:"yellow"`
			Green  int `db:"green"`
		}
		var ic issueCounts
		err = db.QueryRowx(`
			SELECT
				COUNT(*)                                                AS total,
				COUNT(*) FILTER (WHERE delay_status = 'RED')           AS red,
				COUNT(*) FILTER (WHERE delay_status = 'YELLOW')        AS yellow,
				COUNT(*) FILTER (WHERE delay_status = 'GREEN')         AS green
			FROM issues
		`).StructScan(&ic)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch issue counts"})
			return
		}

		// --- Per-organization stats (computed from issues via project_stats) ---
		orgQuery := projectStatusCTE + `
			SELECT
				o.id,
				o.name,
				o.parent_id,
				o.level,
				COUNT(DISTINCT ps.project_id)                                            AS total_projects,
				COUNT(DISTINCT ps.project_id) FILTER (WHERE ps.delay_status = 'RED')    AS red_projects,
				COUNT(DISTINCT ps.project_id) FILTER (WHERE ps.delay_status = 'YELLOW') AS yellow_projects,
				COUNT(DISTINCT ps.project_id) FILTER (WHERE ps.delay_status = 'GREEN')  AS green_projects,
				CASE
					WHEN COUNT(DISTINCT ps.project_id) FILTER (WHERE ps.delay_status = 'RED')    > 0 THEN 'RED'
					WHEN COUNT(DISTINCT ps.project_id) FILTER (WHERE ps.delay_status = 'YELLOW') > 0 THEN 'YELLOW'
					ELSE 'GREEN'
				END AS delay_status
			FROM organizations o
			LEFT JOIN project_stats ps ON ps.organization_id = o.id
			GROUP BY o.id, o.name, o.parent_id, o.level
			ORDER BY o.level, o.name
		`
		var orgs []DashboardOrg
		if err := db.Select(&orgs, orgQuery); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch organization stats"})
			return
		}

		// Compute delay_rate for each org
		for i := range orgs {
			if orgs[i].TotalProjects > 0 {
				orgs[i].DelayRate = float64(orgs[i].RedProjects) / float64(orgs[i].TotalProjects)
			}
		}

		c.JSON(http.StatusOK, DashboardSummaryResponse{
			TotalProjects:  pc.Total,
			RedProjects:    pc.Red,
			YellowProjects: pc.Yellow,
			GreenProjects:  pc.Green,
			TotalIssues:    ic.Total,
			RedIssues:      ic.Red,
			YellowIssues:   ic.Yellow,
			GreenIssues:    ic.Green,
			Organizations:  orgs,
		})
	}
}

// OrgSummaryResponse is the response body for GET /dashboard/organizations/:id.
type OrgSummaryResponse struct {
	Organization DashboardOrg `json:"organization"`
	Projects     []ProjectRow `json:"projects"`
}

// getOrganizationSummaryHandlerWithDB returns summary for a specific organization.
func getOrganizationSummaryHandlerWithDB(db *sqlx.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid organization id"})
			return
		}

		// Fetch org stats
		var org DashboardOrg
		orgQuery := projectStatusCTE + `
			SELECT
				o.id,
				o.name,
				o.parent_id,
				o.level,
				COUNT(DISTINCT ps.project_id)                                            AS total_projects,
				COUNT(DISTINCT ps.project_id) FILTER (WHERE ps.delay_status = 'RED')    AS red_projects,
				COUNT(DISTINCT ps.project_id) FILTER (WHERE ps.delay_status = 'YELLOW') AS yellow_projects,
				COUNT(DISTINCT ps.project_id) FILTER (WHERE ps.delay_status = 'GREEN')  AS green_projects,
				CASE
					WHEN COUNT(DISTINCT ps.project_id) FILTER (WHERE ps.delay_status = 'RED')    > 0 THEN 'RED'
					WHEN COUNT(DISTINCT ps.project_id) FILTER (WHERE ps.delay_status = 'YELLOW') > 0 THEN 'YELLOW'
					ELSE 'GREEN'
				END AS delay_status
			FROM organizations o
			LEFT JOIN project_stats ps ON ps.organization_id = o.id
			WHERE o.id = $1
			GROUP BY o.id, o.name, o.parent_id, o.level
		`
		if err := db.Get(&org, orgQuery, id); err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "organization not found"})
			return
		}
		if org.TotalProjects > 0 {
			org.DelayRate = float64(org.RedProjects) / float64(org.TotalProjects)
		}

		// Fetch projects in this org with issue counts (same pattern as listProjectsHandlerWithDB)
		projectQuery := `
			SELECT
				p.id,
				p.jira_project_id,
				p.key,
				p.name,
				p.lead_account_id,
				p.lead_email,
				p.organization_id,
				COALESCE(COUNT(CASE WHEN i.delay_status = 'RED'    THEN 1 END), 0) AS red_count,
				COALESCE(COUNT(CASE WHEN i.delay_status = 'YELLOW' THEN 1 END), 0) AS yellow_count,
				COALESCE(COUNT(CASE WHEN i.delay_status = 'GREEN'  THEN 1 END), 0) AS green_count,
				COALESCE(COUNT(CASE WHEN i.status_category != 'Done' THEN 1 END), 0) AS open_count,
				COALESCE(COUNT(i.id), 0)                                           AS total_count,
				p.created_at,
				p.updated_at
			FROM projects p
			LEFT JOIN issues i ON i.project_id = p.id
			WHERE p.organization_id = $1
			GROUP BY p.id, p.jira_project_id, p.key, p.name, p.lead_account_id, p.lead_email,
			         p.organization_id, p.created_at, p.updated_at
			ORDER BY red_count DESC, yellow_count DESC, p.name ASC
		`
		projects := make([]ProjectRow, 0)
		if err := db.Select(&projects, projectQuery, id); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch projects"})
			return
		}
		// Compute delay_status in Go (DelayStatus has no db tag in ProjectRow)
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

		c.JSON(http.StatusOK, OrgSummaryResponse{
			Organization: org,
			Projects:     projects,
		})
	}
}
