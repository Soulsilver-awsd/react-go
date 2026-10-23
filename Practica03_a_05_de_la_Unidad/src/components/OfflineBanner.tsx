// src/components/OfflineBanner.tsx
// Banner superior que avisa de que no hay conexión de red.
// Aparece en el Layout, encima de la Navbar.
import { WifiOff } from "lucide-react";
import { usePwa } from "../pwa/pwa";

const OfflineBanner = () => {
  const { isOnline } = usePwa();

  // Con red no hay nada que avisar.
  if (isOnline) return null;

  return (
    <div
      role="status"
      className="flex items-center justify-center gap-2 bg-amber-500 px-4 py-2 text-center text-xs font-semibold text-white sm:text-sm"
    >
      <WifiOff size={16} aria-hidden="true" />
      <span>Sin conexión — se muestran los datos guardados en el dispositivo</span>
    </div>
  );
};

export default OfflineBanner;
