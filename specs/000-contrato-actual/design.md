# Diseño 000 — Contrato actual (línea base)

> Plantilla. Copiala a `specs/NNN-slug/design.md` junto a `spec.md`.
> Este documento responde CÓMO se implementa lo que `spec.md` define. No redefine requisitos.

## 1. Resumen de la decisión

<!-- 3-6 líneas: qué se va a construir y en qué capas. -->

## 2. Decisiones (ADR)

<!-- Una decisión por bloque. Cada ADR registra el porque, no solo el qué.
     Si algo se descartó, dejalo escrito: es lo que evita repetir el debate en otra feature. -->

### ADR-001 — <decisión>

- **Estado**: aceptada
- **Contexto**: <fuerzas en juego: restricciones del repo, del frontend, de la guía>
- **Decisión**: <qué se hace>
- **Alternativas descartadas**: <opción> — <por qué no>
- **Consecuencias**: <qué habilita y qué cuesta>

## 3. Capas

<!-- Repartí el trabajo según la separación de responsabilidades de AGENTS.md. -->

| Capa          | Responsabilidad                                  | Archivos                        |
|---------------|--------------------------------------------------|---------------------------------|
| `models`      | Estructuras y contrato JSON                        | `models/models.go`              |
| `repository`  | Acceso a datos, filtro obligatorio por `tenant_id`  | `repository/<recurso>.go`       |
| `controllers` | HTTP: parseo, validación, código de respuesta      | `controllers/<recurso>Controller.go` |
| `routes`      | Registro de endpoints                             | `routes/routes.go`              |

## 4. Modelo de datos

<!-- Tablas, columnas, tipos, índices. Marcá el índice por tenant_id si aplica. -->

| Tabla      | Columna       | Tipo      | Notas                        |
|------------|---------------|-----------|------------------------------|
| `producto` | `tenant_id`   | `uuid/int`| Siempre filtrada en consultas |
| `producto` | `id`          | —         | Único junto con `tenant_id`  |

## 5. Multitenant

- Origen del tenant en el request: <header `X-Tenant-ID` | claim del token>
- Dónde se valida: <middleware | controller | repository>
- Regla: el `tenant_id` nunca se lee del body; el servidor lo asigna.
- Consultas por ID: se resuelven como `404` si el recurso pertenece a otro tenant.

## 6. Migración / Seed

<!-- Datos de arranque que necesita la feature para ser demostrable.
     Aclará si es seed de desarrollo, no datos de producción. -->

## 7. Riesgos

| Riesgo                                            | Mitigación                                    |
|---------------------------------------------------|-----------------------------------------------|
| Romper el JSON que el frontend ya parsea         | Fixture de no-regresión en `validation.md`    |
| Fuga de datos entre tenants en una consulta nueva| Test manual con dos tenants distintos          |