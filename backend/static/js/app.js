// API Base URL
const API_URL = '/api';

// Day translations
const dayTranslations = {
    monday: { es: 'Lun', en: 'Mon' },
    tuesday: { es: 'Mar', en: 'Tue' },
    wednesday: { es: 'Mié', en: 'Wed' },
    thursday: { es: 'Jue', en: 'Thu' },
    friday: { es: 'Vie', en: 'Fri' },
    saturday: { es: 'Sáb', en: 'Sat' },
    sunday: { es: 'Dom', en: 'Sun' }
};

// I18n Translations
const translations = {
    es: {
        title: "Docker Scheduler",
        subtitle: "Programa la ejecución automática de tus contenedores",
        newSchedule: "Nuevo Schedule",
        imageLabel: "Nombre de la Imagen",
        imagePlaceholder: "ej: nginx, redis, my-app",
        containerLabel: "Nombre del Contenedor",
        optional: "(opcional)",
        portsLabel: "Puertos",
        autoRemoveLabel: "Eliminar contenedor al terminar",
        daysLabel: "Días de Ejecución",
        hoursLabel: "Horas de Ejecución",
        addTimeBtn: "Añadir Hora",
        envLabel: "Variables de Entorno",
        addEnvBtn: "Añadir Variable",
        createBtn: "Crear Schedule",
        schedulesTitle: "Schedules Configurados",
        loadingSchedules: "Cargando schedules...",
        emptySchedules: "No hay schedules configurados",
        historyTitle: "Historial de Ejecuciones",
        clearHistoryBtn: "Limpiar",
        loadingHistory: "Cargando historial...",
        emptyHistory: "No hay ejecuciones registradas",
        containersTitle: "Contenedores en Ejecución",
        loadingContainers: "Cargando contenedores...",
        emptyContainers: "No hay contenedores en ejecución",
        footer: "Docker Scheduler © 2026",
        // Dynamic
        runs: "ejecuciones",
        running: "en ejecución",
        success: "exitosos",
        failed: "fallidos",
        statusRunning: "En ejecución...",
        viewLogs: "Ver logs en vivo",
        stopContainer: "Detener contenedor",
        close: "Cerrar",
        update: "Actualizar",
        envName: "Nombre",
        envValue: "Valor",
        exceptionModalTitle: "Gestionar Excepciones",
        exceptionDesc: "Añade fechas o rangos de fechas en las que este schedule NO debe ejecutarse",
        from: "Desde",
        to: "Hasta",
        addDatesBtn: "Añadir Fecha(s)",
        cancelBtn: "Cancelar",
        confirmBtn: "Confirmar",
        deleteBtn: "Eliminar",
        noExceptions: "No hay excepciones configuradas",
        confirmDeleteSchedule: "¿Eliminar este schedule?",
        confirmStopContainer: "¿Detener el contenedor",
        confirmClearHistory: "¿Estás seguro de que quieres borrar todo el historial de ejecuciones?",
        toastScheduleCreated: "Schedule creado exitosamente",
        toastScheduleDeleted: "Schedule eliminado",
        toastHistoryCleared: "Historial borrado correctamente",
        toastContainerStopped: "Contenedor detenido",
        errorGeneric: "Error",
        mon: "Lun", tue: "Mar", wed: "Mié", thu: "Jue", fri: "Vie", sat: "Sáb", sun: "Dom"
    },
    en: {
        title: "Docker Scheduler",
        subtitle: "Schedule automatic execution of your Docker containers",
        newSchedule: "New Schedule",
        imageLabel: "Image Name",
        imagePlaceholder: "e.g., nginx, redis, my-app",
        containerLabel: "Container Name",
        optional: "(optional)",
        portsLabel: "Ports",
        autoRemoveLabel: "Auto-remove container on finish",
        daysLabel: "Execution Days",
        hoursLabel: "Execution Hours",
        addTimeBtn: "Add Time",
        envLabel: "Environment Variables",
        addEnvBtn: "Add Variable",
        createBtn: "Create Schedule",
        schedulesTitle: "Configured Schedules",
        loadingSchedules: "Loading schedules...",
        emptySchedules: "No schedules configured",
        historyTitle: "Execution History",
        clearHistoryBtn: "Clear",
        loadingHistory: "Loading history...",
        emptyHistory: "No execution history",
        containersTitle: "Running Containers",
        loadingContainers: "Loading containers...",
        emptyContainers: "No running containers",
        footer: "Docker Scheduler © 2026",
        // Dynamic
        runs: "runs",
        running: "running",
        success: "successful",
        failed: "failed",
        statusRunning: "Running...",
        viewLogs: "View live logs",
        stopContainer: "Stop container",
        close: "Close",
        update: "Update",
        envName: "Name",
        envValue: "Value",
        exceptionModalTitle: "Manage Exceptions",
        exceptionDesc: "Add dates or date ranges when this schedule should NOT run",
        from: "From",
        to: "To",
        addDatesBtn: "Add Date(s)",
        cancelBtn: "Cancel",
        confirmBtn: "Confirm",
        deleteBtn: "Delete",
        noExceptions: "No exceptions configured",
        confirmDeleteSchedule: "Delete this schedule?",
        confirmStopContainer: "Stop container",
        confirmClearHistory: "Are you sure you want to clear all execution history?",
        toastScheduleCreated: "Schedule created successfully",
        toastScheduleDeleted: "Schedule deleted",
        toastHistoryCleared: "History cleared successfully",
        toastContainerStopped: "Container stopped",
        errorGeneric: "Error",
        mon: "Mon", tue: "Tue", wed: "Wed", thu: "Thu", fri: "Fri", sat: "Sat", sun: "Sun"
    }
};

