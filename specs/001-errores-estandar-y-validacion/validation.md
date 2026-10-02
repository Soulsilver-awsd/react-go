# Validación 001 — Errores estándar y validación de entrada

> Evidencia real ejecutada en el entorno. Toolchain: Go 1.27.1 (el repo declara `go 1.26.0`;
> no había binario de Go instalado, se usó un SDK descargado en `/tmp` solo para verificar).
> Servidor: `go run .` escuchando en `:3000`.

## 1. Criterios de aceptación

| Criterio | Evidencia | Estado |
|----------|-----------|--------|
| CA-001 `APIError` con `Status int`, `Message string`, `Details interface{}` opcional | CV-00 | PASA |
| CA-002 `details` se omite cuando es `nil` | CV-00, CV-07 | PASA |
| CA-003 `email` vacío → 400 | CV-04 | PASA |
| CA-004 `password` vacío → 400 | CV-05 | PASA |
| CA-005 body `{}` → 400 | CV-06 | PASA |
| CA-006 errores de login conservan `{"error": ...}` | CV-04..CV-08 | PASA |
| CA-007 login válido → 200 con `token`, `email`, `rol` idénticos | CV-02, CV-03 | PASA |
| CA-008 `GET /api/productos/2` → 200 con el producto | CV-10 | PASA |
| CA-009 `id` no numérico → 400 | CV-11 | PASA |
| CA-010 `id` inválido no llega al repository | CV-12 | PASA |
| CA-011 producto inexistente → 404 | CV-12 | PASA |

## 2. Evidencia de ejecución

### CV-00 — Estructura `APIError` en `models/models.go`

`grep -A 12 "^type APIError" models/models.go`:

```go
type APIError struct {
	// Status es el código HTTP que corresponde al error (400, 401, 404, 500, etc.).
	Status int `json:"status"`
	// Message es el texto legible y estable que describe qué ocurrió, apto para mostrarse al usuario.
	Message string `json:"message"`
	// Details es un campo opcional para información adicional (por ejemplo, el campo del cuerpo
	// que falló). Al declararse con omitempty, si queda en nil la clave desaparece del JSON,
	// de modo que el payload se mantiene limpio.
	Details interface{} `json:"details,omitempty"`
}
```

### CV-02 — Login admin (no-regresión, CA-007)

```bash
curl -s -X POST http://localhost:3000/api/login -H "Content-Type: application/json" \
  -d '{"email":"admin@upse.edu.ec","password":"123456"}'
```

```json
{"email":"admin@upse.edu.ec","rol":"admin","token":"fake-jwt-token-123"}
```

### CV-03 — Login cliente (no-regresión, CA-007)

```json
{"email":"cliente@upse.edu.ec","rol":"cliente","token":"fake-jwt-token-456"}
```

### CV-04 — `email` vacío → 400 (CA-003, CA-006)

```bash
curl -s -i -X POST http://localhost:3000/api/login -H "Content-Type: application/json" \
  -d '{"email":"","password":"123456"}'
```

```
HTTP/1.1 400 Bad Request
{"error":"Email y password son obligatorios"}
```

### CV-05 — `password` vacío → 400 (CA-004, CA-006)

```
HTTP/1.1 400 Bad Request
{"error":"Email y password son obligatorios"}
```

### CV-06 — Body `{}` → 400 (CA-005)

```
HTTP/1.1 400 Bad Request
{"error":"Email y password son obligatorios"}
```

### CV-07 — `details` omitido en la forma heredada (CA-002)

El `APIError` del 400 de login **sí** trae `details` con la lista de campos faltantes
(`["email"]`, `["password"]`, `["email","password"]`), pero `ResponderErrorLegado` solo
serializa `{"error": ...}` porque esa es la forma congelada. Comprobado: en CV-04, CV-05 y CV-06
la respuesta tiene una única clave `error` y **no** aparece `details` ni `status`.

### CV-08 — No-regresión de errores de login

```bash
curl -s -i -X POST http://localhost:3000/api/login -H "Content-Type: application/json" -d 'no-json'
```

```
HTTP/1.1 400 Bad Request
{"error":"Cuerpo de petición inválido"}
```

```bash
curl -s -i -X POST http://localhost:3000/api/login -H "Content-Type: application/json" \
  -d '{"email":"admin@upse.edu.ec","password":"mala"}'
```

```
HTTP/1.1 401 Unauthorized
{"error":"Credenciales incorrectas"}
```

