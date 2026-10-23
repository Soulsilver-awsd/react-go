// src/components/MiRed.tsx
// Red multinivel real (Fase 4):
// - Subárbol del usuario autenticado (GET /api/red/mia), no el árbol global.
// - Link de referido personal (/registro?ref=CODIGO) con botón copiar.
// - Panel "Mis comisiones" con el ledger y totales por estado.
import { useCallback, useEffect, useState } from "react";
import {
  getMiRed,
  getMisComisiones,
  type NodoRed,
  type MisComisiones,
} from "../services/redService";
import { getMisLiquidaciones, type Liquidacion } from "../services/liquidacionesService";
import { ApiError } from "../services/http";

// Colores y etiquetas por nivel para el árbol.
const NIVEL_ESTILO: Record<number, { badge: string; border: string; tasa: string }> = {
  1: { badge: "bg-indigo-100 text-indigo-700", border: "border-indigo-200", tasa: "10 %" },
  2: { badge: "bg-purple-100 text-purple-700", border: "border-purple-200", tasa: "5 %" },
  3: { badge: "bg-pink-100 text-pink-700", border: "border-pink-200", tasa: "2 %" },
};

// Aplana el subárbol para contar nodos y sumar ventas/comisiones.
const aplanar = (nodo: NodoRed | null | undefined): NodoRed[] => {
  if (!nodo) return [];
  const hijos = (nodo.hijos ?? []).flatMap(aplanar);
  return [nodo, ...hijos];
};

// Componente recursivo del árbol (prop "nodo" porque "ref" está reservada).
const NodoReferido = ({ nodo }: { nodo: NodoRed }) => {
  const estilo = NIVEL_ESTILO[nodo.nivel] ?? NIVEL_ESTILO[1];
  return (
    <li>
      <div className={`bg-white rounded-lg border ${estilo.border} shadow-sm p-4 flex flex-wrap items-center justify-between gap-3`}>
        <div className="flex items-center gap-3">
          <div className="w-10 h-10 rounded-full bg-slate-200 flex items-center justify-center font-bold text-slate-600">
            {nodo.nombre.charAt(0)}
          </div>
          <div>
            <p className="font-semibold text-slate-800">{nodo.nombre}</p>
            <span className={`inline-block mt-1 text-xs font-semibold px-2 py-0.5 rounded-full ${estilo.badge}`}>
              Nivel {nodo.nivel} · {estilo.tasa} comisión
            </span>
          </div>
        </div>
        <div className="text-right">
          <p className="text-sm text-slate-500">Ventas</p>
          <p className="font-bold text-slate-800">${nodo.ventas.toLocaleString()}</p>
        </div>
      </div>

      {nodo.hijos && nodo.hijos.length > 0 && (
        <ul className="mt-3 ml-6 space-y-3 border-l-2 border-slate-200 pl-4">
          {nodo.hijos.map((hijo) => (
            <NodoReferido key={hijo.id} nodo={hijo} />
          ))}
        </ul>
      )}
    </li>
  );
};

