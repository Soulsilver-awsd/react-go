// Declaramos el paquete controllers: reportes del dashboard (Fase 6).
package controllers

import (
	// Fiber maneja las peticiones y respuestas HTTP.
	"github.com/gofiber/fiber/v2"
	// Modelos de respuesta estandarizados.
	"multicatalogo-backend/models"
	// Repositorio con las consultas agregadas.
	"multicatalogo-backend/repository"
)

// ReporteController agrupa las dependencias de los reportes del dashboard.
type ReporteController struct {
	repo repository.DBRepository
}

// NewReporteController construye el controlador con su repositorio inyectado.
func NewReporteController(repo repository.DBRepository) *ReporteController {
	return &ReporteController{repo: repo}
}

// GetResumen (admin) devuelve los KPIs del periodo (?desde&hasta opcionales).
func (rc *ReporteController) GetResumen(c *fiber.Ctx) error {
	desde := c.Query("desde")
	hasta := c.Query("hasta")
	resumen, err := rc.repo.ObtenerResumenReporte(desde, hasta)
	if err != nil {
		return c.Status(500).JSON(models.APIError{Status: 500, Message: "Error al generar el resumen", Details: models.DetalleError(err)})
	}
	return c.JSON(resumen)
}

// GetTopProductos (admin) devuelve los productos más vendidos del periodo.
func (rc *ReporteController) GetTopProductos(c *fiber.Ctx) error {
	top, err := rc.repo.ObtenerTopProductos(c.Query("desde"), c.Query("hasta"))
	if err != nil {
		return c.Status(500).JSON(models.APIError{Status: 500, Message: "Error al generar el ranking de productos", Details: models.DetalleError(err)})
	}
	return c.JSON(top)
}

// GetTopReferidos (admin) devuelve los nodos de la red con más ventas.
func (rc *ReporteController) GetTopReferidos(c *fiber.Ctx) error {
	top, err := rc.repo.ObtenerTopReferidos()
	if err != nil {
		return c.Status(500).JSON(models.APIError{Status: 500, Message: "Error al generar el ranking de referidos", Details: models.DetalleError(err)})
	}
	return c.JSON(top)
}

// GetComisionesPendientes (admin) devuelve la deuda por distribuidor.
func (rc *ReporteController) GetComisionesPendientes(c *fiber.Ctx) error {
	deudas, err := rc.repo.DeudasPorDistribuidor()
	if err != nil {
		return c.Status(500).JSON(models.APIError{Status: 500, Message: "Error al obtener las comisiones pendientes", Details: models.DetalleError(err)})
	}
	return c.JSON(deudas)
}
