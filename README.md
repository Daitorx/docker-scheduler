# 🐳 Docker Scheduler

Aplicación web para programar la ejecución automática de contenedores Docker en horarios y días específicos.

![Go](https://img.shields.io/badge/Go-00ADD8?style=flat&logo=go&logoColor=white)
![Docker](https://img.shields.io/badge/Docker-2496ED?style=flat&logo=docker&logoColor=white)
![SQLite](https://img.shields.io/badge/SQLite-003B57?style=flat&logo=sqlite&logoColor=white)

## ✨ Características

- 📅 **Programación flexible** - Selecciona días específicos de la semana
- ⏰ **Múltiples horarios** - Configura varias horas de ejecución por schedule
- 🔧 **Variables de entorno** - Pasa variables a tus contenedores
- 🚫 **Fechas de excepción** - Excluye fechas específicas (vacaciones, mantenimiento)
- 🎨 **UI minimalista** - Interfaz oscura moderna con agrupación por imagen
- 💾 **Persistencia** - Los schedules se guardan en SQLite

## 📋 Requisitos

- [Go 1.21+](https://go.dev/dl/)
- [Docker](https://docs.docker.com/get-docker/)
- GCC (para SQLite) - en macOS viene con Xcode Command Line Tools

## 🚀 Instalación

```bash
# Clonar el repositorio
git clone <repo-url>
cd docker-scheduler/backend

# Instalar dependencias
go mod tidy

# Ejecutar
go run main.go
```

El servidor se iniciará en **http://localhost:8080**

## 📖 Uso

1. Abre http://localhost:8080 en tu navegador
2. Introduce el nombre de una **imagen Docker** (ej: `nginx`, `redis`, `mi-app:latest`)
3. Selecciona los **días** de ejecución
4. Añade una o más **horas** de ejecución
5. (Opcional) Añade **variables de entorno**
6. Haz clic en **"Crear Schedule"**

### Gestionar Excepciones

Para excluir fechas específicas (ej: festivos):
1. Haz clic en el icono 📅 del schedule
2. Selecciona un rango de fechas
3. Haz clic en "Añadir Fecha(s)"

## 🔌 API Endpoints

| Método | Endpoint | Descripción |
|--------|----------|-------------|
| GET | `/api/schedules` | Listar todos los schedules |
| POST | `/api/schedules` | Crear un nuevo schedule |
| DELETE | `/api/schedules/:id` | Eliminar un schedule |
| PUT | `/api/schedules/:id/toggle` | Activar/desactivar schedule |
| POST | `/api/schedules/:id/exceptions` | Añadir fecha de excepción |
| DELETE | `/api/schedules/:id/exceptions` | Eliminar fecha de excepción |

### Ejemplo: Crear un schedule

```bash
curl -X POST http://localhost:8080/api/schedules \
  -H "Content-Type: application/json" \
  -d '{
    "container_name": "nginx",
    "days": ["monday", "wednesday", "friday"],
    "times": ["09:00", "18:00"],
    "env_vars": {"NODE_ENV": "production"},
    "exception_dates": ["2026-12-25"]
  }'
```

## 📁 Estructura del Proyecto

```
backend/
├── main.go              # Entry point y router
├── go.mod               # Dependencias Go
├── handlers/            # API handlers (CRUD schedules)
├── models/              # Modelos de datos
├── database/            # Operaciones SQLite
├── scheduler/           # Lógica de programación (gocron)
└── static/              # Frontend (HTML/CSS/JS)
    ├── index.html
    ├── css/styles.css
    └── js/app.js
```

## 🐳 Docker

### Construir la imagen

```bash
docker build -t docker-scheduler .
```

### Ejecutar con Docker

```bash
docker run -d \
  -p 8080:8080 \
  -v /var/run/docker.sock:/var/run/docker.sock \
  -v docker-scheduler-data:/app/data \
  --name scheduler \
  docker-scheduler
```

> ⚠️ **Nota**: Se requiere acceso al socket de Docker para ejecutar contenedores.

## 📄 Licencia

MIT
