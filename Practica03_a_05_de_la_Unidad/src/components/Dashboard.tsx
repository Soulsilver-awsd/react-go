// src/components/Dashboard.tsx
// Resumen general con KPIs reales consumidos de /api/reportes/* (Fase 6):
// total vendido, nº de pedidos, ticket promedio y comisiones devengadas,
// más el ranking de referidos y productos.
import { useCallback, useEffect, useState } from "react";
import {
  getResumen,
  getTopProductos,
  getTopReferidos,
  type ResumenReporte,
  type TopProducto,
  type TopReferido,
} from "../services/reportesService";
import { ApiError } from "../services/http";

const Dashboard = () => {
  const [resumen, setResumen] = useState<ResumenReporte | null>(null);
  const [topReferidos, setTopReferidos] = useState<TopReferido[]>([]);
  const [topProductos, setTopProductos] = useState<TopProducto[]>([]);
  const [cargando, setCargando] = useState(true);
  const [error, setError] = useState("");

  const cargar = useCallback(async () => {
    setCargando(true);
    setError("");
    try {
      const [r, refs, prods] = await Promise.all([
        getResumen(),
        getTopReferidos(),
        getTopProductos(),
      ]);
      setResumen(r);
      setTopReferidos(refs);
      setTopProductos(prods);
    } catch (err) {
      setError(err instanceof ApiError ? err.message : "Error al cargar el dashboard");
    } finally {
      setCargando(false);
    }
  }, []);

  useEffect(() => {
    const t = setTimeout(() => void cargar(), 0);
    return () => clearTimeout(t);
  }, [cargar]);

  if (cargando) {
    return (
      <div>
        <h1 className="text-2xl font-bold text-slate-800 mb-6">Resumen General</h1>
        <p className="text-slate-400">Cargando indicadores…</p>
      </div>
    );
  }

  if (error) {
    return (
      <div>
        <h1 className="text-2xl font-bold text-slate-800 mb-6">Resumen General</h1>
        <div className="bg-red-50 text-red-600 border border-red-200 p-3 rounded-lg text-sm">
          {error}
        </div>
      </div>
    );
  }

  const kpis = [
    {
      etiqueta: "Total Vendido",
      valor: `$${(resumen?.totalVendido ?? 0).toLocaleString()}`,
      color: "text-indigo-600",
    },
    {
      etiqueta: "Pedidos",
      valor: (resumen?.numPedidos ?? 0).toString(),
      color: "text-indigo-600",
    },
    {
      etiqueta: "Ticket Promedio",
      valor: `$${(resumen?.ticketPromedio ?? 0).toFixed(2)}`,
      color: "text-amber-600",
    },
    {
      etiqueta: "Comisiones Devengadas",
      valor: `$${(resumen?.comisionesDevengadas ?? 0).toFixed(2)}`,
      color: "text-green-600",
    },
  ];

  return (
    <div>
      <h1 className="text-2xl font-bold text-slate-800 mb-6">Resumen General</h1>

      {/* ===== KPIs REALES ===== */}
      <div className="grid grid-cols-1 md:grid-cols-2 xl:grid-cols-4 gap-6 mb-8">
        {kpis.map((kpi) => (
          <div
            key={kpi.etiqueta}
            className="bg-white p-6 rounded-lg shadow-sm border border-slate-100"
          >
            <p className="text-sm text-slate-500 uppercase font-semibold">{kpi.etiqueta}</p>
            <p className={`text-3xl font-bold mt-2 ${kpi.color}`}>{kpi.valor}</p>
          </div>
        ))}
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
        {/* ===== TOP REFERIDOS ===== */}
        <div className="bg-white p-6 rounded-lg shadow-sm border border-slate-100">
          <h2 className="text-lg font-bold text-slate-800 mb-4">Top Referidos de la Red</h2>
          {topReferidos.length === 0 ? (
            <p className="text-slate-400 text-sm">Aún no hay datos de la red.</p>
          ) : (
            <div className="space-y-3">
              {topReferidos.slice(0, 5).map((ref, idx) => (
                <div
                  key={ref.referidoId}
                  className="flex items-center justify-between p-3 rounded-lg border border-slate-100 hover:bg-slate-50 transition"
                >
                  <div className="flex items-center gap-3">
                    <span className="w-8 h-8 rounded-full bg-indigo-100 text-indigo-700 font-bold flex items-center justify-center">
                      {idx + 1}
                    </span>
                    <div>
                      <p className="font-semibold text-slate-800">{ref.nombre}</p>
                      <p className="text-xs text-slate-500">
                        {ref.codigo ? `Código ${ref.codigo}` : `Nodo #${ref.referidoId}`}
                      </p>
                    </div>
                  </div>
                  <div className="text-right">
                    <p className="font-bold text-slate-800">
                      ${ref.ventas.toLocaleString()}
                    </p>
                    <p className="text-xs text-green-600 font-semibold">
                      Saldo ${ref.saldo.toFixed(2)}
                    </p>
                  </div>
                </div>
              ))}
            </div>
          )}
        </div>

        {/* ===== TOP PRODUCTOS ===== */}
        <div className="bg-white p-6 rounded-lg shadow-sm border border-slate-100">
          <h2 className="text-lg font-bold text-slate-800 mb-4">Productos Más Vendidos</h2>
          {topProductos.length === 0 ? (
            <p className="text-slate-400 text-sm">Aún no hay ventas registradas.</p>
          ) : (
            <div className="space-y-3">
              {topProductos.slice(0, 5).map((p, idx) => (
                <div
                  key={p.productoId}
                  className="flex items-center justify-between p-3 rounded-lg border border-slate-100 hover:bg-slate-50 transition"
                >
                  <div className="flex items-center gap-3">
                    <span className="w-8 h-8 rounded-full bg-amber-100 text-amber-700 font-bold flex items-center justify-center">
                      {idx + 1}
                    </span>
                    <div>
                      <p className="font-semibold text-slate-800">{p.nombre}</p>
                      <p className="text-xs text-slate-500">{p.cantidad} unidades vendidas</p>
                    </div>
                  </div>
                  <p className="font-bold text-slate-800">${p.total.toLocaleString()}</p>
                </div>
              ))}
            </div>
          )}
        </div>
      </div>
    </div>
  );
};

export default Dashboard;
