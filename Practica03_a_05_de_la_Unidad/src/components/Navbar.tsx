import { Link, useNavigate } from "react-router-dom";
import { Menu, PanelLeftClose, ShoppingCart } from "lucide-react";
import { useCart } from "../context/CartContext";
import { useAuth } from "../context/AuthContext";

interface NavbarProps {
  isCollapsed: boolean;
  onToggleSidebar: () => void;
}

const Navbar = ({ isCollapsed, onToggleSidebar }: NavbarProps) => {
  const { totalItems } = useCart();
  const { logout, userEmail } = useAuth();
  const navigate = useNavigate();

  const handleLogout = () => {
    logout();
    navigate("/login");
  };

  return (
    <header className="flex min-h-16 items-center justify-between border-b border-slate-200 bg-white px-4 md:px-8">
      <div className="flex min-w-0 items-center gap-3">
        <button
          type="button"
          onClick={onToggleSidebar}
          aria-label={isCollapsed ? "Expandir menú" : "Colapsar menú"}
          title={isCollapsed ? "Expandir menú" : "Colapsar menú"}
          className="rounded p-2 text-slate-600 transition hover:bg-slate-100 hover:text-slate-900"
        >
          {isCollapsed ? <Menu size={22} /> : <PanelLeftClose size={22} />}
        </button>
        <h2 className="truncate text-base font-medium text-slate-600 md:text-lg">
          Panel de Administración
        </h2>
      </div>
      <div className="flex items-center gap-2 md:gap-6">
        <Link
          to="/carrito"
          aria-label="Abrir carrito"
          title="Abrir carrito"
          className="relative rounded-full p-2 text-slate-600 transition hover:bg-slate-100 hover:text-slate-900"
        >
          <ShoppingCart size={21} />
          {totalItems > 0 && (
            <span className="absolute top-0 right-0 bg-indigo-600 text-white text-xs font-bold w-5 h-5 flex items-center justify-center rounded-full transform translate-x-1 -translate-y-1">
              {totalItems}
            </span>
          )}
        </Link>

        <div className="flex items-center gap-2 md:gap-4">
          <span className="hidden max-w-40 truncate text-sm text-slate-500 sm:inline">
            {userEmail}
          </span>

          {/* Contenedor relativo con la clase 'group' para detectar el hover */}
          <div className="relative group cursor-pointer pb-2">
            {/* Círculo del usuario / Avatar */}
            <div className="w-8 h-8 rounded-full bg-slate-200 overflow-hidden border border-slate-300 flex items-center justify-center">
              {/* 
        NOTA PARA LA API: 
        Aquí reemplazarás el 'src' quemado por la variable de tu estado, 
        por ejemplo: src={userAvatar || defaultImage} 
      */}
              <img
                src="https://fastly.picsum.photos/id/64/4326/2884.jpg?hmac=9_SzX666YRpR_fOyYStXpfSiJ_edO3ghlSRnH2w09Kg"
                alt="Avatar del usuario"
                className="w-full h-full object-cover"
              />
            </div>

            {/* Menú desplegable */}
            <div className="absolute right-0 top-full mt-1 w-36 bg-white border border-slate-200 rounded-lg shadow-lg opacity-0 invisible group-hover:opacity-100 group-hover:visible transition-all duration-200 z-50">
              <button
                onClick={handleLogout}
                className="w-full text-left px-4 py-2 text-sm text-red-600 font-semibold hover:bg-red-50 rounded-md transition-colors"
              >
                Cerrar Sesión
              </button>
            </div>
          </div>
        </div>
      </div>
    </header>
  );
};

export default Navbar;
