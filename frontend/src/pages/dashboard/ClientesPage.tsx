import { useEffect, useState } from 'react';
import { clientesApi } from '../../services/api';
import type { Cliente } from '../../types';
import { useAuth } from '../../context/AuthContext';
import { Plus, Search, Building2, User, X, Phone, Mail, UserCircle } from 'lucide-react';

// Paleta profesional consistente con el dashboard
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
};

export default function ClientesPage() {
  const { isAdmin } = useAuth();
  const [clientes, setClientes] = useState<Cliente[]>([]);
  const [total, setTotal] = useState(0);
  const [busqueda, setBusqueda] = useState('');
  const [loading, setLoading] = useState(true);
  const [showForm, setShowForm] = useState(false);
  const [form, setForm] = useState({
    nombre: '',
    apellido: '',
    tipo: 'individual',
    telefono: '',
    email: '',
  });
  const [guardando, setGuardando] = useState(false);

  const cargar = async () => {
    setLoading(true);
    try {
      const { data } = await clientesApi.listar();
      setClientes(data.data);
      setTotal(data.total);
    } catch (error) {
      console.error('Error al cargar clientes:', error);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    cargar();
  }, []);

  const filtrados = clientes.filter((c) =>
    `${c.nombre} ${c.apellido}`.toLowerCase().includes(busqueda.toLowerCase())
  );

  const handleCrear = async (e: React.FormEvent) => {
    e.preventDefault();
    setGuardando(true);
    try {
      await clientesApi.crear(form);
      setShowForm(false);
      setForm({ nombre: '', apellido: '', tipo: 'individual', telefono: '', email: '' });
      cargar();
    } catch (err: any) {
      alert(err.response?.data?.error || 'Error al crear cliente');
    } finally {
      setGuardando(false);
    }
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
            Clientes
          </h1>
          <p
            style={{
              fontSize: '0.9rem',
              color: C.textSecondary,
              margin: 0,
            }}
          >
            {total} clientes registrados • gestión completa
          </p>
        </div>
        {isAdmin && (
          <button
            onClick={() => setShowForm(true)}
            style={{
              display: 'flex',
              alignItems: 'center',
              gap: '8px',
              backgroundColor: C.primary,
              color: '#fff',
              border: 'none',
              padding: '10px 20px',
              borderRadius: '40px',
              fontSize: '0.9rem',
              fontWeight: '500',
              cursor: 'pointer',
              transition: 'all 0.2s ease',
              boxShadow: '0 2px 6px rgba(0,0,0,0.05)',
            }}
            onMouseEnter={(e) => {
              e.currentTarget.style.backgroundColor = C.primaryDark;
              e.currentTarget.style.transform = 'translateY(-1px)';
            }}
            onMouseLeave={(e) => {
              e.currentTarget.style.backgroundColor = C.primary;
              e.currentTarget.style.transform = 'translateY(0)';
            }}
          >
            <Plus size={18} /> Nuevo cliente
          </button>
        )}
      </div>

      {/* Buscador */}
      <div
        style={{
          position: 'relative',
          marginBottom: '28px',
        }}
      >
        <Search
          size={18}
          style={{
            position: 'absolute',
            left: '16px',
            top: '50%',
            transform: 'translateY(-50%)',
            color: C.textSecondary,
          }}
        />
        <input
          value={busqueda}
          onChange={(e) => setBusqueda(e.target.value)}
          placeholder="Buscar por nombre o apellido..."
          style={{
            width: '100%',
            padding: '12px 16px 12px 44px',
            border: `1px solid ${C.borderLight}`,
            borderRadius: '60px',
            fontSize: '0.9rem',
            backgroundColor: C.surface,
            color: C.textPrimary,
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

      {/* Modal de nuevo cliente */}
      {showForm && (
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
            if (e.target === e.currentTarget) setShowForm(false);
          }}
        >
          <div
            style={{
              backgroundColor: C.surface,
              borderRadius: '28px',
              maxWidth: '500px',
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
              <h2
                style={{
                  fontSize: '1.4rem',
                  fontWeight: '600',
                  color: C.textPrimary,
                  margin: 0,
                }}
              >
                Nuevo cliente
              </h2>
              <button
                onClick={() => setShowForm(false)}
                style={{
                  background: 'none',
                  border: 'none',
                  cursor: 'pointer',
                  color: C.textSecondary,
                  padding: '4px',
                  borderRadius: '50%',
                  display: 'flex',
                }}
              >
                <X size={20} />
              </button>
            </div>
            <form onSubmit={handleCrear} style={{ padding: '24px' }}>
              <div style={{ display: 'grid', gap: '20px' }}>
                <div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr', gap: '16px' }}>
                  <div>
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
                      Nombre *
                    </label>
                    <input
                      required
                      value={form.nombre}
                      onChange={(e) => setForm((p) => ({ ...p, nombre: e.target.value }))}
                      style={inputStyle}
                    />
                  </div>
                  <div>
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
                      Apellido
                    </label>
                    <input
                      value={form.apellido}
                      onChange={(e) => setForm((p) => ({ ...p, apellido: e.target.value }))}
                      style={inputStyle}
                    />
                  </div>
                </div>
                <div>
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
                    Tipo *
                  </label>
                  <select
                    value={form.tipo}
                    onChange={(e) => setForm((p) => ({ ...p, tipo: e.target.value }))}
                    style={inputStyle}
                  >
                    <option value="individual">Individual</option>
                    <option value="empresa">Empresa</option>
                  </select>
                </div>
                <div>
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
                    Teléfono
                  </label>
                  <input
                    value={form.telefono}
                    onChange={(e) => setForm((p) => ({ ...p, telefono: e.target.value }))}
                    style={inputStyle}
                  />
                </div>
                <div>
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
                    Email
                  </label>
                  <input
                    type="email"
                    value={form.email}
                    onChange={(e) => setForm((p) => ({ ...p, email: e.target.value }))}
                    style={inputStyle}
                  />
                </div>
              </div>
              <div
                style={{
                  display: 'flex',
                  gap: '12px',
                  marginTop: '32px',
                }}
              >
                <button
                  type="button"
                  onClick={() => setShowForm(false)}
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
                  type="submit"
                  disabled={guardando}
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
                    opacity: guardando ? 0.7 : 1,
                  }}
                  onMouseEnter={(e) => {
                    if (!guardando) e.currentTarget.style.backgroundColor = C.primaryDark;
                  }}
                  onMouseLeave={(e) => {
                    if (!guardando) e.currentTarget.style.backgroundColor = C.primary;
                  }}
                >
                  {guardando ? 'Guardando...' : 'Crear cliente'}
                </button>
              </div>
            </form>
          </div>
        </div>
      )}

      {/* Lista de clientes */}
      {loading ? (
        <div style={{ display: 'flex', flexDirection: 'column', gap: '12px' }}>
          {[...Array(5)].map((_, i) => (
            <div
              key={i}
              style={{
                height: '80px',
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
      ) : filtrados.length === 0 ? (
        <div
          style={{
            textAlign: 'center',
            padding: '48px 24px',
            backgroundColor: C.surface,
            borderRadius: '28px',
            border: `1px solid ${C.borderLight}`,
          }}
        >
          <UserCircle size={48} style={{ color: C.textSecondary, opacity: 0.5, marginBottom: '16px' }} />
          <p style={{ color: C.textSecondary, fontSize: '1rem', margin: 0 }}>
            {busqueda ? 'No se encontraron clientes con ese nombre' : 'Aún no hay clientes registrados'}
          </p>
          {!busqueda && isAdmin && (
            <button
              onClick={() => setShowForm(true)}
              style={{
                marginTop: '20px',
                backgroundColor: C.primary,
                color: '#fff',
                border: 'none',
                padding: '10px 20px',
                borderRadius: '40px',
                fontSize: '0.85rem',
                cursor: 'pointer',
              }}
            >
              + Agregar primer cliente
            </button>
          )}
        </div>
      ) : (
        <div style={{ display: 'flex', flexDirection: 'column', gap: '12px' }}>
          {filtrados.map((cliente) => (
            <div
              key={cliente.id}
              style={{
                backgroundColor: C.surface,
                borderRadius: '20px',
                border: `1px solid ${C.borderLight}`,
                padding: '16px 20px',
                display: 'flex',
                flexWrap: 'wrap',
                justifyContent: 'space-between',
                alignItems: 'center',
                transition: 'all 0.2s ease',
                cursor: 'pointer',
              }}
              onMouseEnter={(e) => {
                e.currentTarget.style.borderColor = C.primary;
                e.currentTarget.style.boxShadow = '0 8px 20px -8px rgba(0,0,0,0.1)';
                e.currentTarget.style.transform = 'translateY(-2px)';
              }}
              onMouseLeave={(e) => {
                e.currentTarget.style.borderColor = C.borderLight;
                e.currentTarget.style.boxShadow = 'none';
                e.currentTarget.style.transform = 'translateY(0)';
              }}
            >
              <div style={{ display: 'flex', alignItems: 'center', gap: '14px', flex: '2', minWidth: '180px' }}>
                <div
                  style={{
                    width: '48px',
                    height: '48px',
                    borderRadius: '24px',
                    display: 'flex',
                    alignItems: 'center',
                    justifyContent: 'center',
                    backgroundColor:
                      cliente.tipo === 'empresa' ? `${C.secondary}15` : `${C.primary}15`,
                  }}
                >
                  {cliente.tipo === 'empresa' ? (
                    <Building2 size={22} style={{ color: C.secondary }} />
                  ) : (
                    <User size={22} style={{ color: C.primary }} />
                  )}
                </div>
                <div>
                  <p
                    style={{
                      fontWeight: '600',
                      color: C.textPrimary,
                      margin: 0,
                      fontSize: '1rem',
                    }}
                  >
                    {cliente.nombre} {cliente.apellido}
                  </p>
                  <div style={{ display: 'flex', gap: '12px', marginTop: '4px', flexWrap: 'wrap' }}>
                    {cliente.telefono && (
                      <span style={{ fontSize: '0.7rem', color: C.textSecondary, display: 'flex', alignItems: 'center', gap: '4px' }}>
                        <Phone size={12} /> {cliente.telefono}
                      </span>
                    )}
                    {cliente.email && (
                      <span style={{ fontSize: '0.7rem', color: C.textSecondary, display: 'flex', alignItems: 'center', gap: '4px' }}>
                        <Mail size={12} /> {cliente.email}
                      </span>
                    )}
                    {!cliente.telefono && !cliente.email && (
                      <span style={{ fontSize: '0.7rem', color: C.textSecondary }}>Sin contacto</span>
                    )}
                  </div>
                </div>
              </div>
              <div style={{ display: 'flex', alignItems: 'center', gap: '16px', marginTop: '8px' }}>
                {cliente.deuda_total > 0 && (
                  <span
                    style={{
                      backgroundColor: C.alertBg,
                      color: C.alert,
                      fontWeight: '700',
                      fontSize: '0.75rem',
                      padding: '6px 12px',
                      borderRadius: '40px',
                      border: `1px solid ${C.alert}20`,
                    }}
                  >
                    Deuda: S/ {cliente.deuda_total.toFixed(2)}
                  </span>
                )}
                <span
                  style={{
                    fontSize: '0.7rem',
                    fontWeight: '500',
                    color: C.textSecondary,
                    backgroundColor: C.borderLight,
                    padding: '4px 10px',
                    borderRadius: '40px',
                  }}
                >
                  {cliente.tipo === 'empresa' ? 'Empresa' : 'Individual'}
                </span>
              </div>
            </div>
          ))}
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
        @media (max-width: 640px) {
          .shimmer::after {
            animation-duration: 1.5s;
          }
        }
      `}</style>
    </div>
  );
}

// Estilos reutilizables para inputs
const inputStyle: React.CSSProperties = {
  width: '100%',
  padding: '10px 14px',
  border: `1px solid ${C.borderLight}`,
  borderRadius: '16px',
  fontSize: '0.9rem',
  backgroundColor: C.surface,
  color: C.textPrimary,
  outline: 'none',
  transition: 'all 0.2s',
};