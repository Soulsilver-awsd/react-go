// Declaramos el paquete controllers: gestión de usuarios (Fase 7).
package controllers

import (
	// Errors compara los errores centinela del repositorio.
	"errors"
	// Strconv convierte el id de la URL.
	"strconv"

	// Fiber maneja las peticiones y respuestas HTTP.
	"github.com/gofiber/fiber/v2"
	// Modelos de respuesta estandarizados.
	"multicatalogo-backend/models"
	// Repositorio con las consultas de usuarios.
	"multicatalogo-backend/repository"
)

// UsuarioController agrupa las dependencias del back-office de usuarios.
type UsuarioController struct {
	repo repository.DBRepository
}

// NewUsuarioController construye el controlador con su repositorio inyectado.
func NewUsuarioController(repo repository.DBRepository) *UsuarioController {
	return &UsuarioController{repo: repo}
}

// GetUsuarios (admin) lista las cuentas de forma paginada.
func (uc *UsuarioController) GetUsuarios(c *fiber.Ctx) error {
	page, limit := paginacion(c)
	pagina, err := uc.repo.ObtenerUsuarios(page, limit)
	if err != nil {
		return c.Status(500).JSON(models.APIError{Status: 500, Message: "Error al obtener los usuarios", Details: models.DetalleError(err)})
	}
	return c.JSON(pagina)
}

// ActualizarUsuario (admin) cambia el rol y/o activa/desactiva una cuenta.
func (uc *UsuarioController) ActualizarUsuario(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(models.APIError{Status: 400, Message: "El id debe ser un número válido"})
	}

	// Punteros para distinguir campos omitidos (nil) de valores concretos.
	var body struct {
		Rol    *string `json:"rol"`
		Activo *bool   `json:"activo"`
	}
	if err := c.BodyParser(&body); err != nil {
		return c.Status(400).JSON(models.APIError{Status: 400, Message: "Cuerpo de petición inválido"})
	}
	if body.Rol == nil && body.Activo == nil {
		return c.Status(400).JSON(models.APIError{Status: 400, Message: "Debes enviar 'rol' y/o 'activo'"})
	}
	// Regla de negocio: solo existen los roles admin y cliente.
	if body.Rol != nil && *body.Rol != "admin" && *body.Rol != "cliente" {
		return c.Status(400).JSON(models.APIError{Status: 400, Message: "El rol debe ser 'admin' o 'cliente'"})
	}

	usuario, err := uc.repo.ActualizarUsuario(id, body.Rol, body.Activo)
	if err != nil {
		if errors.Is(err, repository.ErrNoEncontrado) {
			return c.Status(404).JSON(models.APIError{Status: 404, Message: "Usuario no encontrado"})
		}
		return c.Status(500).JSON(models.APIError{Status: 500, Message: "Error al actualizar el usuario", Details: models.DetalleError(err)})
	}
	return c.JSON(usuario)
}
