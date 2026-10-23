// src/components/AdminLiquidaciones.tsx
// Back-office de liquidaciones (Fase 5): deuda por distribuidor + botón
// "Liquidar" (confirmación) y historial de pagos realizados.
import { useCallback, useEffect, useState } from "react";
import {
  getDeudas,
  getLiquidaciones,
  crearLiquidacion,
  type DeudaDistribuidor,
  type Liquidacion,
} from "../services/liquidacionesService";
import { ApiError } from "../services/http";

const AdminLiquidaciones = () => {
  const [deudas, setDeudas] = useState<DeudaDistribuidor[]>([]);
  const [historial, setHistorial] = useState<Liquidacion[]>([]);
  const [cargando, setCargando] = useState(true);
  const [error, setError] = useState("");
  const [liquidando, setLiquidando] = useState<number | null>(null);

  const cargar = useCallback(async () => {
    setCargando(true);
    setError("");
    try {
      const [d, h] = await Promise.all([getDeudas(), getLiquidaciones()]);
      setDeudas(d);
      setHistorial(h);
    } catch (err) {
      setError(err instanceof ApiError ? err.message : "Error al cargar liquidaciones");
    } finally {
      setCargando(false);
    }
  }, []);

  useEffect(() => {
    const t = setTimeout(() => void cargar(), 0);
    return () => clearTimeout(t);
  }, [cargar]);

  const liquidar = async (d: DeudaDistribuidor) => {
    const ok = window.confirm(
      `¿Liquidar $${d.pendiente.toFixed(2)} a ${d.nombre}? Esto marca sus comisiones como pagadas.`
    );
    if (!ok) return;
    setError("");
    setLiquidando(d.referidoId);
    try {
      await crearLiquidacion(d.referidoId);
      await cargar();
    } catch (err) {
      setError(err instanceof ApiError ? err.message : "No se pudo liquidar");
    } finally {
      setLiquidando(null);
    }
  };

  return (
    <div>
      <h1 className="text-2xl font-bold text-slate-800 mb-6">Liquidaciones</h1>

      {error && (
        <div className="bg-red-50 text-red-600 border border-red-200 p-3 rounded-lg text-sm mb-4">
          {error}
        </div>
      )}

      {cargando ? (
        <p className="text-slate-400">Cargando…</p>
      ) : (
        <>
          {/* Deuda por distribuidor */}
          <div className="bg-white border border-slate-200 rounded-lg shadow-sm overflow-x-auto mb-8">
            <div className="p-4 border-b border-slate-100">
              <h2 className="text-lg font-bold text-slate-800">
                Deuda pendiente por distribuidor
              </h2>
            </div>
            {deudas.length === 0 ? (
              <p className="p-6 text-slate-400 text-center">
                No hay comisiones pendientes. Todo al día.
              </p>
            ) : (
              <table className="w-full text-sm">
                <thead className="bg-slate-50 text-slate-600 uppercase text-xs">
                  <tr>
                    <th className="text-left px-4 py-3 font-semibold">Distribuidor</th>
                    <th className="text-left px-4 py-3 font-semibold">ID nodo</th>
                    <th className="text-right px-4 py-3 font-semibold">Pendiente</th>
                    <th className="text-right px-4 py-3 font-semibold">Acción</th>
                  </tr>
                </thead>
                <tbody>
                  {deudas.map((d) => (
                    <tr key={d.referidoId} className="border-t border-slate-100">
                      <td className="px-4 py-3 font-semibold text-slate-800">
                        {d.nombre}
                      </td>
                      <td className="px-4 py-3 text-slate-500">#{d.referidoId}</td>
                      <td className="px-4 py-3 text-right font-bold text-slate-800">
                        ${d.pendiente.toFixed(2)}
                      </td>
                      <td className="px-4 py-3 text-right">
                        <button
                          onClick={() => void liquidar(d)}
                          disabled={liquidando === d.referidoId}
                          className="bg-indigo-600 text-white font-semibold px-4 py-2 rounded-lg hover:bg-indigo-700 transition disabled:opacity-50 min-h-[44px]"
                        >
                          {liquidando === d.referidoId ? "Liquidando…" : "Liquidar"}
                        </button>
                      </td>
                    </tr>
                  ))}
                </tbody>
              </table>
            )}
          </div>

          {/* Historial de liquidaciones */}
          <div className="bg-white border border-slate-200 rounded-lg shadow-sm overflow-x-auto">
            <div className="p-4 border-b border-slate-100">
              <h2 className="text-lg font-bold text-slate-800">Historial de pagos</h2>
            </div>
            {historial.length === 0 ? (
              <p className="p-6 text-slate-400 text-center">
                Aún no se han realizado liquidaciones.
              </p>
            ) : (
              <table className="w-full text-sm">
                <thead className="bg-slate-50 text-slate-600 uppercase text-xs">
                  <tr>
                    <th className="text-left px-4 py-3 font-semibold">#</th>
                    <th className="text-left px-4 py-3 font-semibold">Distribuidor</th>
                    <th className="text-left px-4 py-3 font-semibold">Pagada</th>
                    <th className="text-right px-4 py-3 font-semibold">Total</th>
                    <th className="text-center px-4 py-3 font-semibold">Estado</th>
                  </tr>
                </thead>
                <tbody>
                  {historial.map((l) => (
                    <tr key={l.id} className="border-t border-slate-100">
                      <td className="px-4 py-3 text-slate-500">#{l.id}</td>
                      <td className="px-4 py-3 font-semibold text-slate-800">
                        {l.referidoNombre ?? `#${l.referidoId}`}
                      </td>
                      <td className="px-4 py-3 text-slate-500">
                        {l.pagadaEn ? new Date(l.pagadaEn).toLocaleString() : "—"}
                      </td>
                      <td className="px-4 py-3 text-right font-bold text-slate-800">
                        ${l.total.toFixed(2)}
                      </td>
                      <td className="px-4 py-3 text-center">
                        <span className="text-xs font-semibold px-2 py-1 rounded-full bg-green-100 text-green-700 capitalize">
                          {l.estado}
                        </span>
                      </td>
                    </tr>
                  ))}
                </tbody>
              </table>
            )}
          </div>
        </>
      )}
    </div>
  );
};

export default AdminLiquidaciones;
