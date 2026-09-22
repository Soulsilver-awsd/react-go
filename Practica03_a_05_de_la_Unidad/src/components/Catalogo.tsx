// src/components/Catalogo.tsx
import { useEffect, useState } from 'react';
import { getProductos } from '../services/api';
import { useCart, type Producto } from '../context/CartContext';

const Catalogo = () => {
  const { addToCart } = useCart(); // <-- Usamos la función del contexto
  const [productos, setProductos] = useState<Producto[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState('');

  useEffect(() => {
    getProductos()
      .then(setProductos)
      .catch((requestError) => {
        setError(requestError instanceof Error ? requestError.message : 'No se pudo cargar el catálogo');
      })
      .finally(() => setIsLoading(false));
  }, []);

  return (
    <div>
      <h1 className="text-2xl font-bold text-slate-800 mb-6">Catálogo de Productos</h1>
      {isLoading && <p className="text-slate-500">Cargando productos...</p>}
      {error && <p className="text-red-600">{error}</p>}
      {!isLoading && !error && (
      <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-6">
        {productos.map((prod) => (
          <div key={prod.id} className="bg-white rounded-lg overflow-hidden border border-slate-200 shadow-sm flex flex-col">
            <img src={prod.img} alt={prod.nombre} className="w-full h-40 object-cover" />
            <div className="p-4 flex flex-col flex-1">
              <h3 className="font-semibold text-slate-700">{prod.nombre}</h3>
              <p className="text-indigo-600 font-bold mt-2 mb-4">${prod.precio.toFixed(2)}</p>              
             
              <button 
                onClick={() => addToCart(prod)}
                className="mt-auto w-full bg-slate-900 text-white py-2 rounded text-sm hover:bg-indigo-600 transition"
              >
                Añadir al Carrito
              </button>
            </div>
          </div>
        ))}
      </div>
      )}
    </div>
  );
};

export default Catalogo;