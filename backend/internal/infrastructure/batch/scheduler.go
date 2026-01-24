package batch

import (
	"context"
	"log"
	"time"

	"github.com/m19cmjigen/sandbox-project-management/backend/internal/usecase"
)

// Scheduler handles scheduled batch jobs
type Scheduler struct {
	syncUsecase    usecase.SyncUsecase
	organizationID int64
	interval       time.Duration
	stopChan       chan struct{}
}

// NewScheduler creates a new batch job scheduler
func NewScheduler(syncUsecase usecase.SyncUsecase, organizationID int64, interval time.Duration) *Scheduler {
	return &Scheduler{
		syncUsecase:    syncUsecase,
		organizationID: organizationID,
		interval:       interval,
		stopChan:       make(chan struct{}),
	}
}

// Start begins the scheduled sync jobs
func (s *Scheduler) Start(ctx context.Context) {
	log.Printf("Starting batch scheduler with interval: %v", s.interval)

	// Run immediately on start
	s.runSync(ctx)

	ticker := time.NewTicker(s.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			s.runSync(ctx)
		case <-s.stopChan:
			log.Println("Batch scheduler stopped")
			return
		case <-ctx.Done():
			log.Println("Batch scheduler cancelled by context")
			return
		}
	}
}

// Stop stops the scheduler gracefully
func (s *Scheduler) Stop() {
	close(s.stopChan)
}

// runSync executes a single sync operation
func (s *Scheduler) runSync(ctx context.Context) {
	log.Println("Starting scheduled Jira sync...")
	startTime := time.Now()

	syncLog, err := s.syncUsecase.SyncAllProjects(ctx, s.organizationID)
	if err != nil {
		log.Printf("Scheduled sync failed: %v", err)
		return
	}

	duration := time.Since(startTime)
	log.Printf("Scheduled sync completed in %v: Status=%s, Projects=%d, Issues=%d, Errors=%d",
		duration, syncLog.Status, syncLog.ProjectsSynced, syncLog.IssuesSynced, syncLog.ErrorCount)
}

// RunOnce executes a sync job once (useful for manual triggers)
func (s *Scheduler) RunOnce(ctx context.Context) error {
	log.Println("Running one-time Jira sync...")
	syncLog, err := s.syncUsecase.SyncAllProjects(ctx, s.organizationID)
	if err != nil {
		return err
	}

	log.Printf("One-time sync result: Status=%s, Projects=%d, Issues=%d, Errors=%d",
		syncLog.Status, syncLog.ProjectsSynced, syncLog.IssuesSynced, syncLog.ErrorCount)
	return nil
}
