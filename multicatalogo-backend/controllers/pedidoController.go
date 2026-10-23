// Declaramos el paquete controllers: gestión del ciclo de vida de los pedidos.
package controllers

import (
	// Errors compara los errores centinela del repositorio.
	"errors"
	// Strconv convierte los parámetros de la URL (id, page, limit).
	"strconv"

	// Fiber maneja las peticiones y respuestas HTTP.
	"github.com/gofiber/fiber/v2"
	// Modelos de respuesta estandarizados.
	"multicatalogo-backend/models"
	// Repositorio con la lógica transaccional de pedidos.
	"multicatalogo-backend/repository"
)

// PedidoController agrupa las dependencias de la gestión de pedidos.
type PedidoController struct {
	repo repository.DBRepository
}

// NewPedidoController construye el controlador con su repositorio inyectado.
func NewPedidoController(repo repository.DBRepository) *PedidoController {
	return &PedidoController{repo: repo}
}

// paginacion parsea ?page y ?limit con valores por defecto seguros.
func paginacion(c *fiber.Ctx) (int, int) {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	return page, limit
}

// GetMisPedidos lista los pedidos del usuario autenticado (paginado).
func (pc *PedidoController) GetMisPedidos(c *fiber.Ctx) error {
	usuarioID, err := usuarioIDDeCtx(c)
	if err != nil {
		return c.Status(401).JSON(models.APIError{Status: 401, Message: "Token inválido"})
	}
	page, limit := paginacion(c)

	pagina, err := pc.repo.ObtenerPedidosPorUsuario(usuarioID, page, limit)
	if err != nil {
		return c.Status(500).JSON(models.APIError{Status: 500, Message: "Error al obtener tus pedidos", Details: models.DetalleError(err)})
	}
	return c.JSON(pagina)
}

// GetPedidosAdmin lista todos los pedidos con filtro opcional ?estado= (paginado).
func (pc *PedidoController) GetPedidosAdmin(c *fiber.Ctx) error {
	page, limit := paginacion(c)
	estado := c.Query("estado")

	pagina, err := pc.repo.ObtenerPedidosAdmin(estado, page, limit)
	if err != nil {
		return c.Status(500).JSON(models.APIError{Status: 500, Message: "Error al obtener los pedidos", Details: models.DetalleError(err)})
	}
	return c.JSON(pagina)
}

// GetPedido devuelve el detalle de un pedido. Un cliente solo puede ver los
// suyos; el admin puede ver cualquiera (control de acceso a nivel de objeto).
func (pc *PedidoController) GetPedido(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(models.APIError{Status: 400, Message: "El id debe ser un número válido"})
	}

	pedido, err := pc.repo.ObtenerPedidoPorID(id)
	if err != nil {
		if errors.Is(err, repository.ErrNoEncontrado) {
			return c.Status(404).JSON(models.APIError{Status: 404, Message: "Pedido no encontrado"})
		}
		return c.Status(500).JSON(models.APIError{Status: 500, Message: "Error al obtener el pedido", Details: models.DetalleError(err)})
	}

	// Autorización: el dueño o un admin.
	if rolDeCtx(c) != "admin" {
		usuarioID, err := usuarioIDDeCtx(c)
		if err != nil || pedido.UsuarioID == nil || *pedido.UsuarioID != usuarioID {
			return c.Status(403).JSON(models.APIError{Status: 403, Message: "No tienes acceso a este pedido"})
		}
	}
	return c.JSON(pedido)
}

// ActualizarEstado (admin) aplica una transición válida del ciclo de vida.
func (pc *PedidoController) ActualizarEstado(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(models.APIError{Status: 400, Message: "El id debe ser un número válido"})
	}

	var body struct {
		Estado string `json:"estado"`
	}
	if err := c.BodyParser(&body); err != nil || body.Estado == "" {
		return c.Status(400).JSON(models.APIError{Status: 400, Message: "El campo 'estado' es obligatorio"})
	}

	pedido, err := pc.repo.CambiarEstadoPedido(id, body.Estado)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrNoEncontrado):
			return c.Status(404).JSON(models.APIError{Status: 404, Message: "Pedido no encontrado"})
		case errors.Is(err, repository.ErrTransicionInvalida):
			return c.Status(400).JSON(models.APIError{Status: 400, Message: "Transición de estado no permitida"})
		default:
			return c.Status(500).JSON(models.APIError{Status: 500, Message: "Error al actualizar el pedido", Details: models.DetalleError(err)})
		}
	}
	return c.JSON(pedido)
}
