package usecase

import (
	"context"

	"github.com/m19cmjigen/sandbox-project-management/backend/internal/domain"
	"github.com/m19cmjigen/sandbox-project-management/backend/internal/interface/repository"
)

// ProjectUsecase はプロジェクトのユースケース
type ProjectUsecase interface {
	GetAll(ctx context.Context) ([]domain.Project, error)
	GetAllWithStats(ctx context.Context) ([]domain.ProjectWithStats, error)
	GetByID(ctx context.Context, id int64) (*domain.Project, error)
	GetWithStats(ctx context.Context, id int64) (*domain.ProjectWithStats, error)
	GetByOrganization(ctx context.Context, organizationID int64) ([]domain.Project, error)
	GetUnassigned(ctx context.Context) ([]domain.Project, error)
	AssignToOrganization(ctx context.Context, projectID int64, organizationID *int64) error
}

type projectUsecase struct {
	projectRepo repository.ProjectRepository
}

// NewProjectUsecase は新しいProjectUsecaseを作成
func NewProjectUsecase(projectRepo repository.ProjectRepository) ProjectUsecase {
	return &projectUsecase{
		projectRepo: projectRepo,
	}
}

// GetAll は全てのプロジェクトを取得
func (u *projectUsecase) GetAll(ctx context.Context) ([]domain.Project, error) {
	return u.projectRepo.FindAll(ctx)
}

// GetAllWithStats は全てのプロジェクトを統計情報付きで取得
func (u *projectUsecase) GetAllWithStats(ctx context.Context) ([]domain.ProjectWithStats, error) {
	return u.projectRepo.FindAllWithStats(ctx)
}

// GetByID はIDでプロジェクトを取得
func (u *projectUsecase) GetByID(ctx context.Context, id int64) (*domain.Project, error) {
	return u.projectRepo.FindByID(ctx, id)
}

// GetWithStats はIDで統計情報付きプロジェクトを取得
func (u *projectUsecase) GetWithStats(ctx context.Context, id int64) (*domain.ProjectWithStats, error) {
	return u.projectRepo.FindWithStats(ctx, id)
}

// GetByOrganization は組織IDでプロジェクトを取得
func (u *projectUsecase) GetByOrganization(ctx context.Context, organizationID int64) ([]domain.Project, error) {
	return u.projectRepo.FindByOrganizationID(ctx, organizationID)
}

// GetUnassigned は未分類プロジェクトを取得
func (u *projectUsecase) GetUnassigned(ctx context.Context) ([]domain.Project, error) {
	return u.projectRepo.FindUnassigned(ctx)
}

// AssignToOrganization はプロジェクトを組織に紐付け
func (u *projectUsecase) AssignToOrganization(ctx context.Context, projectID int64, organizationID *int64) error {
	return u.projectRepo.AssignToOrganization(ctx, projectID, organizationID)
}
