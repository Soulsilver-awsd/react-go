package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/jackc/pgx/v5"
	"golang.org/x/crypto/bcrypt"
)

type ProductoSeed struct {
	ID          int
	Nombre      string
	Descripcion string
	Precio      float64
	Categoria   string
	Img         string
	Galeria     []string
	Stock       int
}

type UsuarioSeed struct {
	Email      string
	Password   string
	Rol        string
	Nombre     string
	Activo     bool
	ReferidoID *int
}

type PedidoItemSeed struct {
	ProductoID int
	Nombre     string
	Cantidad   int
	Precio     float64
}

type PedidoSeed struct {
	Numero        string
	ClienteNombre string
	ClienteEmail  string
	Ciudad        string
	MetodoPago    string
	Total         float64
	Estado        string
	UsuarioEmail  string
	DiasAtras     int
	Items         []PedidoItemSeed
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	return string(bytes), err
}

func intPtr(i int) *int {
	return &i
}

func main() {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "postgresql://multicatalogo:w2LXrjHg9gp8teZ3DegAAGotP9X5U5Hm@dpg-dbdd9acs728c73bhli3g-a.oregon-postgres.render.com:5432/multicatalogo_ysra"
	}

	ctx := context.Background()
	conn, err := pgx.Connect(ctx, dsn)
	if err != nil {
		log.Fatalf("No se pudo conectar a PostgreSQL: %v", err)
	}
	defer conn.Close(ctx)

	log.Println("==> Conectado a PostgreSQL en Render exitosamente")

	// 1. PRODUCTOS ADICIONALES (IDs 9 al 20)
	nuevosProductos := []ProductoSeed{
		{
			ID:          9,
			Nombre:      "Limpiador Facial Espumoso",
			Descripcion: "Gel limpiador espumoso con ácido salicílico y extracto de té verde. Elimina el exceso de grasa e impurezas sin resecar.",
			Precio:      22.00,
			Categoria:   "Limpieza",
			Img:         "https://picsum.photos/seed/limpiador/600",
			Galeria:     []string{"https://picsum.photos/seed/limpiador1/600", "https://picsum.photos/seed/limpiador2/600"},
			Stock:       85,
		},
		{
			ID:          10,
			Nombre:      "Protector Solar FPS 50+ Invisible",
			Descripcion: "Protección solar facial de amplio espectro UVA/UVB con acabado mate invisible. Resistente al agua y textura no grasa.",
			Precio:      38.50,
			Categoria:   "Solar",
			Img:         "https://picsum.photos/seed/solar/600",
			Galeria:     []string{"https://picsum.photos/seed/solar1/600", "https://picsum.photos/seed/solar2/600"},
			Stock:       120,
		},
		{
			ID:          11,
			Nombre:      "Bálsamo Labial Reparador Karité",
			Descripcion: "Tratamiento intensivo para labios secos o agrietados enriquecido con manteca de karité pura y cera de abeja natural.",
			Precio:      14.00,
			Categoria:   "Cuidado Labial",
			Img:         "https://picsum.photos/seed/labial/600",
			Galeria:     []string{"https://picsum.photos/seed/labial1/600"},
			Stock:       150,
		},
		{
			ID:          12,
			Nombre:      "Serum Niacinamida 10% + Zinc",
			Descripcion: "Fórmula de alta potencia para reducir imperfecciones, regular el sebo y minimizar visiblemente el tamaño de los poros.",
			Precio:      36.00,
			Categoria:   "Serum",
			Img:         "https://picsum.photos/seed/niacinamida/600",
			Galeria:     []string{"https://picsum.photos/seed/niacinamida1/600", "https://picsum.photos/seed/niacinamida2/600"},
			Stock:       70,
		},
		{
			ID:          13,
			Nombre:      "Agua Micelar Calmante con Manzanilla",
			Descripcion: "Desmaquilla, limpia y tonifica en un solo paso. Formulada especialmente para pieles sensibles y propensas al enrojecimiento.",
			Precio:      19.50,
			Categoria:   "Limpieza",
			Img:         "https://picsum.photos/seed/micelar/600",
			Galeria:     []string{"https://picsum.photos/seed/micelar1/600"},
			Stock:       110,
		},
		{
			ID:          14,
			Nombre:      "Exfoliante Corporal Café & Coco",
			Descripcion: "Exfoliante energizante con granos de café orgánico y aceite de coco virgen. Estimula la circulación y suaviza la piel.",
			Precio:      27.00,
			Categoria:   "Corporal",
			Img:         "https://picsum.photos/seed/exfoliantecorporal/600",
			Galeria:     []string{"https://picsum.photos/seed/exfoliantecorporal1/600"},
			Stock:       95,
		},
		{
			ID:          15,
			Nombre:      "Loción Hidratante Corporal Ceramidas",
			Descripcion: "Loción ligera de absorción rápida con triple complejo de ceramidas para restaurar la barrera protectora de la piel.",
			Precio:      35.00,
			Categoria:   "Corporal",
			Img:         "https://picsum.photos/seed/locion/600",
			Galeria:     []string{"https://picsum.photos/seed/locion1/600"},
			Stock:       80,
		},
		{
			ID:          16,
			Nombre:      "Aceite Puro de Rosa Mosqueta",
			Descripcion: "100% puro y prensado en frío. Altamente regenerador para atenuar cicatrices, estrías y líneas de expresión.",
			Precio:      41.00,
			Categoria:   "Tratamiento",
			Img:         "https://picsum.photos/seed/rosamosqueta/600",
			Galeria:     []string{"https://picsum.photos/seed/rosamosqueta1/600"},
			Stock:       55,
		},
		{
			ID:          17,
			Nombre:      "Mascarilla Arcilla Verde Purificante",
			Descripcion: "Tratamiento détox semanal para purificar poros en profundidad y absorber el exceso de brillo sin irritar.",
			Precio:      29.90,
			Categoria:   "Mascarilla",
			Img:         "https://picsum.photos/seed/arcilla/600",
			Galeria:     []string{"https://picsum.photos/seed/arcilla1/600"},
			Stock:       65,
		},
		{
			ID:          18,
			Nombre:      "Bruma Facial Hidratante Lavanda",
			Descripcion: "Rocío facial calmante e hidratante con infusión de lavanda y glicerina vegetal. Ideal para fijar el maquillaje o refrescar.",
			Precio:      21.50,
			Categoria:   "Tónico",
			Img:         "https://picsum.photos/seed/bruma/600",
			Galeria:     []string{"https://picsum.photos/seed/bruma1/600"},
			Stock:       130,
		},
		{
			ID:          19,
			Nombre:      "Crema Reafirmante Cuello y Escote",
			Descripcion: "Tratamiento intensivo con péptidos tensores y extracto de algas para devolver firmeza y elasticidad a cuello y escote.",
			Precio:      48.00,
			Categoria:   "Crema",
			Img:         "https://picsum.photos/seed/cuello/600",
			Galeria:     []string{"https://picsum.photos/seed/cuello1/600"},
			Stock:       50,
		},
		{
			ID:          20,
			Nombre:      "Pack Spa Relajante & Regenerador",
			Descripcion: "Experiencia completa de spa en casa: bruma calmante, mascarilla purificante, exfoliante corporal y loción hidratante.",
			Precio:      115.00,
			Categoria:   "Kit",
			Img:         "https://picsum.photos/seed/packspa/600",
			Galeria:     []string{"https://picsum.photos/seed/packspa1/600", "https://picsum.photos/seed/packspa2/600"},
			Stock:       40,
		},
	}

	for _, p := range nuevosProductos {
		_, err := conn.Exec(ctx, `
			INSERT INTO productos (id, nombre, descripcion, precio, categoria, img, galeria, stock, activo)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, true)
			ON CONFLICT (id) DO UPDATE SET
				nombre = EXCLUDED.nombre,
				descripcion = EXCLUDED.descripcion,
				precio = EXCLUDED.precio,
				categoria = EXCLUDED.categoria,
				img = EXCLUDED.img,
				galeria = EXCLUDED.galeria,
				stock = EXCLUDED.stock,
				activo = true
		`, p.ID, p.Nombre, p.Descripcion, p.Precio, p.Categoria, p.Img, p.Galeria, p.Stock)
		if err != nil {
			log.Fatalf("Error insertando producto %d: %v", p.ID, err)
		}
	}
	_, _ = conn.Exec(ctx, "SELECT setval('productos_id_seq', (SELECT MAX(id) FROM productos))")
	log.Printf("✔ 20 Productos sincronizados correctamente")

	// 2. USUARIOS ADICIONALES
	hashedPass, err := hashPassword("123456")
	if err != nil {
		log.Fatalf("Error hasheando contraseña: %v", err)
	}

	usuarios := []UsuarioSeed{
		{Email: "ana.garcia@upse.edu.ec", Password: hashedPass, Rol: "cliente", Nombre: "Ana García", Activo: true, ReferidoID: intPtr(1)},
		{Email: "marta.sanchez@upse.edu.ec", Password: hashedPass, Rol: "cliente", Nombre: "Marta Sánchez", Activo: true, ReferidoID: intPtr(3)},
		{Email: "carlos.ruiz@upse.edu.ec", Password: hashedPass, Rol: "cliente", Nombre: "Carlos Ruiz", Activo: true, ReferidoID: intPtr(4)},
		{Email: "sofia.leon@upse.edu.ec", Password: hashedPass, Rol: "cliente", Nombre: "Sofía León", Activo: true, ReferidoID: intPtr(5)},
		{Email: "marco.diaz@upse.edu.ec", Password: hashedPass, Rol: "cliente", Nombre: "Marco Díaz", Activo: true, ReferidoID: intPtr(6)},
		{Email: "diana.paz@upse.edu.ec", Password: hashedPass, Rol: "cliente", Nombre: "Diana Paz", Activo: true, ReferidoID: intPtr(7)},
		{Email: "elena.torres@upse.edu.ec", Password: hashedPass, Rol: "cliente", Nombre: "Elena Torres", Activo: true, ReferidoID: nil},
		{Email: "javier.mendoza@upse.edu.ec", Password: hashedPass, Rol: "cliente", Nombre: "Javier Mendoza", Activo: true, ReferidoID: nil},
		{Email: "valeria.castro@upse.edu.ec", Password: hashedPass, Rol: "cliente", Nombre: "Valeria Castro", Activo: true, ReferidoID: nil},
		{Email: "fernando.reyes@upse.edu.ec", Password: hashedPass, Rol: "cliente", Nombre: "Fernando Reyes", Activo: false, ReferidoID: nil},
	}

	for _, u := range usuarios {
		_, err := conn.Exec(ctx, `
			INSERT INTO usuarios (email, password, rol, nombre, activo, referido_id)
			VALUES ($1, $2, $3, $4, $5, $6)
			ON CONFLICT (email) DO UPDATE SET
				nombre = EXCLUDED.nombre,
				activo = EXCLUDED.activo,
				referido_id = EXCLUDED.referido_id
		`, u.Email, u.Password, u.Rol, u.Nombre, u.Activo, u.ReferidoID)
		if err != nil {
			log.Fatalf("Error insertando usuario %s: %v", u.Email, err)
		}
	}
	log.Printf("✔ Usuarios sincronizados correctamente")

	// 3. MAPA DE USUARIOS (Email -> ID)
	userIDs := make(map[string]int)
	rows, err := conn.Query(ctx, "SELECT id, email FROM usuarios")
	if err != nil {
		log.Fatalf("Error consultando usuarios: %v", err)
	}
	for rows.Next() {
		var id int
		var email string
		_ = rows.Scan(&id, &email)
		userIDs[email] = id
	}
	rows.Close()

	// 4. ACTUALIZAR REFERIDOS (Ventas y Códigos)
	referidosActualizados := []struct {
		ID     int
		Nombre string
		Ventas float64
		Saldo  float64
		Codigo string
	}{
		{ID: 0, Nombre: "Tú", Ventas: 3450.00, Saldo: 125.00, Codigo: "TU"},
		{ID: 1, Nombre: "Ana García", Ventas: 1850.00, Saldo: 85.00, Codigo: "ANA"},
		{ID: 2, Nombre: "Luis Poveda", Ventas: 1420.00, Saldo: 65.50, Codigo: "LUIS"},
		{ID: 3, Nombre: "Marta Sánchez", Ventas: 980.00, Saldo: 45.00, Codigo: "MARTA"},
		{ID: 4, Nombre: "Carlos Ruiz", Ventas: 750.00, Saldo: 30.00, Codigo: "CARLOS"},
		{ID: 5, Nombre: "Sofía León", Ventas: 620.00, Saldo: 25.00, Codigo: "SOFIA"},
		{ID: 6, Nombre: "Marco Díaz", Ventas: 480.00, Saldo: 18.00, Codigo: "MARCO"},
		{ID: 7, Nombre: "Diana Paz", Ventas: 350.00, Saldo: 12.00, Codigo: "DIANA"},
		{ID: 8, Nombre: "Cliente Demo", Ventas: 210.00, Saldo: 0.00, Codigo: "CLIENTE"},
	}

	for _, r := range referidosActualizados {
		_, err := conn.Exec(ctx, `
			UPDATE referidos SET ventas = $1, saldo = $2, codigo = $3
			WHERE id = $4
		`, r.Ventas, r.Saldo, r.Codigo, r.ID)
		if err != nil {
			log.Fatalf("Error actualizando referido %d: %v", r.ID, err)
		}
	}
	log.Printf("✔ Red de referidos actualizada con ventas y saldos")

	// 5. PEDIDOS Y PEDIDO_ITEMS
	pedidos := []PedidoSeed{
		{
			Numero:        "PED-2026-1001",
			ClienteNombre: "Cliente Demo",
			ClienteEmail:  "cliente@upse.edu.ec",
			Ciudad:        "Guayaquil",
			MetodoPago:    "Tarjeta de Crédito",
			Total:         105.50,
			Estado:        "entregado",
			UsuarioEmail:  "cliente@upse.edu.ec",
			DiasAtras:     25,
			Items: []PedidoItemSeed{
				{ProductoID: 1, Nombre: "Serum Revitalizante", Cantidad: 1, Precio: 45.00},
				{ProductoID: 2, Nombre: "Crema Hidratante Pro", Cantidad: 1, Precio: 32.50},
				{ProductoID: 3, Nombre: "Tónico Purificante", Cantidad: 1, Precio: 28.00},
			},
		},
		{
			Numero:        "PED-2026-1002",
			ClienteNombre: "Ana García",
			ClienteEmail:  "ana.garcia@upse.edu.ec",
			Ciudad:        "Quito",
			MetodoPago:    "Transferencia Bancaria",
			Total:         179.00,
			Estado:        "entregado",
			UsuarioEmail:  "ana.garcia@upse.edu.ec",
			DiasAtras:     22,
			Items: []PedidoItemSeed{
				{ProductoID: 8, Nombre: "Kit Rutina Completa", Cantidad: 1, Precio: 129.00},
				{ProductoID: 4, Nombre: "Mascarilla Nocturna", Cantidad: 1, Precio: 50.00},
			},
		},
		{
			Numero:        "PED-2026-1003",
			ClienteNombre: "Carlos Ruiz",
			ClienteEmail:  "carlos.ruiz@upse.edu.ec",
			Ciudad:        "Cuenca",
			MetodoPago:    "Tarjeta de Débito",
			Total:         120.50,
			Estado:        "entregado",
			UsuarioEmail:  "carlos.ruiz@upse.edu.ec",
			DiasAtras:     19,
			Items: []PedidoItemSeed{
				{ProductoID: 5, Nombre: "Serum Anti-Edad Retinol", Cantidad: 1, Precio: 58.00},
				{ProductoID: 10, Nombre: "Protector Solar FPS 50+ Invisible", Cantidad: 1, Precio: 38.50},
				{ProductoID: 6, Nombre: "Crema Contorno de Ojos", Cantidad: 1, Precio: 26.00},
			},
		},
		{
			Numero:        "PED-2026-1004",
			ClienteNombre: "Cliente Demo",
			ClienteEmail:  "cliente@upse.edu.ec",
			Ciudad:        "Salinas",
			MetodoPago:    "Efectivo contra entrega",
			Total:         77.00,
			Estado:        "entregado",
			UsuarioEmail:  "cliente@upse.edu.ec",
			DiasAtras:     17,
			Items: []PedidoItemSeed{
				{ProductoID: 10, Nombre: "Protector Solar FPS 50+ Invisible", Cantidad: 2, Precio: 38.50},
			},
		},
		{
			Numero:        "PED-2026-1005",
			ClienteNombre: "Marta Sánchez",
			ClienteEmail:  "marta.sanchez@upse.edu.ec",
			Ciudad:        "Ambato",
			MetodoPago:    "Tarjeta de Crédito",
			Total:         115.00,
			Estado:        "entregado",
			UsuarioEmail:  "marta.sanchez@upse.edu.ec",
			DiasAtras:     15,
			Items: []PedidoItemSeed{
				{ProductoID: 20, Nombre: "Pack Spa Relajante & Regenerador", Cantidad: 1, Precio: 115.00},
			},
		},
		{
			Numero:        "PED-2026-1006",
			ClienteNombre: "Sofía León",
			ClienteEmail:  "sofia.leon@upse.edu.ec",
			Ciudad:        "Manta",
			MetodoPago:    "Transferencia Bancaria",
			Total:         87.90,
			Estado:        "entregado",
			UsuarioEmail:  "sofia.leon@upse.edu.ec",
			DiasAtras:     13,
			Items: []PedidoItemSeed{
				{ProductoID: 5, Nombre: "Serum Anti-Edad Retinol", Cantidad: 1, Precio: 58.00},
				{ProductoID: 17, Nombre: "Mascarilla Arcilla Verde Purificante", Cantidad: 1, Precio: 29.90},
			},
		},
		{
			Numero:        "PED-2026-1007",
			ClienteNombre: "Cliente Demo",
			ClienteEmail:  "cliente@upse.edu.ec",
			Ciudad:        "Guayaquil",
			MetodoPago:    "Tarjeta de Crédito",
			Total:         143.00,
			Estado:        "entregado",
			UsuarioEmail:  "cliente@upse.edu.ec",
			DiasAtras:     11,
			Items: []PedidoItemSeed{
				{ProductoID: 8, Nombre: "Kit Rutina Completa", Cantidad: 1, Precio: 129.00},
				{ProductoID: 11, Nombre: "Bálsamo Labial Reparador Karité", Cantidad: 1, Precio: 14.00},
			},
		},
		{
			Numero:        "PED-2026-1008",
			ClienteNombre: "Diana Paz",
			ClienteEmail:  "diana.paz@upse.edu.ec",
			Ciudad:        "Loja",
			MetodoPago:    "PayPal",
			Total:         96.00,
			Estado:        "enviado",
			UsuarioEmail:  "diana.paz@upse.edu.ec",
			DiasAtras:     9,
			Items: []PedidoItemSeed{
				{ProductoID: 12, Nombre: "Serum Niacinamida 10% + Zinc", Cantidad: 2, Precio: 36.00},
				{ProductoID: 11, Nombre: "Bálsamo Labial Reparador Karité", Cantidad: 1, Precio: 14.00},
				{ProductoID: 10, Nombre: "Protector Solar FPS 50+ Invisible", Cantidad: 1, Precio: 38.50},
			},
		},
		{
			Numero:        "PED-2026-1009",
			ClienteNombre: "Marco Díaz",
			ClienteEmail:  "marco.diaz@upse.edu.ec",
			Ciudad:        "Machala",
			MetodoPago:    "Tarjeta de Crédito",
			Total:         62.00,
			Estado:        "enviado",
			UsuarioEmail:  "marco.diaz@upse.edu.ec",
			DiasAtras:     7,
			Items: []PedidoItemSeed{
				{ProductoID: 14, Nombre: "Exfoliante Corporal Café & Coco", Cantidad: 1, Precio: 27.00},
				{ProductoID: 15, Nombre: "Loción Hidratante Corporal Ceramidas", Cantidad: 1, Precio: 35.00},
			},
		},
		{
			Numero:        "PED-2026-1010",
			ClienteNombre: "Elena Torres",
			ClienteEmail:  "elena.torres@upse.edu.ec",
			Ciudad:        "Guayaquil",
			MetodoPago:    "Transferencia Bancaria",
			Total:         110.00,
			Estado:        "enviado",
			UsuarioEmail:  "elena.torres@upse.edu.ec",
			DiasAtras:     6,
			Items: []PedidoItemSeed{
				{ProductoID: 16, Nombre: "Aceite Puro de Rosa Mosqueta", Cantidad: 1, Precio: 41.00},
				{ProductoID: 19, Nombre: "Crema Reafirmante Cuello y Escote", Cantidad: 1, Precio: 48.00},
				{ProductoID: 18, Nombre: "Bruma Facial Hidratante Lavanda", Cantidad: 1, Precio: 21.50},
			},
		},
		{
			Numero:        "PED-2026-1011",
			ClienteNombre: "Cliente Demo",
			ClienteEmail:  "cliente@upse.edu.ec",
			Ciudad:        "La Libertad",
			MetodoPago:    "Tarjeta de Crédito",
			Total:         58.00,
			Estado:        "enviado",
			UsuarioEmail:  "cliente@upse.edu.ec",
			DiasAtras:     5,
			Items: []PedidoItemSeed{
				{ProductoID: 5, Nombre: "Serum Anti-Edad Retinol", Cantidad: 1, Precio: 58.00},
			},
		},
		{
			Numero:        "PED-2026-1012",
			ClienteNombre: "Javier Mendoza",
			ClienteEmail:  "javier.mendoza@upse.edu.ec",
			Ciudad:        "Portoviejo",
			MetodoPago:    "Efectivo contra entrega",
			Total:         60.50,
			Estado:        "preparando",
			UsuarioEmail:  "javier.mendoza@upse.edu.ec",
			DiasAtras:     4,
			Items: []PedidoItemSeed{
				{ProductoID: 9, Nombre: "Limpiador Facial Espumoso", Cantidad: 1, Precio: 22.00},
				{ProductoID: 10, Nombre: "Protector Solar FPS 50+ Invisible", Cantidad: 1, Precio: 38.50},
			},
		},
		{
			Numero:        "PED-2026-1013",
			ClienteNombre: "Valeria Castro",
			ClienteEmail:  "valeria.castro@upse.edu.ec",
			Ciudad:        "Quito",
			MetodoPago:    "Tarjeta de Crédito",
			Total:         115.00,
			Estado:        "preparando",
			UsuarioEmail:  "valeria.castro@upse.edu.ec",
			DiasAtras:     3,
			Items: []PedidoItemSeed{
				{ProductoID: 20, Nombre: "Pack Spa Relajante & Regenerador", Cantidad: 1, Precio: 115.00},
			},
		},
		{
			Numero:        "PED-2026-1014",
			ClienteNombre: "Cliente Demo",
			ClienteEmail:  "cliente@upse.edu.ec",
			Ciudad:        "Guayaquil",
			MetodoPago:    "Tarjeta de Débito",
			Total:         79.00,
			Estado:        "preparando",
			UsuarioEmail:  "cliente@upse.edu.ec",
			DiasAtras:     2,
			Items: []PedidoItemSeed{
				{ProductoID: 7, Nombre: "Tónico Exfoliante AHA-BHA", Cantidad: 1, Precio: 34.00},
				{ProductoID: 1, Nombre: "Serum Revitalizante", Cantidad: 1, Precio: 45.00},
			},
		},
		{
			Numero:        "PED-2026-1015",
			ClienteNombre: "Ana García",
			ClienteEmail:  "ana.garcia@upse.edu.ec",
			Ciudad:        "Quito",
			MetodoPago:    "Transferencia Bancaria",
			Total:         69.50,
			Estado:        "pendiente",
			UsuarioEmail:  "ana.garcia@upse.edu.ec",
			DiasAtras:     1,
			Items: []PedidoItemSeed{
				{ProductoID: 13, Nombre: "Agua Micelar Calmante con Manzanilla", Cantidad: 1, Precio: 19.50},
				{ProductoID: 4, Nombre: "Mascarilla Nocturna", Cantidad: 1, Precio: 50.00},
			},
		},
		{
			Numero:        "PED-2026-1016",
			ClienteNombre: "Cliente Demo",
			ClienteEmail:  "cliente@upse.edu.ec",
			Ciudad:        "Salinas",
			MetodoPago:    "Tarjeta de Crédito",
			Total:         48.00,
			Estado:        "pendiente",
			UsuarioEmail:  "cliente@upse.edu.ec",
			DiasAtras:     0,
			Items: []PedidoItemSeed{
				{ProductoID: 19, Nombre: "Crema Reafirmante Cuello y Escote", Cantidad: 1, Precio: 48.00},
			},
		},
	}

	// Limpiamos pedidos y comisiones anteriores de seed para recrearlos de forma pura y vinculada
	_, _ = conn.Exec(ctx, "DELETE FROM pedido_items")
	_, _ = conn.Exec(ctx, "DELETE FROM comisiones")
	_, _ = conn.Exec(ctx, "DELETE FROM liquidaciones")
	_, _ = conn.Exec(ctx, "DELETE FROM pedidos")

	now := time.Now()
	pedidoIDs := make(map[string]int)

	for _, p := range pedidos {
		fecha := now.AddDate(0, 0, -p.DiasAtras)
		uid := userIDs[p.UsuarioEmail]

		var pedidoID int
		err := conn.QueryRow(ctx, `
			INSERT INTO pedidos (numero, cliente_nombre, cliente_email, ciudad, metodo_pago, total, estado, usuario_id, creado_en)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
			RETURNING id
		`, p.Numero, p.ClienteNombre, p.ClienteEmail, p.Ciudad, p.MetodoPago, p.Total, p.Estado, uid, fecha).Scan(&pedidoID)
		if err != nil {
			log.Fatalf("Error insertando pedido %s: %v", p.Numero, err)
		}
		pedidoIDs[p.Numero] = pedidoID

		for _, item := range p.Items {
			_, err := conn.Exec(ctx, `
				INSERT INTO pedido_items (pedido_id, producto_id, nombre, cantidad, precio)
				VALUES ($1, $2, $3, $4, $5)
			`, pedidoID, item.ProductoID, item.Nombre, item.Cantidad, item.Precio)
			if err != nil {
				log.Fatalf("Error insertando item de pedido %s: %v", p.Numero, err)
			}
		}
	}
	_, _ = conn.Exec(ctx, "SELECT setval('pedidos_id_seq', (SELECT MAX(id) FROM pedidos))")
	_, _ = conn.Exec(ctx, "SELECT setval('pedido_items_id_seq', (SELECT MAX(id) FROM pedido_items))")
	log.Printf("✔ %d Pedidos y sus ítems insertados correctamente", len(pedidos))

	// 6. HISTORIAL DE LIQUIDACIONES
	type LiqSeed struct {
		ID         int
		ReferidoID int
		Total      float64
		Estado     string
		DiasAtras  int
	}

	historialLiq := []LiqSeed{
		{ID: 1, ReferidoID: 0, Total: 185.00, Estado: "pagada", DiasAtras: 24},
		{ID: 2, ReferidoID: 1, Total: 110.50, Estado: "pagada", DiasAtras: 20},
		{ID: 3, ReferidoID: 2, Total: 95.00, Estado: "pagada", DiasAtras: 16},
		{ID: 4, ReferidoID: 4, Total: 65.00, Estado: "pagada", DiasAtras: 12},
		{ID: 5, ReferidoID: 1, Total: 75.00, Estado: "pagada", DiasAtras: 8},
	}

	for _, l := range historialLiq {
		fecha := now.AddDate(0, 0, -l.DiasAtras)
		_, err := conn.Exec(ctx, `
			INSERT INTO liquidaciones (id, referido_id, total, estado, creado_en, pagada_en)
			VALUES ($1, $2, $3, $4, $5, $6)
			ON CONFLICT (id) DO UPDATE SET
				referido_id = EXCLUDED.referido_id,
				total = EXCLUDED.total,
				estado = EXCLUDED.estado,
				creado_en = EXCLUDED.creado_en,
				pagada_en = EXCLUDED.pagada_en
		`, l.ID, l.ReferidoID, l.Total, l.Estado, fecha, fecha)
		if err != nil {
			log.Fatalf("Error insertando liquidacion %d: %v", l.ID, err)
		}
	}
	_, _ = conn.Exec(ctx, "SELECT setval('liquidaciones_id_seq', (SELECT MAX(id) FROM liquidaciones))")
	log.Printf("✔ Historial de liquidaciones creado")

	// 7. COMISIONES MLM (Históricas liquidadas + Recientes PENDIENTES de pago)
	type ComisionSeed struct {
		PedidoNumero  string
		ReferidoID    int
		Nivel         int
		Porcentaje    float64
		Monto         float64
		Estado        string
		LiquidacionID *int
		DiasAtras     int
	}

	comisiones := []ComisionSeed{
		// Comisiones ya pagadas en liquidaciones anteriores:
		{PedidoNumero: "PED-2026-1001", ReferidoID: 0, Nivel: 0, Porcentaje: 10.0, Monto: 10.55, Estado: "liquidada", LiquidacionID: intPtr(1), DiasAtras: 25},
		{PedidoNumero: "PED-2026-1001", ReferidoID: 1, Nivel: 1, Porcentaje: 5.0, Monto: 5.28, Estado: "liquidada", LiquidacionID: intPtr(2), DiasAtras: 25},
		{PedidoNumero: "PED-2026-1002", ReferidoID: 0, Nivel: 0, Porcentaje: 10.0, Monto: 17.90, Estado: "liquidada", LiquidacionID: intPtr(1), DiasAtras: 22},
		{PedidoNumero: "PED-2026-1003", ReferidoID: 1, Nivel: 1, Porcentaje: 5.0, Monto: 6.03, Estado: "liquidada", LiquidacionID: intPtr(2), DiasAtras: 19},
		{PedidoNumero: "PED-2026-1003", ReferidoID: 2, Nivel: 2, Porcentaje: 3.0, Monto: 3.62, Estado: "liquidada", LiquidacionID: intPtr(3), DiasAtras: 19},
		{PedidoNumero: "PED-2026-1004", ReferidoID: 0, Nivel: 0, Porcentaje: 10.0, Monto: 7.70, Estado: "liquidada", LiquidacionID: intPtr(1), DiasAtras: 17},
		{PedidoNumero: "PED-2026-1005", ReferidoID: 4, Nivel: 1, Porcentaje: 5.0, Monto: 5.75, Estado: "liquidada", LiquidacionID: intPtr(4), DiasAtras: 15},

		// Comisiones PENDIENTES (¡Aparecerán en la pestaña de Deudas para que el admin pueda liquidar!)
		{PedidoNumero: "PED-2026-1006", ReferidoID: 1, Nivel: 1, Porcentaje: 10.0, Monto: 8.79, Estado: "pendiente", LiquidacionID: nil, DiasAtras: 13},
		{PedidoNumero: "PED-2026-1006", ReferidoID: 0, Nivel: 0, Porcentaje: 5.0, Monto: 4.40, Estado: "pendiente", LiquidacionID: nil, DiasAtras: 13},
		{PedidoNumero: "PED-2026-1007", ReferidoID: 2, Nivel: 1, Porcentaje: 10.0, Monto: 14.30, Estado: "pendiente", LiquidacionID: nil, DiasAtras: 11},
		{PedidoNumero: "PED-2026-1007", ReferidoID: 0, Nivel: 0, Porcentaje: 5.0, Monto: 7.15, Estado: "pendiente", LiquidacionID: nil, DiasAtras: 11},
		{PedidoNumero: "PED-2026-1008", ReferidoID: 4, Nivel: 1, Porcentaje: 10.0, Monto: 9.60, Estado: "pendiente", LiquidacionID: nil, DiasAtras: 9},
		{PedidoNumero: "PED-2026-1008", ReferidoID: 1, Nivel: 2, Porcentaje: 5.0, Monto: 4.80, Estado: "pendiente", LiquidacionID: nil, DiasAtras: 9},
		{PedidoNumero: "PED-2026-1009", ReferidoID: 2, Nivel: 1, Porcentaje: 10.0, Monto: 6.20, Estado: "pendiente", LiquidacionID: nil, DiasAtras: 7},
		{PedidoNumero: "PED-2026-1010", ReferidoID: 3, Nivel: 1, Porcentaje: 10.0, Monto: 11.00, Estado: "pendiente", LiquidacionID: nil, DiasAtras: 6},
		{PedidoNumero: "PED-2026-1011", ReferidoID: 1, Nivel: 1, Porcentaje: 10.0, Monto: 5.80, Estado: "pendiente", LiquidacionID: nil, DiasAtras: 5},
		{PedidoNumero: "PED-2026-1012", ReferidoID: 5, Nivel: 1, Porcentaje: 10.0, Monto: 6.05, Estado: "pendiente", LiquidacionID: nil, DiasAtras: 4},
		{PedidoNumero: "PED-2026-1013", ReferidoID: 0, Nivel: 0, Porcentaje: 10.0, Monto: 11.50, Estado: "pendiente", LiquidacionID: nil, DiasAtras: 3},
		{PedidoNumero: "PED-2026-1014", ReferidoID: 2, Nivel: 1, Porcentaje: 10.0, Monto: 7.40, Estado: "pendiente", LiquidacionID: nil, DiasAtras: 2},
	}

	for _, c := range comisiones {
		pedID := pedidoIDs[c.PedidoNumero]
		fecha := now.AddDate(0, 0, -c.DiasAtras)
		_, err := conn.Exec(ctx, `
			INSERT INTO comisiones (pedido_id, referido_id, nivel, porcentaje, monto, estado, liquidacion_id, creado_en)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		`, pedID, c.ReferidoID, c.Nivel, c.Porcentaje, c.Monto, c.Estado, c.LiquidacionID, fecha)
		if err != nil {
			log.Fatalf("Error insertando comision para pedido %s: %v", c.PedidoNumero, err)
		}
	}
	_, _ = conn.Exec(ctx, "SELECT setval('comisiones_id_seq', (SELECT MAX(id) FROM comisiones))")
	log.Printf("✔ Ledger de comisiones poblado (pendientes y liquidadas)")

	fmt.Println("\n=======================================================")
	fmt.Println("🎉 ¡BASE DE DATOS POBLADA EXITOSAMENTE EN RENDER! 🎉")
	fmt.Println("=======================================================")
}
