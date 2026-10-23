// src/services/pagina.ts
// Sobre estándar de paginación que devuelve el backend (Fase 7).
export interface Pagina<T> {
  data: T[];
  page: number;
  limit: number;
  total: number;
}
