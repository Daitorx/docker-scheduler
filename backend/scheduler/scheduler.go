package scheduler

import (
	"fmt"
	"log"
	"math/rand"
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

	// Re-add System Jobs (wiped by Clear)
	scheduler.Every(5).Minutes().Do(CleanupZombieExecutions)

	schedules, err := database.GetActiveSchedules()
	if err != nil {
		log.Printf("Error loading schedules: %v", err)
		return
	}

	for _, s := range schedules {
		for _, day := range s.Days {
			weekday, ok := dayMap[strings.ToLower(day)]
			if !ok {
				continue
			}

			// Process each configured time for this day
			for _, schedTime := range s.Times {
				scheduleID := s.ID
				containerName := s.ContainerName
				runName := s.RunName
				ports := s.Ports
				autoRemove := s.AutoRemove
				envVars := s.EnvVars
				exceptionDates := s.ExceptionDates
				randomDelay := s.RandomDelay

				// If no random delay, schedule normally
				if randomDelay <= 0 {
					scheduleJob(scheduleID, containerName, runName, ports, autoRemove, envVars, exceptionDates, weekday, schedTime, 0)
					continue
				}

				// If Random Delay is set: random time between `schedTime` and `schedTime + delay` minutes
				// 1. Weekly Trigger at Midnight: Pick a random time for THAT day
				scheduler.Every(1).Week().Weekday(weekday).At("00:00").Tag(fmt.Sprintf("sched-gen-%d", scheduleID)).Do(func(sID int, cName, rName string, pts []string, aRem bool, eVars map[string]string, eDates []string, baseTime string, delay int) {
					finalTime := calculateRandomTimeWithDelay(baseTime, delay)
					scheduleOneOff(sID, cName, rName, pts, aRem, eVars, eDates, baseTime, delay, finalTime)
				}, scheduleID, containerName, runName, ports, autoRemove, envVars, exceptionDates, schedTime, randomDelay)

				// 2. Startup Check: If today is the scheduled day
				if time.Now().Weekday() == weekday {
					scheduleRandomDelayOnStartup(scheduleID, containerName, runName, ports, autoRemove, envVars, exceptionDates, schedTime, randomDelay)
				}
			}
		}
	}

	log.Printf("Loaded %d schedules", len(schedules))
}

func calculateRandomTimeWithDelay(baseTimeStr string, delayMinutes int) string {
	baseTime, err := time.Parse("15:04", baseTimeStr)
	if err != nil {
		return baseTimeStr
	}

	rand.Seed(time.Now().UnixNano())
	randomAdd := rand.Intn(delayMinutes + 1) // +1 because we want inclusive [0, delay]
	duration := time.Duration(randomAdd) * time.Minute

	finalTime := baseTime.Add(duration)
	return finalTime.Format("15:04")
}

func scheduleRandomDelayOnStartup(scheduleID int, containerName, runName string, ports []string, autoRemove bool, envVars map[string]string, exceptionDates []string, baseTimeStr string, delayMinutes int) {
	now := time.Now()
	base, _ := time.Parse("15:04", baseTimeStr)

	// Base time today
	start := time.Date(now.Year(), now.Month(), now.Day(), base.Hour(), base.Minute(), 0, 0, now.Location())
	// Max time today
	maxEnd := start.Add(time.Duration(delayMinutes) * time.Minute)

	// If the entire window is in the past, skip
	if maxEnd.Before(now) {
		return // Too late
	}

	// Recalculate window to be [max(now, start), maxEnd].
	effectiveStart := start
	if now.After(start) {
		effectiveStart = now
	}

	windowDuration := maxEnd.Sub(effectiveStart)
	if windowDuration <= 0 {
		return // No time left in window
	}

	rand.Seed(time.Now().UnixNano())
	randomOffset := time.Duration(rand.Int63n(int64(windowDuration)))
	executionTime := effectiveStart.Add(randomOffset)
	executionTimeStr := executionTime.Format("15:04")

	scheduleOneOff(scheduleID, containerName, runName, ports, autoRemove, envVars, exceptionDates, baseTimeStr, delayMinutes, executionTimeStr)
}

func scheduleOneOff(scheduleID int, containerName, runName string, ports []string, autoRemove bool, envVars map[string]string, exceptionDates []string, scheduledTime string, randomDelay int, timeStr string) {
	now := time.Now()
	parsed, _ := time.Parse("15:04", timeStr)
	scheduled := time.Date(now.Year(), now.Month(), now.Day(), parsed.Hour(), parsed.Minute(), 0, 0, now.Location())

	if scheduled.Before(now) {
		return // In the past, skip
	}

	tag := fmt.Sprintf("sched-exec-%d", scheduleID)
	var job *gocron.Job
	var err error

	job, err = scheduler.Every(1).Day().At(timeStr).Tag(tag).Do(func() {
		// Check exceptions
		today := time.Now().Format("2006-01-02")
		for _, exDate := range exceptionDates {
			if exDate == today {
				return
			}
		}

		runContainer(scheduleID, containerName, runName, ports, autoRemove, envVars, scheduledTime, randomDelay)
	})

	if err == nil {
		job.LimitRunsTo(1)
	}
}

func scheduleJob(scheduleID int, containerName, runName string, ports []string, autoRemove bool, envVars map[string]string, exceptionDates []string, weekday time.Weekday, timeStr string, randomDelay int) {
	tag := fmt.Sprintf("sched-exec-%d", scheduleID)
	scheduler.Every(1).Week().Weekday(weekday).At(timeStr).Tag(tag).Do(func() {
		// Check if today is an exception date
		today := time.Now().Format("2006-01-02")
		for _, exDate := range exceptionDates {
			if exDate == today {
				return
			}
		}
		runContainer(scheduleID, containerName, runName, ports, autoRemove, envVars, timeStr, randomDelay)
	})
}

// runContainer runs a Docker container from an image with environment variables
func runContainer(scheduleID int, containerName string, runName string, ports []string, autoRemove bool, envVars map[string]string, scheduledTime string, randomDelay int) {
	log.Printf("Running container: %s", containerName)

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

	// Create initial log entry with scheduled time and delay
	logID, logErr := database.CreateExecutionLog(scheduleID, containerName, actualDockerName, scheduledTime, randomDelay)
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
		log.Printf("Error running container %s: %v", containerName, err)
	}

	// Update log with final status and output
	if logID != 0 {
		database.UpdateExecutionLog(logID, success, string(output), errorMsg)
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
	runningLogs, err := database.GetRunningExecutions()
	if err != nil {
		return
	}

	if len(runningLogs) == 0 {
		return
	}

	// Get currently running containers
	cmd := exec.Command("docker", "ps", "--format", "{{.Names}}")
	output, err := cmd.Output()
	if err != nil {
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
			nameToCheck = execLog.ContainerName
		}

		if nameToCheck != "" && !runningMap[nameToCheck] {
			database.UpdateExecutionLog(execLog.ID, false, "Container execution interrupted", "Container execution interrupted")
		}
	}
}
