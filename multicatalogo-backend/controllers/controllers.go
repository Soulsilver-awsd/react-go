// Declaramos el paquete controllers para agrupar todos los controladores de la API.
package controllers

// Controllers agrupa las instancias de todos los controladores ya inyectados,
// para que la capa de rutas reciba un único objeto con todas las dependencias resueltas.
type Controllers struct {
	// Auth maneja el inicio de sesión.
	Auth *AuthController
	// Register maneja la creación de cuentas de cliente.
	Register *RegisterController
	// Producto maneja el catálogo.
	Producto *ProductoController
	// Red maneja el árbol multinivel de referidos.
	Red *RedController
	// Health expone el estado del servidor.
	Health *HealthController
	// Checkout procesa las compras.
	Checkout *CheckoutController
	// Pedido gestiona el ciclo de vida de los pedidos.
	Pedido *PedidoController
	// Comision expone el ledger MLM del usuario.
	Comision *ComisionController
	// Liquidacion gestiona los pagos a distribuidores.
	Liquidacion *LiquidacionController
	// Reporte expone los KPIs del dashboard.
	Reporte *ReporteController
	// Usuario gestiona las cuentas (back-office).
	Usuario *UsuarioController
}
