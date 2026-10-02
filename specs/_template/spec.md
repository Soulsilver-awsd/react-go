# Estado: borrador

> Reemplazá esta línea por el estado: `borrador` | `en-revisión` | `aprobada` | `en-implementación` | `validada`.
> Este archivo es la plantilla. Copiala a `specs/NNN-slug/` y editá la copia.

# Spec NNN — <Título de la feature o épica>

## 1. Contexto

<!-- Por qué existe esto. 3-5 líneas. Qué problema resuelve, quién lo consume. -->

## 2. Alcance

**Incluye**

- <qué entra>

**Excluye**

- <qué NO entra, aunque parezca tentador>

## 3. Requisitos

<!-- Numerados y verificables. Usá la sintaxis EARS:
     While   <estado>, el sistema shall <respuesta>
     When    <evento>,  el sistema shall <respuesta>
     If      <condición>, el sistema shall <respuesta>
     El sistema shall <comportamiento>            (caso general)
     Evitar "debe", "se necesita", "es fácil de usar". -->

### REQ-001 — <nombre>

<!-- Ubicación: paquete/capa → models | repository | controllers | routes -->

**Criterio de aceptación (CA-001)**: Given `<precondición>`, When `<acción>`, Then `<resultado observable>`.

## 4. Contrato de la API

<!-- Solo si la feature toca endpoints. Tabla por endpoint + ejemplo de request/response REAL.
     Respeta la convención de AGENTS.md: envelope {success, data} | {success, error}. -->

### `METHOD /api/algo`

**Request**

```json
```

**Respuesta 200**

```json
```

**Errores**

| Código | Condición                         | `error`               |
|--------|-----------------------------------|-----------------------|
| 400    | <validación que falla>            | <mensaje legible>     |
| 401    | <sin token / token inválido>      | <mensaje legible>     |
| 404    | <recurso inexistente o ajeno>     | <mensaje legible>     |

## 5. Multitenant

<!-- Cómo se resuelve el tenant (header/claim), qué filtros exige cada consulta,
     y el comportamiento ante un ID de otro tenant (404, no 403). -->

## 6. No-regresión

<!-- Endpoints o estructuras que esta feature NO puede tocar. Referencia a specs/000-contrato-actual/. -->

- [ ] `/api/productos` sigue devolviendo un array plano `{id, nombre, precio, img}`.
- [ ] `/api/login` sigue devolviendo `token`, `email`, `rol`.
- [ ] CORS y puerto 3000 sin cambios.

## 7. Fuera de alcance

- <explícito, para cerrar el perímetro de la spec>

## 8. Trazabilidad

<!-- Tabla que enlaza requisito ↔ tarea ↔ evidencia. Se completa en G3/G5. -->

| Requisito | Tarea     | Evidencia (validation.md) | Estado |
|-----------|-----------|---------------------------|--------|
| REQ-001   | T-01      | CV-01                     | —      |