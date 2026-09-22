import { NavLink } from "react-router-dom";
import { LayoutDashboard, Network, Package } from "lucide-react";

interface SidebarProps {
  isCollapsed: boolean;
}

const menuItems = [
  { to: "/", label: "Dashboard", icon: LayoutDashboard },
  { to: "/catalogo", label: "Catálogo", icon: Package },
  { to: "/mi-red", label: "Mi Red", icon: Network },
];

const Sidebar = ({ isCollapsed }: SidebarProps) => {
  return (
    <aside
      className={`fixed inset-y-0 left-0 z-40 flex shrink-0 flex-col bg-slate-900 text-white transition-[width] duration-300 md:static md:z-auto ${
        isCollapsed ? "w-20" : "w-64"
      }`}
    >
      <div
        className={`flex h-16 items-center border-b border-slate-700 px-4 text-xl font-bold ${
          isCollapsed ? "justify-center" : "justify-start"
        }`}
      >
        {isCollapsed ? "M" : "MultiCatálogo"}
      </div>

      <nav className="flex-1 space-y-2 p-4">
        {menuItems.map(({ to, label, icon: Icon }) => (
          <NavLink
            key={to}
            to={to}
            end={to === "/"}
            title={isCollapsed ? label : undefined}
            className={({ isActive }) =>
              `flex items-center gap-3 rounded p-3 transition ${
                isCollapsed ? "justify-center" : "justify-start"
              } ${
                isActive
                  ? "bg-slate-700 text-white"
                  : "text-slate-300 hover:bg-slate-800 hover:text-white"
              }`
            }
          >
            <Icon size={20} aria-hidden="true" />
            {!isCollapsed && <span>{label}</span>}
          </NavLink>
        ))}
      </nav>
    </aside>
  );
};

export default Sidebar;