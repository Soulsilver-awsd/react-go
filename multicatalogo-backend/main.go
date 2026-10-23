// Declaramos que este es el paquete principal (main), es decir, el archivo ejecutable que inicia el programa.
package main

import (
	// Log imprime mensajes por consola durante el arranque y apagado.
	"log"
	// Os permite capturar señales del sistema para un apagado limpio.
	"os"
	// Signal permite registrar el manejo de señales del sistema operativo.
	"os/signal"
	// Syscall define las señales del sistema que vamos a escuchar.
	"syscall"

	// Importamos el framework principal Fiber.
	"github.com/gofiber/fiber/v2"
	// Importamos el middleware CORS para gestionar la seguridad entre distintos puertos/dominios.
	"github.com/gofiber/fiber/v2/middleware/cors"
	// Importamos godotenv para cargar las variables de entorno desde el archivo .env.
	"github.com/joho/godotenv"
	// Importamos la capa de caché (implementación Redis).
	"multicatalogo-backend/cache"
	// Importamos nuestro paquete de configuración, que administra las conexiones.
	"multicatalogo-backend/config"
	// Importamos los controladores que vamos a instanciar con inyección de dependencias.
	"multicatalogo-backend/controllers"
	// Importamos la capa de repositorio (acceso a PostgreSQL).
	"multicatalogo-backend/repository"
	// Importamos nuestro propio paquete de rutas para delegarle la configuración de los endpoints.
	"multicatalogo-backend/routes"
)

// func main es el punto de entrada de la aplicación en Go. Todo comienza a ejecutarse aquí.
func main() {
	// Cargamos las variables de entorno definidas en el archivo .env (si existe).
	if err := godotenv.Load(); err != nil {
		log.Println("Aviso: no se encontró el archivo .env, se usarán las variables del sistema")
	}

	// Inicializamos el pool de conexiones a PostgreSQL y garantizamos su cierre al finalizar.
	pool := config.ConnectDB()
	defer config.CloseDB()

	// Inicializamos el cliente de Redis (caché) y garantizamos su cierre al finalizar.
	redisClient := config.ConnectRedis()
	defer config.CloseRedis()

	// ===== Inyección de dependencias (Fase 4) =====
	// Construimos el repositorio concreto con el pool de PostgreSQL.
	repo := repository.NewPostgresRepository(pool)
	// Construimos la caché concreta con el cliente de Redis.
	cacheRepo := cache.NewRedisCache(redisClient)
	// Instanciamos los controladores inyectándoles sus dependencias.
	ctrls := &controllers.Controllers{
		Auth:        controllers.NewAuthController(repo),
		Register:    controllers.NewRegisterController(repo),
		Producto:    controllers.NewProductoController(repo, cacheRepo),
		Red:         controllers.NewRedController(repo, cacheRepo),
		Health:      controllers.NewHealthController(pool, redisClient),
		Checkout:    controllers.NewCheckoutController(repo, cacheRepo),
		Pedido:      controllers.NewPedidoController(repo),
		Comision:    controllers.NewComisionController(repo),
		Liquidacion: controllers.NewLiquidacionController(repo),
		Reporte:     controllers.NewReporteController(repo),
		Usuario:     controllers.NewUsuarioController(repo),
	}

	// Instanciamos una nueva aplicación de Fiber y guardamos la variable 'app'.
	app := fiber.New()

	// Orígenes permitidos por CORS. En local usa los puertos de Vite; en la nube
	// se configura CORS_ORIGINS con el dominio del frontend (lista separada por comas).
	allowOrigins := os.Getenv("CORS_ORIGINS")
	if allowOrigins == "" {
		allowOrigins = "http://localhost:5173, http://127.0.0.1:5173, http://localhost:4173, http://127.0.0.1:4173, http://172.17.82.108:5173"
	}

	// Implementamos el middleware CORS a nivel global usando app.Use() para interceptar todas las peticiones entrantes.
	app.Use(cors.New(cors.Config{
		// Lista de orígenes autorizados (localhost en desarrollo, dominio de Vercel en producción).
		AllowOrigins: allowOrigins,
		// Declaramos de forma explícita qué cabeceras (Headers) se permitirán en la comunicación.
		// Authorization es obligatoria: el frontend envía el JWT en cada petición autenticada.
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
	}))

	// Registramos las rutas inyectando la colección de controladores ya construidos.
	routes.SetupRoutes(app, ctrls)

	// Preparamos un canal para capturar señales de interrupción (Ctrl+C) o terminación.
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	// En una goroutine esperamos la señal para apagar el servidor de forma ordenada.
	go func() {
		<-quit
		log.Println("Apagando el servidor...")
		if err := app.Shutdown(); err != nil {
			log.Printf("Error al apagar el servidor: %v", err)
		}
	}()

	// Puerto de escucha configurable con PORT (defecto 3000): permite levantar
	// el servicio en paralelo cuando el 3000 está ocupado por otro proyecto.
	puerto := os.Getenv("PORT")
	if puerto == "" {
		puerto = "3000"
	}

	// Ponemos la aplicación a escuchar peticiones. (Bloquea el hilo de ejecución).
	if err := app.Listen(":" + puerto); err != nil {
		log.Printf("El servidor se detuvo: %v", err)
	}
}

// // main.go
// package main

// import (
// 	"github.com/gofiber/fiber/v2"
// 	"github.com/gofiber/fiber/v2/middleware/cors"
// )

// // Estructuras de datos
// type LoginRequest struct {
// 	Email    string `json:"email"`
// 	Password string `json:"password"`
// }

// type Producto struct {
// 	ID     int     `json:"id"`
// 	Nombre string  `json:"nombre"`
// 	Precio float64 `json:"precio"`
// 	Img    string  `json:"img"`
// }

// func main() {
// 	app := fiber.New()

// 	// Middleware CORS para permitir peticiones desde React (Vite corre en el puerto 5173)
// 	app.Use(cors.New(cors.Config{
// 		AllowOrigins: "http://localhost:5173",
// 		AllowHeaders: "Origin, Content-Type, Accept",
// 	}))

// 	// Endpoint para Login
// 	app.Post("/api/login", func(c *fiber.Ctx) error {
// 		var req LoginRequest
// 		if err := c.BodyParser(&req); err != nil {
// 			return c.Status(400).JSON(fiber.Map{"error": "Cuerpo de petición inválido"})
// 		}

// 		if req.Email == "admin@upse.edu.ec" && req.Password == "123456" {
// 			return c.JSON(fiber.Map{"token": "fake-jwt-token-123", "email": req.Email})
// 		}
// 		return c.Status(401).JSON(fiber.Map{"error": "Credenciales incorrectas"})
// 	})

// 	// Endpoint para obtener Productos
// 	app.Get("/api/productos", func(c *fiber.Ctx) error {
// 		productos := []Producto{
// 			{ID: 1, Nombre: "Serum Revitalizante", Precio: 45.00, Img: "https://picsum.photos/seed/serum/150"},
// 			{ID: 2, Nombre: "Crema Hidratante Pro", Precio: 32.50, Img: "https://picsum.photos/seed/crema/150"},
// 			{ID: 3, Nombre: "Tónico Purificante", Precio: 28.00, Img: "https://picsum.photos/seed/tonico/150"},
// 			{ID: 4, Nombre: "Mascarilla Nocturna", Precio: 50.00, Img: "https://picsum.photos/seed/mascarilla/150"},
// 		}
// 		return c.JSON(productos)
// 	})

// 	app.Listen(":3000")
// }
