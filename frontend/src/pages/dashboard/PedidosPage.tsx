import { useEffect, useState } from 'react';
import { pedidosApi, clientesApi, inventarioApi } from '../../services/api';
import type { Pedido, Cliente, Producto } from '../../types';
import { useAuth } from '../../context/AuthContext';
import { ShoppingBag, X, Clock, CheckCircle, Ban, CreditCard, Receipt, Plus, Trash2 } from 'lucide-react';

// Paleta profesional consistente
const C = {
  bg: '#fefaf5',
  surface: '#ffffff',
  primary: '#b4532b',
  primaryDark: '#8b3a1a',
  primaryLight: '#fde4d8',
  secondary: '#3f6b5c',
  gold: '#d4af37',
  goldDark: '#b38f2c',
  textPrimary: '#2c2c2a',
  textSecondary: '#6b5e55',
  borderLight: '#efe5dc',
  alert: '#e53e3e',
  alertBg: '#fff5f5',
  success: '#2e7d64',
  warning: '#d97706',
  warningBg: '#fffbeb',
  info: '#3b82f6',
  infoBg: '#eff6ff',
};

const estadoConfig: Record<string, { label: string; color: string; bg: string; icon: any }> = {
  pendiente: {
    label: 'Pendiente',
    color: C.warning,
    bg: C.warningBg,
    icon: Clock,
  },
  entregado: {
    label: 'Entregado',
    color: C.success,
    bg: '#e6f4ea',
    icon: CheckCircle,
  },
  anulado: {
    label: 'Anulado',
    color: C.alert,
    bg: C.alertBg,
    icon: Ban,
  },
};

const pagoConfig: Record<string, { label: string; color: string; bg: string; icon: any }> = {
  contado: {
    label: 'Contado',
    color: C.info,
    bg: C.infoBg,
    icon: CreditCard,
  },
  credito: {
    label: 'Crédito',
    color: C.primary,
    bg: C.primaryLight,
    icon: Receipt,
  },
};

const formInicial = { cliente_id: '', tipo_comida: 'almuerzo', forma_pago: 'contado', notas: '' };

