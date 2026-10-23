// src/services/pedidosService.ts
// Capa de servicios de pedidos (Fase 3): historial del cliente y back-office admin.
import { apiFetch } from "./http";
import type { Pagina } from "./pagina";

// Re-exportamos el sobre de paginación para los componentes que ya lo importan.
export type { Pagina };

// Estado del ciclo de vida de un pedido.
export type EstadoPedido =
  | "pendiente"
  | "preparando"
  | "enviado"
  | "entregado"
  | "cancelado";

export interface PedidoItem {
  productoId: number;
  nombre: string;
  cantidad: number;
  precio: number;
}

export interface Pedido {
  id: number;
  numero: string;
  clienteNombre: string;
  clienteEmail: string;
  ciudad: string;
  metodoPago: string;
  total: number;
  estado: EstadoPedido;
  usuarioId?: number;
  creadoEn?: string;
  items?: PedidoItem[];
}

// GET /api/pedidos/mios?page&limit (cliente autenticado)
export const getMisPedidos = (page = 1, limit = 10): Promise<Pagina<Pedido>> =>
  apiFetch<Pagina<Pedido>>(`/api/pedidos/mios?page=${page}&limit=${limit}`);

// GET /api/pedidos?estado&page&limit (admin)
export const getPedidosAdmin = (
  estado = "",
  page = 1,
  limit = 10
): Promise<Pagina<Pedido>> =>
  apiFetch<Pagina<Pedido>>(
    `/api/pedidos?page=${page}&limit=${limit}${
      estado ? `&estado=${encodeURIComponent(estado)}` : ""
    }`
  );

// GET /api/pedidos/:id (dueño o admin)
export const getPedidoById = (id: number): Promise<Pedido> =>
  apiFetch<Pedido>(`/api/pedidos/${id}`);

// PATCH /api/pedidos/:id/estado (admin)
export const cambiarEstadoPedido = (
  id: number,
  estado: EstadoPedido
): Promise<Pedido> =>
  apiFetch<Pedido>(`/api/pedidos/${id}/estado`, {
    method: "PATCH",
    body: JSON.stringify({ estado }),
  });
