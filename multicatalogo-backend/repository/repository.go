// Declaramos el paquete repository, responsable exclusivo del acceso a datos (PostgreSQL).
package repository

import (
	// Context define el alcance de cada consulta a la base de datos.
	"context"
	// Errors permite crear y comparar errores centinela del repositorio.
	"errors"
	// Fmt construye el número de pedido y los códigos de referido.
	"fmt"
	// Strings normaliza el nombre para generar el código de referido.
	"strings"
	// Time se usa para generar identificadores únicos de pedido.
	"time"

	// Pgconn expone el tipo de error del protocolo PostgreSQL (códigos SQLSTATE).
	"github.com/jackc/pgx/v5/pgconn"
	// Pgx provee los errores y utilidades del driver PostgreSQL.
	"github.com/jackc/pgx/v5"
	// Pgxpool provee el pool de conexiones inyectado.
	"github.com/jackc/pgx/v5/pgxpool"

	// Importamos los modelos que devuelven las consultas.
	"multicatalogo-backend/models"
	// Importamos la lógica de negocio pura para repartir comisiones.
	"multicatalogo-backend/domain"
)

// ErrNoEncontrado se retorna cuando una consulta no arroja resultados.
var ErrNoEncontrado = errors.New("registro no encontrado")

// ErrEmailDuplicado se retorna al intentar registrar un email ya existente.
var ErrEmailDuplicado = errors.New("el email ya está registrado")

// ErrReferidoInvalido se retorna cuando el código de referido no existe.
var ErrReferidoInvalido = errors.New("el código de referido no existe")

// ErrTransicionInvalida se retorna al intentar un cambio de estado no permitido.
var ErrTransicionInvalida = errors.New("transición de estado no permitida")

// ErrSinComisiones se retorna al liquidar un distribuidor sin saldo pendiente.
var ErrSinComisiones = errors.New("el distribuidor no tiene comisiones pendientes")

// ErrCantidadInvalida se retorna cuando una línea del pedido pide una cantidad
// menor o igual a cero. Es un error de validación del cliente (400), distinto
// del error de stock insuficiente.
var ErrCantidadInvalida = errors.New("la cantidad debe ser mayor que cero")

// transicionesPedido es la máquina de estados del ciclo de vida de un pedido.
// Un estado terminal (entregado/cancelado) no admite más cambios.
var transicionesPedido = map[string][]string{
	"pendiente":  {"preparando", "cancelado"},
	"preparando": {"enviado", "cancelado"},
	"enviado":    {"entregado", "cancelado"},
	"entregado":  {},
	"cancelado":  {},
}

// DBRepository es la interfaz que abstrae el acceso a datos (inyección de dependencias).
type DBRepository interface {
	// BuscarUsuarioPorEmail recupera un usuario por su email (login y checkout).
	// El password devuelto es el hash bcrypt; su verificación ocurre en el controlador.
	BuscarUsuarioPorEmail(email string) (*models.Usuario, error)
	// CrearUsuario registra un cliente y su nodo en la red en una sola transacción.
	// Si refCodigo es válido, el nodo cuelga del referente; si no, queda sin padre.
	CrearUsuario(nombre, email, passwordHash, refCodigo string) (*models.Usuario, error)
	// ObtenerProductos devuelve el catálogo con los productos activos.
	ObtenerProductos() ([]models.Producto, error)
	// ObtenerProductosAdmin devuelve todos los productos (activos e inactivos)
	// de forma paginada, con filtro de búsqueda opcional por nombre/descripción.
	ObtenerProductosAdmin(q string, page, limit int) (*models.Pagina[models.Producto], error)
	// ObtenerProductoPorID devuelve un producto por su identificador.
	ObtenerProductoPorID(id int) (*models.Producto, error)
	// CrearProducto inserta un producto nuevo en el catálogo.
	CrearProducto(p models.Producto) (*models.Producto, error)
	// ActualizarProducto edita un producto existente (edición completa).
	ActualizarProducto(id int, p models.Producto) (*models.Producto, error)
	// DesactivarProducto aplica borrado lógico (activo=false).
	DesactivarProducto(id int) error
	// ObtenerRed reconstruye el árbol multinivel global de referidos.
	ObtenerRed() (*models.Referido, error)
	// ObtenerRedDeUsuario reconstruye el subárbol del nodo asociado al usuario.
	ObtenerRedDeUsuario(usuarioID int) (*models.Referido, error)
	// CrearPedido registra una compra dentro de una transacción ACID e inserta
	// el ledger de comisiones (niveles 1..3) con estado 'pendiente'.
	CrearPedido(req models.CheckoutRequest, usuarioID int) (*models.Pedido, error)
	// ActualizarVentasRed incrementa el contador de ventas de los ascendentes.
	ActualizarVentasRed(total float64, usuarioID int) error
	// ObtenerPedidosPorUsuario lista los pedidos del comprador (paginado).
	ObtenerPedidosPorUsuario(usuarioID, page, limit int) (*models.Pagina[models.Pedido], error)
	// ObtenerPedidosAdmin lista todos los pedidos con filtro opcional de estado.
	ObtenerPedidosAdmin(estado string, page, limit int) (*models.Pagina[models.Pedido], error)
	// ObtenerPedidoPorID devuelve un pedido con su detalle (items).
	ObtenerPedidoPorID(id int) (*models.Pedido, error)
	// CambiarEstadoPedido valida la transición y, si es cancelación, restaura
	// stock y revierte comisiones dentro de una transacción ACID.
	CambiarEstadoPedido(id int, nuevoEstado string) (*models.Pedido, error)
	// ObtenerComisionesPorUsuario lista el ledger de comisiones del usuario.
	ObtenerComisionesPorUsuario(usuarioID int) ([]models.Comision, float64, float64, error)
	// DeudasPorDistribuidor suma las comisiones pendientes agrupadas por referido.
	DeudasPorDistribuidor() ([]models.DeudaDistribuidor, error)
	// PagarLiquidacion liquida (paga) todas las comisiones pendientes de un
	// distribuidor en una transacción idempotente (bloquea las filas).
	PagarLiquidacion(referidoID int) (*models.Liquidacion, error)
	// ObtenerLiquidaciones devuelve el historial global de liquidaciones (admin).
	ObtenerLiquidaciones() ([]models.Liquidacion, error)
	// ObtenerLiquidacionesPorUsuario devuelve los cobros del distribuidor autenticado.
	ObtenerLiquidacionesPorUsuario(usuarioID int) ([]models.Liquidacion, error)
	// ObtenerResumenReporte agrega los KPIs del negocio (Fase 6).
	ObtenerResumenReporte(desde, hasta string) (*models.ResumenReporte, error)
	// ObtenerTopProductos rankea los productos más vendidos.
	ObtenerTopProductos(desde, hasta string) ([]models.TopProducto, error)
	// ObtenerTopReferidos rankea los nodos de la red por ventas.
	ObtenerTopReferidos() ([]models.TopReferido, error)
	// ObtenerUsuarios lista las cuentas de forma paginada (back-office).
	ObtenerUsuarios(page, limit int) (*models.Pagina[models.Usuario], error)
	// ActualizarUsuario cambia rol y/o activo (nils = no tocar ese campo).
	ActualizarUsuario(id int, rol *string, activo *bool) (*models.Usuario, error)
}

