// Declaramos el paquete controllers; este archivo maneja el registro de clientes.
package controllers

import (
	// Errors permite comparar errores centinela del repositorio.
	"errors"
	// Strings valida el formato básico del email.
	"strings"

	// Fiber maneja la respuesta HTTP.
	"github.com/gofiber/fiber/v2"
	// Bcrypt genera el hash de la contraseña antes de persistirla.
	"golang.org/x/crypto/bcrypt"

	// Importamos los modelos (RegisterRequest, APIError).
	"multicatalogo-backend/models"
	// Importamos el repositorio para persistir el usuario y su nodo de red.
	"multicatalogo-backend/repository"
)

// RegisterController agrupa las dependencias del registro (inyección de dependencias).
type RegisterController struct {
	repo repository.DBRepository
}

// NewRegisterController construye el controlador inyectando el repositorio.
func NewRegisterController(repo repository.DBRepository) *RegisterController {
	return &RegisterController{repo: repo}
}

// Register valida y crea una cuenta de cliente, enlazándola a la red si viene
// un código de referido válido. Responde 201 con los datos públicos del usuario.
func (rc *RegisterController) Register(c *fiber.Ctx) error {
	var req models.RegisterRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(models.APIError{Status: 400, Message: "Cuerpo de petición inválido"})
	}

	// Normalizamos entradas para evitar espacios accidentales.
	req.Nombre = strings.TrimSpace(req.Nombre)
	req.Email = strings.ToLower(strings.TrimSpace(req.Email))
	req.Ref = strings.ToUpper(strings.TrimSpace(req.Ref))

	// Validación de campos obligatorios (400).
	if req.Nombre == "" {
		return c.Status(400).JSON(models.APIError{Status: 400, Message: "El campo nombre es obligatorio"})
	}
	if req.Email == "" || !strings.Contains(req.Email, "@") {
		return c.Status(400).JSON(models.APIError{Status: 400, Message: "El campo email es obligatorio y debe ser válido"})
	}
	if len(req.Password) < 6 {
		return c.Status(400).JSON(models.APIError{Status: 400, Message: "La contraseña debe tener al menos 6 caracteres"})
	}

	// Hasheamos la contraseña con bcrypt; nunca se guarda en texto plano.
	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(500).JSON(models.APIError{Status: 500, Message: "Error al proteger la contraseña", Details: models.DetalleError(err)})
	}

	// Creamos usuario + nodo de red en una sola transacción.
	usuario, err := rc.repo.CrearUsuario(req.Nombre, req.Email, string(hash), req.Ref)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrEmailDuplicado):
			return c.Status(409).JSON(models.APIError{Status: 409, Message: "El email ya está registrado"})
		case errors.Is(err, repository.ErrReferidoInvalido):
			return c.Status(400).JSON(models.APIError{Status: 400, Message: "El código de referido no existe"})
		default:
			return c.Status(500).JSON(models.APIError{Status: 500, Message: "Error al registrar el usuario", Details: models.DetalleError(err)})
		}
	}

	// Respuesta sin token: el front redirige a /login para iniciar sesión.
	return c.Status(201).JSON(fiber.Map{
		"id":          usuario.ID,
		"email":       usuario.Email,
		"rol":         usuario.Rol,
		"nombre":      usuario.Nombre,
		"referido_id": usuario.ReferidoID,
		"message":     "Cuenta creada correctamente",
	})
}
