// src/components/Registro.tsx
// Vista pública de creación de cuenta (U1-T5 · Fase 1).
// Mobile-first: un campo por fila, botones ≥44px.
import { useState, type FormEvent } from "react";
import { Link, useNavigate, useSearchParams } from "react-router-dom";
import { useAuth, type Rol } from "../context/auth";
import { apiFetch, ApiError } from "../services/http";

// Respuesta de POST /api/register
interface RegisterResponse {
  id: number;
  email: string;
  rol: string;
  nombre: string;
}

// Respuesta de POST /api/login
interface LoginResponse {
  token: string;
  email: string;
  rol: string;
}

const Registro = () => {
  // El código de referido llega por query string: /registro?ref=LUIS
  const [searchParams] = useSearchParams();
  const refInicial = searchParams.get("ref") ?? "";

  const [nombre, setNombre] = useState("");
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [ref, setRef] = useState(refInicial);
  const [error, setError] = useState("");
  const [loading, setLoading] = useState(false);

  const navigate = useNavigate();
  const { login } = useAuth();

  const handleSubmit = async (e: FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    setLoading(true);
    setError("");

    try {
      // 1) Creamos la cuenta (con código de referido opcional).
      await apiFetch<RegisterResponse>("/api/register", {
        method: "POST",
        body: JSON.stringify({ nombre, email, password, ref }),
      });

      // 2) Login automático tras el registro, para entrar directo.
      const sesion = await apiFetch<LoginResponse>("/api/login", {
        method: "POST",
        body: JSON.stringify({ email, password }),
      });
      const rol: Rol = sesion.rol === "admin" ? "admin" : "cliente";
      login({ email: sesion.email, rol }, sesion.token);

      navigate(rol === "admin" ? "/" : "/tienda");
    } catch (err) {
      setError(
        err instanceof ApiError ? err.message : "No se pudo completar el registro"
      );
      setLoading(false);
    }
  };

  return (
    <div className="min-h-screen flex items-center justify-center bg-slate-100 px-4 py-8">
      <div className="max-w-md w-full bg-white rounded-xl shadow-md p-8 border border-slate-200">
        <div className="text-center mb-8">
          <h2 className="text-3xl font-bold text-slate-900">Crear cuenta</h2>
          <p className="text-slate-500 mt-2">
            Únete a la red MultiCatálogo
          </p>
        </div>

        {error && (
          <div className="bg-red-50 text-red-600 p-3 rounded-lg text-sm mb-6 text-center border border-red-200">
            {error}
          </div>
        )}

        <form onSubmit={handleSubmit} className="space-y-6">
          <div>
            <label className="block text-sm font-medium text-slate-700 mb-2">
              Nombre completo
            </label>
            <input
              type="text"
              value={nombre}
              onChange={(e) => setNombre(e.target.value)}
              className="w-full px-4 py-3 rounded-lg border border-slate-300 focus:ring-2 focus:ring-indigo-600 focus:border-indigo-600 outline-none transition"
              placeholder="Ana García"
              required
            />
          </div>

          <div>
            <label className="block text-sm font-medium text-slate-700 mb-2">
              Correo Electrónico
            </label>
            <input
              type="email"
              value={email}
              onChange={(e) => setEmail(e.target.value)}
              className="w-full px-4 py-3 rounded-lg border border-slate-300 focus:ring-2 focus:ring-indigo-600 focus:border-indigo-600 outline-none transition"
              placeholder="tu@correo.com"
              required
            />
          </div>

          <div>
            <label className="block text-sm font-medium text-slate-700 mb-2">
              Contraseña (mín. 6 caracteres)
            </label>
            <input
              type="password"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              className="w-full px-4 py-3 rounded-lg border border-slate-300 focus:ring-2 focus:ring-indigo-600 focus:border-indigo-600 outline-none transition"
              placeholder="••••••"
              minLength={6}
              required
            />
          </div>

          <div>
            <label className="block text-sm font-medium text-slate-700 mb-2">
              Código de referido (opcional)
            </label>
            <input
              type="text"
              value={ref}
              onChange={(e) => setRef(e.target.value.toUpperCase())}
              className="w-full px-4 py-3 rounded-lg border border-slate-300 focus:ring-2 focus:ring-indigo-600 focus:border-indigo-600 outline-none transition uppercase"
              placeholder="LUIS"
            />
          </div>

          <button
            type="submit"
            disabled={loading}
            className="w-full bg-indigo-600 text-white font-bold py-3 rounded-lg hover:bg-indigo-700 transition disabled:opacity-50 disabled:cursor-not-allowed min-h-[44px]"
          >
            {loading ? "Creando cuenta..." : "Crear cuenta"}
          </button>
        </form>

        <p className="text-center text-sm text-slate-500 mt-6">
          ¿Ya tienes cuenta?{" "}
          <Link to="/login" className="text-indigo-600 font-semibold hover:underline">
            Inicia sesión
          </Link>
        </p>
      </div>
    </div>
  );
};

export default Registro;
