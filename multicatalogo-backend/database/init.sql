-- =============================================================================
-- init.sql - Estructura relacional y datos semilla
-- Fase 2 (persistencia) + Fase 4 (pedidos, comisiones y stock)
-- Se ejecuta automáticamente la primera vez que se crea el contenedor.
-- =============================================================================

-- Eliminamos las tablas si ya existen para permitir reinicializaciones limpias.
DROP TABLE IF EXISTS comisiones CASCADE;
DROP TABLE IF EXISTS liquidaciones CASCADE;
DROP TABLE IF EXISTS pedido_items CASCADE;
DROP TABLE IF EXISTS pedidos CASCADE;
DROP TABLE IF EXISTS usuarios CASCADE;
DROP TABLE IF EXISTS referidos CASCADE;
DROP TABLE IF EXISTS productos CASCADE;
DROP SEQUENCE IF EXISTS referidos_id_seq CASCADE;

-- -----------------------------------------------------------------------------
-- Extensión pgcrypto: provee crypt()/gen_salt('bf') para hash bcrypt.
-- -----------------------------------------------------------------------------
CREATE EXTENSION IF NOT EXISTS pgcrypto;

-- -----------------------------------------------------------------------------
-- Tabla usuarios: credenciales y rol (password con hash bcrypt).
-- nombre/activo/creado_en/referido_id alimentan el back-office (Fase 0).
-- El FK de referido_id se agrega tras crear `referidos` (referencia circular).
-- -----------------------------------------------------------------------------
CREATE TABLE usuarios (
    id         SERIAL PRIMARY KEY,
    email      VARCHAR(150) NOT NULL UNIQUE,
    password   VARCHAR(150) NOT NULL,
    rol        VARCHAR(30)  NOT NULL DEFAULT 'cliente',
    nombre     VARCHAR(150) NOT NULL DEFAULT '',
    activo     BOOLEAN      NOT NULL DEFAULT true,
    referido_id INTEGER,
    creado_en  TIMESTAMP    NOT NULL DEFAULT NOW()
);

-- -----------------------------------------------------------------------------
-- Tabla productos: catálogo con stock y galería (ARRAY de texto).
-- activo habilita el borrado lógico (Fase 2).
-- -----------------------------------------------------------------------------
CREATE TABLE productos (
    id          SERIAL PRIMARY KEY,
    nombre      VARCHAR(150)   NOT NULL,
    descripcion TEXT           NOT NULL DEFAULT '',
    precio      NUMERIC(10, 2) NOT NULL DEFAULT 0,
    categoria   VARCHAR(80)    NOT NULL DEFAULT '',
    img         TEXT           NOT NULL DEFAULT '',
    galeria     TEXT[]         NOT NULL DEFAULT '{}',
    stock       INTEGER        NOT NULL DEFAULT 100,
    activo      BOOLEAN        NOT NULL DEFAULT true
);

-- -----------------------------------------------------------------------------
-- Tabla referidos: árbol multinivel (parent_id) con saldo de comisiones y
-- código de invitación único. usuario_id enlaza opcionalmente un nodo con una
-- cuenta de `usuarios` (vínculo inverso de usuarios.referido_id).
-- -----------------------------------------------------------------------------
CREATE TABLE referidos (
    id         INTEGER PRIMARY KEY,
    nombre     VARCHAR(150)   NOT NULL,
    nivel      INTEGER        NOT NULL DEFAULT 0,
    ventas     NUMERIC(12, 2) NOT NULL DEFAULT 0,
    saldo      NUMERIC(12, 2) NOT NULL DEFAULT 0,
    codigo     VARCHAR(40)    UNIQUE,
    usuario_id INTEGER REFERENCES usuarios (id) ON DELETE SET NULL,
    parent_id  INTEGER REFERENCES referidos (id) ON DELETE CASCADE
);

-- Secuencia para que los nodos nuevos (registros) obtengan id automático.
CREATE SEQUENCE referidos_id_seq OWNED BY referidos.id;
ALTER TABLE referidos ALTER COLUMN id SET DEFAULT nextval('referidos_id_seq');

-- FK circular usuario↔referido: se agrega ahora que ambas tablas existen.
ALTER TABLE usuarios
    ADD CONSTRAINT usuarios_referido_id_fkey
    FOREIGN KEY (referido_id) REFERENCES referidos (id);

-- -----------------------------------------------------------------------------
-- Tabla pedidos: registro principal de cada compra (cabecera).
-- estado sigue el flujo pendiente→preparando→enviado→entregado (o cancelado).
-- usuario_id es la FK real al comprador (identidad del token).
-- -----------------------------------------------------------------------------
CREATE TABLE pedidos (
    id             SERIAL PRIMARY KEY,
    numero         VARCHAR(40)    NOT NULL UNIQUE,
    cliente_nombre VARCHAR(150)   NOT NULL DEFAULT '',
    cliente_email  VARCHAR(150)   NOT NULL DEFAULT '',
    ciudad         VARCHAR(120)   NOT NULL DEFAULT '',
    metodo_pago    VARCHAR(80)    NOT NULL DEFAULT '',
    total          NUMERIC(12, 2) NOT NULL DEFAULT 0,
    estado         VARCHAR(40)    NOT NULL DEFAULT 'pendiente',
    usuario_id     INTEGER        REFERENCES usuarios (id),
    creado_en      TIMESTAMP      NOT NULL DEFAULT now()
);

