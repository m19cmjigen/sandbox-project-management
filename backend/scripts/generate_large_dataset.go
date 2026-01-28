package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/m19cmjigen/sandbox-project-management/backend/internal/domain"
	"github.com/m19cmjigen/sandbox-project-management/backend/internal/infrastructure/postgres"
	"github.com/m19cmjigen/sandbox-project-management/backend/internal/interface/repository"
)

const (
	// Target data volumes for performance testing
	TARGET_ORGANIZATIONS = 100
	TARGET_PROJECTS      = 500
	TARGET_ISSUES        = 10000
)

func main() {
	// Get database URL from environment
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgres://admin:admin123@localhost:5432/project_visualization?sslmode=disable"
	}

	// Connect to database
	db, err := sqlx.Connect("postgres", dbURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	log.Println("Connected to database successfully")
	ctx := context.Background()

	// Create repositories
	orgRepo := postgres.NewOrganizationRepository(db)
	projectRepo := postgres.NewProjectRepository(db)
	issueRepo := postgres.NewIssueRepository(db)

	log.Printf("Starting large dataset generation...")
	log.Printf("Target: %d organizations, %d projects, %d issues", TARGET_ORGANIZATIONS, TARGET_PROJECTS, TARGET_ISSUES)

	startTime := time.Now()

	// 1. Generate organizations
	log.Println("Generating organizations...")
	orgs := generateOrganizations(ctx, orgRepo)
	log.Printf("Created %d organizations", len(orgs))

	// 2. Generate projects
	log.Println("Generating projects...")
	projects := generateProjects(ctx, projectRepo, orgs)
	log.Printf("Created %d projects", len(projects))

	// 3. Generate issues
	log.Println("Generating issues (this may take a while)...")
	issues := generateIssues(ctx, issueRepo, projects)
	log.Printf("Created %d issues", len(issues))

	elapsed := time.Since(startTime)
	log.Printf("Large dataset generation completed in %v", elapsed)
	log.Printf("Total records: %d organizations, %d projects, %d issues",
		len(orgs), len(projects), len(issues))
}

func generateOrganizations(ctx context.Context, repo repository.OrganizationRepository) []*domain.Organization {
	orgs := make([]*domain.Organization, 0, TARGET_ORGANIZATIONS)

	// Create root organizations (divisions)
	divisions := []string{
		"事業本部", "営業本部", "技術本部", "管理本部", "マーケティング本部",
		"研究開発本部", "生産本部", "品質保証本部", "人事本部", "財務本部",
	}

	var rootOrgs []*domain.Organization
	for i, divName := range divisions {
		org := &domain.Organization{
			Name:     fmt.Sprintf("%s Division %d", divName, i+1),
			ParentID: nil,
		}
		if err := repo.Create(ctx, org); err != nil {
			log.Printf("Warning: failed to create org %s: %v", org.Name, err)
			continue
		}
		orgs = append(orgs, org)
		rootOrgs = append(rootOrgs, org)
	}

	// Create departments under each division
	for _, parent := range rootOrgs {
		for i := 0; i < 8; i++ {
			dept := &domain.Organization{
				Name:     fmt.Sprintf("%s - 部門 %d", parent.Name, i+1),
				ParentID: &parent.ID,
			}
			if err := repo.Create(ctx, dept); err != nil {
				log.Printf("Warning: failed to create dept: %v", err)
				continue
			}
			orgs = append(orgs, dept)

			// Create teams under some departments
			if i%2 == 0 {
				for j := 0; j < 2; j++ {
					team := &domain.Organization{
						Name:     fmt.Sprintf("%s - チーム %d", dept.Name, j+1),
						ParentID: &dept.ID,
					}
					if err := repo.Create(ctx, team); err != nil {
						continue
					}
					orgs = append(orgs, team)
				}
			}
		}
	}

	return orgs
}

func generateProjects(ctx context.Context, repo repository.ProjectRepository, orgs []*domain.Organization) []*domain.Project {
	projects := make([]*domain.Project, 0, TARGET_PROJECTS)

	projectTypes := []string{
		"システム開発", "業務改善", "インフラ構築", "セキュリティ強化", "データ分析",
		"マーケティング施策", "新規事業", "品質改善", "コスト削減", "研究開発",
	}

	for i := 0; i < TARGET_PROJECTS; i++ {
		// Assign to random organization
		org := orgs[rand.Intn(len(orgs))]
		orgID := org.ID

		projectType := projectTypes[rand.Intn(len(projectTypes))]
		leadEmail := fmt.Sprintf("lead%d@company.com", i%100)

		project := &domain.Project{
			JiraProjectID:  fmt.Sprintf("%d", 20000+i),
			Key:            fmt.Sprintf("PROJ%d", i+1),
			Name:           fmt.Sprintf("%s - プロジェクト%d", projectType, i+1),
			LeadAccountID:  &leadEmail,
			OrganizationID: &orgID,
		}

		if err := repo.Create(ctx, project); err != nil {
			log.Printf("Warning: failed to create project %d: %v", i, err)
			continue
		}

		projects = append(projects, project)

		// Progress indicator
		if (i+1)%50 == 0 {
			log.Printf("  Created %d/%d projects", i+1, TARGET_PROJECTS)
		}
	}

	return projects
}

