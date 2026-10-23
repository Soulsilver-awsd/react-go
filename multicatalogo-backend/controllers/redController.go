// Declaramos el paquete controllers, ya que este archivo maneja la lógica de la red multinivel.
package controllers

import (
	// Importamos la librería estándar 'encoding/json' para serializar/deserializar la caché.
	"encoding/json"
	// Importamos la librería estándar 'time' para definir el TTL de la caché.
	"time"

	// Importamos Fiber para manejar la respuesta HTTP.
	"github.com/gofiber/fiber/v2"
	// Importamos la capa de caché (abstracción inyectada).
	"multicatalogo-backend/cache"
	// Importamos los modelos para usar la estructura Referido y APIError.
	"multicatalogo-backend/models"
	// Importamos el repositorio, que reconstruye el árbol desde PostgreSQL.
	"multicatalogo-backend/repository"
)

// RedController agrupa las dependencias de la red multinivel (inyección de dependencias).
type RedController struct {
	// repo es la abstracción de acceso a datos inyectada.
	repo repository.DBRepository
	// cache es la abstracción de caché inyectada.
	cache cache.CacheRepository
}

// NewRedController construye el controlador inyectando sus dependencias.
func NewRedController(repo repository.DBRepository, cacheRepo cache.CacheRepository) *RedController {
	return &RedController{repo: repo, cache: cacheRepo}
}

// ClaveRed construye la clave de caché de la red para un email dado.
func ClaveRed(email string) string {
	return "red:tree:" + email
}

// GetRed devuelve el árbol multinivel de referidos del usuario autenticado.
func (rc *RedController) GetRed(c *fiber.Ctx) error {
	// Identificamos al usuario mediante el parámetro opcional ?email= para construir su clave de caché.
	// Si no se recibe, usamos "invitado" manteniendo intacta la firma original del endpoint.
	email := c.Query("email")
	if email == "" {
		email = "invitado"
	}
	// Clave de caché específica por usuario, según lo definido en la Fase 3.
	claveRed := ClaveRed(email)

	// ----- Patrón Cache-Aside: intentamos leer el árbol desde la caché -----
	if cached, ok := rc.cache.Get(claveRed); ok {
		// Si existe (Cache Hit), deserializamos el JSON y respondemos sin tocar PostgreSQL.
		var red models.Referido
		if err := json.Unmarshal([]byte(cached), &red); err == nil {
			return c.JSON(red)
		}
		// Si la caché estuviera corrupta, continuamos hacia PostgreSQL como si fuera un miss.
	}

	// ----- Cache Miss: reconstruimos el árbol consultando la tabla referidos -----
	red, err := rc.repo.ObtenerRed()
	if err != nil {
		// Si ocurre un error en la consulta, devolvemos un HTTP 500 estandarizado.
		return c.Status(500).JSON(models.APIError{
			Status:  500,
			Message: "Error al obtener la red de referidos",
			Details: models.DetalleError(err),
		})
	}

	// Serializamos el árbol y lo guardamos en la caché con un TTL de 10 minutos.
	if data, err := json.Marshal(red); err == nil {
		rc.cache.Set(ClaveRed(email), string(data), 10*time.Minute)
	}

	// Devolvemos el árbol en formato JSON con un HTTP 200.
	return c.JSON(red)
}