// PostgresRepository es la implementación de DBRepository sobre PostgreSQL.
type PostgresRepository struct {
	// db es el pool de conexiones inyectado al construir el repositorio.
	db *pgxpool.Pool
}

// NewPostgresRepository construye el repositorio con el pool de conexiones inyectado.
func NewPostgresRepository(db *pgxpool.Pool) *PostgresRepository {
	return &PostgresRepository{db: db}
}

// BuscarUsuarioPorEmail recupera un usuario por su email (login y checkout).
// El campo Password contiene el hash bcrypt almacenado en PostgreSQL.
func (r *PostgresRepository) BuscarUsuarioPorEmail(email string) (*models.Usuario, error) {
	row := r.db.QueryRow(context.Background(),
		`SELECT id, email, password, rol, nombre, activo, referido_id
		 FROM usuarios WHERE email = $1`, email)

	var u models.Usuario
	if err := row.Scan(&u.ID, &u.Email, &u.Password, &u.Rol, &u.Nombre, &u.Activo, &u.ReferidoID); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNoEncontrado
		}
		return nil, err
	}
	return &u, nil
}

// CrearUsuario registra un cliente y su nodo en la red en una sola transacción:
//  1. valida el código de referido (si viene) y resuelve su parent_id
//  2. inserta el usuario con el hash bcrypt (rol cliente)
//  3. inserta el nodo `referidos` colgando del referente
//  4. enlaza usuarios.referido_id con el nodo creado
//
// Si algo falla, el ROLLBACK deja la base intacta.
func (r *PostgresRepository) CrearUsuario(nombre, email, passwordHash, refCodigo string) (*models.Usuario, error) {
	ctx := context.Background()
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	// 1) Resolvemos el parent_id a partir del código de referido (opcional).
	var parentID *int
	if refCodigo != "" {
		var pid int
		err := tx.QueryRow(ctx, "SELECT id FROM referidos WHERE codigo = $1", refCodigo).Scan(&pid)
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrReferidoInvalido
		}
		if err != nil {
			return nil, err
		}
		parentID = &pid
	}

	// 2) Insertamos el usuario; un email duplicado dispara la violación UNIQUE.
	var u models.Usuario
	err = tx.QueryRow(ctx,
		`INSERT INTO usuarios (nombre, email, password, rol)
		 VALUES ($1, $2, $3, 'cliente')
		 RETURNING id, email, rol, nombre, activo`,
		nombre, email, passwordHash).Scan(&u.ID, &u.Email, &u.Rol, &u.Nombre, &u.Activo)
	if err != nil {
		// 23505 = unique_violation (email ya registrado).
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return nil, ErrEmailDuplicado
		}
		return nil, err
	}

	// 3) Creamos el nodo de la red con un código de invitación propio.
	codigo := generarCodigo(nombre, u.ID)
	var nodoID int
	err = tx.QueryRow(ctx,
		`INSERT INTO referidos (nombre, nivel, codigo, usuario_id, parent_id)
		 VALUES ($1, COALESCE((SELECT nivel + 1 FROM referidos WHERE id = $2), 0), $3, $4, $5)
		 RETURNING id`,
		nombre, parentID, codigo, u.ID, parentID).Scan(&nodoID)
	if err != nil {
		return nil, err
	}

	// 4) Enlazamos el vínculo real usuario ↔ nodo (fuente de cálculo de comisiones).
	if _, err := tx.Exec(ctx, "UPDATE usuarios SET referido_id = $1 WHERE id = $2", nodoID, u.ID); err != nil {
		return nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	u.ReferidoID = &nodoID
	return &u, nil
}

// generarCodigo construye un código de invitación único a partir del nombre.
// Ej: "Ana García" + id 12 → "ANA12". Simple y legible para el link de referido.
func generarCodigo(nombre string, id int) string {
	// Nos quedamos solo con letras y pasamos a mayúsculas.
	var b strings.Builder
	for _, r := range strings.ToUpper(nombre) {
		if (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') {
			b.WriteRune(r)
		}
		if b.Len() >= 8 {
			break
		}
	}
	base := b.String()
	if base == "" {
		base = "USER"
	}
	return fmt.Sprintf("%s%d", base, id)
}

