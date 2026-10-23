// src/services/productosService.ts
// Capa de servicios del catálogo (U1-T4): consume la API real vía apiFetch.
// Los componentes que la usan no cambian: la promesa resuelve los mismos tipos.

import { apiFetch } from "./http";
import type { Pagina } from "./pagina";
import type { Producto } from "../data/productos";

// GET /api/productos (lista completa del catálogo)
export const getProductos = (): Promise<Producto[]> =>
  apiFetch<Producto[]>("/api/productos");

// GET /api/productos/:id (detalle de un producto)
export const getProductoById = (id: number): Promise<Producto> =>
  apiFetch<Producto>(`/api/productos/${id}`);

// ===== Back-office (solo admin) =====

// Payload editable de un producto (crear/actualizar).
export interface ProductoPayload {
  nombre: string;
  descripcion: string;
  precio: number;
  categoria: string;
  img: string;
  galeria: string[];
  stock: number;
  activo: boolean;
}

// GET /api/admin/productos?q=&page=&limit= (incluye inactivos; paginado)
export const getProductosAdmin = (
  q = "",
  page = 1,
  limit = 10
): Promise<Pagina<Producto>> => {
  const params = new URLSearchParams({ page: String(page), limit: String(limit) });
  if (q) params.set("q", q);
  return apiFetch<Pagina<Producto>>(`/api/admin/productos?${params.toString()}`);
};

// POST /api/productos (admin)
export const crearProducto = (p: ProductoPayload): Promise<Producto> =>
  apiFetch<Producto>("/api/productos", {
    method: "POST",
    body: JSON.stringify(p),
  });

// PUT /api/productos/:id (admin)
export const actualizarProducto = (
  id: number,
  p: ProductoPayload
): Promise<Producto> =>
  apiFetch<Producto>(`/api/productos/${id}`, {
    method: "PUT",
    body: JSON.stringify(p),
  });

// DELETE /api/productos/:id (admin · borrado lógico)
export const desactivarProducto = (id: number): Promise<{ message: string }> =>
  apiFetch<{ message: string }>(`/api/productos/${id}`, { method: "DELETE" });
