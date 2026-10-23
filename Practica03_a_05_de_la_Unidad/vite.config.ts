import tailwindcss from '@tailwindcss/vite'
import react from '@vitejs/plugin-react'
import { defineConfig } from 'vite'
import { VitePWA } from 'vite-plugin-pwa'

// https://vite.dev/config/
export default defineConfig({
  plugins: [
    react(),
    tailwindcss(),
    // U1-T6 · PWA: convierte la SPA en una app instalable y capaz de
    // mostrar el catálogo sin conexión.
    VitePWA({
      // 'prompt' deja el service worker anterior activo hasta que el usuario
      // acepta la actualización (toast "Nueva versión disponible").
      registerType: 'prompt',
      // El registro lo hacemos nosotros desde React (useRegisterSW), así que
      // el plugin no debe inyectar su propio script de registro.
      injectRegister: null,
      // El SW no se registra en `npm run dev`: evita que el HMR se rompa.
      // La PWA se prueba siempre contra el build (`npm run build` + preview).
      devOptions: { enabled: false },
      // Archivos de public/ que además se precachean.
      includeAssets: [
        'favicon.svg',
        'icons/icon-192.png',
        'icons/icon-512.png',
        'icons/icon-maskable-512.png',
        'icons/apple-touch-icon.png',
      ],
      // ===== Web App Manifest =====
      manifest: {
        name: 'Paulo — Tienda de cosméticos',
        short_name: 'Paulo',
        description:
          'Tienda de cosméticos con red de ventas multinivel: catálogo, carrito, checkout y comisiones.',
        lang: 'es',
        dir: 'ltr',
        // theme_color es obligatorio: sin él el navegador no permite instalar.
        theme_color: '#4f46e5',
        background_color: '#f8fafc',
        display: 'standalone',
        orientation: 'portrait',
        start_url: '/',
        scope: '/',
        categories: ['shopping', 'lifestyle', 'business'],
        icons: [
          {
            src: 'icons/icon-192.png',
            sizes: '192x192',
            type: 'image/png',
            purpose: 'any',
          },
          {
            src: 'icons/icon-512.png',
            sizes: '512x512',
            type: 'image/png',
            purpose: 'any',
          },
          {
            src: 'icons/icon-maskable-512.png',
            sizes: '512x512',
            type: 'image/png',
            purpose: 'maskable',
          },
        ],
      },
      // ===== Service Worker (workbox) =====
      workbox: {
        // Precache del shell: JS, CSS, HTML, imágenes y fuentes.
        globPatterns: ['**/*.{js,css,html,svg,png,ico,woff2}'],
        // Las rutas de la SPA (/catalogo, /producto/1...) se resuelven con
        // index.html también sin conexión.
        navigateFallback: 'index.html',
        // Pero una navegación a la API nunca debe devolver el index.html.
        navigateFallbackDenylist: [/^\/api\//],
        cleanupOutdatedCaches: true,
        clientsClaim: true,
        runtimeCaching: [
          {
            // Solo el catálogo público. El patrón es exacto: nunca captura
            // /api/admin/productos ni ningún endpoint autenticado.
            urlPattern: ({ url }) =>
              url.pathname === '/api/productos' ||
              /^\/api\/productos\/\d+$/.test(url.pathname),
            handler: 'StaleWhileRevalidate',
            options: {
              cacheName: 'api-productos',
              expiration: {
                maxEntries: 20,
                maxAgeSeconds: 60 * 60 * 24,
              },
            },
          },
          {
            // Imágenes remotas del catálogo (picsum.photos): cache first.
            urlPattern: ({ url }) => url.hostname.endsWith('picsum.photos'),
            handler: 'CacheFirst',
            options: {
              cacheName: 'img-picsum',
              cacheableResponse: { statuses: [0, 200] },
              expiration: {
                maxEntries: 60,
                maxAgeSeconds: 60 * 60 * 24 * 30,
              },
            },
          },
        ],
      },
    }),
  ],
})
