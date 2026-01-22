package models

// TimeRange represents a time range for random scheduling

// Schedule represents a scheduled Docker container execution
type Schedule struct {
	ID             int               `json:"id"`
	ContainerName  string            `json:"container_name"`
	RunName        string            `json:"run_name,omitempty"` // Custom container name (--name)
	Ports          []string          `json:"ports,omitempty"`    // Port mappings ["8080:80", "3000:3000"]
	AutoRemove     bool              `json:"auto_remove"`        // Use --rm flag (default true)
	Days           []string          `json:"days"`               // ["monday", "tuesday", ...]
	Times          []string          `json:"times"`              // ["09:00", "14:30", "18:00"]
	RandomDelay    int               `json:"random_delay"`       // Optional random delay in minutes
	ExceptionDates []string          `json:"exception_dates"`    // ["2026-01-20", "2026-02-03"]
	EnvVars        map[string]string `json:"env_vars"`           // {"KEY": "value", ...}
	Active         bool              `json:"active"`
	CreatedAt      string            `json:"created_at"`
}

// CreateScheduleRequest is the request body for creating a schedule
type CreateScheduleRequest struct {
	ContainerName  string            `json:"container_name" binding:"required"`
	RunName        string            `json:"run_name"`
	Ports          []string          `json:"ports"`
	AutoRemove     *bool             `json:"auto_remove"` // Pointer to detect if provided (default true)
	Days           []string          `json:"days" binding:"required"`
	Times          []string          `json:"times"` // Validated manually (required if no random_time_range)
	RandomDelay    int               `json:"random_delay"`
	ExceptionDates []string          `json:"exception_dates"`
	EnvVars        map[string]string `json:"env_vars"`
}
