# Specs — Arnés de Spec-Driven Development (SDD)

Una **spec** es la unidad de trabajo del proyecto. Cada feature/épica vive en su propia carpeta
`specs/NNN-slug-kebab/` y contiene cuatro documentos con una responsabilidad cada uno.

La regla que gobierna todo el arnés:

> **El código nunca es la fuente de verdad. `spec.md` lo es.**
> Si al implementar descubres un requisito que no está en la spec, no lo improvises:
> vuelve a G1, actualiza la spec, recién entonces implementa.

## Estructura

```
specs/
├── README.md                # este archivo: protocolo, estados y reglas
├── _template/               # plantilla base (se copia, nunca se edita)
│   ├── spec.md              # QUÉ: requisitos y criterios de aceptación
│   ├── design.md            # CÓMO: decisiones técnicas y plan de datos
│   ├── tasks.md             # PASOS: tareas atómicas por capa
│   └── validation.md        # PRUEBA: evidencia de que los criterios se cumplen
├── 000-contrato-actual/     # línea base: contrato congelado que el frontend ya consume
└── NNN-<slug>/              # una carpeta por feature
```

## Ciclo de vida y compuertas (gates)

| Gate | Nombre      | Entregable            | Regla de salida                                                                    |
|------|-------------|-----------------------|------------------------------------------------------------------------------------|
| G0   | Contexto    | —                     | Leíste `AGENTS.md`, la guía y `Practica03_a_05_de_la_Unidad/src/` (solo lectura). Sabés qué consume el frontend. |
| G1   | Spec        | `spec.md`             | Requisitos EARS, contrato JSON, casos de error y criterios de aceptación **sin una línea de código**. |
| G2   | Diseño      | `design.md`           | Decisiones con alternativas descartadas, capa `repository`, modelo de datos y `tenant_id`. |
| G3   | Tareas      | `tasks.md`            | Tareas atómicas en orden `models → repository → controllers → routes`, cada una con su verificación. |
| G4   | Código      | `.go`                 | Cada tarea marcada `[x]` en `tasks.md` solo cuando compila y su verificación pasa. |
| G5   | Validación  | `validation.md`       | Evidencia real (comandos + salida), tabla de criterios y no-regresión del contrato congelado. |

Reglas transversales:

- No se salta una compuerta. Si G3 no existe, G4 no arranca.
- G1 es revisable: cambiar el alcance reabre la spec (bump de revisión, no reescritura silenciosa).
- Todo cambio de un endpoint existente exige actualizar primero `spec.md` y `000-contrato-actual/`.
- La spec se considera **validada** solo con evidencia en `validation.md`, no con "funciona en mi máquina".

## Estados

`borrador` → `en-revisión` → `aprobada` → `en-implementación` → `validada` → (archivada)

El estado va como primera línea de `spec.md`.

## Nomenclatura

- Carpeta: `NNN-slug-kebab-case`, `NNN` incremental de tres dígitos, starting en `001`.
- Una spec por feature/épica, no por archivo ni por endpoint suelto.
- Los requisitos se numeran `REQ-001` y los criterios `CA-001`; las decisiones `ADR-001`.
  Las tareas referencian el requisito que implementan (`T-04 → REQ-003`).

## Índice de specs

| Spec                                          | Estado   | Descripción                                                        |
|-----------------------------------------------|----------|--------------------------------------------------------------------|
| `000-contrato-actual`                         | validada | Contrato congelado que consume el frontend (línea base, no se altera). |
| `001-errores-estandar-y-validacion`           | validada | `APIError` reutilizable, validación de login y de `id` de producto (US-01..US-03). |
| `002-health-check`                            | validada | `GET /api/health` y su registro en rutas (US-04, US-05).             |
| `003-red-mlm`                                 | —        | Pendiente: `GET /api/red`, que el frontend ya consume y el backend no expone. |
| `004-tenant-y-persistencia`                   | —        | Pendiente: reemplazar el seed en memoria por SQL con filtro por `tenant_id`. |

## Notas de ejecución

- Verificar siempre desde la raíz del backend: `cd multicatalogo-backend && go build ./... && go vet ./... && gofmt -l .`
- `gofmt -l .` debe imprimir **nada**. Ojo: el repositorio original no tenía salto de línea
  final en sus archivos, así que esa comprobación fallaba antes de la spec 001.
- Para la no-regresión basta con los `curl` de `specs/000-contrato-actual/validation.md`.