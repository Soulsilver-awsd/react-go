// Declaramos el paquete controllers para agrupar las funciones que manejan la lógica de negocio de las rutas.
package controllers

import (
	// Importamos Fiber para tener acceso al contexto (c *fiber.Ctx) de la petición HTTP y al status HTTP.
	"github.com/gofiber/fiber/v2"

	// Importamos nuestro paquete de modelos para usar la estructura LoginRequest y construir APIError.
	"multicatalogo-backend/models"

	// Importamos la capa de acceso a datos; aquí solo se le piden usuarios, nunca se tocan datos directamente.
	"multicatalogo-backend/repository"
)

// Login es la función controladora que se ejecutará cuando el cliente envíe sus credenciales.
func Login(c *fiber.Ctx) error {
	// Creamos una variable 'req' del tipo LoginRequest (ubicada en nuestro paquete models) para almacenar los datos.
	var req models.LoginRequest

	// Intentamos parsear (transformar) el cuerpo JSON entrante y guardarlo en la variable 'req'.
	if err := c.BodyParser(&req); err != nil {
		// Si el JSON está mal formado construimos el error estandarizado y lo devolvemos con estado 400.
		// Usamos la forma congelada porque este endpoint ya lo consume el frontend.
		return ResponderErrorLegado(c, models.NewAPIError(400, "Cuerpo de petición inválido", nil))
	}

	// Validamos la entrada antes de tocar la capa de datos: un email o password vacío no es una
	// credencial, es una petición incompleta, y no debe llegar al repositorio.
	if req.Email == "" || req.Password == "" {
		// Indicamos en 'details' qué campo(s) faltaron, para que el cliente pueda corregir el formulario.
		return ResponderErrorLegado(c, models.NewAPIError(400, "Email y password son obligatorios", camposFaltantes(req)))
	}

	// Resolvemos el tenant de la petición: lo asigna el servidor a partir de la cabecera, nunca el body.
	tenantID := TenantDesdeRequest(c)

	// Consultamos el usuario por email dentro de ese tenant. A este punto los datos ya están validados.
	usuario, err := repository.BuscarPorEmail(tenantID, req.Email)
	if err != nil || usuario.Password != req.Password {
		// Si el usuario no existe en el tenant o la contraseña no coincide, la respuesta es la misma
		// en ambos casos: no le revelamos a un atacante qué parte de la credencial fue incorrecta.
		return ResponderErrorLegado(c, models.NewAPIError(401, "Credenciales incorrectas", nil))
	}

	// Las credenciales son correctas: devolvemos 200 con el token, el email y el rol que espera el frontend.
	// El token sigue siendo ficticio, igual que en la línea base, para no alterar el contrato.
	return c.JSON(fiber.Map{"token": tokenParaRol(usuario.Rol), "email": usuario.Email, "rol": usuario.Rol})
}

// camposFaltantes arma el detalle opcional del error de validación: la lista de campos vacíos
// que provocaron el 400. Devuelve nil si, contraryamente, ambos campos vienen informados.
func camposFaltantes(req models.LoginRequest) interface{} {
	// Declaramos el slice que irá en el campo 'details' del APIError.
	faltantes := []string{}
	// Si el email llegó vacío, lo anotamos.
	if req.Email == "" {
		// Agregamos el nombre del campo en el formato que espera el JSON.
		faltantes = append(faltantes, "email")
	}
	// Si el password llegó vacío, lo anotamos también.
	if req.Password == "" {
		// Agregamos el nombre del campo en el mismo formato.
		faltantes = append(faltantes, "password")
	}
	// Si no hubo ningún campo faltante, devolvemos nil para que 'details' se omita del JSON.
	if len(faltantes) == 0 {
		// Un nil hace que la clave 'details' desaparezca gracias al tag omitempty del APIError.
		return nil
	}
	// Devolvemos la lista de campos que el cliente debe completar.
	return faltantes
}

// tokenParaRol devuelve el token ficticio asociado al rol del usuario. Se mantiene aquí, junto al
// handler de login, porque es el único punto donde el token se emite.
func tokenParaRol(rol string) string {
	// Si el usuario es administrador, entregamos el token de la línea base.
	if rol == "admin" {
		// Token ficticio del administrador, idéntico al que el frontend ya conoce.
		return "fake-jwt-token-123"
	}
	// Cualquier otro rol (por ahora "cliente") recibe su propio token ficticio.
	return "fake-jwt-token-456"
}

// El código anterior validaba las credenciales con un switch hardcodeado directamente en el
// controller. Se conserva aquí como registro histórico; su lógica vive ahora en el repositorio
// repository/usuarioRepo.go y la validación de entrada, en este mismo handler.
//
//	switch {
//	case req.Email == "admin@upse.edu.ec" && req.Password == "123456":
//		return c.JSON(fiber.Map{"token": "fake-jwt-token-123", "email": req.Email, "rol": "admin"})
//	case req.Email == "cliente@upse.edu.ec" && req.Password == "123456":
//		return c.JSON(fiber.Map{"token": "fake-jwt-token-456", "email": req.Email, "rol": "cliente"})
//	default:
//		return c.Status(401).JSON(fiber.Map{"error": "Credenciales incorrectas"})
//	}
