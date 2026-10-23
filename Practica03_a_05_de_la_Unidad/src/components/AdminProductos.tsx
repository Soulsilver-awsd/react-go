// src/components/AdminProductos.tsx
// Back-office de productos (Fase 2): listado admin + crear/editar/desactivar.
// Mobile-first: buscador arriba, tabla con scroll horizontal y formulario
// desplegable sobre la lista. Solo accesible para rol admin (AdminRoute).
import { useCallback, useEffect, useState, type FormEvent } from "react";
import {
  getProductosAdmin,
  crearProducto,
  actualizarProducto,
  desactivarProducto,
  type ProductoPayload,
} from "../services/productosService";
import { ApiError } from "../services/http";
import type { Producto } from "../data/productos";

// Estado vacío del formulario de alta/edición.
const FORM_VACIO: ProductoPayload = {
  nombre: "",
  descripcion: "",
  precio: 0,
  categoria: "",
  img: "",
  galeria: [],
  stock: 0,
  activo: true,
};

// El formulario maneja la galería como texto con comas (UI simple); la API
// recibe un array de URLs.
const galeriaATexto = (g: string[]) => g.join(", ");
const textoAGaleria = (t: string): string[] =>
  t
    .split(",")
    .map((s) => s.trim())
    .filter(Boolean);

// Elementos por página del listado admin.
const LIMIT = 10;

