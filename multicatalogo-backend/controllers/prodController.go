// Seguimos dentro del paquete controllers, ya que maneja la lógica de otra sección del negocio.
package controllers

import (
	// Importamos la librería estándar 'encoding/json' para serializar/deserializar la caché.
	"encoding/json"
	// Importamos la librería estándar 'errors' para comparar errores centinela del repositorio.
	"errors"
	// Importamos la librería estándar 'strconv' para convertir cadenas de texto a números.
	"strconv"
	// Importamos la librería estándar 'time' para definir el TTL de la caché.
	"time"

	// Importamos Fiber para manejar la respuesta HTTP.
	"github.com/gofiber/fiber/v2"
	// Importamos la capa de caché (abstracción inyectada).
	"multicatalogo-backend/cache"
	// Importamos nuestro paquete de modelos para usar la estructura Producto y APIError.
	"multicatalogo-backend/models"
	// Importamos el repositorio, que obtiene los datos reales desde PostgreSQL.
	"multicatalogo-backend/repository"
)

// claveProductos es la clave de caché del catálogo completo.
const claveProductos = "productos:all"

// ProductoController agrupa las dependencias del catálogo (inyección de dependencias).
type ProductoController struct {
	// repo es la abstracción de acceso a datos inyectada.
	repo repository.DBRepository
	// cache es la abstracción de caché inyectada.
	cache cache.CacheRepository
}

// NewProductoController construye el controlador inyectando sus dependencias.
func NewProductoController(repo repository.DBRepository, cacheRepo cache.CacheRepository) *ProductoController {
	return &ProductoController{repo: repo, cache: cacheRepo}
}

// GetProductos es la función controladora encargada de devolver el catálogo de artículos.
func (pc *ProductoController) GetProductos(c *fiber.Ctx) error {
	// ----- Patrón Cache-Aside: intentamos leer el catálogo desde la caché -----
	if cached, ok := pc.cache.Get(claveProductos); ok {
		// Si existe (Cache Hit), deserializamos el JSON y respondemos sin tocar PostgreSQL.
		var productos []models.Producto
		if err := json.Unmarshal([]byte(cached), &productos); err == nil {
			return c.JSON(productos)
		}
		// Si la caché estuviera corrupta, continuamos hacia PostgreSQL como si fuera un miss.
	}

	// ----- Cache Miss: consultamos PostgreSQL a través del repositorio -----
	productos, err := pc.repo.ObtenerProductos()
	if err != nil {
		// Si la consulta falla (ej. base de datos caída), devolvemos un HTTP 500 estandarizado.
		return c.Status(500).JSON(models.APIError{
			Status:  500,
			Message: "Error al obtener los productos",
			Details: models.DetalleError(err),
		})
	}

	// Serializamos el resultado y lo guardamos en la caché con un TTL de 5 minutos.
	if data, err := json.Marshal(productos); err == nil {
		pc.cache.Set(claveProductos, string(data), 5*time.Minute)
	}

	// Devolvemos el catálogo consultado en formato JSON.
	return c.JSON(productos)
}

// GetProductoPorID devuelve un único producto a partir del ID recibido en la URL.
func (pc *ProductoController) GetProductoPorID(c *fiber.Ctx) error {
	// Convertimos el parámetro de ruta "id" (string) a un número entero (int).
	id, err := strconv.Atoi(c.Params("id"))
	// Si la conversión falla (ej. el usuario envía texto en lugar de un número), devolvemos un HTTP 400 con la estructura estandarizada.
	if err != nil {
		return c.Status(400).JSON(models.APIError{
			Status:  400,
			Message: "El id debe ser un número válido",
			Details: models.DetalleError(err),
		})
	}

	// Consultamos el producto por su ID en PostgreSQL.
	producto, err := pc.repo.ObtenerProductoPorID(id)
	if err != nil {
		// Si el producto no existe, devolvemos un HTTP 404.
		if errors.Is(err, repository.ErrNoEncontrado) {
			return c.Status(404).JSON(models.APIError{
				Status:  404,
				Message: "Producto no encontrado",
			})
		}
		// Cualquier otro error se reporta como HTTP 500.
		return c.Status(500).JSON(models.APIError{
			Status:  500,
			Message: "Error al obtener el producto",
			Details: models.DetalleError(err),
		})
	}

	// Devolvemos el producto encontrado en formato JSON con un HTTP 200.
	return c.JSON(producto)
}

