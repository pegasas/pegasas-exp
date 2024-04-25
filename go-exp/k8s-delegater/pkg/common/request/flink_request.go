package request

type FlinkRequest struct {
	DeploymentName string `json:"deployment_name"`
	JobId          string `json:"job_id"`
	TriggerId      string `json:"trigger_id"`
}
