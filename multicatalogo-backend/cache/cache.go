// Declaramos el paquete cache, que define la abstracción de almacenamiento en caché.
package cache

import (
	// Context define el alcance temporal de cada operación contra Redis.
	"context"
	// Time define las duraciones de expiración (TTL).
	"time"

	// Redis es la librería oficial del cliente de Redis para Go.
	"github.com/redis/go-redis/v9"
)

// CacheRepository es la interfaz que abstrae el sistema de caché (inyección de dependencias).
type CacheRepository interface {
	// Get devuelve el valor de una clave y true si existe (cache hit).
	Get(key string) (string, bool)
	// Set guarda un valor con un tiempo de expiración determinado.
	Set(key, value string, ttl time.Duration)
	// Del elimina una clave (invalidación de caché).
	Del(key string)
	// DelByPrefix elimina todas las claves que comienzan con un prefijo.
	DelByPrefix(prefix string)
}

// RedisCache es la implementación de CacheRepository basada en Redis.
type RedisCache struct {
	// client es el cliente de Redis inyectado.
	client *redis.Client
}

// NewRedisCache construye una instancia de RedisCache con el cliente inyectado.
func NewRedisCache(client *redis.Client) *RedisCache {
	return &RedisCache{client: client}
}

// Get intenta leer el valor de una clave; devuelve el valor y true si existe (cache hit).
func (c *RedisCache) Get(key string) (string, bool) {
	// Si Redis no está disponible, se comporta como cache miss.
	if c.client == nil {
		return "", false
	}
	// Leemos el valor de la clave.
	val, err := c.client.Get(context.Background(), key).Result()
	if err != nil {
		// Cualquier error (incluido redis.Nil) se trata como cache miss.
		return "", false
	}
	return val, true
}

// Set guarda un valor en Redis con un tiempo de expiración (TTL) determinado.
func (c *RedisCache) Set(key, value string, ttl time.Duration) {
	// Si Redis no está disponible, no hacemos nada (degradación controlada).
	if c.client == nil {
		return
	}
	// Guardamos el valor ignorando el error: un fallo de caché no debe romper la petición.
	c.client.Set(context.Background(), key, value, ttl)
}

// Del elimina una clave de Redis (usado para invalidar la caché tras un cambio de datos).
func (c *RedisCache) Del(key string) {
	// Si Redis no está disponible, no hacemos nada.
	if c.client == nil {
		return
	}
	// Eliminamos la clave sin propagar errores.
	c.client.Del(context.Background(), key)
}

// DelByPrefix elimina todas las claves que comienzan con un prefijo.
func (c *RedisCache) DelByPrefix(prefix string) {
	// Si Redis no está disponible, no hacemos nada.
	if c.client == nil {
		return
	}
	// Buscamos todas las claves que comienzan con el prefijo.
	keys, err := c.client.Keys(context.Background(), prefix+"*").Result()
	if err != nil || len(keys) == 0 {
		return
	}
	// Eliminamos todas las claves encontradas sin propagar errores.
	c.client.Del(context.Background(), keys...)
}
