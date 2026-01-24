package usecase

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/m19cmjigen/sandbox-project-management/backend/internal/domain"
	"github.com/m19cmjigen/sandbox-project-management/backend/internal/infrastructure/jira"
	"github.com/m19cmjigen/sandbox-project-management/backend/internal/interface/repository"
)

// SyncUsecase handles synchronization with Jira
type SyncUsecase interface {
	SyncAllProjects(ctx context.Context, organizationID int64) (*domain.SyncLog, error)
	SyncProjectIssues(ctx context.Context, projectID int64) error
}

type syncUsecase struct {
	jiraClient     *jira.Client
	projectRepo    repository.ProjectRepository
	issueRepo      repository.IssueRepository
	syncLogRepo    repository.SyncLogRepository
}

// NewSyncUsecase creates a new sync usecase
func NewSyncUsecase(
	jiraClient *jira.Client,
	projectRepo repository.ProjectRepository,
	issueRepo repository.IssueRepository,
	syncLogRepo repository.SyncLogRepository,
) SyncUsecase {
	return &syncUsecase{
		jiraClient:  jiraClient,
		projectRepo: projectRepo,
		issueRepo:   issueRepo,
		syncLogRepo: syncLogRepo,
	}
}

// SyncAllProjects synchronizes all Jira projects and their issues
func (u *syncUsecase) SyncAllProjects(ctx context.Context, organizationID int64) (*domain.SyncLog, error) {
	startTime := time.Now()
	syncLog := &domain.SyncLog{
		StartedAt:      startTime,
		Status:         "RUNNING",
		ProjectsSynced: 0,
		IssuesSynced:   0,
		ErrorCount:     0,
	}

	// Fetch all projects from Jira
	jiraProjects, err := u.jiraClient.GetProjects(ctx)
	if err != nil {
		syncLog.Status = "FAILED"
		errMsg := err.Error()
		syncLog.ErrorMessage = &errMsg
		return syncLog, fmt.Errorf("failed to fetch Jira projects: %w", err)
	}

	log.Printf("Found %d Jira projects", len(jiraProjects))

	// Sync each project
	for _, jiraProject := range jiraProjects {
		// Check if project already exists
		existingProject, err := u.projectRepo.GetByKey(ctx, jiraProject.Key)
		var projectID int64

		if err != nil || existingProject == nil {
			// Create new project
			project := jira.TransformProject(jiraProject, organizationID)
			if err := u.projectRepo.Create(ctx, project); err != nil {
				log.Printf("Failed to create project %s: %v", jiraProject.Key, err)
				syncLog.ErrorCount++
				continue
			}
			// Get the created project to retrieve its ID
			createdProject, err := u.projectRepo.GetByKey(ctx, jiraProject.Key)
			if err != nil || createdProject == nil {
				log.Printf("Failed to retrieve created project %s: %v", jiraProject.Key, err)
				syncLog.ErrorCount++
				continue
			}
			projectID = createdProject.ID
			log.Printf("Created project: %s (%s)", jiraProject.Name, jiraProject.Key)
		} else {
			projectID = existingProject.ID
			// Update project name if changed
			if existingProject.Name != jiraProject.Name {
				existingProject.Name = jiraProject.Name
				if _, err := u.projectRepo.Update(ctx, existingProject); err != nil {
					log.Printf("Failed to update project %s: %v", jiraProject.Key, err)
				}
			}
			log.Printf("Project already exists: %s (%s)", jiraProject.Name, jiraProject.Key)
		}

		// Sync issues for this project
		if err := u.syncProjectIssuesInternal(ctx, projectID, jiraProject.Key, syncLog); err != nil {
			log.Printf("Failed to sync issues for project %s: %v", jiraProject.Key, err)
			syncLog.ErrorCount++
		}

		syncLog.ProjectsSynced++
	}

	// Update sync log
	syncLog.CompletedAt = timePtr(time.Now())
	syncLog.Status = "COMPLETED"
	if syncLog.ErrorCount > 0 {
		syncLog.Status = "COMPLETED_WITH_ERRORS"
	}
	syncLog.ErrorMessage = nil

	// Persist sync log to database
	if err := u.syncLogRepo.Create(ctx, syncLog); err != nil {
		log.Printf("Warning: Failed to persist sync log: %v", err)
	}

	log.Printf("Sync completed: %d projects, %d issues, %d errors",
		syncLog.ProjectsSynced, syncLog.IssuesSynced, syncLog.ErrorCount)

	return syncLog, nil
}

// SyncProjectIssues synchronizes issues for a specific project
func (u *syncUsecase) SyncProjectIssues(ctx context.Context, projectID int64) error {
	project, err := u.projectRepo.GetByID(ctx, projectID)
	if err != nil {
		return fmt.Errorf("failed to get project: %w", err)
	}

	syncLog := &domain.SyncLog{
		StartedAt:      time.Now(),
		Status:         "RUNNING",
		ProjectsSynced: 0,
		IssuesSynced:   0,
		ErrorCount:     0,
	}

	return u.syncProjectIssuesInternal(ctx, projectID, project.Key, syncLog)
}

// syncProjectIssuesInternal is internal helper for syncing project issues
func (u *syncUsecase) syncProjectIssuesInternal(ctx context.Context, projectID int64, projectKey string, syncLog *domain.SyncLog) error {
	// Fetch all issues for this project
	jiraIssues, err := u.jiraClient.GetAllIssuesForProject(ctx, projectKey)
	if err != nil {
		return fmt.Errorf("failed to fetch issues for project %s: %w", projectKey, err)
	}

	log.Printf("Found %d issues for project %s", len(jiraIssues), projectKey)

	// Sync each issue
	for _, jiraIssue := range jiraIssues {
		issue, err := jira.TransformIssue(jiraIssue, projectID)
		if err != nil {
			log.Printf("Failed to transform issue %s: %v", jiraIssue.Key, err)
			syncLog.ErrorCount++
			continue
		}

		// Check if issue already exists
		existingIssue, err := u.issueRepo.GetByJiraKey(ctx, jiraIssue.Key)
		if err != nil || existingIssue == nil {
			// Create new issue
			if _, err := u.issueRepo.Create(ctx, issue); err != nil {
				log.Printf("Failed to create issue %s: %v", jiraIssue.Key, err)
				syncLog.ErrorCount++
				continue
			}
		} else {
			// Update existing issue
			issue.ID = existingIssue.ID
			if _, err := u.issueRepo.Update(ctx, issue); err != nil {
				log.Printf("Failed to update issue %s: %v", jiraIssue.Key, err)
				syncLog.ErrorCount++
				continue
			}
		}

		syncLog.IssuesSynced++
	}

	return nil
}

// failSyncLog logs a sync failure (removed database persistence)

func timePtr(t time.Time) *time.Time {
	return &t
}
