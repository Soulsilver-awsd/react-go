// src/context/cart.ts
// Tipos, contexto y hook del carrito (sin componentes).
// Separado del Provider por la regla de Fast Refresh.
import { createContext, useContext } from 'react';
import type { Producto } from '../data/productos';

// Elemento del carrito: un producto + la cantidad seleccionada
export interface CartItem extends Producto {
  cantidad: number;
}

export interface CartContextType {
  cart: CartItem[];
  addToCart: (producto: Producto) => void;
  removeFromCart: (id: number) => void;
  incrementQuantity: (id: number) => void;
  decrementQuantity: (id: number) => void;
  clearCart: () => void;
  totalItems: number;
  totalPrice: number;
}

// Contexto compartido; el valor real lo provee CartProvider (CartContext.tsx).
export const CartContext = createContext<CartContextType | undefined>(undefined);

// Hook con validación de tipo.
export const useCart = () => {
  const context = useContext(CartContext);
  if (!context) {
    throw new Error('useCart debe ser usado dentro de un CartProvider');
  }
  return context;
};