// ObtenerProductos devuelve el catálogo público: solo productos activos.
func (r *PostgresRepository) ObtenerProductos() ([]models.Producto, error) {
	return r.consultarProductos("SELECT id, nombre, descripcion, precio, categoria, img, galeria, stock, activo FROM productos WHERE activo = true ORDER BY id", nil)
}

// ObtenerProductosAdmin devuelve TODOS los productos (incluidos inactivos) de
// forma paginada; el filtro q busca por nombre o descripción (ILIKE).
func (r *PostgresRepository) ObtenerProductosAdmin(q string, page, limit int) (*models.Pagina[models.Producto], error) {
	ctx := context.Background()
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}
	offset := (page - 1) * limit

	where := ""
	var args []interface{}
	if q != "" {
		where = " WHERE nombre ILIKE $1 OR descripcion ILIKE $1"
		args = append(args, "%"+q+"%")
	}

	var total int
	if err := r.db.QueryRow(ctx, "SELECT COUNT(*) FROM productos"+where, args...).Scan(&total); err != nil {
		return nil, err
	}

	query := "SELECT id, nombre, descripcion, precio, categoria, img, galeria, stock, activo FROM productos" + where +
		fmt.Sprintf(" ORDER BY id LIMIT $%d OFFSET $%d", len(args)+1, len(args)+2)
	rows, err := r.db.Query(ctx, query, append(args, limit, offset)...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	productos := []models.Producto{}
	for rows.Next() {
		var p models.Producto
		if err := rows.Scan(&p.ID, &p.Nombre, &p.Descripcion, &p.Precio, &p.Categoria, &p.Img, &p.Galeria, &p.Stock, &p.Activo); err != nil {
			return nil, err
		}
		productos = append(productos, p)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return &models.Pagina[models.Producto]{Data: productos, Page: page, Limit: limit, Total: total}, nil
}

// consultarProductos centraliza el escaneo de filas de productos.
func (r *PostgresRepository) consultarProductos(query string, args []interface{}) ([]models.Producto, error) {
	rows, err := r.db.Query(context.Background(), query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	productos := []models.Producto{}
	for rows.Next() {
		var p models.Producto
		if err := rows.Scan(&p.ID, &p.Nombre, &p.Descripcion, &p.Precio, &p.Categoria, &p.Img, &p.Galeria, &p.Stock, &p.Activo); err != nil {
			return nil, err
		}
		productos = append(productos, p)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return productos, nil
}

// ObtenerProductoPorID devuelve un único producto según su identificador.
func (r *PostgresRepository) ObtenerProductoPorID(id int) (*models.Producto, error) {
	row := r.db.QueryRow(context.Background(),
		"SELECT id, nombre, descripcion, precio, categoria, img, galeria, stock, activo FROM productos WHERE id = $1", id)

	var p models.Producto
	if err := row.Scan(&p.ID, &p.Nombre, &p.Descripcion, &p.Precio, &p.Categoria, &p.Img, &p.Galeria, &p.Stock, &p.Activo); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNoEncontrado
		}
		return nil, err
	}
	return &p, nil
}

// CrearProducto inserta un producto nuevo y devuelve la fila creada (con su id).
func (r *PostgresRepository) CrearProducto(p models.Producto) (*models.Producto, error) {
	// Galería nil se persiste como array vacío para respetar el NOT NULL.
	if p.Galeria == nil {
		p.Galeria = []string{}
	}
	err := r.db.QueryRow(context.Background(),
		`INSERT INTO productos (nombre, descripcion, precio, categoria, img, galeria, stock, activo)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, true)
		 RETURNING id`,
		p.Nombre, p.Descripcion, p.Precio, p.Categoria, p.Img, p.Galeria, p.Stock).Scan(&p.ID)
	if err != nil {
		return nil, err
	}
	p.Activo = true
	return &p, nil
}

// ActualizarProducto edita todos los campos editables de un producto existente.
// Devuelve ErrNoEncontrado si el id no corresponde a ninguna fila.
func (r *PostgresRepository) ActualizarProducto(id int, p models.Producto) (*models.Producto, error) {
	if p.Galeria == nil {
		p.Galeria = []string{}
	}
	ct, err := r.db.Exec(context.Background(),
		`UPDATE productos
		 SET nombre = $1, descripcion = $2, precio = $3, categoria = $4,
		     img = $5, galeria = $6, stock = $7, activo = $8
		 WHERE id = $9`,
		p.Nombre, p.Descripcion, p.Precio, p.Categoria, p.Img, p.Galeria, p.Stock, p.Activo, id)
	if err != nil {
		return nil, err
	}
	if ct.RowsAffected() == 0 {
		return nil, ErrNoEncontrado
	}
	p.ID = id
	return &p, nil
}

// DesactivarProducto aplica borrado lógico: la fila permanece (no rompe los
// pedido_items históricos) pero deja de aparecer en el catálogo público.
func (r *PostgresRepository) DesactivarProducto(id int) error {
	ct, err := r.db.Exec(context.Background(), "UPDATE productos SET activo = false WHERE id = $1", id)
	if err != nil {
		return err
	}
	if ct.RowsAffected() == 0 {
		return ErrNoEncontrado
	}
	return nil
}

// ObtenerRed reconstruye el árbol multinivel de referidos a partir de la lista plana de PostgreSQL.
func (r *PostgresRepository) ObtenerRed() (*models.Referido, error) {
	// Consultamos todas las filas incluyendo el parent_id para poder enlazar el árbol.
	rows, err := r.db.Query(context.Background(),
		"SELECT id, nombre, nivel, ventas, saldo, parent_id FROM referidos ORDER BY id")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Mapa auxiliar id -> datos planos del nodo.
	nodos := map[int]models.Referido{}
	// Mapa auxiliar id -> lista ordenada de ids de sus hijos.
	hijosIDs := map[int][]int{}
	// Identificador de la raíz (nodo sin padre).
	var raizID int
	var hayRaiz bool

	for rows.Next() {
		var (
			id       int
			nombre   string
			nivel    int
			ventas   float64
			saldo    float64
			parentID *int
		)
		if err := rows.Scan(&id, &nombre, &nivel, &ventas, &saldo, &parentID); err != nil {
			return nil, err
		}
		// Guardamos los datos planos del nodo.
		nodos[id] = models.Referido{ID: id, Nombre: nombre, Nivel: nivel, Ventas: ventas, Saldo: saldo}
		if parentID == nil {
			// El nodo sin padre es la raíz del árbol.
			raizID, hayRaiz = id, true
			continue
		}
		// Registramos el id del nodo como hijo de su padre, preservando el orden de inserción.
		hijosIDs[*parentID] = append(hijosIDs[*parentID], id)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	if !hayRaiz {
		return nil, ErrNoEncontrado
	}

	// Construimos el árbol de forma recursiva a partir de la raíz.
	raiz := construirArbol(raizID, nodos, hijosIDs)
	return &raiz, nil
}

// construirArbol arma recursivamente un nodo y todos sus descendientes.
func construirArbol(id int, nodos map[int]models.Referido, hijosIDs map[int][]int) models.Referido {
	// Partimos de una copia de los datos planos del nodo actual.
	nodo := nodos[id]
	// Recorremos los ids de sus hijos para ensamblarlos como subárboles.
	for _, hijoID := range hijosIDs[id] {
		nodo.Hijos = append(nodo.Hijos, construirArbol(hijoID, nodos, hijosIDs))
	}
	return nodo
}

// CrearPedido registra la cabecera, el detalle y el ledger de comisiones de una
// compra dentro de una transacción ACID. El ledger nace aquí (no en la goroutine):
// cada nivel ascendente 1..3 obtiene una fila 'pendiente' con su porcentaje.
func (r *PostgresRepository) CrearPedido(req models.CheckoutRequest, usuarioID int) (*models.Pedido, error) {
	// Abrimos un contexto para toda la transacción.
	ctx := context.Background()

	// Iniciamos la transacción (BEGIN). A partir de aquí todo es atómico.
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return nil, err
	}
	// Si la función termina sin COMMIT, el defer ejecuta ROLLBACK automáticamente.
	defer tx.Rollback(ctx)

	// Generamos un número de pedido único.
	numero := fmt.Sprintf("MC-%d", time.Now().UnixNano())

	// Paso 1: insertamos la cabecera (estado inicial 'pendiente') con su comprador.
	var pedidoID int
	err = tx.QueryRow(ctx,
		`INSERT INTO pedidos (numero, cliente_nombre, cliente_email, ciudad, metodo_pago, total, estado, usuario_id)
		 VALUES ($1, $2, $3, $4, $5, $6, 'pendiente', $7)
		 RETURNING id`,
		numero, req.Cliente.Nombre, req.Cliente.Email, req.Cliente.Ciudad, req.MetodoPago, req.Total, usuarioID,
	).Scan(&pedidoID)
	if err != nil {
		return nil, err
	}

	// Paso 2: insertamos el detalle y descontamos el stock de cada producto.
	for _, item := range req.Items {
		// Validación del cliente: cantidad debe ser positiva (error propio, 400).
		if item.Cantidad <= 0 {
			return nil, ErrCantidadInvalida
		}
		// Descontamos stock. Si el producto no existe, no se afecta ninguna fila y abortamos.
		ct, err := tx.Exec(ctx,
			"UPDATE productos SET stock = stock - $1 WHERE id = $2",
			item.Cantidad, item.ID)
		if err != nil {
			return nil, err
		}
		if ct.RowsAffected() == 0 {
			// Producto inexistente: provocamos el ROLLBACK de toda la transacción.
			return nil, ErrNoEncontrado
		}

		// Insertamos la línea del detalle (la FK también protege ante productos inválidos).
		_, err = tx.Exec(ctx,
			`INSERT INTO pedido_items (pedido_id, producto_id, nombre, cantidad, precio)
			 VALUES ($1, $2, $3, $4, $5)`,
			pedidoID, item.ID, item.Nombre, item.Cantidad, item.Precio)
		if err != nil {
			return nil, err
		}
	}

	// Paso 3: insertamos el ledger de comisiones dentro de la MISMA transacción.
	// Recorremos los ascendentes del comprador (nivel 1..3) y les asignamos su parte.
	ascendentes, err := r.ascendentes(ctx, tx, usuarioID, len(domain.TasasPorNivel))
	if err != nil {
		return nil, err
	}
	for i, referidoID := range ascendentes {
		c := domain.RepartirComisiones(req.Total, domain.TasasPorNivel)[i]
		_, err := tx.Exec(ctx,
			`INSERT INTO comisiones (pedido_id, referido_id, nivel, porcentaje, monto, estado)
			 VALUES ($1, $2, $3, $4, $5, 'pendiente')`,
			pedidoID, referidoID, c.Nivel, c.Porcentaje*100, c.Monto)
		if err != nil {
			return nil, err
		}
	}

	// Paso 4: confirmamos la transacción (COMMIT) solo si todo salió bien.
	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	// Devolvemos el pedido recién creado.
	return &models.Pedido{
		ID:            pedidoID,
		Numero:        numero,
		ClienteNombre: req.Cliente.Nombre,
		ClienteEmail:  req.Cliente.Email,
		Ciudad:        req.Cliente.Ciudad,
		MetodoPago:    req.MetodoPago,
		Total:         req.Total,
		Estado:        "pendiente",
		UsuarioID:     &usuarioID,
	}, nil
}

// ascendentes devuelve los ids de los nodos referidos ascendentes del usuario
// (nivel 1 = padre directo, nivel 2 = abuelo, ...) hasta `maxNiveles`.
// Recibe un "querier" para reutilizar la misma conexión/transacción.
func (r *PostgresRepository) ascendentes(ctx context.Context, q interface {
	QueryRow(context.Context, string, ...interface{}) pgx.Row
}, usuarioID, maxNiveles int) ([]int, error) {
	// Resolvemos el nodo de la red vinculado al usuario.
	var nodoID *int
	if err := q.QueryRow(ctx, "SELECT referido_id FROM usuarios WHERE id = $1", usuarioID).Scan(&nodoID); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	if nodoID == nil {
		return nil, nil
	}

	// Subimos por la jerarquía recogiendo el id de cada ascendente.
	ascendentes := []int{}
	actual := *nodoID
	for nivel := 1; nivel <= maxNiveles; nivel++ {
		var parentID *int
		if err := q.QueryRow(ctx, "SELECT parent_id FROM referidos WHERE id = $1", actual).Scan(&parentID); err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				break
			}
			return nil, err
		}
		if parentID == nil {
			break
		}
		ascendentes = append(ascendentes, *parentID)
		actual = *parentID
	}
	return ascendentes, nil
}

