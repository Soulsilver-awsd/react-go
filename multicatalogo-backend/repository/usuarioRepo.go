// Declaramos el paquete repository, la capa de acceso a datos de la aplicación.
// Este paquete NO conoce Fiber: recibe y devuelve estructuras del paquete models y errores
// de Go, de modo que la capa de datos sea independiente del framework web.
package repository

import (
	// Importamos errors para declarar el error centinela que usan todos los repositorios.
	"errors"

	// Importamos el paquete models porque los repositorios devuelven las estructuras del dominio.
	"multicatalogo-backend/models"
)

// ErrNoEncontrado es el error centinela que los repositorios devuelven cuando un recurso
// no existe o pertenece a otro tenant. Los controllers lo traducen a un 404 en lugar de
// responder 403, para no filtrar la existencia de datos de otros inquilinos.
var ErrNoEncontrado = errors.New("recurso no encontrado")

// usuarioRepo guarda los usuarios del sistema. En esta iteración los datos viven en memoria
// (seed), pero la firma de los métodos ya es la que usaría una base de datos con filtro por tenant.
type usuarioRepo struct {
	// usuarios es el almacén en memoria; cada registro pertenece a un tenant concreto.
	usuarios []models.Usuario
}

// repositorioUsuario es la instancia del repositorio de usuarios, ya inicializada con el seed
// de arranque. Los métodos solo leen, así que se expone como valor y no necesita puntero.
var repositorioUsuario = usuarioRepo{usuarios: []models.Usuario{
	// Administrador del tenant por defecto: mismas credenciales que la línea base congelada.
	{Email: "admin@upse.edu.ec", Password: "123456", Rol: "admin", TenantID: "default"},
	// Cliente registrado del tenant por defecto.
	{Email: "cliente@upse.edu.ec", Password: "123456", Rol: "cliente", TenantID: "default"},
}}

// BuscarPorEmail devuelve el usuario que coincide con el email dentro del tenant indicado.
// Si no existe en ese tenant devuelve ErrNoEncontrado; nunca devuelve datos de otro inquilino.
func BuscarPorEmail(tenantID string, email string) (models.Usuario, error) {
	// Recorremos el almacén en memoria comparando tenant y email: el filtro por tenant es
	// parte de la condición de búsqueda, no un filtro aplicado después.
	for _, usuario := range repositorioUsuario.usuarios {
		// Solo contemplate el usuario si pertenece al tenant solicitado y el email coincide.
		if usuario.TenantID == tenantID && usuario.Email == email {
			// Devolvemos una copia del registro encontrado.
			return usuario, nil
		}
	}
	// Si llegamos aquí, el usuario no existe en ese tenant.
	return models.Usuario{}, ErrNoEncontrado
}
