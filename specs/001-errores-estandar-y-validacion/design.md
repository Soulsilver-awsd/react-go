# Diseño 001 — Errores estándar y validación de entrada

## 1. Resumen de la decisión

Se introduce `models.APIError` como representación interna del error, una capa `repository`
en memoria para `usuario` y `producto`, y helpers de serialización en `controllers` que deciden
qué JSON emitir: la forma congelada (`{"error": ...}`) en los endpoints de la línea base, y el
envelope estándar (`{"success":false,...}`) en los endpoints nuevos. Los controllers validan
antes de invocar el repository.

## 2. Decisiones (ADR)

### ADR-001 — `APIError` es una estructura en proceso, no un contrato de serialización único

- **Estado**: aceptada
- **Contexto**: la línea base congela `{"error": "..."}` para `/api/login` y `/api/productos`.
  `AGENTS.md` exige además un envelope `{success, error}` para endpoints nuevos. Un único
  `MarshalJSON` en `APIError` no puede satisfacer ambos.
- **Decisión**: `APIError` modela el error (Status/Message/Details) y la serialización se decide
  en `controllers`:
  - `ResponderError` → `{"success":false,"error":msg,"details":...}` (endpoints nuevos).
  - `ResponderErrorLegado` → `{"error":msg}` (endpoints congelados de la línea base).
- **Alternativas descartadas**: cambiar `/api/login` al envelope (rompe el parseo del frontend);
  dejar `fiber.Map` sin estructura (es exactamente el problema que reporta US-01).
- **Consecuencias**: una sola fuente de verdad del mensaje y del status; dos formas de salida,
  documentadas y limitadas por endpoint.

### ADR-002 — Repository en memoria con filtro obligatorio por `tenant_id`

- **Estado**: aceptada
- **Contexto**: hoy los datos están hardcodeados en los controllers. Los AC de US-02/US-03 exigen
  que el repository exista y que los datos inválidos no lleguen a esa capa.
- **Decisión**: `repository/usuarioRepo.go` y `repository/productoRepo.go` con slices en memoria.
  Cada método recibe `tenantID string` y filtra por él. Error sentinel: `ErrNoEncontrado`.
- **Alternativas descartadas**: SQLite/PostgreSQL (dependencia nueva; `go.mod` solo tiene Fiber);
  dejar los datos en el controller (incumple US-01/US-02/US-03 y el patrón de capas).
- **Consecuencias**: reposiciona la API en el patrón sin tocar el contrato; migrar a SQL después
  es cambiar la implementación del repository, no el controller.

### ADR-003 — Tenant por header con valor por defecto `default`

- **Estado**: aceptada
- **Contexto**: el frontend llama `/api/productos` sin cabeceras propias; si el endpoint exigiera
  `X-Tenant-ID` devolvería 401 y rompería el catálogo.
- **Decisión**: `TenantDesdeRequest(c)` resuelve `X-Tenant-ID` y cae a `default` si falta o vacío.
- **Alternativas descartadas**: `401` sin header (rompe no-regresión); claim del token
  (`fake-jwt-token-123` todavía no es un JWT real).
- **Consecuencias**: transición segura; cuando exista JWT real se cambia el origen sin tocar
  controllers ni repositorios.

### ADR-004 — El `id` inválido nunca llega al repository

- **Estado**: aceptada
- **Contexto**: CA-010 lo exige; además evita que un error de conversión se pierda como 500.
- **Decisión**: `strconv.Atoi` en el controller; si falla se responde 400 y se retorna de inmediato
  (`return`), sin llamada al repository.
- **Alternativas descartadas**: validar dentro del repository (mezcla responsabilidades y
  devuelve 500 en vez de 400); ignorar el error con `_`.

## 3. Capas

| Capa          | Responsabilidad                                      | Archivos                                    |
|---------------|------------------------------------------------------|---------------------------------------------|
| `models`      | `APIError`, `TenantID` en `Producto`                  | `models/models.go`                          |
| `repository`  | Datos en memoria, filtro por tenant, `ErrNoEncontrado` | `repository/usuarioRepo.go`, `repository/productoRepo.go` |
| `controllers` | Validación, tenant, serialización de errores          | `controllers/errorController.go`, `authController.go`, `prodController.go` |
| `routes`      | `GET /api/productos/:id`                              | `routes/routes.go`                          |

## 4. Modelo de datos

| Colección  | Campos                                       | Filtro     |
|------------|----------------------------------------------|------------|
| `usuario`  | `Email`, `Password`, `Rol`, `TenantID`        | `TenantID` |
| `producto` | `ID`, `Nombre`, `Precio`, `Img`, `TenantID`   | `TenantID` |

`TenantID` de `Producto` lleva `json:"-"` → invisible en la respuesta congelada.

## 5. Multitenant

- Origen: header `X-Tenant-ID` (ADR-003).
- Validación: `controllers.TenantDesdeRequest`.
- `TenantID` nunca se toma del body ni del path.
- Consultas por id: `404` si el recurso pertenece a otro tenant.

## 6. Datos de arranque

Seed en memoria con las mismas credenciales y los mismos 4 productos de la línea base, de modo
que la respuesta de `/api/productos` sea idéntica antes y después del cambio.

## 7. Riesgos

| Riesgo                                             | Mitigación                                                 |
|----------------------------------------------------|------------------------------------------------------------|
| Cambiar por accidente la forma de `{"error"}`        | Helpers separados (legado vs estándar) + CA-006 + no-regresión |
| `TenantID` asomando en el JSON del producto          | Tag `json:"-"` + CV de `/api/productos`                    |
| Datos inválidos llegando al repository               | `return` inmediato tras validar (ADR-004)                  |