// ActualizarVentasRed incrementa el contador `ventas` de los ascendentes del
// comprador. La goroutine solo actualiza este contador: el dinero ya quedó
// auditado en el ledger (tabla comisiones) durante CrearPedido.
func (r *PostgresRepository) ActualizarVentasRed(total float64, usuarioID int) error {
	ctx := context.Background()

	ascendentes, err := r.ascendentes(ctx, r.db, usuarioID, len(domain.TasasPorNivel))
	if err != nil {
		return err
	}
	for _, referidoID := range ascendentes {
		if _, err := r.db.Exec(ctx,
			"UPDATE referidos SET ventas = ventas + $1 WHERE id = $2",
			total, referidoID); err != nil {
			return err
		}
	}
	return nil
}

// ObtenerRedDeUsuario reconstruye el subárbol multinivel que cuelga del nodo
// asociado al usuario autenticado (no el árbol global del seed).
func (r *PostgresRepository) ObtenerRedDeUsuario(usuarioID int) (*models.Referido, error) {
	ctx := context.Background()

	// Nodo raíz del usuario.
	var raizID *int
	if err := r.db.QueryRow(ctx, "SELECT referido_id FROM usuarios WHERE id = $1", usuarioID).Scan(&raizID); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNoEncontrado
		}
		return nil, err
	}
	if raizID == nil {
		return nil, ErrNoEncontrado
	}

	// Cargamos todos los nodos y armamos el árbol, pero devolvemos solo el subárbol
	// cuya raíz es el nodo del usuario.
	rows, err := r.db.Query(ctx,
		"SELECT id, nombre, nivel, ventas, saldo, codigo, parent_id FROM referidos ORDER BY id")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	nodos := map[int]models.Referido{}
	hijosIDs := map[int][]int{}
	for rows.Next() {
		var (
			id       int
			nombre   string
			nivel    int
			ventas   float64
			saldo    float64
			codigo   *string
			parentID *int
		)
		if err := rows.Scan(&id, &nombre, &nivel, &ventas, &saldo, &codigo, &parentID); err != nil {
			return nil, err
		}
		n := models.Referido{ID: id, Nombre: nombre, Nivel: nivel, Ventas: ventas, Saldo: saldo}
		if codigo != nil {
			n.Codigo = *codigo
		}
		nodos[id] = n
		if parentID != nil {
			hijosIDs[*parentID] = append(hijosIDs[*parentID], id)
		}
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	raiz := construirArbol(*raizID, nodos, hijosIDs)
	return &raiz, nil
}