### CV-09 — No-regresión `GET /api/productos`

```bash
curl -s http://localhost:3000/api/productos
```

```json
[{"id":1,"nombre":"Serum Revitalizante","precio":45,"img":"https://picsum.photos/seed/serum/150"},
 {"id":2,"nombre":"Crema Hidratante Pro","precio":32.5,"img":"https://picsum.photos/seed/crema/150"},
 {"id":3,"nombre":"Tónico Purificante","precio":28,"img":"https://picsum.photos/seed/tonico/150"},
 {"id":4,"nombre":"Mascarilla Nocturna","precio":50,"img":"https://picsum.photos/seed/mascarilla/150"}]
```

Array plano, cuatro claves por objeto, `tenant_id` ausente gracias al tag `json:"-"` (CA §5).

### CV-10 — Detalle válido → 200 (CA-008)

```bash
curl -s -i http://localhost:3000/api/productos/2
```

```
HTTP/1.1 200 OK
{"id":2,"nombre":"Crema Hidratante Pro","precio":32.5,"img":"https://picsum.photos/seed/crema/150"}
```

### CV-11 — `id` no numérico → 400 (CA-009)

```bash
curl -s -i http://localhost:3000/api/productos/abc
```

```
HTTP/1.1 400 Bad Request
{"details":{"parametro":"id","valor":"abc"},"error":"El id debe ser un número entero","success":false}
```

### CV-12 — `id` inválido no alcanza el repository (CA-010) + 404 (CA-011)

Prueba con traza temporal: se insertó una línea `println("[TRACE] repository.BuscarPorID INVOCADO ...")`
al inicio de `repository.BuscarPorID`, elevantó el servidor y se hicieron **dos** peticiones.

```bash
curl -s http://localhost:3000/api/productos/abc   # → 400
curl -s http://localhost:3000/api/productos/3    # → 200
```

Salida del servidor (log):

```
[TRACE] repository.BuscarPorID INVOCADO tenant=default
--- total trazas: 1 ---
```

Una sola traza para dos peticiones: la del `id` no numérico **nunca entró** al repository.
La línea de traza fue removida después de la prueba y el build volvió a pasar limpio.

Producto inexistente (CA-011):

```bash
curl -s -i http://localhost:3000/api/productos/999
```

```
HTTP/1.1 404 Not Found
{"details":{"id":999},"error":"Producto no encontrado","success":false}
```

## 3. Verificación estática

```bash
cd multicatalogo-backend
go build ./...    # sin salida, exit 0
go vet ./...      # sin salida, exit 0
gofmt -l .        # sin salida
```

Nota: los archivos del repo original **no** tenían salto de línea final, por lo que `gofmt -l .`
ya fallaba antes de este trabajo. Se aplicó `gofmt -w .` (solo formato) para dejar el paquete limpio.

## 4. No-regresión del contrato congelado

- [x] `POST /api/login` → `200` con `token`, `email`, `rol` (roles `admin` y `cliente`) — CV-02, CV-03.
- [x] `POST /api/login` cuerpo inválido → `400` `{"error":"Cuerpo de petición inválido"}` — CV-08.
- [x] `POST /api/login` credenciales malas → `401` `{"error":"Credenciales incorrectas"}` — CV-08.
- [x] `GET /api/productos` → array plano `{id, nombre, precio, img}` — CV-09.
- [ ] `GET /api/red` → aún no existe (spec `003`, fuera de esta spec).
- [x] CORS y puerto 3000 sin cambios (`main.go` no se tocó).

## 5. Multitenant

```bash
curl -s http://localhost:3000/api/productos -H "X-Tenant-ID: tenant-b"    # → []
curl -s -i http://localhost:3000/api/productos/2 -H "X-Tenant-ID: tenant-b" # → 404, no 403
curl -s -i -X POST http://localhost:3000/api/login -H "Content-Type: application/json" \
  -H "X-Tenant-ID: tenant-b" -d '{"email":"admin@upse.edu.ec","password":"123456"}' # → 401
```

- [x] Consulta con `tenant_id` de A no devuelve datos de B (`[]` y `404`).
- [x] ID de recurso de otro tenant responde `404`, no `403`.
- [x] `tenant_id` del body se ignora; manda el del request (header `X-Tenant-ID`).
- [x] Sin header → tenant `default`, que preserva el catálogo actual (CV-09).

## 6. Resultado

Todos los criterios de aceptación en PASA, no-regresión limpio y multitenant verificado.
Spec en estado `validada`.