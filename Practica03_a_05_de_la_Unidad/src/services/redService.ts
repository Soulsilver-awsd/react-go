// src/services/redService.ts
// Capa de servicios de la red MLM y el ledger de comisiones (Fase 4).
import { apiFetch } from "./http";

// Nodo del árbol de referidos (subárbol del usuario autenticado).
export interface NodoRed {
  id: number;
  nombre: string;
  nivel: number;
  ventas: number;
  saldo: number;
  codigo?: string;
  hijos?: NodoRed[];
}

export interface Comision {
  id: number;
  pedidoId: number;
  referidoId: number;
  nivel: number;
  porcentaje: number;
  monto: number;
  estado: "pendiente" | "pagada" | "revertida";
  creadoEn?: string;
}

export interface MisComisiones {
  comisiones: Comision[];
  pendiente: number;
  pagado: number;
}

// GET /api/red/mia (subárbol del usuario autenticado, no el global)
export const getMiRed = (): Promise<NodoRed> =>
  apiFetch<NodoRed>("/api/red/mia");

// GET /api/comisiones/mias (ledger + totales por estado)
export const getMisComisiones = (): Promise<MisComisiones> =>
  apiFetch<MisComisiones>("/api/comisiones/mias");
