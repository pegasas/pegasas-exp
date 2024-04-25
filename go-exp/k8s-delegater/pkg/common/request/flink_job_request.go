package request

type FlinkJobRequest struct {
	JobId           string            `json:"job_id"`
	ReleaseName     string            `json:"release_name"`
	ArtifactId      string            `json:"artifact_id"`
	BuildId         string            `json:"build_id"`
	DWLayer         string            `json:"dw_layer"`
	Application     string            `json:"application"`
	ClusterId       string            `json:"cluster_id"`
	FlinkParameters map[string]string `json:"flink_parameters"`
}
