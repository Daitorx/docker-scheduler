# Docker Scheduler 🐳

A self-hosted web application to schedule Docker containers execution, featuring a modern UI, timezone support, and exception management.

![Docker Scheduler Screenshot](https://raw.githubusercontent.com/username/repo/main/screenshot.png)

## ✨ Features

- 📅 **Flexible Scheduling**: Schedule containers by specific days of the week.
- ⏰ **Multiple Times**: Run the same container multiple times per day.
- 🎲 **Random Delay**: Add natural variability to executions, avoiding predictable patterns. Each scheduled time gets its own independent random offset.
- 🔧 **Environment Variables**: Pass custom env vars to your containers for each run.
- 🚫 **Exception Dates**: Exclude specific dates (holidays, maintenance) from the schedule.
- 🧹 **Auto-Remove**: Optional automatic cleanup of containers after execution.
- 📺 **Live Logs (Auto-Refresh)**: View real-time logs of running containers in a modern terminal-style interface.
- 📦 **Container Management**: View running containers, stop them, and configure custom names and ports.
- 📊 **Execution History**: Track detailed status (Running/Success/Failure), grouped by container, with scheduled time and delay info.
- 🌍 **Internationalization (i18n)**: Switch between **English (EN)** and **Spanish (ES)** via the header selector. (English by default).
- 🎨 **Modern Dark UI**: Clean, minimalist interface with dark mode and mobile responsiveness. All sections collapsed by default for a cleaner view.
- 💾 **Persistent Storage**: SQLite database for reliable data storage.
- 🗺️ **Timezone Aware**: Respects your local timezone for precise execution.

## 🚀 Quick Start

Run the container mounting your Docker socket (required to control other containers):

```bash
docker run -d \
  -p 8080:8080 \
  -v /var/run/docker.sock:/var/run/docker.sock \
  -v scheduler-data:/app/data \
  -e TZ=Europe/Madrid \
  --name docker-scheduler \
  daitorx/docker-scheduler:latest
```

Open your browser at **http://localhost:8080**

### 📱 Remote / Mobile Access

To access the scheduler from another device (like your mobile phone) on the same network:

1. Find your computer/server IP address (e.g., `192.168.1.100`).
2. Open the browser on your mobile device.
3. Go to **`http://<YOUR-IP>:8080`** (e.g., `http://192.168.1.100:8080`).

Ensure your firewall allows traffic on port 8080.

## ⚙️ Configuration

### Volumes

| Path | Description | Access Mode |
|------|-------------|-------------|
| `/var/run/docker.sock` | **Required**. Allows the scheduler to start other containers. | `rw` |
| `/app/data` | **Recommended**. Persists your schedules and configuration. | `rw` |

### Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `TZ` | `UTC` | **Important**. Set this to your local timezone (e.g., `Europe/Madrid`, `America/New_York`) to ensure schedules run at the expected local time. |
| `GIN_MODE` | `release` | Web server mode (`debug` or `release`). |

## 🛠️ How it Works

1. **The Web UI**: You define schedules (Container Name, Cron, Env Vars, Auto-remove, Ports, etc.).
2. **The Scheduler**: Uses an internal cron-like system (gocron) to check every minute if any container needs to run.
3. **Execution**: When a match is found:
   - It checks if today is an **Exception Date**.
   - If not, it executes the container.
   - If **Auto-remove** is enabled, it uses `docker run --rm`.
   - Logs and status (Running/Success/Error) are tracked in real-time.

## 📦 Supported Images

This scheduler can run **any Docker image** available on your host or Docker Hub.
- `alpine`
- `python:script`
- `node:worker`
- Any custom image

## 🤝 Troubleshooting

### "Container not running at the correct time"
Ensure you have set the `TZ` environment variable correctly.
Check the container time: `docker exec docker-scheduler date`

### "Permission denied connecting to Docker"
Make sure the user running the container has permissions to access `/var/run/docker.sock`. ensuring the volume is mounted correctly is usually enough.

---
*Built with Go (Gin), SQLite, and Vanilla JS.*
