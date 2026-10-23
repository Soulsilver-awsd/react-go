// Declaramos el paquete controllers, ya que este archivo maneja la simulación transaccional de compras.
package controllers

import (
	// Importamos la librería estándar 'errors' para comparar errores centinela del repositorio.
	"errors"
	// Importamos la librería estándar 'log' para registrar lo ocurrido en la goroutine asíncrona.
	"log"

	// Importamos Fiber para manejar la respuesta HTTP.
	"github.com/gofiber/fiber/v2"
	// Importamos la capa de caché (abstracción inyectada) para invalidar la red.
	"multicatalogo-backend/cache"
	// Importamos los modelos para usar CheckoutRequest y APIError.
	"multicatalogo-backend/models"
	// Importamos el repositorio, que ejecuta la transacción ACID.
	"multicatalogo-backend/repository"
)

// CheckoutController agrupa las dependencias del proceso de compra (inyección de dependencias).
type CheckoutController struct {
	// repo es la abstracción de acceso a datos inyectada.
	repo repository.DBRepository
	// cache es la abstracción de caché inyectada.
	cache cache.CacheRepository
}

// NewCheckoutController construye el controlador inyectando sus dependencias.
func NewCheckoutController(repo repository.DBRepository, cacheRepo cache.CacheRepository) *CheckoutController {
	return &CheckoutController{repo: repo, cache: cacheRepo}
}

// Checkout recibe el JSON de Checkout.tsx, registra el pedido en una transacción ACID y,
// tras confirmar, lanza una goroutine para acreditar comisiones a la red multinivel.
func (cc *CheckoutController) Checkout(c *fiber.Ctx) error {
	// Parseamos el cuerpo JSON enviado por el frontend.
	var req models.CheckoutRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(models.APIError{
			Status:  400,
			Message: "Cuerpo de petición inválido",
			Details: models.DetalleError(err),
		})
	}

	// Validamos que el carrito no venga vacío antes de abrir la transacción.
	if len(req.Items) == 0 {
		return c.Status(400).JSON(models.APIError{
			Status:  400,
			Message: "El pedido debe contener al menos un artículo",
		})
	}

	// ===== Identidad desde el TOKEN (no desde el body) =====
	// RequireAuth ya validó el JWT y publicó el email en Locals. Sobrescribimos
	// el email del body: un email falso en el payload NO se persiste ni atribuye
	// comisiones. La identidad siempre proviene del token firmado por el servidor.
	emailToken, _ := c.Locals("email").(string)
	if emailToken == "" {
		return c.Status(401).JSON(models.APIError{Status: 401, Message: "Token de autenticación requerido"})
	}
	req.Cliente.Email = emailToken

	// El comprador es el id del claim "sub" (string en el JWT → int aquí).
	usuarioID, err := usuarioIDDeCtx(c)
	if err != nil {
		return c.Status(401).JSON(models.APIError{Status: 401, Message: "Token inválido: id de usuario ausente"})
	}

	// Ejecutamos la transacción ACID en el repositorio (BEGIN / INSERT / COMMIT / ROLLBACK).
	pedido, err := cc.repo.CrearPedido(req, usuarioID)
	if err != nil {
		// Si el error fue por un producto inexistente, la transacción se revirtió (ROLLBACK).
		if errors.Is(err, repository.ErrNoEncontrado) {
			return c.Status(400).JSON(models.APIError{
				Status:  400,
				Message: "Producto inexistente: la transacción fue revertida (ROLLBACK)",
			})
		}
		// Cantidad inválida (≤ 0): error de validación del cliente, no de stock.
		if errors.Is(err, repository.ErrCantidadInvalida) {
			return c.Status(400).JSON(models.APIError{
				Status:  400,
				Message: "Cantidad inválida: debe ser mayor que cero",
			})
		}
		// Cualquier otro fallo (ej. base de datos caída) también provoca ROLLBACK.
		return c.Status(500).JSON(models.APIError{
			Status:  500,
			Message: "Error al procesar el pedido: la transacción fue revertida (ROLLBACK)",
			Details: models.DetalleError(err),
		})
	}

	// Tomamos el email del comprador para invalidar la caché de su red.
	email := req.Cliente.Email

	// El contador de ventas de la red se actualiza en segundo plano (concurrencia).
	go cc.procesarVentasAsync(pedido.Total, email, usuarioID)

	// Respondemos inmediatamente con el pedido confirmado (HTTP 201).
	return c.Status(201).JSON(pedido)
}

// procesarVentasAsync actualiza el contador de ventas de la red e invalida la caché.
// El ledger de comisiones ya quedó escrito dentro de la transacción de CrearPedido.
func (cc *CheckoutController) procesarVentasAsync(total float64, email string, usuarioID int) {
	if err := cc.repo.ActualizarVentasRed(total, usuarioID); err != nil {
		log.Printf("Error al actualizar ventas de la red: %v", err)
	}

	// Invalidamos la caché de la red para que el Dashboard muestre KPIs actualizados.
	cc.cache.Del(ClaveRed(email))
	cc.cache.DelByPrefix("red:tree:")
	log.Printf("Ventas de red actualizadas para %s (pedido de %.2f)", email, total)
}
