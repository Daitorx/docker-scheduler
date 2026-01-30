package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// SystemStatus represents the system time and status
type SystemStatus struct {
	ServerTime string `json:"server_time"`
	Timezone   string `json:"timezone"`
	Uptime     string `json:"uptime"` // Placeholder for now
}

// GetSystemStatus returns the current server time and timezone
func GetSystemStatus(c *gin.Context) {
	now := time.Now()
	zoneName, _ := now.Zone()

	status := SystemStatus{
		ServerTime: now.Format("2006-01-02 15:04:05"),
		Timezone:   zoneName,
		Uptime:     "N/A", // Can implement uptime tracking if needed
	}

	c.JSON(http.StatusOK, status)
}
