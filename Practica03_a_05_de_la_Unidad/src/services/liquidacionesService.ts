// src/services/liquidacionesService.ts
// Capa de servicios de liquidaciones (Fase 5): pagos a distribuidores.
import { apiFetch } from "./http";

export interface DeudaDistribuidor {
  referidoId: number;
  nombre: string;
  pendiente: number;
}

export interface Liquidacion {
  id: number;
  referidoId: number;
  total: number;
  estado: string;
  creadoEn?: string;
  pagadaEn?: string;
  referidoNombre?: string;
}

// GET /api/admin/deudas (admin): comisiones pendientes agrupadas por distribuidor.
export const getDeudas = (): Promise<DeudaDistribuidor[]> =>
  apiFetch<DeudaDistribuidor[]>("/api/admin/deudas");

// POST /api/liquidaciones (admin): liquida todo el pendiente de un distribuidor.
export const crearLiquidacion = (referidoId: number): Promise<Liquidacion> =>
  apiFetch<Liquidacion>("/api/liquidaciones", {
    method: "POST",
    body: JSON.stringify({ referidoId }),
  });

// GET /api/liquidaciones (admin): historial global de pagos.
export const getLiquidaciones = (): Promise<Liquidacion[]> =>
  apiFetch<Liquidacion[]>("/api/liquidaciones");

// GET /api/liquidaciones/mias (distribuidor): sus propios cobros.
export const getMisLiquidaciones = (): Promise<Liquidacion[]> =>
  apiFetch<Liquidacion[]>("/api/liquidaciones/mias");
