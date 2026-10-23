// src/context/AuthContext.tsx
// Provider de autenticación (componente). El contexto y el hook viven en auth.ts.
import { useState, type ReactNode } from 'react';
import { AuthContext, type Usuario, type AuthContextType } from './auth';
import {
  guardarSesion,
  limpiarSesion,
  getSesion,
  getToken,
  type SesionPersistida,
} from '../services/http';

interface AuthProviderProps {
  children: ReactNode;
}

export const AuthProvider = ({ children }: AuthProviderProps) => {
  // Rehidratamos la sesión desde localStorage para que un refresh de la SPA
  // no cierre la sesión (persistencia U1-T3).
  const persistida = getSesion();
  const [isAuthenticated, setIsAuthenticated] = useState<boolean>(
    Boolean(persistida && getToken())
  );
  const [user, setUser] = useState<Usuario | null>(
    persistida ? { email: persistida.email, rol: persistida.rol } : null
  );

  const login = (usuario: Usuario, token: string) => {
    setIsAuthenticated(true);
    setUser(usuario);
    // Persistimos usuario + JWT bajo la clave multicatalogo_sesion.
    const sesion: SesionPersistida = {
      email: usuario.email,
      rol: usuario.rol,
      token,
    };
    guardarSesion(sesion);
  };

  const logout = () => {
    setIsAuthenticated(false);
    setUser(null);
    limpiarSesion();
  };

  const value: AuthContextType = { isAuthenticated, user, login, logout };

  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>;
};