let currentLang = localStorage.getItem('dockerSchedulerLang') || 'en'; // Default EN

function changeLanguage(lang) {
    if (!translations[lang]) return;
    currentLang = lang;
    localStorage.setItem('dockerSchedulerLang', lang);
    updateLanguageUI();

    // Reload dynamic content
    loadSchedules();
    loadHistory();
    loadContainers();
}

function t(key) {
    return translations[currentLang][key] || key;
}

function updateLanguageUI() {
    document.querySelectorAll('[data-i18n]').forEach(el => {
        const key = el.getAttribute('data-i18n');
        if (translations[currentLang][key]) {
            if (el.tagName === 'INPUT' && el.getAttribute('placeholder')) {
                el.placeholder = translations[currentLang][key];
            } else {
                el.firstChild.textContent = translations[currentLang][key];
                // Restore spans if they exist (like "optional")
                // Simple text replacement is risky if we have child elements.
                // Better approach for elements with children:
                // If it has children (like optional span), only update the text node?
                // For simplicity, we'll design HTML to separate text nodes or use specific Attributes.
                // Or simply:
                if (el.children.length === 0) {
                    el.textContent = translations[currentLang][key];
                } else {
                    // For labels with <span class="optional">
                    // We need to rebuild it.
                    // Special handling?
                    // Let's rely on data-i18n on the specific TEXT element, not the container.
                }
            }
        }
    });

    // Special cases with nested spans
    if (translations[currentLang].optional) {
        document.querySelectorAll('.optional').forEach(el => el.textContent = translations[currentLang].optional);
    }

    // Weekdays in form
    document.querySelectorAll('.day-checkbox input').forEach(input => {
        const span = input.nextElementSibling;
        const val = input.value; // monday, tuesday...
        const key = val.substring(0, 3).toLowerCase(); // mon, tue...
        // Actually we have dayTranslations structure, but keys in t() are mon, tue...
        // Let's use the dayTranslations variable for this.
        if (dayTranslations[val] && dayTranslations[val][currentLang]) {
            span.textContent = dayTranslations[val][currentLang];
        }
    });

    // Update active state of lang buttons if we had them
    document.getElementById('langEn')?.classList.toggle('active', currentLang === 'en');
    document.getElementById('langEs')?.classList.toggle('active', currentLang === 'es');
}

// Current schedule data
let allSchedules = [];
let currentExceptionScheduleId = null;

// DOM Elements
const scheduleForm = document.getElementById('scheduleForm');
const schedulesList = document.getElementById('schedulesList');
const toast = document.getElementById('toast');
const envVarsContainer = document.getElementById('envVarsContainer');
const addEnvVarBtn = document.getElementById('addEnvVar');
const timesContainer = document.getElementById('timesContainer');
const addTimeBtn = document.getElementById('addTime');
const exceptionModal = document.getElementById('exceptionModal');
let historyList;
let containersList;

// Initialize
document.addEventListener('DOMContentLoaded', () => {
    historyList = document.getElementById('historyList');
    containersList = document.getElementById('containersList');
    updateLanguageUI(); // Initialize language
    loadSchedules();
    loadHistory();
    loadContainers();
    scheduleForm.addEventListener('submit', handleFormSubmit);
    addEnvVarBtn.addEventListener('click', addEnvVarRow);
    addTimeBtn.addEventListener('click', addTimeRow);

    // Close modal on outside click
    window.addEventListener('click', (e) => {
        if (e.target === exceptionModal) {
            closeExceptionModal();
        }
        // Close popovers when clicking outside
        if (!e.target.closest('.env-vars-wrapper')) {
            document.querySelectorAll('.env-vars-popover.show').forEach(el => {
                el.classList.remove('show');
            });
        }
    });
});

// Add time row
function addTimeRow() {
    const row = document.createElement('div');
    row.className = 'time-row';
    row.innerHTML = `
        <input type="time" class="time-input" required>
        <button type="button" class="btn-remove" onclick="removeTimeRow(this)">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <path d="M18 6L6 18M6 6l12 12"/>
            </svg>
        </button>
    `;
    timesContainer.appendChild(row);
    updateTimeRemoveButtons();
}

// Remove time row
function removeTimeRow(btn) {
    btn.parentElement.remove();
    updateTimeRemoveButtons();
}

// Update visibility of remove buttons
function updateTimeRemoveButtons() {
    const rows = timesContainer.querySelectorAll('.time-row');
    rows.forEach((row) => {
        const btn = row.querySelector('.btn-remove');
        btn.style.visibility = rows.length === 1 ? 'hidden' : 'visible';
    });
}