const AdminProductos = () => {
  const [productos, setProductos] = useState<Producto[]>([]);
  const [cargando, setCargando] = useState(true);
  const [error, setError] = useState("");

  const [busqueda, setBusqueda] = useState("");

  // Paginación estándar del backend (Fase 7).
  const [page, setPage] = useState(1);
  const [total, setTotal] = useState(0);
  const totalPaginas = Math.max(1, Math.ceil(total / LIMIT));

  // Formulario visible (null = oculto). Si editandoId no es null, es edición.
  const [mostrandoForm, setMostrandoForm] = useState(false);
  const [editandoId, setEditandoId] = useState<number | null>(null);
  const [form, setForm] = useState<ProductoPayload & { galeriaTexto: string }>({
    ...FORM_VACIO,
    galeriaTexto: "",
  });
  const [guardando, setGuardando] = useState(false);

  const cargar = useCallback(async (q: string, p: number) => {
    setCargando(true);
    setError("");
    try {
      const pagina = await getProductosAdmin(q, p, LIMIT);
      setProductos(pagina.data);
      setTotal(pagina.total);
    } catch (err) {
      setError(err instanceof ApiError ? err.message : "Error al cargar productos");
    } finally {
      setCargando(false);
    }
  }, []);

  // Búsqueda server-side con debounce de 400 ms (?q= ILIKE). Al cambiar el
  // término se vuelve a la primera página.
  useEffect(() => {
    const t = setTimeout(() => {
      setPage(1);
      void cargar(busqueda, 1);
    }, 400);
    return () => clearTimeout(t);
  }, [busqueda, cargar]);

  const cambiarPagina = (p: number) => {
    if (p < 1 || p > totalPaginas) return;
    setPage(p);
    void cargar(busqueda, p);
  };

  const abrirNuevo = () => {
    setEditandoId(null);
    setForm({ ...FORM_VACIO, galeriaTexto: "" });
    setMostrandoForm(true);
    setError("");
  };

  const abrirEdicion = (p: Producto) => {
    setEditandoId(p.id);
    setForm({
      nombre: p.nombre,
      descripcion: p.descripcion,
      precio: p.precio,
      categoria: p.categoria,
      img: p.img,
      galeria: p.galeria ?? [],
      galeriaTexto: galeriaATexto(p.galeria ?? []),
      stock: p.stock ?? 0,
      activo: p.activo ?? true,
    });
    setMostrandoForm(true);
    setError("");
  };

  const guardar = async (e: FormEvent<HTMLFormElement>) => {
    e.preventDefault();

    // Validación móvil-first antes de tocar la red.
    if (!form.nombre.trim()) return setError("El nombre es obligatorio");
    if (!form.categoria.trim()) return setError("La categoría es obligatoria");
    if (!(form.precio > 0)) return setError("El precio debe ser mayor que cero");
    if (form.stock < 0) return setError("El stock no puede ser negativo");

    const payload: ProductoPayload = {
      nombre: form.nombre.trim(),
      descripcion: form.descripcion,
      precio: Number(form.precio),
      categoria: form.categoria.trim(),
      img: form.img,
      galeria: textoAGaleria(form.galeriaTexto),
      stock: Number(form.stock),
      activo: form.activo,
    };

    setGuardando(true);
    setError("");
    try {
      if (editandoId === null) {
        await crearProducto(payload);
      } else {
        await actualizarProducto(editandoId, payload);
      }
      setMostrandoForm(false);
      setEditandoId(null);
      await cargar(busqueda, page); // refresca la lista (y la caché del catálogo)
    } catch (err) {
      setError(err instanceof ApiError ? err.message : "No se pudo guardar");
    } finally {
      setGuardando(false);
    }
  };

  // Borrado lógico con confirmación nativa (impide accidentes en móvil).
  const desactivar = async (p: Producto) => {
    const ok = window.confirm(
      `¿Desactivar "${p.nombre}"? Dejará de aparecer en el catálogo (borrado lógico).`
    );
    if (!ok) return;
    setError("");
    try {
      await desactivarProducto(p.id);
      await cargar(busqueda, page);
    } catch (err) {
      setError(err instanceof ApiError ? err.message : "No se pudo desactivar");
    }
  };

  return (
    <div>
      <div className="flex flex-wrap items-center justify-between gap-3 mb-6">
        <div>
          <h1 className="text-2xl font-bold text-slate-800">Productos</h1>
          <p className="text-slate-500 text-sm">
            Catálogo del back-office (incluye inactivos)
          </p>
        </div>
        <button
          onClick={abrirNuevo}
          className="bg-indigo-600 text-white font-semibold px-4 py-3 rounded-lg hover:bg-indigo-700 transition min-h-[44px]"
        >
          Nuevo producto
        </button>
      </div>

      {/* Buscador server-side */}
      <div className="mb-4">
        <input
          type="search"
          value={busqueda}
          onChange={(e) => setBusqueda(e.target.value)}
          placeholder="Buscar por nombre o descripción…"
          className="w-full max-w-md px-4 py-3 rounded-lg border border-slate-300 focus:ring-2 focus:ring-indigo-600 focus:border-indigo-600 outline-none"
        />
      </div>

      {error && (
        <div className="bg-red-50 text-red-600 border border-red-200 p-3 rounded-lg text-sm mb-4">
          {error}
        </div>
      )}

      {/* Formulario desplegable crear/editar */}
      {mostrandoForm && (
        <form
          onSubmit={guardar}
          className="bg-white border border-slate-200 rounded-lg shadow-sm p-6 mb-6 space-y-4"
        >
          <h2 className="text-lg font-bold text-slate-800">
            {editandoId === null ? "Nuevo producto" : `Editar producto #${editandoId}`}
          </h2>

          <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
            <div>
              <label className="block text-sm font-medium text-slate-700 mb-1">
                Nombre
              </label>
              <input
                value={form.nombre}
                onChange={(e) => setForm({ ...form, nombre: e.target.value })}
                className="w-full px-3 py-3 rounded-lg border border-slate-300 focus:ring-2 focus:ring-indigo-600 outline-none"
                required
              />
            </div>
            <div>
              <label className="block text-sm font-medium text-slate-700 mb-1">
                Categoría
              </label>
              <input
                value={form.categoria}
                onChange={(e) => setForm({ ...form, categoria: e.target.value })}
                className="w-full px-3 py-3 rounded-lg border border-slate-300 focus:ring-2 focus:ring-indigo-600 outline-none"
                required
              />
            </div>
            <div>
              <label className="block text-sm font-medium text-slate-700 mb-1">
                Precio ($)
              </label>
              <input
                type="number"
                min="0"
                step="0.01"
                value={form.precio}
                onChange={(e) => setForm({ ...form, precio: Number(e.target.value) })}
                className="w-full px-3 py-3 rounded-lg border border-slate-300 focus:ring-2 focus:ring-indigo-600 outline-none"
                required
              />
            </div>
            <div>
              <label className="block text-sm font-medium text-slate-700 mb-1">
                Stock
              </label>
              <input
                type="number"
                min="0"
                value={form.stock}
                onChange={(e) => setForm({ ...form, stock: Number(e.target.value) })}
                className="w-full px-3 py-3 rounded-lg border border-slate-300 focus:ring-2 focus:ring-indigo-600 outline-none"
              />
            </div>
          </div>

          <div>
            <label className="block text-sm font-medium text-slate-700 mb-1">
              Descripción
            </label>
            <textarea
              value={form.descripcion}
              onChange={(e) => setForm({ ...form, descripcion: e.target.value })}
              rows={3}
              className="w-full px-3 py-3 rounded-lg border border-slate-300 focus:ring-2 focus:ring-indigo-600 outline-none"
            />
          </div>

          <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
            <div>
              <label className="block text-sm font-medium text-slate-700 mb-1">
                Imagen principal (URL)
              </label>
              <input
                value={form.img}
                onChange={(e) => setForm({ ...form, img: e.target.value })}
                placeholder="https://…"
                className="w-full px-3 py-3 rounded-lg border border-slate-300 focus:ring-2 focus:ring-indigo-600 outline-none"
              />
            </div>
            <div>
              <label className="block text-sm font-medium text-slate-700 mb-1">
                Galería (URLs separadas por coma)
              </label>
              <input
                value={form.galeriaTexto}
                onChange={(e) => setForm({ ...form, galeriaTexto: e.target.value })}
                placeholder="https://a.jpg, https://b.jpg"
                className="w-full px-3 py-3 rounded-lg border border-slate-300 focus:ring-2 focus:ring-indigo-600 outline-none"
              />
            </div>
          </div>

          <label className="flex items-center gap-2 text-sm font-medium text-slate-700">
            <input
              type="checkbox"
              checked={form.activo}
              onChange={(e) => setForm({ ...form, activo: e.target.checked })}
              className="w-5 h-5"
            />
            Activo en el catálogo
          </label>

          <div className="flex gap-3">
            <button
              type="submit"
              disabled={guardando}
              className="bg-indigo-600 text-white font-semibold px-5 py-3 rounded-lg hover:bg-indigo-700 transition disabled:opacity-50 min-h-[44px]"
            >
              {guardando ? "Guardando…" : editandoId === null ? "Crear" : "Guardar"}
            </button>
            <button
              type="button"
              onClick={() => setMostrandoForm(false)}
              className="bg-slate-200 text-slate-700 font-semibold px-5 py-3 rounded-lg hover:bg-slate-300 transition min-h-[44px]"
            >
              Cancelar
            </button>
          </div>
        </form>
      )}

      {/* Listado */}
      <div className="bg-white border border-slate-200 rounded-lg shadow-sm overflow-x-auto">
        <table className="w-full text-sm">
          <thead className="bg-slate-50 text-slate-600 uppercase text-xs">
            <tr>
              <th className="text-left px-4 py-3 font-semibold">Producto</th>
              <th className="text-left px-4 py-3 font-semibold">Categoría</th>
              <th className="text-right px-4 py-3 font-semibold">Precio</th>
              <th className="text-right px-4 py-3 font-semibold">Stock</th>
              <th className="text-center px-4 py-3 font-semibold">Estado</th>
              <th className="text-right px-4 py-3 font-semibold">Acciones</th>
            </tr>
          </thead>
          <tbody>
            {cargando ? (
              <tr>
                <td colSpan={6} className="px-4 py-8 text-center text-slate-400">
                  Cargando productos…
                </td>
              </tr>
            ) : productos.length === 0 ? (
              <tr>
                <td colSpan={6} className="px-4 py-8 text-center text-slate-400">
                  No hay productos que coincidan.
                </td>
              </tr>
            ) : (
              productos.map((p) => (
                <tr key={p.id} className="border-t border-slate-100">
                  <td className="px-4 py-3">
                    <div className="flex items-center gap-3">
                      {p.img ? (
                        <img
                          src={p.img}
                          alt={p.nombre}
                          className="w-10 h-10 rounded object-cover bg-slate-100"
                        />
                      ) : (
                        <div className="w-10 h-10 rounded bg-slate-200 flex items-center justify-center text-slate-400 text-xs font-semibold">
                          N/A
                        </div>
                      )}
                      <div>
                        <p className="font-semibold text-slate-800">{p.nombre}</p>
                        <p className="text-xs text-slate-400">#{p.id}</p>
                      </div>
                    </div>
                  </td>
                  <td className="px-4 py-3 text-slate-600">{p.categoria}</td>
                  <td className="px-4 py-3 text-right font-semibold text-slate-800">
                    ${p.precio}
                  </td>
                  <td className="px-4 py-3 text-right text-slate-600">
                    {p.stock ?? 0}
                  </td>
                  <td className="px-4 py-3 text-center">
                    <span
                      className={`text-xs font-semibold px-2 py-1 rounded-full ${
                        p.activo !== false
                          ? "bg-green-100 text-green-700"
                          : "bg-red-100 text-red-600"
                      }`}
                    >
                      {p.activo !== false ? "Activo" : "Inactivo"}
                    </span>
                  </td>
                  <td className="px-4 py-3 text-right whitespace-nowrap">
                    <button
                      onClick={() => abrirEdicion(p)}
                      className="text-indigo-600 font-semibold hover:underline mr-3 min-h-[44px]"
                    >
                      Editar
                    </button>
                    {p.activo !== false && (
                      <button
                        onClick={() => void desactivar(p)}
                        className="text-red-600 font-semibold hover:underline min-h-[44px]"
                      >
                        Desactivar
                      </button>
                    )}
                  </td>
                </tr>
              ))
            )}
          </tbody>
        </table>
      </div>

      {/* Paginación estándar */}
      <div className="flex items-center justify-between gap-3 mt-4">
        <p className="text-sm text-slate-500">
          {total} resultado{total === 1 ? "" : "s"} · página {page} de {totalPaginas}
        </p>
        <div className="flex gap-2">
          <button
            onClick={() => cambiarPagina(page - 1)}
            disabled={page <= 1 || cargando}
            className="bg-slate-200 text-slate-700 font-semibold px-4 py-2 rounded-lg hover:bg-slate-300 transition disabled:opacity-40 min-h-[44px]"
          >
            Anterior
          </button>
          <button
            onClick={() => cambiarPagina(page + 1)}
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

export default AdminProductos;
