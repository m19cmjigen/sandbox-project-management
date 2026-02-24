package metrics

import "time"

// Recorder records batch sync metrics.
type Recorder interface {
	RecordSync(r SyncResult)
}

// SyncResult holds metrics for a single sync run.
type SyncResult struct {
	// SyncType is the type of sync: "FULL" or "DELTA".
	SyncType string
	// Success is true when the sync completed without error.
	Success bool
	// Duration is the total execution time of the sync.
	Duration time.Duration
	// ProjectsSynced is the number of projects upserted (meaningful for FULL sync only).
	ProjectsSynced int
	// IssuesSynced is the number of issues upserted.
	IssuesSynced int
}