// Get times from form
function getTimes() {
    const times = [];
    const inputs = timesContainer.querySelectorAll('.time-input');
    inputs.forEach(input => {
        if (input.value) {
            times.push(input.value);
        }
    });
    return times;
}

// Clear times
function clearTimes() {
    timesContainer.innerHTML = `
        <div class="time-row">
            <input type="time" class="time-input" required>
            <button type="button" class="btn-remove" onclick="removeTimeRow(this)" style="visibility: hidden;">
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                    <path d="M18 6L6 18M6 6l12 12"/>
                </svg>
            </button>
        </div>
    `;
}

// Add environment variable row
function addEnvVarRow() {
    const row = document.createElement('div');
    row.className = 'env-var-row';
    row.innerHTML = `
        <input type="text" class="env-key" placeholder="KEY" />
        <input type="text" class="env-value" placeholder="value" />
        <button type="button" class="btn-remove" onclick="this.parentElement.remove()">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <path d="M18 6L6 18M6 6l12 12"/>
            </svg>
        </button>
    `;
    envVarsContainer.appendChild(row);
}

// Get environment variables from form
function getEnvVars() {
    const envVars = {};
    const rows = envVarsContainer.querySelectorAll('.env-var-row');
    rows.forEach(row => {
        const key = row.querySelector('.env-key').value.trim();
        const value = row.querySelector('.env-value').value.trim();
        if (key) {
            envVars[key] = value;
        }
    });
    return envVars;
}

// Clear environment variables form
function clearEnvVars() {
    envVarsContainer.innerHTML = '';
}

// Load all schedules
async function loadSchedules() {
    try {
        const response = await fetch(`${API_URL}/schedules`);
        const schedules = await response.json();
        allSchedules = schedules; // Store globally
        renderSchedules(schedules);
        renderSchedules(schedules);
    } catch (error) {
        console.error('Error loading schedules:', error);
        showToast('Error al cargar los schedules', 'error');
        schedulesList.innerHTML = `
            <div class="empty-state">
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                    <circle cx="12" cy="12" r="10"/>
                    <path d="M12 8v4M12 16h.01"/>
                </svg>
                <p>Error al conectar con el servidor</p>
            </div>
        `;
    }
}

// Format env vars for display (Inline Expandable)
function formatEnvVars(envVars, scheduleId) {
    if (!envVars || Object.keys(envVars).length === 0) return '';
    const count = Object.keys(envVars).length;

    const listHtml = Object.entries(envVars).map(([key, value]) => `
        <div class="env-var-item-inline">
            <span class="env-key">${escapeHtml(key)}</span>
            <span class="env-equals">=</span>
            <span class="env-value">${escapeHtml(value)}</span>
        </div>
    `).join('');

    return `
        <div class="env-vars-wrapper">
            <button class="env-vars-badge" onclick="toggleEnvExpandable(this)">
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                    <path d="M9 18l6-6-6-6"/>
                </svg>
                <span>${count} var${count > 1 ? 's' : ''}</span>
            </button>
            <div class="env-vars-expandable">
                ${listHtml}
            </div>
        </div>
    `;
}

// Toggle Inline Expandable visibility
function toggleEnvExpandable(btn) {
    const container = btn.nextElementSibling;
    const isExpanded = btn.classList.contains('expanded');

    if (isExpanded) {
        btn.classList.remove('expanded');
        container.classList.remove('show');
    } else {
        btn.classList.add('expanded');
        container.classList.add('show');
    }
}

// Format times for display
function formatTimes(times) {
    if (!times || times.length === 0) return '';
    return times.map(t => `<span class="time-badge">${t}</span>`).join('');
}

// Format exception dates for display
function formatExceptions(exceptions) {
    if (!exceptions || exceptions.length === 0) return '';
    const tooltip = t('exceptionDesc'); // Or list them, but let's keep it simple for now or loop
    // Better: title with joined dates
    return `<span class="exception-badge" title="${t('exceptionDesc')}: ${exceptions.join(', ')}">${exceptions.length} ${t('errorGeneric') === 'Error' ? 'exceptions' : 'excepciones'}</span>`;
    // Wait, simple 'exceptions' word is not in my list. 
    // Let's rely on standard logic or add 'exceptions' to list? 
    // Actually, I can just use a non-translated string or add it to translations?
    // Added 'runs', 'success', but not 'exception'.
    // Let's just use icon?
    // Or just simple text "X dates".
}