// GetProductosAdmin lista todos los productos (activos e inactivos) para el
// back-office de forma paginada. Acepta búsqueda server-side ?q= (ILIKE).
func (pc *ProductoController) GetProductosAdmin(c *fiber.Ctx) error {
	q := c.Query("q")
	page, limit := paginacion(c)
	productos, err := pc.repo.ObtenerProductosAdmin(q, page, limit)
	if err != nil {
		return c.Status(500).JSON(models.APIError{Status: 500, Message: "Error al obtener los productos", Details: models.DetalleError(err)})
	}
	return c.JSON(productos)
}

// validarProducto comprueba reglas de negocio mínimas del catálogo.
func validarProducto(p models.Producto) string {
	if p.Nombre == "" {
		return "El nombre es obligatorio"
	}
	if p.Precio <= 0 {
		return "El precio debe ser mayor que cero"
	}
	if p.Categoria == "" {
		return "La categoría es obligatoria"
	}
	if p.Stock < 0 {
		return "El stock no puede ser negativo"
	}
	return ""
}

// CrearProducto (admin) inserta un producto y invalida la caché del catálogo.
func (pc *ProductoController) CrearProducto(c *fiber.Ctx) error {
	var p models.Producto
	if err := c.BodyParser(&p); err != nil {
		return c.Status(400).JSON(models.APIError{Status: 400, Message: "Cuerpo de petición inválido"})
	}
	if msg := validarProducto(p); msg != "" {
		return c.Status(400).JSON(models.APIError{Status: 400, Message: msg})
	}

	creado, err := pc.repo.CrearProducto(p)
	if err != nil {
		return c.Status(500).JSON(models.APIError{Status: 500, Message: "Error al crear el producto", Details: models.DetalleError(err)})
	}
	// Invalidate: el catálogo público cacheado ya no refleja la realidad.
	pc.cache.Del(claveProductos)
	return c.Status(201).JSON(creado)
}

// ActualizarProducto (admin) edita un producto y invalida la caché del catálogo.
func (pc *ProductoController) ActualizarProducto(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(models.APIError{Status: 400, Message: "El id debe ser un número válido"})
	}

	var p models.Producto
	if err := c.BodyParser(&p); err != nil {
		return c.Status(400).JSON(models.APIError{Status: 400, Message: "Cuerpo de petición inválido"})
	}
	if msg := validarProducto(p); msg != "" {
		return c.Status(400).JSON(models.APIError{Status: 400, Message: msg})
	}

	actualizado, err := pc.repo.ActualizarProducto(id, p)
	if err != nil {
		if errors.Is(err, repository.ErrNoEncontrado) {
			return c.Status(404).JSON(models.APIError{Status: 404, Message: "Producto no encontrado"})
		}
		return c.Status(500).JSON(models.APIError{Status: 500, Message: "Error al actualizar el producto", Details: models.DetalleError(err)})
	}
	pc.cache.Del(claveProductos)
	return c.JSON(actualizado)
}

// DesactivarProducto (admin) aplica borrado lógico y limpia la caché.
func (pc *ProductoController) DesactivarProducto(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(models.APIError{Status: 400, Message: "El id debe ser un número válido"})
	}

	if err := pc.repo.DesactivarProducto(id); err != nil {
		if errors.Is(err, repository.ErrNoEncontrado) {
			return c.Status(404).JSON(models.APIError{Status: 404, Message: "Producto no encontrado"})
		}
		return c.Status(500).JSON(models.APIError{Status: 500, Message: "Error al desactivar el producto", Details: models.DetalleError(err)})
	}
	pc.cache.Del(claveProductos)
	return c.JSON(fiber.Map{"message": "Producto desactivado"})
}
