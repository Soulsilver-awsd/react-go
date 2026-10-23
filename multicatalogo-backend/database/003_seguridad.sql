-- =============================================================================
-- 003_seguridad.sql - Hash bcrypt de las contraseñas (U2-T4)
-- Idempotente: se puede ejecutar más de una vez sobre un volumen existente.
-- Usa pgcrypto (extensión oficial de PostgreSQL) con crypt()/gen_salt('bf').
-- =============================================================================

-- Habilitamos la extensión que provee crypt() y gen_salt() (necesita permisos
-- de superusuario; el usuario POSTGRES_USER del compose los tiene).
CREATE EXTENSION IF NOT EXISTS pgcrypto;

-- Hasheamos cualquier password que aún esté en texto plano.
-- El patrón detecta hashes bcrypt existentes ($2a$, $2b$, $2y$) para no
-- re-hashearlos (evita doble hash y hace la migración idempotente).
UPDATE usuarios
SET password = crypt(password, gen_salt('bf', 10))
WHERE password !~ '^\$2[aby]\$';
