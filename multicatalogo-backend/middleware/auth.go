// Declaramos el paquete middleware, que agrupa los middlewares HTTP de Fiber.
package middleware

import (
	// Os lee las variables de entorno (JWT_SECRET).
	"os"
	// Strconv convierte el id del usuario al claim "sub" (string).
	"strconv"
	// Strings manipula la cabecera Authorization (prefijo Bearer).
	"strings"
	// Time define la expiración del token (iat/exp).
	"time"

	// Fiber provee el tipo del middleware (c *fiber.Ctx).
	"github.com/gofiber/fiber/v2"
	// JWT es la librería oficial golang-jwt para firmar y validar tokens HS256.
	"github.com/golang-jwt/jwt/v5"

	// Importamos los modelos para tipar los claims del token.
	"multicatalogo-backend/models"
)

// Secret devuelve el secreto de firma JWT leído de JWT_SECRET.
// Si no está definido (entornos locales), usa un secreto de desarrollo con una
// advertencia en el log: en producción JWT_SECRET es obligatorio.
func Secret() string {
	if s := os.Getenv("JWT_SECRET"); s != "" {
		return s
	}
	return "dev-secret-cambiar-en-produccion"
}

// Claims es la carga útil del JWT: sub (id), email, rol + tiempos estándar.
type Claims struct {
	Email string `json:"email"`
	Rol   string `json:"rol"`
	jwt.RegisteredClaims
}

// GenerarToken firma un JWT HS256 válido por 2 horas para el usuario dado.
func GenerarToken(u models.Usuario) (string, error) {
	now := time.Now()
	claims := Claims{
		Email: u.Email,
		Rol:   u.Rol,
		RegisteredClaims: jwt.RegisteredClaims{
			// Sub lleva el id del usuario como string (convención JWT).
			Subject:   strconv.Itoa(u.ID),
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(2 * time.Hour)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(Secret()))
}

// RequireAuth valida el JWT de la cabecera Authorization y expone la identidad
// en c.Locals("email"), c.Locals("rol") y c.Locals("usuarioID").
// Responde 401 con APIError si el token falta, está malformado o expiró.
func RequireAuth(secret string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Extraemos el token del header "Authorization: Bearer <token>".
		auth := c.Get("Authorization")
		if auth == "" {
			return c.Status(401).JSON(models.APIError{Status: 401, Message: "Token de autenticación requerido"})
		}
		tokenStr := strings.TrimPrefix(auth, "Bearer ")
		if tokenStr == auth {
			return c.Status(401).JSON(models.APIError{Status: 401, Message: "Formato de token inválido (se espera Bearer)"})
		}

		// Parseamos y validamos la firma y la expiración del token.
		claims := &Claims{}
		token, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) {
			return []byte(secret), nil
		})
		if err != nil || !token.Valid {
			return c.Status(401).JSON(models.APIError{Status: 401, Message: "Token inválido o expirado"})
		}

		// Publicamos la identidad en Locals para que los controladores la usen
		// (principio: la identidad SIEMPRE viene del token, nunca del body).
		c.Locals("email", claims.Email)
		c.Locals("rol", claims.Rol)
		c.Locals("usuarioID", claims.Subject)
		return c.Next()
	}
}

// RequireAdmin exige rol admin. Debe encadenarse DESPUÉS de RequireAuth.
// Responde 403 con APIError si el rol no es admin.
func RequireAdmin() fiber.Handler {
	return func(c *fiber.Ctx) error {
		if c.Locals("rol") != "admin" {
			return c.Status(403).JSON(models.APIError{Status: 403, Message: "Se requiere rol de administrador"})
		}
		return c.Next()
	}
}
