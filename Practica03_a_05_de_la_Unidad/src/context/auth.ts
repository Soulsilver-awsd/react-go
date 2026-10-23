// src/context/auth.ts
// Tipos, contexto y hook de autenticación (sin componentes).
// Se separa del Provider para cumplir la regla de Fast Refresh:
// un archivo que exporta componentes no debe exportar hooks/utilidades.
import { createContext, useContext } from 'react';

// Roles que maneja la aplicación
export type Rol = 'admin' | 'cliente';

// Usuario autenticado: correo + rol (lo entrega la API en /api/login)
export interface Usuario {
  email: string;
  rol: Rol;
}

export interface AuthContextType {
  isAuthenticated: boolean;
  user: Usuario | null;
  login: (usuario: Usuario, token: string) => void;
  logout: () => void;
}

// Contexto compartido; el valor real lo provee AuthProvider (AuthContext.tsx).
export const AuthContext = createContext<AuthContextType | undefined>(undefined);

// Hook con validación de tipo.
export const useAuth = () => {
  const context = useContext(AuthContext);
  if (!context) {
    throw new Error('useAuth debe ser usado dentro de un AuthProvider');
  }
  return context;
};
