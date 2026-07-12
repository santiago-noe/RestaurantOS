import { useEffect, useState } from 'react';
import { creditosApi } from '../../services/api';
import type { Cliente } from '../../types';
import { CreditCard, DollarSign, Eye, Receipt, TrendingUp } from 'lucide-react';

// Paleta profesional consistente
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
  successBg: '#e6f4ea',
  info: '#3b82f6',
  infoBg: '#eff6ff',
  warning: '#d97706',
  warningBg: '#fffbeb',
};

interface EstadoCuenta {
  deuda_total: number;
  pedidos?: Array<{ id: number; fecha: string; total: number }>;
  pagos_realizados?: Array<{ id: number; fecha: string; metodo: string; monto: number }>;
}

export default function CreditosPage() {
  const [deudores, setDeudores] = useState<Cliente[]>([]);
  const [loading, setLoading] = useState(true);
  const [selectedCliente, setSelectedCliente] = useState<Cliente | null>(null);
  const [estadoCuenta, setEstadoCuenta] = useState<EstadoCuenta | null>(null);
  const [loadingEstado, setLoadingEstado] = useState(false);
  const [monto, setMonto] = useState('');
  const [metodo, setMetodo] = useState('efectivo');
  const [guardando, setGuardando] = useState(false);

  const cargarDeudores = async () => {
    setLoading(true);
    try {
      const { data } = await creditosApi.deudores();
      setDeudores(data.deudores || []);
    } catch (error) {
      console.error('Error al cargar deudores:', error);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    cargarDeudores();
  }, []);

  const verEstado = async (cliente: Cliente) => {
    setSelectedCliente(cliente);
    setLoadingEstado(true);
    try {
      const { data } = await creditosApi.estadoCuenta(cliente.id);
      setEstadoCuenta(data);
    } catch (error) {
      console.error('Error al cargar estado de cuenta:', error);
      setEstadoCuenta(null);
    } finally {
      setLoadingEstado(false);
    }
  };

  const registrarPago = async () => {
    if (!selectedCliente || !monto) return;
    const montoNum = parseFloat(monto);
    if (isNaN(montoNum) || montoNum <= 0) {
      alert('Ingresa un monto válido');
      return;
    }
    if (estadoCuenta && montoNum > estadoCuenta.deuda_total) {
      alert('El monto no puede superar la deuda total');
      return;
    }
    setGuardando(true);
    try {
      await creditosApi.registrarPago({
        cliente_id: selectedCliente.id,
        monto: montoNum,
        metodo,
      });
      setMonto('');
      // Recargar deudores y estado actualizado
      await cargarDeudores();
      const { data } = await creditosApi.estadoCuenta(selectedCliente.id);
      setEstadoCuenta(data);
    } catch (err: any) {
      alert(err.response?.data?.error || 'Error al registrar pago');
    } finally {
      setGuardando(false);
    }
  };

  const formatDate = (dateStr: string) => {
    const date = new Date(dateStr);
    return date.toLocaleDateString('es-PE', { day: '2-digit', month: 'short', year: 'numeric' });
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
      {/* Encabezado */}
      <div style={{ marginBottom: '32px' }}>
        <h1
          style={{
            fontSize: 'clamp(1.75rem, 5vw, 2rem)',
            fontWeight: '700',
            color: C.textPrimary,
            margin: '0 0 6px 0',
            letterSpacing: '-0.02em',
          }}
        >
          Créditos y Pagos
        </h1>
        <p
          style={{
            fontSize: '0.9rem',
            color: C.textSecondary,
            margin: 0,
          }}
        >
          {deudores.length} cliente{deudores.length !== 1 ? 's' : ''} con deuda activa
        </p>
      </div>

      {/* Grid de dos columnas */}
      <div
        style={{
          display: 'grid',
          gridTemplateColumns: 'repeat(auto-fit, minmax(320px, 1fr))',
          gap: '28px',
        }}
      >
        {/* Columna izquierda: Lista de deudores en tabla */}
        <div
          style={{
            backgroundColor: C.surface,
            borderRadius: '24px',
            border: `1px solid ${C.borderLight}`,
            overflow: 'hidden',
            boxShadow: '0 2px 8px rgba(0,0,0,0.02)',
          }}
        >
          <div
            style={{
              padding: '16px 20px',
              borderBottom: `1px solid ${C.borderLight}`,
              backgroundColor: C.bg,
            }}
          >
            <h2
              style={{
                fontSize: '0.85rem',
                fontWeight: '600',
                textTransform: 'uppercase',
                letterSpacing: '0.5px',
                color: C.textSecondary,
                margin: 0,
              }}
            >
              Clientes con deuda
            </h2>
          </div>
          <div style={{ overflowX: 'auto' }}>
            <table style={{ width: '100%', borderCollapse: 'collapse', fontSize: '0.85rem', minWidth: '400px' }}>
              <thead>
                <tr style={{ borderBottom: `1px solid ${C.borderLight}` }}>
                  <th style={thStyle}>Cliente</th>
                  <th style={thStyle}>Teléfono</th>
                  <th style={thStyle}>Deuda</th>
                  <th style={{ ...thStyle, textAlign: 'center' }}>Acción</th>
                </tr>
              </thead>
              <tbody>
                {loading ? (
                  [...Array(3)].map((_, i) => (
                    <tr key={i} style={{ borderBottom: `1px solid ${C.borderLight}` }}>
                      <td colSpan={4} style={{ padding: '12px 20px' }}>
                        <div
                          style={{
                            height: '32px',
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
                ) : deudores.length === 0 ? (
                  <tr>
                    <td colSpan={4} style={{ textAlign: 'center', padding: '48px 20px', color: C.textSecondary }}>
                      <CreditCard size={40} style={{ marginBottom: '12px', opacity: 0.4 }} />
                      <p>No hay deudas activas</p>
                    </td>
                  </tr>
                ) : (
                  deudores.map((cliente) => (
                    <tr
                      key={cliente.id}
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
                        <div>
                          <span style={{ fontWeight: '500', color: C.textPrimary }}>
                            {cliente.nombre} {cliente.apellido}
                          </span>
                        </div>
                      </td>
                      <td style={tdStyle}>
                        <span style={{ color: C.textSecondary, fontSize: '0.8rem' }}>
                          {cliente.telefono || '—'}
                        </span>
                      </td>
                      <td style={tdStyle}>
                        <span
                          style={{
                            fontWeight: '700',
                            color: C.alert,
                            backgroundColor: C.alertBg,
                            padding: '4px 10px',
                            borderRadius: '40px',
                            fontSize: '0.8rem',
                          }}
                        >
                          S/ {cliente.deuda_total?.toFixed(2) ?? '0.00'}
                        </span>
                      </td>
                      <td style={{ ...tdStyle, textAlign: 'center' }}>
                        <button
                          onClick={() => verEstado(cliente)}
                          style={{
                            backgroundColor: C.primaryLight,
                            border: 'none',
                            padding: '6px 14px',
                            borderRadius: '40px',
                            fontSize: '0.75rem',
                            fontWeight: '500',
                            color: C.primary,
                            cursor: 'pointer',
                            transition: 'all 0.2s',
                            display: 'inline-flex',
                            alignItems: 'center',
                            gap: '6px',
                          }}
                          onMouseEnter={(e) => {
                            e.currentTarget.style.backgroundColor = C.primary;
                            e.currentTarget.style.color = '#fff';
                          }}
                          onMouseLeave={(e) => {
                            e.currentTarget.style.backgroundColor = C.primaryLight;
                            e.currentTarget.style.color = C.primary;
                          }}
                        >
                          <Eye size={14} /> Ver
                        </button>
                      </td>
                    </tr>
                  ))
                )}
              </tbody>
            </table>
          </div>
        </div>

        {/* Columna derecha: Estado de cuenta + formulario de pago */}
        <div>
          {selectedCliente && estadoCuenta !== null ? (
            <div
              style={{
                backgroundColor: C.surface,
                borderRadius: '24px',
                border: `1px solid ${C.borderLight}`,
                overflow: 'hidden',
                boxShadow: '0 4px 12px rgba(0,0,0,0.02)',
              }}
            >
              {/* Cabecera del cliente */}
              <div
                style={{
                  padding: '20px 24px',
                  borderBottom: `1px solid ${C.borderLight}`,
                  backgroundColor: C.bg,
                }}
              >
                <h2
                  style={{
                    fontSize: '1.3rem',
                    fontWeight: '600',
                    color: C.textPrimary,
                    margin: '0 0 4px 0',
                  }}
                >
                  {selectedCliente.nombre} {selectedCliente.apellido}
                </h2>
                {selectedCliente.telefono && (
                  <p style={{ fontSize: '0.8rem', color: C.textSecondary, margin: 0 }}>
                    📞 {selectedCliente.telefono}
                  </p>
                )}
              </div>

              {/* Resumen de deuda */}
              <div style={{ padding: '20px 24px', borderBottom: `1px solid ${C.borderLight}` }}>
                <div
                  style={{
                    display: 'flex',
                    justifyContent: 'space-between',
                    alignItems: 'baseline',
                    flexWrap: 'wrap',
                    gap: '12px',
                  }}
                >
                  <div>
                    <p
                      style={{
                        fontSize: '0.7rem',
                        textTransform: 'uppercase',
                        fontWeight: '600',
                        color: C.textSecondary,
                        letterSpacing: '0.5px',
                        margin: '0 0 4px 0',
                      }}
                    >
                      Deuda total
                    </p>
                    <p
                      style={{
                        fontSize: '2rem',
                        fontWeight: '800',
                        color: C.alert,
                        margin: 0,
                        lineHeight: 1.2,
                      }}
                    >
                      S/ {estadoCuenta.deuda_total?.toFixed(2) ?? '0.00'}
                    </p>
                  </div>
                  <div
                    style={{
                      backgroundColor: C.infoBg,
                      padding: '8px 14px',
                      borderRadius: '60px',
                      display: 'flex',
                      alignItems: 'center',
                      gap: '8px',
                    }}
                  >
                    <TrendingUp size={16} color={C.info} />
                    <span style={{ fontSize: '0.75rem', fontWeight: '500', color: C.info }}>
                      {estadoCuenta.pedidos?.length || 0} pedido(s) pendiente(s)
                    </span>
                  </div>
                </div>
              </div>

              {/* Tabla de pedidos pendientes */}
              {estadoCuenta.pedidos && estadoCuenta.pedidos.length > 0 && (
                <div style={{ padding: '0 24px', borderBottom: `1px solid ${C.borderLight}` }}>
                  <p
                    style={{
                      fontSize: '0.7rem',
                      fontWeight: '600',
                      textTransform: 'uppercase',
                      letterSpacing: '0.5px',
                      color: C.textSecondary,
                      margin: '16px 0 12px 0',
                    }}
                  >
                    📋 Pedidos a crédito
                  </p>
                  <div style={{ overflowX: 'auto' }}>
                    <table style={{ width: '100%', fontSize: '0.8rem', borderCollapse: 'collapse' }}>
                      <thead>
                        <tr style={{ borderBottom: `1px solid ${C.borderLight}` }}>
                          <th style={{ ...thStyleInner, textAlign: 'left' }}>ID</th>
                          <th style={{ ...thStyleInner, textAlign: 'left' }}>Fecha</th>
                          <th style={{ ...thStyleInner, textAlign: 'right' }}>Total</th>
                        </tr>
                      </thead>
                      <tbody>
                        {estadoCuenta.pedidos.map((pedido) => (
                          <tr key={pedido.id} style={{ borderBottom: `1px solid ${C.borderLight}` }}>
                            <td style={{ ...tdStyleInner, fontWeight: '500' }}>#{pedido.id}</td>
                            <td style={tdStyleInner}>{formatDate(pedido.fecha)}</td>
                            <td style={{ ...tdStyleInner, textAlign: 'right', fontWeight: '500' }}>
                              S/ {pedido.total.toFixed(2)}
                            </td>
                          </tr>
                        ))}
                      </tbody>
                    </table>
                  </div>
                </div>
              )}

              {/* Tabla de pagos realizados */}
              {estadoCuenta.pagos_realizados && estadoCuenta.pagos_realizados.length > 0 && (
                <div style={{ padding: '0 24px', borderBottom: `1px solid ${C.borderLight}` }}>
                  <p
                    style={{
                      fontSize: '0.7rem',
                      fontWeight: '600',
                      textTransform: 'uppercase',
                      letterSpacing: '0.5px',
                      color: C.textSecondary,
                      margin: '16px 0 12px 0',
                    }}
                  >
                    💰 Pagos registrados
                  </p>
                  <div style={{ overflowX: 'auto' }}>
                    <table style={{ width: '100%', fontSize: '0.8rem', borderCollapse: 'collapse' }}>
                      <thead>
                        <tr style={{ borderBottom: `1px solid ${C.borderLight}` }}>
                          <th style={{ ...thStyleInner, textAlign: 'left' }}>Fecha</th>
                          <th style={{ ...thStyleInner, textAlign: 'left' }}>Método</th>
                          <th style={{ ...thStyleInner, textAlign: 'right' }}>Monto</th>
                        </tr>
                      </thead>
                      <tbody>
                        {estadoCuenta.pagos_realizados.map((pago) => (
                          <tr key={pago.id} style={{ borderBottom: `1px solid ${C.borderLight}` }}>
                            <td style={tdStyleInner}>{formatDate(pago.fecha)}</td>
                            <td style={tdStyleInner}>
                              <span
                                style={{
                                  backgroundColor: C.successBg,
                                  color: C.success,
                                  padding: '2px 10px',
                                  borderRadius: '40px',
                                  fontSize: '0.7rem',
                                  fontWeight: '500',
                                }}
                              >
                                {pago.metodo}
                              </span>
                            </td>
                            <td style={{ ...tdStyleInner, textAlign: 'right', color: C.success, fontWeight: '500' }}>
                              + S/ {pago.monto.toFixed(2)}
                            </td>
                          </tr>
                        ))}
                      </tbody>
                    </table>
                  </div>
                </div>
              )}

              {/* Formulario de pago */}
              <div style={{ padding: '20px 24px' }}>
                <p
                  style={{
                    fontSize: '0.7rem',
                    fontWeight: '600',
                    textTransform: 'uppercase',
                    letterSpacing: '0.5px',
                    color: C.textSecondary,
                    margin: '0 0 12px 0',
                  }}
                >
                  Registrar nuevo abono
                </p>
                <div style={{ display: 'flex', flexDirection: 'column', gap: '14px' }}>
                  <div style={{ display: 'flex', gap: '12px', flexWrap: 'wrap' }}>
                    <div style={{ flex: 2, minWidth: '140px', position: 'relative' }}>
                      <DollarSign
                        size={16}
                        style={{
                          position: 'absolute',
                          left: '12px',
                          top: '50%',
                          transform: 'translateY(-50%)',
                          color: C.textSecondary,
                        }}
                      />
                      <input
                        type="number"
                        step="0.01"
                        min="0.01"
                        max={estadoCuenta.deuda_total}
                        value={monto}
                        onChange={(e) => setMonto(e.target.value)}
                        placeholder="Monto"
                        style={{
                          width: '100%',
                          padding: '10px 12px 10px 36px',
                          border: `1px solid ${C.borderLight}`,
                          borderRadius: '20px',
                          fontSize: '0.9rem',
                          outline: 'none',
                          transition: 'all 0.2s',
                        }}
                        onFocus={(e) => {
                          e.currentTarget.style.borderColor = C.primary;
                          e.currentTarget.style.boxShadow = `0 0 0 3px ${C.primary}20`;
                        }}
                        onBlur={(e) => {
                          e.currentTarget.style.borderColor = C.borderLight;
                          e.currentTarget.style.boxShadow = 'none';
                        }}
                      />
                    </div>
                    <select
                      value={metodo}
                      onChange={(e) => setMetodo(e.target.value)}
                      style={{
                        flex: 1,
                        padding: '10px 12px',
                        border: `1px solid ${C.borderLight}`,
                        borderRadius: '20px',
                        fontSize: '0.85rem',
                        backgroundColor: C.surface,
                        cursor: 'pointer',
                        outline: 'none',
                      }}
                    >
                      {['efectivo', 'yape', 'plin', 'transferencia'].map((m) => (
                        <option key={m} value={m}>
                          {m.charAt(0).toUpperCase() + m.slice(1)}
                        </option>
                      ))}
                    </select>
                  </div>
                  <button
                    onClick={registrarPago}
                    disabled={!monto || guardando}
                    style={{
                      width: '100%',
                      backgroundColor: C.primary,
                      color: '#fff',
                      border: 'none',
                      padding: '12px',
                      borderRadius: '40px',
                      fontSize: '0.9rem',
                      fontWeight: '500',
                      cursor: 'pointer',
                      transition: 'all 0.2s',
                      opacity: !monto || guardando ? 0.6 : 1,
                    }}
                    onMouseEnter={(e) => {
                      if (!guardando && monto)
                        e.currentTarget.style.backgroundColor = C.primaryDark;
                    }}
                    onMouseLeave={(e) => {
                      e.currentTarget.style.backgroundColor = C.primary;
                    }}
                  >
                    {guardando ? 'Registrando...' : 'Registrar pago'}
                  </button>
                </div>
              </div>
            </div>
          ) : loadingEstado ? (
            <div
              style={{
                backgroundColor: C.surface,
                borderRadius: '24px',
                border: `1px solid ${C.borderLight}`,
                padding: '40px',
                textAlign: 'center',
              }}
            >
              <div
                style={{
                  width: '48px',
                  height: '48px',
                  margin: '0 auto 16px',
                  borderRadius: '50%',
                  backgroundColor: C.borderLight,
                  position: 'relative',
                  overflow: 'hidden',
                }}
                className="shimmer"
              />
              <div
                style={{
                  height: '20px',
                  width: '60%',
                  margin: '0 auto',
                  backgroundColor: C.borderLight,
                  borderRadius: '8px',
                  position: 'relative',
                  overflow: 'hidden',
                }}
                className="shimmer"
              />
            </div>
          ) : (
            <div
              style={{
                backgroundColor: C.surface,
                borderRadius: '24px',
                border: `2px dashed ${C.borderLight}`,
                padding: '48px 24px',
                textAlign: 'center',
                color: C.textSecondary,
              }}
            >
              <Receipt size={48} style={{ marginBottom: '16px', opacity: 0.5 }} />
              <p style={{ margin: 0, fontSize: '0.9rem' }}>
                Selecciona un cliente para ver su estado de cuenta
              </p>
            </div>
          )}
        </div>
      </div>

      <style>{`
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

// Estilos para tablas
const thStyle: React.CSSProperties = {
  textAlign: 'left',
  padding: '14px 16px',
  fontWeight: '600',
  fontSize: '0.7rem',
  textTransform: 'uppercase',
  letterSpacing: '0.5px',
  color: C.textSecondary,
  borderBottom: `1px solid ${C.borderLight}`,
};

const tdStyle: React.CSSProperties = {
  padding: '12px 16px',
  color: C.textPrimary,
  verticalAlign: 'middle',
};

// Estilos para tablas internas (pedidos/pagos)
const thStyleInner: React.CSSProperties = {
  padding: '8px 0',
  fontWeight: '500',
  fontSize: '0.7rem',
  color: C.textSecondary,
  borderBottom: `1px solid ${C.borderLight}`,
};

const tdStyleInner: React.CSSProperties = {
  padding: '10px 0',
  borderBottom: `1px solid ${C.borderLight}`,
  verticalAlign: 'middle',
};