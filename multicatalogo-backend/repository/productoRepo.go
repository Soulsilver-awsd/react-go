// Seguimos dentro del paquete repository, ahora para el catálogo de productos.
package repository

import (
	// Importamos el paquete models, de donde salen las estructuras de producto del dominio.
	"multicatalogo-backend/models"
)

// productoRepo guarda el catálogo de productos. En esta iteración los datos viven en memoria
// (seed), replicando los cuatro productos históricos de la línea base para no alterar el
// contrato que ya consume el frontend.
type productoRepo struct {
	// productos es el almacén en memoria; cada registro pertenece a un tenant concreto.
	productos []models.Producto
}

// repositorioProducto es la instancia del repositorio de productos, inicializada con el seed
// de arranque. Los métodos solo leen, así que se expone como valor y no necesita puntero.
var repositorioProducto = productoRepo{productos: []models.Producto{
	// Producto 1: serum facial.
	{ID: 1, Nombre: "Serum Revitalizante", Precio: 45.00, Img: "https://picsum.photos/seed/serum/150", TenantID: "default"},
	// Producto 2: crema hidratante.
	{ID: 2, Nombre: "Crema Hidratante Pro", Precio: 32.50, Img: "https://picsum.photos/seed/crema/150", TenantID: "default"},
	// Producto 3: tónico purificante.
	{ID: 3, Nombre: "Tónico Purificante", Precio: 28.00, Img: "https://picsum.photos/seed/tonico/150", TenantID: "default"},
	// Producto 4: mascarilla nocturna.
	{ID: 4, Nombre: "Mascarilla Nocturna", Precio: 50.00, Img: "https://picsum.photos/seed/mascarilla/150", TenantID: "default"},
}}

// Listar devuelve todos los productos del tenant indicado, en el orden en que fueron sembrados.
// El filtro por tenant es obligatorio: ningún inquilino puede ver el catálogo de otro.
func Listar(tenantID string) []models.Producto {
	// Declaramos el slice de salida, inicialmente vacío, para acumular solo los productos del tenant.
	productos := []models.Producto{}
	// Recorremos el catálogo completo comparando el tenant de cada registro con el solicitado.
	for _, producto := range repositorioProducto.productos {
		// Si el producto pertenece al tenant pedido, lo agregamos a la respuesta.
		if producto.TenantID == tenantID {
			// Append sobre un slice no vacío garantiza que el JSON sea un array (nunca null).
			productos = append(productos, producto)
		}
	}
	// Devolvemos el catálogo filtrado, que puede ser un array vacío si el tenant no tiene productos.
	return productos
}

// BuscarPorID devuelve el producto con ese identificador dentro del tenant indicado.
// Si no existe o pertenece a otro tenant devuelve ErrNoEncontrado, que el controller
// traduce a 404 para no revelar información de otros inquilinos.
func BuscarPorID(tenantID string, id int) (models.Producto, error) {
	// Recorremos el catálogo buscando coincidencia de tenant e identificador.
	for _, producto := range repositorioProducto.productos {
		// La comparación incluye el tenant: el mismo id puede existir en dos inquilinos distintos.
		if producto.TenantID == tenantID && producto.ID == id {
			// Devolvemos una copia del producto encontrado.
			return producto, nil
		}
	}
	// Si llegamos aquí, el producto no existe en ese tenant.
	return models.Producto{}, ErrNoEncontrado
}
