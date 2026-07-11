import axios from 'axios'

const api = axios.create({ baseURL: '/api' })

// Inyecta el JWT en cada request automáticamente
api.interceptors.request.use((config) => {
  const token = localStorage.getItem('token')
  if (token) config.headers.Authorization = `Bearer ${token}`
  return config
})

// Si el servidor devuelve 401, limpia la sesión y redirige al login
api.interceptors.response.use(
  (r) => r,
  (error) => {
    if (error.response?.status === 401) {
      localStorage.removeItem('token')
      localStorage.removeItem('user')
      window.location.href = '/login'
    }
    return Promise.reject(error)
  }
)

// ─── Auth ─────────────────────────────────────────────────────────────────────
export const authApi = {
  login: (email: string, password: string) =>
    api.post('/auth/login', { email, password }),
  me: () => api.get('/auth/me'),
}

// ─── Clientes ─────────────────────────────────────────────────────────────────
export const clientesApi = {
  listar: (page = 1, tipo = '') =>
    api.get(`/empleado/clientes?page=${page}&tipo=${tipo}`),
  obtener: (id: number) => api.get(`/empleado/clientes/${id}`),
  crear: (data: object) => api.post('/admin/clientes', data),
  actualizar: (id: number, data: object) => api.put(`/admin/clientes/${id}`, data),
  desactivar: (id: number) => api.delete(`/admin/clientes/${id}`),
}

// ─── Pedidos ──────────────────────────────────────────────────────────────────
export const pedidosApi = {
  listar: (page = 1, clienteId = 0, estado = '') =>
    api.get(`/empleado/pedidos?page=${page}&cliente_id=${clienteId}&estado=${estado}`),
  obtener: (id: number) => api.get(`/empleado/pedidos/${id}`),
  crear: (data: object) => api.post('/empleado/pedidos', data),
  anular: (id: number) => api.delete(`/admin/pedidos/${id}`),
  marcarEntregado: (id: number) => api.put(`/empleado/pedidos/${id}/entregar`),
}

// ─── Inventario ───────────────────────────────────────────────────────────────
export const inventarioApi = {
  listar: () => api.get('/empleado/productos'),
  alertas: () => api.get('/admin/productos/alertas'),
  obtener: (id: number) => api.get(`/admin/productos/${id}`),
  crear: (data: object) => api.post('/admin/productos', data),
  actualizar: (id: number, data: object) => api.put(`/admin/productos/${id}`, data),
  restock: (id: number, cantidad: number, notas: string) =>
    api.post(`/admin/productos/${id}/restock`, { cantidad, notas }),
}

// ─── Créditos y Pagos ─────────────────────────────────────────────────────────
export const creditosApi = {
  deudores: () => api.get('/admin/creditos'),
  estadoCuenta: (clienteId: number) => api.get(`/admin/creditos/${clienteId}`),
  registrarPago: (data: object) => api.post('/admin/pagos', data),
}

// ─── IA ───────────────────────────────────────────────────────────────────────
export const iaApi = {
  prediccion: () => api.get('/admin/ia/prediccion'),
  alertas: () => api.get('/admin/ia/alertas'),
}

// ─── Reportes ─────────────────────────────────────────────────────────────────
export const reportesApi = {
  ventas: (periodo: string, fecha?: string) =>
    api.get(`/admin/reportes/ventas?periodo=${periodo}${fecha ? `&fecha=${fecha}` : ''}`),
  ventasArchivo: (periodo: string, formato: 'pdf' | 'excel', fecha?: string) =>
    api.get(`/admin/reportes/ventas?periodo=${periodo}&formato=${formato}${fecha ? `&fecha=${fecha}` : ''}`, { responseType: 'blob' }),
  deudores: () => api.get('/admin/reportes/deudores'),
  deudoresArchivo: (formato: 'pdf' | 'excel') =>
    api.get(`/admin/reportes/deudores?formato=${formato}`, { responseType: 'blob' }),
  inventario: (desde?: string, hasta?: string) =>
    api.get(`/admin/reportes/inventario?desde=${desde || ''}&hasta=${hasta || ''}`),
  inventarioArchivo: (formato: 'pdf' | 'excel', desde?: string, hasta?: string) =>
    api.get(`/admin/reportes/inventario?formato=${formato}&desde=${desde || ''}&hasta=${hasta || ''}`, { responseType: 'blob' }),
}

// ─── Reservas ─────────────────────────────────────────────────────────────────
export const reservasApi = {
  crear: (data: object) => api.post('/public/reservas', data),
  listar: (page = 1, estado = '') =>
    api.get(`/empleado/reservas?page=${page}&estado=${estado}`),
  actualizarEstado: (id: number, estado: 'pendiente' | 'confirmada' | 'cancelada') =>
    api.put(`/empleado/reservas/${id}/estado`, { estado }),
  vincularPedido: (id: number, pedidoId: number) =>
    api.put(`/empleado/reservas/${id}/pedido`, { pedido_id: pedidoId }),
}

// ─── Menú público ─────────────────────────────────────────────────────────────
export const menuApi = {
  obtener: () => api.get('/public/menu'),
  listarAdmin: () => api.get('/admin/menu'),
  crear: (data: object) => api.post('/admin/menu', data),
  actualizar: (id: number, data: object) => api.put(`/admin/menu/${id}`, data),
  eliminar: (id: number) => api.delete(`/admin/menu/${id}`),
}
