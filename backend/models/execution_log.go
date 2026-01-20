package models

// ExecutionLog represents a container execution record
type ExecutionLog struct {
	ID            int    `json:"id"`
	ScheduleID    int    `json:"schedule_id,omitempty"`
	ContainerName string `json:"container_name"`
	DockerName    string `json:"docker_name,omitempty"` // Actual docker container name
	ExecutedAt    string `json:"executed_at"`
	Status        string `json:"status,omitempty"` // running, success, error
	Success       bool   `json:"success"`
	Output        string `json:"output,omitempty"`
	Error         string `json:"error,omitempty"`
}