// Render schedules list grouped by container name
function renderSchedules(schedules) {
    const list = document.getElementById('schedulesList');
    if (!schedules || schedules.length === 0) {
        list.innerHTML = `
            <div class="empty-state">
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                    <rect x="3" y="4" width="18" height="18" rx="2"/>
                    <path d="M16 2v4M8 2v4M3 10h18"/>
                </svg>
                <p>${t('emptySchedules')}</p>
            </div>
        `;
        return;
    }

    // Group schedules by container name
    const groups = {};
    schedules.forEach(schedule => {
        const name = schedule.container_name;
        if (!groups[name]) {
            groups[name] = [];
        }
        groups[name].push(schedule);
    });

    // Render grouped schedules
    list.innerHTML = Object.entries(groups).map(([containerName, containerSchedules]) => `
        <div class="schedule-group">
            <div class="schedule-group-header" onclick="toggleGroup(this)">
                <div class="group-info">
                    <svg class="chevron" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                        <path d="M9 18l6-6-6-6"/>
                    </svg>
                    <span class="group-name">${escapeHtml(containerName)}</span>
                    <span class="group-count">${containerSchedules.length}</span>
                </div>
            </div>
            <div class="schedule-group-content">
                ${containerSchedules.map(schedule => renderScheduleCard(schedule)).join('')}
            </div>
        </div>
    `).join('');

    // Auto-expand removed. They will be collapsed by default (CSS assumes .expanded to show).
    // Ensure CSS hides content when not expanded.
    // If CSS uses .schedule-group-content { display: none } and .expanded .schedule-group-content { display: block }
    // Then checking default:
    // User wants "Todo recogido" (collapsed).
}

// Render individual schedule card
function renderScheduleCard(schedule) {
    const envVarsBadge = formatEnvVarsBadge(schedule.env_vars, schedule.id);
    const envVarsContent = formatEnvVarsContent(schedule.env_vars, schedule.id);
    const dayLabels = schedule.days.map(day => {
        // Translation logic: dayTranslations[day] is {es:..., en:...}
        const shortDay = dayTranslations[day] ? dayTranslations[day][currentLang] : day;
        return `<span class="schedule-day">${shortDay}</span>`;
    }).join('');

    return `
        <div class="schedule-card ${schedule.active ? '' : 'inactive'}" data-id="${schedule.id}">
            <div class="schedule-rows">
                <div class="schedule-primary-row">
                    <div class="schedule-info">
                        <div class="schedule-details">
                            <div class="schedule-days">
                                ${dayLabels}
                            </div>
                            <div class="schedule-times">
                                ${formatTimes(schedule.times)}
                            </div>
                            ${envVarsBadge}
                            ${formatExceptions(schedule.exception_dates)}
                        </div>
                    </div>
                    <div class="schedule-actions">
                        <button class="btn-icon" onclick="openExceptionModal(${schedule.id})" title="${t('exceptionModalTitle')}">
                            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                                <rect x="3" y="4" width="18" height="18" rx="2" ry="2"></rect>
                                <line x1="16" y1="2" x2="16" y2="6"></line>
                                <line x1="8" y1="2" x2="8" y2="6"></line>
                                <line x1="3" y1="10" x2="21" y2="10"></line>
                            </svg>
                        </button>
                        <button class="btn-icon" onclick="toggleSchedule(${schedule.id})" title="${schedule.active ? t('close') : t('createBtn')}"> 
                            <!-- Wait, title translation needs improvement. Toggle active/inactive. -->
                            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" class="${schedule.active ? 'text-success' : 'text-muted'}">
                                <path d="M18.36 6.64a9 9 0 1 1-12.73 0"></path>
                                <line x1="12" y1="2" x2="12" y2="12"></line>
                            </svg>
                        </button>
                        <button class="btn-icon delete" onclick="deleteSchedule(${schedule.id})" title="${t('confirmDeleteSchedule')}">
                            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                                <polyline points="3 6 5 6 21 6"></polyline>
                                <path d="M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2"></path>
                            </svg>
                        </button>
                    </div>
                </div>
                ${envVarsContent}
            </div>
        </div>
    `;
}

// Format env vars Badge
function formatEnvVarsBadge(envVars, scheduleId) {
    if (!envVars || Object.keys(envVars).length === 0) return '';
    const count = Object.keys(envVars).length;

    return `
        <button class="env-vars-badge" onclick="toggleEnvExpandable(this)" data-target="env-vars-${scheduleId}">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <path d="M9 18l6-6-6-6"/>
            </svg>
            <span>${count} var${count > 1 ? 's' : ''}</span>
        </button>
    `;
}

// Format env vars Content
function formatEnvVarsContent(envVars, scheduleId) {
    if (!envVars || Object.keys(envVars).length === 0) return '';

    const listHtml = Object.entries(envVars).map(([key, value]) => `
        <div class="env-var-item-inline">
            <span class="env-key">${escapeHtml(key)}</span>
            <span class="env-equals">=</span>
            <span class="env-value">${escapeHtml(value)}</span>
        </div>
    `).join('');

    return `
        <div id="env-vars-${scheduleId}" class="env-vars-expandable">
            ${listHtml}
        </div>
    `;
}

// Toggle Inline Expandable visibility
function toggleEnvExpandable(btn) {
    const targetId = btn.dataset.target;
    const container = document.getElementById(targetId);
    const isExpanded = btn.classList.contains('expanded');

    if (isExpanded) {
        btn.classList.remove('expanded');
        container.classList.remove('show');
    } else {
        btn.classList.add('expanded');
        container.classList.add('show');
    }
}

