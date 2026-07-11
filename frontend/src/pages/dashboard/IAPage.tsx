import { useEffect, useState } from 'react';
import { Brain, AlertTriangle, TrendingUp, RefreshCw, Info, Coffee, Utensils, Wine } from 'lucide-react';

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
  warning: '#d97706',
  warningBg: '#fffbeb',
  info: '#3b82f6',
  infoBg: '#eff6ff',
  purple: '#8b5cf6',
  purpleLight: '#f3e8ff',
};

// Datos de demostración (mismos que original)
const PREDICCIONES_DEMO = [
  { tipo_comida: 'desayuno', porciones_estimadas: 22, base_calculo: 'promedio últimos 8 días similares', icono: Coffee },
  { tipo_comida: 'almuerzo', porciones_estimadas: 48, base_calculo: 'promedio últimos 8 días similares', icono: Utensils },
  { tipo_comida: 'bebida', porciones_estimadas: 15, base_calculo: 'promedio últimos 8 días similares', icono: Wine },
];

const ALERTAS_DEMO = [
  { id: 1, tipo: 'deuda_alta', mensaje: 'Constructora Norte tiene una deuda de S/ 245.00 (umbral: S/ 200)', severidad: 'alta' },
  { id: 2, tipo: 'stock_bajo', mensaje: 'Aceite: 0.5 litros restantes (mínimo: 2 litros)', severidad: 'media' },
  { id: 3, tipo: 'stock_bajo', mensaje: 'Pollo: 1.0 kg restantes (mínimo: 3 kg)', severidad: 'media' },
];

const severidadConfig: Record<string, { label: string; color: string; bg: string }> = {
  alta: { label: 'Alta prioridad', color: C.alert, bg: C.alertBg },
  media: { label: 'Prioridad media', color: C.warning, bg: C.warningBg },
};

