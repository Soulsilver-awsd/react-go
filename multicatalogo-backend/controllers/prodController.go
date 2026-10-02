// Seguimos dentro del paquete controllers, ya que maneja la lógica de otra sección del negocio.
package controllers

import (
	// Importamos el paquete errors para comparar el error centinela del repositorio con errors.Is.
	"errors"

	// Importamos Fiber para manejar la respuesta HTTP y leer los parámetros de la ruta.
	"github.com/gofiber/fiber/v2"

	// Importamos strconv para convertir el parámetro textual "id" de la ruta al entero esperado.
	"strconv"

	// Importamos el paquete models para usar la estructura Producto y construir el APIError.
	"multicatalogo-backend/models"

	// Importamos la capa de acceso a datos, filtrada siempre por tenant.
	"multicatalogo-backend/repository"
)

// GetProductos es la función controladora encargada de devolver el catálogo de artículos.
// Mantiene el contrato congelado: un array plano de productos con id, nombre, precio e img.
func GetProductos(c *fiber.Ctx) error {
	// Resolvemos el tenant de la petición; el frontend actual no manda cabecera, cae en "default".
	tenantID := TenantDesdeRequest(c)

	// Pedimos el catálogo al repositorio, que devuelve únicamente los productos de ese tenant.
	productos := repository.Listar(tenantID)

	// Fiber convierte el slice de estructuras a JSON. El slice nunca es nil (el repositorio lo
	// inicializa vacío), así que la respuesta es siempre un array, tal como el frontend espera.
	return c.JSON(productos)
}

// GetProductoPorID devuelve un único producto a partir del identificador de la ruta (/api/productos/:id).
// Antes de consultar datos valida que el id sea numérico: un id inválido no llega al repositorio.
func GetProductoPorID(c *fiber.Ctx) error {
	// Leemos el parámetro "id" de la ruta; Fiber lo entrega como texto.
	idTexto := c.Params("id")

	// Convertimos el texto al entero que usa el dominio. La ruta puede traer cualquier cadena,
	// así que el error de conversión se maneja explícitamente en lugar de ignorarse.
	id, err := strconv.Atoi(idTexto)
	if err != nil {
		// El id no es numérico: respondemos 400 con la envoltura estándar de la API.
		// El detalle deja claro qué parámetro fue rechazado, y el return impide que el id inválido
		// llegue a la capa de datos.
		return ResponderError(c, models.NewAPIError(400, "El id debe ser un número entero",
			fiber.Map{"parametro": "id", "valor": idTexto}))
	}

	// El id es válido: ahora sí resolvemos el tenant y consultamos el repositorio.
	tenantID := TenantDesdeRequest(c)

	// Buscamos el producto dentro del tenant; la consulta siempre lleva el filtro por tenant_id.
	producto, err := repository.BuscarPorID(tenantID, id)
	if err != nil {
		// Si el error es el centinela "no encontrado" (o el producto es de otro tenant, que se
		// traduce igual), respondemos 404 y no 403 para no revelar la existencia de datos ajenos.
		if errors.Is(err, repository.ErrNoEncontrado) {
			// Construimos el error con la envoltura estándar.
			return ResponderError(c, models.NewAPIError(404, "Producto no encontrado",
				fiber.Map{"id": id}))
		}
		// Cualquier otro fallo es interno: 500 con mensaje genérico, sin filtrar detalles técnicos.
		return ResponderError(c, models.NewAPIError(500, "Error interno del servidor", nil))
	}

	// Producto encontrado: devolvemos el objeto con la misma forma que consume el frontend
	// en DetalleProducto.tsx (id, nombre, precio, img).
	return c.JSON(producto)
}