-- -----------------------------------------------------------------------------
-- Tabla pedido_items: detalle de cada pedido (relación 1:N con pedidos).
-- -----------------------------------------------------------------------------
CREATE TABLE pedido_items (
    id          SERIAL PRIMARY KEY,
    pedido_id   INTEGER        NOT NULL REFERENCES pedidos (id) ON DELETE CASCADE,
    producto_id INTEGER        NOT NULL REFERENCES productos (id),
    nombre      VARCHAR(150)   NOT NULL DEFAULT '',
    cantidad    INTEGER        NOT NULL DEFAULT 1,
    precio      NUMERIC(10, 2) NOT NULL DEFAULT 0
);

-- -----------------------------------------------------------------------------
-- Datos semilla: usuarios (coinciden con el login actual de la API).
-- Los passwords se guardan con hash bcrypt (pgcrypto); el login los valida
-- con bcrypt.CompareHashAndPassword en Go. referido_id se enlaza tras crear la red.
-- -----------------------------------------------------------------------------
INSERT INTO usuarios (email, password, rol, nombre) VALUES
    ('admin@upse.edu.ec',   crypt('123456', gen_salt('bf', 10)), 'admin',   'Administrador'),
    ('cliente@upse.edu.ec', crypt('123456', gen_salt('bf', 10)), 'cliente', 'Cliente Demo');

-- -----------------------------------------------------------------------------
-- Datos semilla: productos (los mismos del catálogo del frontend, con stock).
-- -----------------------------------------------------------------------------
INSERT INTO productos (id, nombre, descripcion, precio, categoria, img, galeria, stock) VALUES
    (1, 'Serum Revitalizante',
     'Serum concentrado con vitamina C y ácido hialurónico. Ilumina la piel, reduce manchas y aporta hidratación profunda desde la primera aplicación.',
     45.00, 'Serum', 'https://picsum.photos/seed/serum/600',
     ARRAY['https://picsum.photos/seed/serum1/600', 'https://picsum.photos/seed/serum2/600', 'https://picsum.photos/seed/serum3/600'], 100),
    (2, 'Crema Hidratante Pro',
     'Crema facial de textura ligera con niacinamida y manteca de karité. Hidratación de 24 horas y barrera cutánea fortalecida para todo tipo de piel.',
     32.50, 'Crema', 'https://picsum.photos/seed/crema/600',
     ARRAY['https://picsum.photos/seed/crema1/600', 'https://picsum.photos/seed/crema2/600', 'https://picsum.photos/seed/crema3/600'], 100),
    (3, 'Tónico Purificante',
     'Tónico sin alcohol con hamamelis y agua de rosas. Limpia poros, equilibra el pH y prepara la piel para el resto de la rutina diaria.',
     28.00, 'Tónico', 'https://picsum.photos/seed/tonico/600',
     ARRAY['https://picsum.photos/seed/tonico1/600', 'https://picsum.photos/seed/tonico2/600', 'https://picsum.photos/seed/tonico3/600'], 100),
    (4, 'Mascarilla Nocturna',
     'Mascarilla de noche con colágeno y vitamina E. Repara la piel mientras duermes y recupera la luminosidad al despertar.',
     50.00, 'Mascarilla', 'https://picsum.photos/seed/mascarilla/600',
     ARRAY['https://picsum.photos/seed/mascarilla1/600', 'https://picsum.photos/seed/mascarilla2/600', 'https://picsum.photos/seed/mascarilla3/600'], 100),
    (5, 'Serum Anti-Edad Retinol',
     'Serum de retinol encapsulado que suaviza líneas de expresión y mejora la firmeza. Uso nocturno recomendado.',
     58.00, 'Serum', 'https://picsum.photos/seed/retinol/600',
     ARRAY['https://picsum.photos/seed/retinol1/600', 'https://picsum.photos/seed/retinol2/600', 'https://picsum.photos/seed/retinol3/600'], 100),
    (6, 'Crema Contorno de Ojos',
     'Contorno de ojos con cafeína y péptidos. Reduce bolsas y ojeras, hidrata la zona más delicada del rostro.',
     26.00, 'Crema', 'https://picsum.photos/seed/ojos/600',
     ARRAY['https://picsum.photos/seed/ojos1/600', 'https://picsum.photos/seed/ojos2/600', 'https://picsum.photos/seed/ojos3/600'], 100),
    (7, 'Tónico Exfoliante AHA-BHA',
     'Exfoliación química suave con ácidos AHA y BHA. Renueva la textura de la piel y desobstruye los poros.',
     34.00, 'Tónico', 'https://picsum.photos/seed/exfoliante/600',
     ARRAY['https://picsum.photos/seed/exfoliante1/600', 'https://picsum.photos/seed/exfoliante2/600', 'https://picsum.photos/seed/exfoliante3/600'], 100),
    (8, 'Kit Rutina Completa',
     'Kit integral con serum, crema, tónico y mascarilla. La rutina perfecta para comenzar: limpieza, tratamiento e hidratación.',
     129.00, 'Kit', 'https://picsum.photos/seed/kit/600',
     ARRAY['https://picsum.photos/seed/kit1/600', 'https://picsum.photos/seed/kit2/600', 'https://picsum.photos/seed/kit3/600'], 100);

