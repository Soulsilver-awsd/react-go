# Tareas NNN — <Título de la feature>

> Plantilla. Copiala a `specs/NNN-slug/tasks.md`.
> Reglas:
> - Orden obligatorio: `models` → `repository` → `controllers` → `routes`.
> - Cada tarea es atómica: un solo archivo o un solo cambio de comportamiento.
> - Toda tarea referencia el requisito (`REQ-00X`) que implementa y cómo se verifica.
> - Solo se marca `[x]` cuando la verificación passa, no cuando el código existe.

## Fase 1 — `models`

- [ ] **T-01** `REQ-001` — Agregar `Transaccion` a `models/models.go`.
  - [ ] Tags `json` alineados con el contrato de `spec.md`.
  - [ ] `gofmt -l .` sin salida.

## Fase 2 — `repository`

- [ ] **T-02** `REQ-001` — `repository/transaccionRepo.go`: `ListarPorUsuario(tenantID, usuarioID)`.
  - [ ] Filtra por `tenant_id` **siempre**.
  - [ ] No importa `fiber`.
  - [ ] Devuelve `ErrNoEncontrado` (no `nil, nil`).

## Fase 3 — `controllers`

- [ ] **T-03** `REQ-001` — `GetTransacciones` en `controllers/`.
  - [ ] `400` si falta el header de tenant.
  - [ ] `401` si falta/expira el token.
  - [ ] Envelope `{success, data}`.
  - [ ] No toca datos directamente: llama al repository.

## Fase 4 — `routes`

- [ ] **T-04** `REQ-001` — Registrar `GET /api/transacciones` en `SetupRoutes`.

## Fase 5 — Verificación

- [ ] `go build ./...` sin errores.
- [ ] `go vet ./...` sin observaciones.
- [ ] `gofmt -l .` sin salida.
- [ ] Casos de `spec.md` probados con `curl` (transcribir salida a `validation.md`).
- [ ] No-regresión: `/api/login`, `/api/productos`, `/api/red` responden igual que antes.

## Dependencias

| Tarea | Depende de |
|-------|------------|
| T-02   | T-01       |
| T-03   | T-02       |
| T-04   | T-03       |