// Declaramos el paquete models para agrupar las estructuras de datos de la aplicación.
package models

import (
	"os"
	"strings"
)

// DetalleError devuelve el detalle técnico de un error SOLO en desarrollo
// (APP_ENV=dev). En producción el campo "details" viaja vacío y el cliente
// recibe únicamente el mensaje sanitizado: no se filtran errores internos
// (SQL, rutas de archivos, etc.) al consumidor de la API.
func DetalleError(err error) interface{} {
	if err == nil {
		return nil
	}
	if strings.EqualFold(os.Getenv("APP_ENV"), "dev") {
		return err.Error()
	}
	return nil
}

// LoginRequest define la estructura esperada para el cuerpo de la petición (JSON) al iniciar sesión.
type LoginRequest struct {
	// Email representa el correo del usuario; la etiqueta `json:"email"` indica cómo se mapea el JSON a esta variable.
	Email string `json:"email"`
	// Password representa la contraseña del usuario; se extrae del campo "password" del JSON entrante.
	Password string `json:"password"`
}

// RegisterRequest define el cuerpo JSON de POST /api/register.
type RegisterRequest struct {
	Nombre   string `json:"nombre"`
	Email    string `json:"email"`
	Password string `json:"password"`
	// Ref es el código de invitación opcional del referente (/registro?ref=CODIGO).
	Ref string `json:"ref"`
}

// Usuario representa un registro de la tabla `usuarios` de PostgreSQL.
type Usuario struct {
	ID       int    `json:"id"`
	Email    string `json:"email"`
	Password string `json:"-"`
	Rol      string `json:"rol"`
	Nombre   string `json:"nombre"`
	Activo   bool   `json:"activo"`
	// ReferidoID es el nodo de la red asociado al usuario (fuente de comisiones).
	ReferidoID *int `json:"referido_id,omitempty"`
}

// APIError estandariza los mensajes de error HTTP de la API en una estructura JSON predecible.
type APIError struct {
	// Status es el código HTTP asociado al error (ej. 400, 401, 404).
	Status int `json:"status"`
	// Message describe en lenguaje humano la causa del error.
	Message string `json:"message"`
	// Details contiene información adicional opcional sobre el error (puede ser nil).
	Details interface{} `json:"details,omitempty"`
}

// Producto define la estructura de los datos de un artículo en nuestro catálogo.
type Producto struct {
	// ID es el identificador único numérico del producto.
	ID int `json:"id"`
	// Nombre es la descripción en texto del producto.
	Nombre string `json:"nombre"`
	// Descripcion amplía la información del producto.
	Descripcion string `json:"descripcion"`
	// Precio es el costo del producto, almacenado como un número con decimales (float64).
	Precio float64 `json:"precio"`
	// Categoria clasifica el producto dentro del catálogo.
	Categoria string `json:"categoria"`
	// Img es la URL o ruta que apunta a la fotografía o imagen principal del producto.
	Img string `json:"img"`
	// Galeria contiene imágenes adicionales para el detalle del producto.
	Galeria []string `json:"galeria"`
	// Stock es la cantidad disponible simulada del producto.
	Stock int `json:"stock"`
	// Activo habilita el borrado lógico (false = retirado del catálogo).
	Activo bool `json:"activo"`
}

// Referido representa un nodo del árbol multinivel de referidos.
type Referido struct {
	ID     int     `json:"id"`
	Nombre string  `json:"nombre"`
	Nivel  int     `json:"nivel"`
	Ventas float64 `json:"ventas"`
	// Saldo es el acumulado de comisiones ganadas por este referido.
	Saldo float64 `json:"saldo"`
	// Codigo es el código de invitación único (/registro?ref=CODIGO).
	Codigo string `json:"codigo,omitempty"`
	// Hijos contiene los referidos del nivel inmediatamente inferior.
	Hijos []Referido `json:"hijos,omitempty"`
}

// CheckoutItem representa un producto dentro del carrito enviado al checkout.
type CheckoutItem struct {
	ID       int     `json:"id"`
	Nombre   string  `json:"nombre"`
	Cantidad int     `json:"cantidad"`
	Precio   float64 `json:"precio"`
}

// CheckoutCliente contiene los datos de envío del comprador.
type CheckoutCliente struct {
	Nombre string `json:"nombre"`
	Email  string `json:"email"`
	Ciudad string `json:"ciudad"`
}