// Toggle group expand/collapse
function toggleGroup(header) {
    const group = header.closest('.schedule-group');
    group.classList.toggle('expanded');
}

// Handle form submission
async function handleFormSubmit(e) {
    e.preventDefault();

    const containerName = document.getElementById('containerName').value.trim();
    const times = getTimes();
    const selectedDays = Array.from(document.querySelectorAll('input[name="days"]:checked'))
        .map(cb => cb.value);
    const envVars = getEnvVars();

    if (!containerName) {
        showToast('Por favor, introduce el nombre del contenedor', 'error');
        return;
    }

    if (selectedDays.length === 0) {
        showToast('Por favor, selecciona al menos un día', 'error');
        return;
    }

    if (times.length === 0) {
        showToast('Por favor, selecciona al menos una hora', 'error');
        return;
    }

    try {
        const runName = document.getElementById('runName').value.trim();
        const portsInput = document.getElementById('ports').value.trim();
        const ports = portsInput ? portsInput.split(',').map(p => p.trim()).filter(p => p) : [];
        const autoRemove = document.getElementById('autoRemove').checked;

        const response = await fetch(`${API_URL}/schedules`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({
                container_name: containerName,
                run_name: runName,
                ports: ports,
                auto_remove: autoRemove,
                days: selectedDays,
                times: times,
                exception_dates: [],
                env_vars: envVars
            })
        });

        if (!response.ok) {
            throw new Error('Error creating schedule');
        }

        scheduleForm.reset();
        clearEnvVars();
        clearTimes();
        await loadSchedules();
        showToast('Schedule creado exitosamente', 'success');
    } catch (error) {
        console.error('Error creating schedule:', error);
        showToast('Error al crear el schedule', 'error');
    }
}

// Toggle schedule active status
async function toggleSchedule(id) {
    try {
        const response = await fetch(`${API_URL}/schedules/${id}/toggle`, {
            method: 'PUT'
        });

        if (!response.ok) {
            throw new Error('Error toggling schedule');
        }

        await loadSchedules();
        showToast('Schedule actualizado', 'success');
    } catch (error) {
        console.error('Error toggling schedule:', error);
        showToast('Error al actualizar el schedule', 'error');
    }
}

// Delete schedule
async function deleteSchedule(id) {
    const confirmed = await showConfirm('¿Eliminar este schedule?');
    if (!confirmed) return;

    try {
        const response = await fetch(`${API_URL}/schedules/${id}`, {
            method: 'DELETE'
        });

        if (!response.ok) {
            throw new Error('Error deleting schedule');
        }

        await loadSchedules();
        showToast('Schedule eliminado', 'success');
    } catch (error) {
        console.error('Error deleting schedule:', error);
        showToast('Error al eliminar el schedule', 'error');
    }
}

// Exception modal functions
async function openExceptionModal(scheduleId) {
    currentExceptionScheduleId = scheduleId;
    exceptionModal.classList.add('show');

    // Load current exceptions
    try {
        const response = await fetch(`${API_URL}/schedules`);
        const schedules = await response.json();
        const schedule = schedules.find(s => s.id === scheduleId);

        if (schedule) {
            renderExceptions(schedule.exception_dates || []);
        }
    } catch (error) {
        console.error('Error loading exceptions:', error);
    }
}

function closeExceptionModal() {
    exceptionModal.classList.remove('show');
    currentExceptionScheduleId = null;
    document.getElementById('exceptionDateFrom').value = '';
    document.getElementById('exceptionDateTo').value = '';
}

function renderExceptions(exceptions) {
    const list = document.getElementById('exceptionsList');

    if (!exceptions || exceptions.length === 0) {
        list.innerHTML = '<p class="no-exceptions">No hay excepciones configuradas</p>';
        return;
    }

    // Sort by date
    exceptions.sort();

    list.innerHTML = exceptions.map(date => `
        <div class="exception-item">
            <span>${formatDate(date)}</span>
            <button class="btn-remove" onclick="removeExceptionDate('${date}')">
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                    <path d="M18 6L6 18M6 6l12 12"/>
                </svg>
            </button>
        </div>
    `).join('');
}

function formatDate(dateStr) {
    // Parse date parts manually to avoid timezone issues
    const [year, month, day] = dateStr.split('-').map(Number);
    const date = new Date(year, month - 1, day);
    return date.toLocaleDateString('es-ES', {
        weekday: 'short',
        year: 'numeric',
        month: 'short',
        day: 'numeric'
    });
}

// Generate array of dates between from and to (inclusive)
function getDateRange(fromStr, toStr) {
    const dates = [];

    // Parse dates manually to avoid timezone issues
    const [fromYear, fromMonth, fromDay] = fromStr.split('-').map(Number);
    const from = new Date(fromYear, fromMonth - 1, fromDay);

    let to = from;
    if (toStr) {
        const [toYear, toMonth, toDay] = toStr.split('-').map(Number);
        to = new Date(toYear, toMonth - 1, toDay);
    }

    const current = new Date(from);
    while (current <= to) {
        // Format as YYYY-MM-DD without using toISOString (which converts to UTC)
        const year = current.getFullYear();
        const month = String(current.getMonth() + 1).padStart(2, '0');
        const day = String(current.getDate()).padStart(2, '0');
        dates.push(`${year}-${month}-${day}`);
        current.setDate(current.getDate() + 1);
    }

    return dates;
}

