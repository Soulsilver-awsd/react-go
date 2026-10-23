// src/components/PwaInstallButton.tsx
// Botón "Instalar app" de la Navbar.
// - En Chrome/Edge el navegador emite 'beforeinstallprompt': lo lanzamos.
// - En iOS no existe ese evento: mostramos cómo instalarla desde Safari.
// - Si la app ya está instalada (standalone) no se muestra nada.
import { useState } from "react";
import { Download, Share, Smartphone, X } from "lucide-react";
import { usePwa } from "../pwa/pwa";

const PwaInstallButton = () => {
  const { installable, install, standalone, isIos } = usePwa();
  const [ayudaVisible, setAyudaVisible] = useState(false);

  // Ya instalada: no tiene sentido ofrecer instalarla otra vez.
  if (standalone) return null;

  // iOS: siempre ofrecemos la guía manual (no hay prompt programático).
  // Chromium: solo cuando el navegador confirma que se puede instalar.
  if (!isIos && !installable) return null;

  if (isIos && !installable) {
    return (
      <div className="relative">
        <button
          onClick={() => setAyudaVisible((visible) => !visible)}
          className="flex min-h-[44px] min-w-[44px] items-center gap-2 rounded-lg border border-indigo-200 bg-indigo-50 px-3 py-2 text-xs font-semibold text-indigo-700 transition hover:bg-indigo-100 focus:outline-none focus:ring-2 focus:ring-indigo-500 sm:text-sm"
          aria-expanded={ayudaVisible}
        >
          <Smartphone size={16} aria-hidden="true" />
          <span className="hidden sm:inline">Añadir a inicio</span>
        </button>
        {ayudaVisible && (
          <div className="absolute right-0 top-full z-50 mt-2 w-64 rounded-lg border border-slate-200 bg-white p-4 text-xs text-slate-600 shadow-lg">
            <div className="mb-2 flex items-start justify-between gap-2">
              <p className="font-semibold text-slate-700">
                Instalar en iPhone/iPad
              </p>
              <button
                onClick={() => setAyudaVisible(false)}
                className="rounded p-1 text-slate-400 hover:bg-slate-100 hover:text-slate-600"
                aria-label="Cerrar ayuda"
              >
                <X size={14} aria-hidden="true" />
              </button>
            </div>
            <ol className="list-decimal space-y-1 pl-4">
              <li>
                Pulsa <Share size={12} className="inline" aria-hidden="true" />{" "}
                (Compartir) en la barra de Safari.
              </li>
              <li>
                Elige <strong>«Añadir a pantalla de inicio»</strong>.
              </li>
              <li>Abre MultiCatálogo desde el icono de tu pantalla de inicio.</li>
            </ol>
          </div>
        )}
      </div>
    );
  }

  return (
    <button
      onClick={() => void install()}
      className="flex min-h-[44px] min-w-[44px] items-center gap-2 rounded-lg border border-indigo-200 bg-indigo-50 px-3 py-2 text-xs font-semibold text-indigo-700 transition hover:bg-indigo-100 focus:outline-none focus:ring-2 focus:ring-indigo-500 sm:text-sm"
      title="Instalar MultiCatálogo en este dispositivo"
    >
      <Download size={16} aria-hidden="true" />
      <span className="hidden sm:inline">Instalar app</span>
    </button>
  );
};

export default PwaInstallButton;
