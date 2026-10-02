# Validación 002 — Health check del servidor

> Evidencia real. Toolchain: Go 1.27.1 (el repo declara `go 1.26.0`; no había binario de Go
> instalado, se usó un SDK descargado en `/tmp` solo para verificar).
> Servidor: `go run .` escuchando en `:3000`.

## 1. Criterios de aceptación

| Criterio | Evidencia | Estado |
|----------|-----------|--------|
| CA-012 `GET /api/health` → 200 | CV-01 | PASA |
| CA-013 cuerpo JSON con exactamente `status` y `message` | CV-01 | PASA |
| CA-014 valores `"ok"` y `"Servidor Go/Fiber operativo"` | CV-01 | PASA |
| CA-015 ruta registrada en `routes/routes.go` | CV-02 | PASA |
| CA-016 accesible tras arrancar la aplicación | CV-01 | PASA |

## 2. Evidencia de ejecución

### CV-01 — `GET /api/health` (CA-012, CA-013, CA-014, CA-016)

```bash
curl -s -i http://localhost:3000/api/health
```

```
HTTP/1.1 200 OK
Date: Fri, 02 Oct 2026 15:49:32 GMT
Content-Type: application/json
Content-Length: 55
Vary: Origin

{"status":"ok","message":"Servidor Go/Fiber operativo"}
```

Dos claves exactas, sin `success`/`data` (desviación documentada en `spec.md` §4, impuesta por US-04).

### CV-02 — Registro en `routes/routes.go` (CA-015)

```bash
grep -n "api/health" routes/routes.go
```

```
24:	app.Get("/api/health", controllers.Health)
```

`controllers.Health` está definido en `controllers/healthController.go`.

## 3. Verificación estática

```bash
cd multicatalogo-backend
go build ./...    # sin salida, exit 0
go vet ./...      # sin salida, exit 0
gofmt -l .        # sin salida
```

## 4. No-regresión del contrato congelado

- [x] `POST /api/login` → `200` con `token`, `email`, `rol`; errores `400`/`401` con `{"error":...}`.
- [x] `GET /api/productos` → array plano `{id, nombre, precio, img}` con los mismos 4 productos.
- [ ] `GET /api/red` → aún no existe (spec `003`).
- [x] `main.go` sin cambios: CORS y puerto 3000 intactos.

## 5. Multitenant

No aplica: el handler no lee `X-Tenant-ID` ni accede a datos de negocio.

## 6. Resultado

Todos los criterios en PASA, no-regresión limpio. Spec en estado `validada`.