export default function IAPage() {
  const [loading, setLoading] = useState(false);
  const manana = new Date();
  manana.setDate(manana.getDate() + 1);
  const fechaFormateada = manana.toLocaleDateString('es-PE', {
    weekday: 'long',
    year: 'numeric',
    month: 'long',
    day: 'numeric',
  });

  const refresh = () => {
    setLoading(true);
    setTimeout(() => setLoading(false), 800);
  };

  useEffect(() => {
    refresh();
  }, []);

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
          <div style={{ display: 'flex', alignItems: 'center', gap: '10px', marginBottom: '6px' }}>
            <Brain size={28} color={C.purple} />
            <h1
              style={{
                fontSize: 'clamp(1.75rem, 5vw, 2rem)',
                fontWeight: '700',
                color: C.textPrimary,
                margin: 0,
                letterSpacing: '-0.02em',
              }}
            >
              IA y Alertas
            </h1>
          </div>
          <p
            style={{
              fontSize: '0.9rem',
              color: C.textSecondary,
              margin: 0,
            }}
          >
            Predicciones para {fechaFormateada}
          </p>
        </div>
        <button
          onClick={refresh}
          disabled={loading}
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
            opacity: loading ? 0.6 : 1,
          }}
          onMouseEnter={(e) => {
            if (!loading) {
              e.currentTarget.style.backgroundColor = C.surface;
              e.currentTarget.style.borderColor = C.purple;
              e.currentTarget.style.color = C.purple;
            }
          }}
          onMouseLeave={(e) => {
            if (!loading) {
              e.currentTarget.style.backgroundColor = 'transparent';
              e.currentTarget.style.borderColor = C.borderLight;
              e.currentTarget.style.color = C.textSecondary;
            }
          }}
        >
          <RefreshCw size={16} className={loading ? 'spin' : ''} /> Actualizar
        </button>
      </div>

      {/* Sección de predicciones */}
      <div style={{ marginBottom: '40px' }}>
        <div
          style={{
            display: 'flex',
            alignItems: 'center',
            gap: '8px',
            marginBottom: '16px',
          }}
        >
          <TrendingUp size={18} color={C.purple} />
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
            Predicción de demanda
          </h2>
        </div>
        {loading ? (
          <div
            style={{
              display: 'grid',
              gridTemplateColumns: 'repeat(auto-fit, minmax(240px, 1fr))',
              gap: '20px',
            }}
          >
            {[...Array(3)].map((_, i) => (
              <div
                key={i}
                style={{
                  height: '180px',
                  backgroundColor: C.surface,
                  borderRadius: '24px',
                  border: `1px solid ${C.borderLight}`,
                  position: 'relative',
                  overflow: 'hidden',
                }}
                className="shimmer"
              />
            ))}
          </div>
        ) : (
          <div
            style={{
              display: 'grid',
              gridTemplateColumns: 'repeat(auto-fit, minmax(240px, 1fr))',
              gap: '20px',
            }}
          >
            {PREDICCIONES_DEMO.map((pred) => {
              const Icono = pred.icono;
              return (
                <div
                  key={pred.tipo_comida}
                  style={{
                    backgroundColor: C.surface,
                    borderRadius: '24px',
                    border: `1px solid ${C.borderLight}`,
                    padding: '20px',
                    transition: 'all 0.25s ease',
                    cursor: 'pointer',
                  }}
                  onMouseEnter={(e) => {
                    e.currentTarget.style.transform = 'translateY(-4px)';
                    e.currentTarget.style.boxShadow = '0 12px 24px -12px rgba(0,0,0,0.12)';
                    e.currentTarget.style.borderColor = C.purple;
                  }}
                  onMouseLeave={(e) => {
                    e.currentTarget.style.transform = 'translateY(0)';
                    e.currentTarget.style.boxShadow = 'none';
                    e.currentTarget.style.borderColor = C.borderLight;
                  }}
                >
                  <div
                    style={{
                      width: '48px',
                      height: '48px',
                      borderRadius: '20px',
                      backgroundColor: C.purpleLight,
                      display: 'flex',
                      alignItems: 'center',
                      justifyContent: 'center',
                      marginBottom: '16px',
                    }}
                  >
                    <Icono size={24} color={C.purple} />
                  </div>
                  <p
                    style={{
                      fontSize: '2.2rem',
                      fontWeight: '800',
                      color: C.purple,
                      margin: '0 0 4px 0',
                      lineHeight: 1.2,
                    }}
                  >
                    {pred.porciones_estimadas}
                  </p>
                  <p
                    style={{
                      fontWeight: '600',
                      color: C.textPrimary,
                      margin: '0 0 6px 0',
                      fontSize: '1rem',
                      textTransform: 'capitalize',
                    }}
                  >
                    {pred.tipo_comida}s
                  </p>
                  <p
                    style={{
                      fontSize: '0.7rem',
                      color: C.textSecondary,
                      margin: 0,
                    }}
                  >
                    {pred.base_calculo}
                  </p>
                </div>
              );
            })}
          </div>
        )}
      </div>

      {/* Sección de alertas activas */}
      <div style={{ marginBottom: '28px' }}>
        <div
          style={{
            display: 'flex',
            alignItems: 'center',
            gap: '8px',
            marginBottom: '16px',
          }}
        >
          <AlertTriangle size={18} color={C.warning} />
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
            Alertas activas ({ALERTAS_DEMO.length})
          </h2>
        </div>

        {loading ? (
          <div
            style={{
              backgroundColor: C.surface,
              borderRadius: '20px',
              border: `1px solid ${C.borderLight}`,
              overflow: 'hidden',
            }}
          >
            <div style={{ padding: '16px 20px' }}>
              <div
                style={{
                  height: '50px',
                  backgroundColor: C.bg,
                  borderRadius: '12px',
                  position: 'relative',
                  overflow: 'hidden',
                }}
                className="shimmer"
              />
            </div>
          </div>
        ) : ALERTAS_DEMO.length === 0 ? (
          <div
            style={{
              backgroundColor: C.surface,
              borderRadius: '20px',
              border: `1px solid ${C.borderLight}`,
              padding: '40px',
              textAlign: 'center',
              color: C.textSecondary,
            }}
          >
            <CheckCircle size={40} style={{ marginBottom: '12px', opacity: 0.5 }} />
            <p>No hay alertas activas. Todo está en orden.</p>
          </div>
        ) : (
          <div
            style={{
              backgroundColor: C.surface,
              borderRadius: '20px',
              border: `1px solid ${C.borderLight}`,
              overflow: 'hidden',
              boxShadow: '0 2px 8px rgba(0,0,0,0.02)',
            }}
          >
            <div style={{ overflowX: 'auto' }}>
              <table
                style={{
                  width: '100%',
                  borderCollapse: 'collapse',
                  fontSize: '0.85rem',
                  minWidth: '500px',
                }}
              >
                <thead>
                  <tr style={{ borderBottom: `1px solid ${C.borderLight}` }}>
                    <th style={thStyle}>Tipo</th>
                    <th style={thStyle}>Mensaje</th>
                    <th style={{ ...thStyle, textAlign: 'center' }}>Prioridad</th>
                  </tr>
                </thead>
                <tbody>
                  {ALERTAS_DEMO.map((alerta) => {
                    const config = severidadConfig[alerta.severidad];
                    return (
                      <tr
                        key={alerta.id}
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
                          <span
                            style={{
                              backgroundColor:
                                alerta.tipo === 'deuda_alta' ? C.alertBg : C.warningBg,
                              color: alerta.tipo === 'deuda_alta' ? C.alert : C.warning,
                              padding: '4px 12px',
                              borderRadius: '40px',
                              fontSize: '0.7rem',
                              fontWeight: '500',
                              textTransform: 'capitalize',
                            }}
                          >
                            {alerta.tipo === 'deuda_alta' ? 'Deuda alta' : 'Stock bajo'}
                          </span>
                        </td>
                        <td style={tdStyle}>
                          <span style={{ color: C.textPrimary }}>{alerta.mensaje}</span>
                        </td>
                        <td style={{ ...tdStyle, textAlign: 'center' }}>
                          <span
                            style={{
                              backgroundColor: config.bg,
                              color: config.color,
                              padding: '4px 10px',
                              borderRadius: '40px',
                              fontSize: '0.7rem',
                              fontWeight: '600',
                            }}
                          >
                            {config.label}
                          </span>
                        </td>
                      </tr>
                    );
                  })}
                </tbody>
              </table>
            </div>
          </div>
        )}
      </div>

      {/* Nota informativa */}
      <div
        style={{
          backgroundColor: C.purpleLight,
          border: `1px solid ${C.purple}30`,
          borderRadius: '20px',
          padding: '16px 20px',
          display: 'flex',
          alignItems: 'flex-start',
          gap: '12px',
        }}
      >
        <Info size={18} color={C.purple} style={{ marginTop: '2px', flexShrink: 0 }} />
        <div>
          <p style={{ fontSize: '0.8rem', color: C.purple, margin: 0 }}>
            <strong>Nota:</strong> Las predicciones se calculan automáticamente con el historial de pedidos de los últimos 30 días,
            agrupado por día de la semana. Sin historial suficiente, la predicción mostrará valores bajos.
          </p>
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
        @keyframes spin {
          from { transform: rotate(0deg); }
          to { transform: rotate(360deg); }
        }
        .spin {
          animation: spin 0.8s linear infinite;
        }
      `}</style>
    </div>
  );
}

// Estilos para la tabla de alertas
const thStyle: React.CSSProperties = {
  textAlign: 'left',
  padding: '14px 20px',
  fontWeight: '600',
  fontSize: '0.7rem',
  textTransform: 'uppercase',
  letterSpacing: '0.5px',
  color: C.textSecondary,
  borderBottom: `1px solid ${C.borderLight}`,
};

const tdStyle: React.CSSProperties = {
  padding: '12px 20px',
  color: C.textPrimary,
  verticalAlign: 'middle',
};

// Necesitamos importar CheckCircle para el estado vacío
import { CheckCircle } from 'lucide-react';