const MiRed = () => {
  const [red, setRed] = useState<NodoRed | null>(null);
  const [comisiones, setComisiones] = useState<MisComisiones | null>(null);
  const [misLiquidaciones, setMisLiquidaciones] = useState<Liquidacion[]>([]);
  const [cargando, setCargando] = useState(true);
  const [error, setError] = useState("");
  const [copiado, setCopiado] = useState(false);

  const cargar = useCallback(async () => {
    setCargando(true);
    setError("");
    try {
      const [r, c, l] = await Promise.all([
        getMiRed(),
        getMisComisiones(),
        getMisLiquidaciones(),
      ]);
      setRed(r);
      setComisiones(c);
      setMisLiquidaciones(l);
    } catch (err) {
      setError(err instanceof ApiError ? err.message : "Error al cargar tu red");
    } finally {
      setCargando(false);
    }
  }, []);

  useEffect(() => {
    const t = setTimeout(() => void cargar(), 0);
    return () => clearTimeout(t);
  }, [cargar]);

  // Link de invitación absoluto con el código del nodo raíz.
  const link = red?.codigo
    ? `${window.location.origin}/registro?ref=${red.codigo}`
    : "";

  const copiar = async () => {
    if (!link) return;
    try {
      await navigator.clipboard.writeText(link);
      setCopiado(true);
      setTimeout(() => setCopiado(false), 2000);
    } catch {
      setError("No se pudo copiar el enlace");
    }
  };

  if (cargando) return <p className="text-slate-400">Cargando tu red…</p>;

  // Nodos del subárbol sin contar la raíz (los referidos reales).
  const nodos = aplanar(red).filter((n) => n.id !== red?.id);
  const ventasRed = nodos.reduce((acc, n) => acc + n.ventas, 0);

  const resumen = [
    { etiqueta: "Referidos Activos", valor: nodos.length.toString(), color: "text-indigo-600" },
    { etiqueta: "Ventas de la Red", valor: `$${ventasRed.toLocaleString()}`, color: "text-indigo-600" },
    { etiqueta: "Comisiones Pendientes", valor: `$${(comisiones?.pendiente ?? 0).toFixed(2)}`, color: "text-green-600" },
    { etiqueta: "Comisiones Pagadas", valor: `$${(comisiones?.pagado ?? 0).toFixed(2)}`, color: "text-slate-600" },
  ];

  return (
    <div>
      <h1 className="text-2xl font-bold text-slate-800 mb-2">Mi Red de Referidos</h1>
      <p className="text-slate-500 mb-6">
        Gana el 10 % de nivel 1, 5 % de nivel 2 y 2 % de nivel 3 sobre las compras de tu red.
      </p>

      {error && (
        <div className="bg-red-50 text-red-600 border border-red-200 p-3 rounded-lg text-sm mb-4">
          {error}
        </div>
      )}

      {/* Link de referido */}
      {link && (
        <div className="bg-white border border-slate-200 rounded-lg shadow-sm p-4 mb-6">
          <p className="text-sm font-semibold text-slate-700 mb-2">
            Tu enlace de invitación
          </p>
          <div className="flex flex-wrap gap-2">
            <input
              readOnly
              value={link}
              className="flex-1 min-w-0 px-3 py-3 rounded-lg border border-slate-300 bg-slate-50 text-slate-600 text-sm"
            />
            <button
              onClick={() => void copiar()}
              className="bg-indigo-600 text-white font-semibold px-4 py-3 rounded-lg hover:bg-indigo-700 transition min-h-[44px]"
            >
              {copiado ? "¡Copiado!" : "Copiar"}
            </button>
          </div>
        </div>
      )}

      {/* Resumen */}
      <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-6 mb-8">
        {resumen.map((item) => (
          <div key={item.etiqueta} className="bg-white p-6 rounded-lg shadow-sm border border-slate-100">
            <p className="text-sm text-slate-500 uppercase font-semibold">{item.etiqueta}</p>
            <p className={`text-3xl font-bold mt-2 ${item.color}`}>{item.valor}</p>
          </div>
        ))}
      </div>

      {/* Árbol de referidos */}
      <div className="bg-white rounded-lg border border-slate-200 shadow-sm p-6 mb-8">
        <h2 className="text-lg font-bold text-slate-800 mb-4">Jerarquía de tu Red</h2>

        {red && (
          <>
            <div className="bg-slate-900 text-white rounded-lg p-4 mb-4 flex flex-wrap items-center justify-between gap-3">
              <div className="flex items-center gap-3">
                <div className="w-10 h-10 rounded-full bg-indigo-500 flex items-center justify-center font-bold">
                  {red.nombre.charAt(0)}
                </div>
                <p className="font-semibold">{red.nombre} (tú)</p>
              </div>
              <p className="text-xs text-slate-300">Código: {red.codigo ?? "—"}</p>
            </div>

            {red.hijos && red.hijos.length > 0 ? (
              <ul className="space-y-3">
                {red.hijos.map((hijo) => (
                  <NodoReferido key={hijo.id} nodo={hijo} />
                ))}
              </ul>
            ) : (
              <p className="text-slate-400 text-sm">
                Aún no tienes referidos. ¡Comparte tu enlace para empezar!
              </p>
            )}
          </>
        )}
      </div>

      {/* Panel de comisiones */}
      <div className="bg-white rounded-lg border border-slate-200 shadow-sm p-6">
        <h2 className="text-lg font-bold text-slate-800 mb-4">Mis comisiones</h2>
        {!comisiones || comisiones.comisiones.length === 0 ? (
          <p className="text-slate-400 text-sm">Todavía no has generado comisiones.</p>
        ) : (
          <div className="overflow-x-auto">
            <table className="w-full text-sm">
              <thead className="bg-slate-50 text-slate-600 uppercase text-xs">
                <tr>
                  <th className="text-left px-4 py-3 font-semibold">Pedido</th>
                  <th className="text-left px-4 py-3 font-semibold">Nivel</th>
                  <th className="text-right px-4 py-3 font-semibold">%</th>
                  <th className="text-right px-4 py-3 font-semibold">Monto</th>
                  <th className="text-center px-4 py-3 font-semibold">Estado</th>
                </tr>
              </thead>
              <tbody>
                {comisiones.comisiones.map((c) => (
                  <tr key={c.id} className="border-t border-slate-100">
                    <td className="px-4 py-3 text-slate-600">#{c.pedidoId}</td>
                    <td className="px-4 py-3 text-slate-600">Nivel {c.nivel}</td>
                    <td className="px-4 py-3 text-right text-slate-600">
                      {c.porcentaje}%
                    </td>
                    <td className="px-4 py-3 text-right font-semibold text-slate-800">
                      ${c.monto.toFixed(2)}
                    </td>
                    <td className="px-4 py-3 text-center">
                      <span
                        className={`text-xs font-semibold px-2 py-1 rounded-full capitalize ${
                          c.estado === "pendiente"
                            ? "bg-amber-100 text-amber-700"
                            : c.estado === "pagada"
                            ? "bg-green-100 text-green-700"
                            : "bg-slate-100 text-slate-500"
                        }`}
                      >
                        {c.estado}
                      </span>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        )}
      </div>
      {/* Panel de liquidaciones recibidas */}
      <div className="bg-white rounded-lg border border-slate-200 shadow-sm p-6 mt-8">
        <h2 className="text-lg font-bold text-slate-800 mb-4">Liquidaciones recibidas</h2>
        {misLiquidaciones.length === 0 ? (
          <p className="text-slate-400 text-sm">Aún no has recibido liquidaciones.</p>
        ) : (
          <div className="overflow-x-auto">
            <table className="w-full text-sm">
              <thead className="bg-slate-50 text-slate-600 uppercase text-xs">
                <tr>
                  <th className="text-left px-4 py-3 font-semibold">#</th>
                  <th className="text-left px-4 py-3 font-semibold">Pagada</th>
                  <th className="text-right px-4 py-3 font-semibold">Total</th>
                  <th className="text-center px-4 py-3 font-semibold">Estado</th>
                </tr>
              </thead>
              <tbody>
                {misLiquidaciones.map((l) => (
                  <tr key={l.id} className="border-t border-slate-100">
                    <td className="px-4 py-3 text-slate-500">#{l.id}</td>
                    <td className="px-4 py-3 text-slate-500">
                      {l.pagadaEn ? new Date(l.pagadaEn).toLocaleString() : "—"}
                    </td>
                    <td className="px-4 py-3 text-right font-semibold text-slate-800">
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
          </div>
        )}
      </div>
    </div>
  );
};

export default MiRed;
