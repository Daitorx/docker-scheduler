package models

// ExecutionLog represents a container execution record
type ExecutionLog struct {
	ID            int    `json:"id"`
	ScheduleID    int    `json:"schedule_id,omitempty"`
	ContainerName string `json:"container_name"`
	DockerName    string `json:"docker_name,omitempty"`    // Actual docker container name
	ScheduledTime string `json:"scheduled_time,omitempty"` // Original scheduled time (e.g., "09:00")
	RandomDelay   int    `json:"random_delay,omitempty"`   // Random delay in minutes if applicable
	ExecutedAt    string `json:"executed_at"`
	Status        string `json:"status,omitempty"` // running, success, error
	Success       bool   `json:"success"`
	Output        string `json:"output,omitempty"`
	Error         string `json:"error,omitempty"`
}
