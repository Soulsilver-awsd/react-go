// src/components/AdminUsuarios.tsx
// Back-office de usuarios (Fase 7): tabla paginada de cuentas con cambio de
// rol (admin/cliente) y activar/desactivar. Solo rol admin (AdminRoute).
import { useCallback, useEffect, useState } from "react";
import {
  getUsuarios,
  actualizarUsuario,
  type Usuario,
  type RolUsuario,
} from "../services/usuariosService";
import { ApiError } from "../services/http";

const LIMIT = 10;

const AdminUsuarios = () => {
  const [usuarios, setUsuarios] = useState<Usuario[]>([]);
  const [cargando, setCargando] = useState(true);
  const [error, setError] = useState("");
  const [page, setPage] = useState(1);
  const [total, setTotal] = useState(0);
  const [ocupado, setOcupado] = useState<number | null>(null);
  const totalPaginas = Math.max(1, Math.ceil(total / LIMIT));

  const cargar = useCallback(async (p: number) => {
    setCargando(true);
    setError("");
    try {
      const pagina = await getUsuarios(p, LIMIT);
      setUsuarios(pagina.data);
      setTotal(pagina.total);
    } catch (err) {
      setError(err instanceof ApiError ? err.message : "Error al cargar usuarios");
    } finally {
      setCargando(false);
    }
  }, []);

  useEffect(() => {
    const t = setTimeout(() => void cargar(page), 0);
    return () => clearTimeout(t);
  }, [cargar, page]);

  const cambiarRol = async (u: Usuario, rol: RolUsuario) => {
    setError("");
    setOcupado(u.id);
    try {
      const actualizado = await actualizarUsuario(u.id, { rol });
      setUsuarios((prev) => prev.map((x) => (x.id === u.id ? actualizado : x)));
    } catch (err) {
      setError(err instanceof ApiError ? err.message : "No se pudo cambiar el rol");
    } finally {
      setOcupado(null);
    }
  };

  const alternarActivo = async (u: Usuario) => {
    const ok = window.confirm(
      `¿${u.activo ? "Desactivar" : "Activar"} la cuenta ${u.email}?`
    );
    if (!ok) return;
    setError("");
    setOcupado(u.id);
    try {
      const actualizado = await actualizarUsuario(u.id, { activo: !u.activo });
      setUsuarios((prev) => prev.map((x) => (x.id === u.id ? actualizado : x)));
    } catch (err) {
      setError(err instanceof ApiError ? err.message : "No se pudo actualizar");
    } finally {
      setOcupado(null);
    }
  };

  return (
    <div>
      <div className="mb-6">
        <h1 className="text-2xl font-bold text-slate-800">Usuarios</h1>
        <p className="text-slate-500 text-sm">Cuentas registradas y su rol</p>
      </div>

      {error && (
        <div className="bg-red-50 text-red-600 border border-red-200 p-3 rounded-lg text-sm mb-4">
          {error}
        </div>
      )}

      <div className="bg-white border border-slate-200 rounded-lg shadow-sm overflow-x-auto">
        <table className="w-full text-sm">
          <thead className="bg-slate-50 text-slate-600 uppercase text-xs">
            <tr>
              <th className="text-left px-4 py-3 font-semibold">Usuario</th>
              <th className="text-left px-4 py-3 font-semibold">Rol</th>
              <th className="text-center px-4 py-3 font-semibold">Estado</th>
              <th className="text-right px-4 py-3 font-semibold">Acciones</th>
            </tr>
          </thead>
          <tbody>
            {cargando ? (
              <tr>
                <td colSpan={4} className="px-4 py-8 text-center text-slate-400">
                  Cargando usuarios…
                </td>
              </tr>
            ) : usuarios.length === 0 ? (
              <tr>
                <td colSpan={4} className="px-4 py-8 text-center text-slate-400">
                  No hay usuarios.
                </td>
              </tr>
            ) : (
              usuarios.map((u) => (
                <tr key={u.id} className="border-t border-slate-100">
                  <td className="px-4 py-3">
                    <p className="font-semibold text-slate-800">
                      {u.nombre || "(sin nombre)"}
                    </p>
                    <p className="text-xs text-slate-400">
                      {u.email} · #{u.id}
                    </p>
                  </td>
                  <td className="px-4 py-3">
                    <select
                      value={u.rol}
                      disabled={ocupado === u.id}
                      onChange={(e) => void cambiarRol(u, e.target.value as RolUsuario)}
                      className="px-3 py-2 rounded-lg border border-slate-300 focus:ring-2 focus:ring-indigo-600 outline-none min-h-[44px]"
                    >
                      <option value="cliente">cliente</option>
                      <option value="admin">admin</option>
                    </select>
                  </td>
                  <td className="px-4 py-3 text-center">
                    <span
                      className={`text-xs font-semibold px-2 py-1 rounded-full ${
                        u.activo
                          ? "bg-green-100 text-green-700"
                          : "bg-red-100 text-red-600"
                      }`}
                    >
                      {u.activo ? "Activo" : "Inactivo"}
                    </span>
                  </td>
                  <td className="px-4 py-3 text-right whitespace-nowrap">
                    <button
                      onClick={() => void alternarActivo(u)}
                      disabled={ocupado === u.id}
                      className={`font-semibold hover:underline min-h-[44px] disabled:opacity-40 ${
                        u.activo ? "text-red-600" : "text-green-600"
                      }`}
                    >
                      {u.activo ? "Desactivar" : "Activar"}
                    </button>
                  </td>
                </tr>
              ))
            )}
          </tbody>
        </table>
      </div>

      <div className="flex items-center justify-between gap-3 mt-4">
        <p className="text-sm text-slate-500">
          {total} usuario{total === 1 ? "" : "s"} · página {page} de {totalPaginas}
        </p>
        <div className="flex gap-2">
          <button
            onClick={() => setPage((p) => Math.max(1, p - 1))}
            disabled={page <= 1 || cargando}
            className="bg-slate-200 text-slate-700 font-semibold px-4 py-2 rounded-lg hover:bg-slate-300 transition disabled:opacity-40 min-h-[44px]"
          >
            Anterior
          </button>
          <button
            onClick={() => setPage((p) => Math.min(totalPaginas, p + 1))}
            disabled={page >= totalPaginas || cargando}
            className="bg-slate-200 text-slate-700 font-semibold px-4 py-2 rounded-lg hover:bg-slate-300 transition disabled:opacity-40 min-h-[44px]"
          >
            Siguiente
          </button>
        </div>
      </div>
    </div>
  );
};

export default AdminUsuarios;
