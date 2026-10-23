// src/components/PwaInstallButton.tsx
// Botón accesible de instalación de la PWA (Progressive Web App).
// - Si el navegador emite 'beforeinstallprompt', dispara el diálogo nativo de instalación.
// - Si es iOS o un navegador sin prompt automático (como Firefox o Chrome desktop sin banner activo),
//   despliega una ventana modal accesible con las instrucciones paso a paso para instalar o añadir a pantalla de inicio.
// - Si ya está en modo standalone, muestra estado de "App instalada" accesible.
import { useState, useEffect } from "react";
import { Download, Share, Smartphone, X, Monitor, CheckCircle2, Info } from "lucide-react";
import { usePwa } from "../pwa/pwa";

interface PwaInstallButtonProps {
  variant?: "navbar" | "sidebar";
  isCollapsed?: boolean;
}

const PwaInstallButton = ({
  variant = "navbar",
  isCollapsed = false,
}: PwaInstallButtonProps) => {
  const { installable, install, standalone, isIos } = usePwa();
  const [modalAbierto, setModalAbierto] = useState(false);

  // Cerrar modal al presionar la tecla Escape
  useEffect(() => {
    const handleKeyDown = (e: KeyboardEvent) => {
      if (e.key === "Escape" && modalAbierto) {
        setModalAbierto(false);
      }
    };
    window.addEventListener("keydown", handleKeyDown);
    return () => window.removeEventListener("keydown", handleKeyDown);
  }, [modalAbierto]);

  const handleClick = () => {
    if (standalone) {
      setModalAbierto(true);
      return;
    }

    if (installable) {
      void install();
    } else {
      setModalAbierto(true);
    }
  };

  // Botón para la Navbar
  const renderNavbarButton = () => (
    <button
      onClick={handleClick}
      className={`flex min-h-[40px] items-center gap-2 rounded-lg border px-3 py-1.5 text-xs font-semibold transition focus:outline-none focus:ring-2 focus:ring-indigo-500 sm:text-sm ${
        standalone
          ? "border-emerald-200 bg-emerald-50 text-emerald-700 hover:bg-emerald-100"
          : "border-indigo-200 bg-indigo-50 text-indigo-700 hover:bg-indigo-100"
      }`}
      aria-label={standalone ? "Aplicación Paulo instalada" : "Instalar aplicación Paulo (PWA)"}
      aria-haspopup="dialog"
      aria-expanded={modalAbierto}
      title={standalone ? "Aplicación ya instalada" : "Instalar aplicación en este dispositivo"}
    >
      {standalone ? (
        <CheckCircle2 size={16} aria-hidden="true" className="text-emerald-600" />
      ) : (
        <Download size={16} aria-hidden="true" className="text-indigo-600" />
      )}
      <span className="hidden sm:inline">
        {standalone ? "App Instalada" : "Instalar App"}
      </span>
      <span className="inline sm:hidden">
        {standalone ? "Instalada" : "Instalar"}
      </span>
    </button>
  );

  // Botón para la Sidebar
  const renderSidebarButton = () => (
    <button
      onClick={handleClick}
      className={`flex items-center gap-3 rounded p-2.5 transition focus:outline-none focus:ring-2 focus:ring-indigo-400 w-full ${
        standalone
          ? "bg-emerald-950/40 text-emerald-300 hover:bg-emerald-900/40"
          : "bg-indigo-950/40 text-indigo-200 hover:bg-indigo-900/50 hover:text-white"
      } ${isCollapsed ? "md:justify-center" : "justify-start"}`}
      aria-label={standalone ? "Aplicación Paulo instalada" : "Instalar aplicación Paulo (PWA)"}
      aria-haspopup="dialog"
      aria-expanded={modalAbierto}
      title={
        isCollapsed
          ? standalone
            ? "App Instalada"
            : "Instalar App (PWA)"
          : undefined
      }
    >
      {standalone ? (
        <CheckCircle2 size={20} aria-hidden="true" className="text-emerald-400 flex-shrink-0" />
      ) : (
        <Download size={20} aria-hidden="true" className="text-indigo-300 flex-shrink-0" />
      )}
      <span
        className={`whitespace-nowrap text-xs font-medium ${
          isCollapsed ? "md:hidden" : ""
        }`}
      >
        {standalone ? "PWA Instalada ✓" : "Instalar App (PWA)"}
      </span>
    </button>
  );

  return (
    <>
      {variant === "sidebar" ? renderSidebarButton() : renderNavbarButton()}

      {/* Modal accesible con instrucciones */}
      {modalAbierto && (
        <div
          className="fixed inset-0 z-50 flex items-center justify-center p-4 bg-slate-900/60 backdrop-blur-xs animate-in fade-in"
          role="dialog"
          aria-modal="true"
          aria-labelledby="pwa-modal-title"
          onClick={() => setModalAbierto(false)}
        >
          <div
            className="relative w-full max-w-md rounded-2xl bg-white p-6 shadow-2xl text-slate-800"
            onClick={(e) => e.stopPropagation()}
          >
            {/* Cabecera */}
            <div className="flex items-start justify-between pb-4 border-b border-slate-100">
              <div className="flex items-center gap-3">
                <div className="flex h-10 w-10 items-center justify-center rounded-xl bg-indigo-100 text-indigo-600">
                  <Smartphone className="h-5 w-5" aria-hidden="true" />
                </div>
                <div>
                  <h3
                    id="pwa-modal-title"
                    className="text-lg font-bold text-slate-900"
                  >
                    {standalone ? "Aplicación Instalada" : "Instalar Paulo (PWA)"}
                  </h3>
                  <p className="text-xs text-slate-500">
                    {standalone
                      ? "Estás usando la versión web progresiva instalada"
                      : "Acceso rápido y funcionamiento sin conexión"}
                  </p>
                </div>
              </div>
              <button
                onClick={() => setModalAbierto(false)}
                className="rounded-lg p-1.5 text-slate-400 hover:bg-slate-100 hover:text-slate-600 focus:outline-none focus:ring-2 focus:ring-indigo-500"
                aria-label="Cerrar ventana de ayuda"
              >
                <X size={18} aria-hidden="true" />
              </button>
            </div>

            {/* Contenido / Instrucciones según plataforma */}
            <div className="mt-4 space-y-4 text-sm text-slate-600">
              {standalone ? (
                <div className="flex items-start gap-3 rounded-xl bg-emerald-50 p-4 text-emerald-800">
                  <CheckCircle2 className="h-5 w-5 text-emerald-600 flex-shrink-0 mt-0.5" />
                  <div>
                    <p className="font-semibold text-emerald-900">¡App activa en modo autónomo!</p>
                    <p className="text-xs text-emerald-700 mt-1">
                      Ya tienes instalada la Progressive Web App en tu dispositivo con soporte de caché y acceso sin conexión.
                    </p>
                  </div>
                </div>
              ) : isIos ? (
                <div>
                  <p className="font-semibold text-slate-800 mb-2 flex items-center gap-2">
                    <Share size={16} className="text-indigo-600" /> Cómo instalar en iPhone o iPad (Safari):
                  </p>
                  <ol className="list-decimal space-y-2 pl-5 text-xs text-slate-600 leading-relaxed">
                    <li>
                      Toca el botón <strong>Compartir</strong> (<Share size={13} className="inline text-slate-700" />) en la barra inferior de Safari.
                    </li>
                    <li>
                      Desliza hacia abajo en el menú y selecciona <strong>«Añadir a pantalla de inicio»</strong>.
                    </li>
                    <li>
                      Toca <strong>«Añadir»</strong> en la esquina superior derecha.
                    </li>
                    <li>
                      ¡Listo! Abre <strong>Paulo</strong> directamente desde tu pantalla de inicio.
                    </li>
                  </ol>
                </div>
              ) : (
                <div className="space-y-3">
                  <div className="rounded-xl bg-slate-50 p-3.5 border border-slate-100">
                    <p className="font-semibold text-slate-800 text-xs flex items-center gap-2 mb-1.5">
                      <Monitor size={15} className="text-indigo-600" /> En Google Chrome o Microsoft Edge:
                    </p>
                    <p className="text-xs text-slate-600 leading-relaxed">
                      Haz clic en el ícono de instalación que aparece a la derecha de la barra de direcciones URL (icono de computadora o flecha hacia abajo <Download size={12} className="inline text-indigo-600" />), o ve al menú <strong>(⋮) &gt; «Instalar Paulo...»</strong>.
                    </p>
                  </div>

                  <div className="rounded-xl bg-slate-50 p-3.5 border border-slate-100">
                    <p className="font-semibold text-slate-800 text-xs flex items-center gap-2 mb-1.5">
                      <Smartphone size={15} className="text-indigo-600" /> En dispositivos Android:
                    </p>
                    <p className="text-xs text-slate-600 leading-relaxed">
                      Abre el menú de Chrome <strong>(⋮)</strong> en la esquina superior derecha y selecciona <strong>«Instalar aplicación»</strong> o <strong>«Añadir a la pantalla principal»</strong>.
                    </p>
                  </div>

                  <div className="flex items-center gap-2 text-xs text-slate-500 pt-1">
                    <Info size={14} className="text-indigo-500 flex-shrink-0" />
                    <span>Beneficios: Funciona sin conexión a internet y carga ultra rápida.</span>
                  </div>
                </div>
              )}
            </div>

            {/* Botón de acción */}
            <div className="mt-6 flex justify-end">
              <button
                onClick={() => setModalAbierto(false)}
                className="w-full sm:w-auto px-4 py-2 text-xs font-semibold rounded-lg bg-indigo-600 text-white hover:bg-indigo-700 transition focus:outline-none focus:ring-2 focus:ring-indigo-500"
              >
                Entendido
              </button>
            </div>
          </div>
        </div>
      )}
    </>
  );
};

export default PwaInstallButton;
