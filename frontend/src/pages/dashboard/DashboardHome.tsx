import { useEffect, useState } from 'react';
import { useAuth } from '../../context/AuthContext';
import { clientesApi, pedidosApi, inventarioApi, creditosApi } from '../../services/api';
import { TrendingUp, Users, ShoppingBag, AlertTriangle, Calendar, Clock } from 'lucide-react';

// Paleta de colores mejorada
const C = {
  bg: '#fefaf5',
  surface: '#ffffff',
  primary: '#b4532b',
  primaryDark: '#8b3a1a',
  secondary: '#3f6b5c',
  gold: '#d4af37',
  goldDark: '#b38f2c',
  textPrimary: '#2c2c2a',
  textSecondary: '#6b5e55',
  borderLight: '#efe5dc',
  alert: '#e53e3e',
  alertBg: '#fff5f5',
};

interface Stats {
  totalClientes: number;
  pedidosHoy: number;
  deudores: number;
  alertasStock: number;
}

export default function DashboardHome() {
  const { user } = useAuth();
  const [stats, setStats] = useState<Stats>({
    totalClientes: 0,
    pedidosHoy: 0,
    deudores: 0,
    alertasStock: 0,
  });
  const [loading, setLoading] = useState(true);

  // Saludo personalizado según la hora
  const getGreeting = () => {
    const hour = new Date().getHours();
    if (hour < 12) return 'Buenos días';
    if (hour < 18) return 'Buenas tardes';
    return 'Buenas noches';
  };

  useEffect(() => {
    const cargar = async () => {
      try {
        const [clientes, pedidos, deudores, alertas] = await Promise.allSettled([
          clientesApi.listar(),
          pedidosApi.listar(),
          creditosApi.deudores(),
          inventarioApi.alertas(),
        ]);
        setStats({
          totalClientes: clientes.status === 'fulfilled' ? clientes.value.data.total : 0,
          pedidosHoy: pedidos.status === 'fulfilled' ? pedidos.value.data.total : 0,
          deudores: deudores.status === 'fulfilled' ? deudores.value.data.total : 0,
          alertasStock: alertas.status === 'fulfilled' ? alertas.value.data.alertas?.length : 0,
        });
      } finally {
        setLoading(false);
      }
    };
    cargar();
  }, []);

  const cards = [
    {
      label: 'Clientes activos',
      value: stats.totalClientes,
      icon: Users,
      gradient: `linear-gradient(135deg, ${C.primary}15, ${C.primary}05)`,
      iconColor: C.primary,
      borderColor: C.primary,
    },
    {
      label: 'Pedidos totales',
      value: stats.pedidosHoy,
      icon: ShoppingBag,
      gradient: `linear-gradient(135deg, ${C.gold}15, ${C.gold}05)`,
      iconColor: C.goldDark,
      borderColor: C.gold,
    },
    {
      label: 'Clientes con deuda',
      value: stats.deudores,
      icon: TrendingUp,
      gradient: `linear-gradient(135deg, ${C.secondary}15, ${C.secondary}05)`,
      iconColor: C.secondary,
      borderColor: C.secondary,
    },
    {
      label: 'Alertas de stock',
      value: stats.alertasStock,
      icon: AlertTriangle,
      gradient: `linear-gradient(135deg, ${C.alert}15, ${C.alert}05)`,
      iconColor: C.alert,
      borderColor: C.alert,
    },
  ];

  const today = new Date();
  const formattedDate = today.toLocaleDateString('es-PE', {
    weekday: 'long',
    year: 'numeric',
    month: 'long',
    day: 'numeric',
  });

  // Formateo capitalizado
  const capitalizedDate = formattedDate.charAt(0).toUpperCase() + formattedDate.slice(1);

  return (
    <div
      style={{
        fontFamily:
          '"Inter", system-ui, -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, sans-serif',
        maxWidth: '1400px',
        margin: '0 auto',
        padding: '24px 20px',
        backgroundColor: C.bg,
        minHeight: '100vh',
      }}
    >
      {/* Header mejorado */}
      <div
        style={{
          display: 'flex',
          flexWrap: 'wrap',
          justifyContent: 'space-between',
          alignItems: 'flex-end',
          marginBottom: '32px',
          gap: '16px',
        }}
      >
        <div>
          <h1
            style={{
              fontSize: 'clamp(1.75rem, 5vw, 2.25rem)',
              fontWeight: '700',
              color: C.textPrimary,
              margin: '0 0 8px 0',
              letterSpacing: '-0.02em',
            }}
          >
            {getGreeting()}, {user?.nombre?.split(' ')[0] || 'Administrador'} 👋
          </h1>
          <div
            style={{
              display: 'flex',
              alignItems: 'center',
              gap: '8px',
              color: C.textSecondary,
              fontSize: '0.9rem',
            }}
          >
            <Calendar size={18} />
            <span>{capitalizedDate}</span>
            <span style={{ width: '4px', height: '4px', backgroundColor: C.borderLight, borderRadius: '50%' }} />
            <Clock size={18} />
            <span>{today.toLocaleTimeString('es-PE', { hour: '2-digit', minute: '2-digit' })}</span>
          </div>
        </div>
        {/* Badge opcional de versión o empresa */}
        <div
          style={{
            backgroundColor: C.surface,
            padding: '8px 16px',
            borderRadius: '40px',
            fontSize: '0.8rem',
            color: C.primary,
            fontWeight: '500',
            boxShadow: '0 1px 2px rgba(0,0,0,0.05)',
            border: `1px solid ${C.borderLight}`,
          }}
        >
          Vilca Suyo • Dashboard
        </div>
      </div>

      {/* Stats Cards Grid con Skeleton Shimmer */}
      {loading ? (
        <div
          style={{
            display: 'grid',
            gridTemplateColumns: 'repeat(auto-fit, minmax(280px, 1fr))',
            gap: '20px',
            marginBottom: '40px',
          }}
        >
          {[...Array(4)].map((_, i) => (
            <div
              key={i}
              style={{
                height: '150px',
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
      ) : (
        <div
          style={{
            display: 'grid',
            gridTemplateColumns: 'repeat(auto-fit, minmax(280px, 1fr))',
            gap: '20px',
            marginBottom: '40px',
          }}
        >
          {cards.map(({ label, value, icon: Icon, gradient, iconColor, borderColor }) => (
            <div
              key={label}
              style={{
                backgroundColor: C.surface,
                borderRadius: '20px',
                border: `1px solid ${C.borderLight}`,
                padding: '24px',
                transition: 'all 0.25s ease',
                cursor: 'pointer',
                position: 'relative',
                overflow: 'hidden',
              }}
              onMouseEnter={(e) => {
                e.currentTarget.style.transform = 'translateY(-4px)';
                e.currentTarget.style.boxShadow = '0 20px 25px -12px rgba(0,0,0,0.1)';
                e.currentTarget.style.borderColor = borderColor;
              }}
              onMouseLeave={(e) => {
                e.currentTarget.style.transform = 'translateY(0)';
                e.currentTarget.style.boxShadow = 'none';
                e.currentTarget.style.borderColor = C.borderLight;
              }}
            >
              <div
                style={{
                  display: 'flex',
                  justifyContent: 'space-between',
                  alignItems: 'center',
                  marginBottom: '20px',
                }}
              >
                <div
                  style={{
                    width: '48px',
                    height: '48px',
                    borderRadius: '16px',
                    background: gradient,
                    display: 'flex',
                    alignItems: 'center',
                    justifyContent: 'center',
                  }}
                >
                  <Icon size={24} color={iconColor} strokeWidth={1.8} />
                </div>
                {/* Mini indicador decorativo */}
                <div
                  style={{
                    width: '8px',
                    height: '8px',
                    borderRadius: '50%',
                    backgroundColor: borderColor,
                    opacity: 0.4,
                  }}
                />
              </div>
              <p
                style={{
                  fontSize: '2.5rem',
                  fontWeight: '800',
                  color: C.textPrimary,
                  margin: '0 0 6px 0',
                  lineHeight: 1.2,
                  letterSpacing: '-0.02em',
                }}
              >
                {value.toLocaleString()}
              </p>
              <p
                style={{
                  fontSize: '0.9rem',
                  color: C.textSecondary,
                  margin: 0,
                  fontWeight: '500',
                }}
              >
                {label}
              </p>
            </div>
          ))}
        </div>
      )}

      {/* Banner rediseñado con dos columnas y mejor jerarquía */}
      <div
        style={{
          borderRadius: '24px',
          background: `linear-gradient(135deg, ${C.primary} 0%, ${C.primaryDark} 100%)`,
          color: '#ffffff',
          overflow: 'hidden',
          position: 'relative',
        }}
      >
        {/* Patrón de fondo sutil */}
        <div
          style={{
            position: 'absolute',
            top: 0,
            right: 0,
            bottom: 0,
            left: 0,
            backgroundImage: `radial-gradient(circle at 20% 40%, rgba(255,255,255,0.08) 2%, transparent 2.5%)`,
            backgroundSize: '24px 24px',
            pointerEvents: 'none',
          }}
        />
        <div
          style={{
            padding: '28px 32px',
            display: 'flex',
            flexWrap: 'wrap',
            justifyContent: 'space-between',
            alignItems: 'center',
            gap: '24px',
            position: 'relative',
            zIndex: 1,
          }}
        >
          <div style={{ flex: 2, minWidth: '200px' }}>
            <h2
              style={{
                fontFamily: '"Cal Sans", "Inter", system-ui, sans-serif',
                fontSize: 'clamp(1.4rem, 4vw, 1.8rem)',
                fontWeight: '600',
                margin: '0 0 8px 0',
                letterSpacing: '-0.01em',
              }}
            >
              Vilca Suyo 🍽️
            </h2>
            <p
              style={{
                fontSize: '0.95rem',
                opacity: 0.92,
                margin: 0,
                lineHeight: 1.5,
                maxWidth: '500px',
              }}
            >
              Panel central para gestionar clientes, pedidos, inventario y créditos.
              Todo lo que necesitas en un solo lugar.
            </p>
          </div>
          <div
            style={{
              display: 'flex',
              flexWrap: 'wrap',
              gap: '12px',
              flex: 1,
              justifyContent: 'flex-end',
            }}
          >
            {stats.alertasStock > 0 && (
              <div
                style={{
                  backgroundColor: 'rgba(255,255,255,0.18)',
                  backdropFilter: 'blur(2px)',
                  borderRadius: '60px',
                  padding: '10px 20px',
                  display: 'flex',
                  alignItems: 'center',
                  gap: '10px',
                  fontWeight: '500',
                  fontSize: '0.85rem',
                  border: '1px solid rgba(255,255,255,0.25)',
                }}
              >
                <AlertTriangle size={18} />
                <span>
                  <strong>{stats.alertasStock}</strong> alertas de stock
                </span>
              </div>
            )}
            {stats.deudores > 0 && (
              <div
                style={{
                  backgroundColor: 'rgba(255,255,255,0.18)',
                  backdropFilter: 'blur(2px)',
                  borderRadius: '60px',
                  padding: '10px 20px',
                  display: 'flex',
                  alignItems: 'center',
                  gap: '10px',
                  fontWeight: '500',
                  fontSize: '0.85rem',
                  border: '1px solid rgba(255,255,255,0.25)',
                }}
              >
                <TrendingUp size={18} />
                <span>
                  <strong>{stats.deudores}</strong> deudores activos
                </span>
              </div>
            )}
            {stats.alertasStock === 0 && stats.deudores === 0 && (
              <div
                style={{
                  backgroundColor: 'rgba(255,255,255,0.18)',
                  backdropFilter: 'blur(2px)',
                  borderRadius: '60px',
                  padding: '10px 20px',
                  fontSize: '0.85rem',
                }}
              >
                Todo en orden 🎉
              </div>
            )}
          </div>
        </div>
      </div>

      {/* Animaciones globales */}
      <style>{`
        @keyframes shimmer {
          0% {
            transform: translateX(-100%);
          }
          100% {
            transform: translateX(100%);
          }
        }
        .shimmer::after {
          content: '';
          position: absolute;
          top: 0;
          left: 0;
          width: 100%;
          height: 100%;
          background: linear-gradient(
            90deg,
            transparent,
            rgba(255, 255, 255, 0.5),
            transparent
          );
          animation: shimmer 1.2s infinite;
        }
        @media (max-width: 640px) {
          .shimmer::after {
            animation-duration: 1.5s;
          }
        }
      `}</style>
    </div>
  );
}