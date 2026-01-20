package scheduler

import (
	"log"
	"os/exec"
	"strings"
	"time"

	"docker-scheduler/database"

	"github.com/go-co-op/gocron"
)

var scheduler *gocron.Scheduler

// dayMap maps day names to time.Weekday
var dayMap = map[string]time.Weekday{
	"sunday":    time.Sunday,
	"monday":    time.Monday,
	"tuesday":   time.Tuesday,
	"wednesday": time.Wednesday,
	"thursday":  time.Thursday,
	"friday":    time.Friday,
	"saturday":  time.Saturday,
}

// Initialize starts the scheduler
func Initialize() {
	scheduler = gocron.NewScheduler(time.Local)
	scheduler.StartAsync()
	log.Println("Scheduler initialized")

	// Load existing schedules
	RefreshSchedules()

	// initial cleanup
	go CleanupZombieExecutions()

	// Schedule periodic cleanup (every 5 minutes)
	scheduler.Every(5).Minutes().Do(CleanupZombieExecutions)
}

// RefreshSchedules reloads all active schedules from the database
func RefreshSchedules() {
	// Clear all jobs
	scheduler.Clear()

	schedules, err := database.GetActiveSchedules()
	if err != nil {
		log.Printf("Error loading schedules: %v", err)
		return
	}

	for _, s := range schedules {
		for _, day := range s.Days {
			weekday, ok := dayMap[strings.ToLower(day)]
			if !ok {
				log.Printf("Invalid day: %s", day)
				continue
			}

			for _, schedTime := range s.Times {
				scheduleID := s.ID
				containerName := s.ContainerName
				runName := s.RunName
				ports := s.Ports
				autoRemove := s.AutoRemove
				envVars := s.EnvVars
				exceptionDates := s.ExceptionDates
				timeStr := schedTime
				job, err := scheduler.Every(1).Week().Weekday(weekday).At(timeStr).Do(func() {
					// Check if today is an exception date
					today := time.Now().Format("2006-01-02")
					for _, exDate := range exceptionDates {
						if exDate == today {
							log.Printf("Skipping container %s - today (%s) is an exception date", containerName, today)
							return
						}
					}
					runContainer(scheduleID, containerName, runName, ports, autoRemove, envVars)
				})

				if err != nil {
					log.Printf("Error scheduling job: %v", err)
				} else {
					log.Printf("Scheduled %s for %s (%s). Next run: %v", containerName, day, timeStr, job.NextRun())
				}
			}
		}
	}

	log.Printf("Loaded %d schedules", len(schedules))
}

// runContainer runs a Docker container from an image with environment variables
func runContainer(scheduleID int, containerName string, runName string, ports []string, autoRemove bool, envVars map[string]string) {
	log.Printf("Running container: %s (autoRemove: %v)", containerName, autoRemove)

	// If container has a name, stop and remove existing container first
	if runName != "" {
		exec.Command("docker", "stop", runName).Run()
		exec.Command("docker", "rm", runName).Run()
	}

	// Build docker run command
	var args []string
	if autoRemove {
		// Run synchronously with --rm to capture output
		args = []string{"run", "--rm"}
	} else {
		// Run detached without --rm
		args = []string{"run", "-d"}
	}

	// Determine the actual Docker container name to use
	actualDockerName := runName
	if actualDockerName == "" {
		// Generate a unique name for this run to avoid conflicts and allow log retrieval
		// Format: sched-<scheduleID>-<timestamp>
		actualDockerName = "sched-" + strings.ReplaceAll(containerName, " ", "_") + "-" + time.Now().Format("20060102150405")
	}

	// Always use --name so we can track it
	args = append(args, "--name", actualDockerName)

	// Add port mappings
	for _, port := range ports {
		args = append(args, "-p", port)
	}

	// Add environment variables
	for key, value := range envVars {
		args = append(args, "-e", key+"="+value)
	}

	args = append(args, containerName)

	// Create initial log entry with "running" status and the ACTUAL docker name
	logID, logErr := database.CreateExecutionLog(scheduleID, containerName, actualDockerName)
	if logErr != nil {
		log.Printf("Failed to create execution log: %v", logErr)
	}

	cmd := exec.Command("docker", args...)
	output, err := cmd.CombinedOutput()

	// Log the execution
	success := err == nil
	errorMsg := ""
	if err != nil {
		errorMsg = err.Error()
		log.Printf("Error running container %s: %v - %s", containerName, err, string(output))
	} else {
		log.Printf("Container %s executed successfully. Output: %s", containerName, string(output))
	}

	// Update log with final status and output
	if logID != 0 {
		if updateErr := database.UpdateExecutionLog(logID, success, string(output), errorMsg); updateErr != nil {
			log.Printf("Failed to update execution log: %v", updateErr)
		}
	}
}

// Stop stops the scheduler
func Stop() {
	if scheduler != nil {
		scheduler.Stop()
	}
}

// CleanupZombieExecutions checks for executions marked as 'running' that are not actually running in Docker
func CleanupZombieExecutions() {
	log.Println("Checking for zombie executions...")
	runningLogs, err := database.GetRunningExecutions()
	if err != nil {
		log.Printf("Error checking for zombie executions: %v", err)
		return
	}

	if len(runningLogs) == 0 {
		return
	}

	// Get currently running containers
	cmd := exec.Command("docker", "ps", "--format", "{{.Names}}")
	output, err := cmd.Output()
	if err != nil {
		log.Printf("Error getting running containers: %v", err)
		return
	}

	runningContainers := strings.Split(string(output), "\n")
	runningMap := make(map[string]bool)
	for _, name := range runningContainers {
		name = strings.TrimSpace(name)
		if name != "" {
			runningMap[name] = true
		}
	}

	// Check each log
	for _, execLog := range runningLogs {
		nameToCheck := execLog.DockerName
		if nameToCheck == "" {
			// Fallback: mostly creates logs with DockerName, but just in case
			nameToCheck = execLog.ContainerName
		}

		if nameToCheck != "" && !runningMap[nameToCheck] {
			log.Printf("Found zombie execution: ID %d, Container %s (Docker Name: %s). Marking as error.", execLog.ID, execLog.ContainerName, nameToCheck)
			database.UpdateExecutionLog(execLog.ID, false, "Container execution interrupted (Zombie detected)", "Container execution interrupted (Zombie detected)")
		}
	}
}