// ObtenerPedidosPorUsuario lista los pedidos del comprador autenticado (paginado).
func (r *PostgresRepository) ObtenerPedidosPorUsuario(usuarioID, page, limit int) (*models.Pagina[models.Pedido], error) {
	return r.consultarPedidos("WHERE usuario_id = $1", []interface{}{usuarioID}, page, limit)
}

// ObtenerPedidosAdmin lista todos los pedidos, con filtro opcional por estado.
func (r *PostgresRepository) ObtenerPedidosAdmin(estado string, page, limit int) (*models.Pagina[models.Pedido], error) {
	if estado == "" {
		return r.consultarPedidos("", nil, page, limit)
	}
	return r.consultarPedidos("WHERE estado = $1", []interface{}{estado}, page, limit)
}

// consultarPedidos centraliza el listado paginado de pedidos.
func (r *PostgresRepository) consultarPedidos(where string, args []interface{}, page, limit int) (*models.Pagina[models.Pedido], error) {
	ctx := context.Background()
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}
	offset := (page - 1) * limit

	// Total de filas que cumplen el filtro (para el sobre de paginación).
	var total int
	if err := r.db.QueryRow(ctx, "SELECT COUNT(*) FROM pedidos "+where, args...).Scan(&total); err != nil {
		return nil, err
	}

	// Página de resultados.
	q := "SELECT id, numero, cliente_nombre, cliente_email, ciudad, metodo_pago, total, estado, usuario_id, creado_en FROM pedidos " + where +
		fmt.Sprintf(" ORDER BY id DESC LIMIT $%d OFFSET $%d", len(args)+1, len(args)+2)
	rows, err := r.db.Query(ctx, q, append(args, limit, offset)...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	pedidos := []models.Pedido{}
	for rows.Next() {
		var p models.Pedido
		var creado time.Time
		if err := rows.Scan(&p.ID, &p.Numero, &p.ClienteNombre, &p.ClienteEmail, &p.Ciudad, &p.MetodoPago, &p.Total, &p.Estado, &p.UsuarioID, &creado); err != nil {
			return nil, err
		}
		p.CreadoEn = creado.Format(time.RFC3339)
		pedidos = append(pedidos, p)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return &models.Pagina[models.Pedido]{Data: pedidos, Page: page, Limit: limit, Total: total}, nil
}

// ObtenerPedidoPorID devuelve la cabecera de un pedido y su detalle (items).
func (r *PostgresRepository) ObtenerPedidoPorID(id int) (*models.Pedido, error) {
	ctx := context.Background()

	var p models.Pedido
	var creado time.Time
	err := r.db.QueryRow(ctx,
		`SELECT id, numero, cliente_nombre, cliente_email, ciudad, metodo_pago, total, estado, usuario_id, creado_en
		 FROM pedidos WHERE id = $1`, id).
		Scan(&p.ID, &p.Numero, &p.ClienteNombre, &p.ClienteEmail, &p.Ciudad, &p.MetodoPago, &p.Total, &p.Estado, &p.UsuarioID, &creado)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNoEncontrado
		}
		return nil, err
	}
	p.CreadoEn = creado.Format(time.RFC3339)

	// Detalle con el nombre del producto (JOIN; el nombre se conserva aunque el
	// producto se borre lógicamente).
	rows, err := r.db.Query(ctx,
		`SELECT pi.producto_id, pi.nombre, pi.cantidad, pi.precio
		 FROM pedido_items pi WHERE pi.pedido_id = $1 ORDER BY pi.id`, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	p.Items = []models.PedidoItem{}
	for rows.Next() {
		var it models.PedidoItem
		if err := rows.Scan(&it.ProductoID, &it.Nombre, &it.Cantidad, &it.Precio); err != nil {
			return nil, err
		}
		p.Items = append(p.Items, it)
	}
	return &p, rows.Err()
}

