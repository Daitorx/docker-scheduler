package main

import (
	"log"

	"docker-scheduler/database"
	"docker-scheduler/handlers"
	"docker-scheduler/scheduler"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize database
	if err := database.Initialize("./data/schedules.db"); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// Initialize scheduler
	scheduler.Initialize()
	defer scheduler.Stop()

	// Setup Gin router
	r := gin.Default()

	// CORS configuration
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type"},
		AllowCredentials: true,
	}))

	// Serve static frontend files
	r.Static("/css", "./static/css")
	r.Static("/js", "./static/js")
	r.StaticFile("/", "./static/index.html")
	r.StaticFile("/index.html", "./static/index.html")

	// API routes
	api := r.Group("/api")
	{
		api.GET("/schedules", handlers.GetSchedules)
		api.POST("/schedules", handlers.CreateSchedule)
		api.DELETE("/schedules/:id", handlers.DeleteSchedule)
		api.PUT("/schedules/:id/toggle", handlers.ToggleSchedule)
		api.POST("/schedules/:id/run", handlers.RunSchedule)
		api.POST("/schedules/:id/exceptions", handlers.AddException)
		api.DELETE("/schedules/:id/exceptions", handlers.RemoveException)
		api.GET("/history", handlers.GetExecutionHistory)
		api.DELETE("/history", handlers.ClearHistory)
		api.GET("/containers", handlers.GetRunningContainers)
		api.GET("/containers/:name/logs", handlers.GetContainerLogs)
		api.POST("/containers/:name/stop", handlers.StopContainer)
		api.DELETE("/containers/:name", handlers.RemoveContainer)
		api.GET("/status", handlers.GetSystemStatus)
	}

	// Trust all proxies (required for Docker/Tailscale)
	r.SetTrustedProxies(nil)

	log.Println("Server starting on 0.0.0.0:8080")
	r.Run("0.0.0.0:8080")
}
