# Estado: validada

# Spec 002 — Health check del servidor (US-04, US-05)

## 1. Contexto

No hay forma de verificar si el proceso Go/Fiber está vivo sin conocer credenciales:
`POST /api/login` responde 401 si el body está mal formado, y `GET /api/productos` no distingue
"server caído" de "catálogo vacío". Un endpoint público de salud permite orquestación y depuración.

## 2. Alcance

**Incluye**: US-04 (`controllers/healthController.go`) y US-05 (registro en `routes.go`).

**Excluye**: readiness contra dependencias (no hay base de datos todavía), métricas, ping a
recursos externos, autenticación del endpoint.

## 3. Requisitos

### REQ-004 — Endpoint de salud (US-04)

**Ubicación**: `controllers`.

EARS: *When `GET /api/health` recibe una petición, el sistema shall responder `200` con el cuerpo
JSON de estado operativo.*

- **CA-012**: Given la aplicación arrancada, When `GET /api/health`, Then responde `200`.
- **CA-013**: Given la respuesta de `/api/health`, When se parsea el cuerpo, Then es JSON válido y
  contiene exactamente `status` y `message`.
- **CA-014**: Given la respuesta de `/api/health`, When se lee `status`, Then vale `"ok"` y
  `message` vale `"Servidor Go/Fiber operativo"`.

### REQ-005 — Registro de la ruta (US-05)

**Ubicación**: `routes`.

EARS: *While la aplicación está inicializada, el sistema shall registrar `GET /api/health` en
`SetupRoutes` de modo que quede accesible por HTTP.*

- **CA-015**: Given `SetupRoutes` ejecutado, When se inspecciona `routes/routes.go`, Then registra
  `GET /api/health` con el handler de health.
- **CA-016**: Given la aplicación escuchando en el puerto 3000, When `GET /api/health`, Then
  devuelve el mismo `200` y el mismo payload que en CA-014.

## 4. Contrato de la API

### `GET /api/health` (endpoint nuevo)

Sin request.

Respuesta 200

```json
{ "status": "ok", "message": "Servidor Go/Fiber operativo" }
```

| Código | Condición          | Cuerpo |
|--------|--------------------|--------|
| 200    | Proceso operativo  | `{"status":"ok","message":"Servidor Go/Fiber operativo"}` |

Nota: este endpoint **no** usa el envelope `{success, data}` porque US-04 fija la forma exacta
del payload. La desviación está acotada a este endpoint y documentada aquí.

## 5. Multitenant

No aplica: el endpoint no devuelve ni acepta datos de negocio. No lee `X-Tenant-ID`.

## 6. No-regresión

- [ ] `/api/login` y `/api/productos` sin cambios respecto de `specs/000-contrato-actual/`.
- [ ] `/api/health` es aditivo: no altera CORS ni el puerto 3000.

## 7. Fuera de alcance

- Sondas de dependencias (DB, caché).
- Versionado del payload, headers custom ni health check autenticado.

## 8. Trazabilidad

| Requisito | Historia | Tarea  | Evidencia | Estado |
|-----------|----------|--------|-----------|--------|
| REQ-004   | US-04    | T-01, T-02 | CV-01  | validada |
| REQ-005   | US-05    | T-03   | CV-01, CV-02 | validada |