// src/services/usuariosService.ts
// Capa de servicios del back-office de usuarios (Fase 7).
import { apiFetch } from "./http";
import type { Pagina } from "./pagina";

export type RolUsuario = "admin" | "cliente";

export interface Usuario {
  id: number;
  email: string;
  rol: RolUsuario;
  nombre: string;
  activo: boolean;
  referido_id?: number;
}

// GET /api/usuarios?page&limit (admin)
export const getUsuarios = (page = 1, limit = 10): Promise<Pagina<Usuario>> =>
  apiFetch<Pagina<Usuario>>(`/api/usuarios?page=${page}&limit=${limit}`);

// PATCH /api/usuarios/:id (admin · cambia rol y/o activo)
export const actualizarUsuario = (
  id: number,
  cambios: { rol?: RolUsuario; activo?: boolean }
): Promise<Usuario> =>
  apiFetch<Usuario>(`/api/usuarios/${id}`, {
    method: "PATCH",
    body: JSON.stringify(cambios),
  });
