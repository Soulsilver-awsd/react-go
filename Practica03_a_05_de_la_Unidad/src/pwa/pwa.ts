// src/pwa/pwa.ts
// Tipos, contexto y hook del estado PWA (sin componentes).
// Se separa del Provider para cumplir la regla de Fast Refresh:
// un archivo que exporta componentes no debe exportar hooks/utilidades.
import { createContext, useContext } from "react";

export interface PwaContextType {
  /** El navegador permite instalar la app (evento beforeinstallprompt capturado). */
  installable: boolean;
  /** Lanza el diálogo de instalación del navegador. */
  install: () => Promise<void>;
  /** La app ya está instalada y se ejecuta en modo standalone. */
  standalone: boolean;
  /** Es iOS/iPadOS: allí no existe beforeinstallprompt y se instala desde Safari. */
  isIos: boolean;
  /** Hay una versión nueva del service worker esperando a activarse. */
  updateAvailable: boolean;
  /** Activa la versión nueva del service worker. */
  applyUpdate: () => Promise<void>;
  /** El shell ya quedó cacheado: la app puede abrirse sin conexión. */
  offlineReady: boolean;
  /** El navegador reporta conexión de red. */
  isOnline: boolean;
}

export const PwaContext = createContext<PwaContextType | undefined>(undefined);

// Hook con validación de tipo.
export const usePwa = () => {
  const context = useContext(PwaContext);
  if (!context) {
    throw new Error("usePwa debe ser usado dentro de un PwaProvider");
  }
  return context;
};
