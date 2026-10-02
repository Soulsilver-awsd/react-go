# Diseño 002 — Health check del servidor

## 1. Resumen de la decisión

Un controller propio con un único handler sin dependencias, registrado como ruta pública. El
payload se define en `models` para que el contrato no viva suelto en el handler.

## 2. Decisiones (ADR)

### ADR-005 — Modelo `HealthResponse` en `models` en vez de `fiber.Map` en el handler

- **Estado**: aceptada
- **Contexto**: la forma del payload está congelada por US-04 (`status`, `message`).
- **Decisión**: `models.HealthResponse` con tags `json:"status"` y `json:"message"`; el handler
  devuelve `c.JSON(respuesta)`.
- **Alternativas descartadas**: `fiber.Map{"status":"ok",...}` en el controller (contrato
  implícito, se rompe en silencio si alguien edita el map); string constante suelto.
- **Consecuencias**: el JSON de salud es verificable por compilación.

### ADR-006 — Handler sin dependencias ni estado

- **Estado**: aceptada
- **Contexto**: un health check que consulta dependencias se convierte en el punto de falla más
  común; hoy no hay dependencias que consultar.
- **Decisión**: `Health` responde con el estado del proceso, sin tocar `repository`, sin leer
  headers, sin I/O.
- **Alternativas descartadas**: verificar la base de datos en cada llamada (no existe todavía);
  reportar versiones (fuera de alcance).
- **Consecuencias**: si mañana hay DB, se agrega un readiness aparte (`/api/health/ready`)
  sin tocar este endpoint.

## 3. Capas

| Capa          | Responsabilidad                          | Archivos                              |
|---------------|------------------------------------------|---------------------------------------|
| `models`      | Forma del payload de salud                | `models/models.go`                    |
| `controllers` | Handler `Health`                          | `controllers/healthController.go`     |
| `routes`      | `GET /api/health`                          | `routes/routes.go`                    |

## 4. Modelo de datos

| Tipo              | Campo     | Tipo   | JSON      |
|-------------------|-----------|--------|-----------|
| `HealthResponse`  | `Status`  | string | `status`  |
| `HealthResponse`  | `Message` | string | `message` |

## 5. Multitenant

No aplica. El handler no lee `X-Tenant-ID` ni accede a datos de negocio.

## 6. Migración / Seed

Ninguna. El endpoint no depende de datos.

## 7. Riesgos

| Riesgo                                             | Mitigación                                              |
|----------------------------------------------------|---------------------------------------------------------|
| Que un consumidor asuma el envelope `{success,...}`  | Contrato documentado en `spec.md` §4 y fijado por US-04  |
| Crear dependencia con `repository` en el health check | ADR-006: handler sin estado ni I/O                      |