package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/m19cmjigen/sandbox-project-management/backend/internal/domain"
	"github.com/m19cmjigen/sandbox-project-management/backend/internal/infrastructure/postgres"
	"github.com/m19cmjigen/sandbox-project-management/backend/internal/interface/repository"
	"golang.org/x/crypto/bcrypt"
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
	userRepo := postgres.NewUserRepository(db)

	log.Println("Starting seed data creation...")

	// 1. Create users
	log.Println("Creating users...")
	users := createUsers(ctx, userRepo)
	log.Printf("Created %d users", len(users))

	// 2. Create organizations
	log.Println("Creating organization hierarchy...")
	orgs := createOrganizations(ctx, orgRepo)
	log.Printf("Created %d organizations", len(orgs))

	// 3. Create projects
	log.Println("Creating projects...")
	projects := createProjects(ctx, projectRepo, orgs)
	log.Printf("Created %d projects", len(projects))

	// 4. Create issues
	log.Println("Creating issues...")
	issues := createIssues(ctx, issueRepo, projects)
	log.Printf("Created %d issues", len(issues))

	log.Println("Seed data creation completed successfully!")
}

func createUsers(ctx context.Context, repo domain.UserRepository) []*domain.User {
	// Hash password for all users
	hash, _ := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)

	users := []*domain.User{
		{
			Username:     "admin",
			Email:        "admin@company.com",
			PasswordHash: string(hash),
			FullName:     sql.NullString{String: "System Administrator", Valid: true},
			Role:         domain.RoleAdmin,
			IsActive:     true,
		},
		{
			Username:     "pmo_manager",
			Email:        "pmo@company.com",
			PasswordHash: string(hash),
			FullName:     sql.NullString{String: "PMO Manager", Valid: true},
			Role:         domain.RoleManager,
			IsActive:     true,
		},
		{
			Username:     "dept_manager",
			Email:        "dept.manager@company.com",
			PasswordHash: string(hash),
			FullName:     sql.NullString{String: "Department Manager", Valid: true},
			Role:         domain.RoleManager,
			IsActive:     true,
		},
		{
			Username:     "viewer",
			Email:        "viewer@company.com",
			PasswordHash: string(hash),
			FullName:     sql.NullString{String: "Regular Viewer", Valid: true},
			Role:         domain.RoleViewer,
			IsActive:     true,
		},
	}

	for _, user := range users {
		if err := repo.Create(ctx, user); err != nil {
			log.Printf("Warning: failed to create user %s: %v", user.Username, err)
		}
	}

	return users
}

func createOrganizations(ctx context.Context, repo repository.OrganizationRepository) []*domain.Organization {
	// Root organization
	headquarters := &domain.Organization{
		Name:     "本社 (Headquarters)",
		ParentID: nil,
	}
	if err := repo.Create(ctx, headquarters); err != nil {
		log.Fatalf("Failed to create headquarters: %v", err)
	}

	// Departments
	salesDept := &domain.Organization{
		Name:     "営業本部 (Sales Division)",
		ParentID: &headquarters.ID,
	}
	if err := repo.Create(ctx, salesDept); err != nil {
		log.Fatalf("Failed to create sales dept: %v", err)
	}

	devDept := &domain.Organization{
		Name:     "開発本部 (Development Division)",
		ParentID: &headquarters.ID,
	}
	if err := repo.Create(ctx, devDept); err != nil {
		log.Fatalf("Failed to create dev dept: %v", err)
	}

	pmoDept := &domain.Organization{
		Name:     "PMO部 (PMO Department)",
		ParentID: &headquarters.ID,
	}
	if err := repo.Create(ctx, pmoDept); err != nil {
		log.Fatalf("Failed to create PMO dept: %v", err)
	}

	// Sales sub-departments
	tokyoSales := &domain.Organization{
		Name:     "東京営業所 (Tokyo Sales Office)",
		ParentID: &salesDept.ID,
	}
	repo.Create(ctx, tokyoSales)

	osakaSales := &domain.Organization{
		Name:     "大阪営業所 (Osaka Sales Office)",
		ParentID: &salesDept.ID,
	}
	repo.Create(ctx, osakaSales)

	// Development sub-departments
	frontendTeam := &domain.Organization{
		Name:     "フロントエンドチーム (Frontend Team)",
		ParentID: &devDept.ID,
	}
	repo.Create(ctx, frontendTeam)

	backendTeam := &domain.Organization{
		Name:     "バックエンドチーム (Backend Team)",
		ParentID: &devDept.ID,
	}
	repo.Create(ctx, backendTeam)

	infraTeam := &domain.Organization{
		Name:     "インフラチーム (Infrastructure Team)",
		ParentID: &devDept.ID,
	}
	repo.Create(ctx, infraTeam)

	return []*domain.Organization{
		headquarters, salesDept, devDept, pmoDept,
		tokyoSales, osakaSales, frontendTeam, backendTeam, infraTeam,
	}
}

