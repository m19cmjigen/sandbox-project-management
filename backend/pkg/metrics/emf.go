package metrics

import (
	"encoding/json"
	"fmt"
	"io"
	"time"
)

// EMFRecorder emits metrics using CloudWatch Embedded Metrics Format (EMF).
// When running on ECS/Lambda with CloudWatch Logs Agent, metrics are automatically
// extracted as CloudWatch custom metrics without requiring direct API calls.
//
// Reference: https://docs.aws.amazon.com/AmazonCloudWatch/latest/monitoring/CloudWatch_Embedded_Metric_Format_Specification.html
type EMFRecorder struct {
	namespace string
	out       io.Writer
}

// NewEMFRecorder creates an EMFRecorder that writes to out.
// namespace is the CloudWatch metrics namespace (e.g. "SandboxProjectManagement/Batch").
func NewEMFRecorder(namespace string, out io.Writer) *EMFRecorder {
	return &EMFRecorder{namespace: namespace, out: out}
}

// RecordSync emits a structured JSON log line in EMF format.
func (r *EMFRecorder) RecordSync(m SyncResult) {
	successVal := 0
	if m.Success {
		successVal = 1
	}

	entry := map[string]interface{}{
		"_aws": map[string]interface{}{
			"Timestamp": time.Now().UnixMilli(),
			"CloudWatchMetrics": []map[string]interface{}{
				{
					"Namespace":  r.namespace,
					"Dimensions": [][]string{{"SyncType"}},
					"Metrics": []map[string]string{
						{"Name": "SyncSuccess", "Unit": "Count"},
						{"Name": "DurationSeconds", "Unit": "Seconds"},
						{"Name": "IssuesSynced", "Unit": "Count"},
						{"Name": "ProjectsSynced", "Unit": "Count"},
					},
				},
			},
		},
		"SyncType":        m.SyncType,
		"SyncSuccess":     successVal,
		"DurationSeconds": m.Duration.Seconds(),
		"IssuesSynced":    m.IssuesSynced,
		"ProjectsSynced":  m.ProjectsSynced,
	}

	data, err := json.Marshal(entry)
	if err != nil {
		// EMF出力失敗はバッチ処理を止めないため警告のみ
		fmt.Fprintf(r.out, `{"level":"warn","msg":"failed to marshal EMF metrics","error":%q}`+"\n", err.Error())
		return
	}
	fmt.Fprintln(r.out, string(data))
}
