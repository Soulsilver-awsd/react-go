// Declaramos el paquete controllers y sumamos el handler del endpoint de salud del servidor.
package controllers

import (
	// Importamos Fiber para poder responder la petición con un 200 y un cuerpo JSON.
	"github.com/gofiber/fiber/v2"

	// Importamos el paquete models, donde vive la forma exacta del payload de salud.
	"multicatalogo-backend/models"
)

// Health es la función controladora del endpoint GET /api/health.
// Responde siempre 200 con el estado del proceso: no consulta datos, no lee cabeceras y no
// depende de la base de datos, para que un fallo externo no pueda reportarse como "servidor caído".
func Health(c *fiber.Ctx) error {
	// Construimos el payload usando la estructura definida en models, de modo que el contrato
	// del endpoint queda documentado en un solo lugar.
	respuesta := models.HealthResponse{
		// Status indica que el proceso está escuchando y atendiendo peticiones.
		Status: "ok",
		// Message es la descripción legible del estado, en español como el resto del proyecto.
		Message: "Servidor Go/Fiber operativo",
	}

	// Fiber serializa la estructura a JSON y responde con el 200 por defecto.
	return c.JSON(respuesta)
}
