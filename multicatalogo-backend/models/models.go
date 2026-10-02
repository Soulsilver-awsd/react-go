// Declaramos el paquete models para agrupar las estructuras de datos de la aplicación.
package models

// LoginRequest define la estructura esperada para el cuerpo de la petición (JSON) al iniciar sesión.
type LoginRequest struct {
	// Email representa el correo del usuario; la etiqueta `json:"email"` indica cómo se mapea el JSON a esta variable.
	Email string `json:"email"`
	// Password representa la contraseña del usuario; se extrae del campo "password" del JSON entrante.
	Password string `json:"password"`
}

// Producto define la estructura de los datos de un artículo en nuestro catálogo.
type Producto struct {
	// ID es el identificador único numérico del producto.
	ID int `json:"id"`
	// Nombre es la descripción en texto del producto.
	Nombre string `json:"nombre"`
	// Precio es el costo del producto, almacenado como un número con decimales (float64).
	Precio float64 `json:"precio"`
	// Img es la URL o ruta que apunta a la fotografía o imagen del producto.
	Img string `json:"img"`
	// TenantID identifica el inquilino (tenant) al que pertenece el producto en este sistema multitenant.
	// Lleva `json:"-"` a propósito: el dato es interno y NO debe aparecer en el JSON que consume
	// el frontend, que espera exactamente las claves id, nombre, precio e img.
	TenantID string `json:"-"`
}

// Usuario representa a una persona capaz de iniciar sesión, junto con el rol que determina
// qué secciones de la interfaz puede ver.
type Usuario struct {
	// Email es el identificador de acceso del usuario y también la clave de búsqueda en el repositorio.
	Email string
	// Password es la contraseña del usuario. En una iteración posterior se reemplaza por su hash.
	Password string
	// Rol define los permisos: "admin" o "cliente". El frontend lo lee para ocultar o mostrar secciones.
	Rol string
	// TenantID indica el inquilino al que pertenece el usuario; toda consulta se filtra por este campo.
	TenantID string
}

// APIError representa un error de la API de forma estandarizada, para que todos los controllers
// construyan sus respuestas de error a partir de una misma estructura y no de mapas sueltos.
// Se define en models porque la comparten varias capas: la construyen los controllers,
// la interpretan los repositories y la serializan las respuestas.
type APIError struct {
	// Status es el código HTTP que corresponde al error (400, 401, 404, 500, etc.).
	Status int `json:"status"`
	// Message es el texto legible y estable que describe qué ocurrió, apto para mostrarse al usuario.
	Message string `json:"message"`
	// Details es un campo opcional para información adicional (por ejemplo, el campo del cuerpo
	// que falló). Al declararse con omitempty, si queda en nil la clave desaparece del JSON,
	// de modo que el payload se mantiene limpio.
	Details interface{} `json:"details,omitempty"`
}

// NewAPIError es un constructor que centraliza la creación de errores, para que ningún controller
// tenga que armar la estructura a mano ni olvidarse de algún campo.
func NewAPIError(status int, message string, details interface{}) APIError {
	// Devolvemos el error ya armado con los tres campos provistos por el llamador.
	return APIError{Status: status, Message: message, Details: details}
}

// HealthResponse define la forma exacta del payload del endpoint de salud GET /api/health.
// Vive en models para que el contrato sea verificable y no dependa de un mapa suelto en el controller.
type HealthResponse struct {
	// Status indica el estado operativo del proceso; vale "ok" cuando el servidor atiende peticiones.
	Status string `json:"status"`
	// Message es la descripción legible del estado, en español como el resto del proyecto.
	Message string `json:"message"`
}
