package database

import (
	"database/sql"
	"encoding/json"
	"log"
	"time"

	"docker-scheduler/models"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

// Initialize creates the database connection and tables
func Initialize(dbPath string) error {
	var err error
	DB, err = sql.Open("sqlite3", dbPath)
	if err != nil {
		return err
	}

	// Create schedules table
	createTable := `
	CREATE TABLE IF NOT EXISTS schedules (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		container_name TEXT NOT NULL,
		run_name TEXT DEFAULT '',
		ports TEXT DEFAULT '[]',
		auto_remove BOOLEAN DEFAULT 1,
		days TEXT NOT NULL,
		times TEXT NOT NULL DEFAULT '[]',
		exception_dates TEXT DEFAULT '[]',
		env_vars TEXT DEFAULT '{}',
		active BOOLEAN DEFAULT 1,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);`

	_, err = DB.Exec(createTable)
	if err != nil {
		return err
	}

	// Create Execution Logs table
	createLogsTableSQL := `CREATE TABLE IF NOT EXISTS execution_logs (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		schedule_id INTEGER,
		container_name TEXT,
		docker_name TEXT,
		executed_at TEXT,
		status TEXT DEFAULT 'completed',
		success BOOLEAN,
		output TEXT,
		error TEXT,
		FOREIGN KEY(schedule_id) REFERENCES schedules(id)
	);`
	if _, err := DB.Exec(createLogsTableSQL); err != nil {
		log.Fatal(err)
	}

	// Migrations for existing databases
	// Note: The original migrations for schedules table are kept,
	// and new migrations for execution_logs are added.
	// The 'status' column migration for execution_logs is now handled more robustly.
	DB.Exec("ALTER TABLE schedules ADD COLUMN times TEXT DEFAULT '[]'")
	DB.Exec("ALTER TABLE schedules ADD COLUMN exception_dates TEXT DEFAULT '[]'")
	DB.Exec("ALTER TABLE schedules ADD COLUMN env_vars TEXT DEFAULT '{}'")
	DB.Exec("ALTER TABLE schedules ADD COLUMN run_name TEXT DEFAULT ''")
	DB.Exec("ALTER TABLE schedules ADD COLUMN ports TEXT DEFAULT '[]'")
	DB.Exec("ALTER TABLE schedules ADD COLUMN auto_remove BOOLEAN DEFAULT 1")
	DB.Exec("ALTER TABLE schedules ADD COLUMN random_delay INTEGER DEFAULT 0")

	// Add status column if it doesn't exist (migration)
	_, err = DB.Exec("ALTER TABLE execution_logs ADD COLUMN status TEXT DEFAULT 'completed'")
	if err != nil && err.Error() != "duplicate column name: status" {
		// Ignore error if column exists
		log.Printf("Migration warning (status): %v", err)
	}

	// Add docker_name column if it doesn't exist (migration)
	_, err = DB.Exec("ALTER TABLE execution_logs ADD COLUMN docker_name TEXT DEFAULT ''")
	if err != nil && err.Error() != "duplicate column name: docker_name" {
		log.Printf("Migration warning (docker_name): %v", err)
	}

	// Add scheduled_time column (migration)
	DB.Exec("ALTER TABLE execution_logs ADD COLUMN scheduled_time TEXT DEFAULT ''")

	// Add random_delay column (migration)
	DB.Exec("ALTER TABLE execution_logs ADD COLUMN random_delay INTEGER DEFAULT 0")

	log.Println("Database initialized successfully")
	return nil
}

// parseScheduleRow parses a schedule row into a Schedule struct
func parseScheduleRow(s *models.Schedule, daysJSON, timesJSON, exceptionDatesJSON, envVarsJSON, portsJSON string) {
	json.Unmarshal([]byte(daysJSON), &s.Days)
	json.Unmarshal([]byte(timesJSON), &s.Times)
	json.Unmarshal([]byte(exceptionDatesJSON), &s.ExceptionDates)
	json.Unmarshal([]byte(portsJSON), &s.Ports)
	if s.Times == nil {
		s.Times = []string{}
	}
	if s.ExceptionDates == nil {
		s.ExceptionDates = []string{}
	}
	if s.Ports == nil {
		s.Ports = []string{}
	}
	if envVarsJSON != "" {
		json.Unmarshal([]byte(envVarsJSON), &s.EnvVars)
	}
	if s.EnvVars == nil {
		s.EnvVars = make(map[string]string)
	}
}

// GetAllSchedules returns all schedules from the database
func GetAllSchedules() ([]models.Schedule, error) {
	rows, err := DB.Query(`
		SELECT id, container_name, COALESCE(run_name, ''), COALESCE(ports, '[]'), COALESCE(auto_remove, 1), days, COALESCE(times, '[]'), COALESCE(random_delay, 0), COALESCE(exception_dates, '[]'), COALESCE(env_vars, '{}'), active, created_at 
		FROM schedules 
		ORDER BY created_at DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var schedules []models.Schedule
	for rows.Next() {
		var s models.Schedule
		var daysJSON, timesJSON, exceptionDatesJSON, envVarsJSON, portsJSON string
		err := rows.Scan(&s.ID, &s.ContainerName, &s.RunName, &portsJSON, &s.AutoRemove, &daysJSON, &timesJSON, &s.RandomDelay, &exceptionDatesJSON, &envVarsJSON, &s.Active, &s.CreatedAt)
		if err != nil {
			return nil, err
		}
		parseScheduleRow(&s, daysJSON, timesJSON, exceptionDatesJSON, envVarsJSON, portsJSON)
		schedules = append(schedules, s)
	}

	return schedules, nil
}

// GetActiveSchedules returns only active schedules
func GetActiveSchedules() ([]models.Schedule, error) {
	rows, err := DB.Query(`
		SELECT id, container_name, COALESCE(run_name, ''), COALESCE(ports, '[]'), COALESCE(auto_remove, 1), days, COALESCE(times, '[]'), COALESCE(random_delay, 0), COALESCE(exception_dates, '[]'), COALESCE(env_vars, '{}'), active, created_at 
		FROM schedules 
		WHERE active = 1
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var schedules []models.Schedule
	for rows.Next() {
		var s models.Schedule
		var daysJSON, timesJSON, exceptionDatesJSON, envVarsJSON, portsJSON string
		err := rows.Scan(&s.ID, &s.ContainerName, &s.RunName, &portsJSON, &s.AutoRemove, &daysJSON, &timesJSON, &s.RandomDelay, &exceptionDatesJSON, &envVarsJSON, &s.Active, &s.CreatedAt)
		if err != nil {
			return nil, err
		}
		parseScheduleRow(&s, daysJSON, timesJSON, exceptionDatesJSON, envVarsJSON, portsJSON)
		schedules = append(schedules, s)
	}

	return schedules, nil
}

// CreateSchedule adds a new schedule to the database
func CreateSchedule(containerName string, runName string, ports []string, autoRemove bool, days []string, times []string, exceptionDates []string, envVars map[string]string, randomDelay int) (*models.Schedule, error) {
	daysJSON, err := json.Marshal(days)
	if err != nil {
		return nil, err
	}

	timesJSON, err := json.Marshal(times)
	if err != nil {
		return nil, err
	}

	if exceptionDates == nil {
		exceptionDates = []string{}
	}
	exceptionDatesJSON, err := json.Marshal(exceptionDates)
	if err != nil {
		return nil, err
	}

	if envVars == nil {
		envVars = make(map[string]string)
	}
	envVarsJSON, err := json.Marshal(envVars)
	if err != nil {
		return nil, err
	}

	if ports == nil {
		ports = []string{}
	}
	portsJSON, err := json.Marshal(ports)
	if err != nil {
		return nil, err
	}

	result, err := DB.Exec(
		"INSERT INTO schedules (container_name, run_name, ports, auto_remove, days, times, random_delay, exception_dates, env_vars) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)",
		containerName, runName, string(portsJSON), autoRemove, string(daysJSON), string(timesJSON), randomDelay, string(exceptionDatesJSON), string(envVarsJSON),
	)
	if err != nil {
		return nil, err
	}

	id, _ := result.LastInsertId()
	return &models.Schedule{
		ID:             int(id),
		ContainerName:  containerName,
		RunName:        runName,
		Ports:          ports,
		AutoRemove:     autoRemove,
		Days:           days,
		Times:          times,
		RandomDelay:    randomDelay,
		ExceptionDates: exceptionDates,
		EnvVars:        envVars,
		Active:         true,
	}, nil
}

// DeleteSchedule removes a schedule by ID
func DeleteSchedule(id int) error {
	_, err := DB.Exec("DELETE FROM schedules WHERE id = ?", id)
	return err
}

// ToggleSchedule toggles the active status of a schedule
func ToggleSchedule(id int) error {
	_, err := DB.Exec("UPDATE schedules SET active = NOT active WHERE id = ?", id)
	return err
}

// GetScheduleByID returns a single schedule by ID
func GetScheduleByID(id int) (*models.Schedule, error) {
	var s models.Schedule
	var daysJSON, timesJSON, exceptionDatesJSON, envVarsJSON, portsJSON string
	err := DB.QueryRow(
		"SELECT id, container_name, COALESCE(run_name, ''), COALESCE(ports, '[]'), COALESCE(auto_remove, 1), days, COALESCE(times, '[]'), COALESCE(random_delay, 0), COALESCE(exception_dates, '[]'), COALESCE(env_vars, '{}'), active, created_at FROM schedules WHERE id = ?",
		id,
	).Scan(&s.ID, &s.ContainerName, &s.RunName, &portsJSON, &s.AutoRemove, &daysJSON, &timesJSON, &s.RandomDelay, &exceptionDatesJSON, &envVarsJSON, &s.Active, &s.CreatedAt)
	if err != nil {
		return nil, err
	}
	parseScheduleRow(&s, daysJSON, timesJSON, exceptionDatesJSON, envVarsJSON, portsJSON)
	return &s, nil
}

// AddExceptionDate adds an exception date to a schedule
func AddExceptionDate(id int, date string) error {
	schedule, err := GetScheduleByID(id)
	if err != nil {
		return err
	}

	// Check if date already exists
	for _, d := range schedule.ExceptionDates {
		if d == date {
			return nil // Already exists
		}
	}

	schedule.ExceptionDates = append(schedule.ExceptionDates, date)
	exceptionDatesJSON, _ := json.Marshal(schedule.ExceptionDates)

	_, err = DB.Exec("UPDATE schedules SET exception_dates = ? WHERE id = ?", string(exceptionDatesJSON), id)
	return err
}

// RemoveExceptionDate removes an exception date from a schedule
func RemoveExceptionDate(id int, date string) error {
	schedule, err := GetScheduleByID(id)
	if err != nil {
		return err
	}

	// Filter out the date
	newDates := []string{}
	for _, d := range schedule.ExceptionDates {
		if d != date {
			newDates = append(newDates, d)
		}
	}

	exceptionDatesJSON, _ := json.Marshal(newDates)
	_, err = DB.Exec("UPDATE schedules SET exception_dates = ? WHERE id = ?", string(exceptionDatesJSON), id)
	return err
}

// CreateExecutionLog creates a new log entry with "running" status
func CreateExecutionLog(scheduleID int, containerName string, dockerName string, scheduledTime string, randomDelay int) (int, error) {
	statement, err := DB.Prepare("INSERT INTO execution_logs (schedule_id, container_name, docker_name, scheduled_time, random_delay, executed_at, status, success, output, error) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		return 0, err
	}
	executedAt := time.Now().Format("2006-01-02 15:04:05")
	res, err := statement.Exec(scheduleID, containerName, dockerName, scheduledTime, randomDelay, executedAt, "running", false, "", "")
	if err != nil {
		return 0, err
	}
	id, err := res.LastInsertId()
	return int(id), err
}

// UpdateExecutionLog updates an existing execution log with final status and output
func UpdateExecutionLog(id int, success bool, output string, errorMsg string) error {
	status := "success"
	if !success {
		status = "error"
	}

	// Truncate output to 1000 chars
	if len(output) > 1000 {
		output = output[:1000] + "...[truncated]"
	}
	if len(errorMsg) > 500 {
		errorMsg = errorMsg[:500] + "...[truncated]"
	}

	_, err := DB.Exec(
		"UPDATE execution_logs SET status = ?, success = ?, output = ?, error = ? WHERE id = ?",
		status, success, output, errorMsg, id,
	)
	return err
}

// LogExecution records a container execution in the database (Legacy wrapper)
func LogExecution(scheduleID int, containerName string, success bool, output string, errorMsg string) error {
	id, err := CreateExecutionLog(scheduleID, containerName, "", "", 0)
	if err != nil {
		return err
	}
	return UpdateExecutionLog(id, success, output, errorMsg)
}

// GetExecutionHistory returns the latest execution logs
func GetExecutionHistory(limit int) ([]models.ExecutionLog, error) {
	if limit <= 0 {
		limit = 50
	}

	rows, err := DB.Query("SELECT id, COALESCE(schedule_id, 0), container_name, COALESCE(docker_name, ''), COALESCE(scheduled_time, ''), COALESCE(random_delay, 0), executed_at, COALESCE(status, 'completed'), success, COALESCE(output, ''), COALESCE(error, '') FROM execution_logs ORDER BY executed_at DESC LIMIT ?", limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var logs []models.ExecutionLog
	for rows.Next() {
		var log models.ExecutionLog
		if err := rows.Scan(&log.ID, &log.ScheduleID, &log.ContainerName, &log.DockerName, &log.ScheduledTime, &log.RandomDelay, &log.ExecutedAt, &log.Status, &log.Success, &log.Output, &log.Error); err != nil {
			return nil, err
		}
		logs = append(logs, log)
	}

	return logs, nil
}

// ClearExecutionHistory deletes all execution logs
func ClearExecutionHistory() error {
	_, err := DB.Exec("DELETE FROM execution_logs")
	return err
}

// GetRunningExecutions returns all execution logs with status 'running'
func GetRunningExecutions() ([]models.ExecutionLog, error) {
	rows, err := DB.Query("SELECT id, COALESCE(schedule_id, 0), container_name, COALESCE(docker_name, ''), COALESCE(scheduled_time, ''), COALESCE(random_delay, 0), executed_at, COALESCE(status, 'completed'), success, COALESCE(output, ''), COALESCE(error, '') FROM execution_logs WHERE status = 'running'")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var logs []models.ExecutionLog
	for rows.Next() {
		var log models.ExecutionLog
		if err := rows.Scan(&log.ID, &log.ScheduleID, &log.ContainerName, &log.DockerName, &log.ScheduledTime, &log.RandomDelay, &log.ExecutedAt, &log.Status, &log.Success, &log.Output, &log.Error); err != nil {
			return nil, err
		}
		logs = append(logs, log)
	}
	return logs, nil
}
