import { useState } from "react";
import { Outlet } from "react-router-dom";
import Sidebar from "./Sidebar";
import Navbar from "./Navbar";

const Layout = () => {
  const [isCollapsed, setIsCollapsed] = useState(false);
  const [isMobileOpen, setIsMobileOpen] = useState(false);

  // En escritorio alterna el colapso del sidebar; en móvil abre/cierra el drawer
  const toggleSidebar = () => {
    if (window.innerWidth < 768) {
      setIsMobileOpen((open) => !open);
    } else {
      setIsCollapsed((collapsed) => !collapsed);
    }
  };

  const closeMobileMenu = () => setIsMobileOpen(false);

  return (
    <div className="flex min-h-screen bg-slate-50">
      <Sidebar
        isCollapsed={isCollapsed}
        isMobileOpen={isMobileOpen}
        closeMobileMenu={closeMobileMenu}
      />

      {/* Área de Contenido Principal */}
      <div className="flex min-w-0 flex-1 flex-col overflow-hidden">
        <Navbar toggleSidebar={toggleSidebar} />

        <main className="flex-1 overflow-x-hidden overflow-y-auto p-4 md:p-8">
          <Outlet />
        </main>
      </div>
    </div>
  );
};

export default Layout;