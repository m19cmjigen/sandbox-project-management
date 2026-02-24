package metrics

// NoopRecorder is a Recorder that discards all metrics.
// Use it in local development and unit tests.
type NoopRecorder struct{}

func (NoopRecorder) RecordSync(_ SyncResult) {}
