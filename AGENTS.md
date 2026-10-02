# Contexto del Proyecto: Motor Backend para Billetera/Catálogo (Unidad 2)

Actúas como un desarrollador backend Senior en Go. Estamos construyendo el core transaccional de un sistema SaaS Multitenant usando Go y Fiber.

El frontend en React (Vite + Tailwind) ya está finalizado y consume la API en el puerto 3000. Tu objetivo estricto es modificar SOLO el backend (`multicatalogo-backend/`).

## Reglas Arquitectónicas

1. **Separación de Responsabilidades:** Mantén el patrón de diseño actual. `routes.go` define los endpoints, `controllers` maneja la lógica HTTP (req/res), `models` define las estructuras, y `repository` maneja los datos.
2. **Estandarización de Respuestas:** Todas las respuestas HTTP deben seguir una estructura JSON predecible.
3. **Estabilidad del Frontend:** Está estrictamente prohibido cambiar las firmas de los endpoints existentes (`/api/login`, `/api/productos`, `/api/red`) o la estructura de los JSON que el frontend ya espera.

## Estructura del Backend

```
multicatalogo-backend/
├── go.mod                  # module multicatalogo-backend (Go 1.26, Fiber v2)
├── main.go                 # package main: app Fiber, CORS, puerto 3000
├── routes/
│   └── routes.go           # SetupRoutes(app *fiber.App) — registro de endpoints
├── controllers/
│   ├── authController.go   # Login
│   └── prodController.go   # GetProductos
├── models/
│   └── models.go           # LoginRequest, Producto (tags `json:"..."`)
└── repository/             # (a crear) acceso a datos — sin dependencia de Fiber
```

Dependencias actuales: **solo** `github.com/gofiber/fiber/v2` + `middleware/cors`. No introduzcas dependencias nuevas sin necesidad; si hace falta persistencia, justifica la elección (driver SQL, ORM) antes de añadirla a `go.mod`.

## Contrato de la API (no romper)

| Método | Ruta              | Request                  | Response 200                                              | Errores                       |
|--------|-------------------|--------------------------|-----------------------------------------------------------|-------------------------------|
| POST   | `/api/login`      | `{email, password}`      | `{token, email, rol}`                                      | 400 parseo, 401 credenciales  |
| GET    | `/api/productos`  | —                        | `Producto[]` = `{id, nombre, precio, img}[]` (array plano)  | —                             |
| GET    | `/api/red`        | —                        | Estructura de red MLM que el frontend ya consume           | —                             |

Reglas de compatibilidad:
- `/api/productos` debe seguir devolviendo un **array plano** de objetos con exactamente `id`, `nombre`, `precio`, `img`. No lo envuelvas en `{data, total, page}`.
- `/api/login` debe devolver `token`, `email` y `rol` (`"admin"` / `"cliente"`). El frontend guarda `userEmail` y lee el rol del contexto de autenticación.
- No cambies puertos (3000), rutas, verbos HTTP ni nombres de campos. Si necesitas datos adicionales, **agrégalos**; nunca renombres ni elimines los existentes.
- CORS en `main.go` permite `localhost:5173`, `127.0.0.1:5173` y la IP del aula. Si cambias orígenes, confirma primero.

## Estandarización de Respuestas

Toda respuesta debe ser predecible. Para endpoints nuevos usa un envelope consistente con códigos HTTP semánticos:

```json
{ "success": true, "data": { } }
{ "success": false, "error": "mensaje legible" }
```

- Éxito de lectura: `200`. Recurso creado: `201`. Petición mal formada / validación: `400`. Sin credenciales / token inválido: `401`. Sin permisos: `403`. No encontrado: `404`. Conflicto (duplicado, saldo insuficiente): `409`. Error interno: `500` (nunca filtres stack traces ni credenciales).
- El campo `error` es texto legible y estable; no lo consumas en el frontend para lógica, solo para mostrar.
- Todas las fechas en ISO-8601 UTC. Todos los montos con dos decimales y la misma representación en toda la API.

## Multitenant

El sistema es SaaS multitenant: todo dato de negocio (usuarios, productos, billetera, transacciones, red) pertenece a un tenant.
- Identifica el tenant en el request (header `X-Tenant-ID` o claim del token) y **filtra siempre por `tenant_id`** en la capa `repository`. Ninguna consulta puede devolver datos de otro tenant.
- El `tenant_id` lo asigna el servidor; nunca lo tomes del body de la petición.
- Impide el acceso cruzado entre tenants también en las validaciones de recursos por ID (un `producto_id` de otro tenant responde `404`, no `403`, para no filtrar existencia).

