// src/components/PwaUpdateToast.tsx
// Aviso de "Nueva versión disponible".
// El service worker se registró con registerType:'prompt', así que el worker
// anterior sigue activo hasta que el usuario acepta la actualización.
import { RefreshCw } from "lucide-react";
import { usePwa } from "../pwa/pwa";

const PwaUpdateToast = () => {
  const { updateAvailable, applyUpdate } = usePwa();

  if (!updateAvailable) return null;

  return (
    <div
      role="alert"
      className="fixed bottom-4 left-4 right-4 z-50 flex items-center justify-between gap-3 rounded-xl border border-indigo-200 bg-white p-4 shadow-xl sm:left-auto sm:right-6 sm:w-96"
    >
      <div>
        <p className="text-sm font-semibold text-slate-800">
          Nueva versión disponible
        </p>
        <p className="text-xs text-slate-500">
          Actualiza para cargar los últimos cambios de MultiCatálogo.
        </p>
      </div>
      <button
        onClick={() => void applyUpdate()}
        className="flex min-h-[44px] shrink-0 items-center gap-2 rounded-lg bg-indigo-600 px-4 py-2 text-sm font-semibold text-white transition hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-indigo-500"
      >
        <RefreshCw size={16} aria-hidden="true" />
        Actualizar
      </button>
    </div>
  );
};

export default PwaUpdateToast;
