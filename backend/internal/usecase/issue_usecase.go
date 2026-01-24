package usecase

import (
	"context"

	"github.com/m19cmjigen/sandbox-project-management/backend/internal/domain"
	"github.com/m19cmjigen/sandbox-project-management/backend/internal/interface/repository"
)

// IssueUsecase はチケットのユースケース
type IssueUsecase interface {
	GetAll(ctx context.Context) ([]domain.Issue, error)
	GetByID(ctx context.Context, id int64) (*domain.Issue, error)
	GetByProjectID(ctx context.Context, projectID int64) ([]domain.Issue, error)
	GetByFilter(ctx context.Context, filter domain.IssueFilter) ([]domain.Issue, error)
	GetByDelayStatus(ctx context.Context, status domain.DelayStatus) ([]domain.Issue, error)
}

type issueUsecase struct {
	issueRepo repository.IssueRepository
}

// NewIssueUsecase は新しいIssueUsecaseを作成
func NewIssueUsecase(issueRepo repository.IssueRepository) IssueUsecase {
	return &issueUsecase{
		issueRepo: issueRepo,
	}
}

// GetAll は全てのチケットを取得
func (u *issueUsecase) GetAll(ctx context.Context) ([]domain.Issue, error) {
	return u.issueRepo.FindAll(ctx)
}

// GetByID はIDでチケットを取得
func (u *issueUsecase) GetByID(ctx context.Context, id int64) (*domain.Issue, error) {
	return u.issueRepo.FindByID(ctx, id)
}

// GetByProjectID はプロジェクトIDでチケットを取得
func (u *issueUsecase) GetByProjectID(ctx context.Context, projectID int64) ([]domain.Issue, error) {
	return u.issueRepo.FindByProjectID(ctx, projectID)
}

// GetByFilter はフィルタ条件でチケットを取得
func (u *issueUsecase) GetByFilter(ctx context.Context, filter domain.IssueFilter) ([]domain.Issue, error) {
	return u.issueRepo.FindByFilter(ctx, filter)
}

// GetByDelayStatus は遅延ステータスでチケットを取得
func (u *issueUsecase) GetByDelayStatus(ctx context.Context, status domain.DelayStatus) ([]domain.Issue, error) {
	return u.issueRepo.FindByDelayStatus(ctx, status)
}