async function addExceptionDates() {
    const fromInput = document.getElementById('exceptionDateFrom');
    const toInput = document.getElementById('exceptionDateTo');
    const fromDate = fromInput.value;
    const toDate = toInput.value;

    if (!fromDate) {
        showToast('Por favor, selecciona al menos la fecha de inicio', 'error');
        return;
    }

    // Validate date range
    if (toDate && toDate < fromDate) {
        showToast('La fecha de fin debe ser posterior a la de inicio', 'error');
        return;
    }

    const dates = getDateRange(fromDate, toDate);

    try {
        // Add each date
        let lastSchedule = null;
        for (const date of dates) {
            const response = await fetch(`${API_URL}/schedules/${currentExceptionScheduleId}/exceptions`, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify({ date: date })
            });

            if (!response.ok) {
                throw new Error('Error adding exception');
            }
            lastSchedule = await response.json();
        }

        renderExceptions(lastSchedule?.exception_dates || []);
        fromInput.value = '';
        toInput.value = '';
        await loadSchedules();
        showToast(`${dates.length} ${t('success')} (${t('added')})`, 'success');
    } catch (error) {
        console.error('Error adding exceptions:', error);
        showToast(t('errorGeneric'), 'error');
    }
}

async function removeExceptionDate(date) {
    try {
        const response = await fetch(`${API_URL}/schedules/${currentExceptionScheduleId}/exceptions`, {
            method: 'DELETE',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({ date: date })
        });

        if (!response.ok) {
            throw new Error('Error removing exception');
        }

        const schedule = await response.json();
        renderExceptions(schedule.exception_dates || []);
        await loadSchedules();
        showToast('Excepción eliminada', 'success');
    } catch (error) {
        console.error('Error removing exception:', error);
        showToast('Error al eliminar la excepción', 'error');
    }
}

// Show toast notification
function showToast(message, type = 'info') {
    toast.textContent = message;
    toast.className = `toast ${type} show`;

    setTimeout(() => {
        toast.classList.remove('show');
    }, 3000);
}

// Escape HTML to prevent XSS
function escapeHtml(text) {
    const div = document.createElement('div');
    div.textContent = text;
    return div.innerHTML;
}

// Custom confirm modal
let confirmResolve = null;

function showConfirm(message) {
    return new Promise((resolve) => {
        confirmResolve = resolve;
        document.getElementById('confirmMessage').textContent = message;
        document.getElementById('confirmModal').classList.add('show');
    });
}

function closeConfirmModal(result) {
    document.getElementById('confirmModal').classList.remove('show');
    if (confirmResolve) {
        confirmResolve(result);
        confirmResolve = null;
    }
}

// Load execution history
async function loadHistory() {
    try {
        const response = await fetch(`${API_URL}/history?limit=20`);
        const logs = await response.json();
        renderHistory(logs);
    } catch (error) {
        console.error('Error loading history:', error);
        historyList.innerHTML = `<p class="error-message">${t('errorGeneric')}</p>`;
    }
}

// Clear history
async function clearHistory() {
    const confirmed = await showConfirm(t('confirmClearHistory'));
    if (!confirmed) return;

    try {
        const response = await fetch(`${API_URL}/history`, {
            method: 'DELETE'
        });

        if (!response.ok) {
            throw new Error('Failed to clear history');
        }

        showToast(t('toastHistoryCleared'), 'success');
        await loadHistory();
    } catch (error) {
        console.error('Error clearing history:', error);
        showToast(t('errorGeneric'), 'error');
    }
}

