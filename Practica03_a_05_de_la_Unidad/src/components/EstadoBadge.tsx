// src/components/EstadoBadge.tsx
// Badge de color según el estado del pedido (reutilizado por cliente y admin).
import type { EstadoPedido } from "../services/pedidosService";

const ESTILOS: Record<EstadoPedido, string> = {
  pendiente: "bg-amber-100 text-amber-700",
  preparando: "bg-blue-100 text-blue-700",
  enviado: "bg-indigo-100 text-indigo-700",
  entregado: "bg-green-100 text-green-700",
  cancelado: "bg-red-100 text-red-600",
};

const EstadoBadge = ({ estado }: { estado: EstadoPedido }) => (
  <span
    className={`text-xs font-semibold px-2 py-1 rounded-full capitalize ${
      ESTILOS[estado] ?? "bg-slate-100 text-slate-600"
    }`}
  >
    {estado}
  </span>
);

export default EstadoBadge;