-- Ajustamos la secuencia del SERIAL tras insertar IDs explícitos.
SELECT setval('productos_id_seq', (SELECT MAX(id) FROM productos));

-- -----------------------------------------------------------------------------
-- Datos semilla: red multinivel de referidos (con código de invitación).
-- El nodo raíz (id 0, "Tú") es la cuenta admin. El nodo id 2 ("Luis Poveda")
-- representa a cliente@upse.edu.ec y cuelga directo de "Tú": al comprar,
-- solo su ascendiente inmediato (la raíz) gana comisión.
-- -----------------------------------------------------------------------------
INSERT INTO referidos (id, nombre, nivel, ventas, saldo, codigo, usuario_id, parent_id) VALUES
    (0, 'Tú',            0, 2400, 0, 'TU',      1, NULL),
    (1, 'Ana García',    1, 1200, 0, 'ANA',     NULL, 0),
    (2, 'Luis Poveda',   1,  850, 0, 'LUIS',    2, 0),
    (3, 'Marta Sánchez', 1,  430, 0, 'MARTA',   NULL, 0),
    (4, 'Carlos Ruiz',   2,  500, 0, 'CARLOS',  NULL, 1),
    (5, 'Sofía León',    2,  430, 0, 'SOFIA',   NULL, 1),
    (6, 'Marco Díaz',    2,  380, 0, 'MARCO',   NULL, 2),
    (7, 'Diana Paz',     3,  300, 0, 'DIANA',   NULL, 4),
    (8, 'Cliente Demo',  4,    0, 0, 'CLIENTE', NULL, 7);

-- Ajustamos la secuencia de referidos tras insertar IDs explícitos.
SELECT setval('referidos_id_seq', (SELECT MAX(id) FROM referidos));

-- Enlazamos el vínculo real usuario↔nodo (admin→raíz, cliente→Luis Poveda).
UPDATE usuarios SET referido_id = 0 WHERE email = 'admin@upse.edu.ec';
UPDATE usuarios SET referido_id = 2 WHERE email = 'cliente@upse.edu.ec';

-- -----------------------------------------------------------------------------
-- Tabla liquidaciones: pagos realizados a distribuidores (Fase 5).
-- -----------------------------------------------------------------------------
CREATE TABLE liquidaciones (
    id          SERIAL PRIMARY KEY,
    referido_id INTEGER        NOT NULL REFERENCES referidos (id) ON DELETE CASCADE,
    total       NUMERIC(12, 2) NOT NULL DEFAULT 0,
    estado      VARCHAR(30)    NOT NULL DEFAULT 'pagada',
    creado_en   TIMESTAMP      NOT NULL DEFAULT NOW(),
    pagada_en   TIMESTAMP
);

-- -----------------------------------------------------------------------------
-- Tabla comisiones: ledger auditable de comisiones MLM (Fase 4).
-- Cada fila = quién cobró qué porcentaje por qué pedido.
-- -----------------------------------------------------------------------------
CREATE TABLE comisiones (
    id             SERIAL PRIMARY KEY,
    pedido_id      INTEGER        NOT NULL REFERENCES pedidos (id) ON DELETE CASCADE,
    referido_id    INTEGER        NOT NULL REFERENCES referidos (id) ON DELETE CASCADE,
    nivel          INTEGER        NOT NULL,
    porcentaje     NUMERIC(5, 2)  NOT NULL,
    monto          NUMERIC(12, 2) NOT NULL DEFAULT 0,
    estado         VARCHAR(30)    NOT NULL DEFAULT 'pendiente',
    liquidacion_id INTEGER        REFERENCES liquidaciones (id),
    creado_en      TIMESTAMP      NOT NULL DEFAULT NOW()
);

-- -----------------------------------------------------------------------------
-- Índices de rendimiento (U2-T3).
-- -----------------------------------------------------------------------------
CREATE INDEX idx_comisiones_referido_pedido ON comisiones (referido_id, pedido_id);
CREATE INDEX idx_comisiones_pedido         ON comisiones (pedido_id);
CREATE INDEX idx_referidos_parent          ON referidos (parent_id);
CREATE INDEX idx_pedidos_cliente_email     ON pedidos (cliente_email);
CREATE INDEX idx_pedido_items_pedido       ON pedido_items (pedido_id);
CREATE INDEX idx_productos_categoria       ON productos (categoria);
