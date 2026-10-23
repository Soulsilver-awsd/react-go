# MultiCatálogo — SPA (React + Vite + Tailwind + PWA)

Tienda de cosméticos con red de ventas multinivel (MLM).
Stack: React 19 · TypeScript · Vite 8 · Tailwind CSS 4 · React Router 7 · **PWA** (`vite-plugin-pwa` + Workbox).

## Comandos

```bash
npm install
npm run dev       # SPA en http://localhost:5173 (sin service worker)
npm run build     # genera dist/ + sw.js + manifest.webmanifest
npm run preview   # sirve dist/ con fallback SPA (http://localhost:4173)
npm run lint      # eslint (0 errores, 0 warnings)
```

La API (Go/Fiber) debe estar en `:3000`; se configura con `VITE_API_URL` en `.env`.

## Estructura

| Ruta (`src/`) | Responsabilidad |
|---|---|
| `App.tsx` | Router, guardas `ProtectedRoute`/`AdminRoute`, `CartBoundary`. |
| `context/` | `AuthContext` (sesión JWT) y `CartContext` (carrito por usuario). |
| `services/` | Capa de servicios: `http.ts` (cliente único `apiFetch`), `productosService`, `redService`, `pedidosService`, `reportesService`, `liquidacionesService`, `usuariosService`. |
| `components/` | Vistas (`Storefront`, `Catalogo`, `DetalleProducto`, `Carrito`, `Checkout`, `Dashboard`, back-office…) y componentes transversales (`Navbar`, `Sidebar`, `EstadoBadge`, `EstadoError`). |
| `pwa/` | Estado PWA: `pwa.ts` (tipos + contexto + hook) y `PwaProvider.tsx` (registro del SW, prompt de instalación, estado de red). |
| `data/` | Tipos y categorías (los datos vienen de la API). |

## PWA (U1-T6)

La app es instalable y muestra el **catálogo sin conexión**.

| Pieza | Dónde |
|---|---|
| Web App Manifest | `vite.config.ts` → `VitePWA({ manifest })` → `dist/manifest.webmanifest` (el `<link rel="manifest">` lo inyecta el plugin) |
| Service worker | `vite.config.ts` → `workbox` → `dist/sw.js` (precache del shell + `navigateFallback: index.html`) |
| Iconos | `public/icons/` (`icon-192`, `icon-512`, `icon-maskable-512`, `apple-touch-icon`, PNG generados desde `logo.svg`) |
| Registro y estado | `src/pwa/PwaProvider.tsx` (envuelto en `main.tsx`) |
| UI | `OfflineBanner` (banner "Sin conexión"), `PwaInstallButton` (Navbar), `PwaUpdateToast` ("Nueva versión disponible") |

### Cómo probarla

El service worker **solo** se registra en el build de producción (`devOptions.enabled = false` para no romper el HMR):

```bash
npm run build
npm run preview          # http://localhost:4173
```

1. Abre `http://localhost:4173` y regístrate/inicia sesión.
2. DevTools → **Application → Manifest**: debe aparecer "MultiCatálogo — Tienda de cosméticos" con los 3 iconos y *Installable: yes*.
3. DevTools → **Application → Service Workers**: estado *activated and running*; en **Cache Storage** se ven `workbox-precache-*`, `api-productos` e `img-picsum`.
4. Instala la app desde el banner de la barra de direcciones o con el botón **Instalar app** de la Navbar (en iOS no hay instalación automática: el botón explica "Compartir → Añadir a pantalla de inicio").
5. Corta la red (DevTools → Network → *Offline*) y recarga: la app abre, el catálogo se lista desde la caché y aparece el banner "Sin conexión".

### Qué se cachea (y qué no)

| Recurso | Estrategia |
|---|---|
| Shell (JS/CSS/HTML/iconos) | Precache (Workbox) |
| `GET /api/productos`, `GET /api/productos/:id` | `StaleWhileRevalidate` (1 día) |
| Imágenes `picsum.photos` | `CacheFirst` (30 días, 60 entradas) |
| `/api/red*`, `/api/pedidos*`, `/api/usuarios*`, `/api/reportes*`, `/api/liquidaciones*` | **sin caché** (datos autenticados) |
| `POST/PUT/PATCH/DELETE` (checkout, back-office) | **sin caché** (escrituras) |

> Motivo de seguridad: nada autenticado se guarda en el dispositivo, así que una PWA instalada en un equipo compartido no filtra datos entre cuentas.

### Límites conocidos

- El checkout y las vistas de back-office **requieren conexión** (el servidor valida stock, precios y comisiones).
- La instalación exige **HTTPS o `localhost`**; en pruebas por IP hace falta un túnel con TLS.
- Las actualizaciones del SW se aplican a petición del usuario (aviso "Nueva versión disponible").
- El proxy inverso con SSL (Caddy/Nginx) es parte de la Unidad 4.