// Render execution history
function renderHistory(logs) {
    if (!logs || logs.length === 0) {
        historyList.innerHTML = `
            <div class="empty-state">
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                    <circle cx="12" cy="12" r="10"/>
                    <polyline points="12 6 12 12 16 14"/>
                </svg>
                <p>${t('emptyHistory')}</p>
            </div>
        `;
        return;
    }

    // Group logs by container name
    const groupedLogs = logs.reduce((acc, log) => {
        if (!acc[log.container_name]) {
            acc[log.container_name] = [];
        }
        acc[log.container_name].push(log);
        return acc;
    }, {});

    historyList.innerHTML = Object.entries(groupedLogs).map(([containerName, containerLogs]) => {
        const lastRun = containerLogs[0]; // First one is most recent
        const successCount = containerLogs.filter(l => l.success).length;
        const failCount = containerLogs.filter(l => !l.success && l.status !== 'running').length;
        const runningCount = containerLogs.filter(l => l.status === 'running').length;

        return `
        <div class="history-group">
            <div class="history-group-header" onclick="toggleHistoryGroup(this)">
                <div class="history-group-info">
                    <div class="history-group-title">
                        <span class="group-name">${escapeHtml(containerName)}</span>
                        <span class="group-count">${containerLogs.length} ${t('runs')}</span>
                    </div>
                    <div class="history-group-stats">
                        ${runningCount > 0 ? `<span class="stat-running">${runningCount} ${t('running')}</span>` : ''}
                        ${successCount > 0 ? `<span class="stat-success">${successCount} ${t('success')}</span>` : ''}
                        ${failCount > 0 ? `<span class="stat-error">${failCount} ${t('failed')}</span>` : ''}
                    </div>
                </div>
                <div class="group-chevron">
                    <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                        <path d="M6 9l6 6 6-6"/>
                    </svg>
                </div>
            </div>
            <div class="history-group-content" style="display: none;">
                ${containerLogs.map(log => {
            let statusClass = log.status === 'running' ? 'running' : (log.success ? 'success' : 'error');
            let statusIcon;

            if (log.status === 'running') {
                statusIcon = '<div class="spinner-small"></div>';
            } else if (log.success) {
                statusIcon = '<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M22 11.08V12a10 10 0 1 1-5.93-9.14"/><path d="M22 4 12 14.01l-3-3"/></svg>';
            } else {
                statusIcon = '<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><circle cx="12" cy="12" r="10"/><path d="M15 9l-6 6M9 9l6 6"/></svg>';
            }

            return `
                    <div class="history-item ${statusClass}">
                        <div class="history-status">
                            ${statusIcon}
                        </div>
                        <div class="history-info">
                            <span class="history-time">${formatHistoryDate(log.executed_at)}</span>
                            ${log.status === 'running'
                    ? `<span class="history-running-tag">${t('statusRunning')}</span>`
                    : (log.error ? `<span class="history-error-msg" title="${escapeHtml(log.error)}">${escapeHtml(log.error)}</span>` : '')
                }
                        </div>
                        ${log.status === 'running' ? `
                            <button class="btn-view-log" onclick="viewLiveLogs('${escapeHtml(log.docker_name || log.container_name)}')" title="${t('viewLogs')}">
                                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                                    <path d="M1 12s4-8 11-8 11 8 11 8-4 8-11 8-11-8-11-8z"></path>
                                    <circle cx="12" cy="12" r="3"></circle>
                                </svg>
                            </button>
                        ` : ''}
                    </div>
                `}).join('')}
            </div>
        </div>
    `}).join('');
}

// Toggle history group expansion
function toggleHistoryGroup(header) {
    const content = header.nextElementSibling;
    const chevron = header.querySelector('.group-chevron svg');
    const isHidden = content.style.display === 'none';

    content.style.display = isHidden ? 'block' : 'none';
    chevron.style.transform = isHidden ? 'rotate(180deg)' : 'rotate(0deg)';
    header.classList.toggle('expanded', isHidden);
}

// View Execution Log
function viewLog(btn, output) {
    const historyItem = btn.closest('.history-item');
    const logViewer = historyItem.nextElementSibling;
    const codeBlock = logViewer.querySelector('code');

    // Toggle visibility
    if (logViewer.style.display === 'none') {
        codeBlock.textContent = output || 'No hay output disponible';
        logViewer.style.display = 'block';
        btn.classList.add('active');
    } else {
        logViewer.style.display = 'none';
        btn.classList.remove('active');
    }
}

// Format history date
function formatHistoryDate(dateStr) {
    const date = new Date(dateStr);
    return date.toLocaleString(currentLang === 'es' ? 'es-ES' : 'en-US', {
        day: '2-digit',
        month: 'short',
        hour: '2-digit',
        minute: '2-digit'
    });
}

// Load running containers
async function loadContainers() {
    try {
        const response = await fetch(`${API_URL}/containers`);
        const containers = await response.json();
        renderContainers(containers);
    } catch (error) {
        console.error('Error loading containers:', error);
        containersList.innerHTML = `<p class="error-message">${t('errorGeneric')}</p>`;
    }
}

// Render running containers
function renderContainers(containers) {
    if (!containers || containers.length === 0) {
        containersList.innerHTML = `
            <div class="empty-state">
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                    <rect x="2" y="4" width="20" height="16" rx="2"/>
                    <path d="M7 8h10M7 12h6"/>
                </svg>
                <p>${t('emptyContainers')}</p>
            </div>
        `;
        return;
    }

    containersList.innerHTML = containers.map(container => `
        <div class="container-item">
            <div class="container-info">
                <span class="container-name">${escapeHtml(container.name)}</span>
                <span class="container-image">${escapeHtml(container.image)}</span>
                ${container.ports ? `<span class="container-ports">${escapeHtml(container.ports)}</span>` : ''}
            </div>
            <div class="container-actions">
                <button class="btn-view-log" onclick="viewLiveLogs('${escapeHtml(container.name)}')" title="${t('viewLogs')}">
                    <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                         <path d="M1 12s4-8 11-8 11 8 11 8-4 8-11 8-11-8-11-8z"></path>
                         <circle cx="12" cy="12" r="3"></circle>
                    </svg>
                </button>
                <button class="btn-stop" onclick="stopContainer('${escapeHtml(container.name)}')" title="${t('stopContainer')}">
                    <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                        <rect x="6" y="6" width="12" height="12" rx="2"/>
                    </svg>
                </button>
            </div>
        </div>
    `).join('');
}

