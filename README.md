# 🐳 Docker Scheduler

Aplicación web para programar la ejecución automática de contenedores Docker en horarios y días específicos.

![Go](https://img.shields.io/badge/Go-00ADD8?style=flat&logo=go&logoColor=white)
![Docker](https://img.shields.io/badge/Docker-2496ED?style=flat&logo=docker&logoColor=white)
![SQLite](https://img.shields.io/badge/SQLite-003B57?style=flat&logo=sqlite&logoColor=white)

## ✨ Características

- 📅 **Programación flexible** - Selecciona días específicos de la semana
- ⏰ **Múltiples horarios** - Configura varias horas de ejecución por schedule
- 🎲 **Retraso aleatorio** - Añade variabilidad natural a las ejecuciones para evitar patrones predecibles
- 🔧 **Variables de entorno** - Pasa variables a tus contenedores
- 🚫 **Fechas de excepción** - Excluye fechas específicas (vacaciones, mantenimiento)
- 📺 **Logs en vivo** - Visualiza logs de contenedores en tiempo real
- 📊 **Historial de ejecuciones** - Registro detallado con hora programada, delay y estado
- 🌍 **Multiidioma** - Interfaz disponible en Español e Inglés
- 🎨 **UI minimalista** - Interfaz oscura moderna con agrupación por imagen
- 💾 **Persistencia** - Los schedules se guardan en SQLite
- 🗺️ **Soporte de Timezone** - Respeta tu zona horaria local

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
6. (Opcional) Añade un **retraso aleatorio** para variabilidad
7. Haz clic en **"Crear Schedule"**

### Gestionar Excepciones

Para excluir fechas específicas (ej: festivos):
1. Haz clic en el icono 📅 del schedule
2. Selecciona un rango de fechas
3. Haz clic en "Añadir Fecha(s)"

### 🎲 Retraso Aleatorio (Random Delay)

Esta funcionalidad permite añadir una variabilidad natural a la ejecución de tus tareas, evitando patrones predecibles.

**¿Cómo funciona?**
Al configurar un schedule, defines un "Retraso Aleatorio" (en minutos). El sistema elegirá **un nuevo minuto al azar** cada día de ejecución, dentro de la ventana especificada.

**Ejemplo Práctico:**
Imagina que configuras una tarea para ejecutarse todos los **Jueves** a las **08:00** con un retraso aleatorio de **30 minutos**.
- **Jueves semana 1**: A medianoche, el sistema lanza los dados y elige las `08:12`. Tu contenedor se ejecuta exactamente a las 08:12.
- **Jueves semana 2**: El sistema vuelve a calcular y elige las `08:27`. Se ejecuta a las 08:27.
- **Jueves semana 3**: El sistema elige las `08:05`. Se ejecuta a las 08:05.

**Características Clave:**
- **Recálculo Diario**: El sistema calcula una nueva hora aleatoria cada día de ejecución programado (a las 00:00). Nunca es la misma.
- **Múltiples Horas Independientes**: Si configuras varias horas en el mismo día (ej. 08:00 y 14:00), **cada una recibe su propio cálculo aleatorio independiente**.
- **Historial Detallado**: En el historial de ejecuciones verás la hora programada original junto con el delay aplicado (ej. `09:00 (+5m)`).

**Ejemplo con Múltiples Horas:**
Configuras un schedule para los **Jueves** a las **08:00** y **14:00** con un retraso de **15 minutos**:
| Hora Base | Jueves 1 | Jueves 2 | Jueves 3 |
|-----------|----------|----------|----------|
| 08:00     | 08:07    | 08:12    | 08:03    |
| 14:00     | 14:11    | 14:02    | 14:09    |

Cada hora se calcula de forma totalmente independiente.

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