export default function PedidosPage() {
  const { isAdmin } = useAuth();
  const [pedidos, setPedidos] = useState<Pedido[]>([]);
  const [total, setTotal] = useState(0);
  const [filtroEstado, setFiltroEstado] = useState('');
  const [loading, setLoading] = useState(true);

  const [showForm, setShowForm] = useState(false);
  const [clientes, setClientes] = useState<Cliente[]>([]);
  const [productos, setProductos] = useState<Producto[]>([]);
  const [form, setForm] = useState(formInicial);
  const [items, setItems] = useState<{ producto_id: string; cantidad: string }[]>([{ producto_id: '', cantidad: '1' }]);
  const [guardando, setGuardando] = useState(false);

  const cargar = async () => {
    setLoading(true);
    try {
      const { data } = await pedidosApi.listar(1, 0, filtroEstado);
      setPedidos(data.data);
      setTotal(data.total);
    } catch (error) {
      console.error('Error al cargar pedidos:', error);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    cargar();
  }, [filtroEstado]);

  const abrirCrear = async () => {
    setForm(formInicial);
    setItems([{ producto_id: '', cantidad: '1' }]);
    setShowForm(true);
    try {
      const [clientesRes, productosRes] = await Promise.all([
        clientesApi.listar(1, ''),
        inventarioApi.listar(),
      ]);
      setClientes(clientesRes.data.data);
      setProductos(productosRes.data);
    } catch (error) {
      console.error('Error al cargar clientes/productos:', error);
    }
  };

  const agregarItem = () => setItems((p) => [...p, { producto_id: '', cantidad: '1' }]);
  const quitarItem = (i: number) => setItems((p) => p.filter((_, idx) => idx !== i));
  const actualizarItem = (i: number, campo: 'producto_id' | 'cantidad', valor: string) => {
    setItems((p) => p.map((it, idx) => (idx === i ? { ...it, [campo]: valor } : it)));
  };

  const totalEstimado = items.reduce((acc, it) => {
    const producto = productos.find((p) => p.id === Number(it.producto_id));
    if (!producto) return acc;
    return acc + producto.precio_venta * (Number(it.cantidad) || 0);
  }, 0);

  const handleCrearPedido = async (e: React.FormEvent) => {
    e.preventDefault();
    const itemsValidos = items.filter((it) => it.producto_id && Number(it.cantidad) > 0);
    if (itemsValidos.length === 0) {
      alert('Agrega al menos un producto al pedido');
      return;
    }
    setGuardando(true);
    try {
      await pedidosApi.crear({
        cliente_id: Number(form.cliente_id),
        tipo_comida: form.tipo_comida,
        forma_pago: form.forma_pago,
        notas: form.notas,
        items: itemsValidos.map((it) => {
          const producto = productos.find((p) => p.id === Number(it.producto_id))!;
          return {
            producto_id: producto.id,
            cantidad: Number(it.cantidad),
            precio_unitario: producto.precio_venta,
          };
        }),
      });
      setShowForm(false);
      cargar();
    } catch (err: any) {
      alert(err.response?.data?.error || 'Error al crear el pedido');
    } finally {
      setGuardando(false);
    }
  };

  const marcarEntregado = async (id: number) => {
    try {
      await pedidosApi.marcarEntregado(id);
      cargar();
    } catch (err: any) {
      alert(err.response?.data?.error || 'No se pudo marcar como entregado');
    }
  };

  const anular = async (id: number) => {
    if (!confirm('¿Anular este pedido? Se devolverá el stock.')) return;
    try {
      await pedidosApi.anular(id);
      cargar();
    } catch (err: any) {
      alert(err.response?.data?.error || 'No se pudo anular');
    }
  };

  const opcionesFiltro = [
    { value: '', label: 'Todos' },
    { value: 'pendiente', label: 'Pendientes' },
    { value: 'entregado', label: 'Entregados' },
    { value: 'anulado', label: 'Anulados' },
  ];

  return (
    <div
      style={{
        maxWidth: '1400px',
        margin: '0 auto',
        padding: '24px 20px',
        backgroundColor: C.bg,
        minHeight: '100vh',
        fontFamily: '"Inter", system-ui, sans-serif',
      }}
    >
      {/* Encabezado */}
      <div
        style={{
          display: 'flex',
          flexWrap: 'wrap',
          justifyContent: 'space-between',
          alignItems: 'center',
          marginBottom: '32px',
          gap: '16px',
        }}
      >
        <div>
          <h1
            style={{
              fontSize: 'clamp(1.75rem, 5vw, 2rem)',
              fontWeight: '700',
              color: C.textPrimary,
              margin: '0 0 6px 0',
              letterSpacing: '-0.02em',
            }}
          >
            Pedidos
          </h1>
          <p
            style={{
              fontSize: '0.9rem',
              color: C.textSecondary,
              margin: 0,
            }}
          >
            {total} pedido{total !== 1 ? 's' : ''} registrado{total !== 1 ? 's' : ''}
          </p>
        </div>
        <button
          onClick={abrirCrear}
          style={{
            display: 'flex', alignItems: 'center', gap: '8px', backgroundColor: C.primary, color: '#fff',
            border: 'none', padding: '10px 20px', borderRadius: '40px', fontSize: '0.9rem', fontWeight: 500,
            cursor: 'pointer', transition: 'all 0.2s ease', boxShadow: '0 2px 6px rgba(0,0,0,0.05)',
          }}
          onMouseEnter={(e) => { e.currentTarget.style.backgroundColor = C.primaryDark; }}
          onMouseLeave={(e) => { e.currentTarget.style.backgroundColor = C.primary; }}
        >
          <Plus size={18} /> Nuevo pedido
        </button>
      </div>

      {showForm && (
        <div
          style={{ position: 'fixed', inset: 0, backgroundColor: 'rgba(0,0,0,0.5)', backdropFilter: 'blur(4px)', zIndex: 1000, display: 'flex', alignItems: 'center', justifyContent: 'center', padding: '20px' }}
          onClick={(e) => { if (e.target === e.currentTarget) setShowForm(false); }}
        >
          <div style={{ backgroundColor: C.surface, borderRadius: '28px', maxWidth: '560px', width: '100%', maxHeight: '90vh', overflowY: 'auto', boxShadow: '0 25px 40px -12px rgba(0,0,0,0.3)' }}>
            <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', padding: '20px 24px', borderBottom: `1px solid ${C.borderLight}` }}>
              <h2 style={{ fontSize: '1.4rem', fontWeight: 600, color: C.textPrimary, margin: 0 }}>Nuevo pedido</h2>
              <button onClick={() => setShowForm(false)} style={{ background: 'none', border: 'none', cursor: 'pointer', color: C.textSecondary, padding: '4px', display: 'flex' }}>
                <X size={20} />
              </button>
            </div>
            <form onSubmit={handleCrearPedido} style={{ padding: '24px' }}>
              <div style={{ display: 'grid', gap: '20px' }}>
                <div>
                  <label style={labelStyle}>Cliente *</label>
                  <select required value={form.cliente_id} onChange={(e) => setForm((p) => ({ ...p, cliente_id: e.target.value }))} style={inputStyle}>
                    <option value="">Selecciona un cliente</option>
                    {clientes.map((c) => (
                      <option key={c.id} value={c.id}>{c.nombre} {c.apellido}</option>
                    ))}
                  </select>
                </div>

                <div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr', gap: '16px' }}>
                  <div>
                    <label style={labelStyle}>Tipo de comida *</label>
                    <select required value={form.tipo_comida} onChange={(e) => setForm((p) => ({ ...p, tipo_comida: e.target.value }))} style={inputStyle}>
                      <option value="desayuno">Desayuno</option>
                      <option value="almuerzo">Almuerzo</option>
                      <option value="cena">Cena</option>
                    </select>
                  </div>
                  <div>
                    <label style={labelStyle}>Forma de pago *</label>
                    <select required value={form.forma_pago} onChange={(e) => setForm((p) => ({ ...p, forma_pago: e.target.value }))} style={inputStyle}>
                      <option value="contado">Contado</option>
                      <option value="credito">Crédito</option>
                    </select>
                  </div>
                </div>

                <div>
                  <label style={labelStyle}>Productos *</label>
                  <div style={{ display: 'flex', flexDirection: 'column', gap: '10px' }}>
                    {items.map((it, i) => {
                      const producto = productos.find((p) => p.id === Number(it.producto_id));
                      return (
                        <div key={i} style={{ display: 'flex', gap: '8px', alignItems: 'center' }}>
                          <select value={it.producto_id} onChange={(e) => actualizarItem(i, 'producto_id', e.target.value)} style={{ ...inputStyle, flex: 2 }}>
                            <option value="">Producto...</option>
                            {productos.map((p) => (
                              <option key={p.id} value={p.id}>{p.nombre} (S/ {p.precio_venta.toFixed(2)}/{p.unidad}) — stock: {p.stock_actual}</option>
                            ))}
                          </select>
                          <input type="number" min="0.1" step="0.1" value={it.cantidad}
                            onChange={(e) => actualizarItem(i, 'cantidad', e.target.value)}
                            style={{ ...inputStyle, flex: 1 }} placeholder="Cant." />
                          <span style={{ fontSize: '0.8rem', color: C.textSecondary, minWidth: '70px' }}>
                            {producto ? `S/ ${(producto.precio_venta * (Number(it.cantidad) || 0)).toFixed(2)}` : ''}
                          </span>
                          {items.length > 1 && (
                            <button type="button" onClick={() => quitarItem(i)} style={{ background: 'none', border: 'none', cursor: 'pointer', color: C.alert, display: 'flex' }}>
                              <Trash2 size={16} />
                            </button>
                          )}
                        </div>
                      );
                    })}
                  </div>
                  <button type="button" onClick={agregarItem}
                    style={{ display: 'flex', alignItems: 'center', gap: '6px', marginTop: '10px', background: 'none', border: `1px dashed ${C.borderLight}`, borderRadius: '12px', padding: '8px 14px', color: C.primary, fontSize: '0.85rem', cursor: 'pointer' }}>
                    <Plus size={14} /> Agregar producto
                  </button>
                </div>

                <div>
                  <label style={labelStyle}>Notas (opcional)</label>
                  <input value={form.notas} onChange={(e) => setForm((p) => ({ ...p, notas: e.target.value }))} style={inputStyle} />
                </div>

                <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', borderTop: `1px solid ${C.borderLight}`, paddingTop: '16px' }}>
                  <span style={{ color: C.textSecondary, fontSize: '0.9rem' }}>Total estimado</span>
                  <span style={{ fontWeight: 700, fontSize: '1.3rem', color: C.gold }}>S/ {totalEstimado.toFixed(2)}</span>
                </div>
              </div>

              <div style={{ display: 'flex', gap: '12px', marginTop: '24px' }}>
                <button type="button" onClick={() => setShowForm(false)}
                  style={{ flex: 1, padding: '12px', borderRadius: '40px', border: `1px solid ${C.borderLight}`, backgroundColor: 'transparent', color: C.textSecondary, fontWeight: 500, cursor: 'pointer' }}>
                  Cancelar
                </button>
                <button type="submit" disabled={guardando}
                  style={{ flex: 1, padding: '12px', borderRadius: '40px', border: 'none', backgroundColor: C.primary, color: '#fff', fontWeight: 500, cursor: 'pointer', opacity: guardando ? 0.7 : 1 }}>
                  {guardando ? 'Guardando...' : 'Crear pedido'}
                </button>
              </div>
            </form>
          </div>
        </div>
      )}

      {/* Filtros tipo Pills */}
      <div
        style={{
          display: 'flex',
          flexWrap: 'wrap',
          gap: '10px',
          marginBottom: '28px',
          borderBottom: `1px solid ${C.borderLight}`,
          paddingBottom: '16px',
        }}
      >
        {opcionesFiltro.map((op) => (
          <button
            key={op.value}
            onClick={() => setFiltroEstado(op.value)}
            style={{
              padding: '8px 20px',
              borderRadius: '40px',
              fontSize: '0.85rem',
              fontWeight: '500',
              cursor: 'pointer',
              transition: 'all 0.2s ease',
              border: 'none',
              backgroundColor: filtroEstado === op.value ? C.primary : 'transparent',
              color: filtroEstado === op.value ? '#fff' : C.textSecondary,
              boxShadow: filtroEstado === op.value ? `0 2px 8px ${C.primary}40` : 'none',
            }}
            onMouseEnter={(e) => {
              if (filtroEstado !== op.value) {
                e.currentTarget.style.backgroundColor = C.borderLight;
                e.currentTarget.style.color = C.textPrimary;
              }
            }}
            onMouseLeave={(e) => {
              if (filtroEstado !== op.value) {
                e.currentTarget.style.backgroundColor = 'transparent';
                e.currentTarget.style.color = C.textSecondary;
              }
            }}
          >
            {op.label}
          </button>
        ))}
      </div>

      {/* Lista de pedidos */}
      {loading ? (
        <div style={{ display: 'flex', flexDirection: 'column', gap: '14px' }}>
          {[...Array(4)].map((_, i) => (
            <div
              key={i}
              style={{
                height: '110px',
                backgroundColor: C.surface,
                borderRadius: '20px',
                border: `1px solid ${C.borderLight}`,
                position: 'relative',
                overflow: 'hidden',
              }}
              className="shimmer"
            />
          ))}
        </div>
      ) : pedidos.length === 0 ? (
        <div
          style={{
            textAlign: 'center',
            padding: '60px 24px',
            backgroundColor: C.surface,
            borderRadius: '28px',
            border: `1px solid ${C.borderLight}`,
          }}
        >
          <ShoppingBag
            size={56}
            style={{ color: C.textSecondary, opacity: 0.4, marginBottom: '16px' }}
          />
          <p style={{ color: C.textSecondary, fontSize: '1rem', margin: 0 }}>
            {filtroEstado ? `No hay pedidos ${filtroEstado}s` : 'No hay pedidos registrados'}
          </p>
          <p style={{ color: C.textSecondary, fontSize: '0.85rem', marginTop: '8px' }}>
            Los pedidos aparecerán aquí cuando se realicen.
          </p>
        </div>
      ) : (
        <div style={{ display: 'flex', flexDirection: 'column', gap: '14px' }}>
          {pedidos.map((pedido) => {
            const EstadoIcon = estadoConfig[pedido.estado]?.icon || Clock;
            const PagoIcon = pagoConfig[pedido.forma_pago]?.icon || CreditCard;
            return (
              <div
                key={pedido.id}
                style={{
                  backgroundColor: C.surface,
                  borderRadius: '20px',
                  border: `1px solid ${C.borderLight}`,
                  padding: '18px 20px',
                  transition: 'all 0.25s ease',
                  cursor: 'pointer',
                }}
                onMouseEnter={(e) => {
                  e.currentTarget.style.borderColor = C.primary;
                  e.currentTarget.style.boxShadow = '0 12px 24px -12px rgba(0,0,0,0.12)';
                  e.currentTarget.style.transform = 'translateY(-2px)';
                }}
                onMouseLeave={(e) => {
                  e.currentTarget.style.borderColor = C.borderLight;
                  e.currentTarget.style.boxShadow = 'none';
                  e.currentTarget.style.transform = 'translateY(0)';
                }}
              >
                <div
                  style={{
                    display: 'flex',
                    flexWrap: 'wrap',
                    justifyContent: 'space-between',
                    alignItems: 'center',
                    gap: '16px',
                  }}
                >
                  {/* Lado izquierdo: icono + info principal */}
                  <div style={{ display: 'flex', alignItems: 'center', gap: '16px', flex: '2', minWidth: '200px' }}>
                    <div
                      style={{
                        width: '52px',
                        height: '52px',
                        borderRadius: '20px',
                        display: 'flex',
                        alignItems: 'center',
                        justifyContent: 'center',
                        backgroundColor: C.primaryLight,
                      }}
                    >
                      <ShoppingBag size={24} color={C.primary} />
                    </div>
                    <div>
                      <div style={{ display: 'flex', alignItems: 'center', gap: '10px', flexWrap: 'wrap', marginBottom: '6px' }}>
                        <span
                          style={{
                            fontWeight: '700',
                            fontSize: '1rem',
                            color: C.textPrimary,
                            letterSpacing: '-0.01em',
                          }}
                        >
                          Pedido #{pedido.id}
                        </span>
                        <span
                          style={{
                            fontSize: '0.7rem',
                            padding: '2px 8px',
                            borderRadius: '40px',
                            backgroundColor: C.borderLight,
                            color: C.textSecondary,
                          }}
                        >
                          {pedido.tipo_comida}
                        </span>
                      </div>
                      <div
                        style={{
                          display: 'flex',
                          flexWrap: 'wrap',
                          gap: '12px',
                          fontSize: '0.75rem',
                          color: C.textSecondary,
                        }}
                      >
                        <span>
                          {pedido.cliente?.nombre || 'Cliente no especificado'}
                        </span>
                        <span>•</span>
                        <span>
                          {new Date(pedido.fecha).toLocaleDateString('es-PE', {
                            day: 'numeric',
                            month: 'short',
                            hour: '2-digit',
                            minute: '2-digit',
                          })}
                        </span>
                      </div>
                    </div>
                  </div>

                  {/* Lado derecho: badges + total + botón anular */}
                  <div
                    style={{
                      display: 'flex',
                      alignItems: 'center',
                      gap: '12px',
                      flexWrap: 'wrap',
                    }}
                  >
                    {/* Badge estado */}
                    <div
                      style={{
                        display: 'flex',
                        alignItems: 'center',
                        gap: '6px',
                        backgroundColor: estadoConfig[pedido.estado]?.bg || '#f3f4f6',
                        padding: '6px 12px',
                        borderRadius: '40px',
                        color: estadoConfig[pedido.estado]?.color || C.textSecondary,
                        fontSize: '0.7rem',
                        fontWeight: '600',
                      }}
                    >
                      <EstadoIcon size={12} />
                      <span>{estadoConfig[pedido.estado]?.label || pedido.estado}</span>
                    </div>

                    {/* Badge forma pago */}
                    <div
                      style={{
                        display: 'flex',
                        alignItems: 'center',
                        gap: '6px',
                        backgroundColor: pagoConfig[pedido.forma_pago]?.bg || '#f3f4f6',
                        padding: '6px 12px',
                        borderRadius: '40px',
                        color: pagoConfig[pedido.forma_pago]?.color || C.textSecondary,
                        fontSize: '0.7rem',
                        fontWeight: '600',
                      }}
                    >
                      <PagoIcon size={12} />
                      <span>{pagoConfig[pedido.forma_pago]?.label || pedido.forma_pago}</span>
                    </div>

                    {/* Total */}
                    <span
                      style={{
                        fontWeight: '800',
                        fontSize: '1.15rem',
                        color: C.textPrimary,
                        letterSpacing: '-0.01em',
                      }}
                    >
                      S/ {pedido.total.toFixed(2)}
                    </span>

                    {/* Botón marcar entregado (solo si está pendiente) */}
                    {pedido.estado === 'pendiente' && (
                      <button
                        onClick={(e) => {
                          e.stopPropagation();
                          marcarEntregado(pedido.id);
                        }}
                        style={{
                          width: '36px',
                          height: '36px',
                          borderRadius: '12px',
                          backgroundColor: '#e6f4ea',
                          border: `1px solid ${C.success}30`,
                          display: 'flex',
                          alignItems: 'center',
                          justifyContent: 'center',
                          cursor: 'pointer',
                          transition: 'all 0.2s',
                        }}
                        onMouseEnter={(e) => {
                          e.currentTarget.style.backgroundColor = C.success;
                          const svg = e.currentTarget.querySelector('svg');
                          if (svg) svg.style.color = '#fff';
                        }}
                        onMouseLeave={(e) => {
                          e.currentTarget.style.backgroundColor = '#e6f4ea';
                          const svg = e.currentTarget.querySelector('svg');
                          if (svg) svg.style.color = C.success;
                        }}
                        title="Marcar como entregado"
                      >
                        <CheckCircle size={16} color={C.success} />
                      </button>
                    )}

                    {/* Botón anular (solo admin y pedido no anulado) */}
                    {isAdmin && pedido.estado !== 'anulado' && (
                      <button
                        onClick={(e) => {
                          e.stopPropagation();
                          anular(pedido.id);
                        }}
                        style={{
                          width: '36px',
                          height: '36px',
                          borderRadius: '12px',
                          backgroundColor: C.alertBg,
                          border: `1px solid ${C.alert}30`,
                          display: 'flex',
                          alignItems: 'center',
                          justifyContent: 'center',
                          cursor: 'pointer',
                          transition: 'all 0.2s',
                        }}
                        onMouseEnter={(e) => {
                          e.currentTarget.style.backgroundColor = C.alert;
                          e.currentTarget.style.borderColor = C.alert;
                          const svg = e.currentTarget.querySelector('svg');
                          if (svg) svg.style.color = '#fff';
                        }}
                        onMouseLeave={(e) => {
                          e.currentTarget.style.backgroundColor = C.alertBg;
                          e.currentTarget.style.borderColor = `${C.alert}30`;
                          const svg = e.currentTarget.querySelector('svg');
                          if (svg) svg.style.color = C.alert;
                        }}
                        title="Anular pedido"
                      >
                        <X size={16} color={C.alert} />
                      </button>
                    )}
                  </div>
                </div>
              </div>
            );
          })}
        </div>
      )}

    </div>
  );
}

const labelStyle: React.CSSProperties = {
  display: 'block', fontSize: '0.75rem', fontWeight: 500, color: C.textSecondary,
  marginBottom: '6px', textTransform: 'uppercase', letterSpacing: '0.5px',
};

const inputStyle: React.CSSProperties = {
  width: '100%', padding: '10px 14px', border: `1px solid ${C.borderLight}`, borderRadius: '16px',
  fontSize: '0.9rem', backgroundColor: C.surface, color: C.textPrimary, outline: 'none',
};