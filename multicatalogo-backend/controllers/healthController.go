// Declaramos el paquete controllers, ya que este archivo expone el estado del servidor.
package controllers

import (
	// Context define el alcance de los pings a dependencias.
	"context"
	// Time acota el tiempo de espera de cada ping.
	"time"

	// Fiber provee el tipo del controlador.
	"github.com/gofiber/fiber/v2"
	// Pgxpool es el pool de PostgreSQL inyectado para el ping.
	"github.com/jackc/pgx/v5/pgxpool"
	// Redis es el cliente de Redis inyectado para el ping.
	"github.com/redis/go-redis/v9"

	// Importamos los modelos para estandarizar la respuesta de error.
	"multicatalogo-backend/models"
)

// HealthController agrupa las dependencias necesarias para verificar el estado
// del servidor y de sus dependencias críticas (Postgres y Redis).
type HealthController struct {
	// db es el pool de PostgreSQL con el que se hace el ping.
	db *pgxpool.Pool
	// rdb es el cliente Redis con el que se hace el ping.
	rdb *redis.Client
}

// NewHealthController construye el controlador de salud inyectando sus dependencias.
func NewHealthController(db *pgxpool.Pool, rdb *redis.Client) *HealthController {
	return &HealthController{db: db, rdb: rdb}
}

// Health devuelve el estado real del servicio:
//   - ok        → PostgreSQL y Redis responden
//   - degraded  → PostgreSQL responde pero Redis no (la API sigue funcionando sin caché)
//   - 503       → PostgreSQL no responde (sin persistencia no hay servicio)
func (hc *HealthController) Health(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	// Ping a PostgreSQL: sin él no hay servicio, respondemos 503.
	if hc.db == nil {
		return c.Status(503).JSON(models.APIError{Status: 503, Message: "PostgreSQL no configurado"})
	}
	if err := hc.db.Ping(ctx); err != nil {
		return c.Status(503).JSON(models.APIError{
			Status:  503,
			Message: "PostgreSQL no responde",
			Details: models.DetalleError(err),
		})
	}

	// Ping a Redis: si falla, el servicio se reporta "degraded" pero sigue vivo.
	estado := "ok"
	if hc.rdb == nil {
		estado = "degraded"
	} else if err := hc.rdb.Ping(ctx).Err(); err != nil {
		estado = "degraded"
	}

	return c.Status(200).JSON(fiber.Map{
		"status":  estado,
		"message": "Servidor Go/Fiber operativo",
		"deps": fiber.Map{
			"postgres": "ok",
			"redis":    map[bool]string{true: "ok", false: "caído"}[estado == "ok"],
		},
	})
}
