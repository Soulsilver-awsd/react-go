# Tareas 000 — Contrato actual (línea base)

> Registro de lo que hoy existe en `multicatalogo-backend/`. No hay tareas pendientes aquí:
> esta spec congela el estado, no propone trabajo. El trabajo nuevo vive en `001-...`.

## Fase 1 — `models`

- [x] **T-01** Contrato base — `models/models.go` con `LoginRequest` (`email`, `password`)
  y `Producto` (`id`, `nombre`, `precio`, `img`).
  - [x] Tags `json` alineados con el contrato del frontend.

## Fase 2 — `repository`

- [ ] **T-02** No existe capa `repository`. Toda lectura de datos está hardcodeada en el controller.
  Requiere spec propia con diseño de multitenant.

## Fase 3 — `controllers`

- [x] **T-03** `controllers/authController.go` — `Login`: `400` por parseo, `401` por credenciales,
  `200` con `token`/`email`/`rol`.
- [x] **T-04** `controllers/prodController.go` — `GetProductos`: devuelve el array plano de 4 productos.

## Fase 4 — `routes`

- [x] **T-05** `routes/routes.go` — `SetupRoutes` registra `POST /api/login` y `GET /api/productos`.
- [ ] **T-06** `GET /api/red` **no registrado**, aunque `src/components/MiRed.tsx` existe en el frontend.
  Requiere spec propia (`001-red-mlm`).

## Fase 5 — Verificación

- [ ] `go build ./...` sin errores.
- [ ] `go vet ./...` sin observaciones.
- [ ] `gofmt -l .` sin salida.
- [ ] No-regresión: `/api/login` y `/api/productos` responden según `spec.md`.