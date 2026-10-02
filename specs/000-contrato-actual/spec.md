# 000 — Contrato actual (línea base congelada)

> Estado: validada
> Esta spec **no describe features**: congela lo que el frontend React ya consume hoy.
> Sirve de referencia de no-regresión. Cualquier endpoint listado acá solo puede **crecer**,
> nunca cambiar de forma. Para agregar comportamiento, se crea una spec nueva (`001-...`).

## 1. Origen

- Frontend consumidor: `Practica03_a_05_de_la_Unidad/` (solo lectura, ya finalizado).
- Backend: `multicatalogo-backend/` — Fiber v2, puerto 3000, datos hardcodeados.
- El frontend guarda `userEmail` en `localStorage` y el `rol` desde el contexto de autenticación
  (`src/context/AuthContext.tsx`), y hace `fetch` directo a `http://localhost:3000/api/login`
  (`src/components/Login.tsx`).

## 2. Contrato congelado

| Método | Ruta             | Ubicación en el backend         | Request             | Response 200                             | Errores                |
|--------|------------------|--------------------------------|---------------------|------------------------------------------|------------------------|
| POST   | `/api/login`     | `controllers/authController.go` | `{email, password}` | `{token, email, rol}`                    | 400 parseo, 401 credenciales |
| GET    | `/api/productos` | `controllers/prodController.go` | —                   | `Producto[]` plano `{id, nombre, precio, img}` | —                 |
| GET    | `/api/red`       | pendiente de implementar       | —                   | Estructura de red MLM que ya consume el frontend (`src/components/MiRed.tsx`) | — |

### Formas exactas que el frontend espera

```json
// POST /api/login → 200
{ "token": "fake-jwt-token-123", "email": "admin@upse.edu.ec", "rol": "admin" }

// POST /api/login → 400
{ "error": "Cuerpo de petición inválido" }

// POST /api/login → 401
{ "error": "Credenciales incorrectas" }

// GET /api/productos → 200  (array plano, NO envolver en {data,total,page})
[
  { "id": 1, "nombre": "Serum Revitalizante", "precio": 45.00, "img": "https://picsum.photos/seed/serum/150" },
  { "id": 2, "nombre": "Crema Hidratante Pro", "precio": 32.50, "img": "https://picsum.photos/seed/crema/150" },
  { "id": 3, "nombre": "Tónico Purificante", "precio": 28.00, "img": "https://picsum.photos/seed/tonico/150" },
  { "id": 4, "nombre": "Mascarilla Nocturna", "precio": 50.00, "img": "https://picsum.photos/seed/mascarilla/150" }
]
```

## 3. Credenciales de prueba

| Email                    | Password | Rol      | Token               |
|--------------------------|----------|----------|---------------------|
| `admin@upse.edu.ec`      | `123456` | `admin`  | `fake-jwt-token-123`|
| `cliente@upse.edu.ec`    | `123456` | `cliente`| `fake-jwt-token-456`|

## 4. Invariantes

- `rol` ∈ {`admin`, `cliente`}. Es el campo que el frontend usa para ocultar/mostrar UI.
- `/api/productos` devuelve un **array plano**. Cambiarlo a objetorompe el catálogo.
- CORS en `main.go`: `http://localhost:5173`, `http://127.0.0.1:5173`, `http://172.17.82.108:5173`.
  Puerto 3000.
- Para endpoints **nuevos** se usa el envelope `{success, data}` / `{success, error}`
  según `AGENTS.md`; los tres endpoints de la tabla 2 mantienen su forma original.

## 5. Trabajo pendiente detectado

| Ítem                                        | Nota                                                     |
|---------------------------------------------|----------------------------------------------------------|
| `/api/red` no está registrado en `routes/routes.go` | `MiRed.tsx` ya existe en el frontend; falta el endpoint. Requiere spec propia (`001-...`). |
| Datos hardcodeados en el controller          | Requiere capa `repository` con filtro por `tenant_id`.    |