func createProjects(ctx context.Context, repo repository.ProjectRepository, orgs []*domain.Organization) []*domain.Project {
	projects := []*domain.Project{}

	// Find organizations for assignment
	var salesDept, devDept, pmo *domain.Organization
	for _, org := range orgs {
		if org.Name == "営業本部 (Sales Division)" {
			salesDept = org
		} else if org.Name == "開発本部 (Development Division)" {
			devDept = org
		} else if org.Name == "PMO部 (PMO Department)" {
			pmo = org
		}
	}

	// Sales projects
	leadEmail := "tanaka@company.com"
	p1 := &domain.Project{
		JiraProjectID:  "10001",
		Key:            "SALES",
		Name:           "新規CRMシステム導入 (New CRM Implementation)",
		LeadAccountID:  &leadEmail,
		OrganizationID: &salesDept.ID,
	}
	repo.Create(ctx, p1)
	projects = append(projects, p1)

	leadEmail2 := "suzuki@company.com"
	p2 := &domain.Project{
		JiraProjectID:  "10002",
		Key:            "MARKET",
		Name:           "マーケティングオートメーション (Marketing Automation)",
		LeadAccountID:  &leadEmail2,
		OrganizationID: &salesDept.ID,
	}
	repo.Create(ctx, p2)
	projects = append(projects, p2)

	// Development projects
	leadEmail3 := "yamada@company.com"
	p3 := &domain.Project{
		JiraProjectID:  "10003",
		Key:            "WEBAPP",
		Name:           "社内Webアプリ刷新 (Internal Web App Renewal)",
		LeadAccountID:  &leadEmail3,
		OrganizationID: &devDept.ID,
	}
	repo.Create(ctx, p3)
	projects = append(projects, p3)

	leadEmail4 := "sato@company.com"
	p4 := &domain.Project{
		JiraProjectID:  "10004",
		Key:            "MOBILE",
		Name:           "モバイルアプリ開発 (Mobile App Development)",
		LeadAccountID:  &leadEmail4,
		OrganizationID: &devDept.ID,
	}
	repo.Create(ctx, p4)
	projects = append(projects, p4)

	// PMO project
	leadEmail5 := "kato@company.com"
	p5 := &domain.Project{
		JiraProjectID:  "10005",
		Key:            "PMOBASE",
		Name:           "プロジェクト管理基盤構築 (PM Platform Construction)",
		LeadAccountID:  &leadEmail5,
		OrganizationID: &pmo.ID,
	}
	repo.Create(ctx, p5)
	projects = append(projects, p5)

	return projects
}

func createIssues(ctx context.Context, repo repository.IssueRepository, projects []*domain.Project) []*domain.Issue {
	issues := []*domain.Issue{}
	now := time.Now()

	for i, project := range projects {
		// Create 5-10 issues per project with different statuses
		numIssues := 5 + i%6

		for j := 0; j < numIssues; j++ {
			var dueDate sql.NullTime
			var delayStatus domain.DelayStatus
			var status string
			var statusCategory domain.StatusCategory
			assignee := fmt.Sprintf("user%d@company.com", (j%3)+1)
			priority := []string{"High", "Medium", "Low"}[j%3]

			// Mix of different delay statuses
			switch j % 5 {
			case 0: // Overdue (RED)
				dueDate = sql.NullTime{Time: now.Add(-7 * 24 * time.Hour), Valid: true}
				delayStatus = domain.DelayStatusRed
				status = "In Progress"
				statusCategory = domain.StatusCategoryInProgress
			case 1: // Due soon (YELLOW)
				dueDate = sql.NullTime{Time: now.Add(2 * 24 * time.Hour), Valid: true}
				delayStatus = domain.DelayStatusYellow
				status = "In Progress"
				statusCategory = domain.StatusCategoryInProgress
			case 2: // No due date (YELLOW)
				dueDate = sql.NullTime{Valid: false}
				delayStatus = domain.DelayStatusYellow
				status = "To Do"
				statusCategory = domain.StatusCategoryToDo
			case 3: // On track (GREEN)
				dueDate = sql.NullTime{Time: now.Add(14 * 24 * time.Hour), Valid: true}
				delayStatus = domain.DelayStatusGreen
				status = "In Progress"
				statusCategory = domain.StatusCategoryInProgress
			case 4: // Done
				dueDate = sql.NullTime{Time: now.Add(-1 * 24 * time.Hour), Valid: true}
				delayStatus = domain.DelayStatusGreen
				status = "Done"
				statusCategory = domain.StatusCategoryDone
			}

			issue := &domain.Issue{
				JiraIssueID:       fmt.Sprintf("%d%03d", 10000+int64(i), j+1),
				JiraIssueKey:      fmt.Sprintf("%s-%d", project.Key, j+1),
				ProjectID:         project.ID,
				Summary:           fmt.Sprintf("Task %d: %s", j+1, generateTaskName(j)),
				Status:            status,
				StatusCategory:    statusCategory,
				DueDate:           dueDate,
				DelayStatus:       delayStatus,
				AssigneeName:      &assignee,
				AssigneeAccountID: &assignee,
				Priority:          &priority,
			}

			if err := repo.Create(ctx, issue); err != nil {
				log.Printf("Warning: failed to create issue %s: %v", issue.JiraIssueKey, err)
			} else {
				issues = append(issues, issue)
			}
		}
	}

	return issues
}

func generateTaskName(index int) string {
	tasks := []string{
		"要件定義 (Requirements Definition)",
		"設計書作成 (Design Documentation)",
		"実装 (Implementation)",
		"単体テスト (Unit Testing)",
		"結合テスト (Integration Testing)",
		"リリース準備 (Release Preparation)",
		"本番デプロイ (Production Deployment)",
		"運用移管 (Operations Handover)",
		"ドキュメント整備 (Documentation)",
		"レビュー対応 (Review Response)",
	}
	return tasks[index%len(tasks)]
}
