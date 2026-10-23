// Declaramos el paquete routes, cuya única responsabilidad es registrar los endpoints de la API.
package routes

import (
	// Importamos el framework Fiber para poder recibir y configurar la instancia principal de la aplicación.
	"github.com/gofiber/fiber/v2"
	// Importamos el middleware de autenticación (RequireAuth / RequireAdmin).
	"multicatalogo-backend/middleware"
	// Importamos nuestro paquete de controladores, que ya llegan con sus dependencias inyectadas.
	"multicatalogo-backend/controllers"
)

// SetupRoutes registra los endpoints en la aplicación Fiber usando los controladores inyectados.
// Toda ruta de escritura o de datos privados nace protegida (seguridad por defecto).
func SetupRoutes(app *fiber.App, c *controllers.Controllers) {
	// Secreto de firma JWT (se inyecta una sola vez en el middleware).
	secret := middleware.Secret()

	// Autenticación y registro (públicas).
	app.Post("/api/login", c.Auth.Login)
	app.Post("/api/register", c.Register.Register)

	// Catálogo de productos (público de lectura).
	app.Get("/api/productos", c.Producto.GetProductos)
	app.Get("/api/productos/:id", c.Producto.GetProductoPorID)

	// Back-office de productos (solo admin): crear, editar y borrado lógico.
	admin := middleware.RequireAdmin()
	app.Get("/api/admin/productos", middleware.RequireAuth(secret), admin, c.Producto.GetProductosAdmin)
	app.Post("/api/productos", middleware.RequireAuth(secret), admin, c.Producto.CrearProducto)
	app.Put("/api/productos/:id", middleware.RequireAuth(secret), admin, c.Producto.ActualizarProducto)
	app.Delete("/api/productos/:id", middleware.RequireAuth(secret), admin, c.Producto.DesactivarProducto)

	// Red multinivel (requiere token: cada usuario solo debe ver datos autenticados).
	app.Get("/api/red", middleware.RequireAuth(secret), c.Red.GetRed)
	// Subárbol de la red del usuario autenticado (no el global del seed).
	app.Get("/api/red/mia", middleware.RequireAuth(secret), c.Red.GetMiRed)

	// Verificación de estado del servidor (pública: la usa el monitor/uptime).
	app.Get("/api/health", c.Health.Health)

	// Procesamiento de compras (requiere token: la identidad viene del JWT).
	app.Post("/api/checkout", middleware.RequireAuth(secret), c.Checkout.Checkout)

	// Pedidos: el cliente ve los suyos; el admin gestiona todos.
	app.Get("/api/pedidos/mios", middleware.RequireAuth(secret), c.Pedido.GetMisPedidos)
	app.Get("/api/pedidos", middleware.RequireAuth(secret), admin, c.Pedido.GetPedidosAdmin)
	app.Get("/api/pedidos/:id", middleware.RequireAuth(secret), c.Pedido.GetPedido)
	app.Patch("/api/pedidos/:id/estado", middleware.RequireAuth(secret), admin, c.Pedido.ActualizarEstado)

	// Ledger de comisiones del usuario autenticado.
	app.Get("/api/comisiones/mias", middleware.RequireAuth(secret), c.Comision.GetMisComisiones)

	// Liquidaciones: el admin paga y ve el historial; el distribuidor ve sus cobros.
	app.Get("/api/admin/deudas", middleware.RequireAuth(secret), admin, c.Liquidacion.GetDeudas)
	app.Post("/api/liquidaciones", middleware.RequireAuth(secret), admin, c.Liquidacion.CrearLiquidacion)
	app.Get("/api/liquidaciones", middleware.RequireAuth(secret), admin, c.Liquidacion.GetLiquidaciones)
	app.Get("/api/liquidaciones/mias", middleware.RequireAuth(secret), c.Liquidacion.GetMisLiquidaciones)

	// Reportes del dashboard (solo admin).
	app.Get("/api/reportes/resumen", middleware.RequireAuth(secret), admin, c.Reporte.GetResumen)
	app.Get("/api/reportes/top-productos", middleware.RequireAuth(secret), admin, c.Reporte.GetTopProductos)
	app.Get("/api/reportes/top-referidos", middleware.RequireAuth(secret), admin, c.Reporte.GetTopReferidos)
	app.Get("/api/reportes/comisiones-pendientes", middleware.RequireAuth(secret), admin, c.Reporte.GetComisionesPendientes)

	// Back-office de usuarios (solo admin): listar y cambiar rol/estado.
	app.Get("/api/usuarios", middleware.RequireAuth(secret), admin, c.Usuario.GetUsuarios)
	app.Patch("/api/usuarios/:id", middleware.RequireAuth(secret), admin, c.Usuario.ActualizarUsuario)
}
