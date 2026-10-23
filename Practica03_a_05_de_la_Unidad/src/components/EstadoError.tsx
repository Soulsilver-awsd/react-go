// src/components/EstadoError.tsx
// Estado de error reutilizable (U1-T6 · PWA): mensaje claro + reintento.
// Evita que una vista sin conexión muestre un "catálogo vacío" sin explicación.
import { CloudOff, RefreshCw } from "lucide-react";

interface EstadoErrorProps {
  mensaje: string;
  onReintentar?: () => void;
}

const EstadoError = ({ mensaje, onReintentar }: EstadoErrorProps) => {
  return (
    <div className="rounded-lg border border-slate-200 bg-white p-10 text-center">
      <CloudOff
        size={32}
        className="mx-auto mb-4 text-slate-400"
        aria-hidden="true"
      />
      <p className="text-lg text-slate-700">{mensaje}</p>
      {onReintentar && (
        <button
          onClick={onReintentar}
          className="mt-4 inline-flex min-h-[44px] items-center gap-2 rounded-lg bg-slate-900 px-4 py-2 text-sm font-semibold text-white transition hover:bg-indigo-600 focus:outline-none focus:ring-2 focus:ring-indigo-500"
        >
          <RefreshCw size={16} aria-hidden="true" />
          Reintentar
        </button>
      )}
    </div>
  );
};

export default EstadoError;
