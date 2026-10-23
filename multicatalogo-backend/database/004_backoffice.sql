-- =============================================================================
-- 004_backoffice.sql - Migración del back-office MLM (Fase 0)
-- Idempotente (IF NOT EXISTS / SET DEFAULT): aplicable a volúmenes existentes.
-- Añade: datos de usuario, borrado lógico, FK real usuario↔referido, códigos de
-- invitación, estados de pedido y el ledger de comisiones + liquidaciones.
-- =============================================================================

-- -----------------------------------------------------------------------------
-- usuarios: nombre, activo (login rechaza inactivos), vínculo real al nodo de
-- la red (referido_id) y fecha de creación.
-- -----------------------------------------------------------------------------
ALTER TABLE usuarios ADD COLUMN IF NOT EXISTS nombre      VARCHAR(150) NOT NULL DEFAULT '';
ALTER TABLE usuarios ADD COLUMN IF NOT EXISTS activo      BOOLEAN      NOT NULL DEFAULT true;
ALTER TABLE usuarios ADD COLUMN IF NOT EXISTS referido_id INTEGER      REFERENCES referidos (id);
ALTER TABLE usuarios ADD COLUMN IF NOT EXISTS creado_en   TIMESTAMP    NOT NULL DEFAULT NOW();

-- -----------------------------------------------------------------------------
-- referidos: código único de invitación (link /registro?ref=CODIGO).
-- -----------------------------------------------------------------------------
ALTER TABLE referidos ADD COLUMN IF NOT EXISTS codigo VARCHAR(40);
CREATE UNIQUE INDEX IF NOT EXISTS referidos_codigo_key ON referidos (codigo);

-- referidos.id no era SERIAL (los nodos del seed se insertan con id explícito).
-- Añadimos una secuencia para que los registros nuevos obtengan id automático.
CREATE SEQUENCE IF NOT EXISTS referidos_id_seq;
ALTER TABLE referidos ALTER COLUMN id SET DEFAULT nextval('referidos_id_seq');
ALTER SEQUENCE referidos_id_seq OWNED BY referidos.id;
SELECT setval('referidos_id_seq', (SELECT COALESCE(MAX(id), 0) FROM referidos));

-- -----------------------------------------------------------------------------
-- productos: borrado lógico (no rompe pedido_items históricos).
-- -----------------------------------------------------------------------------
ALTER TABLE productos ADD COLUMN IF NOT EXISTS activo BOOLEAN NOT NULL DEFAULT true;

-- -----------------------------------------------------------------------------
-- pedidos: estado del ciclo de vida + FK real al comprador.
-- -----------------------------------------------------------------------------
ALTER TABLE pedidos ADD COLUMN IF NOT EXISTS usuario_id INTEGER REFERENCES usuarios (id);
-- El estado por defecto pasa a ser 'pendiente' (máquina de estados de la Fase 3).
ALTER TABLE pedidos ALTER COLUMN estado SET DEFAULT 'pendiente';
UPDATE pedidos SET estado = 'pendiente' WHERE estado = 'confirmado';

-- -----------------------------------------------------------------------------
-- liquidaciones: pagos realizados a distribuidores.
-- -----------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS liquidaciones (
    id          SERIAL PRIMARY KEY,
    referido_id INTEGER        NOT NULL REFERENCES referidos (id) ON DELETE CASCADE,
    total       NUMERIC(12, 2) NOT NULL DEFAULT 0,
    estado      VARCHAR(30)    NOT NULL DEFAULT 'pagada',
    creado_en   TIMESTAMP      NOT NULL DEFAULT NOW(),
    pagada_en   TIMESTAMP
);

-- -----------------------------------------------------------------------------
-- comisiones: ledger auditable (fuente de verdad del dinero MLM).
-- Cada fila = quién cobró qué porcentaje por qué pedido.
-- -----------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS comisiones (
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
-- Índices de rendimiento (U2-T3): comisiones, red, pedidos y catálogo.
-- -----------------------------------------------------------------------------
CREATE INDEX IF NOT EXISTS idx_comisiones_referido_pedido ON comisiones (referido_id, pedido_id);
CREATE INDEX IF NOT EXISTS idx_comisiones_pedido         ON comisiones (pedido_id);
CREATE INDEX IF NOT EXISTS idx_referidos_parent          ON referidos (parent_id);
CREATE INDEX IF NOT EXISTS idx_pedidos_cliente_email     ON pedidos (cliente_email);
CREATE INDEX IF NOT EXISTS idx_pedido_items_pedido       ON pedido_items (pedido_id);
CREATE INDEX IF NOT EXISTS idx_productos_categoria       ON productos (categoria);

-- -----------------------------------------------------------------------------
-- Códigos de invitación para el seed + vínculo real usuario↔nodo.
-- admin (id 1) → nodo 0 "Tú"  (raíz)
-- cliente (id 2) → nodo 2 "Luis Poveda" (hijo directo de "Tú")
-- -----------------------------------------------------------------------------
UPDATE referidos SET codigo = CASE id
    WHEN 0 THEN 'TU'
    WHEN 1 THEN 'ANA'
    WHEN 2 THEN 'LUIS'
    WHEN 3 THEN 'MARTA'
    WHEN 4 THEN 'CARLOS'
    WHEN 5 THEN 'SOFIA'
    WHEN 6 THEN 'MARCO'
    WHEN 7 THEN 'DIANA'
    WHEN 8 THEN 'CLIENTE'
    ELSE codigo
END
WHERE codigo IS NULL;

UPDATE usuarios SET referido_id = 0 WHERE email = 'admin@upse.edu.ec';
UPDATE usuarios SET referido_id = 2 WHERE email = 'cliente@upse.edu.ec';
UPDATE usuarios SET nombre = 'Administrador'        WHERE email = 'admin@upse.edu.ec'   AND nombre = '';
UPDATE usuarios SET nombre = 'Cliente Demo'         WHERE email = 'cliente@upse.edu.ec' AND nombre = '';

-- Alineamos el vínculo inverso: el nodo de Luis Poveda (id 2) es la cuenta
-- cliente; el nodo 8 "Cliente Demo" queda como referido de muestra sin cuenta.
UPDATE referidos SET usuario_id = 2    WHERE id = 2;
UPDATE referidos SET usuario_id = NULL WHERE id = 8;