func generateIssues(ctx context.Context, repo repository.IssueRepository, projects []*domain.Project) []*domain.Issue {
	issues := make([]*domain.Issue, 0, TARGET_ISSUES)
	now := time.Now()

	issueTypes := []string{"Story", "Task", "Bug", "Epic", "Sub-task"}
	priorities := []string{"High", "Medium", "Low", "Critical", "Trivial"}
	statuses := []string{"To Do", "In Progress", "Done", "Blocked", "Review"}
	statusCategories := []domain.StatusCategory{
		domain.StatusCategoryToDo,
		domain.StatusCategoryInProgress,
		domain.StatusCategoryDone,
	}

	batchSize := 100
	batch := make([]*domain.Issue, 0, batchSize)

	for i := 0; i < TARGET_ISSUES; i++ {
		// Assign to random project
		project := projects[rand.Intn(len(projects))]

		// Random delay status distribution: 20% RED, 30% YELLOW, 50% GREEN
		var delayStatus domain.DelayStatus
		var dueDate sql.NullTime
		var status string
		var statusCategory domain.StatusCategory

		randValue := rand.Intn(100)
		switch {
		case randValue < 20: // 20% RED (overdue)
			delayStatus = domain.DelayStatusRed
			dueDate = sql.NullTime{Time: now.Add(-time.Duration(rand.Intn(30)+1) * 24 * time.Hour), Valid: true}
			status = statuses[1] // In Progress
			statusCategory = domain.StatusCategoryInProgress
		case randValue < 50: // 30% YELLOW (due soon or no date)
			delayStatus = domain.DelayStatusYellow
			if rand.Intn(2) == 0 {
				dueDate = sql.NullTime{Time: now.Add(time.Duration(rand.Intn(3)+1) * 24 * time.Hour), Valid: true}
			} else {
				dueDate = sql.NullTime{Valid: false}
			}
			status = statuses[rand.Intn(2)] // To Do or In Progress
			statusCategory = statusCategories[rand.Intn(2)]
		default: // 50% GREEN (on track or done)
			delayStatus = domain.DelayStatusGreen
			if rand.Intn(3) == 0 {
				// Done
				dueDate = sql.NullTime{Time: now.Add(-time.Duration(rand.Intn(7)+1) * 24 * time.Hour), Valid: true}
				status = "Done"
				statusCategory = domain.StatusCategoryDone
			} else {
				// On track
				dueDate = sql.NullTime{Time: now.Add(time.Duration(rand.Intn(30)+7) * 24 * time.Hour), Valid: true}
				status = statuses[1]
				statusCategory = domain.StatusCategoryInProgress
			}
		}

		issueType := issueTypes[rand.Intn(len(issueTypes))]
		priority := priorities[rand.Intn(len(priorities))]
		assignee := fmt.Sprintf("user%d@company.com", rand.Intn(50))

		issue := &domain.Issue{
			JiraIssueID:       fmt.Sprintf("%d", 30000+i),
			JiraIssueKey:      fmt.Sprintf("%s-%d", project.Key, i%1000+1),
			ProjectID:         project.ID,
			Summary:           fmt.Sprintf("[%s] %s対応 - タスク%d", issueType, project.Name, i+1),
			Status:            status,
			StatusCategory:    statusCategory,
			DueDate:           dueDate,
			DelayStatus:       delayStatus,
			AssigneeName:      &assignee,
			AssigneeAccountID: &assignee,
			Priority:          &priority,
			IssueType:         &issueType,
		}

		batch = append(batch, issue)

		// Insert in batches for better performance
		if len(batch) >= batchSize {
			for _, iss := range batch {
				if err := repo.Create(ctx, iss); err != nil {
					log.Printf("Warning: failed to create issue: %v", err)
				}
			}
			issues = append(issues, batch...)
			batch = make([]*domain.Issue, 0, batchSize)

			// Progress indicator
			log.Printf("  Created %d/%d issues", i+1, TARGET_ISSUES)
		}
	}

	// Insert remaining batch
	if len(batch) > 0 {
		for _, iss := range batch {
			if err := repo.Create(ctx, iss); err != nil {
				log.Printf("Warning: failed to create issue: %v", err)
			}
		}
		issues = append(issues, batch...)
	}

	return issues
}
