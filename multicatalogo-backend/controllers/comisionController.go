// Declaramos el paquete controllers: ledger de comisiones y red propia (Fase 4).
package controllers

import (
	// Fiber maneja las peticiones y respuestas HTTP.
	"github.com/gofiber/fiber/v2"
	// Modelos de respuesta estandarizados.
	"multicatalogo-backend/models"
	// Repositorio con el acceso al ledger de comisiones.
	"multicatalogo-backend/repository"
)

// ComisionController agrupa las dependencias del ledger MLM.
type ComisionController struct {
	repo repository.DBRepository
}

// NewComisionController construye el controlador con su repositorio inyectado.
func NewComisionController(repo repository.DBRepository) *ComisionController {
	return &ComisionController{repo: repo}
}

// GetMisComisiones devuelve el ledger del usuario autenticado + totales por estado.
func (cc *ComisionController) GetMisComisiones(c *fiber.Ctx) error {
	usuarioID, err := usuarioIDDeCtx(c)
	if err != nil {
		return c.Status(401).JSON(models.APIError{Status: 401, Message: "Token inválido"})
	}

	comisiones, pendiente, pagado, err := cc.repo.ObtenerComisionesPorUsuario(usuarioID)
	if err != nil {
		return c.Status(500).JSON(models.APIError{Status: 500, Message: "Error al obtener tus comisiones", Details: models.DetalleError(err)})
	}

	return c.JSON(fiber.Map{
		"comisiones": comisiones,
		"pendiente":  pendiente,
		"pagado":     pagado,
	})
}

// GetMiRed devuelve el subárbol de referidos del usuario autenticado (no el global).
func (rc *RedController) GetMiRed(c *fiber.Ctx) error {
	usuarioID, err := usuarioIDDeCtx(c)
	if err != nil {
		return c.Status(401).JSON(models.APIError{Status: 401, Message: "Token inválido"})
	}

	red, err := rc.repo.ObtenerRedDeUsuario(usuarioID)
	if err != nil {
		if err == repository.ErrNoEncontrado {
			// Un usuario sin nodo de red devuelve un árbol vacío, no un error.
			return c.JSON(models.Referido{Nombre: emailDeCtx(c)})
		}
		return c.Status(500).JSON(models.APIError{Status: 500, Message: "Error al obtener tu red", Details: models.DetalleError(err)})
	}
	return c.JSON(red)
}
