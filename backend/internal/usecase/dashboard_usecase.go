package usecase

import (
	"context"

	"github.com/m19cmjigen/sandbox-project-management/backend/internal/domain"
	"github.com/m19cmjigen/sandbox-project-management/backend/internal/interface/repository"
)

// DashboardSummary はダッシュボード全体のサマリ
type DashboardSummary struct {
	TotalProjects    int                          `json:"total_projects"`
	DelayedProjects  int                          `json:"delayed_projects"`
	WarningProjects  int                          `json:"warning_projects"`
	NormalProjects   int                          `json:"normal_projects"`
	TotalIssues      int                          `json:"total_issues"`
	RedIssues        int                          `json:"red_issues"`
	YellowIssues     int                          `json:"yellow_issues"`
	GreenIssues      int                          `json:"green_issues"`
	ProjectsByStatus []domain.ProjectWithStats    `json:"projects_by_status"`
}

// OrganizationSummary は組織別のサマリ
type OrganizationSummary struct {
	Organization     domain.Organization          `json:"organization"`
	TotalProjects    int                          `json:"total_projects"`
	DelayedProjects  int                          `json:"delayed_projects"`
	WarningProjects  int                          `json:"warning_projects"`
	Projects         []domain.ProjectWithStats    `json:"projects"`
}

// DashboardUsecase はダッシュボードのユースケース
type DashboardUsecase interface {
	GetSummary(ctx context.Context) (*DashboardSummary, error)
	GetOrganizationSummary(ctx context.Context, organizationID int64) (*OrganizationSummary, error)
	GetProjectSummary(ctx context.Context, projectID int64) (*domain.ProjectWithStats, error)
}

type dashboardUsecase struct {
	orgRepo     repository.OrganizationRepository
	projectRepo repository.ProjectRepository
	issueRepo   repository.IssueRepository
}

// NewDashboardUsecase は新しいDashboardUsecaseを作成
func NewDashboardUsecase(
	orgRepo repository.OrganizationRepository,
	projectRepo repository.ProjectRepository,
	issueRepo repository.IssueRepository,
) DashboardUsecase {
	return &dashboardUsecase{
		orgRepo:     orgRepo,
		projectRepo: projectRepo,
		issueRepo:   issueRepo,
	}
}

// GetSummary はダッシュボード全体のサマリを取得
func (u *dashboardUsecase) GetSummary(ctx context.Context) (*DashboardSummary, error) {
	// 全プロジェクトを統計情報付きで取得
	projects, err := u.projectRepo.FindAllWithStats(ctx)
	if err != nil {
		return nil, err
	}

	summary := &DashboardSummary{
		TotalProjects:    len(projects),
		ProjectsByStatus: projects,
	}

	// 各プロジェクトの統計を集計
	for _, p := range projects {
		summary.TotalIssues += p.TotalIssues
		summary.RedIssues += p.RedIssues
		summary.YellowIssues += p.YellowIssues
		summary.GreenIssues += p.GreenIssues

		// プロジェクトのステータス判定
		if p.RedIssues > 0 {
			summary.DelayedProjects++
		} else if p.YellowIssues > 0 {
			summary.WarningProjects++
		} else {
			summary.NormalProjects++
		}
	}

	return summary, nil
}

// GetOrganizationSummary は組織別のサマリを取得
func (u *dashboardUsecase) GetOrganizationSummary(ctx context.Context, organizationID int64) (*OrganizationSummary, error) {
	// 組織情報を取得
	org, err := u.orgRepo.FindByID(ctx, organizationID)
	if err != nil {
		return nil, err
	}
	if org == nil {
		return nil, nil
	}

	// 組織配下のプロジェクトを取得
	projects, err := u.projectRepo.FindByOrganizationID(ctx, organizationID)
	if err != nil {
		return nil, err
	}

	// 各プロジェクトの統計を取得
	var projectsWithStats []domain.ProjectWithStats
	delayedCount := 0
	warningCount := 0

	for _, p := range projects {
		stats, err := u.projectRepo.FindWithStats(ctx, p.ID)
		if err != nil {
			return nil, err
		}
		if stats != nil {
			projectsWithStats = append(projectsWithStats, *stats)

			if stats.RedIssues > 0 {
				delayedCount++
			} else if stats.YellowIssues > 0 {
				warningCount++
			}
		}
	}

	summary := &OrganizationSummary{
		Organization:    *org,
		TotalProjects:   len(projects),
		DelayedProjects: delayedCount,
		WarningProjects: warningCount,
		Projects:        projectsWithStats,
	}

	return summary, nil
}

// GetProjectSummary はプロジェクト別のサマリを取得
func (u *dashboardUsecase) GetProjectSummary(ctx context.Context, projectID int64) (*domain.ProjectWithStats, error) {
	return u.projectRepo.FindWithStats(ctx, projectID)
}
