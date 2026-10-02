# Estado: validada

# Spec 001 — Errores estándar y validación de entrada (US-01, US-02, US-03)

## 1. Contexto

Los controllers de `multicatalogo-backend/` arman hoy sus errores a mano con `fiber.Map{"error": "..."}`
y validan nada: `/api/login` acepta cuerpos vacíos y no existe endpoint de detalle de producto.
El frontend ya consume `/api/login` y `/api/productos`, así que **el JSON observable no puede cambiar**.

## 2. Alcance

**Incluye**

- US-01: estructura `APIError` reutilizable en `models`.
- US-02: validación de `email` y `password` en `Login` antes de tocar datos.
- US-03: `GetProductoPorID` con parsing explícito del `id` de la ruta.

**Excluye**

- US-04/US-05 (health check) → spec `002-health-check`.
- Cambiar la forma de las respuestas de `/api/login` y `/api/productos` (ver No-regresión).
- Persistencia real (SQL/ORM). Los repositories nacen en memoria.

## 3. Requisitos

### REQ-001 — Estructura `APIError` (US-01)

**Ubicación**: `models`.

EARS: *El sistema shall exponer una estructura `APIError` en `models/models.go` con los campos `Status int`, `Message string` y `Details interface{}` opcional.*

- **CA-001**: Given el paquete `models`, When se inspecciona `models.go`, Then existe `APIError` con exactamente los campos `Status int`, `Message string`, `Details interface{}` y `Details` es opcional (se omite cuando es `nil`).
- **CA-002**: Given un `APIError` con `Details == nil`, When se serializa como envelope estándar, Then la clave `details` no aparece en el JSON.

### REQ-002 — Login rechaza credenciales incompletas (US-02)

**Ubicación**: `controllers` + `repository`.

EARS: *While la petición llega a `POST /api/login`, el sistema shall validar `email` y `password` antes de invocar la capa `repository`.*

- **CA-003**: Given `{"email":"","password":"123456"}`, When `POST /api/login`, Then responde `400`.
- **CA-004**: Given `{"email":"admin@upse.edu.ec","password":""}`, When `POST /api/login`, Then responde `400`.
- **CA-005**: Given `{}` (ambos vacíos), When `POST /api/login`, Then responde `400`.
- **CA-006**: Given un `400` o `401` de `/api/login`, When el frontend parsea el cuerpo, Then recibe la forma congelada `{"error": "<mensaje>"}` construida a partir de `APIError`.
- **CA-007**: Given credenciales completas y válidas, When `POST /api/login`, Then responde `200` con `token`, `email`, `rol` idénticos a la línea base.

### REQ-003 — Validación del id de producto (US-03)

**Ubicación**: `controllers` (parseo) + `repository` (consulta).

EARS: *When `GET /api/productos/:id` recibe un `id`, el sistema shall convertirlo a entero y manejar el error de conversión explícitamente, sin invocar `repository` cuando la conversión falla.*

- **CA-008**: Given `GET /api/productos/2`, When la petición es válida, Then responde `200` con el objeto `Producto` de `id` 2.
- **CA-009**: Given `GET /api/productos/abc`, When el `id` no es numérico, Then responde `400` con mensaje legible.
- **CA-010**: Given un `id` no numérico, When se procesa la petición, Then la capa `repository` no es invocada.
- **CA-011**: Given `GET /api/productos/999`, When el producto no existe, Then responde `404`.

## 4. Contrato de la API

### `POST /api/login` (endpoint congelado — forma intacta)

Request

```json
{ "email": "admin@upse.edu.ec", "password": "123456" }
```

Respuesta 200 (sin cambios)

```json
{ "token": "fake-jwt-token-123", "email": "admin@upse.edu.ec", "rol": "admin" }
```

Errores (sin cambios de forma)

| Código | Condición                       | Cuerpo                                              |
|--------|---------------------------------|-----------------------------------------------------|
| 400    | JSON mal formado                | `{"error":"Cuerpo de petición inválido"}`           |
| 400    | `email` o `password` vacío      | `{"error":"Email y password son obligatorios"}`     |
| 401    | Credenciales incorrectas        | `{"error":"Credenciales incorrectas"}`              |

### `GET /api/productos/:id` (endpoint nuevo — aditivo)

Respuesta 200

```json
{ "id": 2, "nombre": "Crema Hidratante Pro", "precio": 32.5, "img": "https://picsum.photos/seed/crema/150" }
```

Errores

| Código | Condición                          | Cuerpo                                                                   |
|--------|------------------------------------|--------------------------------------------------------------------------|
| 400    | `id` no numérico                   | `{"success":false,"error":"El id debe ser un número entero"}`             |
| 404    | Producto inexistente o de otro tenant | `{"success":false,"error":"Producto no encontrado"}`                |
| 500    | Error interno                      | `{"success":false,"error":"Error interno del servidor"}`                 |

## 5. Multitenant

- Origen del tenant: header `X-Tenant-ID`. Ausente → tenant por defecto `default`
  (transicional; preserva la no-regresión de `/api/productos` que el frontend llama sin header).
- El `tenant_id` **nunca** se lee del body ni del path.
- Toda consulta de `repository` filtra por `tenant_id`.
- Un `id` de producto de otro tenant responde `404`, no `403`.
- `models.Producto` gana `TenantID string` con tag `json:"-"`: no altera el JSON congelado.

## 6. No-regresión

- [ ] `/api/productos` sigue devolviendo un array plano `{id, nombre, precio, img}`.
- [ ] `/api/login` sigue devolviendo `token`, `email`, `rol`; roles `admin` y `cliente`.
- [ ] Los errores de `/api/login` mantienen la forma `{"error": "..."}`.
- [ ] `go build/vet/gofmt` limpios.

## 7. Fuera de alcance

- Hash de contraseñas, JWT real, base de datos.
- Validación de formato de email (solo presencia no vacía).
- `GET /api/red`.

## 8. Trazabilidad

| Requisito | Historia | Tarea     | Evidencia | Estado |
|-----------|----------|-----------|-----------|--------|
| REQ-001   | US-01    | T-01      | CV-00     | validada |
| REQ-002   | US-02    | T-02, T-05 | CV-02..CV-08 | validada |
| REQ-003   | US-03    | T-03, T-07, T-08 | CV-10..CV-12 | validada |