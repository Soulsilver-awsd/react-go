// Declaramos el paquete config, que además de la base de datos administra la conexión a Redis.
package config

import (
	// Context define el alcance temporal de cada operación contra Redis.
	"context"
	// Log imprime mensajes informativos durante la inicialización.
	"log"
	// Os permite leer las variables de entorno del sistema.
	"os"

	// Redis es la librería oficial del cliente de Redis para Go.
	"github.com/redis/go-redis/v9"
)

// Redis es el cliente global de Redis reutilizado al construir la capa de caché.
var Redis *redis.Client

// ConnectRedis crea e inicializa el cliente de Redis usando variables de entorno.
func ConnectRedis() *redis.Client {
	// Leemos la dirección del servidor Redis; por defecto apuntamos al contenedor local.
	addr := os.Getenv("REDIS_ADDR")
	if addr == "" {
		addr = "localhost:6379"
	}

	// Creamos el cliente con la dirección, contraseña (opcional) y base de datos lógica 0.
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       0,
	})

	// Verificamos la conectividad haciendo un ping.
	if err := client.Ping(context.Background()).Err(); err != nil {
		// Si Redis no está disponible, degradamos de forma controlada: la aplicación
		// sigue funcionando y simplemente consultará siempre PostgreSQL (sin caché).
		log.Printf("Aviso: no se pudo conectar a Redis (%v); se operará sin caché", err)
		Redis = nil
		return client
	}

	// Guardamos el cliente en la variable global y devolvemos éxito.
	Redis = client
	log.Println("Conexión a Redis establecida correctamente")
	return client
}

// CloseRedis cierra el cliente de Redis liberando los recursos al finalizar la aplicación.
func CloseRedis() {
	if Redis != nil {
		Redis.Close()
		log.Println("Conexión a Redis cerrada")
	}
}
