# Validación 000 — Contrato actual (línea base)

> Evidencia del contrato congelado. Correla de nuevo después de cualquier feature:
> si alguno de estos casos cambia, la feature rompió el frontend.

## 1. Criterios de aceptación

| Criterio                                          | Evidencia | Estado |
|---------------------------------------------------|-----------|--------|
| CA-001 `POST /api/login` admin → `200` con rol    | CV-01     | PASA   |
| CA-002 `POST /api/login` cliente → `200` con rol  | CV-02     | PASA   |
| CA-003 `POST /api/login` cuerpo inválido → `400`   | CV-03     | PASA   |
| CA-004 `POST /api/login` credenciales malas → `401`| CV-04     | PASA   |
| CA-005 `GET /api/productos` → array plano          | CV-05     | PASA   |
| CA-006 `GET /api/red` → estructura MLM             | —         | PENDIENTE (no implementado) |

## 2. Evidencia de ejecución

### CV-01 — Login admin

```bash
cd multicatalogo-backend && go run . &
curl -s -X POST http://localhost:3000/api/login \
  -H "Content-Type: application/json" \
  -d '{"email":"admin@upse.edu.ec","password":"123456"}'
```

```json
{"token":"fake-jwt-token-123","email":"admin@upse.edu.ec","rol":"admin"}
```

### CV-02 — Login cliente

```bash
curl -s -X POST http://localhost:3000/api/login \
  -H "Content-Type: application/json" \
  -d '{"email":"cliente@upse.edu.ec","password":"123456"}'
```

```json
{"token":"fake-jwt-token-456","email":"cliente@upse.edu.ec","rol":"cliente"}
```

### CV-03 — Cuerpo inválido

```bash
curl -s -i -X POST http://localhost:3000/api/login \
  -H "Content-Type: application/json" -d 'no-json'
```

```
HTTP/1.1 400 Bad Request
{"error":"Cuerpo de petición inválido"}
```

### CV-04 — Credenciales incorrectas

```bash
curl -s -i -X POST http://localhost:3000/api/login \
  -H "Content-Type: application/json" \
  -d '{"email":"admin@upse.edu.ec","password":"mala"}'
```

```
HTTP/1.1 401 Unauthorized
{"error":"Credenciales incorrectas"}
```

### CV-05 — Catálogo plano

```bash
curl -s http://localhost:3000/api/productos
```

```json
[{"id":1,"nombre":"Serum Revitalizante","precio":45,"img":"https://picsum.photos/seed/serum/150"},
 {"id":2,"nombre":"Crema Hidratante Pro","precio":32.5,"img":"https://picsum.photos/seed/crema/150"},
 {"id":3,"nombre":"Tónico Purificante","precio":28,"img":"https://picsum.photos/seed/tonico/150"},
 {"id":4,"nombre":"Mascarilla Nocturna","precio":50,"img":"https://picsum.photos/seed/mascarilla/150"}]
```

Si la respuesta deja de ser un array de primer nivel, el catálogo del frontend se rompe.

## 3. Verificación estática

```bash
cd multicatalogo-backend
go build ./...
go vet ./...
gofmt -l .
```

## 4. No-regresión del contrato congelado

- [x] `POST /api/login` → `200` con `token`, `email`, `rol` (roles `admin` y `cliente`).
- [x] `POST /api/login` con cuerpo inválido → `400`.
- [x] `POST /api/login` con credenciales malas → `401`.
- [x] `GET /api/productos` → array plano `{id, nombre, precio, img}`.
- [ ] `GET /api/red` → aún no existe (ver `tasks.md` T-06).
- [x] CORS y puerto 3000 sin cambios.

## 5. Multitenant

No aplica a la línea base: no hay `tenant_id` en ningún modelo ni consulta.
Todo dato de negocio necesita tenant antes de entrar a producción; se define en la spec 001.

## 6. Resultado

Línea base congelada y verificada. Las specs `001-...` en adelante deben pasar primero por
`spec.md: sección 6 — No-regresión` y revalidar estos mismos casos.