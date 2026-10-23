// Declaramos el paquete controllers para agrupar las funciones que manejan la lógica de negocio de las rutas.
package controllers

import (
	// Importamos la librería estándar 'errors' para comparar errores centinela del repositorio.
	"errors"

	// Importamos el framework Fiber para tener acceso al contexto (c *fiber.Ctx) de la petición HTTP.
	"github.com/gofiber/fiber/v2"
	// Bcrypt implementa la verificación de contraseñas con hash bcrypt.
	"golang.org/x/crypto/bcrypt"
	// Importamos el middleware para firmar el JWT de sesión.
	"multicatalogo-backend/middleware"
	// Importamos el paquete de modelos para poder usar la estructura LoginRequest y APIError.
	"multicatalogo-backend/models"
	// Importamos el repositorio, que ahora obtiene los datos reales desde PostgreSQL.
	"multicatalogo-backend/repository"
)

// AuthController agrupa las dependencias necesarias para la autenticación (inyección de dependencias).
type AuthController struct {
	// repo es la abstracción de acceso a datos inyectada.
	repo repository.DBRepository
}

// NewAuthController construye el controlador inyectando sus dependencias.
func NewAuthController(repo repository.DBRepository) *AuthController {
	return &AuthController{repo: repo}
}

// Login es la función controladora que se ejecutará cuando el cliente envíe sus credenciales.
func (ac *AuthController) Login(c *fiber.Ctx) error {
	// Creamos una variable 'req' del tipo LoginRequest (ubicada en nuestro paquete models) para almacenar los datos.
	var req models.LoginRequest

	// Intentamos parsear (transformar) el cuerpo JSON entrante y guardarlo en la variable 'req'.
	if err := c.BodyParser(&req); err != nil {
		// Si ocurre un error al parsear (ej. JSON mal formado), retornamos un estado HTTP 400 (Bad Request) con la estructura estandarizada.
		return c.Status(400).JSON(models.APIError{
			Status:  400,
			Message: "Cuerpo de petición inválido",
		})
	}

	// Validamos explícitamente que el campo Email no venga vacío antes de continuar.
	if req.Email == "" {
		return c.Status(400).JSON(models.APIError{
			Status:  400,
			Message: "El campo email es obligatorio",
		})
	}

	// Validamos explícitamente que el campo Password no venga vacío antes de continuar.
	if req.Password == "" {
		return c.Status(400).JSON(models.APIError{
			Status:  400,
			Message: "El campo password es obligatorio",
		})
	}

	// Consultamos el usuario por email (el password se compara con bcrypt abajo).
	usuario, err := ac.repo.BuscarUsuarioPorEmail(req.Email)
	if err != nil {
		// Si el usuario no existe, las credenciales son incorrectas (HTTP 401).
		if errors.Is(err, repository.ErrNoEncontrado) {
			return c.Status(401).JSON(models.APIError{
				Status:  401,
				Message: "Credenciales incorrectas",
			})
		}
		// Cualquier otro error se reporta como HTTP 500.
		return c.Status(500).JSON(models.APIError{
			Status:  500,
			Message: "Error al validar las credenciales",
			Details: models.DetalleError(err),
		})
	}

	// Comparamos la contraseña enviada contra el hash bcrypt almacenado.
	// CompareHashAndPassword es la única comparación válida: nunca se guarda
	// ni se envía la contraseña en texto plano por ningún canal.
	if err := bcrypt.CompareHashAndPassword([]byte(usuario.Password), []byte(req.Password)); err != nil {
		return c.Status(401).JSON(models.APIError{
			Status:  401,
			Message: "Credenciales incorrectas",
		})
	}

	// Rechazamos cuentas desactivadas por el administrador (soft-disable).
	if !usuario.Activo {
		return c.Status(403).JSON(models.APIError{
			Status:  403,
			Message: "Cuenta desactivada",
		})
	}

	// Firmamos un JWT real (HS256, 2h) con claims sub/email/rol/iat/exp.
	token, err := middleware.GenerarToken(*usuario)
	if err != nil {
		return c.Status(500).JSON(models.APIError{
			Status:  500,
			Message: "Error al generar el token de sesión",
			Details: models.DetalleError(err),
		})
	}

	// Devolvemos el token y los datos del usuario manteniendo la estructura esperada por el frontend.
	return c.JSON(fiber.Map{
		"token": token,
		"email": usuario.Email,
		"rol":   usuario.Rol,
	})
}