// CheckoutRequest es el cuerpo JSON que envía el componente Checkout.tsx.
type CheckoutRequest struct {
	Items      []CheckoutItem  `json:"items"`
	Total      float64         `json:"total"`
	Cliente    CheckoutCliente `json:"cliente"`
	MetodoPago string          `json:"metodoPago"`
}

// Pedido representa la cabecera de un pedido confirmado.
type Pedido struct {
	ID            int     `json:"id"`
	Numero        string  `json:"numero"`
	ClienteNombre string  `json:"clienteNombre"`
	ClienteEmail  string  `json:"clienteEmail"`
	Ciudad        string  `json:"ciudad"`
	MetodoPago    string  `json:"metodoPago"`
	Total         float64 `json:"total"`
	Estado        string  `json:"estado"`
	// UsuarioID es el comprador autenticado (identidad del token, no del body).
	UsuarioID *int `json:"usuarioId,omitempty"`
	// CreadoEn es la fecha de creación en formato RFC3339.
	CreadoEn string `json:"creadoEn,omitempty"`
	// Items es el detalle; solo se rellena en la consulta de un pedido concreto.
	Items []PedidoItem `json:"items,omitempty"`
}

// PedidoItem es una línea del detalle de un pedido (con nombre de producto).
type PedidoItem struct {
	ProductoID int     `json:"productoId"`
	Nombre     string  `json:"nombre"`
	Cantidad   int     `json:"cantidad"`
	Precio     float64 `json:"precio"`
}

// Pagina es el sobre estándar de listados paginados (Fase 7).
type Pagina[T any] struct {
	Data  []T `json:"data"`
	Page  int `json:"page"`
	Limit int `json:"limit"`
	Total int `json:"total"`
}

// Comision es una fila del ledger MLM (quién cobró qué por qué pedido).
type Comision struct {
	ID         int     `json:"id"`
	PedidoID   int     `json:"pedidoId"`
	ReferidoID int     `json:"referidoId"`
	Nivel      int     `json:"nivel"`
	Porcentaje float64 `json:"porcentaje"`
	Monto      float64 `json:"monto"`
	Estado     string  `json:"estado"`
	CreadoEn   string  `json:"creadoEn,omitempty"`
}

// Liquidacion es un pago agrupado de comisiones a un distribuidor (Fase 5).
type Liquidacion struct {
	ID         int     `json:"id"`
	ReferidoID int     `json:"referidoId"`
	Total      float64 `json:"total"`
	Estado     string  `json:"estado"`
	CreadoEn   string  `json:"creadoEn,omitempty"`
	PagadaEn   string  `json:"pagadaEn,omitempty"`
	// ReferidoNombre es el nombre del distribuidor (solo en el historial admin).
	ReferidoNombre string `json:"referidoNombre,omitempty"`
}

// DeudaDistribuidor resume lo que se le debe a un distribuidor (comisiones pendientes).
type DeudaDistribuidor struct {
	ReferidoID int     `json:"referidoId"`
	Nombre     string  `json:"nombre"`
	Pendiente  float64 `json:"pendiente"`
}

// ResumenReporte agrega los KPIs del negocio para el dashboard (Fase 6).
type ResumenReporte struct {
	// TotalVendido es la suma de pedidos no cancelados del periodo.
	TotalVendido float64 `json:"totalVendido"`
	// NumPedidos es el número de pedidos no cancelados del periodo.
	NumPedidos int `json:"numPedidos"`
	// TicketPromedio es TotalVendido / NumPedidos (0 si no hay pedidos).
	TicketPromedio float64 `json:"ticketPromedio"`
	// ComisionesDevengadas suma las comisiones no revertidas del periodo.
	ComisionesDevengadas float64 `json:"comisionesDevengadas"`
}

// TopProducto es un producto rankeado por ventas (cantidad y dinero).
type TopProducto struct {
	ProductoID int     `json:"productoId"`
	Nombre     string  `json:"nombre"`
	Cantidad   int     `json:"cantidad"`
	Total      float64 `json:"total"`
}

// TopReferido es un nodo de la red rankeado por ventas de su subárbol.
type TopReferido struct {
	ReferidoID int     `json:"referidoId"`
	Nombre     string  `json:"nombre"`
	Codigo     string  `json:"codigo,omitempty"`
	Ventas     float64 `json:"ventas"`
	// Saldo es el acumulado de comisiones del nodo.
	Saldo float64 `json:"saldo"`
}