// CambiarEstadoPedido valida la transición y, si el nuevo estado es 'cancelado',
// restaura el stock y revierte las comisiones del ledger en la misma transacción.
func (r *PostgresRepository) CambiarEstadoPedido(id int, nuevoEstado string) (*models.Pedido, error) {
	ctx := context.Background()

	tx, err := r.db.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	// Bloqueamos la fila para evitar carreras sobre el estado.
	var estadoActual string
	if err := tx.QueryRow(ctx, "SELECT estado FROM pedidos WHERE id = $1 FOR UPDATE", id).Scan(&estadoActual); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNoEncontrado
		}
		return nil, err
	}

	// Validamos que la transición esté permitida por la máquina de estados.
	permitida := false
	for _, s := range transicionesPedido[estadoActual] {
		if s == nuevoEstado {
			permitida = true
			break
		}
	}
	if !permitida {
		return nil, ErrTransicionInvalida
	}

	// Cancelación: espejo ACID de la compra (stock + clawback de comisiones).
	if nuevoEstado == "cancelado" {
		// Devolvemos el stock de cada ítem al catálogo.
		if _, err := tx.Exec(ctx,
			`UPDATE productos p SET stock = p.stock + pi.cantidad
			 FROM pedido_items pi WHERE pi.pedido_id = $1 AND p.id = pi.producto_id`, id); err != nil {
			return nil, err
		}
		// Las comisiones de este pedido dejan de ser exigibles.
		if _, err := tx.Exec(ctx,
			"UPDATE comisiones SET estado = 'revertida' WHERE pedido_id = $1 AND estado = 'pendiente'", id); err != nil {
			return nil, err
		}
	}

	// Aplicamos el nuevo estado.
	if _, err := tx.Exec(ctx, "UPDATE pedidos SET estado = $1 WHERE id = $2", nuevoEstado, id); err != nil {
		return nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}
	return r.ObtenerPedidoPorID(id)
}

// ObtenerComisionesPorUsuario devuelve el ledger del usuario (su nodo de la red),
// junto con los totales pendiente y pagado.
func (r *PostgresRepository) ObtenerComisionesPorUsuario(usuarioID int) ([]models.Comision, float64, float64, error) {
	ctx := context.Background()

	// Resolvemos el nodo de la red del usuario.
	var nodoID *int
	if err := r.db.QueryRow(ctx, "SELECT referido_id FROM usuarios WHERE id = $1", usuarioID).Scan(&nodoID); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return []models.Comision{}, 0, 0, nil
		}
		return nil, 0, 0, err
	}
	if nodoID == nil {
		return []models.Comision{}, 0, 0, nil
	}

	rows, err := r.db.Query(ctx,
		`SELECT id, pedido_id, referido_id, nivel, porcentaje, monto, estado, creado_en
		 FROM comisiones WHERE referido_id = $1 ORDER BY id DESC`, *nodoID)
	if err != nil {
		return nil, 0, 0, err
	}
	defer rows.Close()

	comisiones := []models.Comision{}
	var totalPendiente, totalPagado float64
	for rows.Next() {
		var c models.Comision
		var creado time.Time
		if err := rows.Scan(&c.ID, &c.PedidoID, &c.ReferidoID, &c.Nivel, &c.Porcentaje, &c.Monto, &c.Estado, &creado); err != nil {
			return nil, 0, 0, err
		}
		c.CreadoEn = creado.Format(time.RFC3339)
		switch c.Estado {
		case "pendiente":
			totalPendiente += c.Monto
		case "pagada":
			totalPagado += c.Monto
		}
		comisiones = append(comisiones, c)
	}
	return comisiones, totalPendiente, totalPagado, rows.Err()
}

