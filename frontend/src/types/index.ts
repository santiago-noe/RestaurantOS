export interface User {
  id: number
  nombre: string
  apellido: string
  email: string
  rol: 'admin' | 'empleado'
}

export interface Cliente {
  id: number
  nombre: string
  apellido: string
  tipo: 'individual' | 'empresa'
  telefono: string
  email: string
  deuda_total: number
  activo: boolean
}

export interface Producto {
  id: number
  nombre: string
  unidad: string
  stock_actual: number
  stock_minimo: number
  precio_venta: number
  activo: boolean
}

export interface ItemPedido {
  producto_id: number
  cantidad: number
  precio_unitario: number
  subtotal: number
  producto?: Producto
}

export interface Pedido {
  id: number
  cliente_id: number
  fecha: string
  tipo_comida: string
  estado: 'pendiente' | 'entregado' | 'anulado'
  forma_pago: 'contado' | 'credito'
  total: number
  notas: string
  cliente?: Cliente
  items?: ItemPedido[]
}

export interface Pago {
  id: number
  cliente_id: number
  monto: number
  metodo: string
  fecha: string
}

export interface Alerta {
  tipo: string
  mensaje: string
  severidad: 'alta' | 'media' | 'baja'
  ref_id?: number
}

export interface Prediccion {
  tipo_comida: string
  porciones_estimadas: number
  base_calculo: string
}

export interface MenuPublico {
  id: number
  categoria: string
  nombre: string
  descripcion: string
  precio: number
  imagen_url: string
  disponible: boolean
  orden: number
  producto_id?: number | null
}

export interface MovimientoStock {
  id: number
  producto_id: number
  tipo: 'entrada' | 'salida'
  cantidad: number
  fecha: string
  producto?: Producto
}

export interface VentaDia {
  fecha: string
  total: number
  cantidad_pedidos: number
}

export interface ReporteVentas {
  periodo: string
  desde: string
  hasta: string
  total_ventas: number
  total_pedidos: number
  por_dia: VentaDia[]
}

export interface ReporteMovimientos {
  desde: string
  hasta: string
  total_entradas: number
  total_salidas: number
  movimientos: MovimientoStock[]
}

export interface Reserva {
  id: number
  nombre: string
  whatsapp: string
  fecha: string
  personas: string
  ocasion: string
  estado: 'pendiente' | 'confirmada' | 'cancelada'
  pedido_id?: number | null
  created_at: string
}

export interface PaginatedResponse<T> {
  data: T[]
  total: number
  page: number
  per_page: number
}
