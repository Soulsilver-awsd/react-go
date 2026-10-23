package controllers

import (
	"errors"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

// usuarioIDDeCtx extrae el id del usuario autenticado desde c.Locals("usuarioID")
// (el claim "sub" del JWT, publicado como string por RequireAuth).
func usuarioIDDeCtx(c *fiber.Ctx) (int, error) {
	raw, _ := c.Locals("usuarioID").(string)
	if raw == "" {
		return 0, errors.New("usuarioID ausente en el contexto")
	}
	return strconv.Atoi(raw)
}

// emailDeCtx extrae el email del usuario autenticado desde Locals.
func emailDeCtx(c *fiber.Ctx) string {
	email, _ := c.Locals("email").(string)
	return email
}

// rolDeCtx extrae el rol del usuario autenticado desde Locals.
func rolDeCtx(c *fiber.Ctx) string {
	rol, _ := c.Locals("rol").(string)
	return rol
}
