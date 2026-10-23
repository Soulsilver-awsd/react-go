// src/components/AdminPedidos.tsx
// Back-office de pedidos (Fase 3): filtro por estado + cambio de estado.
// El admin solo puede aplicar transiciones válidas (el backend las valida igual).
import { useCallback, useEffect, useState } from "react";
import {
  getPedidosAdmin,
  cambiarEstadoPedido,
  getPedidoById,
  type Pedido,
  type EstadoPedido,
} from "../services/pedidosService";
import { ApiError } from "../services/http";
import EstadoBadge from "./EstadoBadge";

const LIMIT = 10;

// Estados ofrecidos en el filtro.
const FILTROS: (EstadoPedido | "")[] = [
  "",
  "pendiente",
  "preparando",
  "enviado",
  "entregado",
  "cancelado",
];

// Transiciones permitidas (espejo de la máquina de estados del backend).
const SIGUIENTES: Record<EstadoPedido, EstadoPedido[]> = {
  pendiente: ["preparando", "cancelado"],
  preparando: ["enviado", "cancelado"],
  enviado: ["entregado", "cancelado"],
  entregado: [],
  cancelado: [],
};

const AdminPedidos = () => {
  const [pedidos, setPedidos] = useState<Pedido[]>([]);
  const [page, setPage] = useState(1);
  const [total, setTotal] = useState(0);
  const [estado, setEstado] = useState<EstadoPedido | "">("");
  const [cargando, setCargando] = useState(true);
  const [error, setError] = useState("");
  const [detalle, setDetalle] = useState<Pedido | null>(null);

  const cargar = useCallback(async () => {
    setCargando(true);
    setError("");
    try {
      const p = await getPedidosAdmin(estado, page, LIMIT);
      setPedidos(p.data);
      setTotal(p.total);
    } catch (err) {
      setError(err instanceof ApiError ? err.message : "Error al cargar pedidos");
    } finally {
      setCargando(false);
    }
  }, [estado, page]);

  useEffect(() => {
    const t = setTimeout(() => void cargar(), 0);
    return () => clearTimeout(t);
  }, [cargar]);

  const cambiar = async (id: number, nuevo: EstadoPedido) => {
    setError("");
    try {
      await cambiarEstadoPedido(id, nuevo);
      await cargar();
    } catch (err) {
      setError(err instanceof ApiError ? err.message : "No se pudo cambiar el estado");
    }
  };

  const verDetalle = async (id: number) => {
    if (detalle?.id === id) return setDetalle(null);
    try {
      setDetalle(await getPedidoById(id));
    } catch (err) {
      setError(err instanceof ApiError ? err.message : "Error al cargar detalle");
    }
  };

  const totalPaginas = Math.max(1, Math.ceil(total / LIMIT));

  return (
    <div>
      <h1 className="text-2xl font-bold text-slate-800 mb-6">Pedidos</h1>

      <div className="mb-4">
        <label className="text-sm font-medium text-slate-700 mr-2">
          Filtrar por estado:
        </label>
        <select
          value={estado}
          onChange={(e) => {
            setEstado(e.target.value as EstadoPedido | "");
            setPage(1);
          }}
          className="px-3 py-2 rounded-lg border border-slate-300 focus:ring-2 focus:ring-indigo-600 outline-none min-h-[44px]"
        >
          {FILTROS.map((f) => (
            <option key={f} value={f}>
              {f === "" ? "Todos" : f}
            </option>
          ))}
        </select>
      </div>

      {error && (
        <div className="bg-red-50 text-red-600 border border-red-200 p-3 rounded-lg text-sm mb-4">
          {error}
        </div>
      )}

      {cargando ? (
        <p className="text-slate-400">Cargando pedidos…</p>
      ) : pedidos.length === 0 ? (
        <div className="bg-white border border-slate-200 rounded-lg p-8 text-center text-slate-400">
          No hay pedidos con ese filtro.
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
                    {p.clienteEmail} · {p.ciudad} · #{p.id}
                  </p>
                </div>
                <div className="flex flex-wrap items-center gap-3">
                  <EstadoBadge estado={p.estado} />
                  <p className="font-bold text-slate-800">
                    ${p.total.toFixed(2)}
                  </p>
                  <button
                    onClick={() => void verDetalle(p.id)}
                    className="text-indigo-600 font-semibold hover:underline min-h-[44px]"
                  >
                    {detalle?.id === p.id ? "Ocultar" : "Detalle"}
                  </button>
                  {SIGUIENTES[p.estado].map((sig) => (
                    <button
                      key={sig}
                      onClick={() => void cambiar(p.id, sig)}
                      className={`font-semibold px-3 py-2 rounded-lg text-sm min-h-[44px] ${
                        sig === "cancelado"
                          ? "bg-red-100 text-red-600 hover:bg-red-200"
                          : "bg-indigo-600 text-white hover:bg-indigo-700"
                      }`}
                    >
                      {sig === "cancelado" ? "Cancelar" : `Pasar a ${sig}`}
                    </button>
                  ))}
                </div>
              </div>

              {detalle?.id === p.id && (
                <div className="border-t border-slate-100 p-4 bg-slate-50">
                  <p className="text-sm text-slate-500 mb-2">
                    {detalle.clienteNombre} · {detalle.metodoPago}
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

export default AdminPedidos;
