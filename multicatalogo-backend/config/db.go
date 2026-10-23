// Declaramos el paquete config, encargado de la configuración de infraestructura del backend.
package config

import (
	// Context permite limitar el tiempo de vida de las operaciones contra la base de datos.
	"context"
	// Log imprime mensajes informativos y fatales durante la inicialización.
	"log"
	// Os permite leer las variables de entorno del sistema.
	"os"

	// Pgxpool provee un pool de conexiones concurrente para PostgreSQL.
	"github.com/jackc/pgx/v5/pgxpool"
)

// DB es el pool de conexiones global que reutiliza todas las capas de la aplicación.
var DB *pgxpool.Pool

// ConnectDB crea e inicializa el pool de conexiones a PostgreSQL usando la variable DATABASE_URL.
func ConnectDB() *pgxpool.Pool {
	// Leemos la cadena de conexión desde el entorno; si no existe, usamos una por defecto para desarrollo.
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "postgres://multicatalogo:multicatalogo@localhost:5432/multicatalogo?sslmode=disable"
	}

	// Creamos el pool de conexiones con la cadena de conexión.
	pool, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		// Si no se puede construir el pool, detenemos la aplicación.
		log.Fatalf("No se pudo crear el pool de conexiones: %v", err)
	}

	// Verificamos que la conexión realmente funcione haciendo un ping a la base de datos.
	if err := pool.Ping(context.Background()); err != nil {
		log.Fatalf("No se pudo conectar a PostgreSQL: %v", err)
	}

	// Guardamos el pool en la variable global para que el repositorio lo utilice.
	DB = pool
	log.Println("Conexión a PostgreSQL establecida correctamente")
	return pool
}

// CloseDB cierra el pool de conexiones liberando los recursos al finalizar la aplicación.
func CloseDB() {
	if DB != nil {
		DB.Close()
		log.Println("Conexión a PostgreSQL cerrada")
	}
}