## Flujo de trabajo

1. Lee `Guia_Practica_Tema5_leccion.md` para el contexto de la unidad; el frontend está en `Practica03_a_05_de_la_Unidad/` y es **solo lectura**.
2. Antes de escribir código, revisa `Practica03_a_05_de_la_Unidad/src/` (especialmente `components/`, `services/`, `context/`) para confirmar qué campos y endpoints consume el frontend.
3. Implementa de abajo hacia arriba: `models` → `repository` → `controllers` → `routes`.
4. Los controllers no tocan datos directamente; los repositories no conocen `fiber.Ctx`.
5. Mantén el estilo del código existente: español en nombres, tipos y comentarios; comentarios explicativos en cada declaración (el código actual lo hace así).

## Arnés SDD (obligatorio)

El trabajo se dirige con specs. **Ninguna tarea arranca sin su `spec.md`**, y el código nunca es la fuente de verdad: si al implementar aparece un requisito que no está en la spec, volvés a G1 y actualizás la spec antes de seguir codificando.

```
specs/
├── README.md              # protocolo, estados y reglas del arnés
├── _template/             # plantilla base (se copia, nunca se edita)
│   ├── spec.md            # QUÉ  — requisitos EARS, contrato JSON, criterios de aceptación
│   ├── design.md          # CÓMO — ADRs, capas, modelo de datos, multitenant
│   ├── tasks.md           # PASOS — tareas atómicas models→repository→controllers→routes
│   └── validation.md      # PRUEBA — evidencia real (comandos + salida) de cada criterio
└── NNN-slug-kebab/        # una spec por feature/épica
```

Compuertas: **G0 Contexto** → **G1 Spec** → **G2 Diseño** → **G3 Tareas** → **G4 Código** → **G5 Validación**.

- **G0** Leíste esta guía, la guía de la unidad y el frontend.
- **G1** `spec.md` completo: requisitos numerados (`REQ-00X`) con sintaxis EARS, contrato JSON de cada endpoint tocado, tabla de errores y criterios `CA-00X`. Sin código.
- **G2** `design.md` con decisiones y alternativas descartadas; incluye cómo entra el `tenant_id`.
- **G3** `tasks.md` con tareas atómicas, cada una referenciando su requisito y su verificación.
- **G4** Solo se marca `[x]` en `tasks.md` cuando la tarea compila y su verificación pasa.
- **G5** `validation.md` con comandos, salida real y tabla de criterios. Una spec está `validada` solo con evidencia escrita, nunca con "funciona".

Reglas transversales:

- `specs/000-contrato-actual/` congela los endpoints que el frontend ya consume. Toda feature debe pasar primero por su sección **No-regresión** y revalidar esos casos al terminar.
- No se salta una compuerta. Si `tasks.md` no existe, no se escribe código.
- Cambiar un endpoint existente obliga a actualizar primero `spec.md` y la línea base, nunca a "ajustar sobre la marcha".
- Nomenclatura: `specs/NNN-slug-kebab-case/`, `NNN` incremental desde `001`; estados `borrador → en-revisión → aprobada → en-implementación → validada` en la primera línea de `spec.md`.
- Si el trabajo no cabe en un cambio chico y está claramente especificado (p. ej. una corrección de un typo en un dato), no levantes una spec completa; mantené el criterio y dejalo registrado en el commit/PR.

## Verificación

```bash
cd multicatalogo-backend
go build ./...
go vet ./...
gofmt -l .        # no debe imprimir nada
```

Prueba los endpoints a mano antes de dar por terminada una tarea:

```bash
go run . &
curl -s -X POST http://localhost:3000/api/login \
  -H "Content-Type: application/json" \
  -d '{"email":"admin@upse.edu.ec","password":"123456"}'
curl -s http://localhost:3000/api/productos
```

## Fuera de alcance

- No modifiques `Practica03_a_05_de_la_Unidad/` (frontend React) ni `Guia_Practica_Tema5_leccion.md`.
- No cambies la configuración de CORS, el puerto 3000 ni las credenciales de prueba sin avisar.
- No elimines código existente commented-out o "legacy": engánchalo al patrón nuevo en lugar de borrarlo.