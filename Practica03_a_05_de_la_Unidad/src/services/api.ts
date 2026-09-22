import type { Producto } from '../context/CartContext';

const API_URL = import.meta.env.VITE_API_URL ?? 'http://localhost:3000/api';

interface LoginResponse {
  token: string;
  email: string;
}

interface ApiError {
  error?: string;
}

const request = async <T>(path: string, options?: RequestInit): Promise<T> => {
  const response = await fetch(`${API_URL}${path}`, {
    headers: {
      'Content-Type': 'application/json',
      ...options?.headers,
    },
    ...options,
  });

  if (!response.ok) {
    const body = (await response.json().catch(() => ({}))) as ApiError;
    throw new Error(body.error ?? 'No se pudo completar la solicitud');
  }

  return response.json() as Promise<T>;
};

export const loginRequest = (email: string, password: string) =>
  request<LoginResponse>('/login', {
    method: 'POST',
    body: JSON.stringify({ email, password }),
  });

export const getProductos = () => request<Producto[]>('/productos');