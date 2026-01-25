package usecase

import (
	"context"
	"log"

	"github.com/m19cmjigen/sandbox-project-management/backend/internal/domain"
	"github.com/m19cmjigen/sandbox-project-management/backend/internal/infrastructure/jira"
	"github.com/m19cmjigen/sandbox-project-management/backend/internal/interface/repository"
	"github.com/m19cmjigen/sandbox-project-management/backend/pkg/retry"
)

// SyncUsecaseWithRetry wraps SyncUsecase with retry logic
type SyncUsecaseWithRetry struct {
	baseUsecase SyncUsecase
	retryConfig *retry.RetryConfig
}

// NewSyncUsecaseWithRetry creates a new sync usecase with retry capability
func NewSyncUsecaseWithRetry(
	jiraClient *jira.Client,
	projectRepo repository.ProjectRepository,
	issueRepo repository.IssueRepository,
	syncLogRepo repository.SyncLogRepository,
	retryConfig *retry.RetryConfig,
) SyncUsecase {
	if retryConfig == nil {
		retryConfig = retry.DefaultRetryConfig()
	}

	baseUsecase := NewSyncUsecase(jiraClient, projectRepo, issueRepo, syncLogRepo)

	return &SyncUsecaseWithRetry{
		baseUsecase: baseUsecase,
		retryConfig: retryConfig,
	}
}

// SyncAllProjects synchronizes all Jira projects with retry
func (u *SyncUsecaseWithRetry) SyncAllProjects(ctx context.Context, organizationID int64) (*domain.SyncLog, error) {
	var result *domain.SyncLog
	var syncErr error

	err := retry.WithRetry(ctx, u.retryConfig, "SyncAllProjects", func(ctx context.Context) error {
		syncLog, err := u.baseUsecase.SyncAllProjects(ctx, organizationID)
		result = syncLog
		syncErr = err

		// Only retry if it's a retryable error
		if err != nil && retry.IsRetryableError(err) {
			return err
		}

		// Don't retry on success or non-retryable errors
		return nil
	})

	// If retry wrapper returned error, it means all retries failed
	if err != nil {
		log.Printf("SyncAllProjects failed after retries: %v", err)
		return result, syncErr
	}

	return result, syncErr
}

// SyncProjectIssues synchronizes project issues with retry
func (u *SyncUsecaseWithRetry) SyncProjectIssues(ctx context.Context, projectID int64) error {
	return retry.WithRetry(ctx, u.retryConfig, "SyncProjectIssues", func(ctx context.Context) error {
		err := u.baseUsecase.SyncProjectIssues(ctx, projectID)

		// Only retry if it's a retryable error
		if err != nil && retry.IsRetryableError(err) {
			return err
		}

		return nil
	})
}
