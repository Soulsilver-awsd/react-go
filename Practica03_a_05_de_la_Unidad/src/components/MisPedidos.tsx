// src/components/MisPedidos.tsx
// Historial de pedidos del cliente (Fase 3): lista paginada con detalle
// expandible. Mobile-first: tarjetas apiladas en móvil, tabla en escritorio.
import { useCallback, useEffect, useState } from "react";
import {
  getMisPedidos,
  getPedidoById,
  type Pedido,
} from "../services/pedidosService";
import { ApiError } from "../services/http";
import EstadoBadge from "./EstadoBadge";

const LIMIT = 10;

const MisPedidos = () => {
  const [pedidos, setPedidos] = useState<Pedido[]>([]);
  const [page, setPage] = useState(1);
  const [total, setTotal] = useState(0);
  const [cargando, setCargando] = useState(true);
  const [error, setError] = useState("");

  // Pedido expandido (con su detalle cargado bajo demanda).
  const [detalle, setDetalle] = useState<Pedido | null>(null);
  const [cargandoDetalle, setCargandoDetalle] = useState(false);

  const cargar = useCallback(async () => {
    setCargando(true);
    setError("");
    try {
      const p = await getMisPedidos(page, LIMIT);
      setPedidos(p.data);
      setTotal(p.total);
    } catch (err) {
      setError(err instanceof ApiError ? err.message : "Error al cargar pedidos");
    } finally {
      setCargando(false);
    }
  }, [page]);

  useEffect(() => {
    const t = setTimeout(() => void cargar(), 0);
    return () => clearTimeout(t);
  }, [cargar]);

  const verDetalle = async (id: number) => {
    if (detalle?.id === id) return setDetalle(null); // toggle
    setCargandoDetalle(true);
    try {
      setDetalle(await getPedidoById(id));
    } catch (err) {
      setError(err instanceof ApiError ? err.message : "Error al cargar el detalle");
    } finally {
      setCargandoDetalle(false);
    }
  };

  const totalPaginas = Math.max(1, Math.ceil(total / LIMIT));

  return (
    <div>
      <h1 className="text-2xl font-bold text-slate-800 mb-6">Mis Pedidos</h1>

      {error && (
        <div className="bg-red-50 text-red-600 border border-red-200 p-3 rounded-lg text-sm mb-4">
          {error}
        </div>
      )}

      {cargando ? (
        <p className="text-slate-400">Cargando pedidos…</p>
      ) : pedidos.length === 0 ? (
        <div className="bg-white border border-slate-200 rounded-lg p-8 text-center text-slate-400">
          Aún no tienes pedidos.
        </div>
      ) : (
        <div className="space-y-4">
          {pedidos.map((p) => (
            <div
              key={p.id}
              className="bg-white border border-slate-200 rounded-lg shadow-sm"
            >
              <div className="p-4 flex flex-wrap items-center justify-between gap-3">
                <div>
                  <p className="font-semibold text-slate-800">{p.numero}</p>
                  <p className="text-xs text-slate-400">
                    {p.creadoEn
                      ? new Date(p.creadoEn).toLocaleString()
                      : "—"}{" "}
                    · #{p.id}
                  </p>
                </div>
                <div className="flex items-center gap-4">
                  <EstadoBadge estado={p.estado} />
                  <p className="font-bold text-slate-800">
                    ${p.total.toFixed(2)}
                  </p>
                  <button
                    onClick={() => void verDetalle(p.id)}
                    className="text-indigo-600 font-semibold hover:underline min-h-[44px]"
                  >
                    {detalle?.id === p.id ? "Ocultar" : "Ver detalle"}
                  </button>
                </div>
              </div>

              {detalle?.id === p.id && (
                <div className="border-t border-slate-100 p-4 bg-slate-50">
                  <p className="text-sm text-slate-500 mb-2">
                    {detalle.clienteNombre} · {detalle.ciudad} ·{" "}
                    {detalle.metodoPago}
                  </p>
                  <ul className="divide-y divide-slate-200">
                    {detalle.items?.map((it, i) => (
                      <li
                        key={i}
                        className="flex items-center justify-between py-2 text-sm"
                      >
                        <span className="text-slate-700">
                          {it.nombre} × {it.cantidad}
                        </span>
                        <span className="font-semibold text-slate-800">
                          ${(it.precio * it.cantidad).toFixed(2)}
                        </span>
                      </li>
                    ))}
                  </ul>
                </div>
              )}
            </div>
          ))}
        </div>
      )}

      {cargandoDetalle && (
        <p className="text-slate-400 text-sm mt-2">Cargando detalle…</p>
      )}

      {/* Paginación */}
      {totalPaginas > 1 && (
        <div className="flex items-center justify-center gap-4 mt-6">
          <button
            disabled={page <= 1}
            onClick={() => setPage((p) => p - 1)}
            className="px-4 py-2 rounded-lg bg-slate-200 text-slate-700 font-semibold disabled:opacity-40 min-h-[44px]"
          >
            Anterior
          </button>
          <span className="text-sm text-slate-500">
            {page} / {totalPaginas}
          </span>
          <button
            disabled={page >= totalPaginas}
            onClick={() => setPage((p) => p + 1)}
            className="px-4 py-2 rounded-lg bg-slate-200 text-slate-700 font-semibold disabled:opacity-40 min-h-[44px]"
          >
            Siguiente
          </button>
        </div>
      )}
    </div>
  );
};

export default MisPedidos;