// Stop a running container
async function stopContainer(name) {
    const confirmed = await showConfirm(`${t('confirmStopContainer')} "${name}"?`);
    if (!confirmed) return;

    try {
        const response = await fetch(`${API_URL}/containers/${name}/stop`, {
            method: 'POST'
        });

        if (!response.ok) {
            throw new Error('Failed to stop container');
        }

        showToast(t('toastContainerStopped'), 'success');
        await loadContainers();
    } catch (error) {
        console.error('Error stopping container:', error);
        showToast(t('errorGeneric'), 'error');
    }
}

// Close Log Modal (Global)
function closeLogModal() {
    const modal = document.getElementById('logModal');
    if (modal) {
        modal.classList.remove('show');
    }
    // Clear auto-refresh if any
    if (window.logRefreshInterval) {
        clearInterval(window.logRefreshInterval);
        window.logRefreshInterval = null;
    }
}

// Open Log Viewer (Unified implementation for live and static logs)
async function openLogViewer(title, sourceOrContent, isLive = true) {
    // Check if modal exists, if not create it (fallback)
    let modal = document.getElementById('logModal');
    if (!modal) {
        // This fallback should rarely run if index.html is correct
        modal = document.createElement('div');
        modal.id = 'logModal';
        modal.className = 'modal log-modal-modern';
        modal.innerHTML = `
            <div class="modal-content log-modal-content">
                <div class="modal-header">
                    <div class="header-title">
                        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                            <path d="M4 17l6-6-6-6M12 19h8"/>
                        </svg>
                        <h3>Logs: <span id="logModalTitle"></span></h3>
                    </div>
                    <button class="close-modal-btn" onclick="closeLogModal()" aria-label="Cerrar">&times;</button>
                </div>
                <div class="modal-body terminal-body">
                    <div class="terminal-header">
                        <div class="terminal-dots">
                            <span class="dot red"></span>
                            <span class="dot yellow"></span>
                            <span class="dot green"></span>
                        </div>
                        <span class="terminal-title">bash</span>
                    </div>
                    <pre class="live-log-viewer"><code id="logModalContent"></code></pre>
                </div>
                <div class="modal-footer" style="border-top: 1px solid #334155; padding: 1rem; display: flex; justify-content: flex-end;">
                   <button class="btn-secondary" onclick="closeLogModal()" data-i18n="close">Close</button>
                </div>
            </div>
        `;
        document.body.appendChild(modal);

        // Close on outside click (attach once)
        window.addEventListener('click', (e) => {
            if (e.target === modal) {
                closeLogModal();
            }
        });
    }

    const modalTitle = document.getElementById('logModalTitle');
    const content = document.getElementById('logModalContent');

    modalTitle.textContent = title;
    content.textContent = 'Cargando...';

    // Show modal
    modal.classList.add('show');

    // Always live mode now
    const containerName = sourceOrContent;
    const fetchLogs = async () => {
        try {
            const response = await fetch(`${API_URL}/containers/${containerName}/logs`);
            if (!response.ok) throw new Error('Error fetching logs');
            const data = await response.json();
            content.textContent = data.logs || 'No hay logs disponibles (contenedor nuevo o silencioso)';
        } catch (error) {
            content.textContent = `Error al conectar con el contenedor "${containerName}".\n\nPosibles causas:\n1. El contenedor ya ha finalizado y ha sido eliminado (auto-remove).\n2. El contenedor se está reiniciando.\n\nDetalle error: ${error.message}`;
        }
    };

    await fetchLogs();

    // Clear previous interval if exists
    if (window.logRefreshInterval) clearInterval(window.logRefreshInterval);
    // Auto-refresh every 2s
    window.logRefreshInterval = setInterval(fetchLogs, 2000);
}

// Backwards compatibility wrappers
function viewLiveLogs(containerName) {
    openLogViewer(containerName, containerName, true);
}

// Removed
// function viewLog(btn, output) {
//     // btn is ignored
//     openLogViewer("Log Finalizado", output, false);
// }

// Make functions global for onclick handlers
window.viewLiveLogs = viewLiveLogs;
// Removed
// window.viewLog = viewLog;
window.stopContainer = stopContainer;
window.toggleHistoryGroup = toggleHistoryGroup;
window.clearHistory = clearHistory;
window.removeTimeRow = removeTimeRow;
window.removeExceptionDate = removeExceptionDate;
window.addExceptionDates = addExceptionDates;
window.closeExceptionModal = closeExceptionModal;
window.addEnvVarRow = addEnvVarRow;
window.toggleSchedule = toggleSchedule;
window.deleteSchedule = deleteSchedule;
window.openExceptionModal = openExceptionModal;
window.toggleGroup = toggleGroup;
window.toggleEnvExpandable = toggleEnvExpandable;
window.closeConfirmModal = closeConfirmModal;
window.closeLogModal = closeLogModal;

