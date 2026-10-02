# Tareas 001 — Errores estándar y validación de entrada

## Fase 1 — `models`

- [x] **T-01** `REQ-001` — Agregar `APIError` a `models/models.go` (`Status int`, `Message string`,
  `Details interface{}` con `json:"details,omitempty"`).
  - [ ] Agregar `TenantID string` con tag `json:"-"` a `Producto`.
  - [ ] Los tags `json` de `LoginRequest` y `Producto` quedan intactos.

## Fase 2 — `repository`

- [x] **T-02** `REQ-001` — `repository/usuarioRepo.go`: `ErrNoEncontrado`, `BuscarPorEmail(tenantID, email)`,
  seed con las 2 credenciales de la línea base.
  - [ ] Filtra por `tenant_id`.
  - [ ] No importa `fiber`.
- [x] **T-03** `REQ-002` — `repository/productoRepo.go`: `Listar(tenantID)`, `BuscarPorID(tenantID, id)`,
  seed con los mismos 4 productos.
  - [ ] Ambos métodos filtran por `tenant_id`.
  - [ ] `BuscarPorID` devuelve `ErrNoEncontrado` si no existe o es de otro tenant.

## Fase 3 — `controllers`

- [x] **T-04** `REQ-001` — `controllers/errorController.go`: `TenantDesdeRequest`,
  `ResponderError` (envelope estándar), `ResponderErrorLegado` (forma congelada).
  - [ ] `details` se omite cuando es `nil`.
- [x] **T-05** `REQ-002` — `Login` valida `email` y `password` antes del repository.
  - [ ] Vacío → `400` con `ResponderErrorLegado` y mensaje "Email y password son obligatorios".
  - [ ] `400` por parseo y `401` por credenciales con `ResponderErrorLegado`.
  - [ ] Consulta usuarios vía repository; `200` con `token`, `email`, `rol` idénticos a la línea base.
- [x] **T-06** `REQ-002` — `GetProductos` consume `repository.Listar` con el tenant del request
  y sigue devolviendo el array plano.
- [x] **T-07** `REQ-003` — `GetProductoPorID` en `controllers/prodController.go`.
  - [ ] Lee `id` de `c.Params("id")`.
  - [ ] `strconv.Atoi`; error de conversión → `400` y `return` sin llamar al repository.
  - [ ] No encontrado o de otro tenant → `404` con `ResponderError`.
- [x] **T-08** `REQ-003` — Registrar `GET /api/productos/:id` en `routes/routes.go`.

## Fase 4 — Verificación

- [ ] `go build ./...` sin errores.
- [ ] `go vet ./...` sin observaciones.
- [ ] `gofmt -l .` sin salida.
- [x] Casos CV-01..CV-04 de `validation.md` probados con `curl`.
- [x] No-regresión: `/api/login` y `/api/productos` idénticos a `specs/000-contrato-actual/`.

## Dependencias

| Tarea | Depende de |
|-------|------------|
| T-02   | T-01       |
| T-03   | T-01       |
| T-04   | T-01       |
| T-05   | T-02, T-04  |
| T-06   | T-03, T-04  |
| T-07   | T-03, T-04  |
| T-08   | T-07       |