// DeudasPorDistribuidor agrupa las comisiones pendientes por referido (lo que
// se le debe a cada distribuidor). Base de la vista admin de liquidaciones.
func (r *PostgresRepository) DeudasPorDistribuidor() ([]models.DeudaDistribuidor, error) {
	rows, err := r.db.Query(context.Background(),
		`SELECT c.referido_id, r.nombre, COALESCE(SUM(c.monto), 0) AS pendiente
		 FROM comisiones c
		 JOIN referidos r ON r.id = c.referido_id
		 WHERE c.estado = 'pendiente'
		 GROUP BY c.referido_id, r.nombre
		 HAVING COALESCE(SUM(c.monto), 0) > 0
		 ORDER BY pendiente DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	deudas := []models.DeudaDistribuidor{}
	for rows.Next() {
		var d models.DeudaDistribuidor
		if err := rows.Scan(&d.ReferidoID, &d.Nombre, &d.Pendiente); err != nil {
			return nil, err
		}
		deudas = append(deudas, d)
	}
	return deudas, rows.Err()
}

// PagarLiquidacion paga todas las comisiones 'pendiente' de un distribuidor en
// una transacción idempotente. Bloquea las filas con FOR UPDATE para impedir
// que dos liquidaciones concurrentes paguen la misma comisión dos veces.
func (r *PostgresRepository) PagarLiquidacion(referidoID int) (*models.Liquidacion, error) {
	ctx := context.Background()

	tx, err := r.db.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	// Bloqueamos las comisiones pendientes y calculamos su total.
	rows, err := tx.Query(ctx,
		"SELECT monto FROM comisiones WHERE referido_id = $1 AND estado = 'pendiente' FOR UPDATE", referidoID)
	if err != nil {
		return nil, err
	}
	var total float64
	for rows.Next() {
		var monto float64
		if err := rows.Scan(&monto); err != nil {
			rows.Close()
			return nil, err
		}
		total += monto
	}
	rows.Close()
	if err := rows.Err(); err != nil {
		return nil, err
	}
	if total == 0 {
		return nil, ErrSinComisiones
	}

	// Creamos la liquidación (ya pagada, con su fecha).
	var liquidacionID int
	var creado, pagada time.Time
	if err := tx.QueryRow(ctx,
		`INSERT INTO liquidaciones (referido_id, total, estado, pagada_en)
		 VALUES ($1, $2, 'pagada', NOW())
		 RETURNING id, creado_en, pagada_en`,
		referidoID, total).Scan(&liquidacionID, &creado, &pagada); err != nil {
		return nil, err
	}

	// Marcamos las comisiones como pagadas y las enlazamos a la liquidación.
	if _, err := tx.Exec(ctx,
		"UPDATE comisiones SET estado = 'pagada', liquidacion_id = $1 WHERE referido_id = $2 AND estado = 'pendiente'",
		liquidacionID, referidoID); err != nil {
		return nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	return &models.Liquidacion{
		ID:         liquidacionID,
		ReferidoID: referidoID,
		Total:      total,
		Estado:     "pagada",
		CreadoEn:   creado.Format(time.RFC3339),
		PagadaEn:   pagada.Format(time.RFC3339),
	}, nil
}

// ObtenerLiquidaciones devuelve el historial completo de liquidaciones (admin),
// con el nombre del distribuidor.
func (r *PostgresRepository) ObtenerLiquidaciones() ([]models.Liquidacion, error) {
	return r.consultarLiquidaciones(
		`SELECT l.id, l.referido_id, l.total, l.estado, l.creado_en, l.pagada_en, r.nombre
		 FROM liquidaciones l JOIN referidos r ON r.id = l.referido_id
		 ORDER BY l.id DESC`, nil)
}

// ObtenerLiquidacionesPorUsuario devuelve los cobros del distribuidor autenticado.
func (r *PostgresRepository) ObtenerLiquidacionesPorUsuario(usuarioID int) ([]models.Liquidacion, error) {
	return r.consultarLiquidaciones(
		`SELECT l.id, l.referido_id, l.total, l.estado, l.creado_en, l.pagada_en, r.nombre
		 FROM liquidaciones l JOIN referidos r ON r.id = l.referido_id
		 WHERE l.referido_id = (SELECT referido_id FROM usuarios WHERE id = $1)
		 ORDER BY l.id DESC`, []interface{}{usuarioID})
}

// consultarLiquidaciones centraliza el escaneo de liquidaciones.
func (r *PostgresRepository) consultarLiquidaciones(query string, args []interface{}) ([]models.Liquidacion, error) {
	rows, err := r.db.Query(context.Background(), query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	liquidaciones := []models.Liquidacion{}
	for rows.Next() {
		var l models.Liquidacion
		var creado time.Time
		var pagada *time.Time
		if err := rows.Scan(&l.ID, &l.ReferidoID, &l.Total, &l.Estado, &creado, &pagada, &l.ReferidoNombre); err != nil {
			return nil, err
		}
		l.CreadoEn = creado.Format(time.RFC3339)
		if pagada != nil {
			l.PagadaEn = pagada.Format(time.RFC3339)
		}
		liquidaciones = append(liquidaciones, l)
	}
	return liquidaciones, rows.Err()
}

// fechaArg convierte una fecha opcional "" en nil (NULL SQL) para los filtros.
func fechaArg(s string) interface{} {
	if s == "" {
		return nil
	}
	return s
}

// ObtenerResumenReporte agrega ventas, nº de pedidos, ticket promedio y
// comisiones devengadas en el periodo [desde, hasta] (ambos opcionales).
func (r *PostgresRepository) ObtenerResumenReporte(desde, hasta string) (*models.ResumenReporte, error) {
	ctx := context.Background()
	rep := &models.ResumenReporte{}

	if err := r.db.QueryRow(ctx,
		`SELECT COALESCE(SUM(total), 0), COUNT(*)
		 FROM pedidos
		 WHERE estado <> 'cancelado'
		   AND ($1::timestamptz IS NULL OR creado_en >= $1::timestamptz)
		   AND ($2::timestamptz IS NULL OR creado_en <= $2::timestamptz)`,
		fechaArg(desde), fechaArg(hasta)).Scan(&rep.TotalVendido, &rep.NumPedidos); err != nil {
		return nil, err
	}

	if err := r.db.QueryRow(ctx,
		`SELECT COALESCE(SUM(monto), 0)
		 FROM comisiones
		 WHERE estado <> 'revertida'
		   AND ($1::timestamptz IS NULL OR creado_en >= $1::timestamptz)
		   AND ($2::timestamptz IS NULL OR creado_en <= $2::timestamptz)`,
		fechaArg(desde), fechaArg(hasta)).Scan(&rep.ComisionesDevengadas); err != nil {
		return nil, err
	}

	if rep.NumPedidos > 0 {
		rep.TicketPromedio = rep.TotalVendido / float64(rep.NumPedidos)
	}
	return rep, nil
}

// ObtenerTopProductos rankea los productos más vendidos por dinero en el periodo.
func (r *PostgresRepository) ObtenerTopProductos(desde, hasta string) ([]models.TopProducto, error) {
	rows, err := r.db.Query(context.Background(),
		`SELECT pi.producto_id, COALESCE(p2.nombre, MAX(pi.nombre)) AS nombre,
		        SUM(pi.cantidad)::int, SUM(pi.cantidad * pi.precio)
		 FROM pedido_items pi
		 JOIN pedidos p ON p.id = pi.pedido_id
		 LEFT JOIN productos p2 ON p2.id = pi.producto_id
		 WHERE p.estado <> 'cancelado'
		   AND ($1::timestamptz IS NULL OR p.creado_en >= $1::timestamptz)
		   AND ($2::timestamptz IS NULL OR p.creado_en <= $2::timestamptz)
		 GROUP BY pi.producto_id, p2.nombre
		 ORDER BY SUM(pi.cantidad * pi.precio) DESC
		 LIMIT 10`, fechaArg(desde), fechaArg(hasta))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	top := []models.TopProducto{}
	for rows.Next() {
		var t models.TopProducto
		if err := rows.Scan(&t.ProductoID, &t.Nombre, &t.Cantidad, &t.Total); err != nil {
			return nil, err
		}
		top = append(top, t)
	}
	return top, rows.Err()
}

// ObtenerTopReferidos rankea los nodos de la red según las ventas de su subárbol.
func (r *PostgresRepository) ObtenerTopReferidos() ([]models.TopReferido, error) {
	rows, err := r.db.Query(context.Background(),
		`SELECT id, nombre, COALESCE(codigo, ''), ventas, saldo
		 FROM referidos
		 ORDER BY ventas DESC
		 LIMIT 10`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	top := []models.TopReferido{}
	for rows.Next() {
		var t models.TopReferido
		if err := rows.Scan(&t.ReferidoID, &t.Nombre, &t.Codigo, &t.Ventas, &t.Saldo); err != nil {
			return nil, err
		}
		top = append(top, t)
	}
	return top, rows.Err()
}

// ObtenerUsuarios lista las cuentas de forma paginada (back-office de usuarios).
func (r *PostgresRepository) ObtenerUsuarios(page, limit int) (*models.Pagina[models.Usuario], error) {
	ctx := context.Background()
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}
	offset := (page - 1) * limit

	var total int
	if err := r.db.QueryRow(ctx, "SELECT COUNT(*) FROM usuarios").Scan(&total); err != nil {
		return nil, err
	}

	rows, err := r.db.Query(ctx,
		`SELECT id, email, rol, nombre, activo, referido_id
		 FROM usuarios ORDER BY id LIMIT $1 OFFSET $2`, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	usuarios := []models.Usuario{}
	for rows.Next() {
		var u models.Usuario
		if err := rows.Scan(&u.ID, &u.Email, &u.Rol, &u.Nombre, &u.Activo, &u.ReferidoID); err != nil {
			return nil, err
		}
		usuarios = append(usuarios, u)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return &models.Pagina[models.Usuario]{Data: usuarios, Page: page, Limit: limit, Total: total}, nil
}

// ActualizarUsuario cambia el rol y/o el estado activo de una cuenta. Se pasan
// punteros para distinguir "no cambiar" (nil) de un valor concreto.
func (r *PostgresRepository) ActualizarUsuario(id int, rol *string, activo *bool) (*models.Usuario, error) {
	var u models.Usuario
	err := r.db.QueryRow(context.Background(),
		`UPDATE usuarios
		 SET rol = COALESCE($2, rol), activo = COALESCE($3, activo)
		 WHERE id = $1
		 RETURNING id, email, rol, nombre, activo, referido_id`,
		id, rol, activo).
		Scan(&u.ID, &u.Email, &u.Rol, &u.Nombre, &u.Activo, &u.ReferidoID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNoEncontrado
		}
		return nil, err
	}
	return &u, nil
}
