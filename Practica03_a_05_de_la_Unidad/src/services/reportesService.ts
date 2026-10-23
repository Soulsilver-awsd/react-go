// src/services/reportesService.ts
// Capa de servicios de reportes del dashboard (Fase 6).
import { apiFetch } from "./http";
import type { DeudaDistribuidor } from "./liquidacionesService";

export interface ResumenReporte {
  totalVendido: number;
  numPedidos: number;
  ticketPromedio: number;
  comisionesDevengadas: number;
}

export interface TopProducto {
  productoId: number;
  nombre: string;
  cantidad: number;
  total: number;
}

export interface TopReferido {
  referidoId: number;
  nombre: string;
  codigo?: string;
  ventas: number;
  saldo: number;
}

export interface FiltroFechas {
  desde?: string;
  hasta?: string;
}

const qs = (f: FiltroFechas): string => {
  const p = new URLSearchParams();
  if (f.desde) p.set("desde", f.desde);
  if (f.hasta) p.set("hasta", f.hasta);
  const s = p.toString();
  return s ? `?${s}` : "";
};

// GET /api/reportes/resumen?desde&hasta
export const getResumen = (f: FiltroFechas = {}): Promise<ResumenReporte> =>
  apiFetch<ResumenReporte>(`/api/reportes/resumen${qs(f)}`);

// GET /api/reportes/top-productos?desde&hasta
export const getTopProductos = (f: FiltroFechas = {}): Promise<TopProducto[]> =>
  apiFetch<TopProducto[]>(`/api/reportes/top-productos${qs(f)}`);

// GET /api/reportes/top-referidos
export const getTopReferidos = (): Promise<TopReferido[]> =>
  apiFetch<TopReferido[]>("/api/reportes/top-referidos");

// GET /api/reportes/comisiones-pendientes
export const getComisionesPendientes = (): Promise<DeudaDistribuidor[]> =>
  apiFetch<DeudaDistribuidor[]>("/api/reportes/comisiones-pendientes");
