// src/services/http.ts
// Cliente HTTP único de la SPA (U1-T4).
// - Centraliza la URL base (VITE_API_URL) y el header Authorization: Bearer.
// - Tipa los errores del backend ({status,message,details}) como ApiError.
// - Ofrece helpers de sesión (token en localStorage) para no duplicar claves.

// URL base de la API; se configura por entorno y cae a localhost:3000 en dev.
export const API_URL: string =
  import.meta.env.VITE_API_URL ?? "http://localhost:3000";

// Clave de localStorage donde vive la sesión {usuario, token}.
const CLAVE_SESION = "multicatalogo_sesion";

// ApiError representa un error HTTP del backend con su status y detalles.
export class ApiError extends Error {
  status: number;
  details?: unknown;
  // true cuando el fallo es de red (sin conexión o servidor inaccesible),
  // no del backend. Lo usa la UI para ofrecer reintentos y avisos.
  offline: boolean;

  constructor(status: number, message: string, details?: unknown, offline = false) {
    super(message);
    this.name = "ApiError";
    this.status = status;
    this.details = details;
    this.offline = offline;
  }
}

// SesionPersistida es la forma exacta del objeto guardado en localStorage.
export interface SesionPersistida {
  email: string;
  rol: "admin" | "cliente";
  token: string;
}

// getToken devuelve el JWT de la sesión actual (o null si no hay sesión).
export const getToken = (): string | null => {
  try {
    const raw = localStorage.getItem(CLAVE_SESION);
    if (!raw) return null;
    return (JSON.parse(raw) as SesionPersistida).token ?? null;
  } catch {
    return null;
  }
};

// getSesion devuelve la sesión completa persistida, o null.
export const getSesion = (): SesionPersistida | null => {
  try {
    const raw = localStorage.getItem(CLAVE_SESION);
    return raw ? (JSON.parse(raw) as SesionPersistida) : null;
  } catch {
    return null;
  }
};

// guardarSesion persiste usuario + token en localStorage.
export const guardarSesion = (sesion: SesionPersistida): void => {
  localStorage.setItem(CLAVE_SESION, JSON.stringify(sesion));
};

// limpiarSesion elimina la sesión persistida (logout).
export const limpiarSesion = (): void => {
  localStorage.removeItem(CLAVE_SESION);
};

// apiFetch es el cliente único: arma la URL, adjunta el Bearer si hay sesión y
// normaliza los errores del backend ({status,message,details}) a ApiError.
export async function apiFetch<T>(
  path: string,
  options: RequestInit = {}
): Promise<T> {
  const token = getToken();

  const headers = new Headers(options.headers);
  if (!headers.has("Content-Type") && options.body) {
    headers.set("Content-Type", "application/json");
  }
  if (token) {
    headers.set("Authorization", `Bearer ${token}`);
  }

  let response: Response;
  try {
    response = await fetch(`${API_URL}${path}`, { ...options, headers });
  } catch {
    throw new ApiError(
      0,
      "Sin conexión: no se pudo contactar con el servidor",
      undefined,
      true
    );
  }

  // 204 u otras respuestas sin cuerpo.
  if (response.status === 204) {
    return undefined as T;
  }

  const data = await response.json().catch(() => null);

  if (!response.ok) {
    const mensaje =
      (data && (data.message as string)) || "Error del servidor";
    throw new ApiError(response.status, mensaje, data?.details);
  }

  return data as T;
}
