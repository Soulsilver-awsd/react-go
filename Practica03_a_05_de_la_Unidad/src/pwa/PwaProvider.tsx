// src/pwa/PwaProvider.tsx
// Provider del estado PWA (U1-T6):
// - registra el service worker (vite-plugin-pwa, modo 'prompt'),
// - captura el evento de instalación del navegador,
// - vigila el estado de la conexión de red.
import { useCallback, useEffect, useMemo, useState, type ReactNode } from "react";
import { useRegisterSW } from "virtual:pwa-register/react";
import { PwaContext, type PwaContextType } from "./pwa";

interface PwaProviderProps {
  children: ReactNode;
}

// Evento que emite Chromium cuando la app ya puede instalarse.
interface BeforeInstallPromptEvent extends Event {
  prompt: () => Promise<void>;
  userChoice: Promise<{ outcome: "accepted" | "dismissed" }>;
}

// Detecta iOS/iPadOS: allí no existe beforeinstallprompt y la instalación
// se hace desde el menú "Compartir" de Safari.
const detectarIos = (): boolean => {
  if (typeof navigator === "undefined") return false;
  const agente = navigator.userAgent;
  const iPadOsComoMac = /Macintosh/.test(agente) && "ontouchend" in document;
  return /iPad|iPhone|iPod/.test(agente) || iPadOsComoMac;
};

// Detecta si la app ya se ejecuta instalada (modo standalone).
const detectarStandalone = (): boolean => {
  if (typeof window === "undefined") return false;
  const safariIos = (navigator as Navigator & { standalone?: boolean }).standalone;
  return (
    window.matchMedia("(display-mode: standalone)").matches || safariIos === true
  );
};

export const PwaProvider = ({ children }: PwaProviderProps) => {
  const {
    needRefresh: [needRefresh, setNeedRefresh],
    offlineReady: [offlineReady, setOfflineReady],
    updateServiceWorker,
  } = useRegisterSW({
    onNeedRefresh: () => setNeedRefresh(true),
    onOfflineReady: () => setOfflineReady(true),
  });

  // Prompt de instalación pendiente (lo lanza el botón de la Navbar).
  const [deferredPrompt, setDeferredPrompt] =
    useState<BeforeInstallPromptEvent | null>(null);
  const [isOnline, setIsOnline] = useState<boolean>(() =>
    typeof navigator === "undefined" ? true : navigator.onLine
  );
  const [standalone, setStandalone] = useState<boolean>(detectarStandalone);
  const [isIos] = useState<boolean>(detectarIos);

  // El navegador avisa cuándo la app puede instalarse.
  useEffect(() => {
    const onBeforeInstallPrompt = (event: Event) => {
      // Sin preventDefault el navegador mostraría su propio aviso.
      event.preventDefault();
      setDeferredPrompt(event as BeforeInstallPromptEvent);
    };
    const onAppInstalled = () => {
      setDeferredPrompt(null);
      setStandalone(true);
    };
    window.addEventListener("beforeinstallprompt", onBeforeInstallPrompt);
    window.addEventListener("appinstalled", onAppInstalled);
    return () => {
      window.removeEventListener("beforeinstallprompt", onBeforeInstallPrompt);
      window.removeEventListener("appinstalled", onAppInstalled);
    };
  }, []);

  // Estado de la conexión de red.
  useEffect(() => {
    const onOnline = () => setIsOnline(true);
    const onOffline = () => setIsOnline(false);
    window.addEventListener("online", onOnline);
    window.addEventListener("offline", onOffline);
    return () => {
      window.removeEventListener("online", onOnline);
      window.removeEventListener("offline", onOffline);
    };
  }, []);

  // Chequeo explícito de actualizaciones: el navegador solo revalida el
  // service worker al navegar y con un throttle de ~24 h, por lo que una
  // versión nueva desplegada mientras la app está abierta puede no avisar.
  // Forzamos la comprobación al arrancar y al volver a la pestaña.
  useEffect(() => {
    if (!("serviceWorker" in navigator)) return;
    let cancelado = false;

    const comprobar = async () => {
      try {
        const registro = await navigator.serviceWorker.getRegistration();
        if (registro && !cancelado) await registro.update();
      } catch {
        // Sin SW registrado (p. ej. en dev) no hay nada que actualizar.
      }
    };

    void comprobar();
    const onVisible = () => {
      if (document.visibilityState === "visible") void comprobar();
    };
    document.addEventListener("visibilitychange", onVisible);
    return () => {
      cancelado = true;
      document.removeEventListener("visibilitychange", onVisible);
    };
  }, []);

  const install = useCallback(async () => {
    if (!deferredPrompt) return;
    await deferredPrompt.prompt();
    const eleccion = await deferredPrompt.userChoice;
    if (eleccion.outcome === "accepted") {
      setDeferredPrompt(null);
    }
  }, [deferredPrompt]);

  const applyUpdate = useCallback(async () => {
    // Envía SKIP_WAITING al service worker en espera; el plugin recarga la
    // página cuando el nuevo worker toma el control.
    await updateServiceWorker(true);
  }, [updateServiceWorker]);

  const value = useMemo<PwaContextType>(
    () => ({
      installable: deferredPrompt !== null,
      install,
      standalone,
      isIos,
      updateAvailable: needRefresh,
      applyUpdate,
      offlineReady,
      isOnline,
    }),
    [
      deferredPrompt,
      install,
      standalone,
      isIos,
      needRefresh,
      applyUpdate,
      offlineReady,
      isOnline,
    ]
  );

  return <PwaContext.Provider value={value}>{children}</PwaContext.Provider>;
};
