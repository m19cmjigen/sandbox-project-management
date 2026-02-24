package metrics

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
	"time"
)

func TestEMFRecorder_RecordSync_Success(t *testing.T) {
	var buf bytes.Buffer
	rec := NewEMFRecorder("TestNamespace/Batch", &buf)

	rec.RecordSync(SyncResult{
		SyncType:       "FULL",
		Success:        true,
		Duration:       10 * time.Second,
		ProjectsSynced: 3,
		IssuesSynced:   42,
	})

	line := strings.TrimSpace(buf.String())
	if line == "" {
		t.Fatal("expected non-empty output")
	}

	var m map[string]interface{}
	if err := json.Unmarshal([]byte(line), &m); err != nil {
		t.Fatalf("output is not valid JSON: %v\noutput: %s", err, line)
	}

	// SyncSuccess should be 1 for success
	if got := m["SyncSuccess"].(float64); got != 1 {
		t.Errorf("expected SyncSuccess=1, got %v", got)
	}
	if got := m["SyncType"].(string); got != "FULL" {
		t.Errorf("expected SyncType=FULL, got %s", got)
	}
	if got := m["DurationSeconds"].(float64); got != 10.0 {
		t.Errorf("expected DurationSeconds=10, got %v", got)
	}
	if got := m["IssuesSynced"].(float64); got != 42 {
		t.Errorf("expected IssuesSynced=42, got %v", got)
	}
	if got := m["ProjectsSynced"].(float64); got != 3 {
		t.Errorf("expected ProjectsSynced=3, got %v", got)
	}
}

func TestEMFRecorder_RecordSync_Failure(t *testing.T) {
	var buf bytes.Buffer
	rec := NewEMFRecorder("TestNamespace/Batch", &buf)

	rec.RecordSync(SyncResult{
		SyncType: "DELTA",
		Success:  false,
		Duration: 2 * time.Second,
	})

	var m map[string]interface{}
	if err := json.Unmarshal([]byte(strings.TrimSpace(buf.String())), &m); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}

	// SyncSuccess should be 0 for failure
	if got := m["SyncSuccess"].(float64); got != 0 {
		t.Errorf("expected SyncSuccess=0, got %v", got)
	}
	if got := m["SyncType"].(string); got != "DELTA" {
		t.Errorf("expected SyncType=DELTA, got %s", got)
	}
}

func TestEMFRecorder_ContainsEMFMetadata(t *testing.T) {
	var buf bytes.Buffer
	rec := NewEMFRecorder("MyApp/Batch", &buf)
	rec.RecordSync(SyncResult{SyncType: "FULL", Success: true, Duration: time.Second})

	var m map[string]interface{}
	if err := json.Unmarshal([]byte(strings.TrimSpace(buf.String())), &m); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}

	aws, ok := m["_aws"].(map[string]interface{})
	if !ok {
		t.Fatal("missing _aws field in EMF output")
	}
	if _, ok := aws["Timestamp"]; !ok {
		t.Error("missing _aws.Timestamp")
	}
	if _, ok := aws["CloudWatchMetrics"]; !ok {
		t.Error("missing _aws.CloudWatchMetrics")
	}
}

func TestNoopRecorder_DoesNotPanic(t *testing.T) {
	var rec NoopRecorder
	// NoopRecorder should not panic on any input
	rec.RecordSync(SyncResult{SyncType: "FULL", Success: true})
	rec.RecordSync(SyncResult{SyncType: "DELTA", Success: false})
}
