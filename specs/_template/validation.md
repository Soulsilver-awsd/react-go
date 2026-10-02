# Validación NNN — <Título de la feature>

> Plantilla. Copiala a `specs/NNN-slug/validation.md`.
> Este documento se llena con EVIDENCIA real en G5. Nada de "funciona", "probado" o "listo".

## 1. Criterios de aceptación

<!-- Una fila por criterio de spec.md. Un criterio sin evidencia queda PENDIENTE y bloquea el gate. -->

| Criterio | Evidencia | Estado |
|----------|-----------|--------|
| CA-001   | CV-01     | PASA   |
| CA-002   | CV-02     | FALLA  |

## 2. Evidencia de ejecución

<!-- Comandos exactos + salida observada. Si la salida es larga, recortá a la parte que prueba el criterio. -->

### CV-01 — <qué prueba>

```bash
cd multicatalogo-backend
go run . &
curl -s -X POST http://localhost:3000/api/login \
  -H "Content-Type: application/json" \
  -d '{"email":"admin@upse.edu.ec","password":"123456"}'
```

```json
{"token":"fake-jwt-token-123","email":"admin@upse.edu.ec","rol":"admin"}
```

### CV-02 — <caso de error>

```bash
curl -s -i -X POST http://localhost:3000/api/login \
  -H "Content-Type: application/json" \
  -d '{"email":"admin@upse.edu.ec","password":"incorrecta"}'
```

```
HTTP/1.1 401 Unauthorized
{"error":"Credenciales incorrectas"}
```

## 3. Verificación estática

```bash
cd multicatalogo-backend
go build ./...   # <salida>
go vet ./...     # <salida>
gofmt -l .       # <vacío = OK>
```

## 4. No-regresión del contrato congelado

<!-- Comparar contra specs/000-contrato-actual/spec.md. -->

- [ ] `POST /api/login` → `200` con `token`, `email`, `rol`. Roles: `admin`, `cliente`.
- [ ] `POST /api/login` con cuerpo inválido → `400`.
- [ ] `POST /api/login` con credenciales malas → `401`.
- [ ] `GET /api/productos` → array plano `[{id, nombre, precio, img}]`, sin envolver.
- [ ] `GET /api/red` → misma estructura que el frontend ya consume.
- [ ] CORS y puerto 3000 sin cambios.

## 5. Multitenant

- [ ] Consulta con `tenant_id` de A no devuelve datos de B.
- [ ] ID de recurso de otro tenant responde `404` (no `403`).
- [ ] `tenant_id` del body se ignora; el del request manda.

## 6. Resultado

<!-- Veredicto final de G5. Si algo quedó fuera, declaralo acá; no lo ocultes. -->

- [ ] Todos los criterios de aceptación en PASA.
- [ ] No-regresión limpio.
- [ ] Listo para archivar como `validada`.