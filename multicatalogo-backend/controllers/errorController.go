// Declaramos el paquete controllers, que agrupa las funciones que manejan la lógica de negocio de las rutas.
// Aquí viven también las utilidades comunes a todos los handlers: resolución del tenant y
// serialización de errores a partir de la estructura models.APIError.
package controllers

import (
	// Importamos Fiber para acceder al contexto de la petición y al tipo de respuesta JSON.
	"github.com/gofiber/fiber/v2"

	// Importamos el paquete models para construir y serializar la estructura de error estandarizada.
	"multicatalogo-backend/models"
)

// tenantPorDefecto es el inquilino que se asume cuando la petición no envía la cabecera
// X-Tenant-ID. Existe para no romper el contrato congelado: el frontend llama a /api/productos
// sin cabeceras propias y una respuesta 401 dejaría el catálogo en blanco.
// Cuando exista un JWT real, el tenant pasará a leerse del claim del token sin tocar los controllers.
const tenantPorDefecto = "default"

// cabeceraTenant es el nombre de la cabecera HTTP que transporta el identificador de inquilino.
const cabeceraTenant = "X-Tenant-ID"

// TenantDesdeRequest resuelve el tenant de la petición a partir de la cabecera X-Tenant-ID.
// El tenant lo asigna siempre el servidor: nunca se toma del cuerpo de la petición ni de la ruta.
func TenantDesdeRequest(c *fiber.Ctx) string {
	// Leemos la cabecera de la petición entrante.
	tenantID := c.Get(cabeceraTenant)
	// Si la cabecera no vino o vino vacía, caemos al tenant por defecto para preservar la no-regresión.
	if tenantID == "" {
		// Devolvemos el tenant por defecto, de modo que el frontend actual siga funcionando sin cabeceras.
		return tenantPorDefecto
	}
	// Hay cabecera válida: ese es el tenant que rige para las consultas.
	return tenantID
}

// ResponderError escribe un error con la envoltura estándar de la API, pensada para los endpoints
// nuevos: { "success": false, "error": "<mensaje>", "details": <opcional> }.
// El status HTTP sale del propio APIError, de modo que el código y el mensaje nunca se desincronizan.
func ResponderError(c *fiber.Ctx, apiError models.APIError) error {
	// Construimos el cuerpo de la respuesta a partir de la estructura estandarizada.
	cuerpo := fiber.Map{
		"success": false,
		"error":   apiError.Message,
	}
	// El campo details es opcional: solo lo agregamos cuando el error trae información adicional.
	if apiError.Details != nil {
		// Copiamos el detalle al cuerpo para que el cliente sepa qué campo o valor provocó el fallo.
		cuerpo["details"] = apiError.Details
	}
	// Devolvemos la respuesta con el status definido en el APIError.
	return c.Status(apiError.Status).JSON(cuerpo)
}

// ResponderErrorLegado escribe un error en la forma congelada por la línea base: { "error": "<mensaje>" }.
// Se usa únicamente en los endpoints que el frontend ya consume (/api/login), donde cambiar la forma
// del cuerpo rompería el parseo del cliente. La fuente del mensaje sigue siendo models.APIError.
func ResponderErrorLegado(c *fiber.Ctx, apiError models.APIError) error {
	// Devolvemos solo el mensaje, con el status tomado del APIError para no perder el código semántico.
	return c.Status(apiError.Status).JSON(fiber.Map{"error": apiError.Message})
}
