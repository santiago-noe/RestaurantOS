import { useEffect, useState } from 'react';
import { inventarioApi } from '../../services/api';
import type { Producto } from '../../types';
import { Package, AlertTriangle, Plus, RefreshCw, TrendingDown, CheckCircle, X } from 'lucide-react';

// Paleta profesional
const C = {
  bg: '#fefaf5',
  surface: '#ffffff',
  primary: '#b4532b',
  primaryDark: '#8b3a1a',
  primaryLight: '#fde4d8',
  secondary: '#3f6b5c',
  textPrimary: '#2c2c2a',
  textSecondary: '#6b5e55',
  borderLight: '#efe5dc',
  alert: '#e53e3e',
  alertBg: '#fff5f5',
  success: '#2e7d64',
  warning: '#d97706',
  warningBg: '#fffbeb',
};

export default function InventarioPage() {
  const [productos, setProductos] = useState<Producto[]>([]);
  const [alertas, setAlertas] = useState<any[]>([]);
  const [loading, setLoading] = useState(true);
  const [showRestock, setShowRestock] = useState<number | null>(null);
  const [cantidad, setCantidad] = useState('');
  const [notas, setNotas] = useState('');

  const cargar = async () => {
    setLoading(true);
    try {
      const [prods, alts] = await Promise.all([
        inventarioApi.listar(),
        inventarioApi.alertas(),
      ]);
      setProductos(prods.data);
      setAlertas(alts.data.alertas || []);
    } catch (error) {
      console.error('Error al cargar inventario:', error);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    cargar();
  }, []);

  const handleRestock = async (id: number) => {
    const cant = parseFloat(cantidad);
    if (!cant || cant <= 0) {
      alert('Ingresa una cantidad válida');
      return;
    }
    try {
      await inventarioApi.restock(id, cant, notas);
      setShowRestock(null);
      setCantidad('');
      setNotas('');
      cargar();
    } catch (err: any) {
      alert(err.response?.data?.error || 'Error al reabastecer');
    }
  };

  const getStockStatus = (producto: Producto) => {
    const ratio = producto.stock_actual / producto.stock_minimo;
    if (ratio <= 0.5) return { label: 'Crítico', color: C.alert, bg: C.alertBg, icon: AlertTriangle };
    if (ratio < 1) return { label: 'Bajo stock', color: C.warning, bg: C.warningBg, icon: TrendingDown };
    return { label: 'Normal', color: C.success, bg: '#e6f4ea', icon: CheckCircle };
  };

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
      {/* Header */}
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
            Inventario
          </h1>
          <p
            style={{
              fontSize: '0.9rem',
              color: C.textSecondary,
              margin: 0,
            }}
          >
            {productos.length} producto{productos.length !== 1 ? 's' : ''} registrado{productos.length !== 1 ? 's' : ''} · {alertas.length} alerta{alertas.length !== 1 ? 's' : ''}
          </p>
        </div>
        <button
          onClick={cargar}
          style={{
            display: 'flex',
            alignItems: 'center',
            gap: '8px',
            backgroundColor: 'transparent',
            border: `1px solid ${C.borderLight}`,
            padding: '8px 18px',
            borderRadius: '40px',
            fontSize: '0.85rem',
            fontWeight: '500',
            color: C.textSecondary,
            cursor: 'pointer',
            transition: 'all 0.2s',
          }}
          onMouseEnter={(e) => {
            e.currentTarget.style.backgroundColor = C.surface;
            e.currentTarget.style.borderColor = C.primary;
            e.currentTarget.style.color = C.primary;
          }}
          onMouseLeave={(e) => {
            e.currentTarget.style.backgroundColor = 'transparent';
            e.currentTarget.style.borderColor = C.borderLight;
            e.currentTarget.style.color = C.textSecondary;
          }}
        >
          <RefreshCw size={16} /> Actualizar
        </button>
      </div>

      {/* Tabla de alertas (si existen) */}
      {alertas.length > 0 && (
        <div
          style={{
            marginBottom: '32px',
            backgroundColor: C.surface,
            borderRadius: '20px',
            border: `1px solid ${C.alert}30`,
            overflow: 'hidden',
            boxShadow: '0 2px 8px rgba(0,0,0,0.02)',
          }}
        >
          <div
            style={{
              backgroundColor: C.alertBg,
              padding: '14px 20px',
              borderBottom: `1px solid ${C.alert}20`,
              display: 'flex',
              alignItems: 'center',
              gap: '10px',
            }}
          >
            <AlertTriangle size={18} color={C.alert} />
            <span style={{ fontWeight: '700', fontSize: '0.9rem', color: C.alert }}>
              Productos bajo stock mínimo
            </span>
            <span
              style={{
                marginLeft: 'auto',
                backgroundColor: C.alert,
                color: '#fff',
                borderRadius: '40px',
                padding: '2px 10px',
                fontSize: '0.7rem',
                fontWeight: '600',
              }}
            >
              {alertas.length}
            </span>
          </div>
          <div style={{ padding: '8px 0' }}>
            {alertas.map((a: any, idx: number) => (
              <div
                key={idx}
                style={{
                  padding: '12px 20px',
                  borderBottom: idx !== alertas.length - 1 ? `1px solid ${C.borderLight}` : 'none',
                  display: 'flex',
                  alignItems: 'center',
                  gap: '12px',
                  fontSize: '0.85rem',
                }}
              >
                <TrendingDown size={14} color={C.alert} />
                <span style={{ color: C.textPrimary }}>{a.mensaje}</span>
              </div>
            ))}
          </div>
        </div>
      )}

      {/* Tabla de productos */}
      <div
        style={{
          backgroundColor: C.surface,
          borderRadius: '24px',
          border: `1px solid ${C.borderLight}`,
          overflow: 'hidden',
          boxShadow: '0 4px 12px rgba(0,0,0,0.02)',
        }}
      >
        <div style={{ overflowX: 'auto' }}>
          <table
            style={{
              width: '100%',
              borderCollapse: 'collapse',
              fontSize: '0.85rem',
              minWidth: '600px',
            }}
          >
            <thead>
              <tr
                style={{
                  backgroundColor: C.bg,
                  borderBottom: `1px solid ${C.borderLight}`,
                }}
              >
                <th style={thStyle}>Producto</th>
                <th style={thStyle}>Unidad</th>
                <th style={thStyle}>Stock actual</th>
                <th style={thStyle}>Stock mínimo</th>
                <th style={thStyle}>Estado</th>
                <th style={{ ...thStyle, textAlign: 'center' }}>Acciones</th>
              </tr>
            </thead>
            <tbody>
              {loading ? (
                // Skeleton rows
                [...Array(4)].map((_, i) => (
                  <tr key={i} style={{ borderBottom: `1px solid ${C.borderLight}` }}>
                    <td colSpan={6} style={{ padding: '16px 20px' }}>
                      <div
                        style={{
                          height: '36px',
                          backgroundColor: C.bg,
                          borderRadius: '8px',
                          position: 'relative',
                          overflow: 'hidden',
                        }}
                        className="shimmer"
                      />
                    </td>
                  </tr>
                ))
              ) : productos.length === 0 ? (
                <tr>
                  <td colSpan={6} style={{ textAlign: 'center', padding: '60px 20px', color: C.textSecondary }}>
                    <Package size={40} style={{ marginBottom: '12px', opacity: 0.4 }} />
                    <p>No hay productos en el inventario</p>
                  </td>
                </tr>
              ) : (
                productos.map((producto) => {
                  const status = getStockStatus(producto);
                  const StatusIcon = status.icon;
                  return (
                    <tr
                      key={producto.id}
                      style={{
                        borderBottom: `1px solid ${C.borderLight}`,
                        transition: 'background 0.2s',
                      }}
                      onMouseEnter={(e) => {
                        e.currentTarget.style.backgroundColor = C.bg;
                      }}
                      onMouseLeave={(e) => {
                        e.currentTarget.style.backgroundColor = 'transparent';
                      }}
                    >
                      <td style={tdStyle}>
                        <div style={{ display: 'flex', alignItems: 'center', gap: '12px' }}>
                          <div
                            style={{
                              width: '36px',
                              height: '36px',
                              borderRadius: '12px',
                              backgroundColor: C.primaryLight,
                              display: 'flex',
                              alignItems: 'center',
                              justifyContent: 'center',
                            }}
                          >
                            <Package size={18} color={C.primary} />
                          </div>
                          <span style={{ fontWeight: '500', color: C.textPrimary }}>
                            {producto.nombre}
                          </span>
                        </div>
                      </td>
                      <td style={tdStyle}>
                        <span
                          style={{
                            backgroundColor: C.borderLight,
                            padding: '4px 10px',
                            borderRadius: '40px',
                            fontSize: '0.7rem',
                            fontWeight: '500',
                            color: C.textSecondary,
                          }}
                        >
                          {producto.unidad}
                        </span>
                      </td>
                      <td style={{ ...tdStyle, fontWeight: '700', color: status.color }}>
                        {producto.stock_actual}
                      </td>
                      <td style={tdStyle}>{producto.stock_minimo}</td>
                      <td style={tdStyle}>
                        <span
                          style={{
                            display: 'inline-flex',
                            alignItems: 'center',
                            gap: '6px',
                            backgroundColor: status.bg,
                            color: status.color,
                            padding: '4px 12px',
                            borderRadius: '40px',
                            fontSize: '0.7rem',
                            fontWeight: '600',
                          }}
                        >
                          <StatusIcon size={12} />
                          {status.label}
                        </span>
                      </td>
                      <td style={{ ...tdStyle, textAlign: 'center' }}>
                        <button
                          onClick={() => setShowRestock(producto.id)}
                          style={{
                            backgroundColor: C.primaryLight,
                            border: 'none',
                            width: '34px',
                            height: '34px',
                            borderRadius: '12px',
                            display: 'inline-flex',
                            alignItems: 'center',
                            justifyContent: 'center',
                            cursor: 'pointer',
                            transition: 'all 0.2s',
                          }}
                          onMouseEnter={(e) => {
                            e.currentTarget.style.backgroundColor = C.primary;
                            const svg = e.currentTarget.querySelector('svg');
                            if (svg) svg.style.color = '#fff';
                          }}
                          onMouseLeave={(e) => {
                            e.currentTarget.style.backgroundColor = C.primaryLight;
                            const svg = e.currentTarget.querySelector('svg');
                            if (svg) svg.style.color = C.primary;
                          }}
                          title="Reabastecer"
                        >
                          <Plus size={16} color={C.primary} />
                        </button>
                      </td>
                    </tr>
                  );
                })
              )}
            </tbody>
          </table>
        </div>
      </div>

      {/* Modal reabastecer (sin cambios en funcionalidad, mismo diseño moderno) */}
      {showRestock && (
        <div
          style={{
            position: 'fixed',
            inset: 0,
            backgroundColor: 'rgba(0,0,0,0.5)',
            backdropFilter: 'blur(4px)',
            zIndex: 1000,
            display: 'flex',
            alignItems: 'center',
            justifyContent: 'center',
            padding: '20px',
            animation: 'fadeIn 0.2s ease',
          }}
          onClick={(e) => {
            if (e.target === e.currentTarget) {
              setShowRestock(null);
              setCantidad('');
              setNotas('');
            }
          }}
        >
          <div
            style={{
              backgroundColor: C.surface,
              borderRadius: '28px',
              maxWidth: '420px',
              width: '100%',
              boxShadow: '0 25px 40px -12px rgba(0,0,0,0.3)',
              overflow: 'hidden',
            }}
          >
            <div
              style={{
                display: 'flex',
                justifyContent: 'space-between',
                alignItems: 'center',
                padding: '20px 24px',
                borderBottom: `1px solid ${C.borderLight}`,
              }}
            >
              <h3
                style={{
                  fontSize: '1.3rem',
                  fontWeight: '600',
                  color: C.textPrimary,
                  margin: 0,
                }}
              >
                Reabastecer producto
              </h3>
              <button
                onClick={() => {
                  setShowRestock(null);
                  setCantidad('');
                  setNotas('');
                }}
                style={{
                  background: 'none',
                  border: 'none',
                  cursor: 'pointer',
                  padding: '4px',
                  borderRadius: '50%',
                  display: 'flex',
                }}
              >
                <X size={20} color={C.textSecondary} />
              </button>
            </div>
            <div style={{ padding: '24px' }}>
              <div style={{ marginBottom: '20px' }}>
                <label
                  style={{
                    display: 'block',
                    fontSize: '0.75rem',
                    fontWeight: '500',
                    color: C.textSecondary,
                    marginBottom: '6px',
                    textTransform: 'uppercase',
                    letterSpacing: '0.5px',
                  }}
                >
                  Cantidad a agregar *
                </label>
                <input
                  type="number"
                  step="0.1"
                  min="0.1"
                  value={cantidad}
                  onChange={(e) => setCantidad(e.target.value)}
                  placeholder="Ej: 10.5"
                  style={inputStyle}
                  autoFocus
                />
              </div>
              <div style={{ marginBottom: '24px' }}>
                <label
                  style={{
                    display: 'block',
                    fontSize: '0.75rem',
                    fontWeight: '500',
                    color: C.textSecondary,
                    marginBottom: '6px',
                    textTransform: 'uppercase',
                    letterSpacing: '0.5px',
                  }}
                >
                  Notas (opcional)
                </label>
                <input
                  value={notas}
                  onChange={(e) => setNotas(e.target.value)}
                  placeholder="Compra del mercado, proveedor..."
                  style={inputStyle}
                />
              </div>
              <div style={{ display: 'flex', gap: '12px' }}>
                <button
                  onClick={() => {
                    setShowRestock(null);
                    setCantidad('');
                    setNotas('');
                  }}
                  style={{
                    flex: 1,
                    padding: '12px',
                    borderRadius: '40px',
                    border: `1px solid ${C.borderLight}`,
                    backgroundColor: 'transparent',
                    color: C.textSecondary,
                    fontWeight: '500',
                    cursor: 'pointer',
                    transition: 'all 0.2s',
                  }}
                  onMouseEnter={(e) => {
                    e.currentTarget.style.backgroundColor = '#f7f7f5';
                  }}
                  onMouseLeave={(e) => {
                    e.currentTarget.style.backgroundColor = 'transparent';
                  }}
                >
                  Cancelar
                </button>
                <button
                  onClick={() => handleRestock(showRestock)}
                  style={{
                    flex: 1,
                    padding: '12px',
                    borderRadius: '40px',
                    border: 'none',
                    backgroundColor: C.primary,
                    color: '#fff',
                    fontWeight: '500',
                    cursor: 'pointer',
                    transition: 'all 0.2s',
                  }}
                  onMouseEnter={(e) => {
                    e.currentTarget.style.backgroundColor = C.primaryDark;
                  }}
                  onMouseLeave={(e) => {
                    e.currentTarget.style.backgroundColor = C.primary;
                  }}
                >
                  Agregar stock
                </button>
              </div>
            </div>
          </div>
        </div>
      )}

      <style>{`
        @keyframes fadeIn {
          from { opacity: 0; }
          to { opacity: 1; }
        }
        @keyframes shimmer {
          0% { transform: translateX(-100%); }
          100% { transform: translateX(100%); }
        }
        .shimmer::after {
          content: '';
          position: absolute;
          top: 0;
          left: 0;
          width: 100%;
          height: 100%;
          background: linear-gradient(90deg, transparent, rgba(255,255,255,0.6), transparent);
          animation: shimmer 1.2s infinite;
        }
      `}</style>
    </div>
  );
}

// Estilos para celdas de tabla
const thStyle: React.CSSProperties = {
  textAlign: 'left',
  padding: '16px 20px',
  fontWeight: '600',
  color: C.textSecondary,
  fontSize: '0.75rem',
  textTransform: 'uppercase',
  letterSpacing: '0.5px',
  borderBottom: `1px solid ${C.borderLight}`,
};

const tdStyle: React.CSSProperties = {
  padding: '14px 20px',
  color: C.textPrimary,
  verticalAlign: 'middle',
};

const inputStyle: React.CSSProperties = {
  width: '100%',
  padding: '12px 14px',
  border: `1px solid ${C.borderLight}`,
  borderRadius: '18px',
  fontSize: '0.9rem',
  backgroundColor: C.surface,
  color: C.textPrimary,
  outline: 'none',
  transition: 'all 0.2s',
};