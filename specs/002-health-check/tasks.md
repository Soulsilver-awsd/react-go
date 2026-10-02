# Tareas 002 — Health check del servidor

## Fase 1 — `models`

- [x] **T-01** `REQ-004` — Agregar `HealthResponse` a `models/models.go` con tags
  `json:"status"` y `json:"message"`.

## Fase 2 — `controllers`

- [x] **T-02** `REQ-004` — Crear `controllers/healthController.go` con `Health`:
  - [ ] Devuelve `c.JSON(models.HealthResponse{Status: "ok", Message: "Servidor Go/Fiber operativo"})`.
  - [ ] Sin dependencias, sin I/O, sin leer headers.

## Fase 3 — `routes`

- [x] **T-03** `REQ-005` — Registrar `app.Get("/api/health", controllers.Health)` en `SetupRoutes`.

## Fase 4 — Verificación

- [ ] `go build ./...` sin errores.
- [ ] `go vet ./...` sin observaciones.
- [ ] `gofmt -l .` sin salida.
- [x] `curl -s -i http://localhost:3000/api/health` devuelve `200` y el payload de CA-014.
- [x] No-regresión: `/api/login` y `/api/productos` idénticos a la línea base.

## Dependencias

| Tarea | Depende de |
|-------|------------|
| T-02   | T-01       |
| T-03   | T-02       |