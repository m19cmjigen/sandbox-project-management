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

// IssueRow represents a ticket with its project info.
type IssueRow struct {
	ID               int64      `db:"id" json:"id"`
	JiraIssueID      string     `db:"jira_issue_id" json:"jira_issue_id"`
	JiraIssueKey     string     `db:"jira_issue_key" json:"jira_issue_key"`
	ProjectID        int64      `db:"project_id" json:"project_id"`
	ProjectKey       string     `db:"project_key" json:"project_key"`
	ProjectName      string     `db:"project_name" json:"project_name"`
	Summary          string     `db:"summary" json:"summary"`
	Status           string     `db:"status" json:"status"`
	StatusCategory   string     `db:"status_category" json:"status_category"`
	DueDate          *string    `db:"due_date" json:"due_date"`
	AssigneeName     *string    `db:"assignee_name" json:"assignee_name"`
	AssigneeAccountID *string   `db:"assignee_account_id" json:"assignee_account_id"`
	DelayStatus      string     `db:"delay_status" json:"delay_status"`
	Priority         *string    `db:"priority" json:"priority"`
	IssueType        *string    `db:"issue_type" json:"issue_type"`
	LastUpdatedAt    time.Time  `db:"last_updated_at" json:"last_updated_at"`
	CreatedAt        time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt        time.Time  `db:"updated_at" json:"updated_at"`
}

// IssueListResponse is the response body for GET /issues.
type IssueListResponse struct {
	Data       []IssueRow     `json:"data"`
	Pagination PaginationMeta `json:"pagination"`
}

// listIssuesHandlerWithDB returns a Gin handler for listing issues with filters.
func listIssuesHandlerWithDB(db *sqlx.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// --- Parse query params ---
		projectIDStr := c.Query("project_id")
		delayStatus := c.Query("delay_status")
		noDueDateStr := c.Query("no_due_date")
		statusCategory := c.Query("status_category")
		assigneeName := c.Query("assignee_name")
		sortParam := c.DefaultQuery("sort", "due_date")
		orderParam := c.DefaultQuery("order", "asc")

		page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
		if err != nil || page < 1 {
			page = 1
		}
		perPage, err := strconv.Atoi(c.DefaultQuery("per_page", "25"))
		if err != nil || perPage < 1 || perPage > 100 {
			perPage = 25
		}

		// --- Build WHERE conditions ---
		var conditions []string
		var args []interface{}
		idx := 1

		if projectIDStr != "" {
			pid, err := strconv.ParseInt(projectIDStr, 10, 64)
			if err == nil {
				conditions = append(conditions, fmt.Sprintf("i.project_id = $%d", idx))
				args = append(args, pid)
				idx++
			}
		}
		if delayStatus != "" {
			conditions = append(conditions, fmt.Sprintf("i.delay_status = $%d", idx))
			args = append(args, delayStatus)
			idx++
		}
		if noDueDateStr == "true" {
			conditions = append(conditions, "i.due_date IS NULL")
		}
		if statusCategory != "" {
			conditions = append(conditions, fmt.Sprintf("i.status_category = $%d", idx))
			args = append(args, statusCategory)
			idx++
		}
		if assigneeName != "" {
			conditions = append(conditions, fmt.Sprintf("i.assignee_name ILIKE $%d", idx))
			args = append(args, "%"+assigneeName+"%")
			idx++
		}

		whereClause := ""
		if len(conditions) > 0 {
			whereClause = "WHERE " + strings.Join(conditions, " AND ")
		}

		// --- ORDER BY clause ---
		validSortCols := map[string]string{
			"due_date":        "i.due_date",
			"last_updated_at": "i.last_updated_at",
			"jira_issue_key":  "i.jira_issue_key",
			"delay_status":    "i.delay_status",
		}
		sortCol, ok := validSortCols[sortParam]
		if !ok {
			sortCol = "i.due_date"
		}
		dir := "ASC"
		if strings.ToLower(orderParam) == "desc" {
			dir = "DESC"
		}
		// Nulls last for due_date ascending
		nullsClause := ""
		if sortParam == "due_date" && dir == "ASC" {
			nullsClause = " NULLS LAST"
		}
		orderBy := fmt.Sprintf("%s %s%s", sortCol, dir, nullsClause)

		// --- COUNT query ---
		countQuery := fmt.Sprintf(`
			SELECT COUNT(*)
			FROM issues i
			JOIN projects p ON i.project_id = p.id
			%s
		`, whereClause)

		var total int
		if err := db.QueryRowx(countQuery, args...).Scan(&total); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to count issues"})
			return
		}

		// --- Main data query ---
		offset := (page - 1) * perPage
		mainQuery := fmt.Sprintf(`
			SELECT
				i.id,
				i.jira_issue_id,
				i.jira_issue_key,
				i.project_id,
				p.key  AS project_key,
				p.name AS project_name,
				i.summary,
				i.status,
				i.status_category,
				TO_CHAR(i.due_date, 'YYYY-MM-DD') AS due_date,
				i.assignee_name,
				i.assignee_account_id,
				i.delay_status,
				i.priority,
				i.issue_type,
				i.last_updated_at,
				i.created_at,
				i.updated_at
			FROM issues i
			JOIN projects p ON i.project_id = p.id
			%s
			ORDER BY %s
			LIMIT $%d OFFSET $%d
		`, whereClause, orderBy, idx, idx+1)

		queryArgs := append(args, perPage, offset)

		issues := make([]IssueRow, 0)
		if err := db.Select(&issues, mainQuery, queryArgs...); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch issues"})
			return
		}

		totalPages := int(math.Ceil(float64(total) / float64(perPage)))
		if totalPages == 0 {
			totalPages = 1
		}

		c.JSON(http.StatusOK, IssueListResponse{
			Data: issues,
			Pagination: PaginationMeta{
				Page:       page,
				PerPage:    perPage,
				Total:      total,
				TotalPages: totalPages,
			},
		})
	}
}

// getIssueHandlerWithDB returns a Gin handler for fetching a single issue by ID.
func getIssueHandlerWithDB(db *sqlx.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid issue id"})
			return
		}

		query := `
			SELECT
				i.id,
				i.jira_issue_id,
				i.jira_issue_key,
				i.project_id,
				p.key  AS project_key,
				p.name AS project_name,
				i.summary,
				i.status,
				i.status_category,
				TO_CHAR(i.due_date, 'YYYY-MM-DD') AS due_date,
				i.assignee_name,
				i.assignee_account_id,
				i.delay_status,
				i.priority,
				i.issue_type,
				i.last_updated_at,
				i.created_at,
				i.updated_at
			FROM issues i
			JOIN projects p ON i.project_id = p.id
			WHERE i.id = $1
		`

		var issue IssueRow
		if err := db.Get(&issue, query, id); err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "issue not found"})
			return
		}

		c.JSON(http.StatusOK, issue)
	}
}
