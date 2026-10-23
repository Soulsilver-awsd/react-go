// Declaramos el paquete controllers: liquidaciones a distribuidores (Fase 5).
package controllers

import (
	// Errors compara los errores centinela del repositorio.
	"errors"

	// Fiber maneja las peticiones y respuestas HTTP.
	"github.com/gofiber/fiber/v2"
	// Modelos de respuesta estandarizados.
	"multicatalogo-backend/models"
	// Repositorio con la transacción de liquidación.
	"multicatalogo-backend/repository"
)

// LiquidacionController agrupa las dependencias de las liquidaciones.
type LiquidacionController struct {
	repo repository.DBRepository
}

// NewLiquidacionController construye el controlador con su repositorio inyectado.
func NewLiquidacionController(repo repository.DBRepository) *LiquidacionController {
	return &LiquidacionController{repo: repo}
}

// GetDeudas (admin) lista lo que se le debe a cada distribuidor (comisiones pendientes).
func (lc *LiquidacionController) GetDeudas(c *fiber.Ctx) error {
	deudas, err := lc.repo.DeudasPorDistribuidor()
	if err != nil {
		return c.Status(500).JSON(models.APIError{Status: 500, Message: "Error al obtener las deudas", Details: models.DetalleError(err)})
	}
	return c.JSON(deudas)
}

// CrearLiquidacion (admin) paga todas las comisiones pendientes de un distribuidor.
func (lc *LiquidacionController) CrearLiquidacion(c *fiber.Ctx) error {
	var body struct {
		// *int para distinguir "omitido" (nil) de "0" (nodo raíz "Tú").
		ReferidoID *int `json:"referidoId"`
	}
	if err := c.BodyParser(&body); err != nil {
		return c.Status(400).JSON(models.APIError{Status: 400, Message: "Cuerpo de petición inválido"})
	}
	if body.ReferidoID == nil || *body.ReferidoID < 0 {
		return c.Status(400).JSON(models.APIError{Status: 400, Message: "El campo 'referidoId' es obligatorio"})
	}

	liquidacion, err := lc.repo.PagarLiquidacion(*body.ReferidoID)
	if err != nil {
		if errors.Is(err, repository.ErrSinComisiones) {
			return c.Status(400).JSON(models.APIError{Status: 400, Message: "El distribuidor no tiene comisiones pendientes"})
		}
		return c.Status(500).JSON(models.APIError{Status: 500, Message: "Error al crear la liquidación", Details: models.DetalleError(err)})
	}
	return c.Status(201).JSON(liquidacion)
}

// GetLiquidaciones (admin) devuelve el historial global de pagos.
func (lc *LiquidacionController) GetLiquidaciones(c *fiber.Ctx) error {
	liquidaciones, err := lc.repo.ObtenerLiquidaciones()
	if err != nil {
		return c.Status(500).JSON(models.APIError{Status: 500, Message: "Error al obtener las liquidaciones", Details: models.DetalleError(err)})
	}
	return c.JSON(liquidaciones)
}

// GetMisLiquidaciones devuelve los cobros recibidos por el distribuidor autenticado.
func (lc *LiquidacionController) GetMisLiquidaciones(c *fiber.Ctx) error {
	usuarioID, err := usuarioIDDeCtx(c)
	if err != nil {
		return c.Status(401).JSON(models.APIError{Status: 401, Message: "Token inválido"})
	}
	liquidaciones, err := lc.repo.ObtenerLiquidacionesPorUsuario(usuarioID)
	if err != nil {
		return c.Status(500).JSON(models.APIError{Status: 500, Message: "Error al obtener tus liquidaciones", Details: models.DetalleError(err)})
	}
	return c.JSON(liquidaciones)
}
