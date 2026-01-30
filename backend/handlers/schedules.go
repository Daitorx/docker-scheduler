package handlers

import (
	"fmt"
	"net/http"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"docker-scheduler/database"
	"docker-scheduler/models"
	"docker-scheduler/scheduler"

	"github.com/gin-gonic/gin"
)

// GetSchedules returns all schedules
func GetSchedules(c *gin.Context) {
	schedules, err := database.GetAllSchedules()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if schedules == nil {
		schedules = []models.Schedule{}
	}

	c.JSON(http.StatusOK, schedules)
}

// CreateSchedule creates a new schedule
func CreateSchedule(c *gin.Context) {
	var req models.CreateScheduleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Manual validation for times
	if len(req.Times) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "At least one time is required"})
		return
	}

	// Default autoRemove to true if not provided
	autoRemove := true
	if req.AutoRemove != nil {
		autoRemove = *req.AutoRemove
	}

	// Generate random name if not provided
	runName := req.RunName
	if runName == "" {
		runName = fmt.Sprintf("%s-%d", req.ContainerName, time.Now().Unix())
	}

	schedule, err := database.CreateSchedule(req.ContainerName, runName, req.Ports, autoRemove, req.Days, req.Times, req.ExceptionDates, req.EnvVars, req.RandomDelay)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Refresh scheduler to pick up the new schedule
	scheduler.RefreshSchedules()

	c.JSON(http.StatusCreated, schedule)
}

// DeleteSchedule removes a schedule by ID
func DeleteSchedule(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	if err := database.DeleteSchedule(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Refresh scheduler
	scheduler.RefreshSchedules()

	c.JSON(http.StatusOK, gin.H{"message": "Schedule deleted"})
}

// ToggleSchedule toggles the active status
func ToggleSchedule(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	if err := database.ToggleSchedule(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Refresh scheduler
	scheduler.RefreshSchedules()

	schedule, _ := database.GetScheduleByID(id)
	c.JSON(http.StatusOK, schedule)
}

// RunSchedule triggers a schedule immediately
func RunSchedule(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	if err := scheduler.RunScheduleNow(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Schedule triggered successfully"})
}

// ExceptionRequest is the request body for exception date operations
type ExceptionRequest struct {
	Date string `json:"date" binding:"required"`
}

// AddException adds an exception date to a schedule
func AddException(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var req ExceptionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := database.AddExceptionDate(id, req.Date); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Refresh scheduler
	scheduler.RefreshSchedules()

	schedule, _ := database.GetScheduleByID(id)
	c.JSON(http.StatusOK, schedule)
}

// RemoveException removes an exception date from a schedule
func RemoveException(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var req ExceptionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := database.RemoveExceptionDate(id, req.Date); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Refresh scheduler
	scheduler.RefreshSchedules()

	schedule, _ := database.GetScheduleByID(id)
	c.JSON(http.StatusOK, schedule)
}

// GetExecutionHistory returns the execution history
func GetExecutionHistory(c *gin.Context) {
	limitStr := c.DefaultQuery("limit", "50")
	limit, _ := strconv.Atoi(limitStr)

	logs, err := database.GetExecutionHistory(limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if logs == nil {
		logs = []models.ExecutionLog{}
	}

	c.JSON(http.StatusOK, logs)
}

// ClearHistory deletes all execution history
func ClearHistory(c *gin.Context) {
	if err := database.ClearExecutionHistory(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "History cleared"})
}

// ContainerInfo represents a running container
type ContainerInfo struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Image string `json:"image"`
	State string `json:"state"`
	Ports string `json:"ports"`
}

// GetRunningContainers returns a list of running Docker containers
func GetRunningContainers(c *gin.Context) {
	cmd := exec.Command("docker", "ps", "-a", "--format", "{{.ID}}|{{.Names}}|{{.Image}}|{{.State}}|{{.Ports}}")
	output, err := cmd.Output()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list containers"})
		return
	}

	var containers []ContainerInfo
	lines := strings.Split(strings.TrimSpace(string(output)), "\n")
	for _, line := range lines {
		if line == "" {
			continue
		}
		parts := strings.SplitN(line, "|", 5)
		if len(parts) >= 4 {
			ports := ""
			if len(parts) == 5 {
				ports = parts[4]
			}
			containers = append(containers, ContainerInfo{
				ID:    parts[0],
				Name:  parts[1],
				Image: parts[2],
				State: parts[3],
				Ports: ports,
			})
		}
	}

	if containers == nil {
		containers = []ContainerInfo{}
	}

	c.JSON(http.StatusOK, containers)
}

// StopContainer stops a running Docker container by name or ID
func StopContainer(c *gin.Context) {
	name := c.Param("name")
	if name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Container name is required"})
		return
	}

	cmd := exec.Command("docker", "stop", name)
	output, err := cmd.CombinedOutput()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": string(output)})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Container stopped", "container": name})
}

// RemoveContainer removes a stopped container
func RemoveContainer(c *gin.Context) {
	name := c.Param("name")
	if name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Container name is required"})
		return
	}

	cmd := exec.Command("docker", "rm", "-f", name)
	output, err := cmd.CombinedOutput()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": string(output)})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Container removed", "container": name})
}

// GetContainerLogs returns logs for a specific container
func GetContainerLogs(c *gin.Context) {
	name := c.Param("name")
	if name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Container name is required"})
		return
	}

	cmd := exec.Command("docker", "logs", name)
	output, err := cmd.CombinedOutput()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error(), "output": string(output)})
		return
	}

	c.JSON(http.StatusOK, gin.H{"logs": string(output)})
}
